package config

type PhpipamConfig struct {
	Url string `mapstructure:"url"`
}
type GitlabConfig struct {
	Url string `mapstructure:"url"`
}
type AnsibleConfig struct {
	Url string `mapstructure:"url"`
}

type Config struct {
	PhpipamConfig PhpipamConfig `mapstructure:"phpipam"`
	GitlabConfig  GitlabConfig  `mapstructure:"gitlab"`
	AnsibleConfig AnsibleConfig `mapstructure:"ansible"`

	OutputDirectory      string            `mapstructure:"output_directory"`
	DataDirectory        string            `mapstructure:"data_directory"`
	DataDirectoryMappers map[string]string `mapstructure:"data_directory_mappers"`
	DataFileExtensions   []string          `mapstructure:"data_file_extensions"`
}
