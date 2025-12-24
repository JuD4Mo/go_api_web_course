package course

import (
	"context"
	"log"
	"time"

	"github.com/JuD4Mo/go_api_web_domain/domain"
)

type (
	Service interface {
		Create(ctx context.Context, name string, startDate, endDate string) (*domain.Course, error)
		Get(ctx context.Context, id string) (*domain.Course, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Update(ctx context.Context, id string, name, startDate, endDate *string) error
		Delete(ctx context.Context, id string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}

	Filters struct {
		Name string
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, name string, startDate, endDate string) (*domain.Course, error) {

	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, ErrInvalidStartDate
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, ErrInvalidEndDate
	}

	if startDateParsed.After(endDateParsed) {
		s.log.Println(ErrEndLessStart)
		return nil, ErrEndLessStart
	}

	course := &domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}
	err = s.repo.Create(ctx, course)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (s service) Get(ctx context.Context, id string) (*domain.Course, error) {
	course, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	courses, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (s service) Update(ctx context.Context, id string, name, startDate, endDate *string) error {

	var startDateParsed, endDateParsed *time.Time

	courseObj, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if startDate != nil && *startDate != "" {
		date, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			return ErrInvalidStartDate
		}
		if date.After(courseObj.EndDate) {
			s.log.Println(ErrEndLessStart)
			return ErrEndLessStart
		}
		startDateParsed = &date
	}

	if endDate != nil && *endDate != "" {
		date, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			return ErrInvalidEndDate
		}

		if courseObj.StartDate.After(date) {
			s.log.Println(ErrEndLessStart)
			return ErrEndLessStart
		}
		endDateParsed = &date
	}

	return s.repo.Update(ctx, id, name, startDateParsed, endDateParsed)
}

func (s service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	num, err := s.repo.Count(ctx, filters)
	if err != nil {
		return 0, err
	}
	return num, nil
}
