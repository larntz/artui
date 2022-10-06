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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/larntz/artui/argo"
	"github.com/larntz/artui/ui"
	"github.com/larntz/artui/ui/models"
)

var cfgFile string

// Cluster gets set to current kubeconfig context
var Cluster string

var argocdClientOptions = apiclient.ClientOptions{
	ServerAddr:           "port-forward",
	Insecure:             false,
	PlainText:            true,
	UserAgent:            "ArTUI 0.0.1",
	PortForward:          true,
	PortForwardNamespace: "",
}

var sessionRequest = session.SessionCreateRequest{
	Username: "",
	Password: "",
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cobra",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

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
	argoClient.Login(sessionRequest)

	appEventChan := make(chan models.AppEvent, 250)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := new(sync.WaitGroup)
	wg.Add(1)

	log.Println("UI Start")
	p := tea.NewProgram(ui.InitializeModel(Cluster, appEventChan), tea.WithAltScreen(), tea.WithMouseAllMotion())

	go func() {
		go argoClient.WatchApplication(ctx, wg, appEventChan)

		for {
			select {
			case msg := <-appEventChan:
				p.Send(msg)
			}

		}
	}()

	if err := p.Start(); err != nil {
		panic(err)
	}
	log.Println("Application Exit")
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/config/artui/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// try to get current k8s context
	kubeConfigPath, err := os.UserHomeDir()
	cobra.CheckErr(err)
	kubeConfigPath += "/.kube/config"
	config := clientcmd.GetConfigFromFileOrDie(kubeConfigPath)
	fmt.Printf("Current kube context: %s\n", config.CurrentContext)
	Cluster = config.CurrentContext
	artuiConfigPrefix := "argocd.contexts." + config.CurrentContext + "."

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		config, err := os.UserConfigDir()
		cobra.CheckErr(err)
		config = filepath.Join(config, "artui")

		// Search config in home directory with name ".cobra" (without extension).
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
		if !viper.IsSet(artuiConfigPrefix + "username") {
			// change prefix to use default config options
			fmt.Printf("Changing config context to default (%s not found)\n", config.CurrentContext)
			artuiConfigPrefix = "argocd.default."
		}

		if host := viper.GetString(artuiConfigPrefix + "host"); host == "" {
			fmt.Println("Unable to get argocd host configuration.")
			// os.Exit(1)
		} else {
			argocdClientOptions.ServerAddr = host
		}

		if ns := viper.GetString(artuiConfigPrefix + "namespace"); ns == "" {
			fmt.Println("Unable to get argocd ns configuration.")
			os.Exit(1)
		} else {
			argocdClientOptions.PortForwardNamespace = ns
		}

		argocdClientOptions.PortForward = viper.GetBool(artuiConfigPrefix + "port-forward")
		argocdClientOptions.Insecure = viper.GetBool(artuiConfigPrefix + "insecure")
		argocdClientOptions.PlainText = viper.GetBool(artuiConfigPrefix + "plaintext")

		if user := viper.GetString(artuiConfigPrefix + "username"); user == "" {
			fmt.Println("Unable to get argocd user configuration.")
			os.Exit(1)
		} else {
			sessionRequest.Username = user
		}

		if password := viper.GetString("password"); password == "" {
			fmt.Println("Unable to get password. Did you set the env variable ARTUI_PASSWORD?")
			os.Exit(1)
		} else {
			sessionRequest.Password = password
		}
	} else {
    fmt.Printf("configuration failed\n")
  }
}
