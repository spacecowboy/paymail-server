package v1

import (
	"testing"
	"gitlab.com/spacecowboy/paymail-server/config"
)

func TestWellKnownNoHost(t *testing.T) {
	config := config.Config {}
	_, err := getWellKnownBsvAliasResponse(&config)
	
	if err == nil {
		t.Fatalf("Expected error but got no error")
	}
}

func TestWellKnownNoVerify(t *testing.T) {
	config := config.Config {
		Server: config.ServerConfig{
			Host: "example.org",
		},
		Bsv: config.BsvConfig{
			Capabilities: config.Capabilities{
				VerifyPublicKeyOwner: false,
			},
		},
	}
	js, err := GetWellKnownBsvAliasResponse(&config)
	
	if err != nil {
		t.Fatalf("Got error %s", err)
	}

	expected := `{"bsvalias":"1.0","capabilities":{"pki":"https://example.org/bsvalias/v1/pki/{alias}@{domain.tld}","paymentDestination":"https://example.org/bsvalias/v1/address/{alias}@{domain.tld}"}}`

	if string(js) != expected {
		t.Errorf("Expected %s\n but got %s", expected, string(js))
	}
}

/*
	Example output of curl -L moneybutton.com/.well-known/bsvalias

	{
  "bsvalias": "1.0",
  "capabilities": {
    "6745385c3fc0": false,
    "pki": "https://www.moneybutton.com/api/v1/bsvalias/id/{alias}@{domain.tld}",
    "paymentDestination": "https://www.moneybutton.com/api/v1/bsvalias/address/{alias}@{domain.tld}",
    "a9f510c16bde": "https://www.moneybutton.com/api/v1/bsvalias/verifypubkey/{alias}@{domain.tld}/{pubkey}",
    "f12f968c92d6": "https://www.moneybutton.com/api/v1/bsvalias/public-profile/{alias}@{domain.tld}",
    "5f1323cddf31": "https://www.moneybutton.com/api/v1/bsvalias/receive-transaction/{alias}@{domain.tld}",
    "2a40af698840": "https://www.moneybutton.com/api/v1/bsvalias/p2p-payment-destination/{alias}@{domain.tld}"
  }
}
	*/

func TestWellKnownWithVerify(t *testing.T) {
	config := config.Config {
		Server: config.ServerConfig{
			Host: "example.org",
		},
		Bsv: config.BsvConfig{
			Capabilities: config.Capabilities{
				VerifyPublicKeyOwner: true,
			},
		},
	}
	js, err := GetWellKnownBsvAliasResponse(&config)
	
	if err != nil {
		t.Fatalf("Got error %s", err)
	}

	expected := `{"bsvalias":"1.0","capabilities":{"pki":"https://example.org/bsvalias/v1/pki/{alias}@{domain.tld}","paymentDestination":"https://example.org/bsvalias/v1/address/{alias}@{domain.tld}","a9f510c16bde":"https://example.org/bsvalias/v1/verifypubkey/{alias}@{domain.tld}/{pubkey}"}}`

	if string(js) != expected {
		t.Errorf("Expected %s\n but got %s", expected, string(js))
	}
}

func TestPki(t *testing.T) {
	account := config.Account{
					Address: "bob@example.org",
					PublicKey: "pkiaddress",
				}
	js, err := getPkiResponse(&account)
	
	if err != nil {
		t.Fatalf("Got error %s", err)
	}

	expected := `{"bsvalias":"1.0","handle":"bob@example.org","pubkey":"pkiaddress"}`

	if string(js) != expected {
		t.Errorf("Expected %s\n but got %s", expected, string(js))
	}
}

func TestPaymentDestinationCompressedKey(t *testing.T) {
	// NewAddressFromPublicKeyString on the large key from https://bsvalias.org/04-01-basic-address-resolution.html
	// Generated from public key 027c1404c3ecb034053e6dd90bc68f7933284559c7d0763367584195a8796d9b0e
	account := config.Account{
			PaymentDestination: "1jSiW7mnRLspAEX8UsFtQSxXWGFuUufMG",
		}
	js, err := getPaymentDestinationResponse(&account)
	
	if err != nil {
		t.Fatalf("Got error %s", err)
	}

	expected := `{"output":"76a9140806efc8bedc8afb37bf484f352e6f79bff1458c88ac"}`

	if string(js) != expected {
		t.Errorf("Expected %s\n but got %s", expected, string(js))
	}
}

func TestVerify(t *testing.T) {
	js, err := getVerifyPublicKeyOwnerResponse("bob@example.org", "pubkey", false)

	if err != nil {
		t.Fatalf("Got error %s", err)
	}

	expected := `{"bsvalias":"1.0","handle":"bob@example.org","pubkey":"pkiaddress"}`

	if string(js) != expected {
		t.Errorf("Expected %s\n but got %s", expected, string(js))
	}

	js, err = getVerifyPublicKeyOwnerResponse("bob@example.org", "pubkey", true)

	if err != nil {
		t.Fatalf("Got error %s", err)
	}

	expected = `{"bsvalias":"1.0","handle":"bob@example.org","pubkey":"pkiaddress"}`

	if string(js) != expected {
		t.Errorf("Expected %s\n but got %s", expected, string(js))
	}
}