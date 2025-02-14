package ssls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/log"
	"github.com/go-acme/lego/v4/providers/http/webroot"
	"github.com/go-acme/lego/v4/registration"
	"os"
	"time"
)

var (
	CADirURL           = lego.LEDirectoryProduction
	DefaultCertTimeout = 30
)

func CreateNonExistingFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0o700)
	} else if err != nil {
		return err
	}
	return nil
}

func SetWebRootChallenge(client *lego.Client, webRootPath string) error {
	provider, err := webroot.NewHTTPProvider(webRootPath)
	if err != nil {
		return err
	}
	return client.Challenge.SetHTTP01Provider(provider)
}

func LegoClient(email string) (account *Account, client *lego.Client, err error) {
	storage, err := NewAccountsStorage(email, CADirURL)
	if err != nil {
		return nil, nil, err
	}
	account, err = storage.LoadAccount()
	if err != nil {
		log.Infof("Failed to load account  %s", err.Error())
	}
	if account == nil {
		account = &Account{Email: email}
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return account, nil, err
		}
		account.Key = privateKey
	}
	config := lego.NewConfig(account)
	config.CADirURL = CADirURL
	config.Certificate = lego.CertificateConfig{
		KeyType:             certcrypto.RSA2048,
		Timeout:             time.Duration(DefaultCertTimeout) * time.Second,
		OverallRequestLimit: certificate.DefaultOverallRequestLimit,
	}
	config.UserAgent = defaultUserAgent
	client, err = lego.NewClient(config)
	if err != nil {
		return account, nil, err
	}
	account.Registration, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return account, nil, err
	}
	_ = storage.Save(account)
	return account, client, nil
}

func ObtainCertificate(client *lego.Client, domains []string) (*certificate.Resource, error) {
	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}
	return client.Certificate.Obtain(request)
}
