package main

import (
	"gopkg.in/yaml.v2"
	"regexp"
	"path/filepath"
	"strings"
	"k8s.io/apimachinery/pkg/runtime"
	v13 "k8s.io/api/rbac/v1"
	"k8s.io/api/core/v1"
	v12 "k8s.io/api/networking/v1"
	v14 "k8s.io/api/storage/v1"
	"k8s.io/api/settings/v1alpha1"
	"k8s.io/api/policy/v1beta1"
	"os"
	"fmt"
)


type Configuration struct {
	Config ConfigurationConfig `yaml:"config"`
	k8sService Kubernetes
}

type ConfigurationConfig struct {
	Cluster ConfigurationCluster                      `yaml:"cluster"`
	Namespaces ConfigurationNamespace                 `yaml:"namespaces"`

	// cluster scope
	ClusterRoles ConfigurationSubItem                 `yaml:"clusterroles"`
	ClusterRoleBindings ConfigurationSubItem          `yaml:"clusterrolebindings"`

	// namespace scope
	ConfigMaps ConfigurationSubItem                   `yaml:"configmaps"`
	ServiceAccounts ConfigurationSubItem              `yaml:"serviceaccounts"`
	Roles ConfigurationSubItem                        `yaml:"roles"`
	RoleBindings ConfigurationSubItem                 `yaml:"rolebindings"`
	ResourceQuotas ConfigurationSubItem               `yaml:"resourcequotas"`
	NetworkPolicies ConfigurationSubItem              `yaml:"networkpolicies"`
	StorageClasses ConfigurationSubItem               `yaml:"storageclasses"`
	PodPresets ConfigurationSubItem                   `yaml:"podpresets"`
	PodSecurityPolicies ConfigurationSubItem          `yaml:"podsecuritypolicies"`
	PodDisruptionBudgets ConfigurationSubItem         `yaml:"poddisruptionbudgets"`
	LimitRanges ConfigurationSubItem                  `yaml:"limitranges"`
}

type ConfigurationNamespace struct {
	Path []string                `yaml:"path"`
	DefaultPath []string         `yaml:"defaultpath"`
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
	PodSecurityPolicies map[string]cfgObject
	StorageClasses map[string]cfgObject
}

type cfgNamespace struct {
	Name string
	Path string
	Labels map[string]string

	ConfigMaps map[string]cfgObject
	ServiceAccounts map[string]cfgObject
	Roles map[string]cfgObject
	RoleBindings map[string]cfgObject
	ResourceQuotas map[string]cfgObject
	NetworkPolicies map[string]cfgObject
	PodPresets map[string]cfgObject
	PodDisruptionBudgets map[string]cfgObject
	LimitRanges map[string]cfgObject
}

type cfgObject struct {
	Name string
	Path string
	Object runtime.Object
}

var (
	globRegexp = regexp.MustCompile("{[^}]+}")
)

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
	clusterConfig.StorageClasses = map[string]cfgObject{}
	clusterConfig.PodSecurityPolicies = map[string]cfgObject{}

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
			case "StorageClass":
				item.Name = item.Object.(*v14.StorageClass).Name
				clusterConfig.StorageClasses[item.Name] = item
			case "PodSecurityPolicy":
				item.Name = item.Object.(*v1beta1.PodSecurityPolicy).Name
				clusterConfig.StorageClasses[item.Name] = item
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
		glob := buildGlobPathFromPatternPath(configPath)

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

				// check if path is a default path
				if IsNamespaceDefaultPath(c.Config.Namespaces.DefaultPath, fsEntry) {
					continue
				}

				namespace := cfgNamespace{}
				namespace.Name = filepath.Base(fsEntry)
				namespace.Path = fsEntry

				// labels
				namespace.Labels = map[string]string{}
				match := labelExtractRegexp.FindStringSubmatch(fsEntry)
				for i, name := range labelExtractRegexp.SubexpNames() {
					if i != 0 && name != "" && match[i] != "" {
						namespace.Labels[name] = match[i]
					}
				}

				// init
				namespace.ConfigMaps = map[string]cfgObject{}
				namespace.ServiceAccounts = map[string]cfgObject{}
				namespace.Roles = map[string]cfgObject{}
				namespace.RoleBindings = map[string]cfgObject{}
				namespace.ResourceQuotas = map[string]cfgObject{}
				namespace.NetworkPolicies = map[string]cfgObject{}
				namespace.PodPresets = map[string]cfgObject{}
				namespace.PodDisruptionBudgets = map[string]cfgObject{}
				namespace.LimitRanges = map[string]cfgObject{}

				// default
				for _, defaultPath := range c.Config.Namespaces.DefaultPath {
					if defaultPath := BuildDefaultPath(defaultPath, &namespace); defaultPath != "" {
						c.collectConfigurationObjects(&namespace, defaultPath)
					}
				}

				// namespace config
				c.collectConfigurationObjects(&namespace, namespace.Path)

				namespaceList[namespace.Name] = namespace
			}
		}
	}

	return
}

func (c *Configuration) collectConfigurationObjects(namespace *cfgNamespace, path string) () {
	fileList := recursiveFileListByPath(path)

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
		case "NetworkPolicy":
			item.Name = item.Object.(*v12.NetworkPolicy).Name
			namespace.NetworkPolicies[item.Name] = item
		case "LimitRange":
			item.Name = item.Object.(*v1.LimitRange).Name
			namespace.LimitRanges[item.Name] = item
		case "PodPreset":
			item.Name = item.Object.(*v1alpha1.PodPreset).Name
			namespace.PodPresets[item.Name] = item
		case "ResourceQuota":
			item.Name = item.Object.(*v1.ResourceQuota).Name
			namespace.ResourceQuotas[item.Name] = item
		default:
			panic("Not allowed object found: " + item.Object.GetObjectKind().GroupVersionKind().Kind)
		}
	}
}

func buildGlobPathFromPatternPath(path string) (string) {
	return globRegexp.ReplaceAllString(path, "*")
}

func buildRegExpPathFromPatternPath(path string) (*regexp.Regexp) {
	regexpPath := buildGlobPathFromPatternPath(path)
	regexpPath = "^" + regexp.QuoteMeta(regexpPath)
	regexpPath = strings.Replace(regexpPath, "\\*", "[^/]+",-1)
	regexpPath = strings.Replace(regexpPath, "\\?", "[^/]",-1)

	return regexp.MustCompile(regexpPath)
}


func IsNamespaceDefaultPath(configPathList []string, path string) (bool) {
	if len(configPathList) > 0 {
		for _, configPath := range configPathList {
			regexpPath := buildRegExpPathFromPatternPath(ensureAbsConfigPath(configPath))
			if regexpPath.MatchString(path) {
				return true
			}
		}
	}

	return false
}

func BuildDefaultPath(configPath string, namespace *cfgNamespace) (string) {
	if configPath != "" {
		path := ensureAbsConfigPath(configPath)

		// replacement markers
		path = strings.Replace(path, "{namespace}", namespace.Name, -1)
		for labelName, labelValue := range namespace.Labels {
			path = strings.Replace(path, fmt.Sprintf("{label=%s}", labelName), labelValue, -1)
		}

		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			return path
		}
	}

	return ""
}
