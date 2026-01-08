package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type TescoStore struct {
	Name        string `json:"name"`
	City        string `json:"city"`
	Country     string `json:"country,omitempty"`
	Address     string `json:"address"`
	Website     string `json:"website"`
	Phone       string `json:"phone"`
}

func main() {

	// Create MCP server
	s := server.NewMCPServer(
		"mcp-tesco-server",
		"0.0.0",
	)

	// Pizzerias by city tool
	storesByCityTool := mcp.NewTool("get_stores_by_city",
		mcp.WithDescription("Get list of Tesco stores in a specific city"),
		mcp.WithString("city",
			mcp.Required(),
			mcp.Description("Name of the city to search for Tesco stores"),
		),
	)
	s.AddTool(storesByCityTool, storesByCityHandler)

	// Pizzerias by country tool
	storesByCountryTool := mcp.NewTool("get_stores_by_country",
		mcp.WithDescription("Get list of Tesco stores in a specific country"),
		mcp.WithString("country",
			mcp.Required(),
			mcp.Description("Name of the country to search for Tesco stores"),
		),
	)
	s.AddTool(storesByCountryTool, storesByCountryHandler)


	// Start the HTTP server
	httpPort := os.Getenv("MCP_HTTP_PORT")
	if httpPort == "" {
		httpPort = "9090"
	}

	log.Println("MCP Files Server is running on port", httpPort)

	// Create a custom mux to handle both MCP and health endpoints
	mux := http.NewServeMux()

	// Add healthcheck endpoint
	mux.HandleFunc("/health", healthCheckHandler)

	// Add MCP endpoint
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	)

	// Register MCP handler with the mux
	mux.Handle("/mcp", httpServer)

	// Start the HTTP server with custom mux
	log.Fatal(http.ListenAndServe(":"+httpPort, mux))
}

func loadTescoStoresData() ([]TescoStore, error) {
	dataPath := "data/tesco_directory.json"

	content, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("error reading Tesco stores data: %v", err)
	}

	var stores []TescoStore
	err = json.Unmarshal(content, &stores)
	if err != nil {
		return nil, fmt.Errorf("error parsing Tesco stores data: %v", err)
	}

	return stores, nil
}

func storesByCityHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	cityArg, exists := args["city"]
	if !exists || cityArg == nil {
		return nil, fmt.Errorf("missing required parameter 'city'")
	}

	city, ok := cityArg.(string)
	if !ok {
		return nil, fmt.Errorf("parameter 'city' must be a string")
	}

	// Convert to lowercase for comparison
	cityLower := strings.ToLower(city)

	stores, err := loadTescoStoresData()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var matchingStores []TescoStore
	for _, store := range stores {
		if strings.ToLower(store.City) == cityLower {
			matchingStores = append(matchingStores, store)
		}
	}

	result, err := json.MarshalIndent(matchingStores, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error formatting results: %v", err)), nil
	}

	log.Printf("Found %d Tesco stores in city: %s", len(matchingStores), city)
	return mcp.NewToolResultText(string(result)), nil
}

func storesByCountryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	countryArg, exists := args["country"]
	if !exists || countryArg == nil {
		return nil, fmt.Errorf("missing required parameter 'country'")
	}

	country, ok := countryArg.(string)
	if !ok {
		return nil, fmt.Errorf("parameter 'country' must be a string")
	}

	// Convert to lowercase for comparison
	countryLower := strings.ToLower(country)

	stores, err := loadTescoStoresData()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var matchingStores []TescoStore
	for _, store := range stores {
		// Check both explicit country field and infer from city names
		storeCountry := store.Country


		if strings.ToLower(storeCountry) == countryLower {
			// Set the inferred country for consistency
			if store.Country == "" {
				store.Country = storeCountry
			}
			matchingStores = append(matchingStores, store)
		}
	}

	result, err := json.MarshalIndent(matchingStores, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error formatting results: %v", err)), nil
	}

	log.Printf("Found %d Tesco stores in country: %s", len(matchingStores), country)
	return mcp.NewToolResultText(string(result)), nil
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status": "healthy",
		"server": "mcp-pizzerias-server",
	}
	json.NewEncoder(w).Encode(response)
}
