# Docker Desktop MCP Gateway + AI Agent

The code example demonstrates how to create an AI agent that interacts with the Docker Desktop MCP Gateway to utilize its tools.

The system instruction and the user message have been externalized as environment variables in the `compose.yml` file for easier configuration.

In this example, the AI agent is instructed to search for good pizzerias in Lyon, France, using the tools available (DuckDuckGo MCP server) in the MCP Gateway.

## Requirements
- Docker Desktop installed with MCP Toolkit
- Golang installed

###  MCP Servers
- Open Docker Desktop
- Install MCP Servers, for example:
  - DuckDuckGo MCP Server (you don't need credentials for this one, but there is a limit of requests per day)
  - Brave Search, you need an API key (you can get a free one with limited requests per month: https://brave.com/search/api/)
  
## Run

```bash
docker compose up --build --no-log-prefix
```