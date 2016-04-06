package apierr

import (
	"encoding/json"
	"fmt"
	"io"
)

// APIErr represent single API error as defined in
// https://github.com/optiopay/documentation/blob/master/api/source/includes/_errors.md
type Error struct {
	Type    string `json:"type"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Param   string `json:"param,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("API error: %s, %s", e.Type, e.Code)
}

// APIErrors defines collection of APIErr structures. Thanks to this alias,
// defining errors declaratevly is way nicer!
//
//     var errs Errors = {
//        {"request_error", "not_found", "entity_not_found", ""},
//     }
//
// Adding more errors can be done by either passing Error instance or using
// helper methods:
//
//    var errs Errors
//    errs = errs.WithErr(Error{Type: "validation_error", Code: "not_integer"})
//    errs = errs.WithNotList("users", "\"users\" has to be list")
//
//    // create new errors list
//    errs = Errors{}.WithNotFound("")
//
//    // print JSON representation
//    errs.WriteTo(os.Stdout)
//
type Errors []Error

// WriteTo write JSON serialized Errors into given writer.
func (errs Errors) WriteTo(w io.Writer) (int64, error) {
	b, err := json.MarshalIndent(errs, "", "\t")
	if err != nil {
		return 0, err
	}
	n, err := w.Write(b)
	return int64(n), err
}

func (errs Errors) With(e Error) Errors {
	return append(errs, e)
}

func (errs Errors) WithRateLimit(message string) Errors {
	return append(errs, Error{
		Type:    "request_error",
		Code:    "rate_limit",
		Message: message,
	})
}

func (errs Errors) WithUnauthorized(message string) Errors {
	return append(errs, Error{
		Type:    "request_error",
		Code:    "unauthorized",
		Message: message,
	})
}

func (errs Errors) WithForbidden(message string) Errors {
	return append(errs, Error{
		Type:    "request_error",
		Code:    "forbidden",
		Message: message,
	})
}

func (errs Errors) WithNotFound(message string) Errors {
	return append(errs, Error{
		Type:    "request_error",
		Code:    "not_found",
		Message: message,
	})
}

func (errs Errors) WithMalformedJSON(message string) Errors {
	return append(errs, Error{
		Type:    "request_error",
		Code:    "malformed_json",
		Message: message,
	})
}

func (errs Errors) WithRequired(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "required",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotString(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_string",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotBoolean(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_boolean",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotNumber(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_number",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotInteger(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_integer",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotList(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_list",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotObject(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_object",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotDate(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_date",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotTime(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_time",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotDatetime(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_datetime",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotEq(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_eq",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotGt(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_gt",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotGte(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_gte",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotLt(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_lt",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotLte(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_lte",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotSingleValue(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_single_value",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotAllowed(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_allowed",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithEmailNotAvailable(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "email_not_available",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotMoney(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_money",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotCurrency(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_currency",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithInvalidMoneyAmount(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "invalid_money_amount",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotDuration(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_duration",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNoSuchEntity(param, message string) Errors {
	if message == "" {
		message = "entity does not exist"
	}
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "no_such_entity",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNoSuchFild(param, message string) Errors {
	if message == "" {
		message = "unknown field"
	}
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "no_such_field",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithInvalidChoice(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "invalid_choice",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithInvalidPage(param, message string) Errors {
	if param == "" {
		param = "page"
	}
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "invalid_page",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithInvalidPageSize(param, message string) Errors {
	if param == "" {
		param = "pageSize"
	}
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "invalid_page_size",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotInRange(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_in_range",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithInvalidCreditCard(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "invalid_credit_card",
		Message: message,
		Param:   param,
	})
}

func (errs Errors) WithNotAllowedBIN(param, message string) Errors {
	return append(errs, Error{
		Type:    "validation_error",
		Code:    "not_allowed_BIN",
		Message: message,
		Param:   param,
	})
}
