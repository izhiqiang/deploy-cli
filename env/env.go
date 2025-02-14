package env

import (
	"os"
	"path"
	"runtime"
)

var HomeDir string

const (
	DEPLOY_HOSTS = "DEPLOY_HOSTS"
	HOST_GROUP   = "HOST_GROUP"
	Email        = "email"
	Version      = "v0.0.x-dev"
)

func init() {
	HomeDir = path.Join(Get("HOME"), ".deploy")
}

// Get get environment variable value
func Get(key string) string {
	return os.Getenv(key)
}

// Set set environment variable value
func Set(key string, value string) {
	_ = os.Setenv(key, value)
}

// All get all environment variables
func All() []string {
	return os.Environ()
}

func GetEmail() string {
	envMail := Get(Email)
	if envMail == "" {
		envMail = "deploy-cli@" + runtime.GOOS + ".com"
		Set(Email, envMail)
	}
	return envMail
}
