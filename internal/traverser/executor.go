package traverser

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"rpa-dfs-engine/internal/logger"
)

// Executor orchestrates the execution of workflow nodes
type Executor struct {
	workflow *Workflow
	context  *UserContext
	browser  Browser
	resolver *TemplateResolver
	reader   *bufio.Reader
}

// NewExecutor creates a new executor instance
func NewExecutor(workflow *Workflow, context *UserContext, browser Browser) *Executor {
	resolver := NewTemplateResolver(context)
	reader := bufio.NewReader(os.Stdin)

	return &Executor{
		workflow: workflow,
		context:  context,
		browser:  browser,
		resolver: resolver,
		reader:   reader,
	}
}

// Execute starts the workflow execution
func (e *Executor) Execute() error {
	logger.LogInfo("Starting workflow execution: %s", e.workflow.Metadata.Name)
	logger.LogInfo("Version: %s", e.workflow.Metadata.Version)

	if e.workflow.Metadata.Description != "" {
		logger.LogInfo("Description: %s", e.workflow.Metadata.Description)
	}

	defer func() {
		if e.browser != nil {
			e.browser.Close()
		}
	}()

	return e.executeNode(e.workflow.Graph)
}

// executeNode executes a single node based on its type
func (e *Executor) executeNode(node *Node) error {
	if node == nil {
		logger.LogInfo("Workflow execution completed")
		return nil
	}

	logger.LogDebug("Executing node: %s", node.NodeType)

	switch node.NodeType {
	case "moveToPage":
		return e.executeMoveToPage(node)
	case "fillField":
		return e.executeFillField(node)
	case "clickButton":
		return e.executeClickButton(node)
	case "sendFile":
		return e.executeSendFile(node)
	case "conditional":
		return e.executeConditional(node)
	case "question":
		return e.executeQuestion(node)
	case "sequence":
		return e.executeSequence(node)
	case "forEach":
		return e.executeForEach(node)
	case "wait":
		return e.executeWait(node)
	default:
		return fmt.Errorf("unknown node type: %s", node.NodeType)
	}
}

// executeMoveToPage navigates to a URL
func (e *Executor) executeMoveToPage(node *Node) error {
	url := e.resolver.Resolve(node.URL)

	if err := e.browser.NavigateTo(url); err != nil {
		return fmt.Errorf("moveToPage failed: %w", err)
	}

	return e.executeNode(node.Next)
}

// executeFillField fills a form field
func (e *Executor) executeFillField(node *Node) error {
	selector := e.resolver.ResolveSelector(node.Selector)
	value := e.resolver.Resolve(node.Value)

	if err := e.browser.FillField(selector, value); err != nil {
		return fmt.Errorf("fillField failed: %w", err)
	}

	return e.executeNode(node.Next)
}

// executeClickButton clicks a button or element
func (e *Executor) executeClickButton(node *Node) error {
	selector := e.resolver.ResolveSelector(node.Selector)

	if err := e.browser.ClickElement(selector); err != nil {
		return fmt.Errorf("clickButton failed: %w", err)
	}

	return e.executeNode(node.Next)
}

// executeSendFile uploads a file
func (e *Executor) executeSendFile(node *Node) error {
	selector := e.resolver.ResolveSelector(node.Selector)
	filePath := e.resolver.Resolve(node.FilePath)

	if err := e.browser.SendFile(selector, filePath); err != nil {
		return fmt.Errorf("sendFile failed: %w", err)
	}

	return e.executeNode(node.Next)
}

// executeConditional handles conditional branching
func (e *Executor) executeConditional(node *Node) error {
	condition := e.resolver.EvaluateCondition(node.ConditionExpression)

	logger.LogDebug("Conditional result: %t for expression: %s", condition, node.ConditionExpression)

	if condition {
		logger.LogInfo("Taking YES branch")
		return e.executeNode(node.Branches.Yes)
	} else {
		logger.LogInfo("Taking NO branch")
		return e.executeNode(node.Branches.No)
	}
}

// executeQuestion handles data-based questions
func (e *Executor) executeQuestion(node *Node) error {
	result := e.resolver.EvaluateDataCheck(node.Check)

	logger.LogDebug("Question result: %t for check: %+v", result, node.Check)

	if result {
		logger.LogInfo("Taking YES branch")
		return e.executeNode(node.Branches.Yes)
	} else {
		logger.LogInfo("Taking NO branch")
		return e.executeNode(node.Branches.No)
	}
}

// executeSequence executes a sequence of nodes
func (e *Executor) executeSequence(node *Node) error {
	logger.LogInfo("Executing sequence with %d nodes", len(node.Sequence))

	for i, child := range node.Sequence {
		logger.LogDebug("Executing sequence item %d/%d", i+1, len(node.Sequence))
		if err := e.executeNode(&child); err != nil {
			return fmt.Errorf("sequence item %d failed: %w", i+1, err)
		}
	}

	logger.LogSuccess("Sequence completed successfully")
	return e.executeNode(node.Next)
}

// executeForEach executes a loop with user questions
func (e *Executor) executeForEach(node *Node) error {

	iteratorParam, exists := e.context.data[node.DataSourceIteratorParam].([]any)

	if !exists {
		return fmt.Errorf("dataSourceIteratorParam not found: %s", node.DataSourceIteratorParam)
	}
	for i := range iteratorParam {
		e.context.data[node.DataSourceIteratorParam] = iteratorParam[i]
		if err := e.executeNode(node.Next); err != nil {
			return fmt.Errorf("forEach item %d failed: %w", err)
		}
	}
	e.context.data[node.DataSourceIteratorParam] = nil
	logger.LogSuccess("ForEach loop completed")
	return nil
}

// executeWait pauses execution
func (e *Executor) executeWait(node *Node) error {
	duration := time.Duration(node.Duration) * time.Millisecond
	logger.LogInfo("Waiting for %v", duration)

	time.Sleep(duration)

	logger.LogSuccess("Wait completed")
	return e.executeNode(node.Next)
}

// GetWorkflowInfo returns information about the loaded workflow
func (e *Executor) GetWorkflowInfo() *WorkflowMetadata {
	if e.workflow != nil {
		return &e.workflow.Metadata
	}
	return nil
}
