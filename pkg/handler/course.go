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
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"

	httptransport "github.com/go-kit/kit/transport/http"
)

func NewCourseHttpServer(ctx context.Context, endpoints course.Endpoints) http.Handler {
	//r := mux.NewRouter()

	r := gin.Default()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodedError),
	}

	r.POST("/courses", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateCourse,
		encodeResponse,
		opts...,
	)))

	r.GET("/courses/:id", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetCourse,
		encodeResponse,
		opts...,
	)))

	r.GET("/courses", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAll,
		encodeResponse,
		opts...,
	)))

	r.PATCH("/courses/:id", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCourse,
		encodeResponse,
		opts...,
	)))

	r.DELETE("/courses/:id", ginDecode, gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteCourse,
		encodeResponse,
		opts...,
	)))

	return r
}

func ginDecode(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), "params", c.Params)
	c.Request = c.Request.WithContext(ctx)
}

func decodeCreateCourse(_ context.Context, r *http.Request) (interface{}, error) {
	var courseReq course.CreateReq

	err := json.NewDecoder(r.Body).Decode(&courseReq)
	if err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: %v", err.Error()))
	}

	return courseReq, nil
}

func decodeGetCourse(ctx context.Context, r *http.Request) (interface{}, error) {
	var getReq course.GetReq

	// path := mux.Vars(r)

	params := ctx.Value("params").(gin.Params)

	id := params.ByName("id")

	if id == "" {
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

func decodeUpdateCourse(ctx context.Context, r *http.Request) (interface{}, error) {
	var updateReq course.UpdateReq

	err := json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: %v", err.Error()))
	}

	// path := mux.Vars(r)
	params := ctx.Value("params").(gin.Params)
	updateReq.ID = params.ByName("id")

	return updateReq, nil
}

func decodeDeleteCourse(ctx context.Context, r *http.Request) (interface{}, error) {

	err := authorization(r.Header.Get("Authorization"))
	if err != nil {
		return nil, response.Forbidden(err.Error())
	}

	var deleteReq course.DeleteReq
	// path := mux.Vars(r)
	params := ctx.Value("params").(gin.Params)
	deleteReq.ID = params.ByName("id")

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
