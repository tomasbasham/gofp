package gofp

import "fmt"

// Either is a type that represents a value of one of two possible types. It is
// either Left or Right. It is typical that the left value represents a failure
// (an error) and the right value represents a success, but this is not a
// requirement.
//
// Type parameter T represents the left value type.
// Type parameter U represents the right value type.
type Either[T, U any] struct {
	left   T
	right  U
	isLeft bool
}

// Map applies a function to transform the right value, or otherwise preserves
// the left value.
func (e Either[T, U]) Map(fn func(U) U) Either[T, U] {
	return EitherMap(e, fn)
}

// MapLeft applies a function to transform the left value, or otherwise
// preserves the right value.
func (e Either[T, U]) MapLeft(fn func(T) T) Either[T, U] {
	return EitherMapLeft(e, fn)
}

// FlatMap composes two [Either] values by using the right value of the first to
// create the second, or otherwise preserves the left value.
func (e Either[T, U]) FlatMap(fn func(U) Either[T, U]) Either[T, U] {
	return EitherFlatMap(e, fn)
}

// FlatMapLeft composes two [Either] values by using the left value of the first
// to create the second, or otherwise preserves the right value.
func (e Either[T, U]) FlatMapLeft(fn func(T) Either[T, U]) Either[T, U] {
	return EitherFlatMapLeft(e, fn)
}

// Left returns an [Either] with a left value.
func Left[T, U any](value T) Either[T, U] {
	return Either[T, U]{left: value, isLeft: true}
}

// Right returns an [Either] with a right value.
func Right[T, U any](value U) Either[T, U] {
	return Either[T, U]{right: value}
}

// FromResult returns an [Either] from a [Result]. As is convention, the left
// value represents an error and the right value represents a success.
func FromResult[T any](r Result[T]) Either[error, T] {
	if r.IsErr() {
		return Left[error, T](r.UnwrapErr())
	}
	return Right[error](r.Unwrap())
}

// EitherMap applies a function to transform the right type of an [Either], or
// otherwise preserves the left value. Similar to the [Either.Map] method but
// allows changing the value type.
func EitherMap[T, U, V any](e Either[T, U], fn func(U) V) Either[T, V] {
	if e.isLeft {
		return Left[T, V](e.left)
	}
	return Right[T](fn(e.right))
}

// EitherMapLeft applies a function to transform the left type of an [Either],
// or otherwise preserves the right value. Similar to the [Either.MapLeft]
// method but allows changing the value type.
func EitherMapLeft[T, U, V any](e Either[T, U], fn func(T) V) Either[V, U] {
	if e.isLeft {
		return Left[V, U](fn(e.left))
	}
	return Right[V](e.right)
}

// EitherApply applies an [Either] containing a function to an [Either]
// containing a value. This is useful for combining multiple [Either] values
// when the function to combine them is itself an [Either].
func EitherApply[T, U, V any](e Either[T, U], efn Either[T, func(U) V]) Either[T, V] {
	if efn.isLeft {
		return Left[T, V](efn.left)
	}
	if e.isLeft {
		return Left[T, V](e.left)
	}
	return Right[T](efn.right(e.right))
}

// EitherApplyMap applies an [Either] containing a function to an [Either]
// containing a value, applying a combining function to the left values if both
// are Left.
//
// This is a special type of Apply that allows combining left values. It is
// useful for accumulating errors.
func EitherApplyMap[T, U, V any](e Either[T, U], efn Either[T, func(U) V], combine func(T, T) T) Either[T, V] {
	if efn.isLeft && e.isLeft {
		return Left[T, V](combine(efn.left, e.left))
	}
	if efn.isLeft {
		return Left[T, V](efn.left)
	}
	if e.isLeft {
		return Left[T, V](e.left)
	}
	return Right[T](efn.right(e.right))
}

// EitherFlatMap composes two [Either] values by using the result of the first
// right value to create the second, or otherwise preserves the left
// value. Similar to the [Either.FlatMap] method but allows changing the value
// type.
func EitherFlatMap[T, U, V any](e Either[T, U], fn func(U) Either[T, V]) Either[T, V] {
	if e.isLeft {
		return Left[T, V](e.left)
	}
	return fn(e.right)
}

