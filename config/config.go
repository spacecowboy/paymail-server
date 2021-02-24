package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/libsv/go-bt/bscript"
	toml "github.com/pelletier/go-toml"
)

type Capabilities struct {
	SenderValidation     bool `toml:"sender_validation" default:"true"`
	VerifyPublicKeyOwner bool `toml:"verify_public_key_owner"`
}

type Account struct {
	Address            string
	PaymentDestination string   `toml:"payment_destination"`
	PublicKey          string   `toml:"public_key"`
	PublicKeys         []string `toml:"public_keys"`
}

type BsvConfig struct {
	Capabilities Capabilities
	Accounts     []Account
}

type ServerConfig struct {
	Host          string
	ListenAddress string `toml:"listen_address" default:"localhost:26245"`
}

type Config struct {
	Server ServerConfig `toml:"server"`
	Bsv    BsvConfig    `toml:"bsv"`
}

func ReadConfig(filePath string) (*Config, error) {
	config := Config{}

	b, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(b, &config)

	if err != nil {
		return nil, err
	}

	if len(config.Server.Host) == 0 {
		return nil, errors.New("Missing host in config")
	}

	if len(config.Bsv.Accounts) == 0 {
		return nil, errors.New("No accounts defined in config")
	}

	for _, account := range config.Bsv.Accounts {
		if len(account.Address) == 0 {
			return nil, errors.New("Missing address in an account")
		}

		_, err := bscript.ValidateAddress(account.PaymentDestination)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Not a valid bitcoin address (%s): %s", err, account.PaymentDestination))
		}

		err = verifyPublicKey(account.PublicKey)
		if err != nil {
			return nil, err
		}
		for _, pubkey := range account.PublicKeys {
			err = verifyPublicKey(pubkey)
			if err != nil {
				return nil, err
			}
		}
	}

	return &config, err
}

func verifyPublicKey(publicKey string) error {
	if len(publicKey) != 66 {
		return errors.New(fmt.Sprintf("Not a valid public key (not 66 characters in length): %s", publicKey))
	}

	_, err := bscript.NewAddressFromPublicKeyString(publicKey, true)
	if err != nil {
		return errors.New(fmt.Sprintf("Not a valid public key (%s): %s", err, publicKey))
	}

	return nil
}
