package argo

import (
	"context"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/session"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

// GetApplications gets argocd apps...
func GetApplications(sessionRequest session.SessionCreateRequest, clientOptions apiclient.ClientOptions) v1alpha1.ApplicationList {

	argoClient := apiclient.NewClientOrDie(&clientOptions)
	sessionCloser, sessionClient := argoClient.NewSessionClientOrDie()
	defer sessionCloser.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	session, err := sessionClient.Create(ctx, &sessionRequest)
	if err != nil {
		panic(err)
	}

	clientOptions.AuthToken = session.Token
	argoClient = apiclient.NewClientOrDie(&clientOptions)

	appCloser, appClient := argoClient.NewApplicationClientOrDie()
	defer appCloser.Close()
	apps, err := appClient.List(context.TODO(), &application.ApplicationQuery{})
	if err != nil {
		panic(err)
	}
	return *apps
}
