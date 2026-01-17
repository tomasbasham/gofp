package state_test

import (
	"fmt"
	"testing"

	"github.com/tomasbasham/gofp"
	"github.com/tomasbasham/gofp/state"
)

// Environment is a test environment type.
type Environment struct {
	Debug bool
	Name  string
	Value int
}

func TestPure(t *testing.T) {
	env := Environment{Debug: true, Name: "test", Value: 42}
	s := state.Pure[Environment](5)

	value, finalState := s.Run(env)
	if value != 5 {
		t.Errorf("expected value 5, got %v", value)
	}

	// Pure should not modify the state
	if !environmentEquals(env, finalState) {
		t.Errorf("expected state to be unchanged, got %v", finalState)
	}
}

func TestGet(t *testing.T) {
	env := Environment{Debug: true, Name: "test", Value: 42}
	s := state.Get[Environment]()

	value, finalState := s.Run(env)
	if !environmentEquals(env, value) {
		t.Errorf("expected value to be the environment, got %v", value)
	}

	// Get should not modify the state
	if !environmentEquals(env, finalState) {
		t.Errorf("expected state to be unchanged, got %v", finalState)
	}
}

func TestGets(t *testing.T) {
	env := Environment{Debug: true, Name: "test", Value: 42}
	s := state.Gets(func(e Environment) string {
		return e.Name
	})

	value, finalState := s.Run(env)
	if value != "test" {
		t.Errorf("expected value 'test', got %v", value)
	}

	// Gets should not modify the state
	if !environmentEquals(env, finalState) {
		t.Errorf("expected state to be unchanged, got %v", finalState)
	}
}

func TestPut(t *testing.T) {
	env := Environment{Debug: true, Name: "test", Value: 42}
	newEnv := Environment{Debug: false, Name: "prod", Value: 100}
	s := state.Put(newEnv)

	value, finalState := s.Run(env)
	if value != gofp.UnitValue {
		t.Errorf("expected UnitValue, got %v", value)
	}

	// Put should replace the state
	if !environmentEquals(newEnv, finalState) {
		t.Errorf("expected state to be replaced, got %v", finalState)
	}
}

func TestModify(t *testing.T) {
	env := Environment{Debug: true, Name: "test", Value: 42}
	modifiedEnv := Environment{Debug: true, Name: "modified", Value: 52}

	s := state.Modify(func(e Environment) Environment {
		e.Value += 10
		e.Name = "modified"
		return e
	})

	value, finalState := s.Run(env)
	if value != gofp.UnitValue {
		t.Errorf("expected UnitValue, got %v", value)
	}

	// Modify should transform the state
	if !environmentEquals(modifiedEnv, finalState) {
		t.Errorf("expected state to be modified, got %v", finalState)
	}
}

func TestMap(t *testing.T) {
	t.Run("maps value only", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		s := state.Pure[Environment](5)

		got := state.Map(s, func(n int) int {
			return n * 2
		})

		value, finalState := got.Run(env)
		if value != 10 {
			t.Errorf("expected value 10, got %v", value)
		}

		// Map should not affect state transitions
		if !environmentEquals(env, finalState) {
			t.Errorf("expected state to be unchanged, got %v", finalState)
		}
	})

	t.Run("changes value type", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		s := state.Pure[Environment](5)

		got := state.Map(s, func(n int) string {
			return fmt.Sprintf("Number: %d", n)
		})

		value, _ := got.Run(env)
		if value != "Number: 5" {
			t.Errorf("expected 'Number: 5', got %v", value)
		}
	})
}

func TestState_Map(t *testing.T) {
	t.Run("maps value only", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		s := state.Pure[Environment](5)

		got := s.Map(func(n int) int {
			return n * 2
		})

		value, finalState := got.Run(env)
		if value != 10 {
			t.Errorf("expected value 10, got %v", value)
		}

		// Map should not affect state transitions
		if !environmentEquals(env, finalState) {
			t.Errorf("expected state to be unchanged, got %v", finalState)
		}
	})
}

