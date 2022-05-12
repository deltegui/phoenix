package validator

import (
	"reflect"

	playgroundValidator "github.com/go-playground/validator/v10"
)

type ValidationResult []ValidationError

type ValidationError struct {
	// Tag is the condition that have failed
	Tag string

	// Complete path to the field that have the error.
	Path string

	// Field is the name (and only the name) of the failing field
	Field string

	// Error is the stringified error
	Error string

	Value interface{}
	Kind  reflect.Kind
}

type PlaygroundValidator struct {
	validator *playgroundValidator.Validate
}

func New() PlaygroundValidator {
	return PlaygroundValidator{validator: playgroundValidator.New()}
}

func (val PlaygroundValidator) Validate(target interface{}) (ValidationResult, error) {
	err := val.validator.Struct(target)
	if err != nil {
		switch err.(type) {
		case playgroundValidator.ValidationErrors:
			{
				valErr := err.(playgroundValidator.ValidationErrors)
				return errorsToResult(valErr), nil
			}
		default:
			return nil, err
		}
	}
	return make(ValidationResult, 0), nil
}

func errorsToResult(ee playgroundValidator.ValidationErrors) ValidationResult {
	result := make(ValidationResult, len(ee))
	for i, e := range ee {
		result[i] = ValidationError{
			Tag:   e.ActualTag(),
			Path:  e.StructNamespace(),
			Field: e.Field(),
			Error: e.Error(),
			Value: e.Value(),
			Kind:  e.Kind(),
		}
	}
	return result
}
