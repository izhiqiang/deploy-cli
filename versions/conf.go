package versions

import (
	"deploy-cli/conf"
	"deploy-cli/env"
	v1 "deploy-cli/versions/v1"
	"fmt"
	"path"
)

const (
	VersionOne = iota + 1
)

var (
	deployFile string
	vConf      VersionConf
)

type VersionConf struct {
	Version int `yaml:"version"`
}

func AddLookupPath() (err error) {
	files := []string{
		path.Join(env.Get("PWD"), ".deploy.yml"),
		path.Join(env.Get("PWD"), ".deploy.yaml"),
	}
	deployFile, err = conf.LoadFiles(files...)
	if err != nil {
		return
	}
	err = conf.Unmarshal(deployFile, &vConf)
	return err
}

func (y *VersionConf) New() (Version, error) {
	switch y.Version {
	case VersionOne:
		return v1.New(deployFile)
	default:
		return nil, fmt.Errorf("version %d not support", y.Version)
	}
}
