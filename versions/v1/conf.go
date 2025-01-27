package v1

import (
	"deploy-cli/conf"
	"os"
	"path"
	"regexp"
)

const (
	// FilterRuleFull 全量
	FilterRuleFull = 1
	// FilterRuleIncrement 增量
	FilterRuleIncrement = 2
)

type Conf struct {
	UUID       string     `yaml:"uuid" validate:"required"`
	FilterRule FilterRule `yaml:"filter_rule" validate:"required"` //压缩文件规则

	HookPreServer  []string `yaml:"hook_pre_server"`  // 代码检出前执行
	HookPostServer []string `yaml:"hook_post_server"` //代码检出后执行

	DstRepo     string `yaml:"dst_repo" validate:"required,dir"` //项目部署路径
	DstDir      string `yaml:"dst_dir" validate:"required,dir"`  //存储路径
	MaxVersions int    `yaml:"max_versions" validate:"required"` //版本备份的数量最大值

	HookPreHost  []string `yaml:"hook_pre_host"`  //应用发布前执行
	HookPostHost []string `yaml:"hook_post_host"` //应用发布后执行
}
type FilterRule struct {
	Mode    int      `yaml:"mode" validate:"required,oneof=1 2"`
	Include []string `yaml:"include"` //包含文件
	Exclude []string `yaml:"exclude"` //排除文件
}

func ReadConf(file string) (vConf Conf, err error) {
	err = conf.Unmarshal(file, &vConf)
	if err != nil {
		return
	}
	err = conf.ValidateStruct(&vConf)
	if err != nil {
		return
	}
	return
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func BaseDir(fileName string) string {
	re := regexp.MustCompile(`\.(tar\.gz|zip|tar)$`)
	file := path.Base(fileName)
	return re.ReplaceAllString(file, "")
}
