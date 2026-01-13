package generator

import (
	"fmt"
	"strconv"
	"strings"
	helpers "erp6-be-golang/core/helpers"
)

func init() {
	RegisterComponent("transform", func(ctx *WorkflowContext) error {
		return handleTransform(ctx)
	})
}

// handleTransform transforms data from previous node
// Supports: array_to_first, format_message, format_array
func handleTransform(ctx *WorkflowContext) error {
	// Get previous node result from wfEngine
	var wfEngine []WorkflowEngine
	var ok bool
	
	if ctx.FiberCtx != nil {
		wfEngine, ok = ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
	} else {
		// WhatsApp context - use Extras
		wfEngine, ok = ctx.Extras["wfEngine"].([]WorkflowEngine)
	}
	
	if !ok || len(wfEngine) == 0 {
		return fmt.Errorf("no previous node result found")
	}

	lastResult := wfEngine[len(wfEngine)-1]
	if lastResult.ResultNode == nil {
		return fmt.Errorf("previous node has no result")
	}

	// Get parameters first to check transform type
	var transformType, messageTemplate, keyField, valueField, separator, dateFormat string
	for _, p := range ctx.Params {
		val := strings.TrimSpace(ResolveParam(ctx.FiberCtx, p.CompValue))
		switch strings.ToLower(p.InputName) {
		case "transform_type":
			transformType = val
		case "message_template":
			messageTemplate = val
		case "key_field":
			keyField = val
		case "value_field":
			valueField = val
		case "separator":
			separator = val
		case "date_format":
			dateFormat = val
		}
	}

	var resultMap map[string]interface{}
	
	// Debug logging
	fmt.Printf("[Transform] Transform type: '%s', Input type: %T\n", transformType, lastResult.ResultNode)
	fmt.Printf("[Transform] Params: key_field='%s', value_field='%s', separator='%s'\n", keyField, valueField, separator)
	
	// Special handling: if input is an array and transform_type is array_to_first, convert it
	if arr, isArray := lastResult.ResultNode.([]interface{}); isArray && transformType == "array_to_first" {
		fmt.Printf("[Transform] Handling array input with %d elements\n", len(arr))
		if len(arr) > 0 {
			// Wrap array in a map so array_to_first can process it
			resultMap = map[string]interface{}{
				"result": arr,
			}
		} else {
			resultMap = map[string]interface{}{
				"result": []interface{}{},
			}
		}
	} else {
		// Normal case: expect a map
		var ok bool
		resultMap, ok = lastResult.ResultNode.(map[string]interface{})
		if !ok {
			err := fmt.Errorf("previous node result is not a map, got type: %T", lastResult.ResultNode)
			appendStepResult(ctx.FiberCtx, 0, 0, "Transform", nil, map[string]interface{}{
				"error": err.Error(),
			}, false, 0, "")
			return err
		}

		// Check for nested "result" map (from Scraper wrapper)
		if nestedResult, ok := resultMap["result"].(map[string]interface{}); ok {
			fmt.Println("[Transform] Detected nested 'result' map, unpacking...")
			// Merge nested result into top level or just use it?
			// Let's use it as the primary map, but maybe keep original keys if needed?
			// For transform purposes (accessing fields), using the nested map is usually what we want.
			resultMap = nestedResult
		}
	}

	var result map[string]interface{}

	switch transformType {
	case "array_to_first":
		// Convert all arrays to their first element
		result = make(map[string]interface{})
		for key, value := range resultMap {
			if arr, ok := getAsArray(value); ok && len(arr) > 0 {
				// If the first element is a map, flatten it to top level
				if firstMap, isMap := arr[0].(map[string]interface{}); isMap {
					// Flatten: copy all fields from the first element to result
					for k, v := range firstMap {
						result[k] = v
					}
				} else {
					result[key] = arr[0]
				}
			} else {
				result[key] = value
			}
		}

	case "format_message":
		// Format data into a message using template
		result = make(map[string]interface{})
		
		// First, convert arrays to first element
		flatData := make(map[string]string)
		for key, value := range resultMap {
			if arr, ok := value.([]interface{}); ok && len(arr) > 0 {
				flatData[key] = fmt.Sprintf("%v", arr[0])
			} else {
				flatData[key] = fmt.Sprintf("%v", value)
			}
		}

		// Replace placeholders in template
		message := messageTemplate
		for key, value := range flatData {
			placeholder := fmt.Sprintf("{{%s}}", key)
			message = strings.ReplaceAll(message, placeholder, value)
		}

		result["message"] = message
		// Also include the flat data
		for k, v := range flatData {
			result[k] = v
		}

	case "format_array":
		fmt.Println("[Transform] Processing format_array...")
		// Format arrays into a message
		if separator == "" {
			separator = " "
		}
		
        // Identify all array fields and find max length
        maxLen := 0
        arrayFields := make(map[string][]interface{})
        
        // If key_field/value_field provided, check them specifically (legacy support)
        if keyField != "" && valueField != "" {
             if kArr, ok := getAsArray(resultMap[keyField]); ok {
                 arrayFields["key"] = kArr
                 if len(kArr) > maxLen { maxLen = len(kArr) }
             }
             if vArr, ok := getAsArray(resultMap[valueField]); ok {
                 arrayFields["value"] = vArr
                 if len(vArr) > maxLen { maxLen = len(vArr) }
             }
        }

        // Also map all other fields for generic access
        for k, v := range resultMap {
            if arr, ok := getAsArray(v); ok {
                arrayFields[k] = arr
                if len(arr) > maxLen { maxLen = len(arr) }
            }
        }

        fmt.Printf("[Transform] Found generic arrays. Max len: %d\n", maxLen)

		// Format the message
		var lines []string
		for i := 0; i < maxLen; i++ {
			line := messageTemplate
            
            // Replace all known array fields
            for field, arr := range arrayFields {
                val := ""
                if i < len(arr) {
                    val = fmt.Sprintf("%v", arr[i])
                }
                
                // Replace {field} and {{field}}
                line = strings.ReplaceAll(line, fmt.Sprintf("{%s}", field), val)
                line = strings.ReplaceAll(line, fmt.Sprintf("{{%s}}", field), val)
            }
            
            // Replace {sep}
			line = strings.ReplaceAll(line, "{sep}", separator)
			lines = append(lines, line)
		}

		message := strings.Join(lines, "\n")

		result = map[string]interface{}{
			"message": message,
			"count":   maxLen,
		}

	case "format_rupiah":
		// Format data into a message using template
		result = make(map[string]interface{})
		
		// First, convert arrays to first element
		flatData := make(map[string]string)
		for key, value := range resultMap {
			if arr, ok := value.([]interface{}); ok && len(arr) > 0 {
				flatData[key] = fmt.Sprintf("%v", arr[0])
			} else {
				flatData[key] = fmt.Sprintf("%v", value)
			}
		}
		
		// If key_field and value_field are provided, inject them into the data
		if keyField != "" && valueField != "" {
			flatData[keyField] = valueField
			fmt.Printf("[Transform] format_rupiah - Injected field: %s = %s\n", keyField, valueField)
		}
		
		// Get keys for debug logging
		availableKeys := make([]string, 0, len(flatData))
		for k := range flatData {
			availableKeys = append(availableKeys, k)
		}
		fmt.Printf("[Transform] format_rupiah - Available fields: %v\n", availableKeys)
		fmt.Printf("[Transform] format_rupiah - Template: '%s'\n", messageTemplate)

		// Replace placeholders in template and inject formatted values into context
		message := messageTemplate
		for key, value := range flatData {
			// Support both {{key}} and {{ key }} formats
			placeholder1 := fmt.Sprintf("{{%s}}", key)
			placeholder2 := fmt.Sprintf("{{ %s }}", key)
			
			// Try to parse as float first (handles both int and float), then convert to int64
			if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
				formatted := helpers.FormatRupiah(int64(floatVal))
				message = strings.ReplaceAll(message, placeholder1, formatted)
				message = strings.ReplaceAll(message, placeholder2, formatted)
				
				// Inject formatted value into context for use in subsequent nodes
				ctx.Extras[key+"_formatted"] = formatted
				fmt.Printf("[Transform] Injected to context: %s_formatted = %s\n", key, formatted)
			} else {
				// If not a number, just use the original value
				message = strings.ReplaceAll(message, placeholder1, value)
				message = strings.ReplaceAll(message, placeholder2, value)
			}
		}

		fmt.Printf("[Transform] format_rupiah - Final message: '%s'\n", message)
		result["message"] = message
		// Also include the flat data (but skip 'message' to avoid overwriting)
		for k, v := range flatData {
			if k != "message" {
				result[k] = v
			}
		}

	case "format_thousand":
		// Format data into a message using template
		result = make(map[string]interface{})
		
		// First, convert arrays to first element
		flatData := make(map[string]string)
		for key, value := range resultMap {
			if arr, ok := value.([]interface{}); ok && len(arr) > 0 {
				flatData[key] = fmt.Sprintf("%v", arr[0])
			} else {
				flatData[key] = fmt.Sprintf("%v", value)
			}
		}
		
		// If key_field and value_field are provided, inject them into the data
		if keyField != "" && valueField != "" {
			flatData[keyField] = valueField
			fmt.Printf("[Transform] format_thousand - Injected field: %s = %s\n", keyField, valueField)
		}
		
		// Get keys for debug logging
		availableKeys := make([]string, 0, len(flatData))
		for k := range flatData {
			availableKeys = append(availableKeys, k)
		}
		fmt.Printf("[Transform] format_thousand - Available fields: %v\n", availableKeys)
		fmt.Printf("[Transform] format_thousand - Template: '%s'\n", messageTemplate)

		// Replace placeholders in template and inject formatted values into context
		message := messageTemplate
		for key, value := range flatData {
			// Support both {{key}} and {{ key }} formats
			placeholder1 := fmt.Sprintf("{{%s}}", key)
			placeholder2 := fmt.Sprintf("{{ %s }}", key)
			
			// Try to parse as float first (handles both int and float), then convert to int64
			if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
				formatted := helpers.FormatThousands(int64(floatVal))
				message = strings.ReplaceAll(message, placeholder1, formatted)
				message = strings.ReplaceAll(message, placeholder2, formatted)
				
				// Inject formatted value into context for use in subsequent nodes
				ctx.Extras[key+"_formatted"] = formatted
				fmt.Printf("[Transform] Injected to context: %s_formatted = %s\n", key, formatted)
			} else {
				// If not a number, just use the original value
				message = strings.ReplaceAll(message, placeholder1, value)
				message = strings.ReplaceAll(message, placeholder2, value)
			}
		}

		result["message"] = message
		// Also include the flat data (but skip 'message' to avoid overwriting)
		for k, v := range flatData {
			if k != "message" {
				result[k] = v
			}
		}

	case "format_date_indonesian":
		// Format date to Indonesian format
		result = make(map[string]interface{})
		
		// First, convert arrays to first element
		flatData := make(map[string]string)
		for key, value := range resultMap {
			if arr, ok := value.([]interface{}); ok && len(arr) > 0 {
				flatData[key] = fmt.Sprintf("%v", arr[0])
			} else {
				flatData[key] = fmt.Sprintf("%v", value)
			}
		}
		
		// If key_field and value_field are provided, inject them into the data
		if keyField != "" && valueField != "" {
			flatData[keyField] = valueField
			fmt.Printf("[Transform] format_date_indonesian - Injected field: %s = %s\n", keyField, valueField)
		}

		// Get keys for debug logging
		availableKeys := make([]string, 0, len(flatData))
		for k := range flatData {
			availableKeys = append(availableKeys, k)
		}
		fmt.Printf("[Transform] format_date_indonesian - Available fields: %v\n", availableKeys)
		fmt.Printf("[Transform] format_date_indonesian - Template: '%s'\n", messageTemplate)

		// Default to "long" format if not specified
		if dateFormat == "" {
			dateFormat = "long"
		}
		
		// Replace placeholders in template and inject formatted values into context
		message := messageTemplate
		for key, value := range flatData {
			// Support both {{key}} and {{ key }} formats
			placeholder1 := fmt.Sprintf("{{%s}}", key)
			placeholder2 := fmt.Sprintf("{{ %s }}", key)
			
			// Try to parse as date with specified format
			formatted := helpers.FormatDateIndonesianWithFormat(value, dateFormat)
			fmt.Printf("[Transform] Replacing '%s' with '%s' (from value '%s', format '%s')\n", placeholder1, formatted, value, dateFormat)
			message = strings.ReplaceAll(message, placeholder1, formatted)
			message = strings.ReplaceAll(message, placeholder2, formatted)
			
			// Inject formatted value into context for use in subsequent nodes
			ctx.Extras[key+"_formatted"] = formatted
			fmt.Printf("[Transform] Injected to context: %s_formatted = %s\n", key, formatted)
		}

		result["message"] = message
		// Also include the flat data
		for k, v := range flatData {
			result[k] = v
		}

	case "split_string":
		// Split a string field into an array
		result = make(map[string]interface{})
		
		targetString := ""
		found := false
		
		if keyField != "" {
			if val, ok := resultMap[keyField]; ok {
				targetString = fmt.Sprintf("%v", val)
				found = true
			}
		} 
		
		if !found {
			// If no key specified or found, try to find "message" or "result" or just take the first value
			if val, ok := resultMap["message"]; ok {
				targetString = fmt.Sprintf("%v", val)
			} else if val, ok := resultMap["result"]; ok {
				targetString = fmt.Sprintf("%v", val)
			} else {
				// Fallback: take first value
				for _, v := range resultMap {
					targetString = fmt.Sprintf("%v", v)
					break
				}
			}
		}

		if separator == "" {
			separator = ","
		}

		parts := strings.Split(targetString, separator)
		finalParts := make([]interface{}, 0, len(parts))
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				finalParts = append(finalParts, trimmed)
			}
		}

		// Return as "result" key for comp_for to detect as columnar
		result["result"] = finalParts
		
		fmt.Printf("[Transform] split_string - Input: '%s', Sep: '%s', Count: %d\n", targetString, separator, len(finalParts))

	default:
		return fmt.Errorf("unknown transform_type: %s", transformType)
	}

	fmt.Printf("[Transform] Output result: %+v\n", result)
	
	// Append result to workflow engine
	wm := WorkflowEngine{
		ComponentName: "Transform",
		ResultNode:    result,
	}
	
	if ctx.FiberCtx != nil {
		appendStepResult(ctx.FiberCtx, 0, 0, "Transform", nil, result, true, 0, "")
		wfEngine, _ = ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
		wfEngine = append(wfEngine, wm)
		ctx.FiberCtx.Locals("wfEngine", wfEngine)
	} else {
		// WhatsApp context - use Extras
		wfEngine, _ = ctx.Extras["wfEngine"].([]WorkflowEngine)
		wfEngine = append(wfEngine, wm)
		ctx.Extras["wfEngine"] = wfEngine
	}

	return nil
}

func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// getAsArray safely converts generic interface to []interface{}
// Handles []interface{}, []string, and single values (wrapping them)
func getAsArray(val interface{}) ([]interface{}, bool) {
	if val == nil {
		return nil, false
	}

	switch v := val.(type) {
	case []interface{}:
		return v, true
	case []string:
		// Convert []string to []interface{}
		res := make([]interface{}, len(v))
		for i, s := range v {
			res[i] = s
		}
		return res, true
	case string:
		// Treat single string as array of 1
		return []interface{}{v}, true
	default:
		return nil, false
	}
}
