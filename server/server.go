package server

import (
	"fmt"
	"net/http"

	"github.com/cbrgm/githubevents/githubevents"
	"github.com/google/go-github/v63/github"
	log "github.com/sirupsen/logrus"
	"github.com/stobias123/daggerserver/pipeline"
)

// Repo is added. Webhook recieved.
// Repo server is created.
// Repo server

type DaggerServer interface {
	Start() error
}

type DaggerServerImpl struct {
	PullRequestPipelines []pipeline.Pipeline
	PushPipelines        []pipeline.Pipeline
}

func NewDaggerServer(pullRequestPipelines []pipeline.Pipeline, pushPipelines []pipeline.Pipeline) DaggerServer {
	return &DaggerServerImpl{
		PullRequestPipelines: pullRequestPipelines,
		PushPipelines:        pushPipelines,
	}
}

func (d *DaggerServerImpl) Start() error {
	// create a new event handler
	handle := githubevents.New("foobar123!")

	handle.OnPullRequestEventAny(
		func(deliveryID string, eventName string, event *github.PullRequestEvent) error {
			log.Infof("Pull Request Event Name: %s  Action: %s", eventName, event.GetAction())
			for _, pipe := range d.PullRequestPipelines {
				opts := pipeline.PipelineRunOpts{
					PullRequestEvent: event,
				}
				err := pipe.Run(opts)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)

	handle.OnPushEventAny(
		func(deliveryID string, eventName string, event *github.PushEvent) error {
			log.Infof("Push Event Name: %s  Action: %s", eventName, event.GetAction())
			for _, pipe := range d.PushPipelines {
				err := pipe.Run(pipeline.PipelineRunOpts{PushEvent: event})
				if err != nil {
					return err
				}
			}
			return nil
		},
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := handle.HandleEventRequest(r)
		if err != nil {
			fmt.Println("error")
		}
	})
	// start the server listening on port 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
