package hw09structvalidator

import "errors"

var (
	errInvalidLength    = errors.New("invalid length")
	errInvalidMinLen    = errors.New("invalid min length")
	errInvalidRegexp    = errors.New("invalid regexp")
	errInvalidMinValue  = errors.New("invalid min value")
	errInvalidMaxValue  = errors.New("invalid max value")
	errUnsupportedType  = errors.New("unsupported field type")
	errUnknownValidator = errors.New("unknown validator")
)
