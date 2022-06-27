package argo

import (
	"context"
	"log"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/session"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

// TODO the idea is to save the client in the application state to prevent having to create a client on every refresh.
// Eventually I'd like to have it referesh automatically every x seconds, and be able to request a refresh.
// Right now we are seeing a port forwarding error message on some refreshses. I am hoping that reusing the client will prevent that.

// GetArgoClient returns an argocd client
// func GetArgoClient(sessionRequest session.SessionCreateRequest, clientOptions apiclient.ClientOptions) apiclient.Client {
//
// 	argoClient := apiclient.NewClientOrDie(&clientOptions)
// 	sessionCloser, sessionClient := argoClient.NewSessionClientOrDie()
// 	defer sessionCloser.Close()
//
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	session, err := sessionClient.Create(ctx, &sessionRequest)
// 	if err != nil {
// 		log.Fatalf("GetArgoClient sessionClient.create() error: %s", err)
// 	}
// 	clientOptions.AuthToken = session.Token
// 	return apiclient.NewClientOrDie(&clientOptions)
// }

// GetApplications gets argocd apps...
func GetApplications(sessionRequest session.SessionCreateRequest, clientOptions apiclient.ClientOptions) v1alpha1.ApplicationList {

	log.Printf("starting apiclient.NewClient")
	argoClient, err := apiclient.NewClient(&clientOptions)
	if err != nil {
		log.Fatalf("apiclient.NewClient err: %s", err)
	}
	sessionCloser, sessionClient := argoClient.NewSessionClientOrDie()
	defer sessionCloser.Close()

	// ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*2))
	log.Printf("Created context.WithTimeout, next setp create session")
	defer cancel()
	session, err := sessionClient.Create(ctx, &sessionRequest)
	if err != nil {
		log.Fatalf("GetApplications sessionClient.create() error: %s", err)
	}

	log.Printf("starting NewClient")
	clientOptions.AuthToken = session.Token
	argoClient, err = apiclient.NewClient(&clientOptions)
	if err != nil {
		log.Fatalf("apiclient.NewClient err: %s", err)
	}

	log.Printf("starting NewApplication")
	appCloser, appClient, err := argoClient.NewApplicationClient()
	if err != nil {
		log.Fatalf("apiClient.NewApplicationClient err: %s", err)
	}
	defer appCloser.Close()

	log.Printf("starting appClient.List")
	apps, err := appClient.List(context.TODO(), &application.ApplicationQuery{})
	if err != nil {
		log.Fatalf("GetApplications apiClient.List() error: %s", err)
	}
	log.Printf("returning from argo.GetApplications")
	return *apps
}

// GetApplications2 gets argocd apps...
// func GetApplications2(argoClient apiclient.Client) v1alpha1.ApplicationList {
// 	appCloser, appClient := argoClient.NewApplicationClientOrDie()
// 	defer appCloser.Close()
// 	apps, err := appClient.List(context.TODO(), &application.ApplicationQuery{})
// 	if err != nil {
// 		log.Fatalf("GetApplications2 apiClient.List() error: %s", err)
// 	}
// 	return *apps
// }
