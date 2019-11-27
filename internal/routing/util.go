package routing

import (
	"encoding/json"
	"log"
	"net/http"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request)

type HTTPEndpoint func(w http.ResponseWriter, r *http.Request) HTTPResponse

type HTTPResponse interface {
	Headers() map[string]string
	Response() interface{}
	StatusCode() int
}

type Response struct {
	Status     int
	Data       interface{}
	HeaderData map[string]string
}

func (e *Response) Response() interface{} {
	return e.Data
}

func (e *Response) StatusCode() int {
	return e.Status
}

func (e *Response) Headers() map[string]string {
	return e.HeaderData
}

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

func Err(status int, data interface{}) *Response {
	return Resp(status, data)
}

func Resp(status int, data interface{}) *Response {
	return &Response{
		Status: status,
		Data:   data,
	}
}

func OK(d interface{}) *Response {
	return Resp(http.StatusOK, d)
}

func Created(d interface{}) *Response {
	return Resp(http.StatusCreated, d)
}

func BadRequestErrResp(err error) *Response {
	return Err(http.StatusBadRequest, struct {
		Error string `json:"error"`
	}{Error: err.Error()})
}
