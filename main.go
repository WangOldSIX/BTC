package main

import "fmt"

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	bc := NewBlockChain()
	cli:=CLI{bc}
	cli.Run()
	
}

func testAddBlock(blockSize int,bc *BlockChain){
	for i := 0; i < blockSize; i++ {
		bc.AddBlock(fmt.Sprintf("这是第%d块区块", i+3))
	}
}


