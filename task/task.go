package task

import "context"

// Task defines a unit of work that can be executed within a context.
// Implementations should respect context cancellation for graceful shutdown.
type Task interface {
	Run(ctx context.Context) error
}
