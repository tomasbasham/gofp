package writer_test

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/tomasbasham/gofp/writer"
)

// StringMonoid implements the Monoid interface for strings.
type StringMonoid struct{}

func (m StringMonoid) Empty() string {
	return ""
}

func (m StringMonoid) Append(a, b string) string {
	return a + b
}

// SliceMonoid implements the Monoid interface for slices.
type SliceMonoid[T any] struct{}

func (m SliceMonoid[T]) Empty() []T {
	return []T{}
}

func (m SliceMonoid[T]) Append(a, b []T) []T {
	return append(a, b...)
}

func TestPure(t *testing.T) {
	t.Run("creates writer with value and empty output", func(t *testing.T) {
		w := writer.Pure[string](42, StringMonoid{})

		value, output := w.Run()
		if value != 42 {
			t.Errorf("expected value 42, got %d", value)
		}
		if output != "" {
			t.Errorf("expected empty output, got %q", output)
		}
	})

	t.Run("creates writer with slice monoid", func(t *testing.T) {
		w := writer.Pure[[]string](42, SliceMonoid[string]{})

		value, output := w.Run()
		if value != 42 {
			t.Errorf("expected value 42, got %d", value)
		}
		if len(output) != 0 {
			t.Errorf("expected empty slice output, got %#v", output)
		}
	})
}

func TestTell(t *testing.T) {
	t.Run("creates writer with only output", func(t *testing.T) {
		w := writer.Tell[string, int]("log message", StringMonoid{})

		value, output := w.Run()
		if value != 0 {
			t.Errorf("expected zero value, got %d", value)
		}
		if output != "log message" {
			t.Errorf(`expected output "log message", got %q`, output)
		}
	})

	t.Run("creates writer with slice output", func(t *testing.T) {
		logs := []string{"error", "warning"}
		w := writer.Tell[[]string, int](logs, SliceMonoid[string]{})

		got, output := w.Run()
		if got != 0 {
			t.Errorf("expected zero value, got %d", got)
		}
		if !reflect.DeepEqual(output, logs) {
			t.Errorf("expected output %#v, got %#v", logs, output)
		}
	})
}

