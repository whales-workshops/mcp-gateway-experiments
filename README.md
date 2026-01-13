# MCP Gateway Experiments

This repository contains three demonstrations exploring the Docker MCP Gateway and its integration with AI agents and MCP servers.

## Demo 01: MCP Toolkit and MCP Gateway

The first demo shows how to connect to Docker Desktop's MCP Toolkit using the MCP Gateway with Streamable HTTP transport. A Go client connects to the gateway (running on port 9011) and retrieves the list of available tools from the DuckDuckGo MCP server.

**Key Features:**
- Docker Compose setup with MCP Gateway
- Streamable HTTP transport
- Connection to Docker Desktop MCP Toolkit
- Tool discovery from DuckDuckGo server

## Demo 02: MCP Toolkit, MCP Gateway and Function Calling

This demo builds an AI agent that uses function calling to intelligently interact with MCP servers through the gateway. The Go-based agent retrieves available tools, analyzes user queries using an LLM (Lucy model), detects which tools to call, and executes them through the MCP Gateway.

**Key Features:**
- AI agent with function calling capabilities
- Integration with OpenAI-compatible LLM (Lucy model)
- Automatic tool detection from user prompts
- Tool execution through MCP Gateway
- Example: "Search for good pizzerias in Lyon, France"

## Demo 03: Docker CE MCP Gateway with Custom Catalog

This demo demonstrates running MCP Gateway with Docker CE (Community Edition) using a custom catalog to connect to custom MCP servers. It includes two custom MCP servers (RAG and Tesco) with the MCP Inspector for debugging.

**Key Features:**
- Custom MCP server catalog configuration
- Custom RAG server with embeddings model
- Custom Tesco server with data management
- MCP Inspector integration for debugging
- Streamable HTTP transport
- Docker CE compatibility
