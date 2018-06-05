package config

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
	"k8s-config-mgmt/src/k8s"
	"io/ioutil"
)

type Configuration struct {
	Path string
	Config ConfigurationConfig         `yaml:"config"`
	Management ConfigurationManagement `yaml:"management"`
	K8sService k8s.Kubernetes
}

type ConfigurationConfig struct {
	Cluster ConfigurationConfigCluster       `yaml:"cluster"`
	Namespaces ConfigurationConfigNamespace  `yaml:"namespaces"`
}

type ConfigurationManagement struct {
	Cluster ConfigurationManagementCluster         `yaml:"cluster"`
	Namespaces []ConfigurationManagementNamespace  `yaml:"namespaces"`
}

type ConfigurationManagementCluster struct {
	ClusterRolebindings ConfigurationManagementItem          `yaml:"clusterrolebindings"`
	ClusterRoles ConfigurationManagementItem                 `yaml:"clusterroles"`
	Namespaces ConfigurationManagementItem                   `yaml:"namespaces"`
	PodSecurityPolicies ConfigurationManagementItem          `yaml:"podsecuritypolicies"`
	StorageClasses ConfigurationManagementItem               `yaml:"storageclasses"`
}

type ConfigurationManagementNamespace struct {
	Name string                                              `yaml:"name"`
	ConfigMaps ConfigurationManagementItem                   `yaml:"configmaps"`
	LimitRanges ConfigurationManagementItem                  `yaml:"limitranges"`
	NetworkPolicies ConfigurationManagementItem              `yaml:"networkpolicies"`
	PodDisruptionBudgets ConfigurationManagementItem         `yaml:"poddisruptionbudgets"`
	PodPresets ConfigurationManagementItem                   `yaml:"podpresets"`
	ResourceQuotas ConfigurationManagementItem               `yaml:"resourcequotas"`
	RoleBindings ConfigurationManagementItem                 `yaml:"rolebindings"`
	Roles ConfigurationManagementItem                        `yaml:"roles"`
	ServiceAccounts ConfigurationManagementItem              `yaml:"serviceaccounts"`
}

type ConfigurationConfigNamespace struct {
	Path []string                `yaml:"path"`
	DefaultPath []string         `yaml:"defaultpath"`
	AutoCleanup bool             `yaml:"cleanup"`
}

type ConfigurationConfigCluster struct {
	Path []string                `yaml:"path"`
}

type ConfigurationManagementItem struct {
	Enabled *bool                `yaml:"enabled"`
	AutoCleanup bool             `yaml:"cleanup"`
}

type ConfigCluster struct {
	ClusterRoles map[string]ConfigObject
	ClusterRoleBindings map[string]ConfigObject
	PodSecurityPolicies map[string]ConfigObject
	StorageClasses map[string]ConfigObject
}

type ConfigNamespace struct {
	Name string
	Path string
	Labels map[string]string

	ConfigMaps map[string]ConfigObject
	ServiceAccounts map[string]ConfigObject
	Roles map[string]ConfigObject
	RoleBindings map[string]ConfigObject
	ResourceQuotas map[string]ConfigObject
	NetworkPolicies map[string]ConfigObject
	PodPresets map[string]ConfigObject
	PodDisruptionBudgets map[string]ConfigObject
	LimitRanges map[string]ConfigObject
}

type ConfigObject struct {
	Name string
	Path string
	Object runtime.Object
}

var (
	globRegexp = regexp.MustCompile("{[^}]+}")
)

func ConfigurationCreateFromFile(path string) (c *Configuration, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(string(data)), &c)
	if err != nil {
		panic(err)
	}

	c.Path = path

	// ensure abs paths (relative to config)
	for key, path := range c.Config.Cluster.Path {
		c.Config.Cluster.Path[key] = ensureAbsConfigPath(c.Path, path)
	}

	for key, path := range c.Config.Namespaces.Path {
		c.Config.Namespaces.Path[key] = ensureAbsConfigPath(c.Path, path)
	}

	return
}

func (c *Configuration) BuildClusterConfiguration() (clusterConfig ConfigCluster, err error) {
	clusterConfig = ConfigCluster{}

	clusterConfig.ClusterRoles = map[string]ConfigObject{}
	clusterConfig.ClusterRoleBindings = map[string]ConfigObject{}
	clusterConfig.StorageClasses = map[string]ConfigObject{}
	clusterConfig.PodSecurityPolicies = map[string]ConfigObject{}

	for _, configPath := range c.Config.Cluster.Path {
		fileList := recursiveFileListByPath(configPath)

		for _, path := range fileList {
			item := ConfigObject{}
			item.Path = path
			item.Object = c.K8sService.ParseConfig(path)

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


func (c *Configuration) BuildNamespaceConfiguration() (namespaceList map[string]ConfigNamespace, err error) {
	namespaceList = map[string]ConfigNamespace{}

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
				if c.IsNamespaceDefaultPath(c.Config.Namespaces.DefaultPath, fsEntry) {
					continue
				}

				namespace := ConfigNamespace{}
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
				namespace.ConfigMaps = map[string]ConfigObject{}
				namespace.ServiceAccounts = map[string]ConfigObject{}
				namespace.Roles = map[string]ConfigObject{}
				namespace.RoleBindings = map[string]ConfigObject{}
				namespace.ResourceQuotas = map[string]ConfigObject{}
				namespace.NetworkPolicies = map[string]ConfigObject{}
				namespace.PodPresets = map[string]ConfigObject{}
				namespace.PodDisruptionBudgets = map[string]ConfigObject{}
				namespace.LimitRanges = map[string]ConfigObject{}

				// default
				for _, defaultPath := range c.Config.Namespaces.DefaultPath {
					if defaultPath := c.BuildDefaultPath(defaultPath, &namespace); defaultPath != "" {
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

func (c *Configuration) collectConfigurationObjects(namespace *ConfigNamespace, path string) () {
	fileList := recursiveFileListByPath(path)

	for _, path := range fileList {
		item := ConfigObject{}
		item.Path = path
		item.Object = c.K8sService.ParseConfig(path)

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


func (c *Configuration) IsNamespaceDefaultPath(configPathList []string, path string) (bool) {
	if len(configPathList) > 0 {
		for _, configPath := range configPathList {
			regexpPath := buildRegExpPathFromPatternPath(ensureAbsConfigPath(c.Path, configPath))
			if regexpPath.MatchString(path) {
				return true
			}
		}
	}

	return false
}

func (c *Configuration) BuildDefaultPath(configPath string, namespace *ConfigNamespace) (string) {
	if configPath != "" {
		path := ensureAbsConfigPath(c.Path, configPath)

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
