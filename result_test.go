package gofp_test

import (
	"errors"
	"testing"

	"github.com/tomasbasham/gofp"
)

func TestOk(t *testing.T) {
	r := gofp.Ok("test")
	if !r.IsOk() || r.IsErr() {
		t.Error("expected Ok")
	}
	if r.Unwrap() != "test" {
		t.Error("expected test")
	}
}

func TestErr(t *testing.T) {
	expectedErr := errors.New("test error")
	r := gofp.Err[string](expectedErr)
	if r.IsOk() || !r.IsErr() {
		t.Error("expected Err")
	}
	if r.UnwrapErr() != expectedErr {
		t.Error("expected test error")
	}
}

func TestFromReturn(t *testing.T) {
	t.Run("returns Ok for non-error", func(t *testing.T) {
		r := gofp.FromReturn("test", nil)
		if !r.IsOk() || r.IsErr() {
			t.Error("expected Ok")
		}
		if r.Unwrap() != "test" {
			t.Error("expected test")
		}
	})

	t.Run("returns Err for error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		r := gofp.FromReturn[any](nil, expectedErr)
		if !r.IsErr() || r.IsOk() {
			t.Error("expected Err")
		}
		if r.UnwrapErr() != expectedErr {
			t.Error("expected test error")
		}
	})
}

func TestResultMap(t *testing.T) {
	t.Run("maps Ok value", func(t *testing.T) {
		r := gofp.Ok("test")
		fn := func(s string) int {
			return len(s)
		}
		got := gofp.ResultMap(r, fn)
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != 4 {
			t.Error("expected 4")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		r := gofp.Err[string](expectedErr)
		fn := func(s string) int {
			return len(s)
		}
		got := gofp.ResultMap(r, fn)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected error to propagate")
		}
	})
}

func TestResult_Map(t *testing.T) {
	t.Run("maps Ok value", func(t *testing.T) {
		r := gofp.Ok("test")
		fn := func(s string) string {
			return s + "_processed"
		}
		got := r.Map(fn)
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != "test_processed" {
			t.Error("expected test_processed")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		r := gofp.Err[string](expectedErr)
		fn := func(s string) string {
			return s + "_processed"
		}
		got := r.Map(fn)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected error to propagate")
		}
	})
}

func TestResultApply(t *testing.T) {
	t.Run("applies Ok function to Ok value", func(t *testing.T) {
		r := gofp.Ok("test")
		fn := gofp.Ok(func(s string) int {
			return len(s)
		})
		got := gofp.ResultApply(r, fn)
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != 4 {
			t.Error("expected 4")
		}
	})

	t.Run("propagates error from value", func(t *testing.T) {
		expectedErr := errors.New("test error")
		r := gofp.Err[string](expectedErr)
		fn := gofp.Ok(func(s string) int {
			return len(s)
		})
		got := gofp.ResultApply(r, fn)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected error to propagate")
		}
	})

	t.Run("propagates error from function", func(t *testing.T) {
		expectedErr := errors.New("function error")
		r := gofp.Ok("test")
		fn := gofp.Err[func(string) int](expectedErr)
		got := gofp.ResultApply(r, fn)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected error to propagate")
		}
	})
}

func TestResultFlatMap(t *testing.T) {
	t.Run("flat maps Ok value to Ok", func(t *testing.T) {
		r := gofp.Ok("test")
		fn := func(s string) gofp.Result[int] {
			return gofp.Ok(len(s))
		}
		got := gofp.ResultFlatMap(r, fn)
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != 4 {
			t.Error("expected 4")
		}
	})

	t.Run("flat maps Ok value to Err", func(t *testing.T) {
		expectedErr := errors.New("mapping error")
		r := gofp.Ok("test")
		fn := func(s string) gofp.Result[int] {
			return gofp.Err[int](expectedErr)
		}
		got := gofp.ResultFlatMap(r, fn)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected mapping error")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		expectedErr := errors.New("initial error")
		r := gofp.Err[string](expectedErr)
		fn := func(s string) gofp.Result[int] {
			return gofp.Ok(len(s))
		}
		got := gofp.ResultFlatMap(r, fn)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected error to propagate")
		}
	})
}

