package v1

import (
	"deploy-cli/cmd"
	"deploy-cli/hosts"
	"deploy-cli/logger"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type Version struct {
	conf      Conf
	safeShell safeShell
}

func New(file string) (*Version, error) {
	cfg, err := ReadConf(file)
	if err != nil {
		return nil, err
	}
	return &Version{conf: cfg, safeShell: safeShell{}}, nil
}

func (v *Version) List(workerPath string, deployHosts []hosts.GHost) error {
	if !DirExists(workerPath) {
		return fmt.Errorf("worker path %s does not exist", workerPath)
	}
	for _, ghost := range deployHosts {
		logHostName := ghost.Base64SSHProtocol()
		defer func() {
			_ = ghost.Close()
		}()
		if err := ghost.Ping(); err != nil {
			logger.Warning(err)
			continue
		}
		sftp, err := ghost.NewSFTP()
		if err != nil {
			logger.Error(err)
			continue
		}
		dir, err := sftp.ReadDir(v.conf.DstRepo)
		if err != nil {
			logger.Error(err)
			continue
		}
		var dirs []string
		for _, d := range dir {
			dirs = append(dirs, d.Name())
		}
		logger.SuccessF("【%s】 exists: %s can be rolled back", logHostName, strings.Join(dirs, ","))
	}
	return nil
}

func (v *Version) Rollback(deployHosts []hosts.GHost, dirName string) error {
	if dirName == "" {
		return fmt.Errorf("missing or missing name")
	}
	sshRepoDir := path.Join(v.conf.DstRepo, dirName)
	for _, ghost := range deployHosts {
		logHostName := ghost.Base64SSHProtocol()
		defer func() {
			_ = ghost.Close()
		}()
		if err := ghost.Ping(); err != nil {
			logger.Warning(err)
			continue
		}
		sftp, err := ghost.NewSFTP()
		if err != nil {
			logger.WarningF("【%s】 SFTP initialization failed err: %s", logHostName, err.Error())
			continue
		}
		_, err = sftp.Stat(sshRepoDir)
		if err != nil {
			logger.WarningF("【%s】 Searching for directory:`%s` Error: %s", logHostName, sshRepoDir, err.Error())
			continue
		}
		releaseLink := v.safeShell.RmLn(sshRepoDir, v.conf)
		_, err = ghost.Command(releaseLink)
		if err != nil {
			logger.Warning(err)
			continue
		}
		hookPostHost := v.safeShell.CdShells(v.conf.DstDir, v.conf.HookPostHost)
		if hookPostHost != "" {
			logger.InfoF("【%s】 Execute `%s` after the application is released", logHostName, hookPostHost)
			if _, err := ghost.Command(hookPostHost); err != nil {
				logger.Warning(err)
				continue
			}
		}
		logger.SuccessF("【%s】 Rollback %s Successfully", logHostName, sshRepoDir)
	}
	logger.Success("All rollbacks were successful")
	return nil
}

func (v *Version) Run(workerPath string, deployHosts []hosts.GHost) error {
	var err error
	logger.InfoF("Enter the working directory %s", workerPath)
	if !DirExists(workerPath) {
		return fmt.Errorf("worker path %s does not exist", workerPath)
	}
	hookPostServer := v.safeShell.CdShells(workerPath, v.conf.HookPostServer)
	if hookPostServer != "" {
		logger.InfoF("After checking out the code, execute `%s`", hookPostServer)
		if _, err = cmd.RunCMD(hookPostServer); err != nil {
			return err
		}
	}
	//本地tar文件
	tarGzFile := fmt.Sprintf("%s-%s.tar.gz", v.conf.UUID, time.Now().Format("20060102150405"))
	//开始进行打包
	tarCommand := v.safeShell.CdTar(workerPath, tarGzFile, v.conf)
	logger.InfoF("Perform packaging command `%s` ...", tarCommand)
	if _, err = cmd.RunCMD(tarCommand); err != nil {
		return err
	}
	logger.SuccessF("The packaging has been completed  %s", tarGzFile)
	//打扫战场
	defer func() {
		_ = os.Remove(tarGzFile)
	}()
	tarPath := BaseDir(tarGzFile)
	//上传到服务器备份目录
	sshRepoDir := path.Join(v.conf.DstRepo, tarPath)
	//上传到服务器到压缩包
	sftpTarFile := path.Join(v.conf.DstRepo, tarGzFile)
	logger.InfoF("Uploaded %s to %s", tarGzFile, sftpTarFile)
	for _, ghost := range deployHosts {
		logHostName := ghost.Base64SSHProtocol()
		defer func() {
			sftp, err := ghost.NewSFTP()
			if err == nil {
				_ = sftp.Remove(sftpTarFile)
			}
			_ = ghost.Close()
		}()
		//创建备份目录和创建发布目录 并校验 发布目录是否为软连接
		_, err := ghost.Command(v.safeShell.CheckDstDir(v.conf))
		if err != nil {
			if erc, ok := err.(*hosts.ErrRunCommand); ok {
				if erc.Status == 1 {
					err = fmt.Errorf("【%s】it is detected that the publishing directory %s of this host already exists. For data security, please back up and delete the directory",
						logHostName,
						v.conf.DstDir,
					)
				}
			}
			return err
		}
		//清理过期文件
		cleanPathFileCmd := v.safeShell.CleanDstRepoFile(v.conf)
		logger.WarningF("【%s】 General instructions and instructions：`%s` ", logHostName, cleanPathFileCmd)
		_, err = ghost.Command(cleanPathFileCmd)
		if err != nil {
			logger.Warning(err)
			continue
		}
		fileByte, err := os.ReadFile(tarGzFile)
		if err != nil {
			logger.Error(err)
			return err
		}
		//上传压缩包
		if err := ghost.Upload(sftpTarFile, fileByte); err != nil {
			logger.Warning(err)
			continue
		}
		logger.SuccessF("【%s】 uploaded `%s` to `%s`", logHostName, tarGzFile, sftpTarFile)
		unTar := v.safeShell.CdUnTar(tarGzFile, sshRepoDir, v.conf)
		_, err = ghost.Command(unTar)
		if err != nil {
			logger.Warning(err)
			continue
		}
		logger.SuccessF("【%s】 unpacking `%s` to `%s` ", logHostName, sftpTarFile, sshRepoDir)
		hookPreHost := v.safeShell.CdShells(sshRepoDir, v.conf.HookPreHost)
		if hookPreHost != "" {
			logger.InfoF("【%s】 Execute before app release `%s` ", logHostName, hookPreHost)
			if _, err := ghost.Command(hookPreHost); err != nil {
				logger.Warning(err)
				continue
			}
		}
		//建立发布软连接
		releaseLink := v.safeShell.RmLn(sshRepoDir, v.conf)
		_, err = ghost.Command(releaseLink)
		if err != nil {
			logger.Warning(err)
			continue
		}
		logger.SuccessF("【%s】 `%s`  Connect the software to `%s`", logHostName, sshRepoDir, v.conf.DstDir)
		//执行发布成功之后的命令
		hookPostHost := v.safeShell.CdShells(v.conf.DstDir, v.conf.HookPostHost)
		if hookPostHost != "" {
			logger.InfoF("【%s】 Execute `%s` after the application is released", logHostName, hookPostHost)
			if _, err := ghost.Command(hookPostHost); err != nil {
				logger.Warning(err)
				continue
			}
		}
		logger.SuccessF("【%s】 published successfully", logHostName)
	}
	logger.Success("All published successfully")
	return nil
}
