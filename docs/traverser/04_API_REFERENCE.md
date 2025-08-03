# Traverser API Reference - Single Action

> Simple Go API with single action per node rule

## 🎯 **Core Types**

### **Workflow Definition**
```go
type Workflow struct {
    Graph    *Node              `json:"graph"`
    Metadata WorkflowMetadata   `json:"metadata"`
}

type WorkflowMetadata struct {
    Name        string `json:"name"`
    Version     string `json:"version"`
    Description string `json:"description"`
}
```

### **Node Structure**
```go
type Node struct {
    NodeType    string      `json:"nodeType"`
    ID          string      `json:"id,omitempty"`
    Next        *Node       `json:"next,omitempty"`
    
    // Navigation
    URL         string      `json:"url,omitempty"`
    
    // Single field
    Selector    string      `json:"selector,omitempty"`
    Value       string      `json:"value,omitempty"`
    FilePath    string      `json:"filePath,omitempty"`
    
    // Conditional
    ConditionExpression string   `json:"conditionExpression,omitempty"`
    Branches           *Branches `json:"branches,omitempty"`
    
    // Question
    Check        *DataCheck `json:"check,omitempty"`
    
    // Sequence
    Sequence []Node `json:"sequence,omitempty"`
    
    // ForEach
    DataSource     string `json:"dataSource,omitempty"`
    QuestionText   string `json:"questionText,omitempty"`
    
    // Wait
    Duration int `json:"duration,omitempty"`
}

type Branches struct {
    Yes *Node `json:"yes,omitempty"`
    No  *Node `json:"no,omitempty"`
}

type DataCheck struct {
    DataPath      string      `json:"dataPath"`
    Operator      string      `json:"operator"`
    ExpectedValue interface{} `json:"expectedValue"`
}
```

## 🚀 **Simple Execution Engine**

### **Main Engine**
```go
type Engine struct {
    workflow *Workflow
    context  *Context
    browser  *Browser
    logger   Logger
}

func NewEngine(browserPath string) *Engine {
    return &Engine{
        browser: NewBrowser(browserPath),
        logger:  NewSimpleLogger(),
    }
}

func (e *Engine) LoadWorkflow(workflowPath string) error {
    data, err := os.ReadFile(workflowPath)
    if err != nil {
        return err
    }
    return json.Unmarshal(data, &e.workflow)
}

func (e *Engine) SetContext(userData map[string]interface{}) {
    e.context = NewContext(userData)
}

func (e *Engine) Execute() error {
    e.logger.Info("Starting: %s", e.workflow.Metadata.Name)
    return e.executeNode(e.workflow.Graph)
}

func (e *Engine) executeNode(node *Node) error {
    if node == nil {
        return nil
    }
    
    e.logger.Debug("Executing: %s", node.NodeType)
    
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
        return fmt.Errorf("unknown node: %s", node.NodeType)
    }
}
```

## 🎛️ **Node Execution**

### **Basic Actions**
```go
func (e *Engine) executeMoveToPage(node *Node) error {
    url := e.resolveString(node.URL)
    e.logger.Info("Navigate to: %s", url)
    
    if err := e.browser.NavigateTo(url); err != nil {
        return err
    }
    
    return e.executeNode(node.Next)
}

func (e *Engine) executeFillField(node *Node) error {
    selector := e.resolveString(node.Selector)
    value := e.resolveString(node.Value)
    e.logger.Info("Fill field %s: %s", selector, value)
    
    if err := e.browser.FillField(selector, value); err != nil {
        return err
    }
    
    return e.executeNode(node.Next)
}

func (e *Engine) executeClickButton(node *Node) error {
    selector := e.resolveString(node.Selector)
    e.logger.Info("Click: %s", selector)
    
    if err := e.browser.ClickButton(selector); err != nil {
        return err
    }
    
    return e.executeNode(node.Next)
}

func (e *Engine) executeSendFile(node *Node) error {
    selector := e.resolveString(node.Selector)
    filePath := e.resolveString(node.FilePath)
    e.logger.Info("Upload %s to %s", filePath, selector)
    
    if err := e.browser.UploadFile(selector, filePath); err != nil {
        return err
    }
    
    return e.executeNode(node.Next)
}
```

