package hosts

import (
	"errors"
)

type GHosts struct {
	groups     []string
	groupHosts map[string][]Host
	ghosts     []GHost
}

func NewGHosts() (*GHosts, error) {
	var confs []Conf
	ghs := &GHosts{
		groups:     make([]string, 0),
		groupHosts: make(map[string][]Host),
		ghosts:     make([]GHost, 0),
	}
	if err := LoadConfig(&confs); err != nil {
		return nil, err
	}
	for _, host := range confs {
		ghs.AddGroupHosts(host.Group, host.Hosts)
	}
	return ghs, nil
}

func (ghs *GHosts) AddGroupHosts(group string, hosts []Host) {
	ghs.groups = append(ghs.groups, group)
	ghs.groupHosts[group] = append(ghs.groupHosts[group], hosts...)
	for _, host := range hosts {
		ghs.ghosts = append(ghs.ghosts, GHost{Host: host, Group: group})
	}
}

func (ghs *GHosts) PingGroup(group string) error {
	errs := NewErrGHostConnects()
	for _, ghost := range ghs.ghosts {
		if group != "" && ghost.Group == group {
			continue
		}
		if err := ghost.Ping(); err != nil {
			errs.Add(ghost, err)
			continue
		}
	}
	return errs.Err()
}

func (ghs *GHosts) GetDeployGHosts(group string) ([]GHost, error) {
	var deployGhosts []GHost
	for _, ghost := range ghs.ghosts {
		if group != "" {
			if ghost.Group == group {
				deployGhosts = append(deployGhosts, ghost)
			}
		} else {
			deployGhosts = append(deployGhosts, ghost)
		}
	}
	if len(deployGhosts) <= 0 {
		return nil, errors.New("no deploy hosts found")
	}
	return deployGhosts, nil
}
