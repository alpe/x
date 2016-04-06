package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONResp(t *testing.T) {
	testcases := []struct {
		content  interface{}
		code     int
		expected map[string]interface{}
	}{
		{
			map[string]interface{}{"name": "john doe", "age": 123},
			http.StatusOK,
			map[string]interface{}{"name": "john doe", "age": 123, "statusCode": http.StatusOK},
		},
		{
			nil,
			http.StatusGone,
			nil,
		},
		{
			struct {
				A string
				B string
			}{"foo", "bar"},
			http.StatusCreated,
			map[string]interface{}{"A": "foo", "B": "bar", "statusCode": http.StatusCreated},
		},
	}

	for i, tc := range testcases {
		w := httptest.NewRecorder()
		JSONResp(w, tc.content, tc.code)

		b, err := json.MarshalIndent(tc.content, "", "\t")
		if err != nil {
			t.Fatalf("%d: cannot serialize %v: %s", i, tc.content, err)
		}

		if w.Header().Get("Content-Type") != "application/json; charset=UTF-8" {
			t.Fatalf("invalid content type: %v", w.Header()["Content-Type"])
		}

		if w.Code != tc.code {
			t.Fatalf("%d: expected %d, got %d", i, tc.code, w.Code)
		}
		if !bytes.Equal(w.Body.Bytes(), b) {
			t.Fatalf("%d: expected %q, got %q", i, b, w.Body.Bytes())
		}
	}
}
