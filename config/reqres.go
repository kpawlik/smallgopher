package config

import (
	"strconv"

	"fmt"

	"github.com/kpawlik/geojson"
)

//FeatureID id of feature
type FeatureID interface{}

func NewFeatureID(id interface{}) interface{} {
	return id
}

// BodyElement is a type which is a part of JSON response
type BodyElement map[string]interface{}

// Body list of BodyElements. JSON response object.
type Body []BodyElement

// FeaturesResponse struct
// Body - result map (field, value) to json
type FeaturesResponse struct {
	Body  *geojson.FeatureCollection
	Error error
}

// GetError return response error
func (r *FeaturesResponse) GetError() error {
	return r.Error
}

//GetBody return response body
func (r *FeaturesResponse) GetBody() interface{} {
	return r.Body
}

//SetError set response error
func (r *FeaturesResponse) SetError(err error) {
	r.Error = err
}

// SetBody set response body
func (r *FeaturesResponse) SetBody(body interface{}) {
	r.Body = body.(*geojson.FeatureCollection)
}

// Request interface
type Request interface {
	GetIDFieldName() string
}

// FeaturesRequest struct
type FeaturesRequest struct {
	Dataset           string
	Collection        string
	IDFieldName       string
	IDFieldtype       string
	GeometryFieldName string
	GeometryFieldType string
	BBox              []float32
	Limit             int
}

//NewFeaturesRequest creates new instance of FeaturesRequest
func NewFeaturesRequest(featureConf *FeatureConf, bb []float32, limit int) *FeaturesRequest {
	IDField := featureConf.GetIDField()
	GeomField := featureConf.GetGeomField()
	return &FeaturesRequest{
		Dataset:           featureConf.Dataset,
		Collection:        featureConf.Name,
		IDFieldName:       IDField.Name,
		IDFieldtype:       IDField.Type,
		GeometryFieldName: GeomField.Name,
		GeometryFieldType: GeomField.Type,
		BBox:              bb,
		Limit:             limit,
	}
}

// GetIDFieldName returns name of id field
func (r *FeaturesRequest) GetIDFieldName() string {
	return r.IDFieldName
}

// FeatureRequest struct
type FeatureRequest struct {
	Dataset      string
	Collection   string
	IDFieldName  string
	IDFieldValue FeatureID
	IDFieldType  string
	Fields       []*FieldConf
}

//NewFeatureRequest creates new instance of FeatureRequest
func NewFeatureRequest(featureConf *FeatureConf, id string) (*FeatureRequest, error) {
	var (
		err   error
		recID interface{}
	)
	idField := featureConf.GetIDField()
	if idField.Type == "int" {
		if recID, err = strconv.ParseInt(id, 10, 64); err != nil {
			err = fmt.Errorf("Cannot convert %s to int", id)
			return nil, err
		}
	} else {
		recID = id
	}

	return &FeatureRequest{
		Dataset:      featureConf.Dataset,
		Collection:   featureConf.Name,
		IDFieldName:  idField.Name,
		IDFieldValue: recID,
		IDFieldType:  idField.Type,
		Fields:       featureConf.FieldsWithoutID(),
	}, nil
}

// GetIDFieldName returns name of id field
func (r *FeatureRequest) GetIDFieldName() string {
	return r.IDFieldName
}

// DumpRequest struct
type DumpRequest struct {
	Dataset           string
	Collection        string
	IDFieldName       string
	IDFieldtype       string
	GeometryFieldName string
	GeometryFieldType string
	Fields            []*FieldConf
	BBox              []float32
	Limit             int
}

//NewDumpRequest creates new instance of FeaturesRequest
func NewDumpRequest(featureConf *FeatureConf, bb []float32, limit int) *DumpRequest {
	IDField := featureConf.GetIDField()
	GeomField := featureConf.GetGeomField()
	return &DumpRequest{
		Dataset:           featureConf.Dataset,
		Collection:        featureConf.Name,
		IDFieldName:       IDField.Name,
		IDFieldtype:       IDField.Type,
		GeometryFieldName: GeomField.Name,
		GeometryFieldType: GeomField.Type,
		Fields:            featureConf.FieldsWithoutID(),
		BBox:              bb,
		Limit:             limit,
	}
}

// SearchRequest for search feature
type SearchRequest struct {
	Datset            string
	Collection        string
	SearchFiledName   string
	SearchFiledType   string
	IDFieldName       string
	SearchFiledValues []string
}

// NewSearchRequest return new search request
func NewSearchRequest(featureConf *FeatureConf, name, fieldType string, values []string) *SearchRequest {
	return &SearchRequest{
		Datset:            featureConf.Dataset,
		Collection:        featureConf.Name,
		SearchFiledName:   name,
		SearchFiledType:   fieldType,
		SearchFiledValues: values,
		IDFieldName:       featureConf.GetIDField().Name,
	}
}

// GetIDFieldName returns name of id field
func (r *SearchRequest) GetIDFieldName() string {
	return r.IDFieldName
}
