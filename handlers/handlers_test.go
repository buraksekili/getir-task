package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/buraksekili/getir-task/persistence/inmemory"
	"github.com/buraksekili/getir-task/persistence/mocks"
)

type testRequest struct {
	client      *http.Client
	method      string
	url         string
	contentType string
	body        io.Reader
}

func (tr testRequest) make() (*http.Response, error) {
	req, err := http.NewRequest(tr.method, tr.url, tr.body)
	if err != nil {
		return nil, err
	}
	if tr.contentType != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Referer", "http://localhost")
	return tr.client.Do(req)
}

func newServer() *httptest.Server {
	mongoMock := mocks.New(map[string]string{})
	inMemory := inmemory.New(map[string]string{})
	agent := NewHTTPAgent(log.New(os.Stdout, "test ", log.LstdFlags), mongoMock, inMemory)

	sm := http.NewServeMux()
	sm.HandleFunc("/", agent.Mongo)
	sm.HandleFunc("/in-memory", agent.InMemory)

	return httptest.NewServer(sm)
}

func toJSON(data interface{}) string {
	jd, _ := json.Marshal(data)
	return string(jd)
}

func TestMongoHandler(t *testing.T) {
	ts := newServer()
	defer ts.Close()
	client := ts.Client()

	data := toJSON(FetchDataReq{StartDate: "2016-01-26", EndDate: "2018-02-02", MinCount: 2700, MaxCount: 3000})

	cases := []struct {
		method      string
		desc        string
		req         string
		contentType string
		status      int
	}{
		{http.MethodPost, "fetch data with valid body", data, contentType, http.StatusOK},
		{http.MethodGet, "fetch data with invalid HTTP Method", data, contentType, http.StatusMethodNotAllowed},
		{http.MethodPost, "fetch data with invalid body", "", contentType, http.StatusBadRequest},
		{http.MethodPost, "fetch data with missing content type", data, "", http.StatusUnsupportedMediaType},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      tc.method,
			contentType: tc.contentType,
			url:         fmt.Sprintf("%s/mongo", ts.URL),
			body:        strings.NewReader(tc.req),
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}

func TestInMemory(t *testing.T) {
	ts := newServer()
	defer ts.Close()
	client := ts.Client()

	data := toJSON(InMemoryCreateReq{Key: "key", Value: "val"})

	cases := []struct {
		method      string
		desc        string
		req         string
		contentType string
		key         string
		status      int
	}{
		{
			method:      http.MethodPost,
			desc:        "create a new data",
			req:         data,
			contentType: contentType,
			status:      http.StatusCreated,
		},
		{
			method:      http.MethodPost,
			desc:        "create a new data with invalid body",
			contentType: contentType,
			status:      http.StatusBadRequest,
		},
		{
			method:      http.MethodGet,
			desc:        "get a valid data",
			req:         data,
			contentType: contentType,
			status:      http.StatusOK,
			key:         "key",
		},
		{
			method:      http.MethodGet,
			desc:        "get a invalid data",
			req:         data,
			contentType: contentType,
			status:      http.StatusNotFound,
			key:         "invalid",
		},
	}

	for _, tc := range cases {
		req := testRequest{
			client:      client,
			method:      tc.method,
			contentType: tc.contentType,
			url:         fmt.Sprintf("%s/in-memory?key=%s", ts.URL, tc.key),
			body:        strings.NewReader(tc.req),
		}
		res, err := req.make()
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, res.StatusCode, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, res.StatusCode))
	}
}
