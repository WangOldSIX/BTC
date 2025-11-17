package main

import (
	"fmt"
	"log"

	"go.etcd.io/bbolt"
)

const (
	DBName       = "blockchain.db"
	BlocksBucket = "Blocks"
)

type BlockChain struct {
	Db   *bbolt.DB
	tail []byte
}

func NewBlockChain(address string) *BlockChain {

	dbName := addTimestampToFilename(DBName)
	db, err := bbolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 创建创世区块并添加到数据库
	genesis := newGenesis(address)
	log.Printf("GeneisBlock:%s\n", genesis)
	err = AddElement2Db(db, string(genesis.Hash), Serialize(genesis))
	if err != nil {
		log.Fatal("Failed to add genesis block:", err)
	}

	return &BlockChain{Db: db, tail: genesis.Hash}
}

func newGenesis(address string) *Block {
	coinBase := NewCoinBaseTx(address, "THIS IS GENESIS BLOCK")
	return NewBlock(
		//"THIS IS GENESIS BLOCK",
		[]*Transaction{coinBase},
		[]byte{},
	)
}

func (bc *BlockChain) AddBlock(txs []*Transaction) error {
	prevBlock := bc.tail
	block := NewBlock(txs, prevBlock)
	bc.tail = block.Hash
	return AddElement2Db(bc.Db, string(block.Hash), Serialize(block))
}
func (bc BlockChain) PrintBC() {
	// log.Printf("TOTAL BLOCKS:%d\n", len(bc.Blocks))
	getAllElementsFromDb(bc.Db)
}

func AddElement2Db(db *bbolt.DB, key string, value []byte) error {
	return db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BlocksBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		bucket.Put([]byte(key), value)
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
			fmt.Printf("  数据内容:   %s\n", block.Transactions[0].TXInputs[0].Sig)
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

// Find all UTXOS where address is prarameter which name is ADDRESS
func (bc *BlockChain) FindUTXOS(address string) []TxOutput {
	var UTXO []TxOutput

	//我们定一个map来保存消费过的utxo
	consumedUTXOs := make(map[string][]int64)

	//1.遍历所有区块
	it := bc.NewIterator()
	for {
		block := it.Next()
		//2.遍历每个区块的交易

		for _, tx := range block.Transactions {
			log.Printf("交易ID:%x\n", tx.TXID)
		OUTPUT: //打标签类似goto语句
			//3.遍历Output，找到自己赚的
			for i, output := range tx.TXOutputs {
				log.Printf("Current Index:%d, PubKeyHash:%s\n", i, output.PubKeyHash)
				//在这里做一个filter，将所有消耗过的outputs和当前output的index进行对比
				//如果当前output的index不在消耗过的outputs中，那么它就是未被消耗的output
				//如果当前output的index在消耗过的outputs中，那么它就是被消耗过的output
				if consumedUTXOs[string(tx.TXID)] == nil {
					for _, j := range consumedUTXOs[string(tx.TXID)] {
						if j == int64(i) {
							//当前准备添加的output已经消耗过了，不用再加了
							log.Printf("当前output的index:%d在消耗过的outputs中，那么它就是被消耗过的output", i)
							continue OUTPUT
						}
					}
				}

				//4.如果输出的PubKeyHash等于address，将其加入UTXO
				if output.PubKeyHash == address {
					UTXO = append(UTXO, output)
				}
			}
			//如果是挖矿交易，则不做遍历直接跳过
			if tx.IsCoinBase() {
				log.Printf("当前交易是挖矿交易，不做遍历直接跳过")
				continue
			} else {
				//3.b 遍历Input，找到自己花的
				for i, input := range tx.TXInputs {
					fmt.Printf("Current Index:%d\n", i)
					indexArray := consumedUTXOs[string(input.TXid)]
					indexArray = append(indexArray, input.Index)
				}
			}
		}

		if len(block.PreHash) == 0 {
			log.Println("到达创世区块,遍历结束")
			break
		}
	}

	return UTXO
}

type BlockChainIterator struct {
	db *bbolt.DB
	//游标，用于不断索引
	currentHashPointer []byte
}

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{
		bc.Db,
		//最初指向区块链的最后一个区块，随着Next的调用，不断变化
		bc.tail,
	}
}

// 迭代器是属于区块链的
// Next方式是属于迭代器的
// 1. 返回当前的区块
// 2. 指针前移
func (it *BlockChainIterator) Next() *Block {
	var block Block
	it.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BlocksBucket))
		if bucket == nil {
			log.Panic("迭代器遍历时bucket不应该为空，请检查!")
		}

		blockTmp := bucket.Get(it.currentHashPointer)
		//解码动作
		block = *Deserialize(blockTmp)
		//游标哈希左移
		it.currentHashPointer = block.PreHash

		return nil
	})

	return &block
}


func (bc *BlockChain) FindNeedUTXOs(from string, amount float64) (map[string][]int64, float64) {
	var utxos map[string][]int64
	//找到utxos里面包含的钱数
	var calc float64
	//TODO: 1. 找到所有未被消耗的UTXO
	//TODO: 2. 从未被消耗的UTXO中找到最合理的组合
	return utxos, calc
}