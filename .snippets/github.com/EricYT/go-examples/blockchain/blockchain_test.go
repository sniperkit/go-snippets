package blockchain

import "testing"

func TestBlockchain1(t *testing.T) {
	BPM := 30
	ancestor := NewAncestorBlock(BPM)
	firshBlock, _ := ancestor.GenerateNextBlock(BPM)
	if err := ancestor.NextBlockValid(firshBlock); err != nil {
		t.Fatalf("first block got a wrong validate error %s", err)
	}
	secondBlock, _ := firshBlock.GenerateNextBlock(BPM)
	if err := ancestor.NextBlockValid(secondBlock); err == nil {
		t.Fatalf("second block is not the next one of ancestor")
	}
}
