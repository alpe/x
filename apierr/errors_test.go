package apierr

import "testing"

func TestErrorsWithCreatesNewInstance(t *testing.T) {
	var errs Errors

	// do not assign result
	errs.WithRateLimit("")
	errs.WithUnauthorized("")

	if len(errs) != 0 {
		t.Fatalf("expected 0 errors, got %d", len(errs))
	}

	errs = errs.With(Error{
		Type: "t1",
		Code: "c1",
	})
	if len(errs) != 1 {
		t.Fatalf("expected 1 errors, got %d", len(errs))
	}

	errs = errs.With(Error{
		Type: "t2",
		Code: "c2",
	})
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}

	if errs[0].Type != "t1" || errs[1].Type != "t2" {
		t.Fatalf("unexpected errors order, got %#v", errs)
	}
}
