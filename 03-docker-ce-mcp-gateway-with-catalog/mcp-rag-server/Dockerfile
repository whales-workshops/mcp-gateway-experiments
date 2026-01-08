FROM golang:1.25.4-alpine AS builder
WORKDIR /app
#COPY go.mod .
COPY go.mod go.sum ./

COPY main.go .
COPY rag/ ./rag/

RUN <<EOF
go mod download
go mod tidy 
go build
EOF

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
WORKDIR /app
#COPY data/ ./data/
COPY --from=builder /app/mcp-rag-server .
COPY *.md /app/
ENTRYPOINT ["./mcp-rag-server"]
