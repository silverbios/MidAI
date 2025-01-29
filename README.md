# MidAI - Cloudflare AI Assistant

## Overview
`MidAI` is a cross-platform command-line AI assistant that interacts with Cloudflare's AI API, allowing users to send messages and receive responses in a conversational manner. It runs seamlessly on Windows, macOS, and Linux.
`MidAI` is a command-line AI assistant that interacts with Cloudflare's AI API, allowing users to send messages and receive responses in a conversational manner. The assistant maintains a short history of interactions and supports model selection from available AI models in your Cloudflare account.

## Features
- Authenticate using your Cloudflare API Token and Account ID.
- Retrieve a list of available AI models.
- Select a model dynamically.
- Maintain a short conversation history (up to 6 messages).
- Communicate with Cloudflare AI via API calls.
- Interactive CLI experience.

## Installation
### Prerequisites
- Go 1.18 or later installed.
- A Cloudflare account with API access to AI models.
- Network connectivity to Cloudflare API.

### Clone the Repository
```sh
 git clone https://github.com/your-repo/MidAI.git
 cd MidAI
```

### Build the Application
```sh
go build -o midai
```

## Usage
### Run the Application
```sh
./midai
```

### Initial Configuration
On first launch, the application prompts for:
- Cloudflare Account ID
- Cloudflare API Token

This information is stored securely in the user's home directory.

### Selecting an AI Model
After startup, the application fetches available models from Cloudflare and presents them as a list. Users can select the model they wish to use.

### Interacting with the AI
Once configured, the user can send messages and receive responses in a conversational format:
```sh
Enter your message for the assistant: What is the capital of France?

Assistant's response:
The capital of France is Paris.
```

## Code Structure
- `main.go` - Handles CLI interactions, configuration, and API calls.
- `auth/` - Manages API authentication and configuration.
- `model/` - Fetches and displays available AI models.

## Configuration File
User configuration is stored at:
```
`~/.aiCFtoken.json` (Linux/macOS) or `%USERPROFILE%\.aiCFtoken.json` (Windows)
```
This file contains the API token and account ID for persistent authentication.

## API Request and Response Format
The application interacts with Cloudflare's AI API using JSON payloads:

### Request Format
```json
{
  "messages": [
    { "role": "system", "content": "You are a friendly assistant" },
    { "role": "user", "content": "Hello!" }
  ]
}
```

### Response Format
```json
{
  "result": {
    "response": "Hello! How can I help you today?"
  }
}
```

## Error Handling
If authentication fails, the application prompts for valid credentials. If an API request fails, an error message is displayed, and the user is prompted to retry.

## License
This project is licensed under the MIT License. See `LICENSE` for details.

## Contributions
Contributions are welcome! Please submit a pull request or open an issue for discussion.

## Full Codebase
This project consists of the following main components:

- `main.go`: The main entry point handling user input, API communication, and conversation flow.
- `auth/`: Manages authentication, API tokens, and configuration persistence.
- `model/`: Handles AI model selection, fetching available models, and displaying them in a structured format.

To explore the full source code, check the corresponding files in this repository.


