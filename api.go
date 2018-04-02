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
	Height int
	Value  float64
	N      int
}

func (client Client) Tx(hash string) (Tx, error) {
	tx := Tx{}
	err := client.request("GET", fmt.Sprintf("tx/%v", hash), &tx)
	return tx, err
}

// Block

type Block struct {
	Tx []Tx
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

func (client Client) IsSpent(utxos []Utxo) ([]bool, error) {
	var isSpent []bool

	// checkmempool is "required" because of a bug
	// https://github.com/bitcoin/bitcoin/pull/12717
	url := fmt.Sprintf("getutxos/checkmempool%v", utxoListToUri(utxos))

	response := utxoResponse{}
	err := client.request("GET", url, &response)
	if err != nil {
		return isSpent, nil
	}

	for _, char := range response.Bitmap {
		isSpent = append(isSpent, char == '1')
	}

	return isSpent, err
}

func utxoListToUri(utxos []Utxo) string {
	s := ""
	for _, utxo := range utxos {
		s += fmt.Sprintf("/%v-%d", utxo.TxId, utxo.Vout)
	}
	return s
}
