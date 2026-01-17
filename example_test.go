package gofp_test

import (
	"errors"
	"fmt"

	"github.com/tomasbasham/gofp"
)

func ExampleEither_Map() {
	e := gofp.Right[int]("hello")
	doubled := e.Map(func(s string) string { return s + s })
	value := doubled.Unwrap()
	fmt.Println(value)
	// Output:
	// hellohello
}

func ExampleEither_FlatMap() {
	e := gofp.Right[int]("hello")
	doubled := e.FlatMap(func(s string) gofp.Either[int, string] {
		return gofp.Right[int](s + s)
	})
	value := doubled.Unwrap()
	fmt.Println(value)
	// Output:
	// hellohello
}

func ExampleEither_FlatMapLeft() {
	e := gofp.Left[int, string](5)
	doubled := e.FlatMapLeft(func(x int) gofp.Either[int, string] {
		return gofp.Left[int, string](x * 2)
	})
	value := doubled.UnwrapLeft()
	fmt.Println(value)
	// Output:
	// 10
}

func ExampleLeft() {
	e := gofp.Left[int, string](5)
	value := e.UnwrapLeft()
	fmt.Println(value)
	// Output:
	// 5
}

func ExampleRight() {
	e := gofp.Right[int]("hello")
	value := e.Unwrap()
	fmt.Println(value)
	// Output:
	// hello
}

func ExampleFromResult() {
	r := gofp.Ok(5)
	e := gofp.FromResult(r)
	value := e.Unwrap()
	fmt.Println(value)
	// Output:
	// 5
}

func ExampleEitherMap() {
	e := gofp.Right[int]("hello")
	doubled := gofp.EitherMap(e, func(s string) string { return s + s })
	value := doubled.Unwrap()
	fmt.Println(value)
	// Output:
	// hellohello
}

func ExampleEitherMapLeft() {
	e := gofp.Left[int, string](5)
	doubled := gofp.EitherMapLeft(e, func(x int) int { return x * 2 })
	value := doubled.UnwrapLeft()
	fmt.Println(value)
	// Output:
	// 10
}

func ExampleEitherApply() {
	e := gofp.Right[int](5)
	double := gofp.Right[int](func(x int) int { return x * 2 })
	value := gofp.EitherApply(e, double).Unwrap()
	fmt.Println(value)
	// Output:
	// 10
}

func ExampleEitherApplyMap() {
	e1 := gofp.Left[string, int]("error1")
	e2 := gofp.Left[string, func(int) int]("error2")
	combined := gofp.EitherApplyMap(e1, e2, func(a, b string) string {
		return b + "; " + a
	})
	value := combined.UnwrapLeft()
	fmt.Println(value)
	// Output:
	// error1; error2
}

func ExampleEitherFlatMap() {
	e := gofp.Right[int]("hello")
	doubled := gofp.EitherFlatMap(e, func(s string) gofp.Either[int, string] {
		return gofp.Right[int](s + s)
	})
	value := doubled.Unwrap()
	fmt.Println(value)
	// Output:
	// hellohello
}

func ExampleEitherFlatMapLeft() {
	e := gofp.Left[int, string](5)
	parseInt := gofp.EitherFlatMapLeft(e, func(x int) gofp.Either[string, string] {
		return gofp.Left[string, string](fmt.Sprintf("%d", x))
	})
	value := parseInt.UnwrapLeft()
	fmt.Println(value)
	// Output:
	// 5
}

func ExampleEitherSequence() {
	eithers := []gofp.Either[string, int]{gofp.Right[string](1), gofp.Right[string](2), gofp.Right[string](3)}
	sequenced := gofp.EitherSequence(eithers)
	values := sequenced.Unwrap()
	fmt.Println(values)
	// Output:
	// [1 2 3]
}

func ExampleEitherFold() {
	e := gofp.Right[int]("hello")
	value := gofp.EitherFold(
		e,
		func(left int) string { return fmt.Sprintf("Left: %d", left) },
		func(right string) string { return "Right: " + right },
	)
	fmt.Println(value)
	// Output:
	// Right: hello
}

func ExampleOption_Map() {
	o := gofp.Some(5)
	doubled := o.Map(func(x int) int { return x * 2 })
	value := doubled.Unwrap()
	fmt.Println(value)
	// Output:
	// 10
}

func ExampleOption_FlatMap() {
	getName := gofp.Some("Alice")
	greet := getName.FlatMap(func(name string) gofp.Option[string] {
		return gofp.Some("Hello, " + name)
	})
	value := greet.Unwrap()
	fmt.Println(value)
	// Output:
	// Hello, Alice
}

func ExampleSome() {
	o := gofp.Some(5)
	value := o.Unwrap()
	fmt.Println(value)
	// Output:
	// 5
}

func ExampleNone() {
	o := gofp.None[int]()
	isNone := o.IsNone()
	fmt.Println(isNone)
	// Output: true
}

func ExampleOptionMap() {
	o := gofp.Some(5)
	parseInt := gofp.OptionMap(o, func(x int) string {
		return fmt.Sprintf("%d", x)
	})
	value := parseInt.Unwrap()
	fmt.Println(value)
	// Output:
	// 5
}

func ExampleOptionApply() {
	o := gofp.Some(5)
	double := gofp.Some(func(x int) int { return x * 2 })
	value := gofp.OptionApply(o, double).Unwrap()
	fmt.Println(value)
	// Output:
	// 10
}

func ExampleOptionFlatMap() {
	getName := gofp.Some("Alice")
	greet := gofp.OptionFlatMap(getName, func(name string) gofp.Option[string] {
		return gofp.Some("Hello, " + name)
	})
	value := greet.Unwrap()
	fmt.Println(value)
	// Output:
	// Hello, Alice
}

