# Traverser User Context - Single Action

> Simple context for single-action nodes and parseable workflows

## 🎯 **Context Overview**

Context holds user data for single-action node execution. Simple key-value data for template resolution in parseable workflows.

### **Basic Structure**
```json
{
  "user": {
    "email": "john@example.com",
    "password": "secret123",
    "documents": ["/path/to/file1.pdf", "/path/to/file2.pdf"]
  }
}
```

## 📊 **Context Types**

### **User Data**
All user-specific information under `user` key:

```json
{
  "user": {
    "email": "john@example.com",
    "password": "secret123",
    "firstName": "John",
    "lastName": "Doe",
    "phone": "+1234567890",
    "website": "https://example.com",
    "documents": [
      "/Users/john/resume.pdf",
      "/Users/john/photo.jpg"
    ],
    "products": ["laptop", "mouse", "keyboard"]
  }
}
```

### **Iterator Data**
Auto-managed during forEach loops:

```json
{
  "iterator": {
    "index": 0,     // 0-based index
    "count": 1,     // 1-based count
    "total": 3      // total items
  }
}
```

## 📝 **Template Usage**

### **Single Action Templates**
Access user data in single-action nodes:

```json
{
  "nodeType": "fillField",
  "selector": "#email",
  "value": "{{user.email}}"
}
```

```json
{
  "nodeType": "sendFile",
  "selector": "input[type='file']",
  "filePath": "{{user.documents[iterator.index]}}"
}
```

### **ForEach Templates**
Access loop data in forEach questions:

```json
{
  "nodeType": "forEach",
  "dataSource": "{{user.products}}",
  "questionText": "Search for {{user.products[iterator.index]}}? ({{iterator.count}}/{{iterator.total}})"
}
```

### **Nested Properties**
```json
{
  "user": {
    "profile": {
      "social": {
        "twitter": "@johndoe"
      }
    }
  }
}
```

Template: `{{user.profile.social.twitter}}` → "@johndoe"

## 🔧 **Context Loading**

### **From File**
```bash
rpa-engine --workflow=login.json --context=user.json
```

**user.json:**
```json
{
  "user": {
    "email": "john@example.com",
    "password": "secret123"
  }
}
```

### **Inline JSON**
```bash
rpa-engine --workflow=workflow.json \
  --context='{"user":{"email":"john@example.com","documents":["/path/file.pdf"]}}'
```

### **Environment Variables**
```bash
export USER_EMAIL="john@example.com"
export USER_PASSWORD="secret123"
rpa-engine --workflow=workflow.json --use-env
```

## 🔄 **Context in ForEach**

### **ForEach with Questions**
```json
{
  "nodeType": "forEach",
  "dataSource": "{{user.items}}",
  "questionText": "Process {{user.items[iterator.index]}} (item {{iterator.count}})?",
  "next": {
    "nodeType": "fillField",
    "selector": "#item-name",
    "value": "{{user.items[iterator.index]}}"
  }
}
```

**Context:**
```json
{
  "user": {
    "items": ["task1", "task2", "task3"]
  }
}
```

**Execution:**
- Question: "Process task1 (item 1)?" → User: Enter
- Execute fillField with "task1"
- Question: "Process task2 (item 2)?" → User: 'n' (skip)
- Question: "Process task3 (item 3)?" → User: Enter
- Execute fillField with "task3"

## 🎯 **Context Examples**

### **Login Context**
```json
{
  "user": {
    "email": "user@example.com",
    "password": "mypassword123"
  }
}
```

**Workflow:**
```json
{
  "nodeType": "fillField",
  "selector": "#email",
  "value": "{{user.email}}",
  "next": {
    "nodeType": "fillField",
    "selector": "#password",
    "value": "{{user.password}}"
  }
}
```

### **File Upload Context**
```json
{
  "user": {
    "documents": [
      "/Users/john/resume.pdf",
      "/Users/john/cover-letter.pdf"
    ]
  }
}
```

