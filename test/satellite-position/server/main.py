from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import redis
import json

# FastAPI app instance
app = FastAPI()

# Redis client setup
redis_client = redis.StrictRedis(host="redis-service", port=6379, decode_responses=True)

# Mocked satellite data
mocked_satellites = [
    {"id": "1", "name": "Satellite 1", "latitude": 40.7128, "longitude": -74.0060},
    {"id": "2", "name": "Satellite 2", "latitude": 34.0522, "longitude": -118.2437},
    {"id": "3", "name": "Satellite 3", "latitude": 51.5074, "longitude": -0.1278},
]

# Data model for input validation
class SatelliteUpdate(BaseModel):
    id: str
    latitude: float
    longitude: float


@app.get("/")
def read_root():
    return {"message": "Satellite Position Service"}


@app.post("/propagate")
def propagate_positions():
    """
    Propagate all satellite positions to the Redis message broker.
    """
    try:
        for satellite in mocked_satellites:
            redis_client.publish("satellite_positions", json.dumps(satellite))
        return {"message": "Satellite positions propagated to Redis successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to propagate positions: {e}")


@app.post("/update")
def update_position(update: SatelliteUpdate):
    """
    Update a satellite's position and propagate it to the Redis broker.
    """
    # Find the satellite to update
    satellite = next((sat for sat in mocked_satellites if sat["id"] == update.id), None)
    if not satellite:
        raise HTTPException(status_code=404, detail="Satellite not found.")

    # Update satellite position
    satellite["latitude"] = update.latitude
    satellite["longitude"] = update.longitude

    # Publish the updated position to Redis
    try:
        redis_client.publish("satellite_positions", json.dumps(satellite))
        return {"message": f"Satellite {update.id} position updated and propagated."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to propagate updated position: {e}")
