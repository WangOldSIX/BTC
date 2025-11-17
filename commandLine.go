package main

import "fmt"

//The responsibility of this file is to implement complex command-line functionality.

func (cli *CLI) PrintBlockchain() {
	cli.bc.PrintBC()
}

func (cli *CLI) AddBlock(txs []*Transaction) {
	cli.bc.AddBlock(txs)
}

func (cli *CLI) GetBalance(address string) {
	utxos := cli.bc.FindUTXOS(address)
	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("[%s]的余额为:%f\n", address, total)
}
