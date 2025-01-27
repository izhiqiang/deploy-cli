package v1

import (
	"fmt"
	"path/filepath"
	"strings"
)

func filterEmpty[T any](slice []T, isEmpty func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if !isEmpty(item) {
			result = append(result, item)
		}
	}
	return result
}

type SafeShell struct {
}

func (s *SafeShell) CleanDstRepoFiles(conf Conf) string {
	dstRepo := conf.DstRepo
	var shell string
	switch dstRepo {
	case "/", ".":
		shell = fmt.Sprintf("ls -d /%s-* 2> /dev/null | sort -t _ -rnk2 | tail -n +%d | xargs rm -rf", conf.UUID, conf.MaxVersions+1)
	default:
		shell = fmt.Sprintf("ls -d %s/* 2> /dev/null | sort -t _ -rnk2 | tail -n +%d | xargs rm -rf", strings.TrimSuffix(dstRepo, "/"), conf.MaxVersions+1)
	}
	return shell
}

func (s *SafeShell) deployPathExistShellString(conf Conf) string {
	baseDstDir := ""
	lastSlashIndex := strings.LastIndex(conf.DstDir, "/")
	if lastSlashIndex != -1 {
		baseDstDir = conf.DstDir[:lastSlashIndex]
	}
	return s.Shells([]string{
		fmt.Sprintf("mkdir -p '%s' '%s'", conf.DstRepo, baseDstDir),
		fmt.Sprintf("[ -e '%s' ] ", conf.DstDir),
		fmt.Sprintf("{ [ ! -L '%s' ] && exit 1 || exit 0; } || exit 0 ", conf.DstDir),
	})
}

func (s *SafeShell) CdMultipleShell(path string, shells []string) string {
	shells = filterEmpty[string](shells, func(s string) bool {
		return s == ""
	})
	if len(shells) > 0 {
		return fmt.Sprintf("cd %s && %s", path, s.Shells(shells))
	}
	return ""
}
func (s *SafeShell) lnSfnShellString(repoDir string, conf Conf) string {
	return fmt.Sprintf("rm -f %s && ln -sfn %s %s", conf.DstDir, repoDir, conf.DstDir)
}

func (s *SafeShell) CdTarShellString(workerPath, tarGzFile string, conf Conf) string {
	return s.CdMultipleShell(workerPath, []string{
		fmt.Sprintf("tar -zcf %s %s", tarGzFile, s.tarOption(workerPath, conf.FilterRule)),
	})
}
func (s *SafeShell) CdUnTarShellString(tarGzFile, tarPath string, conf Conf) string {
	return s.CdMultipleShell(conf.DstRepo, []string{
		fmt.Sprintf("mkdir -p %s && tar -xzf %s -C %s", tarPath, tarGzFile, tarPath),
	})
}

func (s *SafeShell) Shells(shells []string) string {
	var validShells []string
	for _, shell := range shells {
		trimmed := strings.TrimSpace(shell)
		if trimmed != "" {
			validShells = append(validShells, trimmed)
		}
	}
	return strings.Join(validShells, " && ")
}
func (s *SafeShell) tarOption(workerPath string, rule FilterRule) string {
	var (
		exclude string
		contain string
	)
	if rule.Mode == FilterRuleIncrement {
		var includeFiles []string
		for _, x := range rule.Include {
			if x == "" {
				continue
			}
			x = strings.TrimPrefix(x, workerPath)
			x = strings.TrimPrefix(x, "/")
			includeFiles = append(includeFiles, x)
		}
		contain = strings.Join(includeFiles, " ")
	} else {
		var excludes []string
		for _, x := range rule.Exclude {
			if x == "" {
				continue
			}
			var excludePath string
			if strings.HasPrefix(x, "/") {
				excludePath = "--exclude=" + filepath.Join(x)
			} else {
				excludePath = "--exclude=" + x
			}
			excludes = append(excludes, excludePath)
		}
		exclude = strings.Join(excludes, " ")
	}
	if contain == "" {
		contain = "."
	}
	return strings.Join([]string{exclude, contain}, " ")
}
