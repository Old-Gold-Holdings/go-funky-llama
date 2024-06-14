package main

import (
	"context"
	"fmt"

	"github.com/Old-Gold-Holdings/go-funky-llama/llm"
)

func main() {
	ctx := context.Background()

	ollama := llm.Ollama{}
	ollama.New(ctx)
	fmt.Println(ollama.Chat(ctx, "What is the meaning of life? Be succinct."))
}
