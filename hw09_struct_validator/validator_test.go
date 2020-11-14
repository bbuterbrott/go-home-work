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
		in          interface{}
		expectedErr error
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
			expectedErr: nil,
		},
		// zero values
		{
			in: User{},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: errors.New("length must be 36")},
				ValidationError{Field: "Age", Err: errors.New("value must be greater than '18'")},
				ValidationError{Field: "Role", Err: errors.New("value must be one of '[admin stuff]'")},
				ValidationError{Field: "Number", Err: errors.New("value must be one of '[10 11]'")},
			},
		},
		// without validation tags
		{
			in: Token{
				Header:    []byte{0, 1, 2},
				Payload:   []byte{1, 2, 3},
				Signature: []byte{2, 3, 4},
			},
			expectedErr: nil,
		},
		// invalid field types
		{
			in: InvalidTypes{
				len:   10,
				regex: 20,
				in:    make(chan int),
			},
			expectedErr: nil,
		},
		// unknown validation
		{
			in: UnknownValidation{
				len: 10,
			},
			expectedErr: nil,
		},
		// len
		{
			in:          App{Version: "123"},
			expectedErr: ValidationErrors{ValidationError{Field: "Version", Err: errors.New("length must be 5")}},
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
			expectedErr: ValidationErrors{ValidationError{Field: "Phones", Err: errors.New("all elements in slice must have length 11")}},
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
			expectedErr: ValidationErrors{ValidationError{Field: "Email", Err: errors.New("must match regex '^\\w+@\\w+\\.\\w+$'")}},
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
			expectedErr: ValidationErrors{ValidationError{Field: "JobEmails", Err: errors.New("all elements in slice must match regex '^\\w+@\\w+\\.\\w+$'")}},
		},
		// in
		{
			in: User{
				ID:     "123456798132416549687946513215649687",
				Age:    19,
				Role:   "a",
				Number: 12,
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Role", Err: errors.New("value must be one of '[admin stuff]'")},
				ValidationError{Field: "Number", Err: errors.New("value must be one of '[10 11]'")},
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
			expectedErr: ValidationErrors{
				ValidationError{Field: "JobRoles", Err: errors.New("all elements in slice must have one value of '[admin stuff]'")},
				ValidationError{Field: "Numbers", Err: errors.New("all elements in slice must have one value of '[10 11]'")},
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
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: errors.New("value must be greater than '18'")},
				ValidationError{Field: "Ages", Err: errors.New("all elements in slice must have value greater than '18'")},
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
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: errors.New("value must be less than '25'")},
				ValidationError{Field: "Ages", Err: errors.New("all elements in slice must have value less than '25'")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
