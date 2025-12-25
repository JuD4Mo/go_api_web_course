package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/JuD4Mo/go_api_web_course/internal/course"
	"github.com/JuD4Mo/go_lib_response/response"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewCourseHttpServer(ctx context.Context, endpoints course.Endpoints) http.Handler {
	r := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodedError),
	}

	r.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateCourse,
		encodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetCourse,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAll,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCourse,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	r.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteCourse,
		encodeResponse,
		opts...,
	)).Methods("DELETE")

	return r
}

func decodeCreateCourse(_ context.Context, r *http.Request) (interface{}, error) {
	var courseReq course.CreateReq

	err := json.NewDecoder(r.Body).Decode(&courseReq)
	if err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: %v", err.Error()))
	}

	return courseReq, nil
}

func decodeGetCourse(_ context.Context, r *http.Request) (interface{}, error) {
	var getReq course.GetReq

	path := mux.Vars(r)
	id, ok := path["id"]

	if !ok || id == "" {
		return nil, response.BadRequest("id is required")
	}

	getReq.ID = id

	return getReq, nil
}

func decodeGetAll(_ context.Context, r *http.Request) (interface{}, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := course.GetAllReq{
		Name:  v.Get("name"),
		Limit: limit,
		Page:  page,
	}

	return req, nil
}

func decodeUpdateCourse(_ context.Context, r *http.Request) (interface{}, error) {
	var updateReq course.UpdateReq

	err := json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: %v", err.Error()))
	}

	path := mux.Vars(r)
	updateReq.ID = path["id"]

	return updateReq, nil
}

func decodeDeleteCourse(_ context.Context, r *http.Request) (interface{}, error) {

	err := authorization(r.Header.Get("Authorization"))
	if err != nil {
		return nil, response.Forbidden(err.Error())
	}

	var deleteReq course.DeleteReq
	path := mux.Vars(r)
	deleteReq.ID = path["id"]

	return deleteReq, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodedError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)

}

func authorization(token string) error {
	if token != os.Getenv("TOKEN") {
		return errors.New("invalid token")
	}

	return nil
}
