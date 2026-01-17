package gofp

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	// Maximum number of stack frames to capture for the stack trace. This value
	// has been chosen to balance performance and the amount of information
	// captured.
	pcCount = 30

	// Skip the first few stack frames to avoid including the runtime and library
	// code that created the Result.
	pcSkip = 3
)

// Result represents a computation that may fail. Unlike [Either], it
// specifically deals with error cases and provides utility methods for working
// with Go's error handling patterns.
//
// Type parameter T represents the value type.
type Result[T any] struct {
	value T
	err   error
	isErr bool
	stack string
}

// Map applies a function to transform the value of a [Result].
func (r Result[T]) Map(fn func(T) T) Result[T] {
	return ResultMap(r, fn)
}

// FlatMap composes two [Result] computations by using the value of the first
// to create the second.
func (r Result[T]) FlatMap(fn func(T) Result[T]) Result[T] {
	return ResultFlatMap(r, fn)
}

// Ok returns a [Result] with a value.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// Err returns a [Result] with an error.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err, isErr: true, stack: callers()}
}

// FromReturn returns a [Result] from a value and an error (Go's typical return
// pattern).
func FromReturn[T any](v T, err error) Result[T] {
	stack := ""
	if err != nil {
		stack = callers()
	}

	return Result[T]{
		value: v,
		err:   err,
		isErr: err != nil,
		stack: stack,
	}
}

func callers() string {
	pc := make([]uintptr, pcCount)
	n := runtime.Callers(pcSkip, pc)
	if n == 0 {
		// Return now to avoid processing the zero Frame that would otherwise be
		// returned by frames.Next below.
		return ""
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	var sb strings.Builder
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&sb, "%s:%d %s\n", frame.File, frame.Line, frame.Function)
		if frame.Function == "main.main" {
			break
		}

		if !more {
			break
		}
	}

	return sb.String()
}

// ResultMap applies a function to transform the value type of a
// [Result]. Similar to the [Result.Map] method but allows changing the value
// type.
func ResultMap[T, U any](r Result[T], fn func(T) U) Result[U] {
	if r.isErr {
		return Result[U]{err: r.err, isErr: true, stack: r.stack}
	}
	return Ok(fn(r.value))
}

// ResultApply applies a [Result] computation containing a function to a
// [Result] computation containing a value. This is useful for combining
// multiple [Result] computations when the function to combine them is itself
// the result of a [Result] computation.
func ResultApply[T, U any](r Result[T], fn Result[func(T) U]) Result[U] {
	if r.isErr {
		return Result[U]{err: r.err, isErr: true, stack: r.stack}
	}
	if fn.isErr {
		return Result[U]{err: fn.err, isErr: true, stack: fn.stack}
	}
	return Ok(fn.value(r.value))
}

// ResultFlatMap composes two [Result] computations by using the value of the
// first to create the second. Similar to the [Result.FlatMap] method but allows
// changing the value type.
func ResultFlatMap[T, U any](r Result[T], fn func(T) Result[U]) Result[U] {
	if r.isErr {
		return Result[U]{err: r.err, isErr: true, stack: r.stack}
	}
	return fn(r.value)
}

// ResultSequence transforms a slice of [Result] values into a single [Result]
// of a slice. If all values are Ok, it returns Ok with a slice of all
// values, preserving order. If any value is Err, it returns Err.
func ResultSequence[T any](results []Result[T]) Result[[]T] {
	values := Ok([]T{})
	for _, r := range results {
		values = ResultFlatMap(values, func(vs []T) Result[[]T] {
			return ResultMap(r, func(v T) []T {
				return append(vs, v)
			})
		})
	}
	return values
}

// ResultFold applies one of two functions to the value of the [Result]
// depending on whether it is an Ok or an Err.
func ResultFold[T, R any](r Result[T], errFn func(error) R, okFn func(T) R) R {
	if r.isErr {
		return errFn(r.err)
	}
	return okFn(r.value)
}

func (r Result[T]) String() string {
	if r.isErr {
		return fmt.Sprintf("Err(%v)", r.err)
	}
	return fmt.Sprintf("Ok(%v)", r.value)
}

