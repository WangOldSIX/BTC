package main

import "fmt"

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	bc := NewBlockChain()
	bc.AddBlock("2")
	for i := 0; i < 1; i++ {
		bc.AddBlock(fmt.Sprintf("这是第%d块区块", i+3))
	}
	bc.PrintBC()
}

func testString() {
	str := generateTargetString(4)
	fmt.Println(len(str), str)
}
