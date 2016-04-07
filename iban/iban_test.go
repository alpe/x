package iban

import "testing"

func TestValidate(t *testing.T) {
	cases := map[string]bool{
		"DE44 # 5001 0517 5407 3249":        false,
		"DE44 5001 0517 5407 3249 231":      false,
		"de44 5001 0517 5407 3249 31":       true,
		"DE44 5001 0517 5407 3249 31":       true,
		"GR16 0110 1250 0000 0001 2300 695": true,
		"GB29 NWBK 6016 1331 9268 19":       true,
		"SA03 8000 0000 6080 1016 7519":     true,
		"CH93 0076 2011 6238 5295 7":        true,
	}
	for iban, ok := range cases {
		if err := Validate(iban); (err == nil) != ok {
			t.Errorf("expected validation of %q to be %v: %s", iban, ok, err)
		}
	}
}
