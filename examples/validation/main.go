// Validation is a simple demonstration of accumulating validation errors using
// the [gopf.Either] monad.
package main

import (
	"fmt"
	"strings"

	"github.com/tomasbasham/gofp"
)

type User struct {
	Username string
	Email    string
	Age      int
}

type ValidationError struct {
	Field   string
	Message string
}

func main() {
	user := User{
		Username: "ab",
		Email:    "invalid-email",
		Age:      15,
	}

	result := validateUser(user)

	fmt.Println(describe(result))
}

func validateUsername(name string) gofp.Either[[]ValidationError, string] {
	if len(name) < 3 {
		return gofp.Left[[]ValidationError, string](
			[]ValidationError{{"username", "must be at least 3 characters"}},
		)
	}
	return gofp.Right[[]ValidationError](name)
}

func validateEmail(email string) gofp.Either[[]ValidationError, string] {
	if !strings.Contains(email, "@") {
		return gofp.Left[[]ValidationError, string](
			[]ValidationError{{"email", "must contain @"}},
		)
	}
	return gofp.Right[[]ValidationError](email)
}

func validateAge(age int) gofp.Either[[]ValidationError, int] {
	if age < 18 {
		return gofp.Left[[]ValidationError, int](
			[]ValidationError{{"age", "must be at least 18"}},
		)
	}
	return gofp.Right[[]ValidationError](age)
}

func validateUser(u User) gofp.Either[[]ValidationError, User] {
	return gofp.EitherApplyMap(
		validateAge(u.Age),
		gofp.EitherApplyMap(
			validateEmail(u.Email),
			gofp.EitherApplyMap(
				validateUsername(u.Username),
				gofp.Right[[]ValidationError](
					func(username string) func(string) func(int) User {
						return func(email string) func(int) User {
							return func(age int) User {
								return User{
									Username: username,
									Email:    email,
									Age:      age,
								}
							}
						}
					}),
				appendErrors,
			),
			appendErrors,
		),
		appendErrors,
	)
}

func appendErrors(a, b []ValidationError) []ValidationError {
	return append(a, b...)
}

func describe(result gofp.Either[[]ValidationError, User]) string {
	return gofp.EitherFold(
		result,
		func(errs []ValidationError) string {
			sb := strings.Builder{}
			sb.WriteString("Validation failed:\n")
			for _, e := range errs {
				fmt.Fprintf(&sb, "  %s: %s\n", e.Field, e.Message)
			}
			return sb.String()
		},
		func(valid User) string {
			return fmt.Sprintf("Valid user: %s", valid.Username)
		},
	)
}
