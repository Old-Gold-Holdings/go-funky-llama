package main

import (
	"context"
	"fmt"

	"github.com/Old-Gold-Holdings/go-funky-llama/llm"
	"github.com/Old-Gold-Holdings/go-funky-llama/plugins"
)

func main() {
	ctx := context.Background()

	ollama := llm.Ollama{}
	ollama.New(ctx)

	git := plugins.Git{}
	git.New()

	git.AddAll()
	status, err := git.Status()
	if err != nil {
		panic(err)
	}

	fmt.Println(status)

	// Commit the changes
	err = git.Commit("fix: remove push check")
	if err != nil {
		panic(err)
	}

	// Push the changes
	err = git.PushCurrentBranch()
	if err != nil {
		panic(err)
	}
}
