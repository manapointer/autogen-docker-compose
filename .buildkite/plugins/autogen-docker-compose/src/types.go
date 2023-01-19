package main

type project struct {
	Executors map[string]executorConfig `yaml:"executors"`
}

type executorConfig struct {
	Name             string            `yaml:"-"`
	WorkingDirectory string            `yaml:"working_directory"`
	Environment      map[string]string `yaml:"environment"`
	DockerConfigs    []dockerConfig    `yaml:"docker"`
	Parameters       map[string]parameterConfig
}

type parameterConfig struct {
	Type    string `yaml:"type"`
	Default string `yaml:"default"`
}

type dockerConfig struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Command     string            `yaml:"command"`
	Environment map[string]string `yaml:"environment"`
}
