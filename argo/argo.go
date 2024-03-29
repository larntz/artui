// Package argo interacts with argocd server
package argo

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/session"

	"github.com/larntz/artui/argo/headless"
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

	if !client.ClientOptions.Core {
		// credentialed login
		fmt.Printf("Starting credentialed login...\n")
		client.SessionRequest = sr
		argoClient, err := apiclient.NewClient(&client.ClientOptions)
		if err != nil {
			fmt.Printf("Failed to create api client: %s\n", err)
			fmt.Printf("clientOptions: %+v\n\n", client.ClientOptions)
			os.Exit(1)
		}

		sessionCloser, sessionClient, err := argoClient.NewSessionClient()
		if err != nil {
			fmt.Printf("Failed to create session client: %s", err)
			fmt.Printf("clientOptions: %+v\n\n", client.ClientOptions)
			os.Exit(1)
		}
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
		fmt.Println("SessionCreate succesful")
		client.ClientOptions.AuthToken = session.Token

		log.Printf("starting NewClient with session.Token")
		client.APIClient, err = apiclient.NewClient(&client.ClientOptions)
		if err != nil {
			fmt.Printf("Error api client: %s", err.Error())
			log.Fatalf("apiclient.NewClient err: %s", err)
		}
		log.Printf("ArgoLogin/Server complete")
	} else {
		// core login
		client.APIClient = headless.NewClientOrDie(&client.ClientOptions)
		log.Printf("ArgoLogin/Core complete (headless url: %s", client.APIClient.ClientOptions().ServerAddr)
	}
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
				if ctx.Err() != nil {
					wg.Done()
					log.Println("WatchApplication: ctx.Done()")
					return
				}
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
	appCloser, appClient, err := client.APIClient.NewApplicationClient()
	if err != nil {
		log.Fatalf("apiClient.NewApplicationClient err: %s", err)
	}
	defer appCloser.Close()

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			appCloser.Close()
			log.Println("ArgoWorker: ctx.Done()")
			return
		case command := <-ch:
			switch {
			case command.Cmd == models.Refresh:
				log.Printf("ArgoWorker received Refresh command for app: %s", command.App.Name)
				refresh := "false"
				_, err := appClient.Get(context.TODO(), &application.ApplicationQuery{
					Name:    &command.App.Name,
					Refresh: &refresh,
				})
				if err != nil {
					log.Fatalf("appClient.Get failed: %s", err)
				}

			case command.Cmd == models.HardRefresh:
				log.Printf("ArgoWorker received HardRefresh command for app: %s", command.App.Name)
				refresh := "true"
				_, err := appClient.Get(context.TODO(), &application.ApplicationQuery{
					Name:    &command.App.Name,
					Refresh: &refresh,
				})
				if err != nil {
					log.Fatalf("appClient.Get failed: %s", err)
				}

			case command.Cmd == models.Sync:
				prune := true
				log.Printf("ArgoWorker received Sync command for app: %s", command.App.Name)
				_, err := appClient.Sync(context.TODO(), &application.ApplicationSyncRequest{
					Name:  &command.App.Name,
					Prune: &prune,
				})
				if err != nil {
					log.Fatalf("appClient.Sync failed: %s", err)
				}

			default:
				log.Printf("ArgoWorker received unknown command for app: %s", command.App.Name)
			}
		}
	}
}
