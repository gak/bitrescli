package bitrescli

import (
	"fmt"
)

type Tx struct {
	TxId string
	Vin  []TxInput
	Vout []TxOutput
}

type TxInput struct {
	Txid string
	Vout int
}

type TxOutput struct {
	Height       int
	Value        float64
	Txid         string
	N            int
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

type ScriptPubKey struct {
	Addresses []string
	Type      string
}

func (client Client) Tx(hash string) (Tx, error) {
	tx := Tx{}
	err := client.request("GET", fmt.Sprintf("tx/%v", hash), &tx)
	return tx, err
}

// Block

type Block struct {
	Height int
	Hash   string
	Tx     []Tx
}

func (client Client) Block(hash string) (Block, error) {
	block := Block{}
	err := client.request("GET", fmt.Sprintf("block/%v", hash), &block)
	return block, err
}

// Utxo (IsSpent)

type Utxo struct {
	TxId string
	Vout int
}

type utxoResponse struct {
	Bitmap string
	Utxos  []TxOutput
}

func (client Client) BulkUtxos(utxos []Utxo) ([]TxOutput, error) {
	var collected []TxOutput

	// This is the most you can request at once
	chunkSize := 15

	for start := 0; start < len(utxos); start += chunkSize {
		end := start + chunkSize

		if end > len(utxos) {
			end = len(utxos)
		}

		chunk := utxos[start:end]

		response, err := client.Utxos(chunk)
		if err != nil {
			return collected, err
		}
		collected = append(collected, response...)
	}

	return collected, nil
}

func (client Client) Utxos(utxos []Utxo) ([]TxOutput, error) {
	var collected []TxOutput

	// checkmempool is "required" because of a bug
	// https://github.com/bitcoin/bitcoin/pull/12717
	url := fmt.Sprintf("getutxos/checkmempool%v", utxoListToUri(utxos))

	response := utxoResponse{}
	err := client.request("GET", url, &response)
	if err != nil {
		return collected, err
	}

	outputIndex := 0

	for idx, char := range response.Bitmap {
		if char == '1' {
			reference := utxos[idx]
			output := response.Utxos[outputIndex]

			// N is not specified when calling utxos
			output.N = reference.Vout

			// Tack on the Txid for reference
			output.Txid = reference.TxId

			collected = append(collected, output)

			outputIndex ++
		}
	}

	return collected, err
}

func utxoListToUri(utxos []Utxo) string {
	s := ""
	for _, utxo := range utxos {
		s += fmt.Sprintf("/%v-%d", utxo.TxId, utxo.Vout)
	}
	return s
}
