package godb

/*
 bolt is used as local key-value storage, focus on following concerns:
 1. read -- goroutine safe, transaction support. can't modified when read,
 2. write -- goroutine safe, transaction supoort. only one write is allowed, write success or roll back
 3. delete --
 4. read all -- goroutine safe, transaction supoort. can't modified when read
 5. purge -- goroutine safe, transaction support. only one goroutines is allowed. purge success or roll back.

 *performance*
 1. read batch
 2. write batch
 3. make key sortable, use range selector

 Note: keys in bolt can be ordered.
*/

import (
	"fmt"
	"log"
	"testing"
)

// add data to bolt for testing
func init() {
	log.Println("init in test....")

	var key, value string
	for k := 0; k < 10; k++ {
		key = fmt.Sprintf("key%d", k+1)
		value = fmt.Sprintf("value%d", k+1)
		addBoltKV(key, []byte(value))
	}

	for k := 0; k < 10; k++ {
		key = fmt.Sprintf("key%d", k+1)
		value = fmt.Sprintf("value%d", k+1)
		addBoltKV(key, []byte(value))
	}
}

func TestInit(t *testing.T) {
	log.Println(boltDB.Info())
}

func TestGetBoltKV(t *testing.T) {
	log.Println(boltDB.Stats().TxStats.Write)
	log.Println(string(getBoltKV("key1")))
}

func TestGetAllBoltKV(t *testing.T) {
	for k, v := range getAllBoltKV() {
		log.Println("key ==> ", k, " value ==> ", string(v))
	}
}

func TestPurgeBolt(t *testing.T) {
	err := EmptyBoltBucket("default")
	if err != nil {
		t.Errorf("Failed to purge bolt bucket")
	}

	value := getBoltKV("key1")
	if value != nil {
		t.Errorf("After purge, there is no <k, v>, value: %v", string(value))
	}
}
