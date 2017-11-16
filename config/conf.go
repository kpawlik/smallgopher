package config

import (
	"encoding/json"
	"io/ioutil"
)

// WorkMode - enumerator with server mode types
type WorkMode int

const (
	// Version number
	Version = "0.9"
	// UnknownMode is unrecognized mode
	UnknownMode WorkMode = iota
	// NormalMode is production mode
	NormalMode
	// TestMode  to test communication between Acp and worker
	TestMode
)

// ServerConf server configuration
type ServerConf struct {
	Port     int
	Center   []float32
	Features []*FeatureConf
}

// WorkerConf workers configuration
type WorkerConf struct {
	Host string
	Name string
	Port int
}

// FieldConf is parameter name and type definition. Type could take values string, unsigned_int, signed_int, etc
type FieldConf struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
	Key         bool   `json:"key"`
	Searchable  bool   `json:"searchable"`
	Geometry    bool   `json:"geom"`
}

// FeatureConf is a definition of protocol. Contains name, list of entry parameters and list of results fields
type FeatureConf struct {
	Name        string       `json:"name"`
	DisplayName string       `json:"display_name"`
	Dataset     string       `json:"dataset"`
	MinZoom     int          `json:"min_zoom"`
	MaxZoom     int          `json:"max_zoom"`
	Fields      []*FieldConf `json:"fields"`
}

// TODO: add handling config errors eg.
// missing ID/GEOM etc

// GetIDField return primary key field name
func (f FeatureConf) GetIDField() *FieldConf {
	for _, field := range f.Fields {
		if field.Key {
			return field
		}
	}
	return nil
}

// GetGeomField return primary key field name
func (f FeatureConf) GetGeomField() *FieldConf {
	for _, field := range f.Fields {
		if field.Geometry {
			return field
		}
	}
	return nil
}

// FieldsWithoutID return list of fields except the key field
func (f FeatureConf) FieldsWithoutID() (fields []*FieldConf) {
	for _, field := range f.Fields {
		if !field.Key {
			fields = append(fields, field)
		}
	}
	return fields
}

//GetField return field configuration with name
func (f FeatureConf) GetField(name string) (field *FieldConf) {
	for _, f := range f.Fields {
		if f.Name == name {
			field = f
			return
		}
	}
	return
}

// Config application configuration structure
type Config struct {
	Server  ServerConf
	Workers []*WorkerConf
}

// GetFeaturesDef returns Protocol definition of nil if not found
func (c Config) GetFeaturesDef(name string) *FeatureConf {
	for _, prot := range c.Server.Features {
		if prot.Name == name {
			return prot
		}
	}
	return nil
}

// GetWorkerDef returns worker connection definition of nil if not found
func (c Config) GetWorkerDef(name string) *WorkerConf {
	for _, worker := range c.Workers {
		if worker.Name == name {
			return worker
		}
	}
	return nil
}

// ReadConf reads and decodes JSON from file
func ReadConf(filePath string) (conf *Config, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filePath); err == nil {
		conf, err = unmarshal(data)
	}
	return
}

func unmarshal(data []byte) (conf *Config, err error) {
	err = json.Unmarshal(data, &conf)
	return
}
