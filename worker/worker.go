package worker

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/kpawlik/geojson"
)

var (
	acp IAcp
)

const (
	// MoreDataToGet more records to get from ACP
	MoreDataToGet = 0
	// NoMoreDataToGet no more records to get from ACP
	NoMoreDataToGet = 1
	// RecordNotFound record was not found
	RecordNotFound = 1
)

// StartWorker register struct and start RPC worker server
func StartWorker(port string, name string, offline bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Panic(err)
		}
	}()
	if offline {
		acp = NewTestAcp(name)
	} else {
		acp = NewAcp(name)
	}
	if err := acp.Connect(name, 0, 1); err != nil {
		log.Panicf("ACP Connection error: %v\n", err)
	}
	// register worker for RPC server
	worker := NewWorker(port, name)
	rpc.Register(worker)
	rpc.HandleHTTP()
	// start listening for requests from HTTP server
	if listener, err := net.Listen("tcp", port); err != nil {
		log.Panicf("Start worker error on port %s. Error: %v\n", port, err)
	} else {
		log.Printf("Worker started at port: %s\n", port)
		log.Fatalf("RPC SERVER ERROR! %s\n", http.Serve(listener, nil))
	}
}

// Worker type to wrap RPC communication
type Worker struct {
	Port       string
	WorkerName string
	Cache      Cache
}

//NewWorker creates new instance of Worker
func NewWorker(port string, name string) *Worker {
	return &Worker{
		Port:       port,
		WorkerName: name,
		Cache:      NewCache(),
	}
}

func (t *Worker) getCollectionCache(name string) FeaturesMap {
	cache, ok := t.Cache[name]
	if !ok {
		cache = NewFeaturesMap()
		t.Cache[name] = cache
	}
	return cache
}

// FeaturesMap struct to cache Features by id
type FeaturesMap map[interface{}]*geojson.Feature

// NewFeaturesMap creates new empty FeaturesMap
func NewFeaturesMap() FeaturesMap {
	return make(FeaturesMap)
}

//Cache cache features by collection name and id
type Cache map[interface{}]FeaturesMap

//NewCache creates new empty map
func NewCache() Cache {
	return make(Cache)
}
