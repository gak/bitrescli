package main

import (
	"fmt"
	"github.com/gak/bitrescli"
	"os"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
	"github.com/davecgh/go-spew/spew"
	"strings"
	"strconv"
)

type CLI struct {
	BaseURI string `required:"true"`
	Debug   bool

	Tx struct {
		Hash string `required:"true" arg:"true"`
	}

	Block struct {
		Hash string `required:"true" arg:"true"`
	}

	Utxo struct {
		TxVout []string `required:"true" arg:"true"`
	}
}

func main() {
	cmd, cli := parse()

	client := bitrescli.Client{
		Debug:          cli.Debug,
		BaseURI:        cli.BaseURI,
		RequestTimeout: 5,
	}
	client.Connect()

	switch cmd {
	case "tx":
		dump(client.Tx(cli.Tx.Hash))
	case "block":
		dump(client.Block(cli.Block.Hash))
	case "utxo":
		var utxos []bitrescli.Utxo

		for _, utxo := range cli.Utxo.TxVout {
			bits := strings.Split(utxo, "-")
			vout, err := strconv.ParseInt(bits[1], 10, 0)
			if err != nil {
				panic(err)
			}
			utxos = append(utxos, bitrescli.Utxo{TxId: bits[0], Vout: int(vout)})
		}

		dump(client.Utxos(utxos))
	}
}

func dump(data interface{}, err error) {
	if err != nil {
		fmt.Println(err)
	}

	spew.Dump(data)
}

func parse() (string, CLI) {
	var king = kingpin.New("bitrescli", "")

	cli := CLI{}
	var err = king.Struct(&cli)
	if err != nil {
		fmt.Println(err)
	}

	command := kingpin.MustParse(king.Parse(os.Args[1:]))

	return command, cli
}