### **Control Flow**
```go
func (e *Engine) executeConditional(node *Node) error {
    result := e.evaluateCondition(node.ConditionExpression)
    
    if result {
        e.logger.Debug("Condition true")
        return e.executeNode(node.Branches.Yes)
    } else {
        e.logger.Debug("Condition false")
        return e.executeNode(node.Branches.No)
    }
}

func (e *Engine) executeQuestion(node *Node) error {
    value, exists := e.context.Get(node.Check.DataPath)
    if !exists {
        return fmt.Errorf("data not found: %s", node.Check.DataPath)
    }
    
    result := e.compareValues(value, node.Check.Operator, node.Check.ExpectedValue)
    
    if result {
        return e.executeNode(node.Branches.Yes)
    } else {
        return e.executeNode(node.Branches.No)
    }
}

func (e *Engine) executeSequence(node *Node) error {
    for i, seqNode := range node.Sequence {
        e.logger.Debug("Sequence %d/%d", i+1, len(node.Sequence))
        err := e.executeNode(&seqNode)
        if err != nil {
            return err
        }
    }
    
    return e.executeNode(node.Next)
}

func (e *Engine) executeForEach(node *Node) error {
    data, exists := e.context.Get(node.DataSource)
    if !exists {
        return fmt.Errorf("data source not found: %s", node.DataSource)
    }
    
    arr, ok := data.([]interface{})
    if !ok {
        return fmt.Errorf("data source is not array: %s", node.DataSource)
    }
    
    for i := 0; i < len(arr); i++ {
        e.context.SetIterator(i, len(arr))
        
        // Ask user
        question := e.resolveString(node.QuestionText)
        e.logger.Info("=== QUESTION ===")
        e.logger.Info(question)
        e.logger.Info("Press Enter to continue, 'n' to skip...")
        
        var response string
        fmt.Scanln(&response)
        
        if response == "n" || response == "no" {
            e.logger.Info("Skipped item %d", i)
            continue
        }
        
        // Execute action for this item
        err := e.executeNode(node.Next)
        if err != nil {
            return err
        }
    }
    
    return nil
}

func (e *Engine) executeWait(node *Node) error {
    e.logger.Debug("Waiting %d ms", node.Duration)
    time.Sleep(time.Duration(node.Duration) * time.Millisecond)
    return e.executeNode(node.Next)
}
```

## 🔄 **Simple Context**

### **Context Implementation**
```go
type Context struct {
    user     map[string]interface{}
    iterator map[string]interface{}
}

func NewContext(userData map[string]interface{}) *Context {
    return &Context{
        user:     userData,
        iterator: make(map[string]interface{}),
    }
}

func (c *Context) Get(path string) (interface{}, bool) {
    if strings.HasPrefix(path, "user.") {
        return c.getFromMap(c.user, strings.TrimPrefix(path, "user."))
    }
    if strings.HasPrefix(path, "iterator.") {
        return c.getFromMap(c.iterator, strings.TrimPrefix(path, "iterator."))
    }
    return nil, false
}

func (c *Context) getFromMap(data map[string]interface{}, path string) (interface{}, bool) {
    keys := strings.Split(path, ".")
    current := data
    
    for i, key := range keys {
        if i == len(keys)-1 {
            val, exists := current[key]
            return val, exists
        }
        
        if next, ok := current[key].(map[string]interface{}); ok {
            current = next
        } else {
            return nil, false
        }
    }
    
    return nil, false
}

func (c *Context) SetIterator(index, total int) {
    c.iterator = map[string]interface{}{
        "index": index,
        "count": index + 1,
        "total": total,
    }
}
```

## 🌐 **Simple Browser**

### **Browser Interface**
```go
type Browser struct {
    ctx    context.Context
    cancel context.CancelFunc
}

func NewBrowser(chromePath string) *Browser {
    opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.ExecPath(chromePath),
        chromedp.Flag("headless", false),
    )
    
    allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
    ctx, cancel := chromedp.NewContext(allocCtx)
    
    return &Browser{
        ctx:    ctx,
        cancel: cancel,
    }
}

func (b *Browser) NavigateTo(url string) error {
    return chromedp.Run(b.ctx, chromedp.Navigate(url))
}

func (b *Browser) FillField(selector, value string) error {
    return chromedp.Run(b.ctx,
        chromedp.Clear(selector),
        chromedp.SendKeys(selector, value),
    )
}

func (b *Browser) ClickButton(selector string) error {
    return chromedp.Run(b.ctx, chromedp.Click(selector))
}

func (b *Browser) UploadFile(selector, filePath string) error {
    return chromedp.Run(b.ctx, chromedp.SendKeys(selector, filePath))
}

func (b *Browser) Close() {
    b.cancel()
}
```

