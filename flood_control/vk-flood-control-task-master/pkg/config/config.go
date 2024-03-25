package config

import (
	"go.uber.org/config"
	"os"
)

func MergeConfigs(paths []string) *config.YAML {
	filesExpected := len(paths)
	configFiles := make([]config.YAMLOption, filesExpected)

	for i := 0; i < filesExpected; i++ {
		fd, err := os.Open(paths[i])
		if err != nil {
			panic(err)
		}
		configFiles[i] = config.Source(fd)
	}

	configFiles = append(configFiles, config.Expand(os.LookupEnv))

	provider, err := config.NewYAML(configFiles...)
	if err != nil {
		panic(err)
	}

	return provider
}

func ParseConfig(provider *config.YAML, path string, container interface{}) error {
	return provider.Get(path).Populate(container)
}