func TestApply(t *testing.T) {
	t.Run("applies function to value", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		s := state.Pure[Environment](5)

		// Create a State containing a function
		sf := state.Pure[Environment](func(n int) string {
			return fmt.Sprintf("Number: %d", n)
		})

		got := state.Apply(s, sf)

		value, finalState := got.Run(env)
		if value != "Number: 5" {
			t.Errorf("expected 'Number: 5', got %v", value)
		}

		// Apply should not affect state transitions
		if !environmentEquals(env, finalState) {
			t.Errorf("expected state to be unchanged, got %v", finalState)
		}
	})

	t.Run("threads state through computations", func(t *testing.T) {
		initialState := Environment{Debug: true, Name: "test", Value: 42}

		// Create a custom state computation that increments value
		incrementValue := state.Pure[Environment](func(s Environment) (int, Environment) {
			s.Value++
			return s.Value, s
		})

		// First state computation increments value and returns it
		s1 := state.FlatMap(incrementValue, func(f func(Environment) (int, Environment)) state.State[Environment, int] {
			return state.FlatMap(state.Get[Environment](), func(e Environment) state.State[Environment, int] {
				value, newState := f(e)
				return state.FlatMap(state.Put(newState), func(_ gofp.Unit) state.State[Environment, int] {
					return state.Pure[Environment](value)
				})
			})
		})

		// Second state computation contains a function
		sf := state.Pure[Environment](func(n int) string {
			return fmt.Sprintf("Number: %d", n)
		})

		got := state.Apply(s1, sf)

		value, finalState := got.Run(initialState)
		if value != "Number: 43" {
			t.Errorf("expected 'Number: 43', got %v", value)
		}

		if finalState.Value != 43 {
			t.Errorf("expected state to be modified, got %v", finalState)
		}
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("chains state computations", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		s := state.Pure[Environment](5)
		fn := func(n int) state.State[Environment, string] {
			return state.FlatMap(state.Get[Environment](), func(e Environment) state.State[Environment, string] {
				e.Value += n
				return state.FlatMap(state.Put(e), func(_ gofp.Unit) state.State[Environment, string] {
					return state.Pure[Environment](fmt.Sprintf("Number: %d", e.Value))
				})
			})
		}

		got := state.FlatMap(s, fn)

		value, finalState := got.Run(env)
		if value != "Number: 47" {
			t.Errorf("expected 'Number: 47', got %v", value)
		}

		if finalState.Value != 47 {
			t.Errorf("expected state value to be 47, got %v", finalState.Value)
		}
	})
}

func TestState_FlatMap(t *testing.T) {
	t.Run("threads state through computations", func(t *testing.T) {
		env := Environment{Debug: true, Name: "Alice", Value: 10}
		s := state.Pure[Environment]("Hello, ")
		fn := func(greeting string) state.State[Environment, string] {
			return state.FlatMap(state.Get[Environment](), func(e Environment) state.State[Environment, string] {
				e.Value += 5
				return state.FlatMap(state.Put(e), func(_ gofp.Unit) state.State[Environment, string] {
					return state.Pure[Environment](greeting + e.Name + "!")
				})
			})
		}

		got := s.FlatMap(fn)

		value, finalState := got.Run(env)
		if value != "Hello, Alice!" {
			t.Errorf("expected 'Hello, Alice!', got %v", value)
		}

		if finalState.Value != 15 {
			t.Errorf("expected state value to be 15, got %v", finalState.Value)
		}
	})
}

func TestZip(t *testing.T) {
	t.Run("combines two state computations", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		s1 := state.Pure[Environment](5)
		s2 := state.Pure[Environment](10)

		sum := state.Zip(s1, s2, func(a, b int) int {
			return a + b
		})

		value, _ := sum.Run(env)
		if value != 15 {
			t.Errorf("expected 15, got %v", value)
		}
	})

	t.Run("threads state through both computations", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}

		// First state increments Value by 5
		s1 := threadEnvironmentValue(func(e Environment) Environment {
			e.Value += 5
			return e
		})

		// Second state increments Value by 10
		s2 := threadEnvironmentValue(func(e Environment) Environment {
			e.Value += 10
			return e
		})

		sum := state.Zip(s1, s2, func(a, b int) int {
			return a + b
		})

		value, finalState := sum.Run(env)
		if value != 104 { // (42 + 5) + (42 + 10) = 47 + 57
			t.Errorf("expected 104, got %v", value)
		}

		if finalState.Value != 57 { // 42 + 5 + 10
			t.Errorf("expected state value to be 57, got %v", finalState.Value)
		}
	})

	t.Run("changes output type", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		s1 := state.Pure[Environment](5)
		s2 := state.Pure[Environment]("test")

		combined := state.Zip(s1, s2, func(a int, b string) string {
			return fmt.Sprintf("%s: %d", b, a)
		})

		value, _ := combined.Run(env)
		if value != "test: 5" {
			t.Errorf("expected 'test: 5', got %v", value)
		}
	})
}

