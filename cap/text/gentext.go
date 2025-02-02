package gentext

import (
	"MidAI/auth"
	model "MidAI/models"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
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

var maxHistory int // Maximum number of messages to keep in the conversation history

func Prompt() {
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

	// Filter only the models with "Text Generation" capability
	var textToImageModels []model.Model
	for _, m := range models {
		if m.Task.Capability == "Text Generation" {
			textToImageModels = append(textToImageModels, m)
		}
	}
	// Print the table of only Text Generation models
	model.PrintModelsTable(textToImageModels)

	selectedModel, err := model.SelectModel(models)
	if err != nil {
		// Select a random model from the list with the "Text Generation" capability
		var ids []int
		for i, model := range models {
			if model.Task.Capability == "Text Generation" {
				ids = append(ids, i)
			}
		}
		if len(ids) == 0 {
			fmt.Println("No models with the 'Text Generation' capability available")
			return
		}
		randomIndex := rand.Intn(len(ids))       // Generate a random index within the range of IDs
		selectedModel = models[ids[randomIndex]] // Select the model at that index
		// If there's an error selecting a model, print the error and exit
		fmt.Printf("\nWe select the \"%s\" for you.\n", path.Base(selectedModel.Name))
		//fmt.Println("Error selecting model:", err)
		//return
	}

	fmt.Println("\n1. History size: 1  2. History size: 2  3. History size: 3  4. History size: 4  5. History size: 5  6. History size: 6 (default)  7. History size: 7  8. History size: 8  9. History size: 9  10. History size: 10")

	_, err = fmt.Scanln(&maxHistory)
	if err != nil || maxHistory < 1 || maxHistory > 10 {
		fmt.Printf("\nInvalid selection. We select the 6 size for you.\n")
		maxHistory = 6
		fmt.Printf("History size set to: %d\n", maxHistory)
		// If there's an error reading the user's input, panic
		//panic(err)
	}

	// Initialize the conversation history
	conversationHistory := make([]Message, 0, maxHistory)

	// Create a reader to read the user's input
	reader := bufio.NewReader(os.Stdin)

	for {
		// Ask the user for their message
		fmt.Print("\nEnter your message for the assistant (or press 'Enter' or type 'q' to exit): ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			// If there's an error reading the user's input, panic
			panic(err)
		}
		userInput = userInput[:len(userInput)-1]

		if userInput == "" || userInput == "q" {
			// If the user presses 'Enter' or types 'q', exit the loop
			fmt.Printf("\nHave a nice time:) Goodbye!\n")
			break
		}

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
