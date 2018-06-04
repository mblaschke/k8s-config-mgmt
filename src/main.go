package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"fmt"
		"k8s-config-mgmt/src/config"
	"k8s-config-mgmt/src/configmanagement"
	"k8s-config-mgmt/src/k8s"
	"k8s-config-mgmt/src/logger"
)

const (
	Author  = "webdevops.io"
	Version = "0.1.0"
)

var (
	argparser *flags.Parser
	args []string
	k8sService = k8s.Kubernetes{}
	Logger *logger.DaemonLogger
	ErrorLogger *logger.DaemonLogger
)

var opts struct {
	Config          string   `           long:"config"                  env:"CONFIG"                  description:"Path to config.yaml"`

	KubeConfig      string   `           long:"kubeconfig"              env:"KUBECONFIG"              description:"Path to .kube/config"`
	KubeContext     string   `           long:"kubecontext"             env:"KUBECONTEXT"             description:"Context of .kube/config"`
	Validate		bool	 `           long:"validate"                env:"VALIDATE"                description:"Validate only mode"`
	DryRun			bool	 `           long:"dry-run"                 env:"DRYRUN"                  description:"Dryrun"`
}

func main() {
	var err error
	var Configuration *config.Configuration
	argparser = flags.NewParser(&opts, flags.Default)
	args, err = argparser.Parse()

	initOpts()

	// Init logger
	Logger = logger.CreateDaemonLogger(0)
	ErrorLogger = logger.CreateDaemonErrorLogger(0)

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println(err)
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}

	k8sService.KubeConfig = opts.KubeConfig
	k8sService.KubeContext = opts.KubeContext
	k8sService.Logger = Logger


	Logger.Main("Configuration")
	Logger.Step("main configuration")
	if opts.Config != "" {
		Configuration, err = config.ConfigurationCreateFromFile(opts.Config)
		if err != nil {
			panic(err)
		}
		Configuration.K8sService = k8sService
	} else {
		panic("No config defined")
	}
	Logger.StepResult("done")


	Logger.Step("parsing cluster and namespace conf")
	configMgmt := configmanagement.K8sConfigManagement{}
	configMgmt.Logger = Logger
	configMgmt.Configuration = *Configuration
	configMgmt.K8sService = k8sService
	configMgmt.DryRun = opts.DryRun
	configMgmt.Validate = opts.Validate
	configMgmt.Init()
	Logger.StepResult("done")

	if !opts.Validate {
		configMgmt.Run()
	} else {
		Logger.Step("validation only run, all fine")
	}

	Logger.Main("finished")
}

func initOpts() {
	if opts.KubeConfig == "" {
		kubeconfigPath := fmt.Sprintf("%s/.kube/config", UserHomeDir())
		if _, err := os.Stat(kubeconfigPath); err == nil {
			opts.KubeConfig = kubeconfigPath
		}
	}
}
