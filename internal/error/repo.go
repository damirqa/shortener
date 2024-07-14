package error

import "fmt"

// todo: где лучше хранить ошибки? да и еще с разным контекстом (бд, работа с файлом)
type UniqueConstraintError struct {
	Err error
}

func NewUniqueConstraintError(err error) *UniqueConstraintError {
	return &UniqueConstraintError{Err: err}
}

func (e *UniqueConstraintError) Error() string {
	return fmt.Sprintf("Unique constraint error: %v", e.Err)
}
