package transform

import "fmt"

// ErrInvalidType indicates a value was not of the expected string type.
type ErrInvalidType struct {
	Key string
	Got any
}

func (e *ErrInvalidType) Error() string {
	return fmt.Sprintf("type is not string for key %q: %T", e.Key, e.Got)
}

func NewErrInvalidType(key string, got any) *ErrInvalidType {
	return &ErrInvalidType{Key: key, Got: got}
}
