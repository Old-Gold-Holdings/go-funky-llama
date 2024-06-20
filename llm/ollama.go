// Ollama is a wrapper around the Ollama API client that provides
// a simple interface for interacting with Ollama locally.
// It also ensures that the required models are available before
// starting the application.

package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
)

var RequiredModels = []string{
	"mistral",
	"codestral",
}

type Ollama struct {
	Client *api.Client
}

func (o *Ollama) New(ctx context.Context) {
	ollamaClient, err := api.ClientFromEnvironment()
	if err != nil {
		panic(err)
	}

	err = ollamaClient.Heartbeat(ctx)
	if err != nil {
		panic(err)
	}

	// Set the client only after checking for a heartbeat
	o.Client = ollamaClient

	missingModels := o.checkForMissingModels(ctx)
	if len(missingModels) > 0 {
		for _, missingModel := range missingModels {
			err := o.pullModel(ctx, missingModel)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (o *Ollama) checkForMissingModels(ctx context.Context) []string {
	models, err := o.Client.List(ctx)
	if err != nil {
		panic(err)
	}

	missingModels := []string{}
	for _, requiredModel := range RequiredModels {
		found := false
		for _, model := range models.Models {
			if strings.Split(model.Name, ":")[0] == requiredModel {
				found = true
				break
			}
		}

		if !found {
			missingModels = append(missingModels, requiredModel)
		}
	}

	return missingModels
}

func (o *Ollama) pullModel(ctx context.Context, modelName string) error {
	fmt.Println("pulling missing model: " + modelName)

	modelPullRequest := &api.PullRequest{
		Model: modelName,
	}

	err := o.Client.Pull(ctx, modelPullRequest, func(modelPullResponse api.ProgressResponse) error {
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *Ollama) Chat(ctx context.Context, chatContext []api.Message, message string) (string, error) {
	responseMessage := api.Message{}

	userMessage := api.Message{
		Role:    "user",
		Content: message,
	}

	messages := append(chatContext, userMessage)

	stream := false
	ollamaChatRequest := &api.ChatRequest{
		Model:    "mistral",
		Messages: messages,
		Stream:   &stream,
	}

	err := o.Client.Chat(ctx, ollamaChatRequest, func(ollamaChatResponse api.ChatResponse) error {
		if ollamaChatResponse.Done {
			responseMessage = ollamaChatResponse.Message
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return responseMessage.Content, nil
}
