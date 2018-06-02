package main

import (
	"gopkg.in/yaml.v2"
		"regexp"
	"path/filepath"
	"strings"
	"os"
				"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/api/core/v1"
)


type Configuration struct {
	Config ConfigurationConfig `yaml:"config"`
	k8sService Kubernetes
}

type ConfigurationConfig struct {
	Namespaces ConfigurationNamespaces           `yaml:"namespaces"`
	ServiceAccounts ConfigurationServiceAccounts `yaml:"serviceaccounts"`
}

type ConfigurationNamespaces struct {
	Path string                  `yaml:"path"`
	AutoCleanup bool             `yaml:"cleanup"`
}

type ConfigurationServiceAccounts struct {
	SubPath string               `yaml:"subpath"`
	AutoCleanup bool             `yaml:"cleanup"`
}

type cfgNamespace struct {
	Name string
	Path string
	Labels map[string]string

	ServiceAccounts map[string]cfgServiceAccount
}

type cfgServiceAccount struct {
	Name string
	Path string
	Object runtime.Object
}

func ConfigurationCreateFromYaml(yamlString string) (c *Configuration, err error) {
	err = yaml.Unmarshal([]byte(yamlString), &c)

	// ensure abs paths (relative to config)
	c.Config.Namespaces.Path = ensureAbsConfigPath(c.Config.Namespaces.Path)

	return
}


func (c *Configuration) BuildNamespaceConfiguration() (namespaceList map[string]cfgNamespace, err error) {
	namespaceList = map[string]cfgNamespace{}

	globRegexp := regexp.MustCompile("{[^}]+}")
	glob := globRegexp.ReplaceAllString(c.Config.Namespaces.Path, "*")


	labelPath := "^" + regexp.QuoteMeta(c.Config.Namespaces.Path)
	labelPath = strings.Replace(labelPath, "\\*", "[^/]+",-1)
	labelPath = strings.Replace(labelPath, "\\?", "[^/]",-1)

	labelBuildRegexp := regexp.MustCompile("\\\\{label=([^}]+)\\\\}")
	labelMatcher := labelBuildRegexp.ReplaceAllString(labelPath, "(?P<$1>[^/]+)")
	labelExtractRegexp := regexp.MustCompile(labelMatcher)

	fsEntries, err := filepath.Glob(glob)
	if err != nil {
		panic(err.Error())
	}

	for _, fsEntry := range fsEntries {
		if IsDirectory(fsEntry) {
			namespace := cfgNamespace{}
			namespace.Name = filepath.Base(fsEntry)
			namespace.Path = fsEntry
			namespace.Labels = map[string]string{}

			match := labelExtractRegexp.FindStringSubmatch(fsEntry)

			for i, name := range labelExtractRegexp.SubexpNames() {
				if i != 0 && name != "" && match[i] != "" {
					namespace.Labels[name] = match[i]
				}
			}

			c.collectServiceAccounts(&namespace)

			namespaceList[namespace.Name] = namespace
		}
	}

	return
}

func (c *Configuration) collectServiceAccounts(namespace *cfgNamespace) () {
	serviceAccountPath := filepath.Join(namespace.Path, c.Config.ServiceAccounts.SubPath)

	fileList := recursiveFileListByPath(serviceAccountPath)

	namespace.ServiceAccounts = map[string]cfgServiceAccount{}
	for _, path := range fileList {
		item := cfgServiceAccount{}
		item.Path = path
		item.Object = c.k8sService.ParseConfig(path)
		item.Name = item.Object.(*v1.ServiceAccount).Name
		namespace.ServiceAccounts[item.Name] = item
	}
}

func recursiveFileListByPath(path string) (list []string) {
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if IsK8sConfigFile(path) {
			list = append(list, path)
		}
		return nil
	})

	return
}

func ensureAbsConfigPath(path string) (absPath string) {
	var err error
	absPath = path

	if !filepath.IsAbs(absPath) {
		absPath, err = filepath.Abs(filepath.Join(filepath.Dir(opts.Config), absPath))
		if err != nil {
			panic(err.Error())
		}
	}

	return
}