## 📝 **Template Resolution**

### **Template Resolver**
```go
func (e *Engine) resolveString(template string) string {
    re := regexp.MustCompile(`\{\{([^}]+)\}\}`)
    
    return re.ReplaceAllStringFunc(template, func(match string) string {
        expr := strings.Trim(match, "{}")
        expr = strings.TrimSpace(expr)
        
        if value, exists := e.context.Get(expr); exists {
            return fmt.Sprintf("%v", value)
        }
        
        return match
    })
}

func (e *Engine) evaluateCondition(expr string) bool {
    resolved := e.resolveString(expr)
    
    if strings.Contains(resolved, ">") {
        parts := strings.Split(resolved, ">")
        if len(parts) == 2 {
            left, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
            right, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
            return left > right
        }
    }
    
    if strings.Contains(resolved, "==") {
        parts := strings.Split(resolved, "==")
        if len(parts) == 2 {
            left := strings.TrimSpace(parts[0])
            right := strings.TrimSpace(strings.Trim(parts[1], `"`))
            return left == right
        }
    }
    
    return false
}

func (e *Engine) compareValues(actual interface{}, operator string, expected interface{}) bool {
    switch operator {
    case "equals":
        return actual == expected
    case "greaterThan":
        if a, ok := actual.(float64); ok {
            if b, ok := expected.(float64); ok {
                return a > b
            }
        }
    case "contains":
        if a, ok := actual.(string); ok {
            if b, ok := expected.(string); ok {
                return strings.Contains(a, b)
            }
        }
    }
    return false
}
```

## 📊 **Simple Logging**

### **Logger Implementation**
```go
type Logger interface {
    Info(format string, args ...interface{})
    Debug(format string, args ...interface{})
    Error(format string, args ...interface{})
}

type SimpleLogger struct{}

func NewSimpleLogger() *SimpleLogger {
    return &SimpleLogger{}
}

func (l *SimpleLogger) Info(format string, args ...interface{}) {
    fmt.Printf("[INFO] "+format+"\n", args...)
}

func (l *SimpleLogger) Debug(format string, args ...interface{}) {
    fmt.Printf("[DEBUG] "+format+"\n", args...)
}

func (l *SimpleLogger) Error(format string, args ...interface{}) {
    fmt.Printf("[ERROR] "+format+"\n", args...)
}
```

## 🎯 **Usage Example**

### **Simple Usage**
```go
func main() {
    engine := NewEngine("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome")
    
    if err := engine.LoadWorkflow("workflow.json"); err != nil {
        log.Fatal(err)
    }
    
    engine.SetContext(map[string]interface{}{
        "user": map[string]interface{}{
            "email":    "user@example.com",
            "password": "secret123",
            "documents": []string{
                "/path/to/doc1.pdf",
                "/path/to/doc2.pdf",
            },
        },
    })
    
    if err := engine.Execute(); err != nil {
        log.Fatal(err)
    }
}
```

### **Single Action Workflow**
```json
{
  "graph": {
    "nodeType": "moveToPage",
    "url": "https://example.com/login",
    "next": {
      "nodeType": "fillField",
      "selector": "#email",
      "value": "{{user.email}}",
      "next": {
        "nodeType": "fillField",
        "selector": "#password",
        "value": "{{user.password}}",
        "next": {
          "nodeType": "clickButton",
          "selector": "#login-btn",
          "next": null
        }
      }
    }
  },
  "metadata": {
    "name": "Simple Login",
    "version": "1.0.0"
  }
}
```

### **ForEach Workflow**
```json
{
  "graph": {
    "nodeType": "forEach",
    "dataSource": "{{user.documents}}",
    "questionText": "Upload document {{iterator.count}} of {{iterator.total}}?",
    "next": {
      "nodeType": "sendFile",
      "selector": "input[type='file']",
      "filePath": "{{user.documents[iterator.index]}}",
      "next": null
    }
  }
}
```

---

> 💡 **Simplicity**: Single thread, single user, simple context, only Enter for confirmation, uses existing browser.