package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errs []string
	for _, ve := range v {
		errs = append(errs, fmt.Sprintf("%s: %s", ve.Field, ve.Err.Error()))
	}
	return strings.Join(errs, "; ")
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	var validationErrors ValidationErrors

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		validators := strings.Split(validateTag, "|")
		for _, validator := range validators {
			if err := validateField(field.Name, fieldValue, validator); err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateString(fieldName, fieldValue, validatorName, validatorArg string) error {
	switch validatorName {
	case "len":
		length, err := strconv.Atoi(validatorArg)
		if err != nil {
			return fmt.Errorf("invalid length: %s", validatorArg)
		}
		if len(fieldValue) != length {
			return fmt.Errorf("length must be %d", length)
		}
	case "minLen":
		minLength, err := strconv.Atoi(validatorArg)
		if err != nil {
			return fmt.Errorf("invalid min length: %s", validatorArg)
		}
		if len(fieldValue) < minLength {
			return fmt.Errorf("minimum length is %d", minLength)
		}
	case "regexp":
		re, err := regexp.Compile(validatorArg)
		if err != nil {
			return fmt.Errorf("invalid regexp: %s", validatorArg)
		}
		if !re.MatchString(fieldValue) {
			return fmt.Errorf("must match regexp %s", validatorArg)
		}
	case "in":
		allowedValues := strings.Split(validatorArg, ",")
		for _, val := range allowedValues {
			if fieldValue == val {
				return nil
			}
		}
		return fmt.Errorf("must be one of %s", validatorArg)
	default:
		return fmt.Errorf("unknown validator: %s", validatorName)
	}
	return nil
}

func validateField(fieldName string, fieldValue reflect.Value, validator string) error {
	parts := strings.SplitN(validator, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid validator format: %s", validator)
	}

	validatorName := parts[0]
	validatorArg := parts[1]

	switch fieldValue.Kind() {
	case reflect.String:
		return validateString(fieldName, fieldValue.String(), validatorName, validatorArg)
	case reflect.Int:
		return validateInt(fieldName, int(fieldValue.Int()), validatorName, validatorArg)
	case reflect.Slice:
		if fieldValue.Type().Elem().Kind() == reflect.String {
			for i := 0; i < fieldValue.Len(); i++ {
				if err := validateString(fieldName, fieldValue.Index(i).String(), validatorName, validatorArg); err != nil {
					return err
				}
			}
		} else if fieldValue.Type().Elem().Kind() == reflect.Int {
			for i := 0; i < fieldValue.Len(); i++ {
				if err := validateInt(fieldName, int(fieldValue.Index(i).Int()), validatorName, validatorArg); err != nil {
					return err
				}
			}
		} else if fieldValue.Type().Elem().Kind() == reflect.Uint8 {
			return validateString(fieldName, string(fieldValue.Bytes()), validatorName, validatorArg)
		}
	default:
		return fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
	}

	return nil
}

func validateInt(fieldName string, fieldValue int, validatorName, validatorArg string) error {
	switch validatorName {
	case "min":
		min, err := strconv.Atoi(validatorArg)
		if err != nil {
			return fmt.Errorf("invalid min value: %s", validatorArg)
		}
		if fieldValue < min {
			return fmt.Errorf("must be at least %d", min)
		}
	case "max":
		max, err := strconv.Atoi(validatorArg)
		if err != nil {
			return fmt.Errorf("invalid max value: %s", validatorArg)
		}
		if fieldValue > max {
			return fmt.Errorf("must be at most %d", max)
		}
	case "in":
		allowedValues := strings.Split(validatorArg, ",")
		for _, val := range allowedValues {
			if intVal, err := strconv.Atoi(val); err == nil && fieldValue == intVal {
				return nil
			}
		}
		return fmt.Errorf("must be one of %s", validatorArg)
	default:
		return fmt.Errorf("unknown validator: %s", validatorName)
	}
	return nil
}
