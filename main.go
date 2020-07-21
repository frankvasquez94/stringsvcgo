package main

import(
	"context"
	"errors"
	"strings"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"encoding/json"
	"log"
	httptransport "github.com/go-kit/kit/transport/http"
)

// 1. definicion de la interfaz del servicio

type StringService interface {
	UpperCase(string) (string, error)
	Count(string) int
}

// 2. Implementation

type stringService struct {}

func (stringService) UpperCase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

var ErrEmpty = errors.New("Empty string")


// DTO

type upperCaseRequest struct {
	S string `json:"s"`
}

type upperCaseResponse struct {
	V string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

// Endpoints

func makeUpperCaseEndpoint(svc stringService) endpoint.Endpoint {
	return func (_ context.Context, request interface{}) (interface{}, error) {
	    req := request.(upperCaseRequest)
	    v, err := svc.UpperCase(req.S)

	    if err != nil {
	    	return upperCaseResponse{v, err.Error()}, nil
		}

		return upperCaseResponse{v,""}, nil
	}
}

func makeCountEndpoint(svc stringService) endpoint.Endpoint {
	return func (_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}

// Transport
func decodeUpperCaseRequest(_ context.Context, r *http.Request) (interface{}, error){
	var request upperCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(&response)
}

func main(){
    svc := stringService{}
    upperCaseHandler := httptransport.NewServer(
    	makeUpperCaseEndpoint(svc),
    	decodeUpperCaseRequest,
    	encodeResponse,
    	)

    countHandler := httptransport.NewServer(
    	makeCountEndpoint(svc),
    	decodeCountRequest,
    	encodeResponse,)

    http.Handle("/uppercase", upperCaseHandler)
    http.Handle("/count",countHandler)

    log.Fatal(http.ListenAndServe(":8080",nil))

}

