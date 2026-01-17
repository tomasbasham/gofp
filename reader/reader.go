// Package reader implements the Reader monad for dependency injection and
// environment-based computations.
//
// The [Reader] monad models computations that read values from a shared
// environment. It enables functional dependency injection by threading an
// environment through a series of computations without explicitly passing it to
// each function.
package reader

// Reader is a monad that models computations which read values from a shared
// environment. It is also known as the environment monad.
//
// Type parameter E represents the environment type.
// Type parameter A represents the value type.
type Reader[E, A any] struct {
	g func(E) A
}

// Map applies a function to transform the value of a [Reader].
func (r Reader[E, A]) Map(f func(A) A) Reader[E, A] {
	return Map(r, f)
}

// FlatMap composes two [Reader] computations by using the result of the first
// to create the second. Both computations share the same environment.
func (r Reader[E, A]) FlatMap(f func(A) Reader[E, A]) Reader[E, A] {
	return FlatMap(r, f)
}

// Run executes the [Reader] computation with the given environment and returns
// the value.
func (r Reader[E, A]) Run(env E) A {
	return r.g(env)
}

// Pure lifts a value into a [Reader] computation. The resulting [Reader] will
// always return the given value regardless of the environment.
func Pure[E, A any](a A) Reader[E, A] {
	return New(func(_ E) A { return a })
}

// New creates a [Reader] from a function.
func New[E, A any](f func(E) A) Reader[E, A] {
	return Reader[E, A]{g: f}
}

// Ask returns a [Reader] computation that provides the environment. This is a
// fundamental operation that allows creating computations that depend on the
// environment.
func Ask[E any]() Reader[E, E] {
	return New(func(e E) E { return e })
}

// Local creates a new [Reader] computation with a modified environment. The
// modification is temporary and only applies to this specific computation.
func Local[E, A any](r Reader[E, A], f func(E) E) Reader[E, A] {
	return New(func(e E) A { return r.Run(f(e)) })
}

// Map applies a function to transform the value type of a [Reader]. Similar to
// the [Reader.Map] method but allows changing the value type.
func Map[E, A, B any](r Reader[E, A], f func(A) B) Reader[E, B] {
	return Reader[E, B]{
		func(e E) B {
			return f(r.g(e))
		},
	}
}

// Apply applies a [Reader] computation containing a function to a [Reader]
// computation containing a value. This is useful for combining multiple
// [Reader] computations when the function to combine them is itself the result
// of a [Reader] computation.
func Apply[E, A, B any](r Reader[E, A], f Reader[E, func(A) B]) Reader[E, B] {
	return Reader[E, B]{
		func(e E) B {
			return f.g(e)(r.g(e))
		},
	}
}

// FlatMap composes two [Reader] computations by using the result of the first
// to create the second. Both computations share the same environment. Similar
// to the [Reader.FlatMap] method but allows changing the value type.
func FlatMap[E, A, B any](r Reader[E, A], f func(A) Reader[E, B]) Reader[E, B] {
	return Reader[E, B]{
		func(e E) B {
			return f(r.g(e)).g(e)
		},
	}
}

// Zip combines two [Reader] computations into one using a combining
// function. Both computations are run sequentially with the same environment
// threaded through them, and their values are combined using the given
// function.
func Zip[E, A, B, U any](ra Reader[E, A], rb Reader[E, B], f func(A, B) U) Reader[E, U] {
	return FlatMap(ra, func(a A) Reader[E, U] {
		return Map(rb, func(b B) U {
			return f(a, b)
		})
	})
}
