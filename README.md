# gofp

A Go module providing monadic structures for functional programming patterns. It
offers type-safe containers for handling optionality, errors, state,
dependencies, and computations in a composable manner.

Whilst Go is primarily an imperative language, certain patterns from functional
programming can reduce boilerplate, improve type safety, and make complex
control flow more explicit. This library provides battle-tested abstractions for
developers who want these benefits without sacrificing Go's simplicity.

## Prerequisites

You will need the following things properly installed on your computer:

- [Go](https://golang.org/): any one of the **three latest major**
  [releases](https://golang.org/doc/devel/release.html)

## Installation

With [Go module](https://go.dev/wiki/Modules) support (Go 1.11+), simply add the
following import

```go
import "github.com/tomasbasham/gofp"
```

to your code, and then `go [build|run|test]` will automatically fetch the
necessary dependencies.

Otherwise, to install the `gofp` module, run the following command:

```bash
go get -u github.com/tomasbasham/gofp
```

## Philosophy

This library embraces functional programming patterns whilst respecting Go's
pragmatic nature. It's a tool, not a religion. Use it where it adds clarity and
safety, and reach for standard Go patterns where they're clearer.

Monads provide:
- **Explicit handling of edge cases** - No more forgotten nil checks
- **Composable operations** - Chain transformations declaratively
- **Type safety** - The compiler ensures you handle all cases
- **Reduced boilerplate** - Especially for error handling pipelines

But remember: Go is not Haskell. This library works best when integrated
thoughtfully into Go codebases, not when used to fight against the language's
idioms.

## Usage

To use this module, import the relevant packages into your Go code. The core
types (`Option`, `Result`, `Either`) are available directly from the main
package, whilst more specialised monads (`Reader`, `State`, `Writer`) are in
their own subpackages:

```go
import (
    "github.com/tomasbasham/gofp"
    "github.com/tomasbasham/gofp/reader"
    "github.com/tomasbasham/gofp/state"
    "github.com/tomasbasham/gofp/writer"
)
```

Each monad provides a consistent interface with `Map` and `FlatMap` methods for
transforming and composing computations. The library follows common functional
programming conventions, where `Map` transforms the contained value and
`FlatMap` (also called "bind" in other languages) chains operations that
themselves return monadic values.

Begin by identifying which monad fits your problem domain, then compose
operations using the provided combinators. The examples below demonstrate
typical usage patterns for each monad.

## Available Monads

### Option[T]

Represents an optional value that may or may not exist. Eliminates nil pointer
errors and makes optionality explicit in your type signatures.

**Mathematical form:**

```text
Option T = Some T | None
```

**When to use:**
- Database queries that may return no results
- Configuration values that might be missing
- Function parameters that are genuinely optional
- Parsing operations that might fail

**Example:**

```go
import "github.com/tomasbasham/gofp"

func FindUser(id string) gofp.Option[User] {
    user, found := db.Get(id)
    if !found {
        return gofp.None[User]()
    }
    return gofp.Some(user)
}

// Chain operations without nil checks
result := FindUser("123").
    Map(func(u User) User { 
        u.LastAccessed = time.Now()
        return u 
    }).
    UnwrapOr(DefaultUser)
```

**Key functions:**
- `Some(value)` - Create an Option containing a value
- `None[T]()` - Create an empty Option
- `Map(fn)` - Transform the contained value
- `FlatMap(fn)` - Chain operations that return Options
- `UnwrapOr(default)` - Extract value with fallback
- `Filter(predicate)` - Convert Some to None if predicate fails

### Result[T]

Represents a computation that may succeed with a value or fail with an
error. Provides structured error handling with stack traces and error wrapping.

**Mathematical form:**

```text
Result T = Ok T | Err error
```

**When to use:**
- Operations that can fail (file I/O, network calls, parsing)
- Validation pipelines
- Replacing multiple `if err != nil` checks
- Building error contexts

**Example:**

```go
func ProcessData(filename string) gofp.Result[Data] {
    return ReadFile(filename).
        FlatMap(ParseJSON).
        FlatMap(ValidateSchema).
        Map(Transform).
        Wrap("failed to process data")
}

func ReadFile(path string) gofp.Result[[]byte] {
    data, err := os.ReadFile(path)
    return gofp.FromReturn(data, err)
}

// Use the result
result := ProcessData("config.json")
if result.IsErr() {
    log.Printf("Error: %v\n%s", result.UnwrapErr(), result.StackTrace())
    return
}
data := result.Unwrap()
```

**Key functions:**
- `Ok(value)` - Create a successful Result
- `Err[T](error)` - Create a failed Result
- `FromReturn(value, err)` - Convert Go's `(T, error)` pattern
- `Map(fn)` - Transform success values
- `FlatMap(fn)` - Chain fallible operations
- `Wrap(msg)` - Add error context
- `Ensure(err, predicate)` - Validate and fail if predicate is false
- `Recover(fn)` - Convert errors to values

### Either[T, U]

Represents a value that can be one of two types. By convention, Left represents
failure and Right represents success, but both can hold any type.

**Mathematical form:**

```text
Either T U = Left T | Right U
```

**When to use:**
- Multiple error types that need different handling
- Accumulating validation errors
- Representing mutually exclusive outcomes
- When Result's error type is too restrictive

**Example:**

```go
type ValidationError struct {
    Field string
    Message string
}

func ValidateAge(age int) gofp.Either[ValidationError, int] {
    if age < 0 {
        return gofp.Left[ValidationError, int](ValidationError{
            Field: "age",
            Message: "must be non-negative",
        })
    }
    return gofp.Right[ValidationError](age)
}

// Accumulate multiple validation errors
func ValidateUser(user User) gofp.Either[[]ValidationError, User] {
    validations := []gofp.Either[ValidationError, gofp.Unit]{
        ValidateAge(user.Age).Map(func(int) gofp.Unit { 
            return gofp.UnitValue
        }),
        ValidateEmail(user.Email).Map(func(string) gofp.Unit {
            return gofp.UnitValue
        }),
    }
    // Use EitherSequence to collect all errors or succeed
    // Implementation depends on your error handling strategy
}
```

**Key functions:**
- `Left[T, U](value)` - Create a Left value
- `Right[T, U](value)` - Create a Right value
- `FromResult(result)` - Convert Result to Either
- `Map(fn)` - Transform Right values
- `MapLeft(fn)` - Transform Left values
- `FlatMap(fn)` - Chain Either-returning operations
- `Swap()` - Exchange Left and Right
- `EitherFold(leftFn, rightFn)` - Handle both cases

### Reader[E, A]

Represents a computation that reads from a shared environment. Provides
dependency injection without explicit parameter passing.

**Mathematical form:**

```text
Reader E A = E -> A
```

**When to use:**
- Dependency injection
- Configuration that flows through many functions
- Testing with different environments
- Avoiding global state

**Example:**

```go
import "github.com/tomasbasham/gofp/reader"

type Config struct {
    Database string
    APIKey   string
    Debug    bool
}

func GetConnection() reader.Reader[Config, *sql.DB] {
    return reader.Map(
        reader.Ask[Config](),
        func(cfg Config) *sql.DB {
            db, _ := sql.Open("postgres", cfg.Database)
            return db
        },
    )
}

func FetchUsers() reader.Reader[Config, []User] {
    return reader.FlatMap(
        GetConnection(),
        func(db *sql.DB) reader.Reader[Config, []User] {
            return reader.Pure[Config](queryUsers(db))
        },
    )
}

// Execute with configuration
config := Config{Database: "postgres://...", Debug: true}
users := FetchUsers().Run(config)
```

**Key functions:**
- `Pure[E, A](value)` - Lift value into Reader
- `Ask[E]()` - Access the environment
- `Map(fn)` - Transform the result
- `FlatMap(fn)` - Chain Reader operations
- `Local(reader, fn)` - Temporarily modify environment

### State[S, A]

Represents a computation that maintains and transforms state. Provides pure
functional state management without mutable variables.

**Mathematical form:**

```text
State S A = S -> (A, S)
```

**When to use:**
- Parser combinators
- State machines
- Game loops
- Any computation requiring sequential state updates

**Example:**

```go
import "github.com/tomasbasham/gofp/state"

type GameState struct {
    Score  int
    Lives  int
    Level  int
}

func AddPoints(points int) state.State[GameState, gofp.Unit] {
    return state.Modify(func(s GameState) GameState {
        s.Score += points
        return s
    })
}

func LoseLife() state.State[GameState, bool] {
    return state.FlatMap(
        state.Modify(func(s GameState) GameState {
            s.Lives--
            return s
        }),
        func(_ gofp.Unit) state.State[GameState, bool] {
            return state.Gets(func(s GameState) bool {
                return s.Lives > 0
            })
        },
    )
}

// Compose state operations
gameLoop := AddPoints(100).
    FlatMap(func(_ gofp.Unit) state.State[GameState, bool] {
        return LoseLife()
    })

initialState := GameState{Score: 0, Lives: 3, Level: 1}
stillAlive, finalState := gameLoop.Run(initialState)
```

**Key functions:**
- `Pure[S, A](value)` - Lift value without changing state
- `Get[S]()` - Access current state
- `Gets(fn)` - Extract value from state
- `Put(state)` - Replace state
- `Modify(fn)` - Transform state
- `Map(fn)` - Transform the result
- `FlatMap(fn)` - Chain state operations

### Writer[W, A]

Represents a computation that accumulates output (logs, events, metrics)
alongside producing a value. Requires a Monoid instance for combining outputs.

**Mathematical form:**

```text
Writer W A = () -> (A, W) where W is a Monoid
```

**When to use:**
- Collecting logs during computation
- Audit trails
- Gathering metrics
- Accumulating warnings

**Example:**

```go
import "github.com/tomasbasham/gofp/writer"

// Define a Monoid for combining string slices.
type SliceMonoid[T any] struct{}

func (SliceMonoid[T]) Empty() []T {
    return []T{}
}

func (SliceMonoid[T]) Append(a, b []T) []T {
    return append(a, b...)
}

func ProcessItem(item string) writer.Writer[[]string, int] {
    return writer.TellWithValue(
        len(item),
        []string{fmt.Sprintf("processed: %s", item)},
        SliceMonoid[string]{},
    )
}

func ProcessBatch(items []string) writer.Writer[[]string, int] {
    total := writer.Pure(0, SliceMonoid[string]{})
    for _, item := range items {
        total = writer.FlatMap(total, func(s int) writer.Writer[[]string, int] {
            return writer.Map(ProcessItem(item), func(length int) int { 
                return s + length
            })
        })
    }
    return total
}

result, logs := ProcessBatch([]string{"hello", "world"}).Run()
// result = 10
// logs = ["processed: hello", "processed: world"]
```

**Key functions:**
- `Pure(value, monoid)` - Create Writer without output
- `Tell(output, monoid)` - Create Writer with only output
- `TellWithValue(value, output, monoid)` - Create Writer with both
- `Map(fn)` - Transform the value
- `FlatMap(fn)` - Chain Writer operations
- `Listen(writer)` - Include output in the value

### Monoids

A Monoid is a mathematical structure that defines how to combine values of the
same type. It consists of:

1. **An identity element (empty)** - A value that, when combined with any other
   value, returns that value unchanged
1. **An associative binary operation (append)** - A way to combine two values
   that satisfies: `append(append(a, b), c) = append(a, append(b, c))`

In this module, Monoids are represented by an interface:

```go
type Monoid[A any] interface {
    Empty() A 
    Append(a, b A) A
}
```

The Writer monad uses Monoids to combine outputs from multiple computations. For
example, if you're accumulating log messages (strings), you'd use a string
concatenation Monoid. If you're collecting events (slices), you'd use a slice
concatenation Monoid.

## Common Patterns

### Sequencing Operations

All monads provide `Sequence` functions to transform slices of monadic values:

```go
// Options: returns None if any element is None
options := []gofp.Option[int]{gofp.Some(1), gofp.Some(2), gofp.Some(3)}
result := gofp.OptionSequence(options) // Some([]int{1, 2, 3})

// Results: returns Err if any operation fails
results := []gofp.Result[int]{gofp.Ok(1), gofp.Ok(2), gofp.Ok(3)}
combined := gofp.ResultSequence(results) // Ok([]int{1, 2, 3})

// Eithers: returns Left if any element is Left
eithers := []gofp.Either[string, int]{
    gofp.Right[string](1),
    gofp.Right[string](2),
}
sequenced := gofp.EitherSequence(eithers) // Right([]int{1, 2})
```

### Folding

Extract values by handling both success and failure cases:

```go
// Option
value := gofp.OptionFold(
    maybeUser,
    func() string { return "No user found" },
    func(u User) string { return u.Name },
)

// Result
message := gofp.ResultFold(
    operation,
    func(err error) string { return fmt.Sprintf("Error: %v", err) },
    func(val int) string { return fmt.Sprintf("Success: %d", val) },
)

// Either
output := gofp.EitherFold(
    validation,
    func(err ValidationError) string { return err.Message },
    func(data Data) string { return "Valid" },
)
```

### Combining Values

Use `Apply` functions to combine multiple monadic values:

```go
// Combine two Options
add := func(a int) func(int) int {
    return func(b int) int { return a + b }
}

opt1 := gofp.Some(5)
opt2 := gofp.Some(3)
optFn := gofp.Some(add)

// This is typically done with curried functions or helper combinators
result := gofp.OptionApply(opt1, gofp.OptionMap(opt2, add))
```

## When NOT to Use This Library

Functional programming patterns aren't always the right choice for Go
projects. Consider avoiding this library when:

1. **Your team is unfamiliar with functional concepts** - The learning curve can
   slow development and reduce code maintainability if the team doesn't
   understand monads.
1. **Simple error handling suffices** - For straightforward operations, Go's
   standard `if err != nil` pattern is clearer and more idiomatic.
1. **Performance is critical** - Monadic composition introduces additional
   function calls and allocations. Benchmark first if you're in a hot path.
1. **You're writing library code for the Go community** - Most Go developers
   expect idiomatic Go patterns. Using monads in public APIs creates friction.
1. **The problem domain is simple** - Don't add abstraction layers when direct,
   imperative code would be clearer.
1. **You can't justify the abstraction** - If you find yourself wrapping and
   unwrapping frequently, or if the monadic code is harder to read than
   imperative code, reconsider.

### Good Use Cases

- **Complex error handling pipelines** with multiple failure points
- **Configuration-heavy applications** where Reader monad reduces parameter
  passing
- **Parser combinators** where State monad shines
- **Validation logic** that accumulates errors
- **Code that benefits from explicit optionality** beyond nil checks

### Anti-patterns

```go
// DON'T: Use Result for simple operations
func Add(a, b int) gofp.Result[int] {
    return gofp.Ok(a + b) // Unnecessary wrapper
}

// DO: Use Result when operations can genuinely fail
func Divide(a, b int) gofp.Result[int] {
    if b == 0 {
        return gofp.Err[int](errors.New("division by zero"))
    }
    return gofp.Ok(a / b)
}

// DON'T: Wrap every nullable value in Option
func GetName(user *User) gofp.Option[string] {
    if user == nil {
        return gofp.None[string]()
    }
    return gofp.Some(user.Name) // Just check for nil normally
}

// DO: Use Option when optionality is semantically meaningful
func GetMiddleName(user User) gofp.Option[string] {
    // Middle name is genuinely optional in the domain
    if user.MiddleName == "" {
        return gofp.None[string]()
    }
    return gofp.Some(user.MiddleName)
}
```

## License

This project is licensed under the [MIT License](LICENSE).
