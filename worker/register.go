package worker

import (
	"encoding/gob"

	"github.com/kpawlik/geojson"
	"github.com/kpawlik/smallgopher/config"
)

func init() {
	geojson.Register()
	gob.Register(NewAcpErr(""))
	gob.Register(&config.FeaturesResponse{})
	gob.Register(make(map[string]interface{}))
	gob.Register(make([]interface{}, 0))

}
