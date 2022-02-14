//
// This package manages interection with the kubernetes api
//

package k8s

import (
	"context"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/larntz/artui/models"
)

var appResource = schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "applications"}

// Connect to the api-server
func Connect(kubeconfig string) []models.Application {
	// tpl, err := template.New("app").Parse("# {{ .Name }}\n\n- {{ .Status }}\n")
	// if err != nil {
	// 	panic("template is broke")
	// }
	apps := make([]models.Application, 0, 100)

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
	appResources, err := clientset.Resource(appResource).Namespace("argocd").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	log.Printf("There are %d apps in the cluster\n", len(appResources.Items))
	for i := 0; i < len(appResources.Items); i++ {
		y, err := yaml.Marshal(appResources.Items[i].Object["status"].(map[string]interface{})["health"].(map[string]interface{})["status"])
		if err != nil {
			panic("yaml unmarshal failure")
		}
		_, err = yaml.Marshal(appResources.Items[i].Object["status"].(map[string]interface{})["sync"].(map[string]interface{})["status"])
		if err != nil {
			panic("yaml unmarshal failure")
		}

		log.Printf("%s", appResources.Items[i].GetName())
		apps = append(apps, models.Application{
			Name:       appResources.Items[i].GetName(),
			Status:     fmt.Sprintf("%s", string(y)),
			LongStatus: string(y),
		})

		// log.Printf("app %d, %s", i, apps.Items[i].GetName())
		// d := fmt.Sprintf("status: %s", apps.Items[i].Object["status"])
		// a := models.Application{
		// 	Name:   apps.Items[i].GetName(),
		// 	Status: d,
		// }
		// tpl.Execute(log.Writer(), a)

	}
	return apps
}
