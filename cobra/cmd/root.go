// Package cmd handles config and app setup
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/session"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/larntz/artui/argo"
	"github.com/larntz/artui/ui"
	"github.com/larntz/artui/ui/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var cfgFile string
var cfgContext string
var cfgHelp bool

// Cluster gets set to current kubeconfig context if using core, otherwise it's set to config-context
var Cluster string

var argocdClientOptions = apiclient.ClientOptions{
	ServerAddr:  "kubernetes",
	Core:        true,
	PlainText:   true,
	PortForward: false,
}

var sessionRequest = session.SessionCreateRequest{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "artui",
	Short: "TUI application for managing ArgoCD applications.",
	Long:  "TUI application for managing ArgoCD applications.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// setup logging
		f, err := tea.LogToFile("/tmp/artui-debug.log", "debug")
		if err != nil {
			fmt.Println("fatal: ", err)
			os.Exit(1)
		}
		defer f.Close()

		stdErr, err := os.OpenFile("/tmp/artui-stderr.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			fmt.Println("fatal: ", err)
			os.Exit(1)
		}
		os.Stderr = stdErr

		// start application
		log.Println("Application Start")

		argoClient := argo.Clients{
			ClientOptions: argocdClientOptions,
		}
		if err := argoClient.Login(sessionRequest); err != nil {
			log.Fatalf("argoClient.Login failed: %s", err)
		}

		appEventChan := make(chan models.AppEvent, 250)
		workerChan := make(chan models.WorkerCmd, 1)
		ctx, cancel := context.WithCancel(context.Background())
		wg := new(sync.WaitGroup)
		wg.Add(2)

		log.Println("UI Start")
		p := tea.NewProgram(ui.InitializeModel(Cluster, appEventChan, workerChan), tea.WithAltScreen(), tea.WithMouseAllMotion())

		go func() {
			go argoClient.ArgoWorker(ctx, wg, workerChan)
			go argoClient.WatchApplications(ctx, wg, appEventChan)

			for {
				msg := <-appEventChan
				p.Send(msg)
			}
		}()

		if err := p.Start(); err != nil {
			panic(err)
		}
		// wait for workers to shutdown
		log.Println("Shutting down workers...")
		fmt.Println("Shutting down workers...")
		cancel()
		wg.Wait()
		log.Println("Application Exit")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/config/artui/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&cfgContext, "config-context", "default", "config file context")
	rootCmd.PersistentFlags().BoolVarP(&cfgHelp, "help", "h", false, "get help")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// try to get current k8s context
	kubeConfigPath, err := os.UserHomeDir()
	cobra.CheckErr(err)
	kubeConfigPath += "/.kube/config"
	config := clientcmd.GetConfigFromFileOrDie(kubeConfigPath)

	var artuiConfigPrefix string

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		config, err := os.UserConfigDir()
		cobra.CheckErr(err)
		config = filepath.Join(config, "artui")

		viper.AddConfigPath(config)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yaml")
	}

	viper.SetEnvPrefix("ARTUI")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())

		// check if config has k8s context specific settings
		if cfgContext != "default" {
			artuiConfigPrefix = "argocd.contexts." + cfgContext + "."
			Cluster = cfgContext
			fmt.Printf("Using config context %s\n", artuiConfigPrefix)
		} else {
			artuiConfigPrefix = "argocd.default."
			Cluster = config.CurrentContext
			fmt.Printf("Using config context %s\n", artuiConfigPrefix)
		}

		argocdClientOptions.PortForward = viper.GetBool(artuiConfigPrefix + "port-forward")
		argocdClientOptions.Insecure = viper.GetBool(artuiConfigPrefix + "insecure")
		argocdClientOptions.PlainText = viper.GetBool(artuiConfigPrefix + "plaintext")
		argocdClientOptions.GRPCWeb = viper.GetBool(artuiConfigPrefix + "grpcweb")
		argocdClientOptions.Core = viper.GetBool(artuiConfigPrefix + "core")

		if argoNs := viper.GetString(artuiConfigPrefix + "argocdNamespace"); argoNs == "" {
			fmt.Println("Nothing set for ArgoCD Namespace.")
		} else {
			context := clientcmdapi.Context{
				Namespace: argoNs,
			}
			configOverrides := clientcmd.ConfigOverrides{
				Context: context,
			}
			argocdClientOptions.KubeOverrides = &configOverrides
		}

		if argocdClientOptions.PortForward {
			if ns := viper.GetString(artuiConfigPrefix + "namespace"); ns == "" {
				fmt.Println("Unable to get port-forward namespace configuration. Using 'argocd'.")
				argocdClientOptions.PortForwardNamespace = "argocd"
			} else {
				argocdClientOptions.PortForwardNamespace = ns
			}
		}

		if !argocdClientOptions.Core {
			// only need these if we're not using Core
			if host := viper.GetString(artuiConfigPrefix + "host"); host == "" {
				fmt.Println("Unable to get argocd host configuration.")
				os.Exit(1)
			} else {
				argocdClientOptions.ServerAddr = host
			}

			if user := viper.GetString(artuiConfigPrefix + "username"); user == "" {
				fmt.Println("Unable to get argocd user configuration.")
				os.Exit(1)
			} else {
				sessionRequest.Username = user
			}

			if password := viper.GetString("password"); password == "" {
				fmt.Println("Unable to get password. Try setting env ARTUI_PASSWORD")
				os.Exit(1)
			} else {
				sessionRequest.Password = password
			}
		}
	} else {
		// set KubeOverrides namespace to 'argocd'
		context := clientcmdapi.Context{
			Namespace: "argocd",
		}
		configOverrides := clientcmd.ConfigOverrides{
			Context: context,
		}
		argocdClientOptions.KubeOverrides = &configOverrides
		// use current kubeconfig contxt as cluster name
		Cluster = config.CurrentContext
	}
}
