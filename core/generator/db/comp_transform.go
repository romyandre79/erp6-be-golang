package generator

import (
	"fmt"
	"strings"
)

func init() {
	RegisterComponent("transform", func(ctx *WorkflowContext) error {
		return handleTransform(ctx)
	})
}

// handleTransform transforms data from previous node
// Supports: array_to_first, format_message, format_array
func handleTransform(ctx *WorkflowContext) error {
	// Get previous node result
	wfEngine, ok := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
	if !ok || len(wfEngine) == 0 {
		err := fmt.Errorf("no previous node result found")
		appendStepResult(ctx.FiberCtx, 0, 0, "Transform", nil, map[string]interface{}{
			"error": err.Error(),
		}, false, 0, "")
		return err
	}

	lastResult := wfEngine[len(wfEngine)-1]
	if lastResult.ResultNode == nil {
		err := fmt.Errorf("previous node has no result")
		appendStepResult(ctx.FiberCtx, 0, 0, "Transform", nil, map[string]interface{}{
			"error": err.Error(),
		}, false, 0, "")
		return err
	}

	// Get parameters first to check transform type
	var transformType, messageTemplate, keyField, valueField, separator string
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
		// Format two parallel arrays into a message (like currency = rate)
		// Defaults
		if separator == "" {
			separator = " = "
		}
		lineFormat := messageTemplate
		if lineFormat == "" {
			lineFormat = "{key}{sep}{value}"
		}

		// Get the arrays using helper that handles []string and []interface{}
		keyArray, keyOk := getAsArray(resultMap[keyField])
		valueArray, valueOk := getAsArray(resultMap[valueField])

		fmt.Printf("[Transform] Array check: key='%s' exists=%v len=%d, value='%s' exists=%v len=%d\n", 
			keyField, keyOk, len(keyArray), valueField, valueOk, len(valueArray))

		if !keyOk || !valueOk {
			err := fmt.Errorf("key_field '%s' or value_field '%s' not found or not arrays. Available fields: %v", keyField, valueField, getMapKeys(resultMap))
			fmt.Printf("[Transform] Error: %v\n", err)
			appendStepResult(ctx.FiberCtx, 0, 0, "Transform", nil, map[string]interface{}{
				"error": err.Error(),
			}, false, 0, "")
			return err
		}

		// Format the message
		var lines []string
		maxLen := len(keyArray)
		if len(valueArray) < maxLen {
			maxLen = len(valueArray)
		}

		for i := 0; i < maxLen; i++ {
			key := fmt.Sprintf("%v", keyArray[i])
			value := fmt.Sprintf("%v", valueArray[i])
			
			// Replace placeholders in format
			line := strings.ReplaceAll(lineFormat, "{key}", key)
			line = strings.ReplaceAll(line, "{value}", value)
			line = strings.ReplaceAll(line, "{sep}", separator)
			
			lines = append(lines, line)
		}

		message := strings.Join(lines, "\n")

		result = map[string]interface{}{
			"message": message,
			"count":   maxLen,
		}

	default:
		return fmt.Errorf("unknown transform_type: %s", transformType)
	}

	fmt.Printf("[Transform] Output result: %+v\n", result)
	appendStepResult(ctx.FiberCtx, 0, 0, "Transform", nil, result, true, 0, "")

	// Append result to workflow engine
	wm := WorkflowEngine{
		ResultNode: result,
	}
	wfEngine, _ = ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, wm)
	ctx.FiberCtx.Locals("wfEngine", wfEngine)

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
