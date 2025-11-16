package main

import (
	"fmt"
	"path/filepath"
	"time"
)

func joinBlock(b *Block) [][]byte {
	var tmp [][]byte
	tmp = [][]byte{
		b.PreHash,
		//b.Data,
		b.MerkleRoot,
		Uint64ToByte(b.Version),
		Uint64ToByte(b.TimeStamp),
		Uint64ToByte(b.Nonce),
		Uint64ToByte(b.Difficulty),
	}
	return tmp
}

// 添加时间戳到文件名的函数
func addTimestampToFilename(originalFilename string) string {
	// 获取文件的扩展名和基本名
	ext := filepath.Ext(originalFilename)
	baseName := originalFilename[:len(originalFilename)-len(ext)]

	// 获取当前时间并格式化为时间戳
	timestamp := time.Now().Format("20060102_150405")

	// 组合新的文件名：基本名_时间戳.扩展名
	newFilename := fmt.Sprintf("%s_%s%s", baseName, timestamp, ext)

	return newFilename
}
