package server

import (
	"encoding/gob"

	"github.com/kpawlik/geojson"
	"github.com/kpawlik/smallgopher/config"
	"github.com/kpawlik/smallgopher/worker"
)

func init() {
	geojson.Register()
	gob.Register(worker.NewAcpErr(""))
	gob.Register(&config.FeaturesResponse{})
	gob.Register(make(map[string]interface{}))
	gob.Register(make([]interface{}, 0))
	gob.Register(config.NewFeatureID(""))
	gob.Register(config.NewFeatureID(uint64(12)))
}
