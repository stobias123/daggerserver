package main

import (
	"context"
	"os"

	"dagger.io/dagger"
	"github.com/google/go-github/v63/github"
	log "github.com/sirupsen/logrus"
	"github.com/stobias123/daggerserver/pipeline"
)

type GithubServerPipeline struct{}

func NewGithubServerPipeline() *GithubServerPipeline {
	return &GithubServerPipeline{}
}

func (p *GithubServerPipeline) Run(opts pipeline.PipelineRunOpts) error {
	log.Infof("Run called for repo: %s", *opts.PushEvent.Repo.Name)
	if opts.PushEvent != nil {
		return p.runPR(opts.PushEvent)
	}
	return nil
}

func (p *GithubServerPipeline) getSrc(client *dagger.Client, repoUrl string, commitSha string) *dagger.Directory {
	if os.Getenv("CI") != "" {
		return client.Git(repoUrl).Commit(commitSha).Tree().Directory(".")
	}
	// Assume we're in CI dir.
	os.Chdir("..")
	return client.Host().Directory(".")
}

func (p *GithubServerPipeline) runPR(pushEvent *github.PushEvent) error {
	log.Infof("Running push pipeline for SHA: %s", pushEvent.HeadCommit.GetSHA())
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	src := p.getSrc(client, *pushEvent.Repo.URL, pushEvent.HeadCommit.GetSHA())
	artifacts := client.Container().From("golang:1.22").
		WithDirectory("/app", src).
		WithWorkdir("/app").
		WithExec([]string{"go", "mod", "tidy"}).
		WithExec([]string{"go", "build"})
	client.Container().From("alpine:3.14").WithFile("/ci", artifacts.File("/app/ci"), dagger.ContainerWithFileOpts{Permissions: 0755}).Sync(ctx)
	return nil
}
