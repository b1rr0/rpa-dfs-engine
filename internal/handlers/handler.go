package handlers

// Handler defines the contract for all request handlers in the RPA DFS Engine.
// Implementations should handle specific types of processing such as setup, testing, or data processing.
type Handler interface {
	// Execute performs the main handler logic and returns an error if the operation fails.
	Execute() error

	// GetDescription returns a human-readable description of what this handler does.
	GetDescription() string
}
