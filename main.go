package main

import (
	"fmt"
	"net/http"

	"github.com/cbrgm/githubevents/githubevents"
	"github.com/google/go-github/v63/github"

	log "github.com/sirupsen/logrus"
)

func main() {

}

// Repo is added. Webhook recieved.
// Repo server is created.
// Repo server

type DaggerServer interface {
	Start() error
}

type DaggerServerImpl struct {
	PullRequestPipelines []Pipeline
	PushPipelines        []Pipeline
}

func NewDaggerServer() DaggerServer {
	return &DaggerServerImpl{
		PullRequestPipelines: []Pipeline{},
		PushPipelines:        []Pipeline{},
	}
}

func (d *DaggerServerImpl) Start() error {
	// create a new event handler
	handle := githubevents.New("foobar123!")

	handle.OnPullRequestEventAny(
		func(deliveryID string, eventName string, event *github.PullRequestEvent) error {
			log.Infof("Pull Request Event Name: %s  Action: %s", eventName, event.GetAction())
			for _, pipeline := range d.PullRequestPipelines {
				err := pipeline.Run(PipelineRunOpts{PullRequestEvent: event})
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
			for _, pipeline := range d.PushPipelines {
				err := pipeline.Run(PipelineRunOpts{PushEvent: event})
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
