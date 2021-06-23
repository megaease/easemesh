package meshclient

import "github.com/pkg/errors"

var (
	ConflictError = errors.Errorf("resource already exists")
	NotFoundError = errors.Errorf("resource not found")
)

func IsConflictError(err error) (result bool) {
	if errors.Cause(err) == ConflictError {
		result = true
	}
	return
}

func IsNotFoundError(err error) (result bool) {
	if errors.Cause(err) == NotFoundError {
		result = true
	}
	return
}
