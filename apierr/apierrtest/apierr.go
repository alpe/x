package apierrtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/optiopay/x/apierr"
)

// alias that provides pretty printing
type APIValidationErrors []error

func (errs APIValidationErrors) String() string {
	var b bytes.Buffer
	for _, e := range errs {
		fmt.Fprintf(&b, "%+v\n", e)
	}
	return b.String()
}

// HasAPIErrors deserialize given response body into api errors and check if
// all expected errors were provided, comparing type, code and param
// attributes. Message attribute is ignored.
func HasAPIErrors(expected apierr.Errors, body io.Reader) APIValidationErrors {
	var resp struct {
		Errors apierr.Errors `json:"errors"`
	}
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return []error{fmt.Errorf("cannot decode response: %s", err)}
	}

	var errs []error

	// special case when we expect no errors
	if expected == nil {
		for _, got := range resp.Errors {
			errs = append(errs, fmt.Errorf("unexpected error: %#v", got))
		}
		return errs
	}

	for _, ex := range expected {
		has := false
		for _, got := range resp.Errors {
			if ex.Param == got.Param {
				errs = append(errs, compareErrors(ex, got)...)
				has = true
				break
			}
		}

		if !has {
			errs = append(errs, fmt.Errorf("missing error for %q field", ex.Param))
		}
	}

	// error on all unexpected errors
	for _, got := range resp.Errors {
		has := false
		for _, ex := range expected {
			if ex.Param == got.Param {
				has = true
				break
			}
		}

		if !has {
			errs = append(errs, fmt.Errorf("unexpected error: %#v", got))
		}
	}

	return errs
}

func compareErrors(expected, got apierr.Error) []error {
	var errs []error
	if expected.Type != got.Type {
		e := fmt.Errorf("expected Type %q, got %q", expected.Type, got.Type)
		errs = append(errs, e)
	}
	if expected.Code != got.Code {
		e := fmt.Errorf("expected Code %q, got %q", expected.Code, got.Code)
		errs = append(errs, e)
	}
	return errs
}
