package hw09structvalidator

import "fmt"

var (
	errInvalidLength    = fmt.Errorf("invalid length")
	errInvalidMinLen    = fmt.Errorf("invalid min length")
	errInvalidRegexp    = fmt.Errorf("invalid regexp")
	errInvalidMinValue  = fmt.Errorf("invalid min value")
	errInvalidMaxValue  = fmt.Errorf("invalid max value")
	errUnsupportedType  = fmt.Errorf("unsupported field type")
	errUnknownValidator = fmt.Errorf("unknown validator")
)
