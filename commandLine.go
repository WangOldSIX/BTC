package main

//The responsibility of this file is to implement complex command-line functionality.

func (cli *CLI)PrintBlockchain(){
	cli.bc.PrintBC()
}

func (cli *CLI)AddBlock(data string){
	cli.bc.AddBlock(data)
}