package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID        string `json:"id" validate:"len:36"`
		Name      string
		Age       int        `validate:"min:18|max:25"`
		Ages      []int      `validate:"min:18|max:25"`
		Email     string     `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role      UserRole   `validate:"in:admin,stuff"`
		Phones    []string   `validate:"len:11"`
		JobEmails []string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		JobRoles  []UserRole `validate:"in:admin,stuff"`
		Number    int        `validate:"in:10,11"`
		Numbers   []int      `validate:"in:10,11"`
		meta      json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	InvalidTypes struct {
		len   int      `validate:"len:10"`
		regex int      `validate:"rexexp:^\\w+@\\w+\\.\\w+$"`
		in    chan int `validate:"in:admin,stuff"`
		min   chan int `validate:"min:10"`
		max   chan int `validate:"max:20"`
	}

	UnknownValidation struct {
		len int `validate:"xxx:100"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in             interface{}
		expectedErrors ValidationErrors
	}{
		// success
		{
			in: User{
				ID:        "123456798132416549687946513215649687",
				Name:      "name",
				Age:       19,
				Email:     "abs@ddd.com",
				Role:      "admin",
				Phones:    []string{"12345678901", "12345678902"},
				JobEmails: []string{"abs@ddd.com", "abs2@ddd.com"},
				JobRoles:  []UserRole{"admin", "stuff"},
				Number:    11,
				Numbers:   []int{10, 11},
				meta:      json.RawMessage{},
			},
			expectedErrors: nil,
		},
		// zero values
		{
			in: User{},
			expectedErrors: ValidationErrors{
				ValidationError{Field: "ID", Err: LenValidationError{36}},
				ValidationError{Field: "Age", Err: MinValidationError{18}},
				ValidationError{Field: "Role", Err: InValidationError{[]string{"admin", "stuff"}}},
				ValidationError{Field: "Number", Err: InValidationError{[]string{"10", "11"}}},
			},
		},
		// without validation tags
		{
			in: Token{
				Header:    []byte{0, 1, 2},
				Payload:   []byte{1, 2, 3},
				Signature: []byte{2, 3, 4},
			},
			expectedErrors: nil,
		},
		// invalid field types
		{
			in: InvalidTypes{
				len:   10,
				regex: 20,
				in:    make(chan int),
			},
			expectedErrors: nil,
		},
		// unknown validation
		{
			in: UnknownValidation{
				len: 10,
			},
			expectedErrors: nil,
		},
		// len
		{
			in:             App{Version: "123"},
			expectedErrors: ValidationErrors{ValidationError{Field: "Version", Err: LenValidationError{5}}},
		},
		// len slice
		{
			in: User{
				ID:     "123456798132416549687946513215649687",
				Age:    19,
				Role:   "admin",
				Number: 11,
				Phones: []string{"12345678901", "1234567890"},
			},
			expectedErrors: ValidationErrors{ValidationError{Field: "Phones[1]", Err: LenValidationError{11}}},
		},
		// regex
		{
			in: User{
				ID:     "123456798132416549687946513215649687",
				Age:    19,
				Role:   "admin",
				Number: 11,
				Email:  "asd",
			},
			expectedErrors: ValidationErrors{ValidationError{Field: "Email", Err: RegexpValidationError{"^\\w+@\\w+\\.\\w+$"}}},
		},
		// regex slice
		{
			in: User{
				ID:        "123456798132416549687946513215649687",
				Age:       19,
				Role:      "admin",
				Number:    11,
				JobEmails: []string{"abs@ddd.com", "asd"},
			},
			expectedErrors: ValidationErrors{ValidationError{Field: "JobEmails[1]", Err: RegexpValidationError{"^\\w+@\\w+\\.\\w+$"}}},
		},
		// in
		{
			in: User{
				ID:     "123456798132416549687946513215649687",
				Age:    19,
				Role:   "a",
				Number: 12,
			},
			expectedErrors: ValidationErrors{
				ValidationError{Field: "Role", Err: InValidationError{[]string{"admin", "stuff"}}},
				ValidationError{Field: "Number", Err: InValidationError{[]string{"10", "11"}}},
			},
		},
		// in slice
		{
			in: User{
				ID:       "123456798132416549687946513215649687",
				Age:      19,
				Role:     "admin",
				JobRoles: []UserRole{"admin", "s"},
				Number:   11,
				Numbers:  []int{11, 12},
			},
			expectedErrors: ValidationErrors{
				ValidationError{Field: "JobRoles[1]", Err: InValidationError{[]string{"admin", "stuff"}}},
				ValidationError{Field: "Numbers[1]", Err: InValidationError{[]string{"10", "11"}}},
			},
		},
		// min, min slice,
		{
			in: User{
				ID:     "123456798132416549687946513215649687",
				Age:    13,
				Ages:   []int{20, 13},
				Role:   "admin",
				Number: 11,
			},
			expectedErrors: ValidationErrors{
				ValidationError{Field: "Age", Err: MinValidationError{18}},
				ValidationError{Field: "Ages[1]", Err: MinValidationError{18}},
			},
		},
		// max, max slice
		{
			in: User{
				ID:     "123456798132416549687946513215649687",
				Age:    30,
				Ages:   []int{20, 30},
				Role:   "admin",
				Number: 11,
			},
			expectedErrors: ValidationErrors{
				ValidationError{Field: "Age", Err: MaxValidationError{25}},
				ValidationError{Field: "Ages[1]", Err: MaxValidationError{25}},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			if tt.expectedErrors == nil {
				require.NoError(t, nil)
				return
			}
			validationErrors, ok := err.(ValidationErrors)
			if !ok {
				require.Fail(t, "Validate() should return ValidationErrors")
			}
			if len(tt.expectedErrors) != len(validationErrors) {
				require.Fail(t, "Error slices len doesn't match. Expected: %#v, Actual %#v", tt.expectedErrors, validationErrors)
			}
			for j, expectedError := range tt.expectedErrors {
				if !errors.Is(validationErrors[j].Err, expectedError.Err) {
					require.Fail(t, fmt.Sprintf("Error should be of type '%T'. Actual '%T'\n", expectedError.Err, validationErrors[j].Err))
				}
			}

			require.Equal(t, tt.expectedErrors, err)
		})
	}
}
