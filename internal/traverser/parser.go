package traverser

import (
	"encoding/json"
	"fmt"
	"os"

	"rpa-dfs-engine/internal/logger"
)

// Parser handles loading and validating workflow JSON files
type Parser struct{}

// NewParser creates a new parser instance
func NewParser() *Parser {
	return &Parser{}
}

// LoadWorkflow loads a workflow from a JSON file
func (p *Parser) LoadWorkflow(workflowPath string) (*Workflow, error) {
	logger.LogInfo("Loading workflow from: %s", workflowPath)

	data, err := os.ReadFile(workflowPath)
	if err != nil {
		logger.LogError("Failed to read workflow file: %v", err)
		return nil, fmt.Errorf("failed to read workflow file: %w", err)
	}

	var workflow Workflow
	if err := json.Unmarshal(data, &workflow); err != nil {
		logger.LogError("Failed to parse workflow JSON: %v", err)
		return nil, fmt.Errorf("failed to parse workflow JSON: %w", err)
	}

	// Validate the workflow
	if err := p.validateWorkflow(&workflow); err != nil {
		logger.LogError("Workflow validation failed: %v", err)
		return nil, fmt.Errorf("workflow validation failed: %w", err)
	}

	logger.LogSuccess("Workflow loaded successfully: %s", workflow.Metadata.Name)
	return &workflow, nil
}

// LoadContext loads user context from a JSON file
func (p *Parser) LoadContext(contextPath string) (*UserContext, error) {
	logger.LogInfo("Loading context from: %s", contextPath)

	data, err := os.ReadFile(contextPath)
	if err != nil {
		logger.LogError("Failed to read context file: %v", err)
		return nil, fmt.Errorf("failed to read context file: %w", err)
	}

	var contextData map[string]interface{}
	if err := json.Unmarshal(data, &contextData); err != nil {
		logger.LogError("Failed to parse context JSON: %v", err)
		return nil, fmt.Errorf("failed to parse context JSON: %w", err)
	}

	context := NewContext(contextData)
	logger.LogSuccess("Context loaded successfully")
	return context, nil
}

// validateWorkflow validates the workflow structure
func (p *Parser) validateWorkflow(workflow *Workflow) error {
	if workflow.Graph == nil {
		return fmt.Errorf("workflow must have a graph")
	}

	return p.validateNode(workflow.Graph)
}

// validateNode validates a single node and its children
func (p *Parser) validateNode(node *Node) error {
	if node == nil {
		return nil
	}

	// Validate node type
	validNodeTypes := map[string]bool{
		"moveToPage":  true,
		"fillField":   true,
		"clickButton": true,
		"sendFile":    true,
		"conditional": true,
		"question":    true,
		"sequence":    true,
		"forEach":     true,
		"wait":        true,
	}

	if !validNodeTypes[node.NodeType] {
		return fmt.Errorf("invalid node type: %s", node.NodeType)
	}

	// Validate required fields based on node type
	switch node.NodeType {
	case "moveToPage":
		if node.URL == "" {
			return fmt.Errorf("moveToPage node requires URL")
		}
	case "fillField":
		if node.Selector == "" || node.Value == "" {
			return fmt.Errorf("fillField node requires selector and value")
		}
	case "clickButton":
		if node.Selector == "" {
			return fmt.Errorf("clickButton node requires selector")
		}
	case "sendFile":
		if node.Selector == "" || node.FilePath == "" {
			return fmt.Errorf("sendFile node requires selector and filePath")
		}
	case "conditional":
		if node.ConditionExpression == "" || node.Branches == nil {
			return fmt.Errorf("conditional node requires conditionExpression and branches")
		}
		if err := p.validateNode(node.Branches.Yes); err != nil {
			return fmt.Errorf("conditional yes branch: %w", err)
		}
		if err := p.validateNode(node.Branches.No); err != nil {
			return fmt.Errorf("conditional no branch: %w", err)
		}
	case "question":
		if node.Check == nil || node.Branches == nil {
			return fmt.Errorf("question node requires check and branches")
		}
		if err := p.validateNode(node.Branches.Yes); err != nil {
			return fmt.Errorf("question yes branch: %w", err)
		}
		if err := p.validateNode(node.Branches.No); err != nil {
			return fmt.Errorf("question no branch: %w", err)
		}
	case "sequence":
		if len(node.Sequence) == 0 {
			return fmt.Errorf("sequence node requires at least one child node")
		}
		for i, child := range node.Sequence {
			if err := p.validateNode(&child); err != nil {
				return fmt.Errorf("sequence child %d: %w", i, err)
			}
		}
	case "forEach":
		if node.DataSourceIteratorParam == "" {
			return fmt.Errorf("forEach node requires dataSourceIteratorParam")
		}
	case "wait":
		if node.Duration <= 0 {
			return fmt.Errorf("wait node requires positive duration")
		}
	}

	return p.validateNode(node.Next)
}
