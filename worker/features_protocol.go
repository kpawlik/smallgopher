package worker

import (
	"log"

	"github.com/kpawlik/geojson"
	"github.com/kpawlik/smallgopher/config"
)

// getFeatureID return feature id from acp
func (w *Worker) getFeatureID(idType string) (recID interface{}) {
	switch idType {
	case "int":
		recID = acp.GetLong()
	case "string":
		recID = acp.GetString()
	}
	return
}

// putFeatureID puts feature id to acp
func (w *Worker) putFeatureID(idType string, recID interface{}) {
	switch idType {
	case "int":
		acp.PutLong(recID.(int64))
	case "string":
		acp.PutString(recID.(string))
	}
}

// getGeometry gets geometry from acp
func (w *Worker) getGeometry(geomType string) (geom geojson.Geometry) {
	var coord geojson.Coordinate
	switch geomType {
	case "point":
		coord = toCoordinate(acp.GetCoord())
		geom = geojson.NewPoint(coord)
	case "chain":
		noOfCoords := acp.GetUint()
		line := geojson.NewLineString(nil)
		for j := 0; j < noOfCoords; j++ {
			line.AddCoordinates(toCoordinate(acp.GetCoord()))
		}
		geom = line
	case "area":
		noOfCoords := acp.GetUint()
		poly := geojson.NewPolygon(nil)
		coordinates := make(geojson.Coordinates, 0, 0)
		for j := 0; j < noOfCoords; j++ {
			coordinates = append(coordinates, toCoordinate(acp.GetCoord()))
		}
		poly.AddCoordinates(coordinates)
		geom = poly
	}
	return
}

// GetFeatures receive from ACP list of features
func (w *Worker) GetFeatures(request *config.FeaturesRequest, resp *config.FeaturesResponse) (respErr error) {
	var (
		geom            geojson.Geometry
		recordsIDs      []interface{}
		recordsIDsToGet []interface{}
	)
	defer func() {
		if r := recover(); r != nil {
			log.Panicf("PANIC in method Worker.Features: %v\n", r.(error))
		}
	}()
	// Send magik method name to execute
	acp.PutString("features()")
	idFieldName := request.IDFieldName
	idFieldType := request.IDFieldtype
	geomFieldName := request.GeometryFieldName
	geomType := request.GeometryFieldType
	collectionName := request.Collection
	// send protocol name to ACP
	acp.PutString(request.Dataset)
	acp.PutString(collectionName)
	acp.PutInt(int32(request.Limit))
	acp.PutFloat(float64(request.BBox[0]))
	acp.PutFloat(float64(request.BBox[1]))
	acp.PutFloat(float64(request.BBox[2]))
	acp.PutFloat(float64(request.BBox[3]))
	acp.PutString(idFieldName)
	acp.PutString(idFieldType)
	acp.PutString(geomFieldName)
	acp.PutString(geomType)
	// initialize goejson.FeatureCollection
	fc := geojson.NewFeatureCollection(nil)
	// get records ids from bounding box
	for {
		isMoreRecords := acp.GetUbyte()
		if isMoreRecords == NoMoreDataToGet {
			break
		}
		recordsIDs = append(recordsIDs, w.getFeatureID(idFieldType))
	}
	// first try to get records from cache
	collectionCache := w.getCollectionCache(collectionName)
	for _, recID := range recordsIDs {
		if feature, ok := collectionCache[recID]; ok {
			fc.AddFeatures(feature)
		} else {
			recordsIDsToGet = append(recordsIDsToGet, recID)
		}
	}
	// send to ACP how many records needs to be receive
	acp.PutUint(uint32(len(recordsIDsToGet)))
	// Get records field values
	for _, recID := range recordsIDsToGet {
		w.putFeatureID(idFieldType, recID)
		fields := make(map[string]interface{})
		fields[idFieldName] = recID
		geom = w.getGeometry(geomType)
		f := geojson.NewFeature(geom, fields, nil)
		fc.AddFeatures(f)
		collectionCache[recID] = f
	}
	resp.SetBody(fc)
	return nil
}

