package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte `validate:"minLen:10"`
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "Valid User",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Ivan",
				Age:    30,
				Email:  "ivan@mail.ru",
				Role:   "admin",
				Phones: []string{"89261234567", "89159876543"},
			},
			expectedErr: nil,
		},
		{
			name: "Invalid User Age",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Ivan",
				Age:    15,
				Email:  "ivan@mail.ru",
				Role:   "admin",
				Phones: []string{"89261234567", "89159876543"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: fmt.Errorf("must be at least 18")},
			},
		},
		{
			name: "Invalid User Phone length",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Ivan",
				Age:    19,
				Email:  "ivan@mail.ru",
				Role:   "admin",
				Phones: []string{"9261234567", "9159876543"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Phones", Err: fmt.Errorf("length must be 11")},
			},
		},
		{
			name: "Valid App",
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			name: "Invalid App Version length",
			in: App{
				Version: "1.0.0.0",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: fmt.Errorf("length must be 5")},
			},
		},
		{
			name: "Valid Response",
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			name: "Invalid Response Code",
			in: Response{
				Code: 400,
				Body: "Bad Request",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: fmt.Errorf("must be one of 200,404,500")},
			},
		},
		{
			name: "Invalid Email format",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Ivan",
				Age:    35,
				Email:  "ivan@mail",
				Role:   "admin",
				Phones: []string{"89261234567", "89159876543"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Email", Err: fmt.Errorf("must match regexp ^\\w+@\\w+\\.\\w+$")},
			},
		},
		{
			name: "Multiple validation errors",
			in: User{
				ID:     "1",
				Name:   "Ivan",
				Age:    15,
				Email:  "ivan@mail",
				Role:   "user",
				Phones: []string{"9261234567", "9159876543"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: fmt.Errorf("length must be 36")},
				ValidationError{Field: "Age", Err: fmt.Errorf("must be at least 18")},
				ValidationError{Field: "Email", Err: fmt.Errorf("must match regexp ^\\w+@\\w+\\.\\w+$")},
				ValidationError{Field: "Role", Err: fmt.Errorf("must be one of admin,stuff")},
				ValidationError{Field: "Phones", Err: fmt.Errorf("length must be 11")},
			},
		},
		{
			name: "Valid Token",
			in: Token{
				Header:    []byte("validheader"),
				Payload:   []byte(""),
				Signature: []byte(""),
			},
			expectedErr: nil,
		},
		{
			name: "Invalid Header length",
			in: Token{
				Header:    []byte("short"),
				Payload:   []byte(""),
				Signature: []byte(""),
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Header", Err: fmt.Errorf("minimum length is 10")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if !reflect.DeepEqual(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			_ = tt
		})
	}
}