func TestResult_FlatMap(t *testing.T) {
	t.Run("flat maps Ok value to Ok", func(t *testing.T) {
		r := gofp.Ok("test")
		fn := func(s string) gofp.Result[string] {
			return gofp.Ok(s + "_processed")
		}
		got := r.FlatMap(fn)
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != "test_processed" {
			t.Error("expected test_processed")
		}
	})

	t.Run("flat maps Ok value to Err", func(t *testing.T) {
		expectedErr := errors.New("mapping error")
		r := gofp.Ok("test")
		fn := func(s string) gofp.Result[string] {
			return gofp.Err[string](expectedErr)
		}
		got := r.FlatMap(fn)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected mapping error")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		expectedErr := errors.New("initial error")
		r := gofp.Err[string](expectedErr)
		fn := func(s string) gofp.Result[string] {
			return gofp.Ok(s + "_processed")
		}
		got := r.FlatMap(fn)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected error to propagate")
		}
	})
}

func TestResult_String(t *testing.T) {
	t.Run("formats Ok value", func(t *testing.T) {
		r := gofp.Ok("test")
		if got := r.String(); got != "Ok(test)" {
			t.Error("expected Ok(test)")
		}
	})

	t.Run("formats Err value", func(t *testing.T) {
		expectedErr := errors.New("test error")
		r := gofp.Err[string](expectedErr)
		if r.String() != "Err(test error)" {
			t.Error("expected Err(test error)")
		}
	})
}

func TestResult_UnwrapOr(t *testing.T) {
	t.Run("returns value for Ok", func(t *testing.T) {
		r := gofp.Ok("test")
		if got := r.UnwrapOr("default"); got != "test" {
			t.Error("expected test")
		}
	})

	t.Run("returns default for Err", func(t *testing.T) {
		r := gofp.Err[string](errors.New("error"))
		if got := r.UnwrapOr("default"); got != "default" {
			t.Error("expected default")
		}
	})
}

func TestResult_UnwrapOrElse(t *testing.T) {
	t.Run("returns value for Ok", func(t *testing.T) {
		r := gofp.Ok("test")
		if got := r.UnwrapOrElse(func() string { return "default" }); got != "test" {
			t.Error("expected test")
		}
	})

	t.Run("returns function result for Err", func(t *testing.T) {
		r := gofp.Err[string](errors.New("error"))
		if got := r.UnwrapOrElse(func() string { return "computed" }); got != "computed" {
			t.Error("expected computed")
		}
	})
}

func TestResult_And(t *testing.T) {
	t.Run("returns second result when first is Ok", func(t *testing.T) {
		r1 := gofp.Ok("first")
		r2 := gofp.Ok("second")
		got := r1.And(r2)
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != "second" {
			t.Error("expected second value")
		}
	})

	t.Run("returns error when first result is Err", func(t *testing.T) {
		expectedErr := errors.New("first error")
		r1 := gofp.Err[string](expectedErr)
		r2 := gofp.Ok("second")
		got := r1.And(r2)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected first error")
		}
	})

	t.Run("returns second result when it is Err", func(t *testing.T) {
		r1 := gofp.Ok("first")
		expectedErr := errors.New("second error")
		r2 := gofp.Err[string](expectedErr)
		got := r1.And(r2)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected second error")
		}
	})
}

func TestResult_AndThen(t *testing.T) {
	t.Run("applies function when Ok", func(t *testing.T) {
		r := gofp.Ok("test")
		fn := func(s string) gofp.Result[string] {
			return gofp.Ok(s + "_processed")
		}
		got := r.AndThen(fn)
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != "test_processed" {
			t.Error("expected processed value")
		}
	})

	t.Run("propagates initial error", func(t *testing.T) {
		expectedErr := errors.New("initial error")
		r := gofp.Err[string](expectedErr)
		fn := func(s string) gofp.Result[string] {
			return gofp.Ok(s + "_processed")
		}
		got := r.AndThen(fn)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected initial error")
		}
	})

	t.Run("propagates function error", func(t *testing.T) {
		expectedErr := errors.New("processing error")
		r := gofp.Ok("test")
		fn := func(s string) gofp.Result[string] {
			return gofp.Err[string](expectedErr)
		}
		got := r.AndThen(fn)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected processing error")
		}
	})
}

