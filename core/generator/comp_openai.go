package generator

import (
	"context"
	"fmt"
	"strings"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"google.golang.org/genai"
)

func init() {
	RegisterComponent("comp_openai", func(ctx *WorkflowContext) error {
		return handleOpenAI(ctx)
	})
	// Alias for easier usage or back-compat logic if needed
	RegisterComponent("OpenAI", func(ctx *WorkflowContext) error {
		return handleOpenAI(ctx)
	})
}

func handleOpenAI(ctx *WorkflowContext) error {
	var (
		provider   = "openai"
		token      string
		baseURL    string
		model      string
		userPrompt string
	)

	// Extract parameters
	for _, p := range ctx.Params {
		val := strings.TrimSpace(ResolveParam(ctx.FiberCtx, p.CompValue))
		if val == "" {
			val = ResolveParam(ctx.FiberCtx, p.InputName)
			if strings.HasPrefix(val, "$") {
				val = ""
			}
		}

		switch strings.ToLower(p.InputName) {
		case "provider":
			if val != "" {
				provider = val
			}
		case "token", "api_key", "apikey":
			token = val
		case "base_url", "baseurl":
			baseURL = val
		case "model":
			model = val
		case "user_prompt", "prompt", "message":
			userPrompt = val
		}
	}

	// Validate required parameters
	if token == "" {
		return fmt.Errorf("token is required")
	}
	if userPrompt == "" {
		return fmt.Errorf("user_prompt is required")
	}

	fmt.Printf("[CompOpenAI] Executing Provider: '%s', Model: '%s'\n", provider, model)

	var result string
	var err error

	switch strings.ToLower(provider) {
	case "gemini":
		result, err = RunGemini(token, model, userPrompt, baseURL)
	case "claude", "anthropic":
		result, err = RunAnthropic(token, model, userPrompt, baseURL)
	default:
		// Default to OpenAI
		if model == "" {
			model = openai.ChatModelGPT3_5Turbo
		}
		result, err = RunOpenAI(token, baseURL, model, userPrompt)
	}

	if err != nil {
		fmt.Printf("[CompOpenAI] Error: %v\n", err)
		return err
	}

	// Store result
	if ctx.Extras == nil {
		ctx.Extras = make(map[string]interface{})
	}
	ctx.Extras["result"] = result
	ctx.Extras["ai_response"] = result // Alias

	// Also store in legacy WorkflowEngine result structure if needed
	wm := WorkflowEngine{
		ResultNode: map[string]interface{}{
			"result": result,
		},
	}
	
	if ctx.FiberCtx != nil {
		wfEngine, _ := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
		wfEngine = append(wfEngine, wm)
		ctx.FiberCtx.Locals("wfEngine", wfEngine)
	} else {
		wfEngine, _ := ctx.Extras["wfEngine"].([]WorkflowEngine)
		wfEngine = append(wfEngine, wm)
		ctx.Extras["wfEngine"] = wfEngine
	}

	return nil
}

func RunOpenAI(token, baseURL, model, userPrompt string) (string, error) {
	opts := []option.RequestOption{
		option.WithAPIKey(token),
	}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	client := openai.NewClient(opts...)

	// Use Chat Completions API for compatibility with Ollama and other OpenAI-compatible endpoints
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(userPrompt),
		},
		Model: openai.ChatModel(model),
	})

	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %v", err)
	}

	if len(chatCompletion.Choices) > 0 {
		return chatCompletion.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from API")
}

func RunGemini(token, model, userPrompt, baseURL string) (string, error) {
	ctx := context.Background()
	// Set default model if empty
	if model == "" {
		model = "gemini-1.5-flash"
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: token,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %v", err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(userPrompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("Gemini API error: %v", err)
	}

	return result.Text(), nil
}

func RunAnthropic(token, model, userPrompt, baseURL string) (string, error) {
	if model == "" {
		model = "claude-3-5-sonnet-20240620"
	}

	client := anthropic.NewClient(token)
	
	resp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.Model(model),
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(userPrompt),
		},
		MaxTokens: 1024,
	})
	
	if err != nil {
		return "", fmt.Errorf("Anthropic API error: %v", err)
	}
	
	if len(resp.Content) > 0 {
		return *resp.Content[0].Text, nil
	}
	return "", fmt.Errorf("no content in response")
}
