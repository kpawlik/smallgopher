package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/kpawlik/geojson"
	"github.com/kpawlik/smallgopher/config"
)

// loadCacheFromFile loads data for collection with name from JSON file with the same name
func (t *Worker) loadCacheFromFile(name string, r config.Request) (*geojson.FeatureCollection, error) {
	var (
		buff       []byte
		ok         bool
		fc         *geojson.FeatureCollection
		err        error
		featureMap FeaturesMap
	)
	idFieldName := r.GetIDFieldName()
	fc = geojson.NewFeatureCollection(nil)
	if featureMap, ok = t.Cache[name]; !ok {
		featureMap = t.getCollectionCache(name)
		buff, err = ioutil.ReadFile(fmt.Sprintf("offline/%s.json", name))
		if err != nil {
			delete(t.Cache, name)
			return nil, err
		}
		err = json.Unmarshal(buff, &fc)
		if err != nil {
			return nil, err
		}
		for _, feature := range fc.Features {
			fid := feature.Properties[idFieldName].(float64)
			featureMap[int64(fid)] = feature
		}
	} else {
		for _, feature := range featureMap {
			fc.AddFeatures(feature)
		}
	}

	return fc, nil
}

//GetTestFeatures get features id and geom
func (t *Worker) GetTestFeatures(request *config.FeaturesRequest, resp *config.FeaturesResponse) error {
	var (
		err error
	)
	defer func() {
		if err := recover(); err != nil {
			log.Panic("PANIC ", err)
		}
	}()
	bbN := float64(request.BBox[0])
	bbE := float64(request.BBox[1])
	bbS := float64(request.BBox[2])
	bbW := float64(request.BBox[3])
	fc := geojson.NewFeatureCollection(nil)

	cachedFc, err := t.loadCacheFromFile(request.Collection, request)
	if err != nil {
		return err
	}
	features := filterFeatures(cachedFc.Features, bbN, bbE, bbS, bbW)
	fc.AddFeatures(features...)
	resp.SetBody(fc)
	return nil
}

//GetTestFeature get feature details
func (t *Worker) GetTestFeature(request *config.FeatureRequest, resp *config.FeaturesResponse) error {
	var (
		err     error
		cacheFC *geojson.FeatureCollection
		fid     float64
		id      int64
		ok      bool
	)
	defer func() {
		if err := recover(); err != nil {
			log.Panic("PANIC ", err)
		}
	}()
	IDName := request.IDFieldName
	IDValue := request.IDFieldValue
	id, ok = IDValue.(int64)
	if !ok {
		switch IDValue.(type) {
		case string:
			id, _ = strconv.ParseInt(IDValue.(string), 10, 64)
		case int:
			id = int64(IDValue.(int))
		case float64, float32:
			id = int64(IDValue.(float64))
		default:
			log.Panicf("Cannot parse %v to int64", IDValue)
		}
	}

	name := request.Collection
	fc := geojson.NewFeatureCollection(nil)
	cacheFC, err = t.loadCacheFromFile(name, request)
	if err != nil {
		log.Printf("Error get features from cache")
		resp.SetBody(fc)
		return nil
	}

	for _, feature := range cacheFC.Features {
		field := feature.Properties[IDName]
		str := fmt.Sprintf("%f", field)
		fid, err = strconv.ParseFloat(str, 64)
		if err != nil {
			fmt.Printf("Cannot convert %v in long", str)
			return nil
		}
		if int64(fid) == id {
			fc.AddFeatures(feature)
			resp.SetBody(fc)
			return nil
		}
	}
	resp.SetBody(fc)
	return nil
}

//TestSearchFeatures get feature details
func (t *Worker) TestSearchFeatures(request *config.SearchRequest, resp *config.FeaturesResponse) error {
	var (
		err     error
		cacheFC *geojson.FeatureCollection
	)
	name := request.Collection
	searchFieldName := request.SearchFiledName
	fc := geojson.NewFeatureCollection(nil)
	cacheFC, err = t.loadCacheFromFile(name, request)
	if err != nil {
		log.Printf("Error get features from cache")
		resp.SetBody(fc)
		return nil
	}

	for _, feature := range cacheFC.Features {
		field := feature.Properties[searchFieldName]
		for _, val := range request.SearchFiledValues {
			if val == field {
				fc.AddFeatures(feature)
			}
		}
	}
	resp.SetBody(fc)
	return nil
}

func filterFeatures(features []*geojson.Feature, n, e, s, w float64) (filtered []*geojson.Feature) {
	var (
		geomBBox *BBox
	)
	mapBBox := NewBBox(n, e, s, w)
	for _, feature := range features {
		geom, _ := feature.GetGeometry()
		coords := geom.GetCoordinates()

		switch geom.GetType() {
		case "Point":
			coord := coords.(geojson.Coordinate)
			geomBBox = NewBBoxFromCoordinate(coord)
		case "LineString":
			coord := coords.(geojson.Coordinates)
			geomBBox = NewBBoxFromCoordinates(coord)
		}
		if geomBBox.Interact(mapBBox) {
			filtered = append(filtered, feature)
		}
	}
	return
}
