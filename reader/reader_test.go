package reader_test

import (
	"fmt"
	"testing"

	"github.com/tomasbasham/gofp/reader"
)

// Environment is a test environment type.
type Environment struct {
	Debug bool
	Name  string
	Value int
}

func TestPure(t *testing.T) {
	env := Environment{Debug: true, Name: "test", Value: 42}
	r := reader.Pure[Environment]("test value")

	if got := r.Run(env); got != "test value" {
		t.Errorf("expected 'test value', got %v", got)
	}
}

func TestNew(t *testing.T) {
	env := Environment{Debug: true, Name: "test", Value: 42}
	r := reader.New(func(e Environment) int {
		return e.Value
	})

	if got := r.Run(env); got != 42 {
		t.Errorf("expected 42, got %v", got)
	}
}

func TestAsk(t *testing.T) {
	env := Environment{Debug: true, Name: "test", Value: 42}
	r := reader.Ask[Environment]()

	if got := r.Run(env); got != env {
		t.Errorf("expected environment, got %v", got)
	}
}

func TestLocal(t *testing.T) {
	t.Run("modifies environment", func(t *testing.T) {
		env := Environment{Debug: false, Name: "prod", Value: 100}
		r := reader.Ask[Environment]()

		modified := reader.Local(r, func(e Environment) Environment {
			e.Debug = true
			e.Name = "dev"
			return e
		})

		result := modified.Run(env)
		if !result.Debug || result.Name != "dev" || result.Value != 100 {
			t.Errorf("expected modified environment, got %v", result)
		}

		// Original environment should be unchanged
		if env.Debug || env.Name != "prod" || env.Value != 100 {
			t.Errorf("expected original environment unchanged, got %v", env)
		}
	})

	t.Run("temporary modification", func(t *testing.T) {
		env := Environment{Debug: false, Name: "prod", Value: 100}
		r := reader.Ask[Environment]()

		// First local modification
		modified1 := reader.Local(r, func(e Environment) Environment {
			e.Debug = true
			return e
		})

		// Second local modification
		modified2 := reader.Local(r, func(e Environment) Environment {
			e.Name = "dev"
			return e
		})

		result1 := modified1.Run(env)
		if !result1.Debug || result1.Name != "prod" || result1.Value != 100 {
			t.Errorf("expected first modified environment, got %v", result1)
		}

		result2 := modified2.Run(env)
		if result2.Debug || result2.Name != "dev" || result2.Value != 100 {
			t.Errorf("expected second modified environment, got %v", result2)
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("changes value type", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		r := reader.Pure[Environment](5)
		fn := func(n int) string {
			return fmt.Sprintf("Number: %d", n)
		}

		got := reader.Map(r, fn)

		if result := got.Run(env); result != "Number: 5" {
			t.Errorf("expected 'Number: 5', got %v", result)
		}
	})
}

func TestReader_Map(t *testing.T) {
	t.Run("maps value", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		r := reader.Pure[Environment]("test")
		fn := func(s string) string {
			return s + "_processed"
		}

		got := r.Map(fn)

		if result := got.Run(env); result != "test_processed" {
			t.Errorf("expected test_processed, got %v", result)
		}
	})

	t.Run("preserves environment", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		r := reader.Ask[Environment]()
		fn := func(e Environment) Environment {
			return Environment{Debug: !e.Debug, Name: e.Name + "_mapped", Value: e.Value * 2}
		}

		got := r.Map(fn)

		result := got.Run(env)
		if !result.Debug && result.Name == "test_mapped" && result.Value == 84 {
			// This is correct behavior - the environment is passed through
		} else {
			t.Errorf("expected mapped environment, got %v", result)
		}
	})
}

func TestApply(t *testing.T) {
	t.Run("applies function to value", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		r := reader.Pure[Environment](5)

		// Create a Reader containing a function
		rf := reader.Pure[Environment](func(n int) string {
			return fmt.Sprintf("Number: %d", n)
		})

		got := reader.Apply(r, rf)

		if result := got.Run(env); result != "Number: 5" {
			t.Errorf("expected 'Number: 5', got %v", result)
		}
	})

	t.Run("applies environment-dependent function", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		r := reader.Pure[Environment](5)

		// Create a Reader containing a function that depends on the environment
		rf := reader.New(func(e Environment) func(int) string {
			prefix := "DEBUG: "
			if !e.Debug {
				prefix = "PROD: "
			}
			return func(n int) string {
				return fmt.Sprintf("%s%d", prefix, n)
			}
		})

		got := reader.Apply(r, rf)

		if result := got.Run(env); result != "DEBUG: 5" {
			t.Errorf("expected 'DEBUG: 5', got %v", result)
		}
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("changes value type", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		r := reader.Pure[Environment](5)
		fn := func(n int) reader.Reader[Environment, string] {
			return reader.New(func(e Environment) string {
				return fmt.Sprintf("Number: %d in %s", n, e.Name)
			})
		}

		got := reader.FlatMap(r, fn)

		if result := got.Run(env); result != "Number: 5 in test" {
			t.Errorf("expected 'Number: 5 in test', got %v", result)
		}
	})
}

func TestReader_FlatMap(t *testing.T) {
	t.Run("flatmaps computations", func(t *testing.T) {
		env := Environment{Debug: true, Name: "Alice", Value: 42}
		r := reader.Pure[Environment]("Hello, ")
		fn := func(s string) reader.Reader[Environment, string] {
			return reader.New(func(e Environment) string {
				return s + e.Name
			})
		}

		got := r.FlatMap(fn)

		if result := got.Run(env); result != "Hello, Alice" {
			t.Errorf("expected 'Hello, Alice', got %v", result)
		}
	})
}

func TestZip(t *testing.T) {
	t.Run("combines two readers", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		r1 := reader.Pure[Environment](5)
		r2 := reader.Pure[Environment](10)

		sum := reader.Zip(r1, r2, func(a, b int) int {
			return a + b
		})

		if result := sum.Run(env); result != 15 {
			t.Errorf("expected 15, got %v", result)
		}
	})

	t.Run("combines with environment", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		r1 := reader.Pure[Environment](5)
		r2 := reader.New(func(e Environment) int {
			return e.Value
		})

		sum := reader.Zip(r1, r2, func(a, b int) int {
			return a + b
		})

		if result := sum.Run(env); result != 47 {
			t.Errorf("expected 47, got %v", result)
		}
	})

	t.Run("changes output type", func(t *testing.T) {
		env := Environment{Debug: true, Name: "test", Value: 42}
		r1 := reader.Pure[Environment](5)
		r2 := reader.Pure[Environment]("test")

		combined := reader.Zip(r1, r2, func(a int, b string) string {
			return fmt.Sprintf("%s: %d", b, a)
		})

		if result := combined.Run(env); result != "test: 5" {
			t.Errorf("expected 'test: 5', got %v", result)
		}
	})
}

func TestComposition(t *testing.T) {
	env := Environment{Debug: true, Name: "Alice", Value: 42}

	// Create a reader that extracts the name from the environment
	getName := reader.Map(reader.Ask[Environment](), func(e Environment) string {
		return e.Name
	})

	// Create a reader that depends on the result of getName
	greet := reader.FlatMap(getName, func(s string) reader.Reader[Environment, string] {
		return reader.Pure[Environment]("Hello, " + s)
	})

	if result := greet.Run(env); result != "Hello, Alice" {
		t.Errorf("expected 'Hello, Alice', got %v", result)
	}
}
