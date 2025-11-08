package main

func joinBlock(b *Block) [][]byte {
	var tmp [][]byte
	tmp = [][]byte{
		b.PreHash,
		b.Hash,
		b.Data,
		b.MerkleRoot,
		Uint64ToByte(b.Version),
		Uint64ToByte(b.Nonce),
		Uint64ToByte(b.Difficulty),
	}
	return tmp
}
