# Docker Compose MCP ToolKit + MCP Gateway (Streamable HTTP Transport)
> - Get the list of the tools available in the Docker Desktop **MCP ToolKit** using Docker Compose and Streamable HTTP transport

## Requirements
- Docker Desktop installed with MCP Toolkit
- Golang installed

###  MCP Servers
- Open Docker Desktop
- Install MCP Servers, for example:
  - DuckDuckGo MCP Server (you don't need credentials for this one, but there is a limit of requests per day)
  - Brave Search, you need an API key (you can get a free one with limited requests per month: https://brave.com/search/api/)

# Run
```bash
docker compose up --build
```

This will display only the DuckDuckGo tools from the MCP Toolkit:
```yaml
  mcp-gateway:
    # mcp-gateway secures your MCP servers
    image: docker/mcp-gateway:latest
    command:
      - --port=9011
      - --transport=streaming
      - --servers=duckduckgo
      - --verbose
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

