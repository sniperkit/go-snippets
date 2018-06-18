// word2vec model
// Be aware the model uses float64 for word vector values in memory
// but it casts them to float32 when writing to io.Writer/file. And
// when reading model from a io.Reader/file, it reads the word vector
// values in float32 and casts them to float64. This is to support
// existing w2v model file in the wild which word vector values are
// float32.
package w2v

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	//"time"
)

type Vector []float64

// word2vec model.
type Model struct {
	// Word vector size
	FeatureSize int

	// A map from word (string) to its id.
	Word2id map[string]int

	// An slice contains word vectors, index is the word id.
	Vectors []Vector
}

// Get the word vector by the word itself.
func (m *Model) GetVectorByWord(w string) Vector {
	if wid, ok := m.Word2id[w]; ok {
		return m.Vectors[wid]
	} else {
		return nil
	}
}

// Get the word vector by the word id.
func (m *Model) GetVectorByWordId(id int) Vector {
	if id >= len(m.Vectors) {
		return nil
	}
	return m.Vectors[id]
}

// Write the model (in binary format) to the given io.Writer.
func (m *Model) Write(w io.Writer) (int64, error) {
	totalBytesWrote := int64(0)
	bytesWrote, err := fmt.Fprintln(w, len(m.Word2id), m.FeatureSize)
	if err != nil {
		return -1, err
	}
	totalBytesWrote += int64(bytesWrote)

	space := []byte{32}
	wv := make([]float32, m.FeatureSize)
	vectorByteCnt := int64(reflect.TypeOf(wv[0]).Size()) * int64(m.FeatureSize)
	for word, wordId := range m.Word2id {
		// write word bytes
		bytesWrote, err = w.Write([]byte(word))
		if err != nil {
			return -1, err
		}
		totalBytesWrote += int64(bytesWrote)

		// write a space
		bytesWrote, err = w.Write(space)
		if err != nil {
			return -1, err
		}
		totalBytesWrote += int64(bytesWrote)

		// write the vector
		v := m.Vectors[wordId]
		for i := 0; i < m.FeatureSize; i++ {
			wv[i] = float32(v[i])
		}
		err = binary.Write(w, binary.LittleEndian, wv)
		if err != nil {
			return -1, err
		}
		totalBytesWrote += vectorByteCnt
	}
	return totalBytesWrote, nil
}

// Save model (in binary format) to the given path.
func (m *Model) WriteFile(path string) (int64, error) {
	w, err := os.Create(path)
	if err != nil {
		return -1, err
	}
	defer w.Close()

	return m.Write(w)
}

// Gzip the model (in binary format) and save it to the given path.
func (m *Model) WriteGzipFile(path string) (int64, error) {
	f, err := os.Create(path)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	return m.Write(gzip.NewWriter(f))
}

// Load word2vec model in binary format from the given io.Reader.
func FromReader(r io.Reader) (*Model, error) {
	var wordCnt, featureSize int
	n, err := fmt.Fscanln(r, &wordCnt, &featureSize)
	if err != nil {
		return nil, err
	}
	if n != 2 {
		return nil, fmt.Errorf("Failed to extract word count and feature size from input reader!")
	}
	//fmt.Printf("From model file: word-count: %v, feature-size: %v\n", wordCnt, featureSize)

	w2id := make(map[string]int, wordCnt)
	vectors := make([]Vector, wordCnt)

	data := make([]float64, featureSize*wordCnt)
	wv := make([]float32, featureSize)

	var savedWordCnt = 0
	var wBuf = make([]byte, 1)
	// byte array for the word, ignore the word if
	// it is longer than 100 bytes
	var wBytes [100]byte
	for i := 0; i < wordCnt; i++ {
		// 1. read the word and space after it
		j := 0
		for {
			_, err := r.Read(wBuf)
			if err != nil {
				return nil, err
			}
			if wBuf[0] == byte(32) {
				break
			}
			if j < 100 {
				wBytes[j] = wBuf[0]
			}
			j++
		}
		// 2. read the word vector
		if err := binary.Read(r, binary.LittleEndian, wv); err != nil {
			//fmt.Printf("failed to scan word vector, err: %v\n", err)
			return nil, err
		}
		// 3. save word and its vector to the model
		if j < 100 {
			w := strings.ToLower(string(wBytes[0:j]))
			// copy []float32 to []float64
			v := Vector(data[i*featureSize : (i+1)*featureSize])
			for k := 0; k < featureSize; k++ {
				v[k] = float64(wv[k])
			}
			w2id[w] = i
			vectors[i] = v
			savedWordCnt = savedWordCnt + 1
		} else {
			//fmt.Printf("ignored word #%v as it is longer than 100 bytes.", i)
		}
	}

	m := &Model{
		FeatureSize: featureSize,
		Word2id:     w2id,
		Vectors:     vectors,
	}

	//fmt.Printf("loaded %v word-vectors in %v.\n", savedWordCnt, time.Now().Sub(t))
	return m, nil
}

// Load word2vec model in binary format from the given gzipped model file.
func FromGzipFile(path string) (*Model, error) {
	//fmt.Printf("start loading %v\n", path)
	//t := time.Now()
	of, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer of.Close()

	f, err := gzip.NewReader(of)
	if err != nil {
		return nil, err
	}

	return FromReader(f)
}

// Load word2vec model in binary format from the given model file.
func FromFile(path string) (*Model, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return FromReader(f)
}
