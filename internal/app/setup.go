package app

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

func setupLog() *log.Logger {
	return log.Default()
	//TODO proper logging setup
}

func SetupEnvConfig() *envConfig {

	cfg := &envConfig{}
	if err := envconfig.Process("", cfg); err != nil {
		fmt.Errorf("could not parse config: %w", err)
	}
	return cfg
}
