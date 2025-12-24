package course

import (
	"context"
	"errors"

	"github.com/JuD4Mo/go_api_web_meta/meta"
	"github.com/JuD4Mo/go_lib_response/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (response interface{}, err error)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	GetReq struct {
		ID string
	}

	GetAllReq struct {
		Name  string
		Limit int
		Page  int
	}

	DeleteReq struct {
		ID string
	}

	UpdateReq struct {
		ID        string
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
	}

	Config struct {
		LimitPage string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateReq)

		if req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.StartDate == "" {
			return nil, response.BadRequest(ErrStartDateRequired.Error())
		}

		if req.EndDate == "" {
			return nil, response.BadRequest(ErrEndDateRequired.Error())
		}

		course, err := s.Create(ctx, req.Name, req.StartDate, req.EndDate)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", course, nil), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReq)

		course, err := s.Get(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrCourseNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", course, nil), nil

	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllReq)

		filters := Filters{
			Name: req.Name,
		}

		num, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, num, config.LimitPage)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		courses, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", courses, meta), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateReq)

		if req.Name != nil && *req.Name == "" {
			return nil, response.BadRequest("name cannot be empty")
		}

		if req.StartDate != nil && *req.StartDate == "" {
			return nil, response.BadRequest("start date is required")
		}

		if req.EndDate != nil && *req.EndDate == "" {
			return nil, response.BadRequest("end date is required")
		}

		err := s.Update(ctx, req.ID, req.Name, req.StartDate, req.EndDate)
		if err != nil {

			if errors.As(err, &ErrCourseNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			if err == ErrInvalidStartDate || err == ErrInvalidEndDate {
				return nil, response.BadRequest(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}
		return response.Created("updated", nil, nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteReq)

		err := s.Delete(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrCourseNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("deleted", nil, nil), nil
	}
}
