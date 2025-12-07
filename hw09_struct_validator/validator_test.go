package hw09structvalidator

import (
	"encoding/json"
	"fmt"
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
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	TestoStruct struct {
		Number int    `json:"sosos" validate:"min:8|max:11"`
		Str    string `validate:"in:some,boam,tam|len:3"`
	}

	TripleValidationStruct struct {
		Data string `validate:"len:5|in:hello,world,lesgo|regexp:^\\w+$"`
	}

	TestForIntSlice struct {
		Number []int `validate:"max:8|in:1,2,3"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		// Valid cases
		{
			in:          App{Version: "12345"},
			expectedErr: nil,
		},
		{
			in:          TestoStruct{Number: 10, Str: "tam"},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     string(make([]byte, 36)),
				Name:   "John",
				Age:    25,
				Email:  "user@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			in:          Response{Code: 200},
			expectedErr: nil,
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in: TripleValidationStruct{
				Data: "hello",
			},
			expectedErr: nil,
		},
		// Invalid cases
		{
			in: App{Version: "1234"},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: ErrLenValidation},
			},
		},
		{
			in: User{Age: 15},
			expectedErr: ValidationErrors{ // Errors include validation failures due to zero init of struct fields
				{Field: "ID", Err: ErrLenValidation},
				{Field: "Age", Err: ErrMinValidation},
				{Field: "Email", Err: ErrRegexpValidation},
				{Field: "UserRole", Err: ErrInValidation},
			},
		},
		{
			in: TestoStruct{Number: 8, Str: "wrong"},
			expectedErr: ValidationErrors{
				{Field: "Str", Err: ErrInValidation},
				{Field: "Str", Err: ErrLenValidation},
			},
		},
		{
			in: TestoStruct{
				Number: 5,
				Str:    "no",
			},
			expectedErr: ValidationErrors{
				{Field: "Number", Err: ErrMinValidation},
				{Field: "Str", Err: ErrInValidation},
				{Field: "Str", Err: ErrLenValidation},
			},
		},
		{
			in: User{
				ID: "123456123456123456123456123456123456", Age: 21, Email: "bad email", Role: "admin",
				Phones: []string{
					"22222222221", "22222222222", "22222222223", "22222222224", "22222222225",
					"22222222226", "72222222222", "8", "9", "10", "11",
				},
			},
			expectedErr: ValidationErrors{
				{Field: "Email", Err: ErrRegexpValidation},
				{Field: "Phones[7]", Err: ErrLenValidation},
				{Field: "Phones[8]", Err: ErrLenValidation},
				{Field: "Phones[9]", Err: ErrLenValidation},
				{Field: "Phones[10]", Err: ErrLenValidation},
			},
		},
		{
			in: TripleValidationStruct{
				Data: "hell",
			},
			expectedErr: ValidationErrors{
				{Field: "Data", Err: ErrLenValidation},
				{Field: "Data", Err: ErrInValidation},
			},
		},
		{
			in: TripleValidationStruct{
				Data: "hello!",
			},
			expectedErr: ValidationErrors{
				{Field: "Data", Err: ErrLenValidation},
				{Field: "Data", Err: ErrInValidation},
				{Field: "Data", Err: ErrRegexpValidation},
			},
		},
		{
			in: TripleValidationStruct{
				Data: "les_go",
			},
			expectedErr: ValidationErrors{
				{Field: "Data", Err: ErrLenValidation},
				{Field: "Data", Err: ErrInValidation},
			},
		},
		{
			in: TestForIntSlice{
				Number: []int{1, 2, 3, 4, 8, 9},
			},
			expectedErr: ValidationErrors{
				{Field: "Number[3]", Err: ErrInValidation},
				{Field: "Number[4]", Err: ErrInValidation},
				{Field: "Number[5]", Err: ErrMaxValidation},
				{Field: "Number[5]", Err: ErrInValidation},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("Validate(%T) error = %v", tt.in, err)
				}
				return
			}

			if err == nil {
				t.Errorf("Validate(%T) error = nil, want error", tt.in)
				return
			}

			if got := err.Error(); got != tt.expectedErr.Error() {
				t.Errorf("Validate(%T) error %q, want %q", tt.in, got, tt.expectedErr.Error())
			}
		})
	}
}
