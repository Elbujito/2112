package polygon

import "fmt"

type Quadkey struct {
	Lat   float64
	Long  float64
	Level int
}

func NewQuadkey(lat float64, long float64, level int) Quadkey {
	return Quadkey{
		Lat:   lat,
		Long:  long,
		Level: level,
	}
}

func (q *Quadkey) Key() string {
	return fmt.Sprintf("%d-%f-%f", q.Level, q.Lat, q.Long)
}
