package gofp_test

import (
	"testing"

	"github.com/tomasbasham/gofp"
)

func TestSome(t *testing.T) {
	o := gofp.Some("test")
	if !o.IsSome() || o.IsNone() {
		t.Error("expected Some")
	}
	if o.Unwrap() != "test" {
		t.Error("expected test")
	}
}

func TestNone(t *testing.T) {
	o := gofp.None[string]()
	if o.IsSome() || !o.IsNone() {
		t.Error("expected None")
	}
}

func TestOptionMap(t *testing.T) {
	t.Run("maps Some value", func(t *testing.T) {
		o := gofp.Some("test")
		fn := func(s string) int {
			return len(s)
		}
		got := gofp.OptionMap(o, fn)
		if !got.IsSome() || got.IsNone() {
			t.Error("expected Some")
		}
		if got.Unwrap() != 4 {
			t.Error("expected 4")
		}
	})

	t.Run("propagates None value", func(t *testing.T) {
		o := gofp.None[string]()
		fn := func(s string) int {
			return len(s)
		}
		got := gofp.OptionMap(o, fn)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})
}

func TestOption_Map(t *testing.T) {
	t.Run("maps Some value", func(t *testing.T) {
		o := gofp.Some("test")
		fn := func(s string) string {
			return s + "_processed"
		}
		got := o.Map(fn)
		if !got.IsSome() || got.IsNone() {
			t.Error("expected Some")
		}
		if got.Unwrap() != "test_processed" {
			t.Error("expected test_processed")
		}
	})

	t.Run("propagates None value", func(t *testing.T) {
		o := gofp.None[string]()
		fn := func(s string) string {
			return s + "_processed"
		}
		got := o.Map(fn)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})
}

func TestOptionApply(t *testing.T) {
	t.Run("applies Some function to Some value", func(t *testing.T) {
		o := gofp.Some("test")
		fn := gofp.Some(func(s string) int {
			return len(s)
		})
		got := gofp.OptionApply(o, fn)
		if !got.IsSome() || got.IsNone() {
			t.Error("expected Some")
		}
		if got.Unwrap() != 4 {
			t.Error("expected 4")
		}
	})

	t.Run("propagates None value from value", func(t *testing.T) {
		o := gofp.None[string]()
		fn := gofp.Some(func(s string) int {
			return len(s)
		})
		got := gofp.OptionApply(o, fn)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})

	t.Run("propagates None value from function", func(t *testing.T) {
		o := gofp.Some("test")
		fn := gofp.None[func(s string) int]()
		got := gofp.OptionApply(o, fn)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})
}

func TestOptionFlatMap(t *testing.T) {
	t.Run("flat maps Some value to Some", func(t *testing.T) {
		o := gofp.Some("test")
		fn := func(s string) gofp.Option[int] {
			return gofp.Some(len(s))
		}
		got := gofp.OptionFlatMap(o, fn)
		if !got.IsSome() || got.IsNone() {
			t.Error("expected Some")
		}
		if got.Unwrap() != 4 {
			t.Error("expected 4")
		}
	})

	t.Run("flat maps Some value to None", func(t *testing.T) {
		o := gofp.Some("test")
		fn := func(s string) gofp.Option[int] {
			return gofp.None[int]()
		}
		got := gofp.OptionFlatMap(o, fn)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})

	t.Run("propagates None value", func(t *testing.T) {
		o := gofp.None[string]()
		fn := func(s string) gofp.Option[int] {
			return gofp.Some(len(s))
		}
		got := gofp.OptionFlatMap(o, fn)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})
}

func TestOptionFlat_Map(t *testing.T) {
	t.Run("flat maps Some value to Some", func(t *testing.T) {
		o := gofp.Some("test")
		fn := func(s string) gofp.Option[string] {
			return gofp.Some(s + "_processed")
		}
		got := o.FlatMap(fn)
		if !got.IsSome() || got.IsNone() {
			t.Error("expected Some")
		}
		if got.Unwrap() != "test_processed" {
			t.Error("expected test_processed")
		}
	})

	t.Run("flat maps Some value to None", func(t *testing.T) {
		o := gofp.Some("test")
		fn := func(s string) gofp.Option[string] {
			return gofp.None[string]()
		}
		got := o.FlatMap(fn)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})

	t.Run("propagates None value", func(t *testing.T) {
		o := gofp.None[string]()
		fn := func(s string) gofp.Option[string] {
			return gofp.Some(s + "_processed")
		}
		got := o.FlatMap(fn)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})
}

func TestOption_String(t *testing.T) {
	t.Run("formats Some value", func(t *testing.T) {
		o := gofp.Some("test")
		if got := o.String(); got != "Some(test)" {
			t.Error("expected Some(test)")
		}
	})

	t.Run("formats None value", func(t *testing.T) {
		o := gofp.None[string]()
		if got := o.String(); got != "None" {
			t.Error("expected None")
		}
	})
}

func TestOption_UnwrapOr(t *testing.T) {
	t.Run("returns value for Some", func(t *testing.T) {
		o := gofp.Some("test")
		if got := o.UnwrapOr("default"); got != "test" {
			t.Error("expected test")
		}
	})

	t.Run("returns default for None", func(t *testing.T) {
		o := gofp.None[string]()
		if got := o.UnwrapOr("default"); got != "default" {
			t.Error("expected default")
		}
	})
}

func TestOption_UnwrapOrElse(t *testing.T) {
	t.Run("returns value for Some", func(t *testing.T) {
		o := gofp.Some("test")
		if got := o.UnwrapOrElse(func() string { return "default" }); got != "test" {
			t.Error("expected test")
		}
	})

	t.Run("returns function result for None", func(t *testing.T) {
		o := gofp.None[string]()
		if got := o.UnwrapOrElse(func() string { return "computed" }); got != "computed" {
			t.Error("expected computed")
		}
	})
}

