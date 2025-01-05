use std::sync::Arc;
use std::time::{SystemTime, UNIX_EPOCH, Duration};
use redis::Client as RedisClient;
use serde_json::json;
use tokio::sync::mpsc;
use tracing::{error, info, debug};
use chrono::{DateTime, Utc};
use sgp4::{Elements, Constants, MinutesSinceEpoch};
use crate::tle::{TLEUpdate, SatellitePosition};
use crate::redis_utils::{publish_to_redis, store_to_redis};

/// Parse eccentricity from the raw TLE string and interpret it correctly.
fn parse_eccentricity(tle_line2: &str) -> f64 {
    tle_line2[26..33]
        .trim()
        .parse::<f64>()
        .unwrap_or(0.0)
        / 10_000_000.0 // Interpret as 0.xxxxxxxx
}

/// Extract inclination from TLE second line as a float.
fn parse_inclination(tle_line2: &str) -> f64 {
    tle_line2[8..16].trim().parse::<f64>().unwrap_or(0.0)
}

/// Convert TLE epoch to UNIX time in seconds.
fn tle_epoch_to_unix(epoch: f64) -> f64 {
    let year = if epoch >= 1_000.0 {
        2000 + ((epoch / 1_000.0).floor() as i32)
    } else {
        1900 + ((epoch / 1_000.0).floor() as i32)
    };
    let day_of_year = epoch % 1_000.0;

    let jan_1 = chrono::NaiveDate::from_ymd_opt(year, 1, 1)
        .expect("Failed to create NaiveDate for January 1st");

    jan_1
        .and_hms_opt(0, 0, 0)
        .expect("Failed to create NaiveDateTime at midnight")
        .and_utc()
        .timestamp() as f64
        + (day_of_year - 1.0) * 86400.0
}

/// Convert SystemTime to a readable ISO 8601 timestamp.
fn system_time_to_iso8601(system_time: SystemTime) -> String {
    let datetime: DateTime<Utc> = system_time.into();
    datetime.to_rfc3339()
}

/// Calculate latitude, longitude, and altitude from Cartesian coordinates.
fn calculate_lat_lon_alt(position: [f64; 3]) -> (f64, f64, f64) {
    let [x, y, z] = position;
    let r = (x * x + y * y + z * z).sqrt();
    let latitude = (z / r).asin().to_degrees();
    let longitude = y.atan2(x).to_degrees();
    let altitude = r - 6371.0; // Earth's radius in km
    (latitude, longitude, altitude)
}

/// Propagate satellite positions using raw TLE strings.
pub async fn propagate_satellite_positions(
    satellite_id: &str,
    tle_line1: &str,
    tle_line2: &str,
    start_time: SystemTime,
    duration_minutes: u64,
    interval_seconds: u64,
) -> Vec<SatellitePosition> {
    let mut positions = Vec::new();
    let mut failures = Vec::new();

    debug!(
        "Processing TLE for satellite {}: line1={}, line2={}",
        satellite_id, tle_line1, tle_line2
    );

    // Correct TLE line2 if necessary
    let eccentricity = parse_eccentricity(tle_line2);
    let inclination = parse_inclination(tle_line2);

    debug!(
        "Parsed TLE for satellite {}: eccentricity={}, inclination={}°",
        satellite_id, eccentricity, inclination
    );

    // Parse the TLE data
    let elements = match sgp4::Elements::from_tle(
        Some(satellite_id.to_owned()),
        tle_line1.as_bytes(),
        tle_line2.as_bytes(),
    ) {
        Ok(e) => e,
        Err(e) => {
            error!("Invalid TLE for satellite {}: {:?}", satellite_id, e);
            return positions;
        }
    };

    // Compute constants for propagation
    let constants = match Constants::from_elements_afspc_compatibility_mode(&elements) {
        Ok(c) => c,
        Err(e) => {
            error!("Failed to create constants for satellite {}: {:?}", satellite_id, e);
            return positions;
        }
    };

    debug!(
        "Successfully completed TLE parsing and constants computation for satellite {}.",
        satellite_id
    );

    let tle_epoch = tle_epoch_to_unix(elements.epoch());
    let end_time = start_time + Duration::from_secs(duration_minutes * 60);
    let mut current_time = start_time;

    debug!(
        "Propagating satellite {} from {:?} to {:?} with interval {} seconds",
        satellite_id,
        system_time_to_iso8601(start_time),
        system_time_to_iso8601(end_time),
        interval_seconds
    );

    while current_time <= end_time {
        let timestamp = match current_time.duration_since(UNIX_EPOCH) {
            Ok(d) => d.as_secs_f64(),
            Err(e) => {
                error!("SystemTime error: {:?}", e);
                break;
            }
        };

        let minutes_since_epoch = MinutesSinceEpoch((timestamp - tle_epoch) / 60.0);

        match constants.propagate_afspc_compatibility_mode(minutes_since_epoch) {
            Ok(prediction) => {
                let (latitude, longitude, altitude) = calculate_lat_lon_alt(prediction.position);
                positions.push(SatellitePosition {
                    id: satellite_id.to_string(),
                    timestamp: system_time_to_iso8601(current_time),
                    latitude,
                    longitude,
                    altitude,
                });
            }
            Err(_) => {
                failures.push(current_time);
            }
        }

        current_time += Duration::from_secs(interval_seconds);
    }

    // Summary log for successful propagation
    if failures.is_empty() {
        info!(
            "Successfully propagated satellite {}: {} positions generated.",
            satellite_id,
            positions.len()
        );
    } else {
        // Log batch failures
        let total_failures = failures.len();
        let sample_failures: Vec<String> = failures
            .iter()
            .take(5)
            .map(|&time| system_time_to_iso8601(time))
            .collect();

        error!(
            "Propagation failed for satellite {}: {} failures out of {} intervals. Sample times: {:?}{}",
            satellite_id,
            total_failures,
            ((end_time.duration_since(start_time).unwrap().as_secs() / interval_seconds) as usize),
            sample_failures,
            if total_failures > 5 {
                format!(" ... (and {} more)", total_failures - 5)
            } else {
                String::new()
            }
        );

        // Return an empty positions list to indicate failure
        return Vec::new();
    }

    positions
}


