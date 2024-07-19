package main

import (
	"context"

	"dagger.io/dagger"
)

type GetSrcOpts struct {
	RepoUrl   string
	CommitSha string
}

type Pipeline interface {
	GetSrc(ctx context.Context, client *dagger.Client, opts GetSrcOpts) *dagger.Directory
	Run(ctx context.Context, client *dagger.Client) error
}
