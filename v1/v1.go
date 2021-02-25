package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/libsv/go-bt/bscript"
	"gitlab.com/spacecowboy/paymail-server/config"
)


type V1Handler struct{
  config *config.Config
}

func GetHandler(config *config.Config) (http.Handler) {
  v1 := new(V1Handler)
  v1.config = config

  r := chi.NewRouter()

  r.Get("/.well-known/bsvalias", v1.wellKnownBsvAlias)
  r.Get("/bsvalias/v1/pki/{alias}", v1.Pki)
  r.Post("/bsvalias/v1/address/{alias}", v1.PaymentDestination)

  if v1.config.Bsv.Capabilities.VerifyPublicKeyOwner {
    r.Get("/bsvalias/v1/verifypubkey/{alias}/{pubkey}", v1.VerifyPublicKeyOwner)
  }

  return r
}

type WellKnownBsvAliasResponse struct {
  Version string `json:"bsvalias"`
  Capabilities BsvAliasCapabilities `json:"capabilities"`
}

type BsvAliasCapabilities struct {
  Pki string `json:"pki"`
  PaymentDestination string `json:"paymentDestination"`
  VerifyPublicKey string `json:"a9f510c16bde,omitempty"`
}

func GetWellKnownBsvAliasResponse(config *config.Config) ([]byte, error) {
  if config == nil {
    return nil, errors.New("Config was nil")
  }

  host := config.Server.Host

  if host == "" {
    return nil, errors.New("Missing Host definition in config")
  }

  data := WellKnownBsvAliasResponse{
    Version: "1.0",
    Capabilities: BsvAliasCapabilities{
      Pki: fmt.Sprintf("https://%s/bsvalias/v1/pki/{alias}@{domain.tld}", host),
      PaymentDestination: fmt.Sprintf("https://%s/bsvalias/v1/address/{alias}@{domain.tld}", host),
    },
  }

  if config.Bsv.Capabilities.VerifyPublicKeyOwner {
    data.Capabilities.VerifyPublicKey = fmt.Sprintf("https://%s/bsvalias/v1/verifypubkey/{alias}@{domain.tld}/{pubkey}", host)
  }

  return json.Marshal(data)
}

func (v1 *V1Handler) wellKnownBsvAlias(w http.ResponseWriter, r *http.Request) {
  js, err := GetWellKnownBsvAliasResponse(v1.config)

  if err != nil {
    http.Error(w, "Server Error", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

type PkiResponse struct{
  Version string `json:"bsvalias"`
  Handle string `json:"handle"`
  PubKey string `json:"pubkey"`
}

func getPkiResponse(account *config.Account) ([]byte, error) {
  data := PkiResponse{
    Version: "1.0",
    Handle: account.Address,
    PubKey: account.PublicKey,
  }

  return json.Marshal(data)
}

// /bsvalias/v1/pki/{alias}
func (v1 *V1Handler) Pki(w http.ResponseWriter, r *http.Request) {
  alias := chi.URLParam(r, "alias")

  for _, account := range v1.config.Bsv.Accounts {
    if account.Address == alias {
      js, err := getPkiResponse(&account)

      if err != nil {
        http.Error(w, "Server Error", http.StatusInternalServerError)
        return
      }

      w.Header().Set("Content-Type", "application/json")
      w.Write(js)
      return
    }
  }

  http.Error(w, "Unknown account", http.StatusNotFound)
}

type PaymentDestinationRequest struct{
  SenderName string `json:"senderName"`
  SenderHandle string `json:"senderHandle"`
  Dt string `json:"dt"`
  Amount int64 `json:"amount"`
  Purpose string `json:"purpose"`
  Signature string `json:"signature"`
}

type PaymentDestinationResponse struct{
  // P2PKH output script
  Output string `json:"output"`
}

func getPaymentDestinationResponse(account *config.Account) ([]byte, error) {
  if account == nil {
    return nil, errors.New("Account was nil")
  }
  
  log.Printf("JTSF")
  // a, _ := bscript.NewAddressFromPublicKeyString(account.PaymentDestination, true)
  // log.Printf("%s = %s", a.AddressString, a.PublicKeyHash)

  script, err := bscript.NewP2PKHFromAddress(account.PaymentDestination)

  if err != nil {
    return nil, err
  }

  data := PaymentDestinationResponse{
    Output: script.ToString(),
  }

  return json.Marshal(data)
}

func getPaymentDestinationRequest(r *http.Request) (*PaymentDestinationRequest, error) {
  result, err := ioutil.ReadAll(r.Body)

  if err != nil {
    return nil, err
  }

  var request PaymentDestinationRequest
  err = json.Unmarshal(result, &request)

  return &request, err
}

// //bsvalias/v1/address/{alias}
func (v1 *V1Handler) PaymentDestination(w http.ResponseWriter, r *http.Request) {
  alias := chi.URLParam(r, "alias")

  _, err := getPaymentDestinationRequest(r)

  if err != nil {
    http.Error(w, "Server Error", http.StatusInternalServerError)
  }

  for _, account := range v1.config.Bsv.Accounts {
    if account.Address == alias {
      js, err := getPaymentDestinationResponse(&account)

      if err != nil {
        http.Error(w, "Server Error", http.StatusInternalServerError)
        return
      }

      w.Header().Set("Content-Type", "application/json")
      w.Write(js)
      return
    }
  }

  http.Error(w, "Unknown account", http.StatusNotFound)
}

type VerifyPublicKeyOwnerResponse struct{
  Handle string `json:"handle"`
  PubKey string `json:"pubkey"`
  Match bool `json:"match"`
}

func getVerifyPublicKeyOwnerResponse(handle string, pubkey string, match bool) ([]byte, error) {
  data := VerifyPublicKeyOwnerResponse{
    Handle: handle,
    PubKey: pubkey,
    Match: match,
  }
  
  return json.Marshal(data)
}

// /bsvalias/v1/verifypubkey/{alias}/{pubkey}
func (v1 *V1Handler) VerifyPublicKeyOwner(w http.ResponseWriter, r *http.Request) {
  alias := chi.URLParam(r, "alias")
  pubkey := chi.URLParam(r, "pubkey")

  match := false

  for _, account := range v1.config.Bsv.Accounts {
    if account.Address == alias {
      if account.PublicKey == pubkey{
        match = true
        break
      }
      for _, key := range account.PublicKeys {
        if key == pubkey {
          match = true
          break
        }
      }
    }
  }

  js, err := getVerifyPublicKeyOwnerResponse(alias, pubkey, match)

  if err != nil {
    http.Error(w, "Server Error", http.StatusInternalServerError)
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}
