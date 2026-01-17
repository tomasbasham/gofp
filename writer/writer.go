// Package writer implements the Writer monad for computations that produce
// output.
//
// The [Writer] monad models computations that accumulate output (such as logs
// or traces) alongside computing values. It requires a Monoid to define how
// outputs are combined.
package writer

// Monoid represents a type that can be combined with other values of the same
// type. It requires an empty value and a way to combine two values.
//
// Type parameter A represents the value type.
type Monoid[A any] interface {
	Empty() A
	Append(A, A) A
}

// Writer is a monad that models computations that produce output. It
// accumulates output alongside computing a value.
//
// Type parameter W represents the output/log type, which must satisfy the
// Monoid interface.
// Type parameter A represents the value type.
type Writer[W, A any] struct {
	g func() (A, W)

	// Monoid is a type that can be combined with other values of the same type.
	monoid Monoid[W]
}

// Map applies a function to transform the value of a [Writer], while preserving
// the output.
func (w Writer[W, A]) Map(f func(A) A) Writer[W, A] {
	return Map(w, f)
}

// FlatMap composes two [Writer] computations by using the result of the first
// to create the second. The outputs from both are combined according to the
// [Monoid].
func (w Writer[W, A]) FlatMap(f func(A) Writer[W, A]) Writer[W, A] {
	return FlatMap(w, f)
}

// Run executes the [Writer] computation and returns both the value and the
// accumulated output.
func (w Writer[W, A]) Run() (A, W) {
	return w.g()
}

// Pure lifts a value into a [Writer] computation with an empty output.
func Pure[W, A any](a A, m Monoid[W]) Writer[W, A] {
	return Writer[W, A]{
		g: func() (A, W) {
			return a, m.Empty()
		},
		monoid: m,
	}
}

// Tell creates a [Writer] computation that only produces output without
// computing a meaningful value. The result will be the zero value for type A.
func Tell[W, A any](w W, m Monoid[W]) Writer[W, A] {
	return Writer[W, A]{
		g: func() (A, W) {
			var zero A
			return zero, w
		},
		monoid: m,
	}
}

// TellWithValue creates a [Writer] computation that produces both a given value
// and output. This is useful when you need to add logs while preserving a
// value.
func TellWithValue[W, A any](a A, w W, m Monoid[W]) Writer[W, A] {
	return Writer[W, A]{
		g: func() (A, W) {
			return a, w
		},
		monoid: m,
	}
}

type listen[A, W any] struct {
	Value A
	Log   W
}

// Listen creates a [Writer] computation that includes its own output in the
// value. This is useful when you need to examine or transform the accumulated
// output.
func Listen[W, A any](w Writer[W, A]) Writer[W, listen[A, W]] {
	return Writer[W, listen[A, W]]{
		g: func() (listen[A, W], W) {
			a, log := w.Run()
			return listen[A, W]{
				Value: a,
				Log:   log,
			}, log
		},
		monoid: w.monoid,
	}
}

// Map applies a function to transform the value type of a [Writer], while
// preserving the output. Similar to the [Writer.Map] method but allows changing
// the value type.
func Map[W, A, B any](w Writer[W, A], f func(A) B) Writer[W, B] {
	return Writer[W, B]{
		g: func() (B, W) {
			a, log := w.g()
			return f(a), log
		},
		monoid: w.monoid,
	}
}

// Apply applies a [Writer] computation containing a function to a [Writer]
// computation containing a value. The outputs from both [Writer] values are
// combined according to the [Monoid]. This is useful for combining multiple
// [Writer] computations when the function to combine them is itself the result
// of a [Writer] computation.
func Apply[W, A, B any](w Writer[W, A], f Writer[W, func(A) B]) Writer[W, B] {
	return Writer[W, B]{
		g: func() (B, W) {
			a, logA := w.Run()  // orignally w.g()
			fn, logF := f.Run() // originally f.g()
			return fn(a), w.monoid.Append(logA, logF)
		},
		monoid: w.monoid,
	}
}

// FlatMap composes two [Writer] computations by using the result of the first
// to create the second. The outputs from both are combined according to the
// [Monoid]. Similar to the [Writer.FlatMap] method but allows changing the
// value type.
func FlatMap[W, A, B any](w Writer[W, A], f func(A) Writer[W, B]) Writer[W, B] {
	return Writer[W, B]{
		g: func() (B, W) {
			a, w1 := w.g()
			wb := f(a)
			b, w2 := wb.g()
			return b, w.monoid.Append(w1, w2)
		},
		monoid: w.monoid,
	}
}

// Zip combines two [Writer] computations into one using a combining
// function. Both computations are run sequentially, and their values are
// combined, along with their outputs, using the given function.
func Zip[W, A, B, U any](wa Writer[W, A], wb Writer[W, B], f func(A, B) U) Writer[W, U] {
	return FlatMap(wa, func(a A) Writer[W, U] {
		return Map(wb, func(b B) U {
			return f(a, b)
		})
	})
}
