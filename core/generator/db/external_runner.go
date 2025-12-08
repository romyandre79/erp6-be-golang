package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// ExternalPluginInput is the data structure sent to the plugin
type ExternalPluginInput struct {
	Params []WorkflowDetailResult `json:"params"`
	// We can add more context here if needed, but Params is usually enough
}

// ExternalPluginOutput is the expected output from the plugin
type ExternalPluginOutput struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
}

// RunExternalPlugin executes a standalone binary as a plugin
func RunExternalPlugin(cmdPath string, ctx *WorkflowContext) error {
	// 1. Prepare Input
	params := make([]WorkflowDetailResult, len(ctx.Params))
	copy(params, ctx.Params)

	// Map of existing param names to avoid duplicates
	existingParams := make(map[string]int)
	for i, p := range params {
		existingParams[p.InputName] = i
	}

	// Helper to upsert parameters
	upsertParam := func(key, value string) {
		if value == "" {
			return
		}
		if idx, exists := existingParams[key]; exists {
			params[idx].CompValue = value
			fmt.Printf("[ExternalPlugin] Substituted param '%s' with runtime value: '%s'\n", key, value)
		} else {
			newParam := WorkflowDetailResult{
				InputName: key,
				CompValue: value,
				Label:     key,
			}
			params = append(params, newParam)
			existingParams[key] = len(params) - 1
		}
	}

	// 1. Update from FormValue (Covers Query and Body for existing keys)
	for _, p := range params {
		val := ctx.FiberCtx.FormValue(p.InputName)
		upsertParam(p.InputName, val)
	}

	// 2. Discover NEW dynamic parameters from Query String
	for k, v := range ctx.FiberCtx.Queries() {
		upsertParam(k, v)
	}

	// 3. Discover NEW dynamic parameters from Multipart Form
	form, err := ctx.FiberCtx.MultipartForm()
	if err == nil && form != nil {
		for k, v := range form.Value {
			if len(v) > 0 {
				upsertParam(k, v[0])
			}
		}
	}

	// 4. Discover NEW dynamic parameters from PostArgs (x-www-form-urlencoded)
	ctx.FiberCtx.Request().PostArgs().VisitAll(func(key, value []byte) {
		upsertParam(string(key), string(value))
	})

	input := ExternalPluginInput{
		Params: params,
	}

	inputJson, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal plugin input: %v", err)
	}

	// 2. Setup Command
	fmt.Printf("[ExternalPlugin] Executing: %s\n", cmdPath)
	fmt.Printf("[ExternalPlugin] Input: %s\n", string(inputJson))

	cmd := exec.Command(cmdPath)
	cmd.Stdin = bytes.NewReader(inputJson)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// 3. Execute
	if err := cmd.Run(); err != nil {
		fmt.Printf("[ExternalPlugin] Error: %v\nStderr: %s\n", err, errBuf.String())
		return fmt.Errorf("plugin execution failed: %v, stderr: %s", err, errBuf.String())
	}

	// 4. Parse Output
	fmt.Printf("[ExternalPlugin] Output: %s\n", outBuf.String())

	var output ExternalPluginOutput
	if err := json.Unmarshal(outBuf.Bytes(), &output); err != nil {
		// If not JSON, maybe it just printed the result directly?
		// Let's assume strict JSON for now to be safe.
		return fmt.Errorf("failed to parse plugin output: %v, raw: %s", err, outBuf.String())
	}

	if output.Error != "" {
		return fmt.Errorf("plugin reported error: %s", output.Error)
	}

	// 5. Update Workflow Engine
	// We need to append the result to the wfEngine local
	// Note: InternalFlow handles appending result if component didn't.
	// But here we are appending. This matches InternalFlow logic.

	// FIX: We shouldn't necessarily append a NEW item if InternalFlow is managing the array.
	// InternalFlow logic:
	// if len(wfEngine) > 0 { lastResult := wfEngine[len(wfEngine)-1] ... }
	// if len(wfEngine) > 0 && wfEngine[len(wfEngine)-1].ComponentName == "" { ... }

	// If checking InternalFlow (lines 313 in executeflow.go), it updates the LAST item if ComponentName is empty.
	// But ExecuteFlow initializes empty wfEngine.
	// Actually InternalFlow creates the entry generally?
	// No, InternalFlow (line 325) appends result if not present.

	// If we append here, we might duplicate?
	// InternalFlow says: "Get any result that might have been set by the component"
	// So we should append here.

	wfEngine := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, WorkflowEngine{
		// DataInputNode: input, // Can potentially save input too
		ResultNode: output.Result,
	})
	ctx.FiberCtx.Locals("wfEngine", wfEngine)

	return nil
}
