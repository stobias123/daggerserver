package main

import (
	"github.com/stobias123/daggerserver/pipeline"
	"github.com/stobias123/daggerserver/server"
)

func main() {
	pushPipe := NewGithubServerPipeline()
	server := server.NewDaggerServer(
		[]pipeline.Pipeline{},
		[]pipeline.Pipeline{pushPipe},
	)
	server.Start()
}
