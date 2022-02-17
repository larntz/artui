package argo

import (
	"context"
	"fmt"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/session"
)

func main() {
	clientOptions := apiclient.ClientOptions{
		ServerAddr:           "argocd.192.168.200.240.nip.io",
		Insecure:             true,
		PlainText:            true,
		UserAgent:            "ArTUI 0.0.1",
		PortForward:          false,
		PortForwardNamespace: "argocd",
	}
	argoClient := apiclient.NewClientOrDie(&clientOptions)
	fmt.Printf("%+v\n\n", argoClient)

	sessionCloser, sessionClient := argoClient.NewSessionClientOrDie()
	defer sessionCloser.Close()
	fmt.Printf("%+v\n\n", sessionClient)

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

	fmt.Printf("token: %s\n\n", session.Token)

	clientOptions.AuthToken = session.Token
	argoClient = apiclient.NewClientOrDie(&clientOptions)

	appCloser, appClient := argoClient.NewApplicationClientOrDie()
	defer appCloser.Close()
	apps, err := appClient.List(context.TODO(), &application.ApplicationQuery{})
	if err != nil {
		panic(err)
	}

	for _, a := range apps.Items {
		fmt.Printf("%s\n\n", a.GetFinalizers())
	}
}
