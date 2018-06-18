package w2v

import (
	"bytes"
	"fmt"
	"testing"
)

var model = &Model{
	FeatureSize: 2,
	Word2id: map[string]int{
		"test": 0,
		"word": 1,
	},
	Vectors: []Vector{
		[]float64{1, 2},
		[]float64{3, 4},
	},
}

func sameModel(a, b *Model) error {
	if a.FeatureSize != b.FeatureSize {
		return fmt.Errorf("model FeatureSize isn't equal: %v != %v.", a.FeatureSize, b.FeatureSize)
	}
	if len(a.Word2id) != len(b.Word2id) {
		return fmt.Errorf("model word2id size doesn't match: %v != %v.", len(a.Word2id), len(b.Word2id))
	}
	for aw, awid := range a.Word2id {
		bwid, ok := b.Word2id[aw]
		if !ok {
			return fmt.Errorf("word %v in model a, but doesn't in model b.", aw)
		}
		if awid != bwid {
			return fmt.Errorf("word id is different for word %v: %v != %v.", aw, awid, bwid)
		}
		awv := a.Vectors[awid]
		bwv := b.Vectors[awid]
		if len(awv) != len(bwv) {
			return fmt.Errorf("word vector length is different for word %v: %v != %v", aw, len(awv), len(bwv))
		}
		for i := 0; i < len(awv); i++ {
			if awv[i] != bwv[i] {
				return fmt.Errorf("for word %v, awv[%v]=%v != bwv[%v]=%v.", aw, i, awv[i], i, bwv[i])
			}
		}
	}
	return nil
}

func TestW2V(t *testing.T) {
	// test model read/write
	var buf bytes.Buffer
	bytesWrote, err := model.Write(&buf)
	if err != nil {
		t.Error("failed to write model:", err)
		return
	}
	if bytesWrote != int64(buf.Len()) {
		t.Error(fmt.Sprintf("bytesWrote=%v, buf.Len()=%v", bytesWrote, buf.Len()))
		return
	}
	m, err := FromReader(&buf)
	if err != nil {
		t.Error("failed to create model from buf:", err)
		return
	}
	//fmt.Println(m)
	if err = sameModel(model, m); err != nil {
		t.Error("Model changed after write/read:", err)
		return
	}
	// test get vector
	v := m.GetVectorByWord("test")
	if v[0] != float64(1) || v[1] != float64(2) {
		t.Error("GetVectorByWord() returns wrong word vector", v)
		return
	}
	v = m.GetVectorByWordId(1) // "word"
	if v[0] != float64(3) || v[1] != float64(4) {
		t.Error("GetVectorByWordId() returns wrong word vector", v)
		return
	}
}
