# Workflow Component Architecture

## Overview

The ERP6 workflow engine uses a **modular component architecture** that supports three types of components:

1. **Native Components** - Built into the core engine (Start, End)
2. **Built-in Plugin Components** - Compiled with the app but in separate files
3. **External Plugin Components** - Loaded dynamically at runtime

## Component Types

### 1. Native Components (Core Engine)
These components are fundamental to the workflow engine and remain in `executeflow.go`:

- **Start** - Initializes workflow execution
- **End** - Terminates workflow execution

### 2. Built-in Plugin Components
These are compiled into the binary but organized in separate files for maintainability:

| Component | File | Description |
|-----------|------|-------------|
| Search | `comp_search.go` | Database queries with paging |
| SearchRow | `comp_search.go` | Single row query |
| SearchSingle | `comp_search.go` | Single value query |
| API | `comp_api.go` | HTTP API calls (GET/POST) |
| Calculator | `comp_calculator.go` | Mathematical expression evaluation |
| Midtrans | `comp_midtrans.go` | Payment gateway integration |
| Brevo | `comp_brevo.go` | Email via SMTP |
| Decision | `comp_decision.go` | Conditional branching |
| Table | `comp_table.go` | Direct table operations (INSERT/UPDATE/DELETE) |
| Workflow | `comp_workflow.go` | Execute sub-workflows |
| Report | `comp_report.go` | Generate reports (PDF/Excel/CSV) |
| ImportData | `comp_import.go` | Excel file import |
| StoreProcedure | `comp_storeproc.go` | Execute stored procedures |
| SendEmail | `comp_email.go` | Generic email sending |
| SaveLog | `comp_savelog.go` | Logging operations |
| SendMessage | `comp_message.go` | User messages |

### 3. External Plugin Components
Loaded dynamically without recompiling:

#### A. Executable Plugins (`.exe`)
- Standalone binaries that communicate via JSON over stdin/stdout
- **No rebuild required** - just upload and use
- Can be written in any language (Go, Python, Node.js, etc.)
- Example: `ExternalRandom` (see `external_plugin.zip`)

#### B. Interpreted Plugins (`.go` via Yaegi)
- Go source code executed by the Yaegi interpreter
- **No rebuild required** - interpreted at runtime
- Must use `package main` and specific function signatures
- Example: `RandomNumber` (see `random_plugin.zip`)

## Component Registration

All components register themselves using the `RegisterComponent` function in their `init()` blocks:

```go
func init() {
    RegisterComponent("componentname", func(ctx *WorkflowContext) error {
        return handleComponentName(ctx.FiberCtx, ctx.Params, ctx.DB)
    })
}
```

## Plugin Upload API

**Endpoint**: `POST /api/plugins/upload`

**Request**: `multipart/form-data` with `plugin` field containing a ZIP file

**ZIP Structure**:
```
plugin.zip
├── plugin.json          # Manifest (required)
├── component.go         # Go source (optional, for Yaegi)
└── component.exe        # Binary (optional, for external)
```

**plugin.json Example**:
```json
{
  "componentname": "MyPlugin",
  "componenttitle": "My Custom Plugin",
  "componentcategoryid": 5,
  "componentclass": "myplugin",
  "input": 1,
  "output": 1,
  "details": [
    {
      "detailtype": "text",
      "lable": "Input Parameter",
      "inputtype": "text",
      "inputname": "param1",
      "inputdesc": "Description",
      "order": 1
    }
  ]
}
```

## Writing External Plugins

### Go Example (Compiled to `.exe`):

```go
package main

import (
    "encoding/json"
    "os"
)

type Input struct {
    Params []struct {
        InputName string `json:"inputname"`
        CompValue string `json:"compvalue"`
    } `json:"params"`
}

type Output struct {
    Result interface{} `json:"result"`
    Error  string      `json:"error"`
}

func main() {
    var input Input
    json.NewDecoder(os.Stdin).Decode(&input)
    
    // Your logic here
    result := "Hello from plugin!"
    
    json.NewEncoder(os.Stdout).Encode(Output{Result: result})
}
```

### Python Example:

```python
import json
import sys

input_data = json.load(sys.stdin)
params = input_data.get('params', [])

# Your logic here
result = "Hello from Python!"

json.dump({"result": result, "error": ""}, sys.stdout)
```

## Component Lifecycle

1. **Registration** - Component registers itself via `init()`
2. **Discovery** - Registry maintains a map of all components
3. **Execution** - `InternalFlow` calls `GetComponent(name)` and executes the handler
4. **Context Passing** - `WorkflowContext` provides access to Fiber context, DB, params, etc.
5. **Result Storage** - Results are appended to `wfEngine` locals

## Benefits of This Architecture

✅ **Modularity** - Each component is self-contained
✅ **Extensibility** - Add new components without modifying core engine
✅ **Hot Reload** - External plugins can be added without restart
✅ **Language Agnostic** - External plugins can be written in any language
✅ **Type Safety** - Built-in components are fully type-checked
✅ **Maintainability** - Easier to find and modify specific components
