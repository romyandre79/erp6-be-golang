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
	input := ExternalPluginInput{
		Params: ctx.Params,
	}

	inputJson, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal plugin input: %v", err)
	}

	// 2. Setup Command
	cmd := exec.Command(cmdPath)
	cmd.Stdin = bytes.NewReader(inputJson)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// 3. Execute
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("plugin execution failed: %v, stderr: %s", err, errBuf.String())
	}

	// 4. Parse Output
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
	wfEngine := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, WorkflowEngine{
		DataInputNode: "",
		ResultNode:    output.Result,
	})
	ctx.FiberCtx.Locals("wfEngine", wfEngine)

	return nil
}