func TestListen(t *testing.T) {
	t.Run("includes output in value", func(t *testing.T) {
		w := writer.Pure[string](42, StringMonoid{}).
			FlatMap(func(x int) writer.Writer[string, int] {
				return writer.TellWithValue[string](x, fmt.Sprintf("processed: %d", x), StringMonoid{})
			})

		listened := writer.Listen(w)
		value, output := listened.Run()

		got := value.Value
		if got != 42 {
			t.Errorf("expected value 42, got %d", got)
		}

		wantLog := "processed: 42"
		if value.Log != wantLog {
			t.Errorf("expected log %q, got %q", wantLog, value.Log)
		}
		if output != wantLog {
			t.Errorf("expected output %q, got %q", wantLog, output)
		}
	})

	t.Run("works with slice logs", func(t *testing.T) {
		w := writer.Pure[[]string](42, SliceMonoid[string]{}).
			FlatMap(func(x int) writer.Writer[[]string, int] {
				return writer.TellWithValue[[]string](x, []string{fmt.Sprintf("processed: %d", x)}, SliceMonoid[string]{})
			})

		listened := writer.Listen(w)
		value, output := listened.Run()

		got := value.Value
		if got != 42 {
			t.Errorf("expected value 42, got %d", got)
		}

		wantLog := []string{"processed: 42"}
		if !slices.Equal(value.Log, wantLog) {
			t.Errorf("expected log %#v, got %#v", wantLog, value.Log)
		}
		if !slices.Equal(output, wantLog) {
			t.Errorf("expected output %#v, got %#v", wantLog, output)
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("transforms value while preserving output", func(t *testing.T) {
		w := writer.Pure[string](5, StringMonoid{}).
			FlatMap(func(x int) writer.Writer[string, int] {
				return writer.TellWithValue[string](x, "original", StringMonoid{})
			})

		doubled := writer.Map(w, func(x int) string {
			return fmt.Sprintf("Number: %d", x*2)
		})

		value, output := doubled.Run()

		if value != "Number: 10" {
			t.Errorf("expected 'Number: 10', got %v", value)
		}

		if output != "original" {
			t.Errorf("expected output 'original', got %v", output)
		}
	})

	t.Run("method version works the same", func(t *testing.T) {
		w := writer.Pure[string](5, StringMonoid{}).
			FlatMap(func(x int) writer.Writer[string, int] {
				return writer.TellWithValue[string](x, "original", StringMonoid{})
			})

		doubled := w.Map(func(x int) int {
			return x * 2
		})

		value, output := doubled.Run()

		if value != 10 {
			t.Errorf("expected 10, got %v", value)
		}

		if output != "original" {
			t.Errorf("expected output 'original', got %v", output)
		}
	})
}

func TestApply(t *testing.T) {
	t.Run("applies function to value", func(t *testing.T) {
		w := writer.Pure[string](5, StringMonoid{})

		fn := func(n int) string {
			return fmt.Sprintf("Number: %d", n)
		}

		wf := writer.Pure[string](fn, StringMonoid{})
		got := writer.Apply(w, wf)

		value, output := got.Run()
		if value != "Number: 5" {
			t.Errorf("expected 'Number: 5', got %v", value)
		}

		if output != "" {
			t.Errorf("expected empty output, got %v", output)
		}
	})

	t.Run("combines output from both writers", func(t *testing.T) {
		w := writer.Pure[[]string](5, SliceMonoid[string]{}).
			FlatMap(func(x int) writer.Writer[[]string, int] {
				return writer.TellWithValue[[]string](x, []string{"value log"}, SliceMonoid[string]{})
			})

		fn := func(n int) string {
			return fmt.Sprintf("Number: %d", n)
		}

		wf := writer.Pure[[]string](fn, SliceMonoid[string]{}).
			FlatMap(func(f func(int) string) writer.Writer[[]string, func(int) string] {
				return writer.TellWithValue[[]string](f, []string{"function log"}, SliceMonoid[string]{})
			})

		got := writer.Apply(w, wf)

		value, output := got.Run()
		if value != "Number: 5" {
			t.Errorf("expected 'Number: 5', got %v", value)
		}

		expectedOutput := []string{"value log", "function log"}
		if !reflect.DeepEqual(output, expectedOutput) {
			t.Errorf("expected %v, got %v", expectedOutput, output)
		}
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("chains computations and combines output", func(t *testing.T) {
		w := writer.Pure[string](5, StringMonoid{})
		fn := func(n int) writer.Writer[string, string] {
			return writer.Pure[string](fmt.Sprintf("Number: %d", n), StringMonoid{}).
				FlatMap(func(s string) writer.Writer[string, string] {
					return writer.TellWithValue[string](s, fmt.Sprintf("processed %d", n), StringMonoid{})
				})
		}

		got := writer.FlatMap(w, fn)

		value, output := got.Run()
		if value != "Number: 5" {
			t.Errorf("expected 'Number: 5', got %v", value)
		}

		if output != "processed 5" {
			t.Errorf("expected 'processed 5', got %v", output)
		}
	})

	t.Run("method version works the same", func(t *testing.T) {
		w := writer.Pure[string](5, StringMonoid{})
		fn := func(n int) writer.Writer[string, int] {
			return writer.Pure[string](n*2, StringMonoid{}).
				FlatMap(func(x int) writer.Writer[string, int] {
					return writer.TellWithValue[string](x, fmt.Sprintf("doubled to %d", x), StringMonoid{})
				})
		}

		got := w.FlatMap(fn)

		value, output := got.Run()
		if value != 10 {
			t.Errorf("expected 10, got %v", value)
		}

		if output != "doubled to 10" {
			t.Errorf("expected 'doubled to 10', got %v", output)
		}
	})

	t.Run("accumulates output according to monoid", func(t *testing.T) {
		w := writer.Pure[[]string](5, SliceMonoid[string]{}).
			FlatMap(func(x int) writer.Writer[[]string, int] {
				return writer.TellWithValue[[]string](x, []string{"first log"}, SliceMonoid[string]{})
			})

		fn := func(n int) writer.Writer[[]string, int] {
			return writer.Pure[[]string](n*2, SliceMonoid[string]{}).
				FlatMap(func(x int) writer.Writer[[]string, int] {
					return writer.TellWithValue[[]string](x, []string{"second log"}, SliceMonoid[string]{})
				})
		}

		got := w.FlatMap(fn)

		value, output := got.Run()
		if value != 10 {
			t.Errorf("expected 10, got %v", value)
		}

		expectedOutput := []string{"first log", "second log"}
		if !reflect.DeepEqual(output, expectedOutput) {
			t.Errorf("expected %v, got %v", expectedOutput, output)
		}
	})
}

func TestZip(t *testing.T) {
	t.Run("combines two writers with a function", func(t *testing.T) {
		w1 := writer.Pure[string](5, StringMonoid{})
		w2 := writer.Pure[string](10, StringMonoid{})

		sum := writer.Zip(w1, w2, func(a, b int) int {
			return a + b
		})

		value, output := sum.Run()
		if value != 15 {
			t.Errorf("expected 15, got %v", value)
		}

		if output != "" {
			t.Errorf("expected empty output, got %v", output)
		}
	})

	t.Run("combines writers with different value types", func(t *testing.T) {
		w1 := writer.Pure[string](5, StringMonoid{})
		w2 := writer.Pure[string]("test", StringMonoid{})

		combined := writer.Zip(w1, w2, func(a int, b string) string {
			return fmt.Sprintf("%s: %d", b, a)
		})

		value, output := combined.Run()
		if value != "test: 5" {
			t.Errorf("expected 'test: 5', got %v", value)
		}

		if output != "" {
			t.Errorf("expected empty output, got %v", output)
		}
	})

	t.Run("combines output from both writers", func(t *testing.T) {
		w1 := writer.Pure[[]string](5, SliceMonoid[string]{}).
			FlatMap(func(x int) writer.Writer[[]string, int] {
				return writer.TellWithValue[[]string](x, []string{"first log"}, SliceMonoid[string]{})
			})

		w2 := writer.Pure[[]string](10, SliceMonoid[string]{}).
			FlatMap(func(x int) writer.Writer[[]string, int] {
				return writer.TellWithValue[[]string](x, []string{"second log"}, SliceMonoid[string]{})
			})

		sum := writer.Zip(w1, w2, func(a, b int) int {
			return a + b
		})

		value, output := sum.Run()
		if value != 15 {
			t.Errorf("expected 15, got %v", value)
		}

		expectedOutput := []string{"first log", "second log"}
		if !reflect.DeepEqual(output, expectedOutput) {
			t.Errorf("expected %v, got %v", expectedOutput, output)
		}
	})
}

func TestComposition(t *testing.T) {
	t.Run("chains multiple operations", func(t *testing.T) {
		// Start with a pure value
		w := writer.Pure[[]string](5, SliceMonoid[string]{})

		// Log the initial value
		loggedStart := writer.FlatMap(w, func(x int) writer.Writer[[]string, int] {
			return writer.Pure[[]string](x, SliceMonoid[string]{}).
				FlatMap(func(n int) writer.Writer[[]string, int] {
					return writer.TellWithValue[[]string](n, []string{fmt.Sprintf("starting with %d", n)}, SliceMonoid[string]{})
				})
		})

		// Transform the value and log again
		doubled := writer.FlatMap(loggedStart, func(x int) writer.Writer[[]string, int] {
			return writer.Pure[[]string](x*2, SliceMonoid[string]{}).
				FlatMap(func(n int) writer.Writer[[]string, int] {
					return writer.TellWithValue[[]string](n, []string{fmt.Sprintf("doubled to %d", n)}, SliceMonoid[string]{})
				})
		})

		// Final transformation
		result := writer.Map(doubled, func(x int) string {
			return fmt.Sprintf("final: %d", x)
		})

		value, output := result.Run()
		expectedValue := "final: 10"
		if value != expectedValue {
			t.Errorf("expected %q, got %q", expectedValue, value)
		}

		expectedOutput := []string{"starting with 5", "doubled to 10"}
		if !reflect.DeepEqual(output, expectedOutput) {
			t.Errorf("expected %#v, got %#v", expectedOutput, output)
		}
	})
}
