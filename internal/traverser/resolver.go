package traverser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"rpa-dfs-engine/internal/logger"
)

// TemplateResolver handles template substitution for workflow values
type TemplateResolver struct {
	context     *UserContext
	selectorMap map[string]string
}

// NewTemplateResolver creates a new template resolver
func NewTemplateResolver(context *UserContext) *TemplateResolver {
	resolver := &TemplateResolver{
		context:     context,
		selectorMap: make(map[string]string),
	}

	// Initialize with default selector mappings
	resolver.initializeDefaultSelectors()

	return resolver
}

// Resolve replaces template variables in a string with actual values
func (r *TemplateResolver) Resolve(template string) string {
	if template == "" {
		return template
	}

	// Regex to match {{expression}}
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	return re.ReplaceAllStringFunc(template, func(match string) string {
		// Extract the expression between {{ }}
		expr := strings.Trim(match, "{}")
		expr = strings.TrimSpace(expr)

		logger.LogDebug("Resolving template expression: %s", expr)

		// Get value from context
		if value, exists := r.context.data[expr]; exists {
			resolved := fmt.Sprintf("%v", value)
			logger.LogDebug("Template resolved: %s -> %s", expr, resolved)
			return resolved
		}

		logger.LogWarning("Template variable not found: %s", expr)
		return match // Return original if not found
	})
}

// EvaluateCondition evaluates a conditional expression
func (r *TemplateResolver) EvaluateCondition(expr string) bool {
	// First resolve any templates in the expression
	resolved := r.Resolve(expr)
	logger.LogDebug("Evaluating condition: %s -> %s", expr, resolved)

	// Handle different operators
	if strings.Contains(resolved, " > ") {
		return r.evaluateComparison(resolved, ">")
	}
	if strings.Contains(resolved, " < ") {
		return r.evaluateComparison(resolved, "<")
	}
	if strings.Contains(resolved, " >= ") {
		return r.evaluateComparison(resolved, ">=")
	}
	if strings.Contains(resolved, " <= ") {
		return r.evaluateComparison(resolved, "<=")
	}
	if strings.Contains(resolved, " == ") {
		return r.evaluateEquality(resolved, "==")
	}
	if strings.Contains(resolved, " != ") {
		return r.evaluateEquality(resolved, "!=")
	}
	if strings.Contains(resolved, " contains ") {
		return r.evaluateContains(resolved)
	}

	// If no operator found, treat as boolean
	return false
}

// evaluateComparison handles numeric comparisons (>, <, >=, <=)
func (r *TemplateResolver) evaluateComparison(expr, operator string) bool {
	parts := strings.Split(expr, " "+operator+" ")
	if len(parts) != 2 {
		logger.LogWarning("Invalid comparison expression: %s", expr)
		return false
	}

	left := strings.TrimSpace(parts[0])
	right := strings.TrimSpace(parts[1])

	leftNum, leftErr := strconv.ParseFloat(left, 64)
	rightNum, rightErr := strconv.ParseFloat(right, 64)

	if leftErr != nil || rightErr != nil {
		logger.LogWarning("Non-numeric comparison: %s %s %s", left, operator, right)
		return false
	}

	switch operator {
	case ">":
		return leftNum > rightNum
	case "<":
		return leftNum < rightNum
	case ">=":
		return leftNum >= rightNum
	case "<=":
		return leftNum <= rightNum
	}

	return false
}

// evaluateEquality handles equality comparisons (==, !=)
func (r *TemplateResolver) evaluateEquality(expr, operator string) bool {
	parts := strings.Split(expr, " "+operator+" ")
	if len(parts) != 2 {
		logger.LogWarning("Invalid equality expression: %s", expr)
		return false
	}

	left := strings.TrimSpace(parts[0])
	right := strings.TrimSpace(parts[1])

	// Remove quotes if present
	left = strings.Trim(left, `"'`)
	right = strings.Trim(right, `"'`)

	switch operator {
	case "==":
		return left == right
	case "!=":
		return left != right
	}

	return false
}

