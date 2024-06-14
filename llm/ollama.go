package llm

import (
	"context"
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

	models, err := ollamaClient.List(ctx)
	if err != nil {
		panic(err)
	}

	missingModels := o.checkForMissingModels(models)
	if len(missingModels) > 0 {
		panic("missing models: " + strings.Join(missingModels, ", "))
	}

	o.Client = ollamaClient
}

func (o *Ollama) checkForMissingModels(models *api.ListResponse) []string {
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

func (o *Ollama) Chat(ctx context.Context, message string) (string, error) {
	responseMessage := api.Message{}

	userMessage := api.Message{
		Role:    "user",
		Content: message,
	}

	stream := false
	ollamaChatRequest := &api.ChatRequest{
		Model:    "mistral",
		Messages: []api.Message{userMessage},
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
