package main

import (
	"fmt"
	"log"
	"go.etcd.io/bbolt"
)

const(
	DBName = "blockchain.db"
	BlocksBucket = "Blocks"
)

type BlockChain struct {
	Db *bbolt.DB
	tail []byte
}

func NewBlockChain() *BlockChain {

	dbName := addTimestampToFilename(DBName)
	db,err:=bbolt.Open(dbName, 0600, nil)
	if err!=nil{
		log.Fatal(err)
	}
	
	// 创建创世区块并添加到数据库
	genesis := newGenesis()
	err = AddElement2Db(db, string(genesis.Hash), Serialize(genesis))
	if err != nil {
		log.Fatal("Failed to add genesis block:", err)
	}
	
	return &BlockChain{Db: db, tail: genesis.Hash}
}

func newGenesis() *Block {
	return NewBlock(
		"THIS IS GENESIS BLOCK",
		[]byte{},
	)
}

func (bc *BlockChain)AddBlock(data string) error {
	prevBlock := bc.tail
	block := NewBlock(data, prevBlock)
	bc.tail = block.Hash
	return AddElement2Db(bc.Db,string(block.Hash),	Serialize(block))
}
func (bc BlockChain) PrintBC() {
	// log.Printf("TOTAL BLOCKS:%d\n", len(bc.Blocks))
	getAllElementsFromDb(bc.Db)
}

func AddElement2Db(db *bbolt.DB,key string,value []byte) error {
	return db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BlocksBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		bucket.Put([]byte(key),value)
		return nil
	})
}

func getAllElementsFromDb(db *bbolt.DB) error {
	return db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BlocksBucket))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		
		fmt.Println("===========================================")
		fmt.Println("         区块链数据库内容                   ")
		fmt.Println("===========================================")
		
		cursor := bucket.Cursor()
		var blockCount int
		
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			blockCount++
			block := Deserialize(v)
			
			fmt.Printf("\n[区块 #%d]\n", blockCount)
			fmt.Println("-------------------------------------------")
			fmt.Printf("区块哈希: %x\n", k)
			fmt.Println("区块详情:")
			fmt.Printf("  前区块哈希: [%x]\n", block.PreHash)
			fmt.Printf("  当前哈希:   [%x]\n", block.Hash)
			fmt.Printf("  数据内容:   %s\n", block.Data)
			fmt.Printf("  时间戳:     %v\n", block.TimeStamp)
			fmt.Printf("  随机数:     %v\n", block.Nonce)
			fmt.Printf("  难度值:     %v\n", block.Difficulty)
			
			if len(block.PreHash) == 0 {
				fmt.Println("  [创世区块]")
			}
		}
		
		fmt.Println("===========================================")
		fmt.Printf("总共 %d 个区块\n", blockCount)
		fmt.Println("===========================================")
		
		return nil
	})
}