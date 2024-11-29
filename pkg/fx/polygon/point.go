package polygon

type Edge struct {
	Start Point
	End   Point
}

type Point struct {
	Latitude  float64 // Latitude (Y-axis)
	Longitude float64 // Longitude (X-axis)
}

// Point3D represents a point in 3D space.
type Point3D struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
}

// Edge3D represents an edge in 3D space.
type Edge3D struct {
	Start Point3D
	End   Point3D
}
