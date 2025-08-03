# Traverser Examples - Single Action

> Single action per node examples with parseable structure

## 🎯 **Basic Examples**

### **Example 1: Simple Login (Separate Actions)**
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

**Context:**
```json
{
  "user": {
    "email": "john@example.com",
    "password": "secret123"
  }
}
```

### **Example 2: File Upload**
```json
{
  "graph": {
    "nodeType": "moveToPage",
    "url": "https://example.com/upload",
    "next": {
      "nodeType": "sendFile",
      "selector": "input[type='file']",
      "filePath": "{{user.filePath}}",
      "next": {
        "nodeType": "clickButton",
        "selector": "#upload-btn",
        "next": null
      }
    }
  },
  "metadata": {
    "name": "File Upload",
    "version": "1.0.0"
  }
}
```

**Context:**
```json
{
  "user": {
    "filePath": "/Users/john/document.pdf"
  }
}
```

## 🔀 **Conditional Examples**

### **Example 3: Age Check**
```json
{
  "graph": {
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
  },
  "metadata": {
    "name": "Age Check",
    "version": "1.0.0"
  }
}
```

### **Example 4: Data Check**
```json
{
  "graph": {
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
        "value": "Active user: {{user.name}}",
        "next": null
      },
      "no": {
        "nodeType": "fillField",
        "selector": "#status",
        "value": "Inactive user: {{user.name}}",
        "next": null
      }
    }
  }
}
```

## 🔄 **ForEach Examples**

### **Example 5: ForEach with User Question**
```json
{
  "graph": {
    "nodeType": "forEach",
    "dataSource": "{{user.items}}",
    "questionText": "Process item {{iterator.count}} of {{iterator.total}}?",
    "next": {
      "nodeType": "fillField",
      "selector": "#item-name",
      "value": "Processing item {{iterator.index}}",
      "next": {
        "nodeType": "wait",
        "duration": 1000,
        "next": null
      }
    }
  },
  "metadata": {
    "name": "Process Items",
    "version": "1.0.0"
  }
}
```

**Context:**
```json
{
  "user": {
    "items": ["item1", "item2", "item3", "item4"]
  }
}
```