/// Process a single TLE update, propagate positions, and store them in Redis.
pub async fn process_tle_update(redis_client: Arc<RedisClient>, tle_update: TLEUpdate) {
    info!("Received TLE update for satellite {}", tle_update.id);

    let start_time = SystemTime::now();
    let duration_minutes = 90;
    let interval_seconds = 15;

    let positions = propagate_satellite_positions(
        &tle_update.id,
        &tle_update.line_1,
        &tle_update.line_2,
        start_time,
        duration_minutes,
        interval_seconds,
    )
    .await;

    for position in &positions {
        let position_json = match serde_json::to_value(position) {
            Ok(json) => json,
            Err(e) => {
                error!("Failed to serialize position for satellite {}: {:?}", tle_update.id, e);
                continue;
            }
        };

        if let Err(e) = store_to_redis(
            &format!("satellite_positions:{}", tle_update.id),
            &position_json,
            position.timestamp.parse().unwrap_or_default(),
            &redis_client,
        ) {
            error!("Failed to store position in Redis for satellite {}: {:?}", tle_update.id, e);
        }
    }

    let summary = json!({
        "event": "event_satellite_positions_updated",
        "satellite_id": tle_update.id,
        "start_time": positions.first().map(|p| &p.timestamp),
        "end_time": positions.last().map(|p| &p.timestamp),
        "positions_count": positions.len(),
    });

    match publish_to_redis("event_satellite_positions_updated", &summary, &redis_client) {
        Ok(_) => {
            info!(
                "Successfully published update for satellite {} to Redis.",
                tle_update.id
            );
        }
        Err(e) => {
            error!("Failed to publish summary for satellite {}: {:?}", tle_update.id, e);
        }
    }
}

/// Subscribe to Redis channel for TLE updates and process them concurrently.
pub async fn subscribe_to_tle_updates(redis_client: Arc<RedisClient>, max_threads: usize) {
    let mut con = match redis_client.get_connection() {
        Ok(c) => c,
        Err(e) => {
            error!("Failed to connect to Redis: {:?}", e);
            return;
        }
    };

    let mut pubsub = con.as_pubsub();

    if let Err(e) = pubsub.subscribe("satellite_tle_updates") {
        error!("Failed to subscribe to channel: {:?}", e);
        return;
    }

    let (tx, mut rx) = mpsc::channel::<TLEUpdate>(max_threads);

    tokio::spawn(async move {
        while let Some(tle_update) = rx.recv().await {
            let redis_client = Arc::clone(&redis_client);
            tokio::spawn(async move {
                process_tle_update(redis_client, tle_update).await;
            });
        }
    });

    loop {
        if let Ok(message) = pubsub.get_message() {
            if let Ok(payload) = message.get_payload::<String>() {
                match serde_json::from_str::<TLEUpdate>(&payload) {
                    Ok(tle_update) => {
                        info!("Processing TLE update payload: {:?}", tle_update);
                        if tx.send(tle_update).await.is_err() {
                            error!("Failed to send TLE update to workers");
                        }
                    }
                    Err(e) => {
                        error!("Failed to deserialize TLE update: {:?}", e);
                    }
                }
            }
        }
    }
}
