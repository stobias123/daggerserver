package main

import (
	"fmt"
	"net/http"

	"github.com/cbrgm/githubevents/githubevents"
	"github.com/google/go-github/v63/github"

	log "github.com/sirupsen/logrus"
)

func main() {
	// create a new event handler
	handle := githubevents.New("foobar123!")

	handle.OnCheckRunEventAny(
		func(deliveryID string, eventName string, event *github.CheckRunEvent) error {
			log.Infof("Check Run Event Name: %s  Requested Action: %s", eventName, event.RequestedAction)
			return nil
		},
	)

	handle.OnPullRequestEventAny(
		func(deliveryID string, eventName string, event *github.PullRequestEvent) error {
			log.Infof("Pull Request Event Name: %s  Action: %s", eventName, event.GetAction())
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
