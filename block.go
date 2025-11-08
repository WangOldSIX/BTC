package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

type Block struct {
	Data    []byte
	Hash    []byte
	PreHash []byte

	//新增的字段
	Version    uint64
	MerkleRoot []byte
	TimeStamp  uint64
	Difficulty uint64
	Nonce      uint64
}

func NewBlock(data string, preHash []byte) *Block {
	b := Block{
		Data:       []byte(data),
		PreHash:    preHash,
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: uint64(5),
	}
	//b.SetHash()

	pow := NewPOW(&b)
	hash, nonce := pow.Run()
	b.Hash = hash
	b.Nonce = nonce

	return &b
}

func (b Block) PrintBlock() {
	fmt.Printf("PREHASH:[%x]\n", b.PreHash)
	fmt.Printf("HASH:[%x]\n", b.Hash)
	fmt.Printf("DATA:\"%s\"\n", b.Data)
	fmt.Printf("TIMESTAMP:\"%v\"\n", b.TimeStamp)
	fmt.Printf("NONCE:\"%v\"\n", b.Nonce)
}

func (b Block) GoString() string {
	blockString := ""
	//blockString += fmt.Sprintf("VERSION:<%v>\n", b.Version)
	blockString += fmt.Sprintf("PREHASH:[%x]\n", b.PreHash)
	blockString += fmt.Sprintf("HASH:[%x]\n", b.Hash)
	blockString += fmt.Sprintf("DATA:\"%s\"\n", b.Data)
	blockString += fmt.Sprintf("TIMESTAMP:\"%v\"\n", b.TimeStamp)
	blockString += fmt.Sprintf("NONCE:\"%v\"\n", b.Nonce)

	return blockString
}

func (b *Block) SetHash() {
	var blockInfo []byte
	tmp := joinBlock(b)
	blockInfo = bytes.Join(tmp, make([]byte, 0))
	//log.Printf("blockInfo:%v", blockInfo)
	hash := sha256.Sum256(blockInfo)
	b.Hash = hash[:]
}

func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}
