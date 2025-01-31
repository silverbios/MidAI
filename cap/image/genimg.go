package gentext

import (
	"MidAI/auth"
	model "MidAI/models"
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
)

// Message represents a single message in the conversation history
//type Message struct {
//	Prompt string `json:"prompt"` // Role of the message (user or assistant)
//}

// RequestBody is the request body for the AI API
type RequestBody struct {
	Prompt string `json:"prompt"` // Conversation history
}

// ApiResponse is the response from the AI API
type ApiResponse struct {
	Result struct {
		Response string `json:"image"` // Response from the assistant
	} `json:"result"`
}

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
	model.PrintModelsTable(models)

	selectedModel, err := model.SelectModel(models)
	if err != nil {
		// Select a random model from the list with the "Text-to-Image" capability
		var ids []int
		for i, model := range models {
			if model.Task.Capability == "Text-to-Image" {
				ids = append(ids, i)
			}
		}
		if len(ids) == 0 {
			fmt.Println("No models with the 'Text-to-Image' capability available")
			return
		}
		randomIndex := rand.Intn(len(ids))       // Generate a random index within the range of IDs
		selectedModel = models[ids[randomIndex]] // Select the model at that index
		// If there's an error selecting a model, print the error and exit
		fmt.Printf("\nWe select the \"%s\" for you.\n", path.Base(selectedModel.Name))
		//fmt.Println("Error selecting model:", err)
		//return
	}

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

		// Build the API URL
		apiURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s", config.AccountID, selectedModel.Name)

		// Build the request body
		requestBody := RequestBody{
			Prompt: userInput,
		}

		// Get the assistant's response
		imageData, err := getAssistantResponse(apiURL, config.Token, requestBody)
		if err != nil {
			// If there's an error getting the assistant's response, panic
			panic(err)
		}

		// Save image
		if err := saveBase64Image(imageData, "generated_image.png"); err != nil {
			fmt.Println("Error saving image:", err)
			continue
		}

		fmt.Println("✅ Image saved as 'generated_image.png'")
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

// Save base64-encoded image to a file
func saveBase64Image(base64Data, filename string) error {
	imageBytes, err := decodeBase64(base64Data)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, imageBytes, 0644)
}

// Decode base64 string
func decodeBase64(base64String string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, errors.New("failed to decode base64 image")
	}
	return decoded, nil
}

// getAssistantResponse makes an API call to the AI API to get the assistant's response
func getAssistantResponse(url, token string, requestBody RequestBody) (string, error) {
	// Marshal the request body to JSON
	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	fmt.Println("\n🟢 Sending Request:")
	fmt.Println("🔹 URL:", url)
	fmt.Println("🔹 Headers: Authorization=Bearer <hidden>, Content-Type=application/json")
	fmt.Println("🔹 Payload:", string(body)) // Print request body as string

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

	// Log the raw response
	fmt.Println("\n🔵 Received Response:")
	fmt.Println("🔹 Status Code:", res.Status)
	fmt.Println("🔹 Headers:", res.Header)
	fmt.Println("🔹 Body:", string(respBody)) // Print full response

	// Unmarshal the response body to the ApiResponse struct
	var apiResponse ApiResponse
	if err := json.Unmarshal(respBody, &apiResponse); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w\nRaw Response: %s", err, string(respBody))
	}

	// Return the assistant's response
	return apiResponse.Result.Response, nil
}
