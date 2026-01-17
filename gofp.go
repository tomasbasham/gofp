// Package gofp provides functional programming primitives for Go.
//
// This package implements common functional types including [Option], [Either],
// and [Result], along with their associated operations. These types enable
// safer, more composable error handling and optional value management.
package gofp

// Unit is a type that has only one value.
type Unit struct{}

// UnitValue is the only value of type [Unit].
var UnitValue = Unit{}
