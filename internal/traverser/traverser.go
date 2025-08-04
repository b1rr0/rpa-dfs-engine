package traverser

import (
	"fmt"

	"rpa-dfs-engine/internal/logger"
)

// Traverser provides a high-level API for workflow execution
type Traverser struct {
	parser   *Parser
	browser  Browser
	executor *Executor
}

func NewWithBrowser(browser Browser) *Traverser {
	return &Traverser{
		parser:  NewParser(),
		browser: browser,
	}
}

func (t *Traverser) ExecuteWorkflowWithContext(workflowPath string, context *UserContext) error {
	logger.LogInfo("=== Traverser Workflow Execution ===")

	workflow, err := t.parser.LoadWorkflow(workflowPath)
	if err != nil {
		return fmt.Errorf("failed to load workflow: %w", err)
	}

	t.executor = NewExecutor(workflow, context, t.browser)
	return t.executor.Execute()
}

func (t *Traverser) Close() error {
	if t.browser != nil {
		return t.browser.Close()
	}
	return nil
}

// New creates a new traverser with default browser
func New() *Traverser {
	return &Traverser{
		parser: NewParser(),
	}
}

// ExecuteWorkflow overloaded version that takes context file path
func (t *Traverser) ExecuteWorkflow(workflowPath, contextPath string) error {
	context, err := t.parser.LoadContext(contextPath)
	if err != nil {
		return fmt.Errorf("failed to load context: %w", err)
	}

	return t.ExecuteWorkflowWithContext(workflowPath, context)
}

// Validate validates a workflow file
func Validate(workflowPath string) error {
	parser := NewParser()
	_, err := parser.LoadWorkflow(workflowPath)
	return err
}
