# Traverser JSON Schema - Single Action

> JSON schema for single-action parseable workflows

## 📋 **Root Schema**

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "graph": {
      "type": "object",
      "$ref": "#/definitions/node"
    },
    "metadata": {
      "type": "object",
      "properties": {
        "name": {"type": "string"},
        "version": {"type": "string"},
        "description": {"type": "string"}
      }
    }
  },
  "required": ["graph"]
}
```

## 🌳 **Node Schema**

```json
{
  "definitions": {
    "node": {
      "type": "object",
      "properties": {
        "nodeType": {
          "type": "string",
          "enum": [
            "moveToPage",
            "fillField",
            "clickButton", 
            "sendFile",
            "conditional",
            "question",
            "sequence",
            "forEach",
            "wait"
          ]
        },
        "id": {"type": "string"},
        "next": {
          "oneOf": [
            {"$ref": "#/definitions/node"},
            {"type": "null"}
          ]
        }
      },
      "required": ["nodeType"]
    }
  }
}
```

## 🎯 **Action Nodes**

### **moveToPage**
```json
{
  "allOf": [{"$ref": "#/definitions/node"}],
  "properties": {
    "nodeType": {"const": "moveToPage"},
    "url": {"type": "string"}
  },
  "required": ["nodeType", "url"]
}
```

### **fillField**
```json
{
  "allOf": [{"$ref": "#/definitions/node"}],
  "properties": {
    "nodeType": {"const": "fillField"},
    "selector": {"type": "string"},
    "value": {"type": "string"}
  },
  "required": ["nodeType", "selector", "value"]
}
```

### **clickButton**
```json
{
  "allOf": [{"$ref": "#/definitions/node"}],
  "properties": {
    "nodeType": {"const": "clickButton"},
    "selector": {"type": "string"}
  },
  "required": ["nodeType", "selector"]
}
```

### **sendFile**
```json
{
  "allOf": [{"$ref": "#/definitions/node"}],
  "properties": {
    "nodeType": {"const": "sendFile"},
    "selector": {"type": "string"},
    "filePath": {"type": "string"}
  },
  "required": ["nodeType", "selector", "filePath"]
}
```

### **wait**
```json
{
  "allOf": [{"$ref": "#/definitions/node"}],
  "properties": {
    "nodeType": {"const": "wait"},
    "duration": {"type": "integer", "minimum": 0}
  },
  "required": ["nodeType", "duration"]
}
```

## 🔀 **Control Nodes**

### **conditional**
```json
{
  "allOf": [{"$ref": "#/definitions/node"}],
  "properties": {
    "nodeType": {"const": "conditional"},
    "conditionExpression": {"type": "string"},
    "branches": {
      "type": "object",
      "properties": {
        "yes": {"$ref": "#/definitions/node"},
        "no": {"$ref": "#/definitions/node"}
      }
    }
  },
  "required": ["nodeType", "conditionExpression", "branches"]
}
```

### **question**
```json
{
  "allOf": [{"$ref": "#/definitions/node"}],
  "properties": {
    "nodeType": {"const": "question"},
    "check": {
      "type": "object",
      "properties": {
        "dataPath": {"type": "string"},
        "operator": {
          "type": "string",
          "enum": ["equals", "greaterThan", "contains"]
        },
        "expectedValue": {}
      },
      "required": ["dataPath", "operator", "expectedValue"]
    },
    "branches": {
      "type": "object",
      "properties": {
        "yes": {"$ref": "#/definitions/node"},
        "no": {"$ref": "#/definitions/node"}
      }
    }
  },
  "required": ["nodeType", "check", "branches"]
}
```

### **sequence**
```json
{
  "allOf": [{"$ref": "#/definitions/node"}],
  "properties": {
    "nodeType": {"const": "sequence"},
    "sequence": {
      "type": "array",
      "items": {"$ref": "#/definitions/node"},
      "minItems": 1
    }
  },
  "required": ["nodeType", "sequence"]
}
```

### **forEach**
```json
{
  "allOf": [{"$ref": "#/definitions/node"}],
  "properties": {
    "nodeType": {"const": "forEach"},
    "dataSource": {"type": "string"},
    "questionText": {"type": "string"}
  },
  "required": ["nodeType", "dataSource", "questionText"]
}
```

## 🔄 **Complete Example**

```json
{
  "graph": {
    "nodeType": "moveToPage",
    "url": "{{user.website}}",
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
          "next": {
            "nodeType": "conditional",
            "conditionExpression": "{{user.age}} > 18",
            "branches": {
              "yes": {
                "nodeType": "fillField",
                "selector": "#age-category",
                "value": "adult",
                "next": null
              },
              "no": {
                "nodeType": "fillField",
                "selector": "#age-category",
                "value": "minor",
                "next": null
              }
            }
          }
        }
      }
    }
  },
  "metadata": {
    "name": "Single Action Login",
    "version": "1.0.0"
  }
}
```

## 📚 **Template Variables**

### **Supported Templates**
- `{{user.propertyName}}` - User data
- `{{user.array[iterator.index]}}` - Array item in forEach
- `{{iterator.index}}` - Current index (0-based)
- `{{iterator.count}}` - Current count (1-based)
- `{{iterator.total}}` - Total items

### **Template Examples**
```json
{
  "url": "{{user.website}}",
  "value": "{{user.firstName}} {{user.lastName}}",
  "questionText": "Process {{user.items[iterator.index]}} ({{iterator.count}}/{{iterator.total}})?"
}
```

## ✅ **Single Action Rules**

### **✅ Valid Schema**
Each node performs exactly one action:

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

### **❌ Invalid Schema**
Multiple actions in one node:

```json
{
  "nodeType": "fillForm",
  "fields": [...],
  "submitButton": "#submit"  // FORBIDDEN
}
```

## 🎯 **ForEach Schema**

### **ForEach with Question**
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

### **Sequence for Multiple Actions**
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
  ]
}
```

---

> 💡 **Simplicity**: One node = one action. ForEach with questions. Parseable JSON schema. 