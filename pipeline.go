package main

import (
	"github.com/google/go-github/v63/github"
)

type GetSrcOpts struct {
	RepoUrl   string
	CommitSha string
}

type PipelineRunOpts struct {
	PullRequestEvent *github.PullRequestEvent
	PushEvent        *github.PushEvent
}

type Pipeline interface {
	Run(opts PipelineRunOpts) error
}
