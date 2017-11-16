package worker

import (
	"fmt"

	"github.com/kpawlik/geojson"
)

//AcpErr ACP error
type AcpErr struct {
	Err string
}

//Error implements error interface
func (err *AcpErr) Error() string {
	return err.Err
}

// NewAcpErr new error
func NewAcpErr(msg string) *AcpErr {
	return &AcpErr{msg}
}

// NewAcpErrf new error from formated string
func NewAcpErrf(format string, args ...interface{}) *AcpErr {
	return &AcpErr{fmt.Sprintf(format, args...)}
}

func toCoordinate(arr [2]float64) geojson.Coordinate {
	return geojson.Coordinate{
		geojson.Coord(arr[0]),
		geojson.Coord(arr[1]),
	}
}

func toCoordinates(arr ...[2]float64) geojson.Coordinates {
	return geojson.Coordinates{geojson.Coordinate{
		geojson.Coord(arr[0]),
		geojson.Coord(arr[1]),
	},
	}
}
