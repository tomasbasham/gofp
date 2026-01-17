package state_test

import (
	"fmt"

	"github.com/tomasbasham/gofp/state"
)

func ExampleState_Map() {
	env := Environment{}
	s := state.Pure[Environment](5)
	doubled := s.Map(func(x int) int { return x * 2 })
	value, finalState := doubled.Run(env)
	fmt.Println(value, finalState)
	// Output:
	// 10 {false  0}
}

func ExampleState_FlatMap() {
	getName := state.Pure[Environment]("Alice")
	greet := getName.FlatMap(func(name string) state.State[Environment, string] {
		return state.Pure[Environment]("Hello, " + name)
	})

	env := Environment{Name: "Alice"}
	value, finalState := greet.Run(env)
	fmt.Println(value, finalState)
	// Output:
	// Hello, Alice {false Alice 0}
}

func ExampleState_Run() {
	env := Environment{}
	s := state.Pure[Environment](5)
	value, finalState := s.Run(env)
	fmt.Println(value, finalState)
	// Output:
	// 5 {false  0}
}

func ExamplePure() {
	env := Environment{}
	s := state.Pure[Environment](5)
	value, finalState := s.Run(env)
	fmt.Println(value, finalState)
	// Output:
	// 5 {false  0}
}

func ExampleGet() {
	env := Environment{Debug: true, Name: "Test", Value: 1}
	s := state.Get[Environment]()
	value, finalState := s.Run(env)
	fmt.Println(value, finalState)
	// Output:
	// {true Test 1} {true Test 1}
}

func GetName() state.State[Environment, string] {
	return state.Gets(func(e Environment) string {
		return e.Name
	})
}

func ExampleGets() {
	env := Environment{Debug: true, Name: "Test", Value: 1}
	s := GetName()
	value, finalState := s.Run(env)
	fmt.Println(value, finalState)
	// Output:
	// Test {true Test 1}
}

func ExamplePut() {
	env := Environment{Debug: false, Name: "Old", Value: 0}
	s := state.Put(Environment{Debug: true, Name: "New", Value: 1})
	_, finalState := s.Run(env)
	fmt.Println(finalState)
	// Output:
	// {true New 1}
}

func ExampleModify() {
	env := Environment{Debug: false, Name: "Test", Value: 5}
	increment := state.Modify(func(env Environment) Environment {
		env.Value += 1
		return env
	})
	_, finalState := increment.Run(env)
	fmt.Println(finalState)
	// Output:
	// {false Test 6}
}

func ExampleMap() {
	env := Environment{}
	s := state.Pure[Environment](5)
	doubled := s.Map(func(x int) int { return x * 2 })
	value, finalState := doubled.Run(env)
	fmt.Println(value, finalState)
	// Output:
	// 10 {false  0}
}

func ExampleApply() {
	env := Environment{}
	s := state.Pure[Environment](5)
	double := state.Pure[Environment](func(x int) int { return x * 2 })
	value, finalState := state.Apply(s, double).Run(env)
	fmt.Println(value, finalState)
	// Output:
	// 10 {false  0}
}

func ExampleFlatMap() {
	getName := state.Pure[Environment]("Alice")
	greet := getName.FlatMap(func(name string) state.State[Environment, string] {
		return state.Pure[Environment]("Hello, " + name)
	})

	env := Environment{Name: "Alice"}
	value, finalState := greet.Run(env)
	fmt.Println(value, finalState)
	// Output:
	// Hello, Alice {false Alice 0}
}

func ExampleZip() {
	env := Environment{Debug: true, Name: "test", Value: 42}
	s1 := state.Pure[Environment](5)
	s2 := state.Pure[Environment](10)

	sum := state.Zip(s1, s2, func(a, b int) int {
		return a + b
	})

	value, finalState := sum.Run(env)
	fmt.Println(value, finalState)
	// Output:
	// 15 {true test 42}
}

func ExampleSequence() {
	env := Environment{Debug: true, Name: "test", Value: 42}
	s1 := state.Pure[Environment]("hello")
	s2 := state.Pure[Environment]("world")
	s3 := state.Pure[Environment]("!")
	sequenced := state.Sequence([]state.State[Environment, string]{s1, s2, s3})
	value, finalState := sequenced.Run(env)
	fmt.Println(value, finalState)
	// Output:
	// [hello world !] {true test 42}
}
