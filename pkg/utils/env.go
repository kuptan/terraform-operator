package utils

import (
	"log"
	"os"
)

type EnvConfig struct {
	DockerRepository        string
	TerraformRunnerImage    string
	TerraformRunnerImageTag string
}

var Env *EnvConfig

func getEnvOrPanic(name string) string {
	env, present := os.LookupEnv(name)

	if !present {
		log.Panicf("environment variable '%s' is required but was not found", name)
	}

	return env
}

func LoadEnv() {
	cfg := &EnvConfig{}

	cfg.DockerRepository = getEnvOrPanic("DOCKER_REPOSITRY")
	cfg.TerraformRunnerImage = getEnvOrPanic("TERRAFORM_RUNNER_IMAGE")
	cfg.TerraformRunnerImageTag = getEnvOrPanic("TERRAFORM_RUNNER_IMAGE_TAG")

	Env = cfg
}
