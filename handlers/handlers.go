package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/buraksekili/getir-task/persistence"
	"github.com/buraksekili/getir-task/persistence/inmemory"
	"github.com/buraksekili/getir-task/persistence/mongo"
)

const contentType = "application/json"

var (
	// ErrUnsupportedMediaType shows invalid content-type for HTTP requests.
	ErrUnsupportedMediaType = errors.New("unsupported content-type")

	// ErrMalformedEntity shows invalid request body for HTTP requests.
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnsupportedMethod shows invalid HTTP Method type for requests.
	ErrUnsupportMethod = errors.New("unsupported HTTP Method")

	// ErrNotFound shows that requested entity could not found.
	ErrNotFound = errors.New("non-existent entity")
)

// Agent represents an HTTP handler.
type Agent struct {
	l  *log.Logger
	db persistence.Database
	im inmemory.InMemory
}

// NewHTTPAgent creates a new HTTP Agent for HTTP requests
func NewHTTPAgent(l *log.Logger, db persistence.Database, memory inmemory.InMemory) *Agent {
	return &Agent{l, db, memory}
}

// Mongo serves HTTP handler for MongoDB related service.
func (a *Agent) Mongo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", contentType)
	if r.Method != http.MethodPost {
		encodeError(ErrUnsupportMethod, w)
		return
	}

	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		encodeError(ErrUnsupportedMediaType, w)
		return
	}

	fdr := FetchDataReq{}
	if err := json.NewDecoder(r.Body).Decode(&fdr); err != nil {
		encodeError(ErrMalformedEntity, w)
		return
	}
	if err := fdr.Validate(); err != nil {
		encodeError(err, w)
		return
	}

	filter := persistence.Filter{
		EndDate:       fdr.EndDate,
		StartDate:     fdr.StartDate,
		MaxTotalCount: fdr.MaxCount,
		MinTotalCount: fdr.MinCount,
	}
	result, err := a.db.FetchData(context.Background(), filter)
	if err != nil {
		encodeError(err, w)
		return
	}

	res := FetchDataRes{Code: 0, Msg: "Success", Records: result}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		encodeError(err, w)
		return
	}
	return
}

// InMemory serves HTTP hanndler for In-Memory.
func (a *Agent) InMemory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", contentType)
	if r.Method == http.MethodGet {
		key := r.URL.Query().Get("key")
		val := a.im.FetchData(key)
		if val == "" {
			encodeError(ErrNotFound, w)
			return
		}

		if err := json.NewEncoder(w).Encode(InMemoryRes{Key: key, Value: val}); err != nil {
			encodeError(err, w)
		}
		return
	}

	req := InMemoryCreateReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		encodeError(ErrMalformedEntity, w)
		return
	}
	a.im.StoreData(persistence.InMemoryKeyValue{Key: req.Key, Value: req.Value})
	w.WriteHeader(http.StatusCreated)
}

func encodeError(err error, w http.ResponseWriter) {
	res := ErrorRes{Error: err.Error()}
	switch {
	case strings.Contains(err.Error(), ErrMalformedEntity.Error()):
		w.WriteHeader(http.StatusBadRequest)
	case strings.Contains(err.Error(), mongo.ErrParseTime.Error()):
		w.WriteHeader(http.StatusBadRequest)
	case strings.Contains(err.Error(), persistence.ErrAggregation.Error()):
		w.WriteHeader(http.StatusBadRequest)
	case strings.Contains(err.Error(), ErrUnsupportMethod.Error()):
		w.WriteHeader(http.StatusMethodNotAllowed)
	case strings.Contains(err.Error(), ErrNotFound.Error()):
		w.WriteHeader(http.StatusNotFound)
	case strings.Contains(err.Error(), ErrUnsupportedMediaType.Error()):
		w.WriteHeader(http.StatusUnsupportedMediaType)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(res)
}
