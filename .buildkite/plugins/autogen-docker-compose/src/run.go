package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	composeTypes "github.com/compose-spec/compose-go/types"
	"gopkg.in/yaml.v3"
)

func run(cfg *config) error {
	b, err := os.ReadFile(cfg.inputPath)
	if err != nil {
		return err
	}

	var project project
	if err := yaml.Unmarshal(b, &project); err != nil {
		return err
	}
	// only translate the specified executor, which must be of the `docker` type
	executor, ok := project.Executors[cfg.executor]
	executor.Name = cfg.executor
	if !ok {
		return fmt.Errorf("executor does not exist: %s", cfg.executor)
	}
	if len(executor.DockerConfigs) == 0 {
		return fmt.Errorf("executor is not of type `docker`: %s", cfg.executor)
	}

	var composeProject composeTypes.Project

	for i, dockerConfig := range executor.DockerConfigs {
		serviceConfig := asComposeServiceConfig(&executor, &dockerConfig, i)
		composeProject.Services = append(composeProject.Services, serviceConfig)
		if i > 0 {
			if composeProject.Services[0].DependsOn == nil {
				composeProject.Services[0].DependsOn = composeTypes.DependsOnConfig{}
			}
			composeProject.Services[0].DependsOn[serviceConfig.Name] = composeTypes.ServiceDependency{
				Condition: composeTypes.ServiceConditionStarted,
			}
		}
	}

	log.Printf("writing compose project to %s\n", cfg.outputPath)

	var buf bytes.Buffer

	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	enc.Encode(composeProject)
	if err := enc.Close(); err != nil {
		return err
	}

	return os.WriteFile(cfg.outputPath, buf.Bytes(), 0644)
}

func asComposeServiceConfig(executorConfig *executorConfig, dockerConfig *dockerConfig, index int) composeTypes.ServiceConfig {
	// merge environments, with container environment taking precedence
	environment := map[string]string{}
	for key, value := range executorConfig.Environment {
		environment[key] = value
	}
	for key, value := range dockerConfig.Environment {
		environment[key] = value
	}
	var pairs []string
	for key, value := range environment {
		pairs = append(pairs, key+"="+value)
	}

	// the first item in the docker configuration is the main container
	name := dockerConfig.Name
	if index == 0 {
		name = executorConfig.Name
	} else if name == "" {
		name = fmt.Sprintf("%s_%d", executorConfig.Name, index)
	}

	return composeTypes.ServiceConfig{
		Name:        name,
		Image:       dockerConfig.Image,
		Command:     composeTypes.ShellCommand(dockerConfig.Command),
		Environment: composeTypes.NewMappingWithEquals(pairs),
	}
}
