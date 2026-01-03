package generator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func init() {
	RegisterComponent("For Each", func(ctx *WorkflowContext) error {
		return handleForEach(ctx)
	})
}

// handleForEach executes a loop over a dataset (array or map of arrays)
// and triggers the connected downstream nodes for each item.
func handleForEach(ctx *WorkflowContext) error {
	var (
		sourceVar string
	)

	// Extract parameters
	for _, p := range ctx.Params {
		val := strings.TrimSpace(ResolveParam(ctx.FiberCtx, p.CompValue))
		switch p.InputName {
		case "source":
			sourceVar = val
		}
	}

	// 1. Resolve Source Data
	var sourceData interface{}
	
	if sourceVar == "" {
		// Try to get from previous node result
		wfEngine, ok := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
		if ok && len(wfEngine) > 0 {
			lastResult := wfEngine[len(wfEngine)-1]
			// Check if lastResult.ResultNode has data
			if lastResult.ResultNode != nil {
				sourceData = lastResult.ResultNode
				fmt.Printf("[ForEach] Using previous node result as source: %T\n", sourceData)
			}
		}
	} else {
		// Check if sourceVar is a variable reference (starts with $)
		if strings.HasPrefix(sourceVar, "$") {
			key := sourceVar[1:]
			if val, exists := ctx.Extras[key]; exists {
				sourceData = val
			} else {
				// Also check Locals directly?
				// or maybe it's just a raw JSON string?
				// For now error if not found
				fmt.Printf("[ForEach] Warning: variable %s not found in extras\n", key)
			}
		} else {
			// It's a direct string value (already resolved or raw input)
			sourceData = sourceVar
			fmt.Printf("[ForEach] Source resolved to value len: %d\n", len(sourceVar))
		}
	}

	if sourceData == nil {
		return fmt.Errorf("source data is nil or empty. Please specify a 'source' variable or ensure previous node returned data")
	}

	fmt.Printf("[ForEach] Looping over data type: %T\n", sourceData)

	// 2. Normalize Data into a list of Maps (rows)
	var rows []map[string]interface{}
	
	switch v := sourceData.(type) {
	case string:
		// Attempt to parse as JSON first
		var jsonObj interface{}
		if strings.TrimSpace(v) != "" {
			// Try to unmarshal into a structure we can loop over
			if err := json.Unmarshal([]byte(v), &jsonObj); err == nil {
				// Recursive call? Or just handle the unpacked data
				// Handling unpacked data is safer to avoid infinite recursion risk if types assume weird things
				switch jVal := jsonObj.(type) {
				case []interface{}:
					for _, item := range jVal {
						if mapItem, ok := item.(map[string]interface{}); ok {
							rows = append(rows, mapItem)
						} else {
							rows = append(rows, map[string]interface{}{"value": item})
						}
					}
				case map[string]interface{}:
					// Treat as single item or iterate keys?
					// Consistency with the map case above:
					// "If values are slices, we assume columnar."
					isColumnar := false
					maxLen := 0
					keys := make([]string, 0)
					
					for k, val := range jVal {
						if valSlice, ok := getAsSlice(val); ok {
							isColumnar = true
							if len(valSlice) > maxLen {
								maxLen = len(valSlice)
							}
							keys = append(keys, k)
						}
					}

					if isColumnar {
						for i := 0; i < maxLen; i++ {
							row := make(map[string]interface{})
							for _, k := range keys {
								valSlice, _ := getAsSlice(jVal[k])
								if i < len(valSlice) {
									row[k] = valSlice[i]
								} else {
									row[k] = nil
								}
							}
							rows = append(rows, row)
						}
					} else {
						rows = append(rows, jVal)
					}
				}
			} else {
				// Not JSON, treat as newline separated string
				// Formats: "KEY=VALUE" or just "VALUE"
				lines := strings.Split(v, "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line == "" {
						continue
					}
					
					// Check for Key=Value (simple assignment)
					// Only if it doesn't look like a URL or other stuff? 
					// Simple rule: first '=' splits it.
					if strings.Contains(line, "=") {
						parts := strings.SplitN(line, "=", 2)
						key := strings.TrimSpace(parts[0])
						val := strings.TrimSpace(parts[1])
						// If user wants specific key/value structured data
						rows = append(rows, map[string]interface{}{
							"key":   key,
							"value": val,
							// Also allow direct access via key? 
							// Logic below flattens result into Extras. 
							// So if we have "USD": "1.0", next node can access $USD.
							key: val, 
						})
					} else {
						// Simple list item
						rows = append(rows, map[string]interface{}{"value": line})
					}
				}
			}
		}

	case map[string]interface{}:
		// Check if it's a "Columnar" map (Map of Arrays) - simpler heuristic
		// or just a single object?
		// If values are slices, we assume columnar.
		isColumnar := false
		maxLen := 0
		keys := make([]string, 0)
		
		for k, val := range v {
			if valSlice, ok := getAsSlice(val); ok {
				isColumnar = true
				if len(valSlice) > maxLen {
					maxLen = len(valSlice)
				}
				keys = append(keys, k)
			}
		}

		if isColumnar {
			// Convert Columnar to Rows
			for i := 0; i < maxLen; i++ {
				row := make(map[string]interface{})
				for _, k := range keys {
					valSlice, _ := getAsSlice(v[k])
					if i < len(valSlice) {
						row[k] = valSlice[i]
					} else {
						row[k] = nil
					}
				}
				rows = append(rows, row)
			}
		} else {
			// Iterate over map keys/values
			for k, val := range v {
				rows = append(rows, map[string]interface{}{
					"key":   k,
					"value": val,
					k:       val,
				})
			}
		}

	case []interface{}:
		// List of items
		for _, item := range v {
			if mapItem, ok := item.(map[string]interface{}); ok {
				rows = append(rows, mapItem)
			} else {
				// Wrap primitives
				rows = append(rows, map[string]interface{}{"value": item})
			}
		}
	
	default:
		// Attempt slice reflection for other types
		if valSlice, ok := getAsSlice(v); ok {
			for _, item := range valSlice {
				if mapItem, ok := item.(map[string]interface{}); ok {
					rows = append(rows, mapItem)
				} else {
					rows = append(rows, map[string]interface{}{"value": item})
				}
			}
		} else {
			return fmt.Errorf("unsupported source type for looping: %T", v)
		}
	}

	fmt.Printf("[ForEach] Total iterations: %d\n", len(rows))

	// 3. Execution Loop
	// Get downstream nodes from the FIRST output (assuming standard usage)
	
	components := ctx.FiberCtx.Locals("components").([]Component)
	workflowId := ctx.CurrentComponent.WorkflowId
	
	// Find current component definition to get outputs
	var outputConnections []Connection
	var completedConnections []Connection

	for _, comp := range components {
		if comp.ID == ctx.CurrentComponent.ID {
			if out, ok := comp.Outputs["output_1"]; ok { // Standard output (Loop Body)
				outputConnections = out.Connections
			}
			
			if out, ok := comp.Outputs["output_2"]; ok { // Finished output
				completedConnections = out.Connections
			}
			break
		}
	}

	if len(outputConnections) == 0 {
		return fmt.Errorf("no connected nodes to loop over")
	}

	db := ctx.DB
	search := ctx.Search

	for i, row := range rows {
		fmt.Printf("[ForEach] Iteration %d\n", i+1)
		
		// Update Extras with row data
		// Flatten row into extra keys
		for k, val := range row {
			ctx.Extras[k] = val
		}
		// Also provide a scoped "row" or "item" object?
		ctx.Extras["item"] = row
		ctx.Extras["loop_index"] = i
		
		// Update Locals to persist (since InternalFlow uses Locals)
		ctx.FiberCtx.Locals("wfExtras", ctx.Extras)
		
		fmt.Printf("[ForEach Debug] Iteration %d Variables: key=%v, value=%v\n", i+1, ctx.Extras["key"], ctx.Extras["value"])

		// Trigger downstream nodes
		for _, conn := range outputConnections {
			nextNodeId := 0
			fmt.Sscanf(conn.Node, "%d", &nextNodeId)
			
			// Find the component
			var nextComp *Component
			// Need to find slice index to modify IsRun
			var compIndex int = -1
			
			for idx := range components {
				if components[idx].ID == nextNodeId && components[idx].WorkflowId == workflowId {
					nextComp = &components[idx]
					compIndex = idx
					break
				}
			}

			if nextComp != nil {
				// Reset IsRun for the ENTIRE downstream branch from this node
				// Otherwise subsequent iterations won't run.
				// This requires graph traversal.
				// We exclude the current ForEach node ID to prevents resetting itself if there's a loop-back connection
				resetDownstreamIsRun(components, nextNodeId, workflowId, ctx.CurrentComponent.ID)

				// Execute
				err := InternalFlow(ctx.FiberCtx, components[compIndex], workflowId, nextNodeId, db, search)
				if err != nil {
					// Continue or break? Usually break on error
					fmt.Printf("[ForEach] Error in iteration %d: %v\n", i, err)
					return err
				}
			}
		}
	}

	// 4. Trigger "Completed" output (output_2) if connected
	if len(completedConnections) > 0 {
		fmt.Printf("[ForEach] Loop finished, triggering output_2\n")
		for _, conn := range completedConnections {
			nextNodeId := 0
			fmt.Sscanf(conn.Node, "%d", &nextNodeId)

			// Execute using same logic (no reset needed for continuation)
			for idx := range components {
				if components[idx].ID == nextNodeId && components[idx].WorkflowId == workflowId {
					// We might need to ensure IsRun is false if it was somehow run before (loop cycle?)
					// But usually 'Completed' path is new content.
					err := InternalFlow(ctx.FiberCtx, components[idx], workflowId, nextNodeId, db, search)
					if err != nil {
						return err
					}
					break
				}
			}
		}
	}

	return nil
}

// Helper to get slice from interface
func getAsSlice(v interface{}) ([]interface{}, bool) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Slice {
		length := val.Len()
		result := make([]interface{}, length)
		for i := 0; i < length; i++ {
			result[i] = val.Index(i).Interface()
		}
		return result, true
	}
	return nil, false
}

// DFS string to reset IsRun
func resetDownstreamIsRun(components []Component, startNodeId int, workflowId int, excludeNodeId int) {
	// Build map for fast lookup
	compMap := make(map[int]*Component)
	for i := range components {
		if components[i].WorkflowId == workflowId {
			compMap[components[i].ID] = &components[i]
		}
	}

	visited := make(map[int]bool)
	var queue []int
	queue = append(queue, startNodeId)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if visited[curr] {
			continue
		}
		visited[curr] = true
		
		// Skip the excludeNodeId (usually the For Each node itself)
		if curr == excludeNodeId {
			continue
		}

		if comp, ok := compMap[curr]; ok {
			comp.IsRun = false // Reset!
			
			// Add children
			for _, out := range comp.Outputs {
				for _, conn := range out.Connections {
					var nextId int
					fmt.Sscanf(conn.Node, "%d", &nextId)
					queue = append(queue, nextId)
				}
			}
		}
	}
}
