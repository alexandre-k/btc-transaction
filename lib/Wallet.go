package lib

import (
"errors"
"github.com/btcsuite/btcd/btcec"
"github.com/btcsuite/btcd/chaincfg"
"github.com/btcsuite/btcutil"
)

type Wallet struct {
     Testnet bool
     Config *chaincfg.Params
     // PrivateKey btcutil.WIF
}

func CreateWallet(testnet bool) *Wallet {
  config := &chaincfg.TestNet3Params
  if testnet {
    config = &chaincfg.TestNet3Params
  } else {
    config = &chaincfg.MainNetParams
  }
  return &Wallet{ Testnet: testnet, Config: config }
}

func (w *Wallet) CreatePrivateKey() (*btcutil.WIF, error) {
     secret, err := btcec.NewPrivateKey(btcec.S256())
     if err != nil {
       return nil, err
     }
     return btcutil.NewWIF(secret, w.Config, true)
}

func (w *Wallet) ImportWIF(wifStr string) (*btcutil.WIF, error) {
  wif, err := btcutil.DecodeWIF(wifStr)
  if err != nil {
    return nil, err
  }

  if !wif.IsForNet(w.Config) {
    return nil, errors.New("Detected invalid WIF")
  }

  return wif, nil
}

func (w *Wallet) GetAddress(wif *btcutil.WIF) (*btcutil.AddressPubKey, error) {
     return btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), w.Config)
}

func (w *Wallet) GetAddressPublicKey (wif *btcutil.WIF) (*btcutil.AddressPubKey, error) {
	addressPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), w.Config)
  if err != nil {
     return nil, err
  }
  return addressPubKey, nil
}

func (w *Wallet) GetDecodedAddress (address string) (btcutil.Address, error) {
	sourceAddress, err := btcutil.DecodeAddress(address, w.Config)
	if err != nil {
		return nil, err
	}
  return sourceAddress, nil
}