**Workflow:**
```json
{
  "nodeType": "forEach",
  "dataSource": "{{user.documents}}",
  "questionText": "Upload {{user.documents[iterator.index]}}?",
  "next": {
    "nodeType": "sendFile",
    "selector": "input[type='file']",
    "filePath": "{{user.documents[iterator.index]}}"
  }
}
```

### **Form Context**
```json
{
  "user": {
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@email.com",
    "phone": "+1-555-0123"
  }
}
```

**Workflow (Sequence):**
```json
{
  "nodeType": "sequence",
  "sequence": [
    {
      "nodeType": "fillField",
      "selector": "#firstName",
      "value": "{{user.firstName}}"
    },
    {
      "nodeType": "fillField",
      "selector": "#lastName",
      "value": "{{user.lastName}}"
    },
    {
      "nodeType": "fillField",
      "selector": "#email",
      "value": "{{user.email}}"
    }
  ]
}
```

### **Shopping Context**
```json
{
  "user": {
    "products": ["laptop", "mouse", "keyboard"],
    "website": "https://shop.example.com"
  }
}
```

**Workflow:**
```json
{
  "nodeType": "moveToPage",
  "url": "{{user.website}}",
  "next": {
    "nodeType": "forEach",
    "dataSource": "{{user.products}}",
    "questionText": "Search for {{user.products[iterator.index]}}?",
    "next": {
      "nodeType": "fillField",
      "selector": "#search",
      "value": "{{user.products[iterator.index]}}"
    }
  }
}
```

## 🎛️ **Simple API**

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

func (c *Context) SetIterator(index, total int) {
    c.iterator = map[string]interface{}{
        "index": index,
        "count": index + 1,
        "total": total,
    }
}
```

## 📋 **Best Practices**

### **1. Simple Structure**
```json
{
  "user": {
    "email": "user@example.com",
    "password": "secret",
    "files": ["/path/file1.pdf", "/path/file2.pdf"]
  }
}
```

### **2. Array for ForEach**
```json
{
  "user": {
    "tasks": ["task1", "task2", "task3"],
    "documents": ["/file1.pdf", "/file2.pdf"]
  }
}
```

### **3. Absolute Paths**
```json
{
  "user": {
    "profilePhoto": "/Users/john/Documents/photo.jpg",
    "resume": "/Users/john/Desktop/resume.pdf"
  }
}
```

### **4. Boolean Values**
```json
{
  "user": {
    "isActive": true,
    "hasPermission": false
  }
}
```

## ❌ **What's NOT Supported**

- Complex state management
- Multi-user contexts
- Dynamic context modification
- External data sources
- Context inheritance
- Encrypted values
- User approval messages (use forEach questions instead)

## 🎯 **Complete Example**

### **Single Action Workflow**
**workflow.json:**
```json
{
  "graph": {
    "nodeType": "moveToPage",
    "url": "{{user.website}}",
    "next": {
      "nodeType": "fillField",
      "selector": "#search",
      "value": "{{user.searchTerm}}",
      "next": {
        "nodeType": "clickButton",
        "selector": "#search-btn",
        "next": {
          "nodeType": "wait",
          "duration": 2000,
          "next": null
        }
      }
    }
  }
}
```

**context.json:**
```json
{
  "user": {
    "website": "https://google.com",
    "searchTerm": "golang tutorial"
  }
}
```

### **ForEach Workflow**
**workflow.json:**
```json
{
  "graph": {
    "nodeType": "forEach",
    "dataSource": "{{user.files}}",
    "questionText": "Upload {{user.files[iterator.index]}} ({{iterator.count}}/{{iterator.total}})?",
    "next": {
      "nodeType": "sendFile",
      "selector": "input[type='file']",
      "filePath": "{{user.files[iterator.index]}}",
      "next": null
    }
  }
}
```

**context.json:**
```json
{
  "user": {
    "files": [
      "/Users/john/doc1.pdf",
      "/Users/john/doc2.pdf",
      "/Users/john/doc3.pdf"
    ]
  }
}
```

**Execution:**
```bash
rpa-engine --workflow=workflow.json --context=context.json
```

---

> 💡 **Simplicity**: One node = one action. ForEach with questions instead of approval. Parseable structure. 