package env

import (
	"os"
	"path"
)

var HomeDir string

const (
	DEPLOY_HOSTS = "DEPLOY_HOSTS"
	HOST_GROUP   = "HOST_GROUP"
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
