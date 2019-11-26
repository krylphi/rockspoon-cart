package routing

import (
	"encoding/json"
	"net/http"
)

type HttpHandler func(w http.ResponseWriter, r *http.Request)

type HttpEndpoint func(w http.ResponseWriter, r *http.Request) HttpResponse

type HttpResponse interface {
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

func Json(fn HttpEndpoint) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		d := fn(w, r)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		for k, v := range d.Headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(d.StatusCode())
		json.NewEncoder(w).Encode(d.Response())
	}
}

func Err(status int, data interface{}) *Response {
	return Resp(status, data)
}

func Resp(status int, data interface{}) *Response {
	resp := &Response{}
	resp.Status = status
	resp.Data = data
	return resp
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
