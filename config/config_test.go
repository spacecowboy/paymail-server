package config

import (
	"strings"
	"testing"
)

func TestNoFile(t *testing.T) {
	_, err := ReadConfig(("/path/to/no_such_file.toml"))

	if err == nil {
		t.Error("Expected error by got Nil")
	}
}

func TestEmptyFile(t *testing.T) {
	_, err := ReadConfig(("../.test/empty.toml"))

	if err == nil {
		t.Fatal("Expected error but got none")
	}

	if err.Error() != "Missing host in config" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func TestErrors(t *testing.T) {
	_, err := ReadConfig(("../.test/bad_payment_destination.toml"))

	if err == nil {
		t.Fatal("Expected error but got none")
	}

	if !strings.Contains(err.Error(), "Not a valid bitcoin address") {
		t.Errorf("Unexpected error message: %s", err)
	}

	_, err = ReadConfig(("../.test/bad_public_key.toml"))

	if err == nil {
		t.Fatal("Expected error but got none")
	}

	if !strings.Contains(err.Error(), "Not a valid public key") {
		t.Errorf("Unexpected error message: %s", err)
	}
}

func TestFullFile(t *testing.T) {
	config, err := ReadConfig(("../.test/full.toml"))

	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	if config.Server.Host == "" {
		t.Errorf("Expected host but got %s", config.Server.Host)
	}

	if config.Server.ListenAddress != "localhost:26245" {
		t.Errorf("Expected ListenAddres %s but got %s", "localhost:26245", config.Server.ListenAddress)
	}

	if !config.Bsv.Capabilities.SenderValidation {
		t.Error("Expected sender validation to be true")
	}

	if !config.Bsv.Capabilities.VerifyPublicKeyOwner {
		t.Error("Expected VerifyPublicKeyOwner to be true")
	}

	if len(config.Bsv.Accounts) != 2 {
		t.Errorf("Expected Accounts to have 2 accounts not %d", len(config.Bsv.Accounts))
	}

	a1 := config.Bsv.Accounts[0]

	if a1.Address != "bob@domain.tld" {
		t.Errorf("Incorrect address: %s", a1.Address)
	}

	if a1.PaymentDestination != "13WLwvqzGNcc14DVssidThWsfTaejZXr6r" {
		t.Errorf("Incorrect payment destination: %s", a1.PaymentDestination)
	}

	if len(a1.PublicKeys) != 2 {
		t.Errorf("Incorrect number of public keys: %d", len(a1.PublicKeys))
	}

	if a1.PublicKeys[0] != "038EA552AB83AD27BDDB42CA9B81E1597DD0BF2B9E24C567A7A6366B698F67F0AA" {
		t.Errorf("Incorrect first public key: %s", a1.PublicKeys[0])
	}

	if a1.PublicKeys[1] != "03AB9E665A4035A6C8B37C028C0331FB808445034172FB95E2B13FC0EB6B71AE9E" {
		t.Errorf("Incorrect first public key: %s", a1.PublicKeys[1])
	}

	a2 := config.Bsv.Accounts[1]

	if a2.Address != "alice@domain.tld" {
		t.Errorf("Incorrect address: %s", a2.Address)
	}

	if a2.PaymentDestination != "1N5WoeWXtf1iAKdAAdkZouxataAKfuGmm5" {
		t.Errorf("Incorrect payment destination: %s", a2.PaymentDestination)
	}

	if len(a2.PublicKeys) != 0 {
		t.Errorf("Incorrect number of public keys: %d", len(a2.PublicKeys))
	}
}
