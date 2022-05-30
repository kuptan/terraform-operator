package utils

import (
	"log"
	"os"
)

// EnvConfig holds the environment variables information
type EnvConfig struct {
	DockerRepository        string
	TerraformRunnerImage    string
	TerraformRunnerImageTag string
	KnownHostsConfigMapName string
}

// Env holds the values of the environment variables
var Env *EnvConfig

// getEnvOrPanic returns a required environment variable and panics if it does not exist
func getEnvOrPanic(name string) string {
	env, present := os.LookupEnv(name)

	if !present {
		log.Panicf("environment variable '%s' is required but was not found", name)
	}

	return env
}

// getEnvOptional returns an optional environment variable if exist
func getEnvOptional(name string) string {
	env, present := os.LookupEnv(name)

	if present {
		return env
	}

	return ""
}

// LoadEnv loads teh environment variables
func LoadEnv() {
	cfg := &EnvConfig{}

	cfg.DockerRepository = getEnvOrPanic("DOCKER_REGISTRY")
	cfg.TerraformRunnerImage = getEnvOrPanic("TERRAFORM_RUNNER_IMAGE")
	cfg.TerraformRunnerImageTag = getEnvOrPanic("TERRAFORM_RUNNER_IMAGE_TAG")
	cfg.TerraformRunnerImageTag = getEnvOrPanic("TERRAFORM_RUNNER_IMAGE_TAG")
	cfg.KnownHostsConfigMapName = getEnvOptional("KNOWN_HOSTS_CONFIGMAP_NAME")

	Env = cfg
}
