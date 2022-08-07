/*
Copyright © 2022 Alexandre Krispin <k.m.alexandre@protonmail.com>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	lib "github.com/alexandre-k/btc-transaction/lib"
	"github.com/btcsuite/btcutil"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"unicode"
)

type InputData struct {
	Destination string
	Amount      int64
	Fee         int64
	Vout        uint32
	LastHash    string
	Testnet     bool
	Utxos       string
	UtxoFile    string
}

func getInputOrFlag(input map[string]interface{}, cmd *cobra.Command, field string) string {
	r := []rune(field)
	r[0] = unicode.ToUpper(r[0])
	capitalizedField := string(r)
	flagInput, _ := cmd.Flags().GetString(field)
	if flagInput != "" {
		return flagInput
	} else {
		val := fmt.Sprint(input[capitalizedField])

		return val
	}
}

// createCmd represents the create command
var (
	utxos     []lib.UTXO
	jsonFile  string
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Bitcoin transaction",
		Long:  "Create transactions for the Bitcoin network from data feed through a json file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var input InputData
			if jsonFile != "" {
				content, err := ioutil.ReadFile(jsonFile)
				if err == nil {
					json.Unmarshal(content, &input)
				}
			}

			// ******* Read input parameters *********

			var inputMap map[string]interface{}
			data, _ := json.Marshal(input)
			json.Unmarshal(data, &inputMap)
			testnet, _ := strconv.ParseBool(getInputOrFlag(inputMap, cmd, "testnet"))
			destination := getInputOrFlag(inputMap, cmd, "destination")
			amount, _ := strconv.ParseInt(getInputOrFlag(inputMap, cmd, "amount"), 10, 64)
			fee, _ := strconv.ParseInt(getInputOrFlag(inputMap, cmd, "fee"), 10, 64)
			utxosContent := getInputOrFlag(inputMap, cmd, "utxos")
			utxoFile := getInputOrFlag(inputMap, cmd, "utxoFile")

			fmt.Println(utxosContent)
			if utxoFile == "" {
				json.Unmarshal([]byte(utxosContent), &utxos)
			} else {
				utxosContent, err := ioutil.ReadFile("utxos.json")
				fmt.Println(err)
				if err == nil {
					err = json.Unmarshal(utxosContent, &utxos)
				}
			}

			fmt.Println("\nInput parameters:\n")
			fmt.Println("\t- Amount: ", amount)
			fmt.Println("\t- Fee: ", fee)
			fmt.Println("\t- Testnet: ", testnet)
			fmt.Println("\t- UTXOs: ", utxos)
			if destination == "" || amount == 0 || len(utxos) == 0 || fee == 0 {
				fmt.Println("Parameter unknown. All parameters are necessary")
				return nil
			}

			// *****************

			wallet := lib.CreateWallet(testnet)

			cwd, _ := os.Getwd()

			var privateKeyFilename = filepath.Join(cwd, "/private.key")

			_, err := os.Stat(privateKeyFilename)

			var privKeyWIF *btcutil.WIF

			if err != nil {
				privKeyWIF, _ = wallet.CreatePrivateKey()
				os.WriteFile(privateKeyFilename, []byte(privKeyWIF.String()), 0600)
			} else {
				privateKeyFile, _ := os.ReadFile(privateKeyFilename)
				privateKeyContent := string(privateKeyFile)
				privKeyWIF, _ = wallet.ImportWIF(privateKeyContent)
			}

			source, _ := wallet.GetAddressPublicKey(privKeyWIF)
			// fmt.Println("\t- Public key uncompressed: ", source)
			// compressedSource, _ := wallet.GetAddress(privKeyWIF)
			// fmt.Println("\t- Public key compressed: ", compressedSource)
			// witnessPubKey := wallet.GetWitnessPubKeyHash(privKeyWIF)
			// fmt.Println("Witness ", witnessPubKey)
			sourceAddress, _ := wallet.GetDecodedAddress(source.EncodeAddress())
			destinationAddress, _ := wallet.GetDecodedAddress(destination)

			fmt.Println("\t- Transaction: ", sourceAddress, " => ", destinationAddress)

			transaction, err := lib.CreateTransaction(
				privKeyWIF, sourceAddress, destinationAddress, amount, fee, utxos)
			if err != nil {
				log.Fatal(err)
				return err
			}
			fmt.Println("\nOutput Transaction:\n")
			fmt.Println("\t", transaction)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVarP(&jsonFile, "input", "i", "", "A file containing all required data to create a transaction")

	createCmd.PersistentFlags().StringP(
		"privateKey", "p", "", "A private key to sign with")

	createCmd.PersistentFlags().StringP(
		"destination", "d", "", "A public key to transfer a given amount of Bitcoin to")

	createCmd.PersistentFlags().StringP(
		"amount", "a", "", "Amount to transfer in transaction")

	createCmd.PersistentFlags().StringP(
		"fee", "f", "", "Fee to pay for the transaction")

	createCmd.PersistentFlags().StringP("utxo", "u", "", "Previous source UTXOs to build a transaction upon")

	createCmd.PersistentFlags().StringP("utxoFile", "t", "", "File containing previous source UTXOs to build a transaction upon")

}
