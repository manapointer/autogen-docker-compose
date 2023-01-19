package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type config struct {
	inputPath  string
	outputPath string
	executor   string
	parameters map[string]string
}

func cfgFromEnv() (*config, error) {
	cfg := &config{parameters: make(map[string]string)}

	require := func(p *string, name string) error {
		if *p = os.Getenv(name); *p == "" {
			return fmt.Errorf("required environment variable: %s", name)
		}
		return nil
	}

	if err := require(&cfg.inputPath, "BUILDKITE_PLUGIN_AUTOGEN_DOCKER_COMPOSE_INPUT_PATH"); err != nil {
		return nil, err
	}
	if err := require(&cfg.outputPath, "BUILDKITE_PLUGIN_AUTOGEN_DOCKER_COMPOSE_OUTPUT_PATH"); err != nil {
		return nil, err
	}
	if err := require(&cfg.executor, "BUILDKITE_PLUGIN_AUTOGEN_DOCKER_COMPOSE_EXECUTOR"); err != nil {
		return nil, err
	}

	// parse parameters from pairs
	parameters := os.Getenv("BUILDKITE_PLUGIN_AUTOGEN_DOCKER_COMPOSE_PARAMETERS")
	if parameters != "" {
		for _, pair := range strings.Split(parameters, ",") {
			parts := strings.SplitN(pair, "=", 2)
			key := strings.Trim(parts[0], " ")
			value := strings.Trim(parts[1], " ")
			cfg.parameters[key] = value
		}
	}

	return cfg, nil
}

func main() {
	cfg, err := cfgFromEnv()
	if err != nil {
		log.Fatalf("failed to initialize config: %s", err)
	}

	if err := run(cfg); err != nil {
		log.Fatalf("failed to generate docker-compose file: %s", err)
	}
}
