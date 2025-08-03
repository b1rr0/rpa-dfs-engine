# Traverser Execution Flow - Single Action

> Simple single-action node execution with parseable structure

## 🎯 **Simple Algorithm**

One node = one action. No complex state, just simple tree traversal.

```
ExecuteNode(node):
    1. Check if node is null → stop
    2. Execute single action
    3. Follow next pointer or branch
    4. Continue until null
```

## 🚀 **Simple Engine**

### **Basic Execution**
```go
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

## 🎛️ **Single Action Execution**

### **Basic Actions (One Action Per Node)**
```go
func (e *Engine) executeMoveToPage(node *Node) error {
    url := e.resolveString(node.URL)
    e.logger.Info("Navigate to: %s", url)
    
    err := e.browser.NavigateTo(url)
    if err != nil {
        return err
    }
    
    return e.executeNode(node.Next)
}

func (e *Engine) executeFillField(node *Node) error {
    selector := e.resolveString(node.Selector)
    value := e.resolveString(node.Value)
    e.logger.Info("Fill %s: %s", selector, value)
    
    err := e.browser.FillField(selector, value)
    if err != nil {
        return err
    }
    
    return e.executeNode(node.Next)
}

func (e *Engine) executeClickButton(node *Node) error {
    selector := e.resolveString(node.Selector)
    e.logger.Info("Click: %s", selector)
    
    err := e.browser.ClickButton(selector)
    if err != nil {
        return err
    }
    
    return e.executeNode(node.Next)
}

func (e *Engine) executeSendFile(node *Node) error {
    selector := e.resolveString(node.Selector)
    filePath := e.resolveString(node.FilePath)
    e.logger.Info("Upload %s to %s", filePath, selector)
    
    err := e.browser.UploadFile(selector, filePath)
    if err != nil {
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
        
        // Ask user question
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

## 📝 **Template Resolution**

### **Simple Templates**
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
```

## 🔄 **Execution Patterns**

### **Linear Chain (Single Actions)**
```
moveToPage → fillField → fillField → clickButton → null
```
1. Execute moveToPage
2. Follow next to fillField (email)
3. Execute fillField (email)
4. Follow next to fillField (password)
5. Execute fillField (password)  
6. Follow next to clickButton
7. Execute clickButton
8. next is null → stop

### **Sequence Block**
```
sequence [
    fillField (field1),
    fillField (field2),
    clickButton (submit)
] → next → null
```
1. Execute sequence
2. Execute fillField (field1)
3. Execute fillField (field2)
4. Execute clickButton (submit)
5. Follow sequence.next
6. Sequence complete

### **Conditional Branch**
```
conditional
├─ yes → fillField → null
└─ no → fillField → null
```
1. Execute conditional
2. Evaluate condition
3. If true: execute yes branch (fillField)
4. If false: execute no branch (fillField)
5. Branch complete

### **ForEach with User Question**
```
forEach (user.items)
├─ question: "Process item 1?"
├─ user answer: Enter or 'n'
└─ if Enter: execute action
   if 'n': skip to next item
```
1. Execute forEach
2. Set iterator context (0, 3)
3. Ask user: "Process item 1 of 3?"
4. User presses Enter → execute action
5. Set iterator context (1, 3)
6. Ask user: "Process item 2 of 3?"
7. User types 'n' → skip action
8. Set iterator context (2, 3)
9. Ask user: "Process item 3 of 3?"
10. User presses Enter → execute action
11. forEach complete

## 📊 **Simple Context**

### **Context Structure**
```go
type Context struct {
    user     map[string]interface{} // {{user.*}}
    iterator map[string]interface{} // {{iterator.*}}
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

func (c *Context) SetIterator(index, total int) {
    c.iterator = map[string]interface{}{
        "index": index,     // 0-based
        "count": index + 1, // 1-based
        "total": total,     // total items
    }
}
```

### **Context Usage**
```json
{
  "user": {
    "email": "john@example.com",
    "password": "secret123",
    "documents": ["file1.pdf", "file2.pdf"]
  }
}
```

Templates:
- `{{user.email}}` → "john@example.com"
- `{{user.documents[iterator.index]}}` → "file1.pdf" (first iteration)
- `{{iterator.count}}` → 1 (first iteration)

## 🎯 **Simple Example**

### **Workflow JSON**
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
  }
}
```

### **Execution Flow**
1. Start at `moveToPage`
2. Navigate to "https://example.com/login"
3. Follow next to `fillField` (email)
4. Fill "#email" with "john@example.com"
5. Follow next to `fillField` (password)
6. Fill "#password" with "secret123"
7. Follow next to `clickButton`
8. Click "#login-btn"
9. next is null → execution complete

## ✅ **Single Action Rules**

### **✅ Correct Flow**
```
moveToPage → fillField → clickButton → null
```
Each node does exactly one action.

### **❌ Wrong Flow**
```
fillFormAndSubmit → null  // FORBIDDEN: Multiple actions
```

**Use sequence instead:**
```
sequence [
    fillField,
    fillField,
    clickButton
] → null
```

## 🔄 **ForEach with Questions**

### **ForEach Execution**
```json
{
  "nodeType": "forEach",
  "dataSource": "{{user.files}}",
  "questionText": "Upload {{user.files[iterator.index]}}?",
  "next": {
    "nodeType": "sendFile",
    "selector": "input[type='file']",
    "filePath": "{{user.files[iterator.index]}}"
  }
}
```

**Execution:**
1. Question: "Upload file1.pdf?" → User: Enter → Upload file1.pdf
2. Question: "Upload file2.pdf?" → User: 'n' → Skip
3. Question: "Upload file3.pdf?" → User: Enter → Upload file3.pdf
4. forEach complete

## 📋 **Node Types Summary**

### **Action Nodes (Single Action)**
- `moveToPage` - Navigate to URL
- `fillField` - Fill one field
- `clickButton` - Click element
- `sendFile` - Upload file
- `wait` - Pause

### **Control Nodes**
- `conditional` - Branch on condition
- `question` - Branch on data check
- `sequence` - Execute nodes in order
- `forEach` - Loop with user questions

---

> 💡 **Simplicity**: One node = one action. ForEach with user questions. No complex states. 