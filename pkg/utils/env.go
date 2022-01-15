package utils

import (
	"log"
	"os"
)

type EnvConfig struct {
	DockerRepository        string
	TerraformRunnerImage    string
	TerraformRunnerImageTag string
	KnownHostsConfigMapName string
}

var Env *EnvConfig

func getEnvOrPanic(name string) string {
	env, present := os.LookupEnv(name)

	if !present {
		log.Panicf("environment variable '%s' is required but was not found", name)
	}

	return env
}

func getEnvOptional(name string) string {
	env, present := os.LookupEnv(name)

	if present {
		return env
	}

	return ""
}

func LoadEnv() {
	cfg := &EnvConfig{}

	cfg.DockerRepository = getEnvOrPanic("DOCKER_REGISTRY")
	cfg.TerraformRunnerImage = getEnvOrPanic("TERRAFORM_RUNNER_IMAGE")
	cfg.TerraformRunnerImageTag = getEnvOrPanic("TERRAFORM_RUNNER_IMAGE_TAG")
	cfg.TerraformRunnerImageTag = getEnvOrPanic("TERRAFORM_RUNNER_IMAGE_TAG")
	cfg.KnownHostsConfigMapName = getEnvOptional("KNOWN_HOSTS_CONFIGMAP_NAME")

	Env = cfg
}
