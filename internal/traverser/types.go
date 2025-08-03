package traverser

// Workflow represents a complete workflow definition
type Workflow struct {
	Graph    *Node            `json:"graph"`
	Metadata WorkflowMetadata `json:"metadata"`
}

// WorkflowMetadata contains metadata about the workflow
type WorkflowMetadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// Node represents a single action or control flow node
type Node struct {
	NodeType string `json:"nodeType"`
	ID       string `json:"id,omitempty"`
	Next     *Node  `json:"next,omitempty"`

	// Navigation
	URL string `json:"url,omitempty"`

	// Single field actions
	Selector string `json:"selector,omitempty"`
	Value    string `json:"value,omitempty"`
	FilePath string `json:"filePath,omitempty"`

	// Conditional
	ConditionExpression string    `json:"conditionExpression,omitempty"`
	Branches            *Branches `json:"branches,omitempty"`

	// Question
	Check *DataCheck `json:"check,omitempty"`

	// Sequence
	Sequence []Node `json:"sequence,omitempty"`

	// ForEach
	DataSource              string `json:"dataSource,omitempty"`
	DataSourceIteratorParam string `json:"dataSourceIteratorParam,omitempty"`
	QuestionText            string `json:"questionText,omitempty"`

	// Wait
	Duration int `json:"duration,omitempty"`
}

// Branches represents conditional branches
type Branches struct {
	Yes *Node `json:"yes,omitempty"`
	No  *Node `json:"no,omitempty"`
}

// DataCheck represents a data validation check
type DataCheck struct {
	DataPath      string      `json:"dataPath"`
	Operator      string      `json:"operator"`
	ExpectedValue interface{} `json:"expectedValue"`
}

// UserContext represents the execution context with user data
type UserContext struct {
	data map[string]interface{}
}

// Browser interface for browser automation
type Browser interface {
	NavigateTo(url string) error
	FillField(selector, value string) error
	ClickElement(selector string) error
	SendFile(selector, filePath string) error
	Close() error
}
