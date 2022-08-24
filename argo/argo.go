package argo

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/session"
	"github.com/larntz/artui/ui/models"
)

// Clients holds argocd clients
type Clients struct {
	APIClient     apiclient.Client
	SessionClient session.SessionServiceClient
	ClientOptions apiclient.ClientOptions
}

// Login performs user and password authentication
func (client *Clients) Login(credentials session.SessionCreateRequest) {
	log.Printf("ArgoLogin apiclient.NewClient")
	argoClient, err := apiclient.NewClient(&client.ClientOptions)
	if err != nil {
		log.Fatalf("apiclient.NewClient err: %s", err)
	}

	sessionCloser, sessionClient := argoClient.NewSessionClientOrDie()
	client.SessionClient = sessionClient
	defer sessionCloser.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	log.Printf("created context.WithTimeout(5s), next step create session")
	defer cancel()
	session, err := sessionClient.Create(ctx, &credentials)
	if err != nil {
		log.Fatalf("GetApplications sessionClient.create() error: %s", err)
	}
	client.ClientOptions.AuthToken = session.Token

	log.Printf("starting NewClient with session.Token")
	client.APIClient, err = apiclient.NewClient(&client.ClientOptions)
	if err != nil {
		log.Fatalf("apiclient.NewClient err: %s", err)
	}
	log.Printf("ArgoLogin complete")
}

// WatchApplication watches an app for changes
func (client Clients) WatchApplication(ctx context.Context, wg *sync.WaitGroup, ch chan<- models.AppEvent) {
	defer wg.Done()
	log.Printf("starting WatchApplication")
	appCloser, appClient, err := client.APIClient.NewApplicationClient()
	if err != nil {
		log.Fatalf("apiClient.NewApplicationClient err: %s", err)
	}
	defer appCloser.Close()

	appWatcher, err := appClient.Watch(ctx, &application.ApplicationQuery{})
	for {
		select {

		case <-ctx.Done():
			log.Println("WatchApplication: ctx.Done()")
			return

		default:
			log.Printf("WatchApplication checking Recv()")
			event, err := appWatcher.Recv()

			if err != nil {
				log.Printf("WatchApplication err: %s", err)
			} else {
				log.Printf("WatchApplication sending %s: %s", event.Application.Name, event.Type)
				ch <- models.AppEvent{Event: *event}
			}
		}
	}
}
