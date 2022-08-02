/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	transaction "github.com/alexandre-k/btc-transaction/lib"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"strconv"
	"unicode"
)

type InputData struct {
	PrivateKey  string
	Destination string
	Amount      int64
	LastHash    string
}

func getInputOrFlag(input map[string]interface{}, cmd *cobra.Command, field string) string {
	r := []rune(field)
	r[0] = unicode.ToUpper(r[0])
	capitalizedField := string(r)
	if input[field] == "" {
		flagInput, _ := cmd.Flags().GetString(field)
		return flagInput
	} else {
		val := fmt.Sprint(input[capitalizedField])

		return val
	}
}

// createCmd represents the create command
var (
	jsonFile  string
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Bitcoin transaction",
		Long:  "Create transactions for the Bitcoin network from data feed through a json file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var input InputData
			if jsonFile != "" {
				content, err := ioutil.ReadFile(jsonFile)
				if err != nil {
					log.Fatal("Error while reading the file ", jsonFile, ". ", err)
				}

				err = json.Unmarshal(content, &input)
				if err != nil {
					log.Fatal("Unable to unmarshal: ", err)
				}
			}

			var inputMap map[string]interface{}
			data, _ := json.Marshal(input)
			json.Unmarshal(data, &inputMap)
			privateKey := getInputOrFlag(inputMap, cmd, "privateKey")
			destination := getInputOrFlag(inputMap, cmd, "destination")
			amount, _ := strconv.ParseInt(getInputOrFlag(inputMap, cmd, "amount"), 10, 64)
			lastHash := getInputOrFlag(inputMap, cmd, "lastHash")

			transaction, err := transaction.CreateTransaction(
				privateKey, destination, amount, lastHash)
			if err != nil {
				fmt.Println(err)
				return err
			}
			tx, _ := json.Marshal(transaction)
			fmt.Println(string(tx))
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
		"lastHash", "l", "", "Previous source UTXO hash to build a transaction upon")
}
