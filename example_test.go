package gofp_test

import (
	"fmt"

	"github.com/tomasbasham/gofp"
)

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
