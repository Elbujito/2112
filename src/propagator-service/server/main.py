from fastapi import FastAPI, HTTPException, BackgroundTasks, Request
import redis
import json
from skyfield.api import EarthSatellite, load, wgs84
from datetime import datetime, timedelta
from typing import List, Dict
import threading
import logging

# FastAPI app instance
app = FastAPI()

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Redis client setup
redis_client = redis.StrictRedis(host="redis-service", port=6379, decode_responses=True)

# Function to propagate satellite positions
def propagate_satellite_position(
    satellite_id: str,
    tle_line1: str,
    tle_line2: str,
    start_time: str,
    duration_minutes: int,
    interval_seconds: int
) -> List[Dict]:
    try:
        # Preprocess the timestamp
        start_time = start_time.split("+")[0]
        if "." in start_time:
            start_time = start_time.split(".")[0] + "." + start_time.split(".")[1][:6] + "+00:00"
        else:
            start_time += "+00:00"

        start_time = datetime.fromisoformat(start_time)

        ts = load.timescale()
        satellite = EarthSatellite(tle_line1, tle_line2, satellite_id, ts)

        end_time = start_time + timedelta(minutes=duration_minutes)
        current_time = start_time
        positions = []

        while current_time <= end_time:
            t = ts.utc(current_time.year, current_time.month, current_time.day,
                       current_time.hour, current_time.minute, current_time.second)
            geocentric = satellite.at(t)
            subpoint = wgs84.subpoint(geocentric)
            positions.append({
                "id": satellite_id,
                "timestamp": current_time.isoformat(),
                "latitude": subpoint.latitude.degrees,
                "longitude": subpoint.longitude.degrees,
                "altitude": subpoint.elevation.km
            })
            current_time += timedelta(seconds=interval_seconds)

        return positions

    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Error in propagating satellite position: {e}")

# Function to publish satellite positions to Redis
def publish_satellite_positions(satellite_id: str, positions: List[Dict]):
    try:
        for pos in positions:
            key = f"satellite_positions:{satellite_id}"
            timestamp = datetime.fromisoformat(pos['timestamp']).timestamp()
            redis_client.zadd(key, {json.dumps(pos): timestamp})
            redis_client.publish("satellite_positions", json.dumps(pos))
            logger.info(f"Published position for {satellite_id} at {pos['timestamp']}")

        update_event = {
            "event": "event_satellite_positions_updated",
            "satellite_id": satellite_id,
            "start_time": positions[0]["timestamp"] if positions else None,
            "end_time": positions[-1]["timestamp"] if positions else None,
            "positions_count": len(positions),
        }

        redis_client.publish("event_satellite_positions_updated", json.dumps(update_event))
        logger.info(f"Published satellite update event for {satellite_id}")

    except Exception as e:
        logger.error(f"Error publishing satellite position to Redis: {e}")

# Function to subscribe to TLE updates and compute one pass
def subscribe_to_tle_updates():
    pubsub = redis_client.pubsub()
    pubsub.subscribe("satellite_tle_updates")

    logger.info("Subscribed to Redis channel: satellite_tle_updates")

    for message in pubsub.listen():
        if message["type"] == "message":
            try:
                satellite_data = json.loads(message["data"])

                tle_line1 = satellite_data.get("line_1")
                tle_line2 = satellite_data.get("line_2")
                satellite_id = satellite_data.get("id")
                epoch = satellite_data.get("epoch", datetime.utcnow().isoformat() + "Z")

                # Validate TLE data
                if not tle_line1 or not tle_line2 or not satellite_id:
                    logger.warning(f"Incomplete TLE data received: {satellite_data}")
                    continue

                # Calculate altitude and orbital period
                ts = load.timescale()
                satellite = EarthSatellite(tle_line1, tle_line2, satellite_id, ts)
                geocentric = satellite.at(ts.now())
                subpoint = wgs84.subpoint(geocentric)
                altitude_km = subpoint.elevation.km

                if altitude_km < 2000:  # LEO
                    pass_duration_minutes = 1440  # 24 hours
                    num_points = 1440  # 1 point per minute
                elif altitude_km < 35786:  # MEO
                    pass_duration_minutes = 180  # 3 hours
                    num_points = 50  # Approx. 1 point every 3.6 minutes
                else:  # GEO
                    pass_duration_minutes = 1440  # 24 hours
                    num_points = 10  # 1 point every 144 minutes

                interval_seconds = (pass_duration_minutes * 60) // num_points

                logger.info(f"Starting propagation for satellite {satellite_id}")
                positions = propagate_satellite_position(
                    satellite_id, tle_line1, tle_line2, epoch,
                    pass_duration_minutes, interval_seconds
                )
                publish_satellite_positions(satellite_id, positions)
                logger.info(f"Finished propagation and publishing for satellite {satellite_id}")

            except Exception as e:
                logger.error(f"Error processing TLE message: {e}")

# Start the Redis subscription in a separate thread when the application starts
@app.on_event("startup")
def start_tle_subscription():
    threading.Thread(target=subscribe_to_tle_updates, daemon=True).start()

# Root endpoint for checking service status
@app.get("/")
def read_root():
    return {"message": "Satellite Propagation Service is running"}

# Health check endpoint
@app.get("/health")
def health_check():
    try:
        redis_client.ping()
        return {"status": "healthy"}
    except Exception as e:
        raise HTTPException(status_code=500, detail="Redis connection failed")

# Endpoint to propagate satellite position based on TLE
@app.post("/satellite/propagate")
async def propagate_satellite(request: Request, background_tasks: BackgroundTasks):
    data = await request.json()
    tle_line1 = data.get("line_1")
    tle_line2 = data.get("line_2")
    start_time = data.get("start_time")
    duration_minutes = data.get("duration_minutes", 90)
    interval_seconds = data.get("interval_seconds", 15)

    if not tle_line1 or not tle_line2 or not start_time:
        raise HTTPException(status_code=400, detail="TLE data and start time are required")

    satellite_id = f"satellite-{hash(tle_line1 + tle_line2)}"

    positions = propagate_satellite_position(
        satellite_id, tle_line1, tle_line2, start_time, duration_minutes, interval_seconds
    )

    background_tasks.add_task(publish_satellite_positions, satellite_id, positions)

    return {"message": "Satellite propagation started", "positions_count": len(positions)}
