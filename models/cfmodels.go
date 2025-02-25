package model

import (
	"MidAI/auth"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// Model struct to hold model data with ID, Name, and Description fields
type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Task        struct {
		Capability string `json:"name"`
	} `json:"task"`
}

// ModelsResponse struct to hold the response from the Cloudflare API
type ModelsResponse struct {
	Success bool    `json:"success"`
	Result  []Model `json:"result"`
}

// GetAvailableModels fetches available models from the Cloudflare API
func GetAvailableModels(config auth.Config) ([]Model, error) {
	// Construct the API URL using the account ID from the config
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/models/search", config.AccountID)

	// Create a new HTTP GET request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the authorization header using the API token from the config
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Token))
	client := &http.Client{}

	// Execute the request
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer response.Body.Close()

	// Check if the response status is OK (200)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", response.Status)
	}

	// Decode the JSON response into the ModelsResponse struct
	var modelsResponse ModelsResponse
	if err = json.NewDecoder(response.Body).Decode(&modelsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if the API response indicates success
	if !modelsResponse.Success {
		return nil, fmt.Errorf("failed to fetch models")
	}

	// Return the list of models from the response
	return modelsResponse.Result, nil
}

// Function to determine the terminal width dynamically
func getTerminalWidth() int {
	width, _, err := term.GetSize(int(syscall.Stdout))
	if err != nil {
		return 120 // Fallback width if terminal size can't be determined
	}
	return width
}

// PrintModelsTable prints the list of models in a table format
func PrintModelsTable(models []Model) {
	// Check if the models slice is nil or empty
	if models == nil {
		fmt.Println("No models available to display.")
		return
	}

	// Create a strings.Builder to build the table output
	builder := &strings.Builder{}
	// Print the table header
	fmt.Fprintln(builder, "+----+---------------------------------------------+---------------------+----------------------------------------------------------------------------------------+")
	fmt.Fprintln(builder, "| #  |                  Model Name                 |     Model Type      |                                      Description                                       |")
	fmt.Fprintln(builder, "+----+---------------------------------------------+---------------------+----------------------------------------------------------------------------------------+")

	// Iterate over the models and print each row
	for i, model := range models {
		printRow(builder, i+1, model.Name, model.Task.Capability, model.Description)
		// Print the table footer
		fmt.Fprintln(builder, "+----+---------------------------------------------+---------------------+----------------------------------------------------------------------------------------+")

	}
	// Output the entire table
	fmt.Print(builder.String())
}

// printRow handles printing a single row of the table, truncating the description if necessary
func printRow(builder *strings.Builder, number int, name string, capability string, description string) {
	// Define fixed column widths
	idWidth := 3          // ID column width (fixed)
	nameWidth := 30       // Name column width (fixed)
	capabilityWidth := 25 // Capability column width (fixed)

	// Dynamically calculate the remaining space for description
	descriptionWidth := getTerminalWidth() - (idWidth + nameWidth + capabilityWidth + 16) // 16 accounts for column dividers

	// Ensure description width is not negative
	if descriptionWidth < 16 {
		descriptionWidth = 16 // Minimum width to prevent overflow
	}

	// Trim and truncate fields
	truncatedName := truncate(strings.TrimSpace(path.Base(name)), nameWidth)
	truncatedCapability := truncate(strings.TrimSpace(capability), capabilityWidth)
	truncatedDescription := truncate(strings.TrimSpace(description), descriptionWidth)

	// Print the formatted row to the builder
	fmt.Fprintf(builder, "| %-*d | %-*s | %-*s | %-*s |\n",
		idWidth, number,
		nameWidth, truncatedName,
		capabilityWidth, truncatedCapability,
		descriptionWidth, truncatedDescription)
}

// truncate ensures that text is cut off with '...' if it's too long
func truncate(str string, maxLength int) string {
	// Check if the string length exceeds the specified max length
	if len(str) > maxLength {
		// Return the truncated string with '...' at the end
		return str[:maxLength-3] + "..."
	}
	// Return the original string if no truncation is needed
	return str
}

// SelectModel handles model selection based on user input
func SelectModel(models []Model) (Model, error) {
	// Check if models slice is nil or empty
	if models == nil {
		return Model{}, fmt.Errorf("no models available for selection")
	}

	var selectedModelIndex int
	// Prompt the user to enter the model number
	fmt.Print("Enter the number corresponding to the model you'd like to use: ")
	_, err := fmt.Scanln(&selectedModelIndex)

	// Validate the user input
	if err != nil || selectedModelIndex < 1 || selectedModelIndex > len(models) {
		return Model{}, fmt.Errorf("invalid selection")
	}

	// TODO: Add validation to prevent selecting a model that was not in the original list
	if selectedModelIndex > len(models) {
		return Model{}, fmt.Errorf("model number %d is not in the list", selectedModelIndex)
	}

	// Return the selected model based on the user's input
	return models[selectedModelIndex-1], nil
}
