## Host configuration`~/.deploy/hosts.yml `

~~~
- group: web
  hosts:
    - host: 127.0.0.1
      port: 22
      user: root
      password: '123456'
~~~

> Define environment `DEPLOY_HOSTS` variables to automatically write host configuration
>
> Format：ssh://username:passwd@host:port
>
> for example: `ssh://root:123456@127.0.0.1:8022,ssh://root:123456@127.0.0.1:22`

project configuration `project/.deploy.yml`

~~~
# Version for executing the release process, default is 1, waiting for further support
version: 1
# Unique representation of the project
uuid: golang
# Compressed file rules
filter_rule:
  mode: 2
  #  Include files
  include:
    - 'deploy'
#During project execution
hook_post_server:
  - "export PATH=$PATH:/opt/homebrew/opt/go/bin"
  - "go env -w GO111MODULE=on"
  - "go env -w GOPROXY=https://goproxy.cn,direct"
  - "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -o deploy main.go"
# Project storage path
dst_repo: /data/deploy/golang
# Maximum number of version backups
max_versions: 5
#Project deployment path
dst_dir: /data/wwwroot/golang
# Execute before application release
hook_pre_host:
  - "chmod +x deploy"
#Execute after application release
hook_post_host:
  - ""
~~~

## Command Parameters

~~~
% ./deploy-cli --help
From development to production, a robust and easy-to-use developer tool
that makes adoption quick and easy for building and deploying  native applications.

Usage:
  deploy-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        Get the dir_name that can be rolled back from the server
  ping        Check server status
  rollback    rollback the code to the specified version
  run         Publish configuration information to the server through `.deploy.yml`

Flags:
  -h, --help                 help for deploy-cli
  -G, --host_group string    Operation host group name
  -H, --hosts string         SSH link hosts, for example: ssh://username:passwd@host:port
  -v, --version              version for deploy-cli
  -W, --working_dir string   Changes the working directory for the console

Use "deploy-cli [command] --help" for more information about a command.
~~~
