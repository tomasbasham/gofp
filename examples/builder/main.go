// Builder is a build system simulation using the [writer.Writer] monad to
// accumulate logs during the build process. The build can either succeed with a
// BuildArtifact or fail with a BuildError.
package main

import (
	"fmt"
	"strings"

	"github.com/tomasbasham/gofp"
	"github.com/tomasbasham/gofp/writer"
)

// BuildArtifact represents the result of a build step.
type BuildArtifact struct {
	Name  string
	Hash  string
	Stage string
}

// BuildError represents an error that occurred during the build process.
type BuildError struct {
	Stage string
	Msg   string
}

// BuildResult represents a build that can either succeed with an artifact or
// fail with an error.
type BuildResult = gofp.Either[BuildError, BuildArtifact]

// Build represents a build computation that accumulates log messages
type Build = writer.Writer[[]string, BuildResult]

// SliceMonoid implements the Monoid interface for string slices.
type SliceMonoid[T any] struct{}

// Empty returns an empty slice.
func (m SliceMonoid[T]) Empty() []T {
	return []T{}
}

// Append appends two slices together.
func (m SliceMonoid[T]) Append(a, b []T) []T {
	return append(a, b...)
}

func main() {
	result, log := build("main.go").Run()

	fmt.Println("Build log:")
	for _, entry := range log {
		fmt.Println(" ", entry)
	}

	fmt.Println(describe(result))
}

func build(source string) Build {
	compiled := compile(source).
		FlatMap(func(r BuildResult) Build {
			return gofp.EitherFold(r, propagateFailure, test)
		})

	// Observe logs so far without changing the result.
	listened := writer.Listen(compiled)
	_, logs := listened.Run()
	fmt.Println("Logs after compile + test:")
	for _, entry := range logs {
		fmt.Println(" ", entry)
	}
	fmt.Println("--------------")

	return compiled.FlatMap(func(r BuildResult) Build {
		return gofp.EitherFold(r, propagateFailure, docs)
	})
}

// Simulate compilation.
func compile(source string) Build {
	if strings.Contains(source, "broken") {
		return log("Compiling " + source).
			FlatMap(func(_ BuildResult) Build {
				return failAt("compile", "syntax error")
			})
	}

	artifact := BuildArtifact{
		Name:  strings.Replace(source, ".go", ".o", 1),
		Hash:  "abc123",
		Stage: "compiled",
	}

	return log("Compiled " + source).
		FlatMap(func(_ BuildResult) Build {
			return ok(artifact)
		})
}

// Simulate testing.
func test(a BuildArtifact) Build {
	if strings.Contains(a.Name, "flaky") {
		return log("Running tests for " + a.Name).
			FlatMap(func(_ BuildResult) Build {
				return failAt("test", "tests failed")
			})
	}

	a.Stage = "tested"

	return log("Tests passed for " + a.Name).
		FlatMap(func(_ BuildResult) Build {
			return ok(a)
		})
}

// Simulate documentation generation.
func docs(a BuildArtifact) Build {
	a.Name = strings.Replace(a.Name, ".o", ".doc", 1)
	a.Stage = "documented"

	return log("Generated docs for " + a.Name).
		FlatMap(func(_ BuildResult) Build {
			return ok(a)
		})
}

func log(msg string) Build {
	return writer.Tell[[]string, BuildResult](
		[]string{msg},
		SliceMonoid[string]{},
	)
}

func ok(artifact BuildArtifact) Build {
	return writer.Pure[[]string](
		gofp.Right[BuildError](artifact),
		SliceMonoid[string]{},
	)
}

func failAt(stage, msg string) Build {
	return writer.Pure[[]string](
		gofp.Left[BuildError, BuildArtifact](BuildError{stage, msg}),
		SliceMonoid[string]{},
	)
}

func propagateFailure(err BuildError) Build {
	return writer.Pure[[]string](
		gofp.Left[BuildError, BuildArtifact](err),
		SliceMonoid[string]{},
	)
}

func describe(result BuildResult) string {
	return gofp.EitherFold(
		result,
		func(err BuildError) string {
			return fmt.Sprintf("Build failed at %s: %s\n", err.Stage, err.Msg)
		},
		func(a BuildArtifact) string {
			return fmt.Sprintf("Build succeeded for: %s, with hash: %s\n", a.Name, a.Hash)
		},
	)
}
