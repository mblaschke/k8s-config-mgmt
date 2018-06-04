package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"fmt"
	"io/ioutil"
	)

const (
	Author  = "webdevops.io"
	Version = "0.1.0"
)

var (
	argparser *flags.Parser
	args []string
	k8sService = Kubernetes{}
	Logger *DaemonLogger
	ErrorLogger *DaemonLogger
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
	var Configuration *Configuration
	argparser = flags.NewParser(&opts, flags.Default)
	args, err = argparser.Parse()

	initOpts()

	// Init logger
	Logger = CreateDaemonLogger(0)
	ErrorLogger = CreateDaemonErrorLogger(0)

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
		data, err := ioutil.ReadFile(opts.Config)
		if err != nil {
			panic(err)
		}

		Configuration, err = ConfigurationCreateFromYaml(string(data))
		if err != nil {
			panic(err)
		}
		Configuration.k8sService = k8sService
	} else {
		panic("No config defined")
	}
	Logger.StepResult("done")


	Logger.Step("parsing cluster and namespace conf")
	configMgmt := K8sConfigManagement{}
	configMgmt.Logger = Logger
	configMgmt.Configuration = *Configuration
	configMgmt.K8sService = k8sService
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
