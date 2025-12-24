package course

import (
	"errors"
	"fmt"
)

var ErrInvalidStartDate = errors.New("invalid start date format")
var ErrInvalidEndDate = errors.New("invalid end date format")
var ErrNameRequired = errors.New("name is required")
var ErrStartDateRequired = errors.New("start date is required")
var ErrEndDateRequired = errors.New("end date is required")

type ErrCourseNotFound struct {
	CourseId string
}

func (e ErrCourseNotFound) Error() string {
	return fmt.Sprintf("course '%s' does not exist", e.CourseId)
}