// EitherFlatMapLeft composes two [Either] values by using the result of the
// first left value to create the second, or otherwise preserves the right
// value. Similar to the [Either.FlatMapLeft] method but allows changing the
// value type.
func EitherFlatMapLeft[T, U, V any](e Either[T, U], fn func(T) Either[V, U]) Either[V, U] {
	if e.isLeft {
		return fn(e.left)
	}
	return Right[V](e.right)
}

// EitherSequence transforms a slice of [Either] values into a single [Either]
// of a slice. If all values are Right, it returns Right with a slice of all
// values, preserving order. If any value is Left, it returns Left.
func EitherSequence[T, U any](eithers []Either[T, U]) Either[T, []U] {
	values := Right[T]([]U{})
	for _, e := range eithers {
		values = EitherFlatMap(values, func(vs []U) Either[T, []U] {
			return EitherMap(e, func(v U) []U {
				return append(vs, v)
			})
		})
	}
	return values
}

// EitherFold applies one of the two functions to the value of the [Either]
// depending on whether it is Left or Right.
func EitherFold[T, U, R any](e Either[T, U], left func(T) R, right func(U) R) R {
	if e.isLeft {
		return left(e.left)
	}
	return right(e.right)
}

func (e Either[T, U]) String() string {
	if e.isLeft {
		return fmt.Sprintf("Left(%v)", e.left)
	}
	return fmt.Sprintf("Right(%v)", e.right)
}

// IsLeft returns true if the [Either] is Left.
func (e Either[T, U]) IsLeft() bool {
	return e.isLeft
}

// IsRight returns true if the [Either] is Right.
func (e Either[T, U]) IsRight() bool {
	return !e.isLeft
}

// TryUnwrap returns the right value of the [Either] and a boolean indicating
// whether it is Right.
func (e Either[T, U]) TryUnwrap() (U, bool) {
	if e.isLeft {
		var zero U
		return zero, false
	}
	return e.right, true
}

// Unwrap returns the right value of the [Either] or panics if the [Either] is
// Left.
func (e Either[T, U]) Unwrap() U {
	if e.isLeft {
		panic(fmt.Sprintf("Cannot unwrap: Either is Left(%v)", e.left))
	}
	return e.right
}

// UnwrapOr returns the right value of the [Either] or a default value if the
// [Either] is Left.
func (e Either[T, U]) UnwrapOr(defaultValue U) U {
	if e.isLeft {
		return defaultValue
	}
	return e.right
}

// UnwrapOrElse returns the right value of the [Either] or the result of the
// given function if the [Either] is Left.
func (e Either[T, U]) UnwrapOrElse(fn func() U) U {
	if e.isLeft {
		return fn()
	}
	return e.right
}

// TryUnwrapLeft returns the left value of the [Either] and a boolean indicating
// whether it is Left.
func (e Either[T, U]) TryUnwrapLeft() (T, bool) {
	if !e.isLeft {
		var zero T
		return zero, false
	}
	return e.left, true
}

// UnwrapLeft returns the left value of the [Either] or panics if the [Either]
// is Right.
func (e Either[T, U]) UnwrapLeft() T {
	if !e.isLeft {
		panic(fmt.Sprintf("Cannot unwrap Left: Either is Right(%v)", e.right))
	}
	return e.left
}

// UnwrapLeftOr returns the left value of the pEither] or a default value if the
// [Either] is Right.
func (e Either[T, U]) UnwrapLeftOr(defaultValue T) T {
	if e.isLeft {
		return e.left
	}
	return defaultValue
}

// UnwrapLeftOrElse returns the left value of the [Either] or the result of the
// given function if the [Either] is Right.
func (e Either[T, U]) UnwrapLeftOrElse(fn func() T) T {
	if e.isLeft {
		return e.left
	}
	return fn()
}

// Swap returns a new [Either] with the left and right values swapped.
func (e Either[T, U]) Swap() Either[U, T] {
	if e.isLeft {
		return Right[U](e.left)
	}
	return Left[U, T](e.right)
}
