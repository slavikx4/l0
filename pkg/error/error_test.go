package error

import (
	"errors"
	"testing"
)

type errorCase struct {
	inputError error
	inputOp    string
	expected   error
}

var cases = []errorCase{

	{&Error{Err: nil, Code: ErrorNotFound, Message: "", Op: ""},
		"*Handler.Get -> ",
		&Error{Err: nil, Code: ErrorNotFound, Message: "", Op: "*Handler.Get -> "},
	},

	{&Error{Err: errors.New(""), Code: ErrorNotFound, Message: "", Op: ""},
		"*Handler.Get -> ",
		&Error{Err: errors.New(""), Code: ErrorNotFound, Message: "", Op: "*Handler.Get -> "},
	},
	{&Error{Err: errors.New(""), Code: ErrorNotFound, Message: "", Op: "*Handler.Get -> "},
		"*Service.Get -> ",
		&Error{Err: errors.New(""), Code: ErrorNotFound, Message: "", Op: "*Handler.Get -> *Service.Get -> "},
	},
}

func errorEqual(t *testing.T) {
	for _, item := range cases {
		result := AddOp(item.inputError, item.inputOp).Error()
		expected := item.expected.Error()
		if result != expected {
			t.Errorf("expected: %v\nresult: %v\nInput: %v, %v",
				expected, result, item.inputError.Error(), item.inputOp)
		}
	}
}

func TestAddOp(t *testing.T) {
	t.Parallel()

	t.Run("errorCheckEqual", errorEqual)
}