func TestOption_And(t *testing.T) {
	t.Run("returns second option when first is Some", func(t *testing.T) {
		o1 := gofp.Some("first")
		o2 := gofp.Some("second")
		got := o1.And(o2)
		if !got.IsSome() || got.IsNone() {
			t.Error("expected Some")
		}
		if got.Unwrap() != "second" {
			t.Error("expected second value")
		}
	})

	t.Run("returns None when first option is None", func(t *testing.T) {
		o1 := gofp.None[string]()
		o2 := gofp.Some("second")
		got := o1.And(o2)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})

	t.Run("returns second option when it is None", func(t *testing.T) {
		o1 := gofp.Some("first")
		o2 := gofp.None[string]()
		got := o1.And(o2)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})
}

func TestOption_AndThen(t *testing.T) {
	t.Run("applies function when Some", func(t *testing.T) {
		o := gofp.Some("test")
		fn := func(s string) gofp.Option[string] {
			return gofp.Some(s + "_processed")
		}
		got := o.AndThen(fn)
		if !got.IsSome() || got.IsNone() {
			t.Error("expected Some")
		}
		if got.Unwrap() != "test_processed" {
			t.Error("expected test_processed")
		}
	})

	t.Run("propagates initial None value", func(t *testing.T) {
		o := gofp.None[string]()
		fn := func(s string) gofp.Option[string] {
			return gofp.Some(s + "_processed")
		}
		got := o.AndThen(fn)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})

	t.Run("propagates function None value", func(t *testing.T) {
		o := gofp.Some("test")
		fn := func(s string) gofp.Option[string] {
			return gofp.None[string]()
		}
		got := o.AndThen(fn)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})
}

func TestOption_Or(t *testing.T) {
	t.Run("returns first option when Some", func(t *testing.T) {
		o1 := gofp.Some("first")
		o2 := gofp.Some("second")
		got := o1.Or(o2)
		if !got.IsSome() || got.IsNone() {
			t.Error("expected Some")
		}
		if got.Unwrap() != "first" {
			t.Error("expected first value")
		}
	})

	t.Run("returns second option when first is None", func(t *testing.T) {
		o1 := gofp.None[string]()
		o2 := gofp.Some("second")
		got := o1.Or(o2)
		if !got.IsSome() || got.IsNone() {
			t.Error("expected Some")
		}
		if got.Unwrap() != "second" {
			t.Error("expected second value")
		}
	})

	t.Run("returns None when both are None", func(t *testing.T) {
		o1 := gofp.None[string]()
		o2 := gofp.None[string]()
		got := o1.Or(o2)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})
}

func TestOption_OrElse(t *testing.T) {
	t.Run("returns original option when Some", func(t *testing.T) {
		o := gofp.Some("test")
		fn := func() gofp.Option[string] {
			return gofp.Some("fallback")
		}
		got := o.OrElse(fn)
		if !got.IsSome() || got.IsNone() {
			t.Error("expected Some")
		}
		if got.Unwrap() != "test" {
			t.Error("expected original value")
		}
	})

	t.Run("returns function result when None", func(t *testing.T) {
		o := gofp.None[string]()
		fn := func() gofp.Option[string] {
			return gofp.Some("fallback")
		}
		got := o.OrElse(fn)
		if !got.IsSome() || got.IsNone() {
			t.Error("expected Some")
		}
		if got.Unwrap() != "fallback" {
			t.Error("expected fallback value")
		}
	})

	t.Run("returns None when both are None", func(t *testing.T) {
		o := gofp.None[string]()
		fn := func() gofp.Option[string] {
			return gofp.None[string]()
		}
		got := o.OrElse(fn)
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})
}

func TestOption_Filter(t *testing.T) {
	t.Run("keeps value when predicate is true", func(t *testing.T) {
		o := gofp.Some(5)
		got := o.Filter(func(i int) bool { return i > 0 })
		if !got.IsSome() || got.IsNone() {
			t.Error("expected Some")
		}
		if got.Unwrap() != 5 {
			t.Error("expected 5")
		}
	})

	t.Run("returns None when predicate is false", func(t *testing.T) {
		o := gofp.Some(5)
		got := o.Filter(func(i int) bool { return i < 0 })
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})

	t.Run("propagates None value", func(t *testing.T) {
		o := gofp.None[int]()
		got := o.Filter(func(i int) bool { return i > 0 })
		if !got.IsNone() || got.IsSome() {
			t.Error("expected None")
		}
	})
}

func TestOption_MarshalJSON(t *testing.T) {
	t.Run("marshals Some value", func(t *testing.T) {
		o := gofp.Some("test")
		got, err := o.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		if string(got) != `"test"` {
			t.Error("expected \"test\"")
		}
	})

	t.Run("marshals None value", func(t *testing.T) {
		o := gofp.None[string]()
		got, err := o.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		if got != nil {
			t.Error("expected nil")
		}
	})
}

func TestOption_UnmarshalJSON(t *testing.T) {
	t.Run("unmarshals Some value", func(t *testing.T) {
		var o gofp.Option[string]
		data := []byte(`"test"`)
		err := o.UnmarshalJSON(data)
		if err != nil {
			t.Error(err)
		}
		if o.Unwrap() != "test" {
			t.Error("expected test")
		}
	})

	t.Run("unmarshals None value", func(t *testing.T) {
		var o gofp.Option[string]
		var data []byte
		err := o.UnmarshalJSON(data)
		if err != nil {
			t.Error(err)
		}
		if o.IsSome() || !o.IsNone() {
			t.Error("expected None")
		}
	})
}
