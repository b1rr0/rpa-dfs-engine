package traverser

import "rpa-dfs-engine/internal/logger"

// Get retrieves a value from the context data
func (c *UserContext) Get(key string) (interface{}, bool) {
	value, exists := c.data[key]
	return value, exists
}

// NewContext creates a new context with user data
func NewContext(userData map[string]interface{}) *UserContext {
	logger.LogInfo("Creating new context with user data: %v", userData)
	return &UserContext{
		data: userData,
	}
}
