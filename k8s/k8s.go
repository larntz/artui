//
// This package manages interection with the kubernetes api
//

package k8s

import (
	"context"
	"fmt"
	"html/template"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

var appResource = schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "applications"}

type app struct {
	Name        string
	Destination string
}

// Connect to the api-server
func Connect(kubeconfig string) {
	tpl, err := template.New("app").Parse("# {{ .Name }}\n\n- {{ .Destination }}\n")
	if err != nil {
		panic("template is broke")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	apps, err := clientset.Resource(appResource).Namespace("argocd").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	log.Printf("There are %d apps in the cluster\n", len(apps.Items))
	for i := 0; i < len(apps.Items); i++ {
		log.Printf("app %d, %s", i, apps.Items[i].GetName())
		d := fmt.Sprintf("%s", apps.Items[i].Object["spec"].(map[string]interface{})["destination"].(map[string]interface{})["namespace"])
		a := app{
			Name:        apps.Items[i].GetName(),
			Destination: d,
		}
		tpl.Execute(log.Writer(), a)
	}
}
