from fastapi import FastAPI, HTTPException, BackgroundTasks
from pydantic import BaseModel
import redis
import json
from skyfield.api import EarthSatellite, load, wgs84
from datetime import datetime, timedelta
from typing import List
import threading
import logging

# FastAPI app instance
app = FastAPI()

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Redis client setup
redis_client = redis.StrictRedis(host="redis-service", port=6379, decode_responses=True)

# Data model for input validation
class TLEData(BaseModel):
    tle_line1: str
    tle_line2: str
    start_time: str
    duration_minutes: int = 90
    interval_seconds: int = 15

class Position(BaseModel):
    time: str
    latitude: float
    longitude: float
    altitude: float

# Function to calculate satellite position based on TLE using Skyfield
def propagate_satellite_position(tle_line1: str, tle_line2: str, start_time: str, duration_minutes: int = 1440, interval_seconds: int = 60) -> List[Position]:
    try:
        # Preprocess start_time to handle "Z" in ISO 8601
        start_time = datetime.fromisoformat(start_time.replace("Z", "+00:00"))

        ts = load.timescale()
        satellite = EarthSatellite(tle_line1, tle_line2, "Satellite", ts)

        end_time = start_time + timedelta(minutes=duration_minutes)
        current_time = start_time
        positions = []

        while current_time <= end_time:
            t = ts.utc(current_time.year, current_time.month, current_time.day,
                       current_time.hour, current_time.minute, current_time.second)
            geocentric = satellite.at(t)
            subpoint = wgs84.subpoint(geocentric)
            positions.append(Position(
                time=current_time.isoformat(),
                latitude=subpoint.latitude.degrees,
                longitude=subpoint.longitude.degrees,
                altitude=subpoint.elevation.km
            ))
            current_time += timedelta(seconds=interval_seconds)

        return positions

    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Error in propagating satellite position: {e}")

# Background task to propagate satellite position to Redis
def publish_satellite_positions(positions: List[Position]):
    try:
        for pos in positions:
            # Publish each position update to the "satellite_positions" Redis channel
            redis_client.publish("satellite_positions", json.dumps(pos.dict()))
            
            # Log the publication of the satellite position
            logger.info(f"Published position: {pos.time} - Latitude: {pos.latitude}, Longitude: {pos.longitude}, Altitude: {pos.altitude} km")
    except Exception as e:
        logger.error(f"Error publishing satellite position to Redis: {e}")

# Background task to handle Redis subscription for incoming satellite updates
def subscribe_to_redis():
    """
    Subscribe to Redis channel and process incoming satellite updates.
    When a new satellite update with TLE is received, start the propagation.
    """
    pubsub = redis_client.pubsub()
    pubsub.subscribe("satellite_tle_data")  # Subscribe to the TLE channel

    for message in pubsub.listen():
        if message["type"] == "message":
            # Decode the message (satellite TLE data)
            satellite_data = json.loads(message["data"])

            # Check if the message contains valid TLE data
            if "tle_line1" in satellite_data and "tle_line2" in satellite_data:
                tle_line1 = satellite_data["tle_line1"]
                tle_line2 = satellite_data["tle_line2"]
                start_time = satellite_data["start_time"]

                # Propagate the satellite position based on the received TLE
                positions = propagate_satellite_position(tle_line1, tle_line2, start_time, 1440, 60)  # Default: 24 hours, 1 min interval

                # Start a background task to publish the propagated positions to Redis
                publish_satellite_positions(positions)

# Root endpoint for checking service status
@app.get("/")
def read_root():
    return {"message": "Satellite Propagation Service is running"}

# Start the Redis subscription in a separate thread when the application starts
@app.on_event("startup")
def start_subscribe_to_redis():
    # Automatically propagate satellite positions for 24 hours from now, with 1 minute intervals
    current_time = datetime.utcnow().isoformat() + "Z"
    tle_line1 = "1 25544U 98067A   21349.25024444  .00001234  00000-0  12345-4 0  9990"  # Example TLE line
    tle_line2 = "2 25544  51.6404 102.7490 0007637  35.6700 324.1674 15.48906415594989"  # Example TLE line
    positions = propagate_satellite_position(tle_line1, tle_line2, current_time, 1440, 60)  # Propagate for 24 hours

    # Start a background task to publish positions to Redis
    threading.Thread(target=publish_satellite_positions, args=(positions,), daemon=True).start()

    # Start subscribing to Redis in a separate thread
    threading.Thread(target=subscribe_to_redis, daemon=True).start()

# Endpoint to propagate satellite position based on TLE
@app.post("/satellite/propagate")
async def propagate_satellite(tle_data: TLEData, background_tasks: BackgroundTasks):
    """
    Propagate satellite position based on TLE data for a specified duration and interval.
    """
    # Propagate satellite position for the given duration
    positions = propagate_satellite_position(
        tle_data.tle_line1,
        tle_data.tle_line2,
        tle_data.start_time,
        tle_data.duration_minutes,
        tle_data.interval_seconds
    )

    # Start a background task to publish positions to Redis
    background_tasks.add_task(publish_satellite_positions, positions)

    return {"message": "Satellite propagation started", "positions_count": len(positions)}
