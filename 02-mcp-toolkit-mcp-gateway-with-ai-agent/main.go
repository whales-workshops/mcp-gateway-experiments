package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	ctx := context.Background()
	mcpClient, err := client.NewStreamableHttpClient(
		os.Getenv("MCP_HOST"), // Use environment variable for MCP host
	)
	//defer mcpClient.Close()
	if err != nil {
		fmt.Println("🔴 Failed to create MCP client:", err)
		panic(err)
	}

	// Start the connection to the server
	err = mcpClient.Start(ctx)
	if err != nil {
		fmt.Println("🔴 Failed to start MCP client:", err)
		panic(err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "bob",
		Version: "0.0.0",
	}

	result, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		fmt.Println("🔴 Failed to initialize MCP client:", err)
		panic(err)
	}
	fmt.Println("Streamable HTTP client connected & initialized with server!", result)

	toolsRequest := mcp.ListToolsRequest{}
	mcpTools, err := mcpClient.ListTools(ctx, toolsRequest)
	if err != nil {
		panic(err)
	}
	fmt.Println("Available Tools:")
	for _, tool := range mcpTools.Tools {
		fmt.Printf("🛠️ Tool: %s\n", tool.Name)
		fmt.Printf("  Description: %s\n", tool.Description)
	}

	modelRunnerBaseUrl := os.Getenv("MODEL_RUNNER_BASE_URL")

	clientEngine := openai.NewClient(
		option.WithBaseURL(modelRunnerBaseUrl),
		option.WithAPIKey(""),
	)

	modelRunnerToolsModel := os.Getenv("MODEL_RUNNER_TOOLS_MODEL")
	if modelRunnerToolsModel == "" {
		panic("MODEL_RUNNER_TOOLS_MODEL environment variable is not set")
	}

	openAITools := ConvertMCPToolsToOpenAITools(mcpTools)

	toolsCompletionParams := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{},
		//ParallelToolCalls: openai.Bool(true),
		ParallelToolCalls: openai.Bool(false),
		Tools:             openAITools,
		Model:             modelRunnerToolsModel,
		Temperature:       openai.Opt(0.0),
	}

	systemToolsInstructions := os.Getenv("SYSTEM_INSTRUCTION")

	toolsCompletionParams.Messages = []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemToolsInstructions),
		openai.UserMessage(os.Getenv("USER_MESSAGE")),
	}

	fmt.Println("⏳ Running tools completion...")
	// Make initial Tool completion request
	// TOOLS COMPLETION:
	completion, err := clientEngine.Chat.Completions.New(ctx, toolsCompletionParams)
	if err != nil {
		fmt.Printf("😡 Tools completion error: %v\n", err)
		panic(err)
	}

	fmt.Println("🛠️ Tools completion received")
	detectedToolCalls := completion.Choices[0].Message.ToolCalls
	// fetch_content
	// search
	if len(detectedToolCalls) > 0 {
		for _, toolCall := range detectedToolCalls {
			// Display the detected tool call
			fmt.Println("💡 tool detection:", toolCall.Function.Name, toolCall.Function.Arguments)

			// Parse the tool arguments from JSON string
			var args map[string]any
			err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
			if err != nil {
				fmt.Println("🔴 Error parsing tool arguments:", err)
				panic(err)
				//continue
			}

			// NOTE: Call the MCP tool with the arguments
			request := mcp.CallToolRequest{}
			request.Params.Name = toolCall.Function.Name
			request.Params.Arguments = args

			toolResponse, err := mcpClient.CallTool(ctx, request)

			if err != nil {
				fmt.Println("🔴 Error calling tool:", err)
				continue
			} else {
				if toolResponse != nil && len(toolResponse.Content) > 0 {
					result := toolResponse.Content[0].(mcp.TextContent).Text
					fmt.Printf("✅ Tool %s executed successfully, result: %s\n", toolCall.Function.Name, result)
				}
			}

		}
	}

}

func ConvertMCPToolsToOpenAITools(tools *mcp.ListToolsResult) []openai.ChatCompletionToolParam {
	openAITools := make([]openai.ChatCompletionToolParam, len(tools.Tools))
	for i, tool := range tools.Tools {

		openAITools[i] = openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        tool.Name,
				Description: openai.String(tool.Description),
				Parameters: openai.FunctionParameters{
					"type":       "object",
					"properties": tool.InputSchema.Properties,
					"required":   tool.InputSchema.Required,
				},
			},
		}
	}
	return openAITools
}
