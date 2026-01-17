package gofp

import (
	"encoding/json"
	"fmt"
)

// Option is a type that represents an optional value. It is either Some or
// None.
//
// Type parameter T represents the value type.
type Option[T any] struct {
	value T
	valid bool
}

// Map applies a function to transform the value of an [Option].
func (o Option[T]) Map(fn func(T) T) Option[T] {
	return OptionMap(o, fn)
}

// FlatMap composes two [Option] values by using the value of the first to
// create the second.
func (o Option[T]) FlatMap(fn func(T) Option[T]) Option[T] {
	return OptionFlatMap(o, fn)
}

// Some returns an [Option] with a value.
func Some[T any](value T) Option[T] {
	return Option[T]{value: value, valid: true}
}

// None returns an [Option] with no value.
func None[T any]() Option[T] {
	return Option[T]{}
}

// OptionMap applies a function to transform the value type of an
// [Option]. Similar to the [Option.Map] method but allows changing the value
// type.
func OptionMap[T, U any](o Option[T], fn func(T) U) Option[U] {
	if !o.valid {
		return None[U]()
	}
	return Some(fn(o.value))
}

// OptionApply applies an [Option] containing a function to an [Option]
// containing a value. This is useful for combining multiple [Option] value when
// the function to combine them is itself an Option.
func OptionApply[T, U any](o Option[T], fn Option[func(T) U]) Option[U] {
	if !o.valid || !fn.valid {
		return None[U]()
	}
	return Some(fn.value(o.value))
}

// OptionFlatMap composes two [Option] value by using the result of the first to
// create the second. Similar to the [Option.FlatMap] method but allows changing
// the value type.
func OptionFlatMap[T, U any](o Option[T], fn func(T) Option[U]) Option[U] {
	if !o.valid {
		return None[U]()
	}
	return fn(o.value)
}

// OptionSequence transforms a slice of [Option] values into a single [Option]
// of a slice. If all values are Some, it returns Some with a slice of all
// values, preserving order. If any value is None, it returns None.
func OptionSequence[T any](options []Option[T]) Option[[]T] {
	values := Some([]T{})
	for _, o := range options {
		values = OptionFlatMap(values, func(vs []T) Option[[]T] {
			return OptionMap(o, func(v T) []T {
				return append(vs, v)
			})
		})
	}
	return values
}

// OptionFold applies one of two functions to the value of the [Option]
// depending on whether it is Some or None.
func OptionFold[T, R any](o Option[T], none func() R, some func(T) R) R {
	if !o.valid {
		return none()
	}
	return some(o.value)
}

func (o Option[T]) String() string {
	if o.valid {
		return fmt.Sprintf("Some(%v)", o.value)
	}
	return "None"
}

// IsSome returns true if the [Option] is Some.
func (o Option[T]) IsSome() bool {
	return o.valid
}

// IsNone returns true if the [Option] is None.
func (o Option[T]) IsNone() bool {
	return !o.valid
}

// TryUnwrap returns the value of the [Option] and a boolean indicating whether
// the [Option] is Some.
func (o Option[T]) TryUnwrap() (T, bool) {
	if !o.valid {
		var zero T
		return zero, false
	}
	return o.value, true
}

// Unwrap returns the value of the [Option] or panics if the [Option] is None.
func (o Option[T]) Unwrap() T {
	if !o.valid {
		panic("unwrapping None")
	}
	return o.value
}

// UnwrapOr returns the value of the [Option] or a default value if the [Option]
// is None.
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if !o.valid {
		return defaultValue
	}
	return o.value
}

// UnwrapOrElse returns the value of the [Option] or the result of the given
// function if the [Option] is None.
func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if !o.valid {
		return fn()
	}
	return o.value
}

// And returns the receiver [Options] if it is None, otherwise it returns the
// given [Option].
func (o Option[T]) And(opt Option[T]) Option[T] {
	if !o.valid {
		return None[T]()
	}
	return opt
}

// AndThen returns the receiver [Option] if it is None, otherwise it returns the
// [Option] produced by the given function.
func (o Option[T]) AndThen(fn func(T) Option[T]) Option[T] {
	if !o.valid {
		return None[T]()
	}
	return fn(o.value)
}

// Or returns the receiver [Option] if it is Some, otherwise it returns the
// given [Option].
func (o Option[T]) Or(opt Option[T]) Option[T] {
	if o.valid {
		return o
	}
	return opt
}

// OrElse returns the receiver [Option] if it is Some, otherwise it returns the
// [Option] produced by the given function.
func (o Option[T]) OrElse(fn func() Option[T]) Option[T] {
	if o.valid {
		return o
	}
	return fn()
}

// Filter converts a Some value to None if it doesn't satisfy the given
// predicate.
func (o Option[T]) Filter(fn func(T) bool) Option[T] {
	if !o.valid || !fn(o.value) {
		return None[T]()
	}
	return o
}

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.valid {
		return nil, nil
	}
	return json.Marshal(o.value)
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*o = None[T]()
		return nil
	}
	var value T
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	*o = Some(value)
	return nil
}

// NullableOption is a variant of [Option] that serializes None as JSON null.
type NullableOption[T any] Option[T]

func (o NullableOption[T]) MarshalJSON() ([]byte, error) {
	if !o.valid {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

func (o *NullableOption[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*o = NullableOption[T](None[T]())
		return nil
	}
	var value T
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	*o = NullableOption[T](Some(value))
	return nil
}
