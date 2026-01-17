package writer_test

import (
	"fmt"

	"github.com/tomasbasham/gofp/writer"
)

func ExampleWriter_Map() {
	w := writer.Pure[[]string](5, SliceMonoid[string]{})
	doubled := w.Map(func(x int) int { return x * 2 })
	value, output := doubled.Run()
	fmt.Println(value, output)
	// Output:
	// 10 []
}

func ExampleWriter_FlatMap() {
	w := writer.TellWithValue[[]string](5, []string{"first"}, SliceMonoid[string]{})
	doubled := w.FlatMap(func(x int) writer.Writer[[]string, int] {
		return writer.TellWithValue[[]string](x*2, []string{"second"}, SliceMonoid[string]{})
	})
	value, output := doubled.Run()
	fmt.Println(value, output)
	// Output:
	// 10 [first second]
}

func ExampleWriter_Run() {
	w := writer.Pure[[]string](5, SliceMonoid[string]{})
	value, output := w.Run()
	fmt.Println(value, output)
	// Output:
	// 5 []
}

func ExamplePure() {
	w := writer.Pure[[]string](5, SliceMonoid[string]{})
	value, output := w.Run()
	fmt.Println(value, output)
	// Output:
	// 5 []
}

func ExampleTell() {
	w := writer.Tell[[]string, int]([]string{"log entry"}, SliceMonoid[string]{})
	value, output := w.Run()
	fmt.Println(value, output)
	// Output:
	// 0 [log entry]
}

func ExampleTellWithValue() {
	w := writer.TellWithValue[[]string, int](5, []string{"log entry"}, SliceMonoid[string]{})
	value, output := w.Run()
	fmt.Println(value, output)
	// Output:
	// 5 [log entry]
}

func ExampleListen() {
	w := writer.TellWithValue[[]string, int](5, []string{"processed: 5"}, SliceMonoid[string]{})
	listened := writer.Listen(w)
	value, output := listened.Run()
	fmt.Println(value, output)
	// Output:
	// {5 [processed: 5]} [processed: 5]
}

func ExampleMap() {
	w := writer.Pure[[]string](5, SliceMonoid[string]{})
	parseInt := writer.Map(w, func(x int) string { return fmt.Sprintf("%d", x) })
	value, output := parseInt.Run()
	fmt.Println(value, output)
	// Output:
	// 5 []
}

func ExampleApply() {
	w := writer.Pure[[]string](5, SliceMonoid[string]{})
	double := writer.Pure[[]string](func(x int) int { return x * 2 }, SliceMonoid[string]{})
	value, output := writer.Apply(w, double).Run()
	fmt.Println(value, output)
	// Output:
	// 10 []
}

func ExampleFlatMap() {
	w := writer.Pure[[]string](5, SliceMonoid[string]{})
	doubled := writer.FlatMap(w, func(x int) writer.Writer[[]string, int] {
		return writer.TellWithValue[[]string](x, []string{"processed: " + fmt.Sprintf("%d", x)}, SliceMonoid[string]{})
	})
	value, output := doubled.Run()
	fmt.Println(value, output)
	// Output:
	// 5 [processed: 5]
}

func ExampleZip() {
	w1 := writer.TellWithValue[[]string](5, []string{"first"}, SliceMonoid[string]{})
	w2 := writer.TellWithValue[[]string](3, []string{"second"}, SliceMonoid[string]{})
	sum := writer.Zip(w1, w2, func(a, b int) int { return a + b })
	value, output := sum.Run()
	fmt.Println(value, output)
	// Output:
	// 8 [first second]
}
