package reader_test

import (
	"fmt"

	"github.com/tomasbasham/gofp/reader"
)

func ExampleReader_Map() {
	env := Environment{}
	r := reader.Pure[Environment](5)
	doubled := r.Map(func(x int) int { return x * 2 })
	value := doubled.Run(env)
	fmt.Println(value)
	// Output:
	// 10
}

func ExampleReader_FlatMap() {
	getName := reader.Map(reader.Ask[Environment](), func(env Environment) string {
		return env.Name
	})

	greet := getName.FlatMap(func(name string) reader.Reader[Environment, string] {
		return reader.Pure[Environment]("Hello, " + name)
	})

	env := Environment{Name: "Alice"}
	value := greet.Run(env)
	fmt.Println(value)
	// Output:
	// Hello, Alice
}

func ExampleReader_Run() {
	env := Environment{}
	r := reader.Pure[Environment]("test")
	value := r.Run(env)
	fmt.Println(value)
	// Output:
	// test
}

func ExamplePure() {
	env := Environment{}
	r := reader.Pure[Environment](5)
	value := r.Run(env)
	fmt.Println(value)
	// Output:
	// 5
}

func IsProd() reader.Reader[Environment, bool] {
	return reader.New(func(env Environment) bool {
		return env.Name == "production"
	})
}

func ExampleNew() {
	env := Environment{Name: "production"}
	value := IsProd().Run(env)
	fmt.Println(value)
	// Output:
	// true
}

func AskIsProd() reader.Reader[Environment, bool] {
	return reader.Map(reader.Ask[Environment](), func(env Environment) bool {
		return env.Name == "production"
	})
}

func ExampleAsk() {
	env := Environment{Name: "production"}
	value := AskIsProd().Run(env)
	fmt.Println(value)
	// Output:
	// true
}

func RunInProd(r reader.Reader[Environment, Environment]) reader.Reader[Environment, Environment] {
	return reader.Local(r, func(env Environment) Environment {
		env.Name = "production"
		env.Value = 10
		return env
	})
}

func ExampleLocal() {
	env := Environment{Name: "test"}
	value := RunInProd(reader.Ask[Environment]()).Run(env)
	fmt.Println(env)
	fmt.Println(value)
	// Output:
	// {false test 0}
	// {false production 10}
}

func ExampleMap() {
	env := Environment{}
	r := reader.Pure[Environment](5)
	parseInt := reader.Map(r, func(x int) string { return fmt.Sprintf("%d", x) })
	value := parseInt.Run(env)
	fmt.Println(value)
	// Output:
	// 5
}

func ExampleApply() {
	env := Environment{}
	r := reader.Pure[Environment](5)
	double := reader.Pure[Environment](func(x int) int { return x * 2 })
	value := reader.Apply(r, double).Run(env)
	fmt.Println(value)
	// Output:
	// 10
}

func ExampleFlatMap() {
	getName := reader.Map(reader.Ask[Environment](), func(env Environment) string {
		return env.Name
	})

	greet := reader.FlatMap(getName, func(name string) reader.Reader[Environment, string] {
		return reader.Pure[Environment]("Hello, " + name)
	})

	env := Environment{Name: "Alice"}
	value := greet.Run(env)
	fmt.Println(value)
	// Output:
	// Hello, Alice
}

func ExampleZip() {
	env := Environment{}
	r1 := reader.Pure[Environment](5)
	r2 := reader.Pure[Environment](3)
	sum := reader.Zip(r1, r2, func(a, b int) int { return a + b })
	value := sum.Run(env)
	fmt.Println(value)
	// Output:
	// 8
}
