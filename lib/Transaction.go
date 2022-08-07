package lib

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	// "errors"
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

type UtxoStatus struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight int64  `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   int64  `json:"block_time"`
}

type UTXO struct {
	TxId   string     `json:"txid"`
	Vout   uint32     `json:"vout"`
	Status UtxoStatus `json:"status"`
	Value  int64      `json:"value"`
}

func debugTx(tx *wire.MsgTx) {
	pkScript, _ := txscript.DisasmString(tx.TxOut[0].PkScript)
	sigScript, _ := txscript.DisasmString(tx.TxIn[0].SignatureScript)
	fmt.Println("\t**** DEBUG ****")
	fmt.Println("\t - pkScript:", pkScript)
	fmt.Println("\t - sigScript:", sigScript)
	fmt.Println("\t***************")
}

// Pay-To-Public-Key-Hash (P2PKH) transaction type.
//  <pubkey> OP_CHECKSIG
func GetPayToAddrScript(address btcutil.Address) []byte {
	script, _ := txscript.PayToAddrScript(address)
	return script
}

func isValidSignature(script []byte, tx *wire.MsgTx, amount int64) error {
	// Validate signature
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(script, tx, 0, flags, nil, nil, amount)
	err = vm.Execute()
	return err
}

func serializeTx(tx *wire.MsgTx) string {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	tx.Serialize(buf)
	return hex.EncodeToString(buf.Bytes())
}

// See https://en.bitcoin.it/wiki/Transaction
func CreateTransaction(wifKey *btcutil.WIF, src btcutil.Address, dst btcutil.Address, amount int64, feeRate int64, utxos []UTXO) (string, error) {

	tx := wire.NewMsgTx(wire.TxVersion)

	var unspentAmount int64 = 0
	for _, utxo := range utxos {
		if unspentAmount > amount {
			break
		}

		// doubled SHA256-hashed of a (previous) to-be-used transaction
		utxoHash, _ := chainhash.NewHashFromStr(utxo.TxId)
		fmt.Println("* Double SHA256-hashed of Previous transaction: ", utxoHash)
		// Add TxIn
		tx.AddTxIn(
			wire.NewTxIn(
				wire.NewOutPoint(
					utxoHash,
					// Previous Txout-index
					utxo.Vout), nil, nil))
		unspentAmount += utxo.Value
	}

	// Add TxOut
	// 1. Amount for destination
	tx.AddTxOut(
		wire.NewTxOut(amount, GetPayToAddrScript(dst)))
	// 2. Remaining for source

	// fee = size in kb * feeRate
	fee := int64(math.Ceil(float64(tx.SerializeSize()/1000))) * int64(feeRate)
	if fee == 0 {
		fee = int64(feeRate)
	}

	subscript := GetPayToAddrScript(src)
	tx.AddTxOut(
		wire.NewTxOut(unspentAmount-amount-fee, subscript))

	for index, txIn := range tx.TxIn {
		signatureScript, err := txscript.SignatureScript(
			tx, index, subscript, txscript.SigHashAll, wifKey.PrivKey, false)

		txIn.SignatureScript = signatureScript
		if err != nil {
			return "", err
		}

	}

	debugTx(tx)

	return serializeTx(tx), nil
}