// IsOk returns true if the [Result] is Ok.
func (r Result[T]) IsOk() bool {
	return !r.isErr
}

// IsErr returns true if the [Result] is an Err.
func (r Result[T]) IsErr() bool {
	return r.isErr
}

// TryUnwrap returns the value of the [Result] and a boolean indicating whether
// the [Result] is an Ok.
func (r Result[T]) TryUnwrap() (T, bool) {
	if r.isErr {
		var zero T
		return zero, false
	}
	return r.value, true
}

// Unwrap returns the value of the [Result] or panics if the [Result] is an Err.
func (r Result[T]) Unwrap() T {
	if r.isErr {
		panic(r.err)
	}
	return r.value
}

// UnwrapOr returns the value of the [Result] or a default value if the [Result]
// is an Err.
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.isErr {
		return defaultValue
	}
	return r.value
}

// UnwrapOrElse returns the value of the [Result] or the result of the given
// function if the [Result] is an Err.
func (r Result[T]) UnwrapOrElse(fn func() T) T {
	if r.isErr {
		return fn()
	}
	return r.value
}

// UnwrapErr returns the error of the [Result] or panics if the [Result] is an
// Ok.
func (r Result[T]) UnwrapErr() error {
	if !r.isErr {
		panic("unwrapping Ok")
	}
	return r.err
}

// And returns the receiver [Result] if it is an Err, otherwise it returns the
// given [Result].
func (r Result[T]) And(res Result[T]) Result[T] {
	if r.isErr {
		return r
	}
	return res
}

// AndThen returns the receiver [Result] if it is an Err, otherwise it returns
// the [Result] produced by the given function.
func (r Result[T]) AndThen(fn func(T) Result[T]) Result[T] {
	if r.isErr {
		return r
	}
	return fn(r.value)
}

// Or returns the receiver [Result] if it is an Ok, otherwise it returns the
// given [Result].
func (r Result[T]) Or(res Result[T]) Result[T] {
	if !r.isErr {
		return r
	}
	return res
}

// OrElse returns the receiver [Result] if it is an Ok, otherwise it returns the
// [Result] produced by the given function.
func (r Result[T]) OrElse(fn func(error) Result[T]) Result[T] {
	if !r.isErr {
		return r
	}
	return fn(r.err)
}

// Ensure converts a value to an Err if it doesn't satisfy the given predicate.
func (r Result[T]) Ensure(err error, pred func(T) bool) Result[T] {
	if r.isErr {
		return r
	}
	if !pred(r.value) {
		return Err[T](err)
	}
	return r
}

// EnsureWith converts a value to an Err if it doesn't satisfy the given
// predicate. The error is generated by the given function.
func (r Result[T]) EnsureWith(pred func(T) bool, errFn func(T) error) Result[T] {
	if r.isErr {
		return r
	}
	if !pred(r.value) {
		return Err[T](errFn(r.value))
	}
	return r
}

// Wrap adds additional context to the error if the [Result] is an Err.
func (r Result[T]) Wrap(msg string) Result[T] {
	if !r.isErr {
		return r
	}

	// Wrap the existing error with additional context, preserving the stack
	// trace.
	return Result[T]{
		err:   fmt.Errorf("%s: %w", msg, r.err),
		isErr: true,
		stack: r.stack,
	}
}

// ToReturn converts the [Result] back to Go's (value, error) pattern.
func (r Result[T]) ToReturn() (T, error) {
	return r.value, r.err
}

// Recover converts an error into a value using the given function if the
// [Result] is an Err.
func (r Result[T]) Recover(fn func(error) T) Result[T] {
	if r.isErr {
		return Ok(fn(r.err))
	}
	return r
}

// RecoverWith converts an error into a [Result] using the given function if the
// [Result] is an Err.
func (r Result[T]) RecoverWith(fn func(error) Result[T]) Result[T] {
	if r.isErr {
		return fn(r.err)
	}
	return r
}

// StackTrace returns the stack trace of the [Result] if it is an Err.
func (r Result[T]) StackTrace() string {
	if r.isErr {
		return r.stack
	}
	return ""
}
