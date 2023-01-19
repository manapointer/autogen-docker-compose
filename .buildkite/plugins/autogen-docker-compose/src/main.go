package main

import (
	"fmt"
	"log"
	"os"
)

type config struct {
	inputPath  string
	outputPath string
	executor   string
	properties map[string]string
}

func cfgFromEnv() (*config, error) {
	cfg := &config{properties: make(map[string]string)}

	require := func(p *string, name string) error {
		if *p = os.Getenv(name); *p == "" {
			return fmt.Errorf("required environment variable: %s", name)
		}
		return nil
	}

	if err := require(&cfg.inputPath, "AUTOGEN_INPUT_PATH"); err != nil {
		return nil, err
	}
	if err := require(&cfg.outputPath, "AUTOGEN_OUTPUT_PATH"); err != nil {
		return nil, err
	}
	if err := require(&cfg.executor, "AUTOGEN_EXECUTOR"); err != nil {
		return nil, err
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