func TestResult_Or(t *testing.T) {
	t.Run("returns first result when Ok", func(t *testing.T) {
		r1 := gofp.Ok("first")
		r2 := gofp.Ok("second")
		got := r1.Or(r2)
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != "first" {
			t.Error("expected first value")
		}
	})

	t.Run("returns second result when first is Err", func(t *testing.T) {
		r1 := gofp.Err[string](errors.New("first error"))
		r2 := gofp.Ok("second")
		got := r1.Or(r2)
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != "second" {
			t.Error("expected second value")
		}
	})

	t.Run("returns second error when both are Err", func(t *testing.T) {
		r1 := gofp.Err[string](errors.New("first error"))
		expectedErr := errors.New("second error")
		r2 := gofp.Err[string](expectedErr)
		got := r1.Or(r2)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected second error")
		}
	})
}

func TestResult_OrElse(t *testing.T) {
	t.Run("returns original result when Ok", func(t *testing.T) {
		r := gofp.Ok("test")
		fn := func(error) gofp.Result[string] {
			return gofp.Ok("fallback")
		}
		got := r.OrElse(fn)
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != "test" {
			t.Error("expected original value")
		}
	})

	t.Run("returns function result when Err", func(t *testing.T) {
		r := gofp.Err[string](errors.New("initial error"))
		fn := func(error) gofp.Result[string] {
			return gofp.Ok("fallback")
		}
		got := r.OrElse(fn)
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != "fallback" {
			t.Error("expected fallback value")
		}
	})

	t.Run("returns function error when both fail", func(t *testing.T) {
		r := gofp.Err[string](errors.New("initial error"))
		expectedErr := errors.New("fallback error")
		fn := func(error) gofp.Result[string] {
			return gofp.Err[string](expectedErr)
		}
		got := r.OrElse(fn)
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected fallback error")
		}
	})
}

func TestResult_Ensure(t *testing.T) {
	t.Run("keeps value when predicate is true", func(t *testing.T) {
		r := gofp.Ok(5)
		unexpectedErr := errors.New("predicate error")
		got := r.Ensure(unexpectedErr, func(i int) bool { return i > 0 })
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != 5 {
			t.Error("expected 5")
		}
	})

	t.Run("returns error when predicate is false", func(t *testing.T) {
		r := gofp.Ok(5)
		expectedErr := errors.New("predicate error")
		got := r.Ensure(expectedErr, func(i int) bool { return i < 0 })
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected predicate error")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		unexpectedErr := errors.New("predicate error")
		r := gofp.Err[int](expectedErr)
		got := r.Ensure(unexpectedErr, func(i int) bool { return i > 0 })
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected error to propagate")
		}
	})
}

func TestResult_ToReturn(t *testing.T) {
	t.Run("returns value for Ok", func(t *testing.T) {
		r := gofp.Ok("test")
		got, err := r.ToReturn()
		if got != "test" {
			t.Error("expected test")
		}
		if err != nil {
			t.Error("expected no error")
		}
	})

	t.Run("returns error for Err", func(t *testing.T) {
		expectedErr := errors.New("test error")
		r := gofp.Err[string](expectedErr)
		got, err := r.ToReturn()
		if got != "" {
			t.Error("expected empty string")
		}
		if err != expectedErr {
			t.Error("expected test error")
		}
	})
}

func TestResult_Recover(t *testing.T) {
	t.Run("returns Ok for non-error", func(t *testing.T) {
		r := gofp.Ok("test")
		got := r.Recover(func(error) string { return "recovered" })
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != "test" {
			t.Error("expected test")
		}
	})

	t.Run("returns Ok for error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		r := gofp.Err[string](expectedErr)
		got := r.Recover(func(error) string { return "recovered" })
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != "recovered" {
			t.Error("expected recovered")
		}
	})
}

func TestResult_RecoverWith(t *testing.T) {
	t.Run("returns Ok for non-error", func(t *testing.T) {
		r := gofp.Ok("test")
		got := r.RecoverWith(func(error) gofp.Result[string] { return gofp.Ok("recovered") })
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != "test" {
			t.Error("expected test")
		}
	})

	t.Run("returns Ok for error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		r := gofp.Err[string](expectedErr)
		got := r.RecoverWith(func(error) gofp.Result[string] { return gofp.Ok("recovered") })
		if !got.IsOk() || got.IsErr() {
			t.Error("expected Ok")
		}
		if got.Unwrap() != "recovered" {
			t.Error("expected recovered")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		r := gofp.Err[string](expectedErr)
		got := r.RecoverWith(func(error) gofp.Result[string] { return gofp.Err[string](expectedErr) })
		if !got.IsErr() || got.IsOk() {
			t.Error("expected Err")
		}
		if got.UnwrapErr() != expectedErr {
			t.Error("expected test error")
		}
	})
}
