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
	Logger *logger.DaemonLogger
	ErrorLogger *logger.DaemonLogger
)

var opts struct {
	Config          string   `           long:"config"                  env:"CONFIG"                  description:"Path to config.yaml"`

	KubeConfig      string   `           long:"kubeconfig"              env:"KUBECONFIG"              description:"Path to .kube/config"`
	KubeContext     string   `           long:"kubecontext"             env:"KUBECONTEXT"             description:"Context of .kube/config"`
	Validate		bool	 `           long:"validate"                env:"VALIDATE"                description:"Validate only mode"`
	DryRun			bool	 `           long:"dry-run"                 env:"DRYRUN"                  description:"Dryrun"`
	Force			bool	 `           long:"force"                   env:"FORCE"                   description:"Force (delete/recreate on error)"`
}

func main() {
	initArgparser()

	// Init logger
	Logger = logger.CreateDaemonLogger(0)
	ErrorLogger = logger.CreateDaemonErrorLogger(0)


	Logger.Main("Init")

	Logger.Step("k8s connection")
	k8sService := initK8sService()
	Logger.StepResult("done")

	Logger.Step("main configuration")
	Configuration := initConfiguration(k8sService)
	Logger.StepResult("done")


	Logger.Step("parsing cluster and namespace conf")
	configMgmt := initConfigManagement(Configuration, k8sService)
	Logger.StepResult("done")

	if !opts.Validate {
		configMgmt.Run()
	} else {
		Logger.Step("validation only run, all fine")
	}

	Logger.Main("finished")
}

func initArgparser() {
	argparser = flags.NewParser(&opts, flags.Default)
	_, err := argparser.Parse()

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

	if opts.KubeConfig == "" {
		kubeconfigPath := fmt.Sprintf("%s/.kube/config", UserHomeDir())
		if _, err := os.Stat(kubeconfigPath); err == nil {
			opts.KubeConfig = kubeconfigPath
		}
	}
}

func initK8sService() (*k8s.Kubernetes) {
	service := k8s.Kubernetes{}
	service.KubeConfig = opts.KubeConfig
	service.KubeContext = opts.KubeContext
	service.Logger = Logger

	return &service
}

func initConfiguration(k8sService *k8s.Kubernetes) (*config.Configuration) {
	var (
		err error
		conf *config.Configuration
	)

	if opts.Config != "" {
		conf, err = config.ConfigurationCreateFromFile(opts.Config)
		if err != nil {
			panic(err)
		}
	} else {
		panic("No config defined")
	}

	conf.K8sService = *k8sService

	return conf
}

func initConfigManagement(config *config.Configuration, k8sService *k8s.Kubernetes) (*configmanagement.K8sConfigManagement) {
	configMgmt := configmanagement.K8sConfigManagement{}
	configMgmt.Logger = Logger
	configMgmt.GlobalConfiguration = *config
	configMgmt.K8sService = k8sService
	configMgmt.DryRun = opts.DryRun
	configMgmt.Validate = opts.Validate
	configMgmt.Force = opts.Force
	configMgmt.Init()

	return &configMgmt
}
