from flask import Flask, request, jsonify
from skyfield.api import EarthSatellite, load, wgs84
from datetime import datetime, timedelta

app = Flask(__name__)

@app.route('/satellite/propagate', methods=['POST'])
def propagate_satellite():
    try:
        data = request.json
        tle_line1 = data['tle_line1']
        tle_line2 = data['tle_line2']
        start_time = datetime.fromisoformat(data['start_time'])
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

    except Exception as e:
        return jsonify({"error": str(e)}), 400

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