func ExampleOptionSequence() {
	options := []gofp.Option[int]{gofp.Some(1), gofp.Some(2), gofp.Some(3)}
	sequenced := gofp.OptionSequence(options)
	values := sequenced.Unwrap()
	fmt.Println(values)
	// Output:
	// [1 2 3]
}

func ExampleOptionFold() {
	o := gofp.Some(5)
	value := gofp.OptionFold(
		o,
		func() string { return "No value" },
		func(x int) string { return fmt.Sprintf("Value is %d", x) },
	)
	fmt.Println(value)
	// Output:
	// Value is 5
}

func ExampleOption_And() {
	o := gofp.Some(5)
	nextOpt := o.And(gofp.Some(10))
	value := nextOpt.Unwrap()
	fmt.Println(value)
	// Output:
	// 10
}

func ExampleOption_AndThen() {
	o := gofp.Some(5)
	incremented := o.AndThen(func(x int) gofp.Option[int] {
		return gofp.Some(x + 1)
	})
	value := incremented.Unwrap()
	fmt.Println(value)
	// Output:
	// 6
}

func ExampleOption_Or() {
	o := gofp.None[int]()
	defaultOpt := o.Or(gofp.Some(10))
	value := defaultOpt.Unwrap()
	fmt.Println(value)
	// Output:
	// 10
}

func ExampleOption_OrElse() {
	o := gofp.None[int]()
	defaultOpt := o.OrElse(func() gofp.Option[int] {
		return gofp.Some(10)
	})
	value := defaultOpt.Unwrap()
	fmt.Println(value)
	// Output:
	// 10
}

func ExampleOption_Filter() {
	o := gofp.Some(5)
	even := o.Filter(func(x int) bool { return x%2 == 0 })
	isNone := even.IsNone()
	fmt.Println(isNone)
	// Output:
	// true
}

func ExampleResult_Map() {
	r := gofp.Ok(5)
	doubled := r.Map(func(x int) int { return x * 2 })
	value := doubled.Unwrap()
	fmt.Println(value)
	// Output:
	// 10
}

func ExampleResult_FlatMap() {
	getName := gofp.Ok("Alice")
	greet := getName.FlatMap(func(name string) gofp.Result[string] {
		return gofp.Ok("Hello, " + name)
	})
	value := greet.Unwrap()
	fmt.Println(value)
	// Output:
	// Hello, Alice
}

func ExampleOk() {
	r := gofp.Ok(5)
	value := r.Unwrap()
	fmt.Println(value)
	// Output:
	// 5
}

func ExampleErr() {
	r := gofp.Err[int](fmt.Errorf("an error"))
	err := r.UnwrapErr()
	fmt.Println(r.IsErr())
	fmt.Println(err)
	// Output:
	// true
	// an error
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func ExampleFromReturn() {
	r := gofp.FromReturn(Divide(10, 2))
	value := r.Unwrap()
	fmt.Println(value)
	// Output:
	// 5
}

func DivideResult(a, b int) gofp.Result[int] {
	if b == 0 {
		return gofp.Err[int](fmt.Errorf("division by zero"))
	}
	return gofp.Ok(a / b)
}

func ExampleResult_ToReturn() {
	value, err := DivideResult(10, 2).ToReturn()
	fmt.Println(value, err)
	// Output:
	// 5 <nil>
}

func ExampleResultMap() {
	r := gofp.Ok(5)
	parseInt := gofp.ResultMap(r, func(x int) string {
		return fmt.Sprintf("%d", x)
	})
	value := parseInt.Unwrap()
	fmt.Println(value)
	// Output:
	// 5
}

func ExampleResultApply() {
	r := gofp.Ok(5)
	double := gofp.Ok(func(x int) int { return x * 2 })
	value := gofp.ResultApply(r, double).Unwrap()
	fmt.Println(value)
	// Output:
	// 10
}

func ExampleResultFlatMap() {
	getName := gofp.Ok("Alice")
	greet := gofp.ResultFlatMap(getName, func(name string) gofp.Result[string] {
		return gofp.Ok("Hello, " + name)
	})
	value := greet.Unwrap()
	fmt.Println(value)
	// Output:
	// Hello, Alice
}

func ExampleResultSequence() {
	results := []gofp.Result[int]{gofp.Ok(1), gofp.Ok(2), gofp.Ok(3)}
	sequenced := gofp.ResultSequence(results)
	values := sequenced.Unwrap()
	fmt.Println(values)
	// Output:
	// [1 2 3]
}

func ExampleResultFold() {
	r := gofp.Ok(5)
	value := gofp.ResultFold(
		r,
		func(err error) string { return "error occurred" },
		func(v int) string { return fmt.Sprintf("value is %d", v) },
	)
	fmt.Println(value)
	// Output:
	// value is 5
}

func ExampleResult_Ensure() {
	r := gofp.Ok(5)
	r2 := r.Ensure(
		errors.New("value is not greater than 10"),
		func(x int) bool { return x > 10 },
	)
	value := r2.UnwrapErr()
	fmt.Println(value)
	// Output:
	// value is not greater than 10
}

func ExampleResult_EnsureWith() {
	r := gofp.Ok(5)
	r2 := r.EnsureWith(
		func(x int) bool { return x > 10 },
		func(x int) error { return fmt.Errorf("%d is not greater than 10", x) },
	)
	value := r2.UnwrapErr()
	fmt.Println(value)
	// Output:
	// 5 is not greater than 10
}
