// Package argo interacts with argocd server
package argo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/session"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/larntz/artui/ui/models"
)

// Clients holds argocd clients
type Clients struct {
	APIClient      apiclient.Client
	SessionClient  session.SessionServiceClient
	ClientOptions  apiclient.ClientOptions
	SessionRequest session.SessionCreateRequest
}

// Login performs user and password authentication
func (client *Clients) Login(sr session.SessionCreateRequest) error {
	log.Printf("ArgoLogin apiclient.NewClient")
	client.SessionRequest = sr
	argoClient, err := apiclient.NewClient(&client.ClientOptions)
	if err != nil {
		fmt.Printf("Error creating argocd client: %s", err.Error())
		return err
		//log.Fatalf("apiclient.NewClient err: %s", err)
	}

	sessionCloser, sessionClient := argoClient.NewSessionClientOrDie()
	client.SessionClient = sessionClient
	defer sessionCloser.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	log.Printf("created context.WithTimeout(5s), next step create session")
	defer cancel()
	session, err := sessionClient.Create(ctx, &client.SessionRequest)
	if err != nil {
		fmt.Printf("Error creating session: %s", err.Error())
		log.Fatalf("GetApplications sessionClient.create() error: %s", err)
	}
	client.ClientOptions.AuthToken = session.Token

	log.Printf("starting NewClient with session.Token")
	client.APIClient, err = apiclient.NewClient(&client.ClientOptions)
	if err != nil {
		fmt.Printf("Error api client: %s", err.Error())
		log.Fatalf("apiclient.NewClient err: %s", err)
	}
	log.Printf("ArgoLogin complete")
	return nil
}

// WatchApplications watches an app for changes
func (client Clients) WatchApplications(ctx context.Context, wg *sync.WaitGroup, ch chan<- models.AppEvent) {
	log.Printf("starting WatchApplication")
	appCloser, appClient, err := client.APIClient.NewApplicationClient()
	if err != nil {
		log.Fatalf("apiClient.NewApplicationClient err: %s", err)
	}

	appWatcher, err := appClient.Watch(ctx, &application.ApplicationQuery{})
	if err != nil {
		log.Printf("appClientWatch(1) err: %s", err)
		fmt.Printf("cannot watch applications, exiting")
		log.Fatal("cannot watch applications, exiting")
	}
	for {
		select {

		case <-ctx.Done():
			appCloser.Close()
			wg.Done()
			log.Println("WatchApplication: ctx.Done()")
			return

		default:
			log.Printf("WatchApplication checking Recv()")
			event, err := appWatcher.Recv()

			if err != nil {
				log.Printf("WatchApplication err: %s", err)
				log.Printf("Attempting to reconnect in 5 seconds...")
				time.Sleep(5 * time.Second)
				err = client.Login(client.SessionRequest)
				if err != nil {
					log.Printf("Argocd login failed...")
				}
				appCloser, appClient, err = client.APIClient.NewApplicationClient()
				if err != nil {
					log.Printf("apiClient.NewApplicationClient err: %s", err)
				}
				appWatcher, err = appClient.Watch(ctx, &application.ApplicationQuery{})
				if err != nil {
					log.Printf("appClientWatch err: %s", err)
				}
				// TODO: issue #16
				// need to login again here, create all new clients.
				// should probably be moved to functions.
				// need to save credentials passed to the Login function or re-read somehow
			} else {
				log.Printf("WatchApplication sending %s: %s", event.Application.Name, event.Type)
				ch <- models.AppEvent{Event: *event}
			}
		}
	}
}

// ArgoWorker waits for commands from the ui
func (client Clients) ArgoWorker(ctx context.Context, wg *sync.WaitGroup, ch <-chan models.WorkerCmd) {
	log.Printf("starting ArgoWorker")
	// 	appCloser, appClient, err := client.APIClient.NewApplicationClient()
	// 	if err != nil {
	// 		log.Fatalf("apiClient.NewApplicationClient err: %s", err)
	// 	}
	// 	defer appCloser.Close()

	select {
	case <-ctx.Done():
		wg.Done()

		// wait for commands and then do some stuff here.
		/*
		   // not sure yet if I should send appClient on a
		   goroutine or have the function create a new client.
		       case refresh:
		         go RefreshApplication(appClient, app)
		       case hardRefresh:
		         go RefreshAPplicatin(appClient, app,hard=true)
		       case Sync:
		         go SyncAplication(appClient, app)
		*/
	}
}

// RefreshApplication checks for application updates, but does not sync unless autoSync is enabled on the application.
func (client Clients) RefreshApplication(ctx context.Context, appClient application.ApplicationServiceClient, app v1alpha1.Application, hardRefresh bool) {
	log.Printf("starting RefreshApplication")

	refresh := fmt.Sprintf("%t", hardRefresh)
	appQuery := application.ApplicationQuery{
		Name:    &app.Name,
		Refresh: &refresh,
	}

	returnedApp, err := appClient.Get(ctx, &appQuery)
	if err != nil {
		fmt.Printf("error getting application: %s", err.Error())
	}
	log.Printf("RefreshApplication: %v", returnedApp)
}
