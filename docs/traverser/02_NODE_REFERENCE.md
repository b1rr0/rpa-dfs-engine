# Traverser Node Reference - Single Action

> One node = one action. No double actions. Parseable structure.

## 📋 **Node Structure**

All nodes follow this simple structure:

```json
{
  "nodeType": "actionName",
  "id": "optional-id",
  "next": {
    // Next node or null
  }
}
```

**RULE: One node = one action only!**

## 🎯 **Action Nodes**

### **moveToPage**
Navigate to URL.

```json
{
  "nodeType": "moveToPage",
  "url": "{{user.website}}",
  "next": {
    "nodeType": "fillField",
    "selector": "#email",
    "value": "{{user.email}}"
  }
}
```

**Properties:**
- `url` (string): Target URL
- `next` (node|null): Next node

### **fillField**
Fill one field.

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

**Properties:**
- `selector` (string): CSS selector
- `value` (string): Value to fill
- `next` (node|null): Next node

### **clickButton**
Click button or element.

```json
{
  "nodeType": "clickButton",
  "selector": "#login-btn",
  "next": {
    "nodeType": "wait",
    "duration": 2000
  }
}
```

**Properties:**
- `selector` (string): Element selector
- `next` (node|null): Next node

### **sendFile**
Upload file.

```json
{
  "nodeType": "sendFile",
  "selector": "input[type='file']",
  "filePath": "{{user.documentPath}}",
  "next": null
}
```

**Properties:**
- `selector` (string): File input selector
- `filePath` (string): File path
- `next` (node|null): Next node

## 🔀 **Control Flow Nodes**

### **conditional**
Branch based on condition.

```json
{
  "nodeType": "conditional",
  "conditionExpression": "{{user.age}} > 18",
  "branches": {
    "yes": {
      "nodeType": "fillField",
      "selector": "#adult-field",
      "value": "adult"
    },
    "no": {
      "nodeType": "fillField",
      "selector": "#minor-field", 
      "value": "minor"
    }
  }
}
```

**Properties:**
- `conditionExpression` (string): Boolean expression
- `branches` (object): yes/no branches

### **question**
Check data and branch.

```json
{
  "nodeType": "question",
  "check": {
    "dataPath": "user.isActive",
    "operator": "equals",
    "expectedValue": true
  },
  "branches": {
    "yes": {
      "nodeType": "fillField",
      "selector": "#status",
      "value": "active"
    },
    "no": {
      "nodeType": "fillField",
      "selector": "#status",
      "value": "inactive"
    }
  }
}
```

**Properties:**
- `check.dataPath` (string): Path to data
- `check.operator` (string): equals, greaterThan, contains
- `check.expectedValue` (any): Expected value
- `branches` (object): yes/no branches

### **sequence**
Execute nodes in order.

```json
{
  "nodeType": "sequence",
  "sequence": [
    {
      "nodeType": "fillField",
      "selector": "#field1",
      "value": "{{user.value1}}"
    },
    {
      "nodeType": "fillField", 
      "selector": "#field2",
      "value": "{{user.value2}}"
    },
    {
      "nodeType": "clickButton",
      "selector": "#submit"
    }
  ],
  "next": {
    "nodeType": "wait",
    "duration": 2000
  }
}
```

**Properties:**
- `sequence` (array): Nodes to execute
- `next` (node|null): Node after sequence

### **forEach**
Iterate over data array and execute node for each item.

```json
{
  "nodeType": "forEach",
  "dataSourceIteratorParam": "currentItem",
  "next": {
    "nodeType": "fillField",
    "selector": "#item-field",
    "value": "{{currentItem}}"
  }
}
```

**Properties:**
- `dataSourceIteratorParam` (string): Context key containing array to iterate over
- `next` (node): Node to execute for each item

**How it works:**
1. Retrieves array from context using `dataSourceIteratorParam` key
2. For each item in the array:
   - Sets the current item as the value of `dataSourceIteratorParam` in context
   - Executes the `next` node with the current item available
3. Cleans up by setting `dataSourceIteratorParam` to null
4. Continues to next node after loop completion

**Context Usage:**
- The array must be pre-loaded in the user context with the key specified in `dataSourceIteratorParam`
- During iteration, each array item becomes accessible via `{{dataSourceIteratorParam}}`

## ⏰ **Utility Nodes**

### **wait**
Pause execution.

```json
{
  "nodeType": "wait",
  "duration": 3000,
  "next": {
    "nodeType": "fillField",
    "selector": "#delayed-field",
    "value": "delayed"
  }
}
```

**Properties:**
- `duration` (number): Milliseconds
- `next` (node|null): Next node

## 📚 **Template System**

### **User Variables**
```json
"value": "{{user.firstName}}"
"value": "{{user.email}}"
```

### **Dynamic Data Variables**
```json
"value": "{{dataSourceIteratorParam}}"
"filePath": "{{currentItem}}"
```

## 🎯 **Examples**

### **Login Flow (Separate Actions)**
```json
{
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
```

### **Form with Sequence**
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
  ],
  "next": {
    "nodeType": "clickButton",
    "selector": "#submit"
  }
}
```

### **ForEach Data Iteration**
```json
{
  "nodeType": "forEach",
  "dataSourceIteratorParam": "documentPath",
  "next": {
    "nodeType": "sendFile",
    "selector": "input[type='file']",
    "filePath": "{{documentPath}}"
  }
}
```

**Note:** The `documentPath` array must be pre-loaded in the user context. During execution, each document path becomes available as `{{documentPath}}`.

### **Conditional Branch**
```json
{
  "nodeType": "conditional",
  "conditionExpression": "{{user.age}} > 18",
  "branches": {
    "yes": {
      "nodeType": "fillField",
      "selector": "#age-category",
      "value": "adult"
    },
    "no": {
      "nodeType": "fillField",
      "selector": "#age-category", 
      "value": "minor"
    }
  }
}
```

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

### **❌ Wrong (Double Action)**
```json
{
  "nodeType": "fillForm",
  "fields": [...],
  "submitButton": "#submit"  // <- FORBIDDEN!
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

---

> 💡 **Simplicity**: One node = one action. Parseable structure. ForEach with user questions. 