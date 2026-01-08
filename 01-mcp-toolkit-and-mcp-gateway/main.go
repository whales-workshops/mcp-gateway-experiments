package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
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

}
