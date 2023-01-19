package main

import (
	"strings"

	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v3"
)

type project struct {
	Executors map[string]executorConfig `yaml:"executors"`
}

type executorConfig struct {
	Name             string            `yaml:"-"`
	WorkingDirectory string            `yaml:"working_directory"`
	Environment      environmentConfig `yaml:"environment"`
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
	Command     command           `yaml:"command"`
	Environment environmentConfig `yaml:"environment"`
}

type command []string

func (cmd *command) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		var command string
		if err := value.Decode(&command); err != nil {
			return err
		}
		words, err := shellwords.Parse(command)
		if err != nil {
			return err
		}
		*cmd = words
		return nil
	}
	return value.Decode((*[]string)(cmd))
}

type environmentConfig map[string]string

func (cfg *environmentConfig) UnmarshalYAML(value *yaml.Node) error {
	*cfg = environmentConfig{}
	if value.Kind == yaml.MappingNode {
		return value.Decode(map[string]string(*cfg))
	}
	var pairs []string
	if err := value.Decode(&pairs); err != nil {
		return err
	}
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		(*cfg)[parts[0]] = parts[1]
	}
	return nil
}