### **Example 6: Upload Multiple Files**
```json
{
  "graph": {
    "nodeType": "moveToPage",
    "url": "https://example.com/upload",
    "next": {
      "nodeType": "forEach",
      "dataSource": "{{user.documents}}",
      "questionText": "Upload {{user.documents[iterator.index]}}?",
      "next": {
        "nodeType": "sendFile",
        "selector": "input[type='file']",
        "filePath": "{{user.documents[iterator.index]}}",
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

**Context:**
```json
{
  "user": {
    "documents": [
      "/Users/john/doc1.pdf",
      "/Users/john/doc2.pdf", 
      "/Users/john/doc3.pdf"
    ]
  }
}
```

## 📋 **Sequence Examples**

### **Example 7: Registration Form (Sequence)**
```json
{
  "graph": {
    "nodeType": "moveToPage",
    "url": "https://example.com/register",
    "next": {
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
        },
        {
          "nodeType": "fillField",
          "selector": "#phone",
          "value": "{{user.phone}}"
        }
      ],
      "next": {
        "nodeType": "clickButton",
        "selector": "#register-btn",
        "next": null
      }
    }
  },
  "metadata": {
    "name": "Registration Form",
    "version": "1.0.0"
  }
}
```

### **Example 8: Search Process**
```json
{
  "graph": {
    "nodeType": "sequence",
    "sequence": [
      {
        "nodeType": "moveToPage",
        "url": "{{user.website}}"
      },
      {
        "nodeType": "fillField",
        "selector": "#search",
        "value": "{{user.searchTerm}}"
      },
      {
        "nodeType": "clickButton",
        "selector": "#search-btn"
      },
      {
        "nodeType": "wait",
        "duration": 2000
      }
    ],
    "next": null
  }
}
```

## 🎯 **Complex Examples**

### **Example 9: Social Media Post with Conditional**
```json
{
  "graph": {
    "nodeType": "moveToPage",
    "url": "https://facebook.com",
    "next": {
      "nodeType": "fillField",
      "selector": "#email",
      "value": "{{user.email}}",
      "next": {
        "nodeType": "fillField",
        "selector": "#pass",
        "value": "{{user.password}}",
        "next": {
          "nodeType": "clickButton",
          "selector": "[name='login']",
          "next": {
            "nodeType": "wait",
            "duration": 3000,
            "next": {
              "nodeType": "conditional",
              "conditionExpression": "{{user.hasMessage}} == true",
              "branches": {
                "yes": {
                  "nodeType": "fillField",
                  "selector": "[contenteditable='true']",
                  "value": "{{user.message}}",
                  "next": {
                    "nodeType": "clickButton",
                    "selector": "[aria-label='Post']",
                    "next": null
                  }
                },
                "no": null
              }
            }
          }
        }
      }
    }
  }
}
```

**Context:**
```json
{
  "user": {
    "email": "user@example.com",
    "password": "secret123",
    "hasMessage": true,
    "message": "Hello, world! 👋"
  }
}
```

### **Example 10: E-commerce Product Search**
```json
{
  "graph": {
    "nodeType": "moveToPage",
    "url": "https://shop.example.com",
    "next": {
      "nodeType": "forEach",
      "dataSource": "{{user.products}}",
      "questionText": "Search for {{user.products[iterator.index]}}?",
      "next": {
        "nodeType": "fillField",
        "selector": "#search-input",
        "value": "{{user.products[iterator.index]}}",
        "next": {
          "nodeType": "clickButton",
          "selector": "#search-btn",
          "next": {
            "nodeType": "wait",
            "duration": 3000,
            "next": {
              "nodeType": "clickButton",
              "selector": ".product-item:first-child .add-to-cart",
              "next": null
            }
          }
        }
      }
    }
  }
}
```

## 📋 **Usage Patterns**

### **Running Examples**
```bash
# Simple login
rpa-engine --workflow=login.json --context=user.json

# With inline context
rpa-engine --workflow=upload.json \
  --context='{"user":{"filePath":"/path/to/file.pdf"}}'

# ForEach example
rpa-engine --workflow=foreach.json --context=items.json
```

### **Context Files**

**user.json:**
```json
{
  "user": {
    "email": "john@example.com",
    "password": "secret123",
    "firstName": "John",
    "lastName": "Doe",
    "phone": "+1234567890",
    "website": "https://google.com",
    "searchTerm": "golang tutorial",
    "age": 30,
    "isActive": true
  }
}
```

**items.json:**
```json
{
  "user": {
    "items": ["task1", "task2", "task3"],
    "documents": [
      "/Users/john/file1.pdf",
      "/Users/john/file2.pdf"
    ],
    "products": ["laptop", "mouse", "keyboard"]
  }
}
```

## ✅ **Single Action Rules**

### **✅ Correct Examples**

**Separate Actions:**
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

**Sequence for Multiple Fields:**
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

### **❌ Wrong Examples**

**Multiple Actions in One Node:**
```json
{
  "nodeType": "fillForm",
  "fields": [...],
  "submitButton": "#submit"  // FORBIDDEN!
}
```

**Mixed Actions:**
```json
{
  "nodeType": "fillAndClick",
  "selector": "#field",
  "value": "text",
  "clickSelector": "#button"  // FORBIDDEN!
}
```

## 🎯 **Execution Flow**

### **Linear Chain**
```
moveToPage → fillField → fillField → clickButton → null
```

### **Sequence Block**
```
sequence [
  fillField,
  fillField, 
  clickButton
] → next → null
```

### **ForEach Block**
```
forEach (ask user for each item) → action → next item...
```

### **Conditional Branch**
```
conditional → yes: fillField → null
           → no: fillField → null
```

---

> 💡 **Simplicity**: One node = one action. ForEach with questions. Parseable structure without complexity. 