package ssl

import (
	"deploy-cli/env"
	"deploy-cli/logger"
	"deploy-cli/ssls"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"os"
	"path"
	"path/filepath"
)

var flag Flags

type Flags struct {
	WebRoot  string
	Email    string
	Domain   string
	SavePath string
}

func init() {
	Cmd.Flags().StringVarP(&flag.Email, "email", "E", env.GetEmail(), "e-mail address")
	Cmd.Flags().StringVarP(&flag.WebRoot, "webroot", "R", "", "directory for project deployment")
	Cmd.Flags().StringVarP(&flag.Domain, "domain", "D", "", "Generate SSL for which domain name")
	Cmd.Flags().StringVarP(&flag.SavePath, "path", "P", "/etc/nginx/ssl/", "Where to save the certificate")
}

var Cmd = &cobra.Command{
	Use:     "ssl",
	Short:   "Generate SSL for sustainable use",
	Example: `deploy-cli ssl`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if flag.Domain == "" {
			return fmt.Errorf("domain can't be empty")
		}
		if flag.WebRoot == "" {
			return fmt.Errorf("webroot can't be empty")
		}
		_, client, err := ssls.LegoClient(flag.Email)
		if err != nil {
			return err
		}
		if err := ssls.SetWebRootChallenge(client, flag.WebRoot); err != nil {
			return err
		}
		certificate, err := ssls.ObtainCertificate(client, []string{flag.Domain})
		if err != nil {
			return err
		}
		var (
			wg errgroup.Group
		)
		KeyPath := path.Join(flag.SavePath, fmt.Sprintf("%s.key", flag.Domain))
		cerPath := path.Join(flag.SavePath, fmt.Sprintf("%s.cer", flag.Domain))
		wg.Go(func() error {
			return copyFile(certificate.Certificate, cerPath)
		})
		wg.Go(func() error {
			return copyFile(certificate.PrivateKey, KeyPath)
		})
		if err := wg.Wait(); err != nil {
			return err
		}
		logger.SuccessF("%s Certificate generated successfully", flag.Domain)
		logger.SuccessF("save key: %s", KeyPath)
		logger.SuccessF("save cer: %s", cerPath)
		return nil
	},
}

func copyFile(content []byte, path string) error {
	dir := filepath.Dir(path)
	// 如果目录不存在，创建目录
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	// 创建或打开文件
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	// 写入内容到文件
	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("fail to write to file: %w", err)
	}
	return nil
}
