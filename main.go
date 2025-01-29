package main

import (
	"MidAI/auth"
	"MidAI/models"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Message represents a single message in the conversation history
type Message struct {
	Role    string `json:"role"`    // Role of the message (user or assistant)
	Content string `json:"content"` // Content of the message
}

// RequestBody is the request body for the AI API
type RequestBody struct {
	Messages []Message `json:"messages"` // Conversation history
}

// ApiResponse is the response from the AI API
type ApiResponse struct {
	Result struct {
		Response string `json:"response"` // Response from the assistant
	} `json:"result"`
}

const maxHistory = 6 // Maximum number of messages to keep in the conversation history

func main() {
	// Load the configuration from the user's home directory
	config, err := auth.LoadConfig()
	if err != nil {
		// If the configuration doesn't exist, prompt the user for the necessary information
		config = promptForConfig()
		if err := auth.SaveConfig(config); err != nil {
			// If there's an error saving the configuration, panic
			panic(err)
		}
	}

	// Fetch the list of available models
	models, err := model.GetAvailableModels(config)
	if err != nil {
		// If there's an error fetching the models, print the error and exit
		fmt.Println("Error fetching models:", err)
		return
	}
	model.PrintModelsTable(models)

	// Ask the user to select a model
	selectedModel, err := model.SelectModel(models)
	if err != nil {
		// If there's an error selecting a model, print the error and exit
		fmt.Println("Error selecting model:", err)
		return
	}

	// Initialize the conversation history
	conversationHistory := make([]Message, 0, maxHistory)

	// Create a reader to read the user's input
	reader := bufio.NewReader(os.Stdin)

	for {
		// Ask the user for their message
		fmt.Print("\nEnter your message for the assistant: ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			// If there's an error reading the user's input, panic
			panic(err)
		}
		userInput = userInput[:len(userInput)-1]

		// Add the user's message to the conversation history
		conversationHistory = appendMessage(conversationHistory, "user", userInput)

		// Build the API URL
		apiURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s", config.AccountID, selectedModel.Name)

		// Build the request body
		requestBody := RequestBody{
			Messages: append([]Message{
				{Role: "system", Content: "You are a friendly assistant"},
			}, conversationHistory...),
		}

		// Get the assistant's response
		assistantResponse, err := getAssistantResponse(apiURL, config.Token, requestBody)
		if err != nil {
			// If there's an error getting the assistant's response, panic
			panic(err)
		}

		// Print the assistant's response
		fmt.Printf("\nAssistant's response:\n%s\n", assistantResponse)

		// Add the assistant's response to the conversation history
		conversationHistory = appendMessage(conversationHistory, "assistant", assistantResponse)
	}
}

// promptForConfig prompts the user for the necessary configuration information
func promptForConfig() auth.Config {
	var config auth.Config
	fmt.Print("Enter your Cloudflare Account ID: ")
	fmt.Scanln(&config.AccountID)
	fmt.Print("Enter your Cloudflare API Token: ")
	fmt.Scanln(&config.Token)
	return config
}

// appendMessage appends a message to the conversation history
func appendMessage(history []Message, role, content string) []Message {
	if len(history) >= maxHistory {
		// If the conversation history is full, remove the oldest message
		history = history[1:]
	}
	return append(history, Message{Role: role, Content: content})
}

// getAssistantResponse makes an API call to the AI API to get the assistant's response
func getAssistantResponse(url, token string, requestBody RequestBody) (string, error) {
	// Marshal the request body to JSON
	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	// Set the authorization header
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Make the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// Unmarshal the response body to the ApiResponse struct
	var apiResponse ApiResponse
	if err := json.Unmarshal(respBody, &apiResponse); err != nil {
		return "", err
	}

	// Return the assistant's response
	return apiResponse.Result.Response, nil
}
