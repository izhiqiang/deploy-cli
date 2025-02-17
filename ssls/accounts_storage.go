package ssls

import (
	"crypto/ecdsa"
	"crypto/x509"
	"deploy-cli/env"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	baseKeysFolderName = "keys"
	accountFileName    = "account.json"
)

const (
	filePerm os.FileMode = 0o600
)

var (
	defaultUserAgent = fmt.Sprintf("deploy-cli/%s", env.Version)
	savePath         = path.Join(env.HomeDir, "ssl")
)

type AccountsStorage struct {
	Email           string
	CADirURL        string
	accountFilePath string
	keysPath        string
}

// NewAccountsStorage Creates a new AccountsStorage.
func NewAccountsStorage(email, CADirURL string) (*AccountsStorage, error) {
	serverURL, err := url.Parse(CADirURL)
	if err != nil {
		return nil, err
	}
	serverPath := strings.NewReplacer(":", "_", "/", string(os.PathSeparator)).Replace(serverURL.Host)
	accountsPath := filepath.Join(savePath, serverPath)
	rootUserPath := filepath.Join(accountsPath, email)
	if err := CreateNonExistingFolder(rootUserPath); err != nil {
		return nil, err
	}
	return &AccountsStorage{
		Email:           email,
		CADirURL:        CADirURL,
		keysPath:        filepath.Join(rootUserPath, baseKeysFolderName),
		accountFilePath: filepath.Join(rootUserPath, accountFileName),
	}, nil
}

func (s *AccountsStorage) Save(account *Account) error {
	var (
		wg errgroup.Group
	)
	wg.Go(func() error {
		jsonBytes, err := json.MarshalIndent(account, "", "\t")
		if err != nil {
			return err
		}
		return os.WriteFile(s.accountFilePath, jsonBytes, filePerm)
	})
	wg.Go(func() error {
		privateKeyBytes, err := x509.MarshalECPrivateKey(account.Key.(*ecdsa.PrivateKey))
		if err != nil {
			return err
		}
		return os.WriteFile(s.keysPath, privateKeyBytes, filePerm)
	})
	return wg.Wait()
}
func (s *AccountsStorage) LoadAccount() (*Account, error) {
	fileBytes, err := os.ReadFile(s.accountFilePath)
	if err != nil {
		return nil, err
	}
	var account Account
	err = json.Unmarshal(fileBytes, &account)
	if err != nil {
		return nil, err
	}
	keyBytes, err := os.ReadFile(s.keysPath)
	if err != nil {
		return nil, err
	}
	privateKey, err := x509.ParseECPrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}
	account.Key = privateKey
	return &account, nil
}
