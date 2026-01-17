package gofp_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tomasbasham/gofp"
)

func TestLeft(t *testing.T) {
	e := gofp.Left[string, int]("test")
	if !e.IsLeft() || e.IsRight() {
		t.Error("expected Left")
	}
	if e.UnwrapLeft() != "test" {
		t.Error("expected test")
	}
}

func TestRight(t *testing.T) {
	e := gofp.Right[string](42)
	if e.IsLeft() || !e.IsRight() {
		t.Error("expected Right")
	}
	if e.Unwrap() != 42 {
		t.Error("expected 42")
	}
}

func TestFromResult(t *testing.T) {
	t.Run("returns Right for Ok result", func(t *testing.T) {
		r := gofp.Ok("test")
		e := gofp.FromResult(r)
		if !e.IsRight() || e.IsLeft() {
			t.Error("expected Right")
		}
		if e.Unwrap() != "test" {
			t.Error("expected test")
		}
	})

	t.Run("returns Left for Err result", func(t *testing.T) {
		expectedErr := errors.New("test error")
		r := gofp.Err[string](expectedErr)
		e := gofp.FromResult(r)
		if !e.IsLeft() || e.IsRight() {
			t.Error("expected Left")
		}
		if e.UnwrapLeft() != expectedErr {
			t.Error("expected test error")
		}
	})
}

func TestEitherMap(t *testing.T) {
	t.Run("maps Right value", func(t *testing.T) {
		e := gofp.Right[string](21)
		fn := func(i int) string {
			return fmt.Sprintf("%d", i)
		}
		got := gofp.EitherMap(e, fn)
		if got.IsLeft() || !got.IsRight() {
			t.Error("expected Right")
		}
		if got.Unwrap() != "21" {
			t.Error("expected \"21\"")
		}
	})

	t.Run("preserves Left value", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		fn := func(i int) string {
			return fmt.Sprintf("%d", i)
		}
		got := gofp.EitherMap(e, fn)
		if !got.IsLeft() || got.IsRight() {
			t.Error("expected Left")
		}
		if got.UnwrapLeft() != "test" {
			t.Error("expected test")
		}
	})
}

func TestEither_Map(t *testing.T) {
	t.Run("maps Right value", func(t *testing.T) {
		e := gofp.Right[string](21)
		fn := func(i int) int {
			return i * 2
		}
		got := e.Map(fn)
		if got.IsLeft() || !got.IsRight() {
			t.Error("expected Right")
		}
		if got.Unwrap() != 42 {
			t.Error("expected 42")
		}
	})

	t.Run("preserves Left value", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		fn := func(i int) int {
			return i * 2
		}
		got := e.Map(fn)
		if !got.IsLeft() || got.IsRight() {
			t.Error("expected Left")
		}
		if got.UnwrapLeft() != "test" {
			t.Error("expected test")
		}
	})
}

func TestEitherMapLeft(t *testing.T) {
	t.Run("maps Left value", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		fn := func(s string) int {
			return len(s)
		}
		got := gofp.EitherMapLeft(e, fn)
		if !got.IsLeft() || got.IsRight() {
			t.Error("expected Left")
		}
		if got.UnwrapLeft() != 4 {
			t.Error("expected 4")
		}
	})

	t.Run("preserves Right value", func(t *testing.T) {
		e := gofp.Right[string](42)
		fn := func(s string) int {
			return len(s)
		}
		got := gofp.EitherMapLeft(e, fn)
		if got.IsLeft() || !got.IsRight() {
			t.Error("expected Right")
		}
		if got.Unwrap() != 42 {
			t.Error("expected 42")
		}
	})
}

func TestEither_MapLeft(t *testing.T) {
	t.Run("maps Left value", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		fn := func(s string) string {
			return s + "_processed"
		}
		got := e.MapLeft(fn)
		if !got.IsLeft() || got.IsRight() {
			t.Error("expected Left")
		}
		if got.UnwrapLeft() != "test_processed" {
			t.Error("expected test_processed")
		}
	})

	t.Run("preserves Right value", func(t *testing.T) {
		e := gofp.Right[string](42)
		fn := func(s string) string {
			return s + "_processed"
		}
		got := e.MapLeft(fn)
		if got.IsLeft() || !got.IsRight() {
			t.Error("expected Right")
		}
		if got.Unwrap() != 42 {
			t.Error("expected 42")
		}
	})
}

func TestEitherFlatMap(t *testing.T) {
	t.Run("flat maps Right value", func(t *testing.T) {
		e := gofp.Right[string](21)
		fn := func(i int) gofp.Either[string, string] {
			return gofp.Right[string](fmt.Sprintf("%d", i*2))
		}
		got := gofp.EitherFlatMap(e, fn)
		if got.IsLeft() || !got.IsRight() {
			t.Error("expected Right")
		}
		if got.Unwrap() != "42" {
			t.Error("expected \"42\"")
		}
	})

	t.Run("preserves Left value", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		fn := func(i int) gofp.Either[string, string] {
			return gofp.Right[string](fmt.Sprintf("%d", i*2))
		}
		got := gofp.EitherFlatMap(e, fn)
		if !got.IsLeft() || got.IsRight() {
			t.Error("expected Left")
		}
		if got.UnwrapLeft() != "test" {
			t.Error("expected test")
		}
	})
}

