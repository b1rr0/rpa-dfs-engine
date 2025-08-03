# Traverser - Single Action

> Simple tree-based workflow automation with single-action nodes

## 🎯 **Overview**

Traverser executes JSON-defined workflows using single-action nodes. Each node performs exactly one action, making workflows easily parseable and maintainable.

### **Key Principles**
- **One Node = One Action** - No complex multi-action nodes
- **ForEach with Questions** - User interaction through simple questions  
- **Parseable Structure** - Clean JSON tree without complexity
- **Single User** - One workflow, one user, one browser

## 🏗️ **Architecture**

```
Workflow JSON → Engine → Browser Actions
     ↓
Single Actions:
├── moveToPage
├── fillField  
├── clickButton
├── sendFile
└── wait

Control Flow:
├── conditional
├── question
├── sequence
└── forEach (with questions)
```

## 📋 **Core Components**

### **1. Single Action Nodes**
Each node performs exactly one browser action:

```json
{
  "nodeType": "fillField",
  "selector": "#email",
  "value": "{{user.email}}",
  "next": {
    "nodeType": "clickButton",
    "selector": "#submit"
  }
}
```

### **2. Control Flow Nodes**
Manage execution flow without performing browser actions:

```json
{
  "nodeType": "conditional",
  "conditionExpression": "{{user.age}} > 18",
  "branches": {
    "yes": { /* adult flow */ },
    "no": { /* minor flow */ }
  }
}
```

### **3. ForEach with Questions**
Loop through arrays with user confirmation:

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

### **4. Simple Context**
User data for template resolution:

```json
{
  "user": {
    "email": "john@example.com",
    "files": ["/path/file1.pdf", "/path/file2.pdf"]
  }
}
```

## 🚀 **Quick Start**

### **1. Simple Login Workflow**
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

### **2. Context Data**
```json
{
  "user": {
    "email": "john@example.com",
    "password": "secret123"
  }
}
```

### **3. Execution**
```bash
rpa-engine --workflow=login.json --context=user.json
```

## 📊 **Node Types**

### **Action Nodes (Single Action)**

| Node Type | Purpose | Properties |
|-----------|---------|------------|
| `moveToPage` | Navigate to URL | `url` |
| `fillField` | Fill one field | `selector`, `value` |
| `clickButton` | Click element | `selector` |
| `sendFile` | Upload file | `selector`, `filePath` |
| `wait` | Pause execution | `duration` |

### **Control Nodes**

| Node Type | Purpose | Properties |
|-----------|---------|------------|
| `conditional` | Branch on condition | `conditionExpression`, `branches` |
| `question` | Branch on data check | `check`, `branches` |
| `sequence` | Execute nodes in order | `sequence` |
| `forEach` | Loop with questions | `dataSource`, `questionText` |

## ✅ **Single Action Rules**

### **✅ Correct (One Action Per Node)**
```json
{
  "nodeType": "fillField",
  "selector": "#email",
  "value": "{{user.email}}",
  "next": {
    "nodeType": "clickButton", 
    "selector": "#submit"
  }
}
```

### **❌ Wrong (Multiple Actions)**
```json
{
  "nodeType": "fillForm",
  "fields": [...],
  "submitButton": "#submit"  // FORBIDDEN!
}
```

**Use sequence instead:**
```json
{
  "nodeType": "sequence",
  "sequence": [
    {"nodeType": "fillField", "selector": "#field1", "value": "value1"},
    {"nodeType": "fillField", "selector": "#field2", "value": "value2"}, 
    {"nodeType": "clickButton", "selector": "#submit"}
  ]
}
```

## 🔄 **ForEach with Questions**

Instead of user approval nodes, use forEach with questions:

```json
{
  "nodeType": "forEach",
  "dataSource": "{{user.documents}}",
  "questionText": "Upload document {{iterator.count}} of {{iterator.total}}?",
  "next": {
    "nodeType": "sendFile",
    "selector": "input[type='file']",
    "filePath": "{{user.documents[iterator.index]}}"
  }
}
```

**User Interaction:**
- Question: "Upload document 1 of 3?" → User: Enter (continue) or 'n' (skip)
- Execute action or skip to next item
- Repeat for all items

## 📚 **Documentation Structure**

```
docs/traverser/
├── 01_OVERVIEW.md           # This overview
├── 02_NODE_REFERENCE.md     # Single action node types
├── 03_JSON_SCHEMA.md        # JSON schema definition  
├── 04_API_REFERENCE.md      # Go implementation
├── 05_EXECUTION_FLOW.md     # How nodes execute
├── 06_EXAMPLES.md           # Workflow examples
└── 07_USER_CONTEXT.md       # Context and templates
```

## 🎯 **Key Features**

### **1. Parseable Structure**
- Clean JSON tree structure
- One action per node
- Clear execution flow
- Easy to validate

### **2. Simple Templates**
- `{{user.propertyName}}` - User data
- `{{user.array[iterator.index]}}` - Array items
- `{{iterator.count}}` - Current iteration

### **3. User Questions**
- No complex approval dialogs
- Simple Enter/n responses
- Per-item confirmation in loops

### **4. Browser Integration**
- Uses existing Chrome installation
- Single browser instance
- Direct CSS selector actions

## 🛠️ **Implementation Strategy**

### **Phase 1: Core Engine**
1. JSON workflow parser
2. Single action node executor  
3. Simple context resolver
4. Basic browser integration

### **Phase 2: Control Flow**
1. Conditional branching
2. Data questions
3. Sequence execution
4. ForEach loops

### **Phase 3: Templates**
1. User data templates
2. Iterator variables
3. Context validation
4. Error handling

## 🎯 **Use Cases**

### **Web Form Automation**
- Login to websites
- Fill registration forms
- Submit applications
- Upload documents

### **Data Entry**
- Process CSV files
- Fill multiple forms
- Batch operations
- User confirmation per item

### **Testing Workflows**
- Automated UI testing
- Regression testing
- User journey validation
- Performance testing

## 🔗 **Integration Points**

### **With Existing Browser**
- Use current Chrome installation
- No browser management complexity
- Direct CSS selector targeting

### **With Go Application**
- Embedded workflow engine
- JSON configuration files
- Simple API integration
- Error handling callbacks

---

> 💡 **Simplicity**: One node = one action. ForEach with questions. Parseable structure for easy automation. 