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

# Function to calculate satellite position based on TLE using Skyfield
def propagate_satellite_position(
    satellite_id: str,
    tle_line1: str,
    tle_line2: str,
    start_time: str,
    duration_minutes: int = 1440,
    interval_seconds: int = 60
) -> List[Dict]:
    try:
        start_time = datetime.fromisoformat(start_time.replace("Z", "+00:00"))

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

# Background task to propagate satellite position to Redis
def publish_satellite_positions(satellite_id: str, positions: List[Dict]):
    try:
        for pos in positions:
            # Use a structured Redis key for querying (e.g., satellite_positions:<id>:<timestamp>)
            key = f"satellite_positions:{satellite_id}:{pos['timestamp']}"
            redis_client.set(key, json.dumps(pos))
            redis_client.publish("satellite_positions", json.dumps(pos))
            logger.info(f"Published position for {satellite_id} at {pos['timestamp']}")
    except Exception as e:
        logger.error(f"Error publishing satellite position to Redis: {e}")

# Background task to handle Redis subscription for incoming TLE updates
def subscribe_to_tle_updates():
    """
    Subscribe to Redis channel for TLE updates and process the data.
    """
    pubsub = redis_client.pubsub()
    pubsub.subscribe("satellite_tle_updates")  # Subscribe to the 'satellite_tle_updates' channel

    logger.info("Subscribed to Redis channel: satellite_tle_updates")

    for message in pubsub.listen():
        if message["type"] == "message":
            try:
                # Log the raw message received
                logger.debug(f"Raw message received: {message}")

                # Parse the TLE update message
                satellite_data = json.loads(message["data"])
                logger.debug(f"Parsed satellite data: {json.dumps(satellite_data, indent=2)}")

                tle_line1 = satellite_data.get("line_1")
                tle_line2 = satellite_data.get("line_2")
                satellite_id = satellite_data.get("id")
                epoch = satellite_data.get("epoch", datetime.utcnow().isoformat() + "Z")

                if not tle_line1 or not tle_line2 or not satellite_id:
                    logger.warning(f"Incomplete TLE data received: {json.dumps(satellite_data, indent=2)}")
                    continue

                # Trigger satellite propagation
                logger.info(f"Starting propagation for satellite {satellite_id}")
                start_time = epoch
                duration_minutes = 90  # Default propagation duration
                interval_seconds = 15  # Default propagation interval

                # Propagate the satellite positions
                positions = propagate_satellite_position(
                    satellite_id,
                    tle_line1,
                    tle_line2,
                    start_time,
                    duration_minutes,
                    interval_seconds
                )

                # Publish the propagated positions to Redis
                publish_satellite_positions(satellite_id, positions)

                logger.info(f"Finished propagation and publishing for satellite {satellite_id}")

            except Exception as e:
                logger.error(f"Error processing TLE message: {e}")



# Start the Redis subscription in a separate thread when the application starts
@app.on_event("startup")
def start_tle_subscription():
    # Start subscribing to Redis in a separate thread
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
    """
    Propagate satellite position based on TLE data for a specified duration and interval.
    """
    data = await request.json()
    tle_line1 = data.get("line_1")
    tle_line2 = data.get("line_2")
    start_time = data.get("start_time")
    duration_minutes = data.get("duration_minutes", 90)
    interval_seconds = data.get("interval_seconds", 15)

    if not tle_line1 or not tle_line2 or not start_time:
        raise HTTPException(status_code=400, detail="TLE data and start time are required")

    satellite_id = f"satellite-{hash(tle_line1 + tle_line2)}"

    # Propagate satellite position
    positions = propagate_satellite_position(
        satellite_id,
        tle_line1,
        tle_line2,
        start_time,
        duration_minutes,
        interval_seconds
    )

    # Start a background task to publish positions to Redis
    background_tasks.add_task(publish_satellite_positions, satellite_id, positions)

    return {"message": "Satellite propagation started", "positions_count": len(positions)}
