package hosts

import (
	"deploy-cli/conf"
	"deploy-cli/env"
	"deploy-cli/logger"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
)

const DefaultGroupName = "_deploy_hosts"

type Conf struct {
	Group string `yaml:"group"`
	Hosts []Host `yaml:"hosts"`
}

func LoadConfig(out any) (err error) {
	homeDir := env.HomeDir
	confFileNames := []string{
		path.Join(homeDir, "hosts.yml"),
		path.Join(homeDir, "hosts.yaml"),
	}
	err = conf.UnmarshalFiles(out, confFileNames...)
	if err != nil {
		if deployHostEnv := env.Get(env.DEPLOY_HOSTS); deployHostEnv != "" {
			// ssh://root:123456@127.0.0.1,ssh://root:123456@127.0.0.1:8022
			sshs := strings.Split(deployHostEnv, ",")
			var hosts []Host
			for _, ssh := range sshs {
				h, err := ParseSSHProtocol(ssh)
				if err != nil {
					logger.WarningF("parse env %s hosts error: %s", env.DEPLOY_HOSTS, err)
					continue
				}
				hosts = append(hosts, h)
			}
			if len(hosts) <= 0 {
				return fmt.Errorf("deploy hosts is empty")
			}
			if ptr, ok := out.(*[]Conf); ok {
				*ptr = []Conf{{Hosts: hosts, Group: DefaultGroupName}}
				return nil
			}
		}
	}
	return err
}

// ParseSSHProtocol
// ssh://username:passwd@host:port
// ssh://root:123456@127.0.0.1
// ssh://root:123456@127.0.0.1:8022
func ParseSSHProtocol(ssh string) (h Host, err error) {
	ssh = strings.Trim(ssh, "\"")
	parse, err := url.Parse(ssh)
	if err != nil {
		return
	}
	h.Host = parse.Hostname()
	h.User = parse.User.Username()
	h.Password, _ = parse.User.Password()
	portString := parse.Port()
	port, _ := strconv.Atoi(portString)
	if port == 0 {
		port = 22
	}
	h.Port = port
	return
}
