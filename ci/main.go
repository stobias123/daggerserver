package main

import (
	"context"
	"os"

	"dagger.io/dagger"
	log "github.com/sirupsen/logrus"
	"github.com/stobias123/daggerserver"
)

func getSrc(client *dagger.Client, opts daggerserver.GetSrcOpts) *dagger.Directory {
	if os.Getenv("CI") != "" {
		return client.Git(opts.RepoUrl).Commit(opts.CommitSha).Tree().Directory(".")
	}
	// Assume we're in CI dir.
	os.Chdir("..")
	return client.Host().Directory(".")
}

func main() {
	runCtx := context.Background()
	client, err := dagger.Connect(runCtx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		log.Fatalf("Failed to connect to Dagger Engine: %v", err)
	}
	srcOpts := GetSrcOpts{}
	artifacts := client.Container().From("go:1.22").
		WithDirectory("/app", getSrc(client, srcOpts)).
		WithExec([]string{"go", "mod", "tidy"}).
		WithExec([]string{"go", "build"})
	client.Container().From("alpine:3.14").WithFile("/ci", artifacts.File("/app/ci"), dagger.ContainerWithFileOpts{Permissions: 0755}).Sync(runCtx)
}

type GithubServerPipeline struct{}

func (p *GithubServerPipeline) GetSrc(client *dagger.Client, opts daggerserver.GetSrcOpts) *dagger.Directory {
	return getSrc(client, opts)
}

func (p *GithubServerPipeline) Run(client *dagger.Client) error {
	srcOpts := daggerserver.GetSrcOpts{}
	artifacts := client.Container().From("go:1.22").
		WithDirectory("/app", getSrc(client, srcOpts)).
		WithExec([]string{"go", "mod", "tidy"}).
		WithExec([]string{"go", "build"})
	client.Container().From("alpine:3.14").WithFile("/ci", artifacts.File("/app/ci"), dagger.ContainerWithFileOpts{Permissions: 0755}).Sync(runCtx)
	return nil
}