func TestSequence(t *testing.T) {
	t.Run("combines multiple state computations", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		s1 := state.Pure[Environment]("hello")
		s2 := state.Pure[Environment]("world")
		s3 := state.Pure[Environment]("!")

		sequenced := state.Sequence([]state.State[Environment, string]{s1, s2, s3})

		values, _ := sequenced.Run(env)
		if len(values) != 3 || values[0] != "hello" || values[1] != "world" || values[2] != "!" {
			t.Errorf("expected ['hello', 'world', '!'], got %v", values)
		}
	})

	t.Run("threads state through all computations", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 10}

		// Each state adds an amount to Value
		s1 := threadEnvironmentValue(func(e Environment) Environment {
			e.Value += 1
			return e
		})

		s2 := threadEnvironmentValue(func(e Environment) Environment {
			e.Value += 2
			return e
		})

		s3 := threadEnvironmentValue(func(e Environment) Environment {
			e.Value += 3
			return e
		})

		sequenced := state.Sequence([]state.State[Environment, int]{s1, s2, s3})

		values, finalState := sequenced.Run(env)
		if len(values) != 3 || values[0] != 11 || values[1] != 13 || values[2] != 16 {
			t.Errorf("expected [11, 13, 16], got %v", values)
		}

		if finalState.Value != 16 {
			t.Errorf("expected final state Value to be 16, got %v", finalState.Value)
		}
	})
}

func TestComposition(t *testing.T) {
	t.Run("complex state transformation", func(t *testing.T) {
		env := Environment{Debug: true, Name: "Alice", Value: 42}

		// Get the current name
		getName := state.Gets(func(e Environment) string {
			return e.Name
		})

		// Use the name to create a greeting and modify the state
		greet := state.FlatMap(getName, func(name string) state.State[Environment, string] {
			return state.FlatMap(state.Get[Environment](), func(e Environment) state.State[Environment, string] {
				e.Debug = false // Toggle debug
				e.Value *= 2    // Double the value
				return state.FlatMap(state.Put(e), func(_ gofp.Unit) state.State[Environment, string] {
					return state.Pure[Environment](fmt.Sprintf("Hello, %s", name))
				})
			})
		})

		value, finalState := greet.Run(env)
		if value != "Hello, Alice" {
			t.Errorf("expected 'Hello, Alice', got %v", value)
		}

		if finalState.Debug != false || finalState.Value != 84 {
			t.Errorf("expected state to be modified, got %v", finalState)
		}
	})

	t.Run("building a counter", func(t *testing.T) {
		// Initial state
		counter := 0

		// Increment operation
		increment := threadInt(func(s int) int {
			return s + 1
		})

		// Decrement operation
		decrement := threadInt(func(s int) int {
			return s - 1
		})

		// Compose operations: increment twice, then decrement once
		operations := state.FlatMap(increment, func(v1 int) state.State[int, int] {
			return state.FlatMap(increment, func(v2 int) state.State[int, int] {
				return decrement
			})
		})

		result, finalState := operations.Run(counter)
		if result != 1 {
			t.Errorf("expected result 1, got %v", result)
		}

		if finalState != 1 {
			t.Errorf("expected final state 1, got %v", finalState)
		}
	})
}

func threadEnvironmentValue(fn func(e Environment) Environment) state.State[Environment, int] {
	return state.FlatMap(state.Get[Environment](), func(e Environment) state.State[Environment, int] {
		e = fn(e)
		return state.FlatMap(state.Put(e), func(_ gofp.Unit) state.State[Environment, int] {
			return state.Pure[Environment](e.Value)
		})
	})
}

func threadInt(fn func(s int) int) state.State[int, int] {
	return state.FlatMap(state.Get[int](), func(s int) state.State[int, int] {
		s = fn(s)
		return state.FlatMap(state.Put(s), func(_ gofp.Unit) state.State[int, int] {
			return state.Pure[int](s)
		})
	})
}

func environmentEquals(a, b Environment) bool {
	return a.Debug == b.Debug && a.Name == b.Name && a.Value == b.Value
}
