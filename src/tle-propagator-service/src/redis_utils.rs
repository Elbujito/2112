use redis:: { Commands, RedisResult };
use serde_json:: Value;
use tracing::{debug, error};
use redis::RedisError;

pub fn publish_to_redis(channel: & str, message: & Value, redis_client: & redis:: Client) -> RedisResult < () > {
    let mut con = redis_client.get_connection() ?;
    let _: () = con.publish(channel, message.to_string()) ?;
    Ok(())
}

pub fn store_to_redis(
    key: &str,
    data: &Value,
    score: i64, // Changed to integer
    redis_client: &redis::Client,
) -> RedisResult<()> {
    // Validate inputs
    if key.is_empty() {
        error!("Empty Redis key provided, aborting operation.");
        return Err(RedisError::from((
            redis::ErrorKind::InvalidClientConfig,
            "Empty Redis key",
        )));
    }

    // Serialize data
    let serialized_data = serde_json::to_string(data).map_err(|e| {
        error!("Serialization error for key '{}': {}", key, e);
        RedisError::from((
            redis::ErrorKind::InvalidClientConfig,
            "Serialization error", // Changed to static str
        ))
    })?;

    // Connect to Redis
    let mut con = redis_client.get_connection().map_err(|e| {
        error!("Failed to connect to Redis for key '{}': {}", key, e);
        RedisError::from((
            redis::ErrorKind::IoError,
            "Failed to connect to Redis", // Changed to static str
        ))
    })?;

    // Execute ZADD
    redis::cmd("ZADD")
        .arg(key)
        .arg(score)
        .arg(&serialized_data)
        .query(&mut con)
        .map_err(|e| {
            error!(
                "Failed to execute ZADD for key '{}', score {}: {}",
                key, score, e
            );
            RedisError::from((
                redis::ErrorKind::IoError,
                "Failed to execute ZADD", // Changed to static str
            ))
        })?;

    debug!(
        "Successfully added to Redis: key='{}', score={}, data={}",
        key, score, serialized_data
    );

    Ok(())
}
