package hosts

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
	"time"
)

type GHost struct {
	Host  Host
	Group string

	sshClient  *ssh.Client  `yaml:"-"`
	sftpClient *sftp.Client `yaml:"-"`
}

func (h *GHost) Ping() error {
	_, err := h.NewSSHClient()
	if err != nil {
		return NewErrGHostConnect(h, err)
	}
	return nil
}

func (h *GHost) Command(command string) (string, error) {
	if h.sshClient == nil {
		if _, err := h.NewSSHClient(); err != nil {
			return "", err
		}
	}
	session, err := h.sshClient.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	var stderrBuf bytes.Buffer
	session.Stderr = &stderrBuf
	err = session.Run(command)
	if err != nil {
		return "", NewErrRunCommand(h, command, err)
	}
	return stdoutBuf.String(), nil
}

func (h *GHost) Upload(path string, content []byte) error {
	if h.sshClient == nil {
		if _, err := h.NewSSHClient(); err != nil {
			return err
		}
	}
	if h.sftpClient == nil {
		if _, err := h.NewSFTP(); err != nil {
			return err
		}
	}
	if err := h.sftpClient.MkdirAll(filepath.Dir(path)); err != nil {
		return err
	}
	file, err := h.sftpClient.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(content); err != nil {
		return err
	}
	return nil
}

func (h *GHost) NewSSHClient() (*ssh.Client, error) {
	var authMethod ssh.AuthMethod
	if h.Host.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(h.Host.PrivateKey))
		if err != nil {
			return nil, err
		}
		authMethod = ssh.PublicKeys(signer)
	} else {
		authMethod = ssh.Password(h.Host.Password)
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", h.Host.Host, h.Host.Port), &ssh.ClientConfig{
		User:            h.Host.User,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(h.Host.Timeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	h.sshClient = client
	return client, err
}

func (h *GHost) NewSFTP() (*sftp.Client, error) {
	if h.sshClient == nil {
		if _, err := h.NewSSHClient(); err != nil {
			return nil, err
		}
	}
	client, err := sftp.NewClient(h.sshClient)
	if err != nil {
		return nil, err
	}
	h.sftpClient = client
	return client, nil
}

func (h *GHost) Close() error {
	if h.sshClient != nil {
		if err := h.sshClient.Close(); err != nil {
			return err
		}
	}
	if h.sftpClient != nil {
		if err := h.sftpClient.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (h *GHost) SSHProtocol() string {
	return fmt.Sprintf("ssh://%s:%s@%s:%d", h.Host.User, h.Host.Password, h.Host.Host, h.Host.Port)
}

func (h *GHost) Base64SSHProtocol() string {
	return fmt.Sprintf("ssh://%s@%s:%d",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", h.Host.User, h.Host.Password))),
		h.Host.Host,
		h.Host.Port,
	)
}
