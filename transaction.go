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
	Value float64
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
	output := TxOutput{PubKeyHash: address, Value: REWARD}
	//对于coinbase交易来说，只有一个input和output
	tx := Transaction{make([]byte, 0), []TxInput{input}, []TxOutput{output}}
	tx.SetHash()
	return &tx
}

func (tx *Transaction) IsCoinBase() bool {
	//1.交易的input只有一个
	//2.交易id为空
	//3.交易index为-1
	return len(tx.TXInputs) == 1 && tx.TXInputs[0].Index == -1 && bytes.Equal(tx.TXInputs[0].TXid, make([]byte, 0))
}

// 3. create common transaction(创建普通转账交易)
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	// 找到最合理的UTXO集合 map[string][]int64
	utxos, resValue := bc.FindNeedUTXOs(from, amount)
	var inputs []TxInput
	var outputs []TxOutput

	if resValue < amount {
		log.Panic("Transaction failed. Not enough balance, resValue:", resValue, "amount:", amount)
	}
	// 将UTXO一个一个转成inputs
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TxInput{[]byte(id), int64(i), from}
			inputs = append(inputs, input)
		}
	}
	// 创建outputs
	output := TxOutput{
		Value:      amount,
		PubKeyHash: to,
	}
	outputs = append(outputs, output)
	//如果有零钱，要找零
	if resValue > amount {
		output := TxOutput{
			Value:      resValue - amount,
			PubKeyHash: from,
		}
		outputs = append(outputs, output)
	}
	tx := Transaction{make([]byte, 0), inputs, outputs}
	tx.SetHash()
	return &tx
}

//4. overwrite main program (DATA->TRANSACTION)
