package main

import (
	"context"
	"fmt"
	"time"

	di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
)

// Example services demonstrating context usage

type Logger interface {
	Log(message string)
}

type contextLogger struct {
	requestID string
}

func (l *contextLogger) Log(message string) {
	fmt.Printf("[RequestID: %s] %s\n", l.requestID, message)
}

type Database struct {
	ConnectionString string
	RequestID        string
}

type Repository struct {
	DB     *Database
	Logger Logger
}

type UserService struct {
	Repo      *Repository
	Logger    Logger
	RequestID string
}

type contextKey string

const requestIDKey contextKey = "requestID"

func main() {
	fmt.Println("=== Context Support in Dependency Injection ===")

	// Build container with context-aware providers
	container := di.NewBuilder().
		Register(func(r di.Register) {
			// Logger that uses context to get request ID
			r.AsScope(new(Logger), func(ctx context.Context) Logger {
				requestID := ctx.Value(requestIDKey).(string)
				return &contextLogger{requestID: requestID}
			}, nil)

			// Database that captures request ID from context
			r.AsScope(new(*Database), func(ctx context.Context) *Database {
				requestID := ctx.Value(requestIDKey).(string)
				return &Database{
					ConnectionString: "postgres://localhost:5432/mydb",
					RequestID:        requestID,
				}
			}, nil)

			// Repository that depends on Database and Logger
			r.AsScope(new(*Repository), func(ctx context.Context, db *Database, logger Logger) *Repository {
				logger.Log("Creating Repository")
				return &Repository{
					DB:     db,
					Logger: logger,
				}
			}, nil)

			// UserService that depends on Repository and Logger
			r.AsScope(new(*UserService), func(ctx context.Context, repo *Repository, logger Logger) *UserService {
				requestID := ctx.Value(requestIDKey).(string)
				logger.Log("Creating UserService")
				return &UserService{
					Repo:      repo,
					Logger:    logger,
					RequestID: requestID,
				}
			}, nil)
		}).
		MustBuild()

	resolver := container.Resolver()

	// Example 1: Simulating different HTTP requests with different request IDs
	fmt.Println("Example 1: Request-scoped context values")
	fmt.Println("------------------------------------------")

	simulateHTTPRequest(resolver, "req-001")
	fmt.Println()
	simulateHTTPRequest(resolver, "req-002")
	fmt.Println()

	// Example 2: Using timeout
	fmt.Println("Example 2: Context with timeout")
	fmt.Println("--------------------------------")
	demonstrateTimeout(container)
	fmt.Println()

	// Example 3: Using cancellation
	fmt.Println("Example 3: Context cancellation")
	fmt.Println("--------------------------------")
	demonstrateCancellation(container)
	fmt.Println()

	// Example 4: Context propagation through nested dependencies
	fmt.Println("Example 4: Context propagation")
	fmt.Println("-------------------------------")
	demonstrateContextPropagation(resolver)
}

func simulateHTTPRequest(resolver di.Resolver, requestID string) {
	// Create context with request ID (like in an HTTP handler)
	ctx := context.WithValue(context.Background(), requestIDKey, requestID)

	// Resolve service with context - all dependencies get the same context
	service := di.ResolveCtx[*UserService](ctx, resolver)

	service.Logger.Log("Processing user request")
	service.Logger.Log(fmt.Sprintf("Database connection: %s", service.Repo.DB.ConnectionString))

	// Verify that context propagated correctly
	fmt.Printf("Service RequestID: %s\n", service.RequestID)
	fmt.Printf("Database RequestID: %s\n", service.Repo.DB.RequestID)
	fmt.Printf("All services share the same request context: %v\n",
		service.RequestID == service.Repo.DB.RequestID)
}

func demonstrateTimeout(container di.Container) {
	// Register a slow service
	container.Register().AsType(new(*SlowService), func(ctx context.Context) (*SlowService, error) {
		fmt.Println("SlowService: Starting initialization (takes 2 seconds)...")
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("SlowService: Initialization completed")
			return &SlowService{}, nil
		case <-ctx.Done():
			return nil, fmt.Errorf("SlowService: Initialization cancelled: %w", ctx.Err())
		}
	}, nil)

	// Try with insufficient timeout
	fmt.Println("Trying to resolve SlowService with 500ms timeout...")
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Caught panic (expected): %v\n", r)
		}
	}()

	container.Resolver().TypeCtx(ctx, new(*SlowService), nil)
}

func demonstrateCancellation(container di.Container) {
	// Register another service
	container.Register().AsType(new(*CancellableService), func(ctx context.Context) *CancellableService {
		fmt.Println("CancellableService: Checking context...")
		if err := ctx.Err(); err != nil {
			panic(fmt.Errorf("context already cancelled: %w", err))
		}
		return &CancellableService{}
	}, nil)

	// Create cancellable context and cancel it immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel before resolution

	fmt.Println("Trying to resolve with already-cancelled context...")

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Caught panic (expected): %v\n", r)
		}
	}()

	container.Resolver().TypeCtx(ctx, new(*CancellableService), nil)
}

func demonstrateContextPropagation(resolver di.Resolver) {
	ctx := context.WithValue(context.Background(), requestIDKey, "req-propagation-test")

	fmt.Println("Resolving UserService (depends on Repository -> Database + Logger)...")
	service := di.ResolveCtx[*UserService](ctx, resolver)

	fmt.Printf("✓ Context value propagated to all levels:\n")
	fmt.Printf("  - UserService.RequestID: %s\n", service.RequestID)
	fmt.Printf("  - Database.RequestID: %s\n", service.Repo.DB.RequestID)
	fmt.Printf("  - Logger has access to the same context\n")
}

// Helper types for examples
type SlowService struct{}
type CancellableService struct{}
