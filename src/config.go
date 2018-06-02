package main

import (
	"gopkg.in/yaml.v2"
	"regexp"
	"path/filepath"
	"strings"
	"k8s.io/apimachinery/pkg/runtime"
	v13 "k8s.io/api/rbac/v1"
	"k8s.io/api/core/v1"
	)


type Configuration struct {
	Config ConfigurationConfig `yaml:"config"`
	k8sService Kubernetes
}

type ConfigurationConfig struct {
	Cluster ConfigurationCluster             `yaml:"cluster"`
	Namespaces ConfigurationNamespace        `yaml:"namespaces"`

	// cluster scope
	ClusterRoles ConfigurationSubItem        `yaml:"clusterroles"`
	ClusterRoleBindings ConfigurationSubItem `yaml:"clusterrolebindings""`

	// namespace scope
	ConfigMaps ConfigurationSubItem          `yaml:"configmaps"`
	ServiceAccounts ConfigurationSubItem     `yaml:"serviceaccounts"`
	Roles ConfigurationSubItem               `yaml:"roles"`
	RoleBindings ConfigurationSubItem        `yaml:"rolebindings"`
	LimitRanges ConfigurationSubItem         `yaml:"limitranges"`
}

type ConfigurationNamespace struct {
	Path []string                `yaml:"path"`
	AutoCleanup bool             `yaml:"cleanup"`
}

type ConfigurationCluster struct {
	Path []string                `yaml:"path"`
}

type ConfigurationSubItem struct {
	AutoCleanup bool             `yaml:"cleanup"`
}

type cfgCluster struct {
	ClusterRoles map[string]cfgObject
	ClusterRoleBindings map[string]cfgObject
}

type cfgNamespace struct {
	Name string
	Path string
	Labels map[string]string

	ConfigMaps map[string]cfgObject
	ServiceAccounts map[string]cfgObject
	Roles map[string]cfgObject
	RoleBindings map[string]cfgObject
	LimitRanges map[string]cfgObject
}

type cfgObject struct {
	Name string
	Path string
	Object runtime.Object
}

func ConfigurationCreateFromYaml(yamlString string) (c *Configuration, err error) {
	err = yaml.Unmarshal([]byte(yamlString), &c)

	// ensure abs paths (relative to config)
	for key, path := range c.Config.Cluster.Path {
		c.Config.Cluster.Path[key] = ensureAbsConfigPath(path)
	}

	for key, path := range c.Config.Namespaces.Path {
		c.Config.Namespaces.Path[key] = ensureAbsConfigPath(path)
	}

	return
}

func (c *Configuration) BuildClusterConfiguration() (clusterConfig cfgCluster, err error) {
	clusterConfig = cfgCluster{}

	clusterConfig.ClusterRoles = map[string]cfgObject{}
	clusterConfig.ClusterRoleBindings = map[string]cfgObject{}

	for _, configPath := range c.Config.Cluster.Path {
		fileList := recursiveFileListByPath(configPath)

		for _, path := range fileList {
			item := cfgObject{}
			item.Path = path
			item.Object = c.k8sService.ParseConfig(path)

			switch(item.Object.GetObjectKind().GroupVersionKind().Kind) {
			case "ClusterRole":
				item.Name = item.Object.(*v13.ClusterRole).Name
				clusterConfig.ClusterRoles[item.Name] = item
			case "ClusterRoleBinding":
				item.Name = item.Object.(*v13.ClusterRoleBinding).Name
				clusterConfig.ClusterRoleBindings[item.Name] = item
			default:
				panic("Not allowed object found: " + item.Object.GetObjectKind().GroupVersionKind().Kind)
			}
		}
	}

	return
}


func (c *Configuration) BuildNamespaceConfiguration() (namespaceList map[string]cfgNamespace, err error) {
	namespaceList = map[string]cfgNamespace{}

	for _, configPath := range c.Config.Namespaces.Path {
		globRegexp := regexp.MustCompile("{[^}]+}")
		glob := globRegexp.ReplaceAllString(configPath, "*")

		labelPath := "^" + regexp.QuoteMeta(configPath)
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

				c.collectConfigurationObjects(&namespace)

				namespaceList[namespace.Name] = namespace
			}
		}
	}

	return
}

func (c *Configuration) collectConfigurationObjects(namespace *cfgNamespace) () {
	fileList := recursiveFileListByPath(namespace.Path)

	namespace.ConfigMaps = map[string]cfgObject{}
	namespace.ServiceAccounts = map[string]cfgObject{}
	namespace.Roles = map[string]cfgObject{}
	namespace.RoleBindings = map[string]cfgObject{}
	namespace.LimitRanges = map[string]cfgObject{}

	for _, path := range fileList {
		item := cfgObject{}
		item.Path = path
		item.Object = c.k8sService.ParseConfig(path)

		switch(item.Object.GetObjectKind().GroupVersionKind().Kind) {
		case "ConfigMap":
			item.Name = item.Object.(*v1.ConfigMap).Name
			namespace.ConfigMaps[item.Name] = item
		case "ServiceAccount":
			item.Name = item.Object.(*v1.ServiceAccount).Name
			namespace.ServiceAccounts[item.Name] = item
		case "Role":
			item.Name = item.Object.(*v13.Role).Name
			namespace.Roles[item.Name] = item
		case "RoleBinding":
			item.Name = item.Object.(*v13.RoleBinding).Name
			namespace.RoleBindings[item.Name] = item
		case "LimitRange":
			item.Name = item.Object.(*v1.LimitRange).Name
			namespace.LimitRanges[item.Name] = item
		default:
			panic("Not allowed object found: " + item.Object.GetObjectKind().GroupVersionKind().Kind)
		}
	}
}
