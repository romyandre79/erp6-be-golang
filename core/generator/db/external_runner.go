package generator

import (
	"bytes"
	"encoding/json"
	"erp6-be-golang/core/configs"
	"fmt"
	"os/exec"
)

// ExternalPluginInput is the data structure sent to the plugin
type ExternalPluginInput struct {
	Params   []WorkflowDetailResult `json:"params"`
	DBConfig *DBConfig              `json:"db_config"`
}

type DBConfig struct {
	Driver string `json:"driver"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Name   string `json:"name"`
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
	// Deep copy and resolve
	for i, p := range ctx.Params {
		params[i] = p
		params[i].CompValue = ResolveParam(ctx.FiberCtx, p.CompValue)
	}

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

	// 1. Use resolved parameters from workflow definition
	for _, p := range params {
		val := ResolveParam(ctx.FiberCtx, p.CompValue)
		upsertParam(p.InputName, val)
	}

	// Note: Sections 2-4 (Query/Form/PostArgs discovery) removed to prevent
	// parameters from previous nodes overwriting current node's parameters

	// Prepare DB Config
	dbConfig := &DBConfig{
		Driver: configs.ConfigApps.DBDriver,
		Host:   configs.ConfigApps.DBHost,
		Port:   configs.ConfigApps.DBPort,
		User:   configs.ConfigApps.DBUser,
		Pass:   configs.ConfigApps.DBPass,
		Name:   configs.ConfigApps.DBName,
	}

	input := ExternalPluginInput{
		Params:   params,
		DBConfig: dbConfig,
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

	// Print stderr even on success (for debug messages)
	if errBuf.Len() > 0 {
		fmt.Printf("[ExternalPlugin] Stderr: %s\n", errBuf.String())
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
