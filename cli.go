package main

import (
	"fmt"
	"os"
)

/*
	This is a file that receives command line arguments
	to control the behavior of the blockchain
*/

type CLI struct {
	bc *BlockChain
}

const Usage = `
Usage:
	addBlock --data [DATA]  "add data to the blockchain"
	printChain  "print all blocks in the blockchain"
	getBalance --address ADDRESS  "get UTXOs where address is ADDRESS"

Example:
	addBlock --data "Hello Blockchain"
`

/*
We put the operation into a function, which can recive parameters
*/
func (cli *CLI) Run() {
	//1.get all command line arguments
	args := os.Args
	if len(args) < 2 {
		fmt.Println(Usage)
		return
	}

	//2.analyze the command line arguments
	cmd := args[1]
	switch cmd {
	/*case "addBlock":
	//add blocks
	fmt.Println("addBlock")
	//ensure the command is valid
	if len(args) == 4 && args[2] == "--data" {

		//a.get the data
		data:=args[3]
		//b.handle the data
		cli.AddBlock(data)
	}else{
		fmt.Println("Invalid command, please check the usage")
		fmt.Println(Usage)
		return
	}*/

	case "printChain":
		//print chain
		fmt.Println("printChain")
		cli.PrintBlockchain()
	case "getBalance":
		fmt.Println("getBalance")
		if len(args) == 4 && args[2] == "--address" {

			address := args[3]
			cli.GetBalance(address)
		} else {
			fmt.Println("Invalid command, please check the usage")
			fmt.Println(Usage)
			return
		}
	default:
		fmt.Println("Invalid command, please check the usage")
		fmt.Println(Usage)
		return
	}

	//3.execute the command

}
