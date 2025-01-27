package versions

import (
	"deploy-cli/env"
	"deploy-cli/hosts"
)

type Version interface {
	List(workerPath string, deployHosts []hosts.GHost) error
	Rollback(deployHosts []hosts.GHost, dirName string) error
	Run(workerPath string, deployHosts []hosts.GHost) error
}

type Versions struct {
	v Version
}

func NewVersions() *Versions {
	return &Versions{}
}
func NewVersionsWithVersion() (*Versions, error) {
	versions := &Versions{}
	v, err := vConf.New()
	if err != nil {
		return nil, err
	}
	versions.v = v
	return versions, nil
}

func (v *Versions) GetHostGroupName() string {
	return env.Get(env.HOST_GROUP)
}
func (v *Versions) GetPWD() string {
	return env.Get("PWD")
}

func (v *Versions) Ping() error {
	ghost, err := hosts.NewGHosts()
	if err != nil {
		return err
	}
	err = ghost.PingGroup(v.GetHostGroupName())
	return err
}
func (v *Versions) getDeployHosts() ([]hosts.GHost, error) {
	ghost, err := hosts.NewGHosts()
	if err != nil {
		return nil, err
	}
	return ghost.GetDeployGHosts(v.GetHostGroupName())
}

func (v *Versions) List() error {
	deployHosts, err := v.getDeployHosts()
	if err != nil {
		return err
	}
	return v.v.List(v.GetPWD(), deployHosts)
}
func (v *Versions) Rollback(dirName string) error {
	deployHosts, err := v.getDeployHosts()
	if err != nil {
		return err
	}
	return v.v.Rollback(deployHosts, dirName)
}

func (v *Versions) Run() error {
	deployHosts, err := v.getDeployHosts()
	if err != nil {
		return err
	}
	return v.v.Run(v.GetPWD(), deployHosts)
}
