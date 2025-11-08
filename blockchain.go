package main

import (
	"fmt"
	"log"
)

type BlockChain struct {
	Blocks []*Block
}

func NewBlockChain() *BlockChain {
	genesisBlock := newGenesis()
	return &BlockChain{Blocks: []*Block{genesisBlock}}
}

func newGenesis() *Block {
	return NewBlock(
		"THIS IS GENESIS BLOCK",
		[]byte{},
	)
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	block := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, block)
}

func (bc BlockChain) PrintBC() {
	log.Printf("TOTAL BLOCKS:%d\n", len(bc.Blocks))
	for i := 0; i < len(bc.Blocks); i++ {
		fmt.Println("CURRENT BLOCK HEIGHT:", i+1)
		fmt.Printf("%#v\n", bc.Blocks[i])
	}
}
