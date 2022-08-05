package lib

import (
//"fmt"
"errors"
	"bytes"
	"encoding/hex"
	// "github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type Transaction struct {
	TxId               string `json:"txid"`
	SourceAddress      string `json:"source_address"`
	DestinationAddress string `json:"destination_address"`
	Amount             int64  `json:"amount"`
	tx                 string `json:"tx"`
}

type UTXO struct {
     Address btcutil.Address
     TxId string
     OutputIndex uint32
     Script []byte
     Satoshis int64
}

// debugTx(tx wire.MsgTx) {
	// sigScriptDisasm, _ := txscript.DisasmString(tx.TxIn[0].SignatureScript)
	// fmt.Println("sigScript:", tx)
// }

// Pay-To-Public-Key-Hash (P2PKH) transaction type.
//  <pubkey> OP_CHECKSIG
func GetPayToAddrScript(address btcutil.Address) []byte {
  receiveScript, _ := txscript.PayToAddrScript(address)
  return receiveScript
}

func isValidSignature(unspentTx UTXO, tx *wire.MsgTx, amount int64) bool {
  // Validate signature
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(unspentTx.Script, tx, 0, flags, nil, nil, amount)

	if err != nil {
		return false
	}

	if err := vm.Execute(); err != nil {
     return false
	}

  return true
}

func serializeTx(tx *wire.MsgTx) string {
  buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
  tx.Serialize(buf)
  return hex.EncodeToString(buf.Bytes())
}

func CreateTransaction(wifKey *btcutil.WIF, src btcutil.Address, dst btcutil.Address, amount int64, fee int64, vout uint32, lastHash string) (string, error) {

  unspentTx := UTXO {
    Address: src,
    TxId: lastHash,
    OutputIndex: vout,
    Script: GetPayToAddrScript(src),
    Satoshis: amount,
  }

	tx := wire.NewMsgTx(wire.TxVersion)
	utxoHash, _ := chainhash.NewHashFromStr(unspentTx.TxId)

  // Add TxIn
	tx.AddTxIn(
    wire.NewTxIn(
      wire.NewOutPoint(
        utxoHash,
        unspentTx.OutputIndex), nil, nil))

  // Add TxOut
	tx.AddTxOut(
    wire.NewTxOut(amount - fee, GetPayToAddrScript(dst)))

	signatureScript, err := txscript.SignatureScript(
    tx, 0, unspentTx.Script, txscript.SigHashAll, wifKey.PrivKey, false)
	if err != nil {
		return "", err
	}
	tx.TxIn[0].SignatureScript = signatureScript

  // debugTx(tx)

  if isValidSignature(unspentTx, tx, amount) != true {
    return "", errors.New("Invalid signature")
  }

  return serializeTx(tx), nil
}
