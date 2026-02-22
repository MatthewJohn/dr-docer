package config

type PhpipamConfig struct {
	Url string `yaml:"url"`
}
type GitlabConfig struct {
	Url string `yaml:"url"`
}
type AnsibleConfig struct {
	Url string `yaml:"url"`
}

type Config struct {
	PhpipamConfig PhpipamConfig `yaml:"phpipam"`
	GitlabConfig  GitlabConfig  `yaml:"gitlab"`
	AnsibleConfig AnsibleConfig `yaml:"ansible"`
}
