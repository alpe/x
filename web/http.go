package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/optiopay/x/apierr"
	"github.com/optiopay/x/log"
)

// JSONResp send to given writer JSON encoded content. If passed content cannot
// be JSON serialized, 500 response with stantard JSON error message will be
// send instead.
//
// To make it easier for the client to understand the response, if not present
// in response body, "code" value is injected if possible. Code is the exact
// code send as part of HTTP response.
//
// Even though it's possible, NEVER send array as top level structure.
func JSONResp(w http.ResponseWriter, content interface{}, code int) {
	b, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		log.Error("cannot JSON serialize response",
			"content", fmt.Sprintf("%T", content),
			"error", err.Error())
		code = http.StatusInternalServerError
		var resp = struct {
			Errors apierr.Errors `json:"errors"`
		}{
			Errors: apierr.Errors{
				{
					Type:    "server_error",
					Message: http.StatusText(http.StatusInternalServerError),
				},
			},
		}

		b, err = json.MarshalIndent(resp, "", "\t")
		if err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	_, _ = w.Write(b)
}

// JSONErr send to given writer JSON encoded response, using always the same
// format for the message, as described in
// https://github.com/husio/documentation/blob/master/api/source/includes/_errors.md
func JSONErr(w http.ResponseWriter, errs apierr.Errors, code int) {
	var resp = struct {
		Errors apierr.Errors `json:"errors"`
	}{
		Errors: errs,
	}
	JSONResp(w, resp, code)
}

// StdJSONErr write standard HTTP response code and message for given status
// code. Content is serialized with JSON and formatted as described in format
// for the message, as described in
// https://github.com/husio/documentation/blob/master/api/source/includes/_errors.md
//
// Only 4xx and 5xx are valid codes for this function. Call with any other call
// will produce incomplete message.
func StdJSONErr(w http.ResponseWriter, code int) {
	err := apierr.Error{
		Message: http.StatusText(code),
	}

	switch code / 100 {
	case 4:
		err.Type = "requests_error"
	case 5:
		err.Type = "server_error"
	}

	// some error codes can provide even more information
	switch code {
	case http.StatusNotFound:
		err.Code = "not_found"
	case http.StatusForbidden:
		err.Code = "forbidden"
	}

	JSONErr(w, apierr.Errors{}.With(err), code)
}
