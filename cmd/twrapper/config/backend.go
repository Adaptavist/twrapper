package config

type Backend struct {
	Type  string    `yaml:"type"`
	Props ConfigMap `yaml:"props"`
}
