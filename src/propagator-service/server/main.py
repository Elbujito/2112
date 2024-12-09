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
    id: str
    timestamp: str
    latitude: float
    longitude: float
    altitude: float

# Function to calculate satellite position based on TLE using Skyfield
def propagate_satellite_position(
    satellite_id: str,
    tle_line1: str,
    tle_line2: str,
    start_time: str,
    duration_minutes: int = 1440,
    interval_seconds: int = 60
) -> List[Position]:
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
            positions.append(Position(
                id=satellite_id,
                timestamp=current_time.isoformat(),
                latitude=subpoint.latitude.degrees,
                longitude=subpoint.longitude.degrees,
                altitude=subpoint.elevation.km
            ))
            current_time += timedelta(seconds=interval_seconds)

        return positions

    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Error in propagating satellite position: {e}")

# Background task to propagate satellite position to Redis
def publish_satellite_positions(satellite_id: str, positions: List[Position]):
    try:
        for pos in positions:
            # Use a structured Redis key for querying (e.g., satellite_positions:<id>:<timestamp>)
            key = f"satellite_positions:{satellite_id}:{pos.timestamp}"
            redis_client.set(key, json.dumps(pos.dict()))
            redis_client.publish("satellite_positions", json.dumps(pos.dict()))
            logger.info(f"Published position for {satellite_id} at {pos.timestamp}")
    except Exception as e:
        logger.error(f"Error publishing satellite position to Redis: {e}")

# Background task to handle Redis subscription for incoming satellite updates
def subscribe_to_redis():
    """
    Subscribe to Redis channel and process incoming satellite updates.
    """
    pubsub = redis_client.pubsub()
    pubsub.subscribe("satellite_tle_data")

    for message in pubsub.listen():
        if message["type"] == "message":
            try:
                satellite_data = json.loads(message["data"])
                if "tle_line1" in satellite_data and "tle_line2" in satellite_data and "id" in satellite_data:
                    tle_line1 = satellite_data["tle_line1"]
                    tle_line2 = satellite_data["tle_line2"]
                    satellite_id = satellite_data["id"]
                    start_time = satellite_data.get("start_time", datetime.utcnow().isoformat() + "Z")

                    # Propagate the satellite position
                    positions = propagate_satellite_position(satellite_id, tle_line1, tle_line2, start_time)

                    # Publish positions to Redis
                    publish_satellite_positions(satellite_id, positions)
            except Exception as e:
                logger.error(f"Error processing Redis message: {e}")

# Root endpoint for checking service status
@app.get("/")
def read_root():
    return {"message": "Satellite Propagation Service is running"}

# Start the Redis subscription in a separate thread when the application starts
@app.on_event("startup")
def start_subscribe_to_redis():
    # Start subscribing to Redis in a separate thread
    threading.Thread(target=subscribe_to_redis, daemon=True).start()

@app.route('/satellite/legacy/propagate', methods=['POST'])
def propagate_satellite():
    try:
        data = request.json
        tle_line1 = data['tle_line1']
        tle_line2 = data['tle_line2']
        
        # Preprocess start_time to handle "Z" in ISO 8601
        start_time = datetime.fromisoformat(data['start_time'].replace("Z", "+00:00"))
        
        duration_minutes = data.get('duration_minutes', 90)
        interval_seconds = data.get('interval_seconds', 15)

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
            positions.append({
                "time": current_time.isoformat(),
                "latitude": subpoint.latitude.degrees,
                "longitude": subpoint.longitude.degrees,
                "altitude": subpoint.elevation.km
            })
            current_time += timedelta(seconds=interval_seconds)

        return jsonify(positions)

# Endpoint to propagate satellite position based on TLE
@app.post("/satellite/propagate")
async def propagate_satellite(tle_data: TLEData, background_tasks: BackgroundTasks):
    """
    Propagate satellite position based on TLE data for a specified duration and interval.
    """
    satellite_id = f"satellite-{hash(tle_data.tle_line1 + tle_data.tle_line2)}"

    # Propagate satellite position
    positions = propagate_satellite_position(
        satellite_id,
        tle_data.tle_line1,
        tle_data.tle_line2,
        tle_data.start_time,
        tle_data.duration_minutes,
        tle_data.interval_seconds
    )

    # Start a background task to publish positions to Redis
    background_tasks.add_task(publish_satellite_positions, satellite_id, positions)

    return {"message": "Satellite propagation started", "positions_count": len(positions)}