func TestEither_FlatMap(t *testing.T) {
	t.Run("flat maps Right value", func(t *testing.T) {
		e := gofp.Right[string](21)
		fn := func(i int) gofp.Either[string, int] {
			return gofp.Right[string](i * 2)
		}
		got := e.FlatMap(fn)
		if got.IsLeft() || !got.IsRight() {
			t.Error("expected Right")
		}
		if got.Unwrap() != 42 {
			t.Error("expected 42")
		}
	})

	t.Run("preserves Left value", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		fn := func(i int) gofp.Either[string, int] {
			return gofp.Right[string](i * 2)
		}
		got := e.FlatMap(fn)
		if !got.IsLeft() || got.IsRight() {
			t.Error("expected Left")
		}
		if got.UnwrapLeft() != "test" {
			t.Error("expected test")
		}
	})
}

func TestEitherFlatMapLeft(t *testing.T) {
	t.Run("flat maps Left value", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		fn := func(s string) gofp.Either[int, int] {
			return gofp.Left[int, int](len(s))
		}
		got := gofp.EitherFlatMapLeft(e, fn)
		if !got.IsLeft() || got.IsRight() {
			t.Error("expected Left")
		}
		if got.UnwrapLeft() != 4 {
			t.Error("expected 4")
		}
	})

	t.Run("preserves Right value", func(t *testing.T) {
		e := gofp.Right[string](42)
		fn := func(s string) gofp.Either[int, int] {
			return gofp.Left[int, int](len(s))
		}
		got := gofp.EitherFlatMapLeft(e, fn)
		if got.IsLeft() || !got.IsRight() {
			t.Error("expected Right")
		}
		if got.Unwrap() != 42 {
			t.Error("expected 42")
		}
	})
}

func TestEither_FlatMapLeft(t *testing.T) {
	t.Run("flat maps Left value", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		fn := func(s string) gofp.Either[string, int] {
			return gofp.Left[string, int](s + "_processed")
		}
		got := e.FlatMapLeft(fn)
		if !got.IsLeft() || got.IsRight() {
			t.Error("expected Left")
		}
		if got.UnwrapLeft() != "test_processed" {
			t.Error("expected test_processed")
		}
	})

	t.Run("preserves Right value", func(t *testing.T) {
		e := gofp.Right[string](42)
		fn := func(s string) gofp.Either[string, int] {
			return gofp.Left[string, int](s + "_processed")
		}
		got := e.FlatMapLeft(fn)
		if got.IsLeft() || !got.IsRight() {
			t.Error("expected Right")
		}
		if got.Unwrap() != 42 {
			t.Error("expected 42")
		}
	})
}

func TestEither_String(t *testing.T) {
	t.Run("formats Left value", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		if e.String() != "Left(test)" {
			t.Error("expected Left(test)")
		}
	})

	t.Run("formats Right value", func(t *testing.T) {
		e := gofp.Right[string](42)
		if e.String() != "Right(42)" {
			t.Error("expected Right(42)")
		}
	})
}

func TestEither_UnwrapOr(t *testing.T) {
	t.Run("returns Right value when Right", func(t *testing.T) {
		e := gofp.Right[string](42)
		got := e.UnwrapOr(0)
		if got != 42 {
			t.Error("expected 42")
		}
	})

	t.Run("returns default value when Left", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		got := e.UnwrapOr(0)
		if got != 0 {
			t.Error("expected 0")
		}
	})
}

func TestEither_UnwrapOrElse(t *testing.T) {
	t.Run("returns Right value when Right", func(t *testing.T) {
		e := gofp.Right[string](42)
		got := e.UnwrapOrElse(func() int { return 0 })
		if got != 42 {
			t.Error("expected 42")
		}
	})

	t.Run("returns function result when Left", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		got := e.UnwrapOrElse(func() int { return 99 })
		if got != 99 {
			t.Error("expected 99")
		}
	})
}

func TestEither_UnwrapLeftOr(t *testing.T) {
	t.Run("returns Left value when Left", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		got := e.UnwrapLeftOr("default")
		if got != "test" {
			t.Error("expected test")
		}
	})

	t.Run("returns default value when Right", func(t *testing.T) {
		e := gofp.Right[string](42)
		got := e.UnwrapLeftOr("default")
		if got != "default" {
			t.Error("expected default")
		}
	})
}

func TestEither_UnwrapLeftOrElse(t *testing.T) {
	t.Run("returns Left value when Left", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		got := e.UnwrapLeftOrElse(func() string { return "computed" })
		if got != "test" {
			t.Error("expected test")
		}
	})

	t.Run("returns function result when Right", func(t *testing.T) {
		e := gofp.Right[string](42)
		got := e.UnwrapLeftOrElse(func() string { return "computed" })
		if got != "computed" {
			t.Error("expected computed")
		}
	})
}

func TestEither_Swap(t *testing.T) {
	t.Run("swaps Left to Right", func(t *testing.T) {
		e := gofp.Left[string, int]("test")
		got := e.Swap()
		if got.IsLeft() || !got.IsRight() {
			t.Error("expected Right")
		}
		if got.Unwrap() != "test" {
			t.Error("expected test")
		}
	})

	t.Run("swaps Right to Left", func(t *testing.T) {
		e := gofp.Right[string](42)
		got := e.Swap()
		if !got.IsLeft() || got.IsRight() {
			t.Error("expected Left")
		}
		if got.UnwrapLeft() != 42 {
			t.Error("expected 42")
		}
	})
}
