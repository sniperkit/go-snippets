package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

// TODO:
// 1. the implement are bad codes, because it's not
//    friendly for testing. time.Now() is a side effect
//    function, we cant generate a same ancestor block,
//    so we cant testing it.

var (
	ErrorWrongIndex    error = errors.New("block chain: index wrong")
	ErrorWrongPrevHash error = errors.New("block chain: previous hash wrong")
	ErrorWrongHash     error = errors.New("block chain: hash wrong")
)

// block chain
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

func NewAncestorBlock(BPM int) Block {
	var ancestor Block

	ancestor.Index = 0
	ancestor.Timestamp = time.Now().String()
	ancestor.BPM = BPM
	ancestor.PrevHash = ""
	ancestor.Hash = calculateHash(ancestor)

	return ancestor
}

func (b *Block) CalculateHash() string {
	b.Hash = calculateHash(*b)
	return b.Hash
}

func (b *Block) GenerateNextBlock(BPM int) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = b.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = b.Hash
	newBlock.CalculateHash()

	return newBlock, nil
}

func (b *Block) NextBlockValid(next Block) error {
	if next.Index != b.Index+1 {
		return ErrorWrongIndex
	}
	if next.PrevHash != b.Hash {
		return ErrorWrongPrevHash
	}
	if next.Hash != calculateHash(next) {
		return ErrorWrongHash
	}
	return nil
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
