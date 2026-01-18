// Database is a user management system simulation using the [reader.Reader]
// monad to handle dependencies like database configuration and user repository.
package main

import (
	"fmt"
	"time"

	"github.com/tomasbasham/gofp"
	"github.com/tomasbasham/gofp/reader"
)

// User represents a user entity.
type User struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

// UserRepository defines an interface for user data operations.
type UserRepository interface {
	FindByID(id int) gofp.Result[User]
	Save(user User) error
}

// DatabaseConfig represents the configuration for a database connection.
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// AppConfig represents the application configuration. It is used as the
// context for the Reader monad.
type AppConfig struct {
	EnableCache bool
	LogRequests bool
	Repository  UserRepository
}

// Mock repository implementation.
type mockUserRepository struct {
	config DatabaseConfig
	users  map[int]User
}

// FindByID simulates fetching a user from the database.
func (m *mockUserRepository) FindByID(id int) gofp.Result[User] {
	fmt.Printf("Using DB at %s:%d\n", m.config.Host, m.config.Port)
	user, ok := m.users[id]
	if !ok {
		return gofp.Err[User](fmt.Errorf("user with ID %d not found", id))
	}
	return gofp.Ok(user)
}

// Save simulates saving a user to the database.
func (m *mockUserRepository) Save(user User) error {
	fmt.Printf("Saving user to %s:%d\n", m.config.Host, m.config.Port)
	m.users[user.ID] = user
	return nil
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(cfg DatabaseConfig) UserRepository {
	return &mockUserRepository{
		config: cfg,
		users:  make(map[int]User),
	}
}

func main() {
	repo := NewUserRepository(DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "admin",
		Password: "secret",
	})

	// Set up environment.
	cfg := AppConfig{
		EnableCache: true,
		LogRequests: true,
		Repository:  repo,
	}

	user := User{
		ID:        1,
		Name:      "John Doe",
		CreatedAt: time.Now(),
	}

	result := userPipeline(user).Run(cfg)

	result.AndThen(func(u User) gofp.Result[User] {
		fmt.Printf("Successfully processed user: %s (created at: %s)\n", u.Name, u.CreatedAt.Format(time.RFC3339))
		return gofp.Ok(u)
	})
}

func userPipeline(user User) reader.Reader[AppConfig, gofp.Result[User]] {
	return reader.FlatMap(saveUser(user), func(err error) reader.Reader[AppConfig, gofp.Result[User]] {
		if err != nil {
			// Lift the error into a Result.
			return reader.Pure[AppConfig](gofp.Err[User](err))
		}
		return getUserByID(user.ID)
	})
}

func getUserByID(id int) reader.Reader[AppConfig, gofp.Result[User]] {
	return reader.New(func(cfg AppConfig) gofp.Result[User] {
		if cfg.LogRequests {
			fmt.Printf("Getting user with ID: %d\n", id)
		}

		result := cfg.Repository.FindByID(id)
		if cfg.EnableCache {
			fmt.Println("Checking cache for user data")
		}

		return result
	})
}

func saveUser(user User) reader.Reader[AppConfig, error] {
	return reader.New(func(cfg AppConfig) error {
		if cfg.LogRequests {
			fmt.Printf("Saving user: %s\n", user.Name)
		}
		return cfg.Repository.Save(user)
	})
}