// evaluateContains handles string contains operations
func (r *TemplateResolver) evaluateContains(expr string) bool {
	parts := strings.Split(expr, " contains ")
	if len(parts) != 2 {
		logger.LogWarning("Invalid contains expression: %s", expr)
		return false
	}

	left := strings.TrimSpace(parts[0])
	right := strings.TrimSpace(parts[1])

	// Remove quotes if present
	left = strings.Trim(left, `"'`)
	right = strings.Trim(right, `"'`)

	return strings.Contains(left, right)
}

func (r *TemplateResolver) EvaluateDataCheck(check *DataCheck) bool {
	if check == nil {
		return false
	}

	value, exists := r.context.data[check.DataPath]
	if !exists {
		logger.LogDebug("Data path not found: %s", check.DataPath)
		return false
	}

	expected := check.ExpectedValue

	switch check.Operator {
	case "equals":
		return fmt.Sprintf("%v", value) == fmt.Sprintf("%v", expected)
	case "greaterThan":
		valueNum, valueErr := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
		expectedNum, expectedErr := strconv.ParseFloat(fmt.Sprintf("%v", expected), 64)
		if valueErr != nil || expectedErr != nil {
			return false
		}
		return valueNum > expectedNum
	case "contains":
		valueStr := fmt.Sprintf("%v", value)
		expectedStr := fmt.Sprintf("%v", expected)
		return strings.Contains(valueStr, expectedStr)
	default:
		logger.LogWarning("Unknown data check operator: %s", check.Operator)
		return false
	}
}

// GetContext returns the context for accessing data
func (r *TemplateResolver) GetContext() *UserContext {
	return r.context
}

// initializeDefaultSelectors sets up the default selector constant mappings
func (r *TemplateResolver) initializeDefaultSelectors() {
	r.selectorMap = SelectorMap
}

// ResolveSelector resolves a selector by checking if it's a constant first, then applying template resolution
func (r *TemplateResolver) ResolveSelector(selector string) string {
	if selector == "" {
		return selector
	}

	logger.LogDebug("Resolving selector: %s", selector)

	// Check if it's a selector constant (uppercase letters, numbers, underscores)
	if r.isSelectorConstant(selector) {
		if mappedSelector, exists := r.selectorMap[selector]; exists {
			logger.LogDebug("Mapped selector constant %s to %s", selector, mappedSelector)
			selector = mappedSelector
		} else {
			logger.LogWarning("Unknown selector constant: %s", selector)
		}
	}

	// Apply template resolution for any {{}} expressions
	resolvedSelector := r.Resolve(selector)

	logger.LogDebug("Final resolved selector: %s", resolvedSelector)
	return resolvedSelector
}

// isSelectorConstant checks if a string looks like a selector constant
func (r *TemplateResolver) isSelectorConstant(selector string) bool {
	// Selector constants should be uppercase with underscores
	if len(selector) == 0 {
		return false
	}

	for _, char := range selector {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}

	return true
}

// AddSelectorMapping adds or updates a selector constant mapping
func (r *TemplateResolver) AddSelectorMapping(constant, cssSelector string) {
	r.selectorMap[constant] = cssSelector
	logger.LogDebug("Added selector mapping: %s -> %s", constant, cssSelector)
}

// RemoveSelectorMapping removes a selector constant mapping
func (r *TemplateResolver) RemoveSelectorMapping(constant string) {
	delete(r.selectorMap, constant)
	logger.LogDebug("Removed selector mapping: %s", constant)
}

// GetAllSelectorMappings returns a copy of all current selector mappings
func (r *TemplateResolver) GetAllSelectorMappings() map[string]string {
	mappings := make(map[string]string)
	for k, v := range r.selectorMap {
		mappings[k] = v
	}
	return mappings
}

// HasSelectorMapping checks if a selector constant has a mapping
func (r *TemplateResolver) HasSelectorMapping(constant string) bool {
	_, exists := r.selectorMap[constant]
	return exists
}
