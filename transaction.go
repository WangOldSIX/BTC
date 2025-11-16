package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

const REWARD = 12.5

// 1. TRANSACTION struct
type Transaction struct {
	TXID      []byte
	TXInputs  []TxInput  //Transaction inputs array
	TXOutputs []TxOutput //Transaction outputs array
}
type TxInput struct {
	//1.ID
	TXid []byte
	//2.Index
	Index int64
	//3.解锁脚本
	Sig string
}

type TxOutput struct {
	//transfer balance
	value float64
	//script
	PubKeyHash string
}

// 设置交易ID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

// 2. method of creating transaction(CoinBase挖矿交易)
func NewCoinBaseTx(address string, data string) *Transaction {
	//Miner 挖矿时无需指定签名，所以sig字段可以由miner自己填写
	input := TxInput{make([]byte, 0), -1, data}
	output := TxOutput{PubKeyHash: address, value: REWARD}
	//对于coinbase交易来说，只有一个input和output
	tx := Transaction{make([]byte, 0), []TxInput{input}, []TxOutput{output}}
	tx.SetHash()
	return &tx
}

//3. create transaction
//4. overwrite main program (DATA->TRANSACTION)
