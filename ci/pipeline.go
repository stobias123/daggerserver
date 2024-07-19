package main

import (
	"os"

	"dagger.io/dagger"
	"github.com/stobias123/daggerserver"
)

type GithubServerPipeline struct {
	srcOpts daggerserver.GetSrcOpts
}

func NewGithubServerPipeline(srcOpts daggerserver.GetSrcOpts) *GithubServerPipeline {
	return &GithubServerPipeline{srcOpts: srcOpts}
}

func (p *GithubServerPipeline) GetSrc(client *dagger.Client) *dagger.Directory {
	if os.Getenv("CI") != "" {
		return client.Git(p.srcOpts.RepoUrl).Commit(p.srcOpts.CommitSha).Tree().Directory(".")
	}
	// Assume we're in CI dir.
	os.Chdir("..")
	return client.Host().Directory(".")
}

func (p *GithubServerPipeline) Run(client *dagger.Client) error {
	src := p.GetSrc(client)
	artifacts := client.Container().From("go:1.22").
		WithDirectory("/app", src).
		WithExec([]string{"go", "mod", "tidy"}).
		WithExec([]string{"go", "build"})
	client.Container().From("alpine:3.14").WithFile("/ci", artifacts.File("/app/ci"), dagger.ContainerWithFileOpts{Permissions: 0755}).Sync(runCtx)
	return nil
}
