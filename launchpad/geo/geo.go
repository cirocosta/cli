package geo

type Point [2]float64

type Line struct {
	Type        string  `json:"type"`
	Coordinates []Point `json:"coordinates"`
}

type BoundingBox struct {
	Type        string  `json:"type"`
	Coordinates []Point `json:"coordinates"`
}

type Circle struct {
	Type        string `json:"type"`
	Coordinates Point  `json:"coordinates"`
	Radius      string `json:"radius"`
}

type Polygon struct {
	Type        string    `json:"type"`
	Coordinates [][]Point `json:"coordinates"`
}

func (p *Polygon) AddHole(coordinates ...Point) {
	p.Coordinates = append(p.Coordinates, coordinates)
}

func NewPoint(lat, lon float64) Point {
	var p Point

	p[0] = lat
	p[1] = lon

	return p
}

func NewLine(coordinates ...Point) Line {
	return Line{
		Type:        "linestring",
		Coordinates: coordinates,
	}
}

func NewBoundingBox(upperLeft, lowerRight Point) BoundingBox {
	return BoundingBox{
		Type:        "envelope",
		Coordinates: []Point{upperLeft, lowerRight},
	}
}

func NewCircle(coordinates Point, radius string) Circle {
	return Circle{
		Type:        "circle",
		Coordinates: coordinates,
		Radius:      radius,
	}
}

func NewPolygon(coordinates ...Point) Polygon {
	var x = [][]Point{coordinates}

	return Polygon{
		Type:        "polygon",
		Coordinates: x,
	}
}
