package main

import (
	"context"
	"fmt"

	"github.com/Old-Gold-Holdings/go-funky-llama/llm"
	"github.com/ollama/ollama/api"
)

func main() {
	ctx := context.Background()

	ollama := llm.Ollama{}
	ollama.New(ctx)
	resp, err := ollama.Chat(ctx, []api.Message{}, "What is the meaning of life? Be succinct.")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}
