package main

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

//2. method of creating transaction
//3. create transaction
//4. overwrite main program (DATA->TRANSACTION)