// GetFeature receive from ACP list of features
func (w *Worker) GetFeature(request *config.FeatureRequest, resp *config.FeaturesResponse) (respErr error) {
	var (
		geom geojson.Geometry
	)
	defer func() {
		if r := recover(); r != nil {
			log.Panicf("PANIC in method Worker.Features: %v\n", r.(error))
		}
	}()
	fc := geojson.NewFeatureCollection(nil)
	idFieldName := request.IDFieldName
	idFieldType := request.IDFieldType
	idFieldValue := request.IDFieldValue
	// send Magik method name to execute
	acp.PutString("feature()")
	// send data necessary to get record
	acp.PutString(request.Dataset)
	acp.PutString(request.Collection)
	acp.PutString(idFieldName)
	acp.PutString(idFieldType)
	// send record ID to get from ACP
	w.putFeatureID(idFieldType, idFieldValue)
	// send no of fields to get
	acp.PutUint(uint32(len(request.Fields)))
	// Send list of field names
	for _, field := range request.Fields {
		name := field.Name
		acp.PutString(name)
	}
	// check if record was found in SW database
	recordExists := acp.GetUbyte()
	if recordExists == RecordNotFound {
		log.Printf("Error Record with ID %s  not found\n", idFieldValue)
		resp.SetBody(fc)
		return nil
	}
	// Get records fields values
	fields := make(map[string]interface{})
	fields[idFieldName] = idFieldValue
	for _, field := range request.Fields {
		name := field.Name
		switch field.Type {
		case "string":
			fields[name] = acp.GetString()
		case "point", "chain", "area":
			geom = w.getGeometry(field.Type)
		}
	}
	fc.AddFeatures(geojson.NewFeature(geom, fields, nil))
	resp.SetBody(fc)
	return
}

// DumpFeatures receive from ACP list of features
func (w *Worker) DumpFeatures(request *config.DumpRequest, resp *config.FeaturesResponse) (respErr error) {
	var (
		geom       geojson.Geometry
		recordsIDs []interface{}
	)
	defer func() {
		if r := recover(); r != nil {
			log.Panicf("PANIC in method Worker.DumpFeatures: %v\n", r.(error))
		}
	}()
	// Send magik method name to execute
	acp.PutString("dump_features()")
	idFieldName := request.IDFieldName
	idFieldType := request.IDFieldtype
	geomFieldName := request.GeometryFieldName
	geomType := request.GeometryFieldType
	collectionName := request.Collection
	// send protocol name to ACP
	acp.PutString(request.Dataset)
	acp.PutString(collectionName)
	acp.PutInt(int32(request.Limit))
	acp.PutFloat(float64(request.BBox[0]))
	acp.PutFloat(float64(request.BBox[1]))
	acp.PutFloat(float64(request.BBox[2]))
	acp.PutFloat(float64(request.BBox[3]))
	acp.PutString(idFieldName)
	acp.PutString(idFieldType)
	acp.PutString(geomFieldName)
	acp.PutString(geomType)
	// initialize goejson.FeatureCollection
	fc := geojson.NewFeatureCollection(nil)
	// get records ids from bounding box
	for {
		isMoreRecords := acp.GetUbyte()
		// log.Printf("More records to get? 0-y, 1-n %d\n", moreRecordsToGet)
		if isMoreRecords == NoMoreDataToGet {
			break
		}
		recordsIDs = append(recordsIDs, w.getFeatureID(idFieldType))
	}
	// send no of fields to get
	acp.PutUint(uint32(len(request.Fields)))
	// Send list of field names
	for _, field := range request.Fields {
		name := field.Name
		acp.PutString(name)
	}
	// first try to get records from cache
	acp.PutUint(uint32(len(recordsIDs)))
	// Get records field values
	for _, recID := range recordsIDs {
		w.putFeatureID(idFieldType, recID)
		// Get records fields values
		fields := make(map[string]interface{})
		fields[idFieldName] = recID
		for _, field := range request.Fields {
			name := field.Name
			switch field.Type {
			case "string":
				fields[name] = acp.GetString()
			case "point", "chain", "area":
				geom = w.getGeometry(field.Type)
			}
		}
		fc.AddFeatures(geojson.NewFeature(geom, fields, nil))
	}
	resp.SetBody(fc)
	return nil
}

// SearchFeatures receive from ACP list of features
func (w *Worker) SearchFeatures(request *config.SearchRequest, resp *config.FeaturesResponse) (respErr error) {
	defer func() {
		if r := recover(); r != nil {
			log.Panicf("PANIC in method Worker.SearchFeatures: %v\n", r.(error))
		}
	}()
	// Send magik method name to execute
	acp.PutString("dump_features()")
	fc := geojson.NewFeatureCollection(nil)
	resp.SetBody(fc)
	return nil
}
