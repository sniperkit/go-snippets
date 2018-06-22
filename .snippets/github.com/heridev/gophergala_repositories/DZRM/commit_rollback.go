package goxa

import (
	"log"
)

func Commit(xa XA, id []byte) (int, error) {

	count, err := xa.ConnHandle.Write(append([]byte(CMIT), id...))
	log.Printf("Sent CMIT id %s.\n", id)

	return count, err
}

func Rollback(xa XA, id []byte) (int, error) {

	count, err := xa.ConnHandle.Write(append([]byte(BACK), id...))
	log.Printf("Sent BACK id %s.\n", id)

	return count, err
}
