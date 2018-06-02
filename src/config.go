package main

import (
	"gopkg.in/yaml.v2"
		"regexp"
	"path/filepath"
	"strings"
	)


type Configuration struct {
	Config ConfigurationConfig `yaml:"config"`
}

type ConfigurationConfig struct {
	Namespaces ConfigurationNamespaces `yaml:"namespaces"`
}

type ConfigurationNamespaces struct {
	Path string                  `yaml:"path"`
	AutoCleanup bool             `yaml:"autocleanup"`
}

type parserNamespace struct {
	Name string
	Path string
	Labels map[string]string
}


func ConfigurationCreateFromYaml(yamlString string) (c *Configuration, err error) {
	err = yaml.Unmarshal([]byte(yamlString), &c)

	// ensure abs paths (relative to config)
	c.Config.Namespaces.Path = ensureAbsConfigPath(c.Config.Namespaces.Path)

	return
}


func (c *ConfigurationNamespaces) GetList() (namespaceList map[string]parserNamespace, err error) {
	namespaceList = map[string]parserNamespace{}

	globRegexp := regexp.MustCompile("{[^}]+}")
	glob := globRegexp.ReplaceAllString(c.Path, "*")


	labelPath := "^" + regexp.QuoteMeta(c.Path)
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
			namespace := parserNamespace{}
			namespace.Name = filepath.Base(fsEntry)
			namespace.Path = fsEntry
			namespace.Labels = map[string]string{}

			match := labelExtractRegexp.FindStringSubmatch(fsEntry)

			for i, name := range labelExtractRegexp.SubexpNames() {
				if i != 0 && name != "" && match[i] != "" {
					namespace.Labels[name] = match[i]
				}
			}

			namespaceList[namespace.Name] = namespace
		}
	}

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
