package main

import (
	"fmt"
	"time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	bc := NewBlockChain()
	cli := CLI{bc}
	cli.Run()
	cli.PrintBlockchain()
	printCurrentTime()
}

func testAddBlock(blockSize int, bc *BlockChain) {
	for i := 0; i < blockSize; i++ {
		bc.AddBlock(fmt.Sprintf("这是第%d块区块", i+3))
	}
}

func printCurrentTime() {
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("当前时间：", t)
}
