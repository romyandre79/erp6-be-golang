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

Transform Component - Number and Date Formatting
This document describes the formatting transform types available in comp_transform for formatting numbers and dates to Indonesian formats.

Context Injection Feature
All formatting transforms automatically inject formatted values into ctx.Extras with a _formatted suffix. This allows subsequent nodes to reference the formatted values using $fieldname_formatted.

Example:

Transform 1: format_date_indonesian with key_field=startdate
→ Injects: startdate_formatted = "Selasa, 6 Januari 2026"
Transform 2: Use in template
message_template: "Date: $startdate_formatted"
→ Output: "Date: Selasa, 6 Januari 2026"
Available Transform Types
1. format_rupiah
Formats numeric values as Indonesian Rupiah currency with thousand separators.

Parameters:

transform_type: format_rupiah
message_template: Template string with placeholders (e.g., "Total: {{ amount }}")
key_field (optional): Field name to inject
value_field (optional): Value to inject
Input Formats Supported:

Integer strings: "25000"
Float strings: "25000.0000", "1500000.50"
Both {{field}} and {{ field }} placeholder formats
Output Format:

"Rp 25.000" (Indonesian thousand separator using dots)
Context injection: fieldname_formatted = "Rp 25.000"
Example:

Input: { amount: "25000.0000" }
Template: "Total pembayaran: {{ amount }}"
Output: "Total pembayaran: Rp 25.000"
Context: amount_formatted = "Rp 25.000"
2. format_thousand
Formats numeric values with Indonesian thousand separators (without currency symbol).

Parameters:

transform_type: format_thousand
message_template: Template string with placeholders
key_field (optional): Field name to inject
value_field (optional): Value to inject
Input Formats Supported:

Integer strings: "1000000"
Float strings: "1000000.50"
Both {{field}} and {{ field }} placeholder formats
Output Format:

"1.000.000" (Indonesian thousand separator using dots)
Context injection: fieldname_formatted = "1.500.000"
Example:

Input: { quantity: "1500000" }
Template: "Jumlah: {{ quantity }} unit"
Output: "Jumlah: 1.500.000 unit"
Context: quantity_formatted = "1.500.000"
3. format_date_indonesian
Formats date strings to Indonesian date format with Indonesian day and month names.

Parameters:

transform_type: format_date_indonesian
message_template: Template string with placeholders
date_format: Output format - long, short, or numeric (default: long)
key_field (optional): Field name to inject
value_field (optional): Value to inject
Input Formats Supported:

2006-01-02T15:04:05Z07:00 - ISO8601 with timezone
2006-01-02T15:04:05Z - ISO8601 UTC
2006-01-02 15:04:05 - MySQL datetime
2006-01-02 - Date only (YYYY-MM-DD)
02/01/2006 - DD/MM/YYYY
01/02/2006 - MM/DD/YYYY
2006/01/02 - YYYY/MM/DD
Output Formats:

long: "Rabu, 8 Januari 2026" (with day name)
short: "8 Januari 2026" (without day name)
numeric: "08/01/2026" (dd/MM/yyyy)
Context injection: fieldname_formatted = "Rabu, 8 Januari 2026"
Indonesian Names:

Days: Minggu, Senin, Selasa, Rabu, Kamis, Jumat, Sabtu
Months: Januari, Februari, Maret, April, Mei, Juni, Juli, Agustus, September, Oktober, November, Desember
Example:

Input: { startdate: "2026-01-06" }
Template: "Tanggal mulai: {{ startdate }}"
Output: "Tanggal mulai: Senin, 6 Januari 2026"
Implementation Details
Files Modified:
core/generator/comp_transform.go

Added three new transform type cases
Added support for both {{key}} and {{ key }} placeholder formats
Uses strconv.ParseFloat() for number parsing (handles decimals)
Added debug logging for troubleshooting
core/helpers/util.go

Added 
FormatThousands(n int64) string
 - formats with thousand separators
Added 
FormatRupiah(n int64) string
 - formats as Rupiah currency
Added 
FormatDateIndonesian(dateStr string) string
 - formats dates to Indonesian
Error Handling:
Numbers: If a value cannot be parsed as a number, the original value is used
Dates: If a date cannot be parsed in any supported format, the original string is returned
Placeholder Support:
All three transform types support both placeholder formats:

{{fieldname}} - no spaces
{{ fieldname }} - with spaces
Usage Tips
Chaining Transforms
When chaining multiple transforms, be aware that each transform outputs a message field and the original data fields. The next transform in the chain will receive this combined data.

Example workflow:

Query node → outputs { amount: "25000", date: "2026-01-06" }
Transform (format_rupiah) → outputs { amount: "25000", message: "Rp 25.000" }
Transform (format_date_indonesian) → receives the output from step 2
Debug Logging
The format_date_indonesian transform includes debug logging that shows:

Available fields in the input data
The template being used
Each replacement being made
Check your console output for lines like:

[Transform] format_date_indonesian - Available fields: [startdate enddate]
[Transform] format_date_indonesian - Template: '{{ startdate }}'
[Transform] Replacing '{{startdate}}' with 'Senin, 6 Januari 2026' (from value '2026-01-06')
Common Issues
Issue: Placeholder not replaced (shows {{fieldname}} in output)
Cause: The field name in the template doesn't match any field in the input data.

Solution:

Check the debug logs to see available fields
Ensure the previous node outputs the field you're trying to format
Verify the field name spelling matches exactly
Issue: Number shows as original value instead of formatted
Cause: The value couldn't be parsed as a number (e.g., contains non-numeric characters).

Solution: Ensure the input value is a valid number string. The parser handles both integers and decimals.

Issue: Date shows as original value instead of formatted
Cause: The date format doesn't match any of the supported formats.

Solution: Convert your date to one of the supported formats before passing to the transform, or add the new format to the formats array in 
FormatDateIndonesian()
.