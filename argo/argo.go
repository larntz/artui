package argo

import (
	"context"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/session"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

// GetApps gets apps...
func GetApplications() v1alpha1.ApplicationList {
	clientOptions := apiclient.ClientOptions{
		ServerAddr:           "argocd.192.168.200.240.nip.io",
		Insecure:             true,
		PlainText:            true,
		UserAgent:            "ArTUI 0.0.1",
		PortForward:          false,
		PortForwardNamespace: "argocd",
	}
	argoClient := apiclient.NewClientOrDie(&clientOptions)
	sessionCloser, sessionClient := argoClient.NewSessionClientOrDie()
	defer sessionCloser.Close()

	sessionRequest := session.SessionCreateRequest{
		Username: "admin",
		Password: "admin123",
	}
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
