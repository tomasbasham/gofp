// Stack is a simple demonstration of error propagation using the [gofp.Result]
// monad.
package main

import (
	"fmt"

	"github.com/tomasbasham/gofp"
)

func main() {
	result := topMethod()
	fmt.Println(describe(result))
}

func topMethod() gofp.Result[int] {
	return middleMethod(gofp.Ok(5))
}

func middleMethod(r gofp.Result[int]) gofp.Result[int] {
	return r.FlatMap(func(i int) gofp.Result[int] {
		return bottomMethod(r)
	})
}

func bottomMethod(_ gofp.Result[int]) gofp.Result[int] {
	return gofp.Err[int](fmt.Errorf("an error"))
}

func describe(result gofp.Result[int]) string {
	return gofp.ResultFold(
		result,
		func(err error) string {
			return fmt.Sprintf("Error: %s\nStack Trace:\n%s", err, result.StackTrace())
		},
		func(v int) string {
			return fmt.Sprintf("Success: %d", v)
		},
	)
}
