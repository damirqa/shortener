package error

import "fmt"

type UniqueConstraintError struct {
	Err error
}

func NewUniqueConstraintError(err error) *UniqueConstraintError {
	return &UniqueConstraintError{Err: err}
}

func (e *UniqueConstraintError) Error() string {
	return fmt.Sprintf("Unique constraint error: %v", e.Err)
}
