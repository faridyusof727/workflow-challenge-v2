package render

import "errors"

var (
	ErrInvalidWorkflowID   = errors.New("invalid workflow id")
	ErrNotFound            = errors.New("not found")
	ErrInternalServerError = errors.New("internal server error")
)

func GetAPIError(err error) error {
	switch err {
	case ErrInvalidWorkflowID:
		return err
	case ErrNotFound:
		return err
	default:
		return ErrInternalServerError
	}
}
