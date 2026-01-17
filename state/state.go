// Package state implements the State monad for stateful computations.
//
// The [State] monad models computations that depend on and modify some global
// state. It threads state through a series of computations, allowing each step
// to both read and update the state whilst producing values.
package state

import "github.com/tomasbasham/gofp"

// State is a monad that models computations that depend on some global state.
//
// Type parameter S represents the state type.
// Type parameter A represents the value type.
type State[S, A any] struct {
	g func(S) (A, S)
}

// Map applies a function to transform the value of a [State], while preserving
// the state transitions.
func (s State[S, A]) Map(f func(A) A) State[S, A] {
	return Map(s, f)
}

// FlatMap composes two [State] computations by using the result of the first to
// create the second. Both computations share the same state. The state changes
// are threaded through both computations sequentially.
func (s State[S, A]) FlatMap(f func(A) State[S, A]) State[S, A] {
	return FlatMap(s, f)
}

// Run executes the [State] computation with the given initial state and returns
// both the value and the final state.
func (s State[S, A]) Run(state S) (A, S) {
	return s.g(state)
}

// Pure lifts a value into a [State] computation. The resulting [State] will
// always return the given value and leave the state unchanged.
func Pure[S, A any](a A) State[S, A] {
	return State[S, A]{
		func(s S) (A, S) {
			return a, s
		},
	}
}

// Get returns a [State] computation that provides the current state as its
// value without modifying the state. This is useful for extracting the state to
// use in further computations and possibly updating the state.
func Get[S any]() State[S, S] {
	return State[S, S]{
		func(state S) (S, S) {
			return state, state
		},
	}
}

// Gets returns a [State] computation that applies a function to the current
// state to extract a value, without modifying the state.
func Gets[S, A any](f func(S) A) State[S, A] {
	return State[S, A]{
		func(s S) (A, S) {
			return f(s), s
		},
	}
}

// Put returns a [State] computation that replaces the current state with the
// given state and returns [gofp.Unit] (a type with only one possible value,
// representing "no value").
func Put[S any](state S) State[S, gofp.Unit] {
	return State[S, gofp.Unit]{
		func(_ S) (gofp.Unit, S) {
			return gofp.UnitValue, state
		},
	}
}

// Modify returns a [State] computation that transforms the current state using
// the provided function and returns [gofp.Unit] (a type with only one possible
// value, representing "no value").
func Modify[S any](f func(S) S) State[S, gofp.Unit] {
	return State[S, gofp.Unit]{
		func(s S) (gofp.Unit, S) {
			return gofp.UnitValue, f(s)
		},
	}
}

// Map applies a function to transform the value type of a [State], while
// preserving the state transitions. Similar to the [State.Map] method but
// allows changing the value type.
func Map[S, A, B any](s State[S, A], f func(A) B) State[S, B] {
	return State[S, B]{
		func(state S) (B, S) {
			a, newState := s.g(state)
			return f(a), newState
		},
	}
}

// Apply applies a [State] computation containing a function to a [State]
// computation containing a value. This is useful for combining multiple [State]
// computations when the function to combine them is itself the result of a
// [State] computation.
func Apply[S, A, B any](s State[S, A], f State[S, func(A) B]) State[S, B] {
	return State[S, B]{
		func(state S) (B, S) {
			a, s1 := s.g(state)
			g, s2 := f.g(s1)
			return g(a), s2
		},
	}
}

// FlatMap composes two [State] computations by using the result of the first to
// create the second. Both computations share the same state. The state changes
// are threaded through both computations sequentially. Similar to the
// [State.FlatMap] method but allows changing the value type.
func FlatMap[S, A, B any](s State[S, A], f func(A) State[S, B]) State[S, B] {
	return State[S, B]{
		func(state S) (B, S) {
			a, newState := s.g(state)
			return f(a).g(newState)
		},
	}
}

// Zip combines two [State] computations into one using a combining
// function. Both computations are run sequentially with the same state threaded
// through them, and their values are combined using the given function.
func Zip[S, A, B, U any](sa State[S, A], sb State[S, B], f func(A, B) U) State[S, U] {
	return FlatMap(sa, func(a A) State[S, U] {
		return Map(sb, func(b B) U {
			return f(a, b)
		})
	})
}

// Sequence transforms a slice of [State] computations into a single [State]
// computation that returns a slice of values. The state is threaded through
// all computations in order.
func Sequence[S, A any](states []State[S, A]) State[S, []A] {
	values := Pure[S]([]A{})
	for _, s := range states {
		values = Zip(values, s, func(vs []A, a A) []A {
			return append(vs, a)
		})
	}
	return values
}
