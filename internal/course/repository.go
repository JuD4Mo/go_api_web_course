package course

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/JuD4Mo/go_api_web_domain/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(ctx context.Context, course *domain.Course) error
		Get(ctx context.Context, id string) (*domain.Course, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Update(ctx context.Context, id string, name *string, startDate, endDate *time.Time) error
		Delete(ctx context.Context, id string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	repo struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepo(db *gorm.DB, log *log.Logger) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(ctx context.Context, course *domain.Course) error {
	result := repo.db.WithContext(ctx).Create(course)
	if result.Error != nil {
		return result.Error
	}
	repo.log.Println("course created with id: ", course.ID)
	return nil
}

func (repo *repo) Get(ctx context.Context, id string) (*domain.Course, error) {
	//Como el id está poblado y GORM detecta que es la PK filtra por ese parámetro
	course := domain.Course{
		ID: id,
	}
	// result := repo.db.Model(&Course{}).Where("id = ?", id).First(&course)
	result := repo.db.WithContext(ctx).First(&course)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrCourseNotFound{CourseId: id}
		}
		return nil, result.Error
	}

	return &course, nil
}

func (repo *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	var courses []domain.Course
	tx := repo.db.WithContext(ctx).Model(&courses)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at DESC").Find(&courses)

	if result.Error != nil {
		return nil, result.Error
	}

	return courses, nil
}

func (repo *repo) Update(ctx context.Context, id string, name *string, startDate, endDate *time.Time) error {
	values := make(map[string]interface{})

	if name != nil {
		values["name"] = *name
	}

	if startDate != nil {
		values["start_date"] = *startDate
	}

	if endDate != nil {
		values["end_date"] = *endDate
	}

	repo.log.Println(values)

	result := repo.db.WithContext(ctx).Model(&domain.Course{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrCourseNotFound{CourseId: id}
	}

	return nil
}

func (repo *repo) Delete(ctx context.Context, id string) error {
	course := domain.Course{
		ID: id,
	}

	result := repo.db.WithContext(ctx).Delete(&course)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrCourseNotFound{CourseId: id}
	}
	return nil
}

func (repo *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(domain.Course{})
	tx = applyFilters(tx, filters)

	result := tx.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}

	return tx
}
