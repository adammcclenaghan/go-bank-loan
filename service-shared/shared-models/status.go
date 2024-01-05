package shared_models

import "github.com/go-playground/validator/v10"

//Enum defines an abstraction for validating whether a status is valid or not.
type Enum interface {
	IsValid() bool
}

type Status string

const (
	Pending   Status = "pending"
	Completed Status = "completed"
	Rejected  Status = "rejected"
)

func (s Status) IsValid() bool {
	switch s {
	case Pending, Completed, Rejected:
		return true
	}

	return false
}

// ValidStatus is a validator function used to ensure that string representations of a loan's status are valid.
var ValidStatus validator.Func = func(f1 validator.FieldLevel) bool {
	statusString, ok := f1.Field().Interface().(Enum)
	if ok {
		return statusString.IsValid()
	}

	return false
}
