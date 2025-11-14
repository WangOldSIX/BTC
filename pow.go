package main

import (
	"bytes"
	"crypto/sha256"
	"log"
	"math/big"
)

type POW struct {
	block  *Block
	target *big.Int
}

func NewPOW(b *Block) *POW {
	pow := POW{
		block: b,
	}
	targetString := generateTargetString(b.Difficulty)
	//log.Printf("难度值:%v,难度字符串:(%d,%s)", b.Difficulty, len(targetString), targetString)
	targetInt := big.Int{}
	targetInt.SetString(targetString, 16)
	pow.target = &targetInt
	return &pow
}

func generateTargetString(difficulty uint64) string {
	res := ""
	// 正确的工作量证明目标：前面有difficulty个0
	for i := 1; i <= 64; i++ {
		if i <= int(difficulty) {
			res += "0"
		} else {
			res += "f"
		}
	}
	return res
}

func (pow *POW) Run() (hash []byte, nonce uint64) {
	block := pow.block
	nonce = 0
	log.Printf("开始挖矿,难度值：%v,%s", block.Difficulty, generateTargetString(block.Difficulty))
	for {
		//在计算哈希前更新区块的nonce
		block.Nonce = nonce
		tmp := joinBlock(block)
		blockInfo := bytes.Join(tmp, make([]byte, 0))

		currentHash := sha256.Sum256(blockInfo)
		tmpInt := big.Int{}
		tmpInt.SetBytes(currentHash[:])

		//作比较
		if tmpInt.Cmp(pow.target) == -1 {
			log.Printf("挖到了，nonce：%v", nonce)
			hash = currentHash[:]
			break
		} else {
			if nonce%1000000 == 0 {
				log.Println("current nonce:", nonce)
			}
			nonce++
		}

	}
	return
}
