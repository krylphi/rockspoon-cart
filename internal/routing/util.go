package routing

import (
	"encoding/json"
	"log"
	"net/http"
)

// HTTPHandler is handler function.
type HTTPHandler func(w http.ResponseWriter, r *http.Request)

// HTTPEndpoint is function for producing an HTTPResponse.
type HTTPEndpoint func(w http.ResponseWriter, r *http.Request) HTTPResponse

// HTTPResponse contains response data, status and headers.
type HTTPResponse interface {
	Headers() map[string]string
	Response() interface{}
	StatusCode() int
}

type response struct {
	Status     int
	Data       interface{}
	HeaderData map[string]string
}

func (e *response) Response() interface{} {
	return e.Data
}

func (e *response) StatusCode() int {
	return e.Status
}

func (e *response) Headers() map[string]string {
	return e.HeaderData
}

// JSON produces jsonified response.
func JSON(fn HTTPEndpoint) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		d := fn(w, r)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		for k, v := range d.Headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(d.StatusCode())
		err := json.NewEncoder(w).Encode(d.Response())
		if err != nil {
			log.Printf("JSON() error, while encoding response: %v", err.Error())
		}
	}
}

// Err produces custom error message HTTPResponse.
func Err(status int, data interface{}) HTTPResponse {
	return Resp(status, data)
}

// Resp produces regular HTTPResponse.
func Resp(status int, data interface{}) HTTPResponse {
	return &response{
		Status: status,
		Data:   data,
	}
}

// OK produces HTTPResponse with 200 status code.
func OK(d interface{}) HTTPResponse {
	return Resp(http.StatusOK, d)
}

// Created produces HTTPResponse with 201 status code.
func Created(d interface{}) HTTPResponse {
	return Resp(http.StatusCreated, d)
}

// BadRequestErrResp produces HTTPResponse woth err description and 401 status code.
func BadRequestErrResp(err error) HTTPResponse {
	return Err(http.StatusBadRequest, struct {
		Error string `json:"error"`
	}{Error: err.Error()})
}
