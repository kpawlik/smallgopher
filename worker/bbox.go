package worker

import (
	"math"

	"github.com/kpawlik/geojson"
)

// BBox struct to support Bounding box operations
type BBox struct {
	x, y          float64
	width, height float64
}

//NewBBoxFromCoordinates creates new BBox from coordinates
func NewBBoxFromCoordinates(coords geojson.Coordinates) *BBox {
	var (
		fx, fy                 float64
		minX, minY, maxX, maxY float64
		x, y                   float64
	)
	min, max := math.Min, math.Max
	for _, coord := range coords {
		fx, fy = float64(coord[0]), float64(coord[1])
		x, y = math.Abs(fx), math.Abs(fy)
		minX = min(minX, x)
		maxX = max(maxX, x)
		minY = min(minY, y)
		maxY = max(maxY, y)
	}

	return &BBox{x: minX, y: minY, width: math.Abs(maxX - minX), height: math.Abs(maxY - minY)}
}

//NewBBox creates new BBox from N, W, S, E data
func NewBBox(n, w, s, e float64) *BBox {
	x := math.Min(math.Abs(n), math.Abs(s))
	width := math.Abs(n - s)
	y := math.Min(math.Abs(w), math.Abs(e))
	height := math.Abs(w - e)
	return &BBox{x: x, y: y, width: width, height: height}
}

//NewBBoxFromCoordinate creates new BBox from coordinate
func NewBBoxFromCoordinate(coord geojson.Coordinate) *BBox {
	x, y := math.Abs(float64(coord[0])), math.Abs(float64(coord[1]))
	return &BBox{
		x:      x,
		y:      y,
		width:  0,
		height: 0,
	}
}

//Interact return true if b and bbox are interacts
func (b *BBox) Interact(bbox *BBox) bool {
	return !((bbox.x > b.x+b.width) || (bbox.x+bbox.width < b.x) || (bbox.y+bbox.height < b.y) || (bbox.y > b.y+b.height))
}
