package goxa

import (
	"code.google.com/p/go-uuid/uuid"
	"log"
)

func Wait(previousXa XA, currentXa XA) {

	log.Printf("Waiting for CMIT/BACK.\n")
	buffer := make([]byte, 1024)
	count, err := currentXa.ConnHandle.Read(buffer)

	if count > 0 && err == nil {

		if string(buffer)[:len(CMIT)] == CMIT {

			log.Printf("CMIT id %s Received.\n", buffer[len(CMIT):])

			if currentXa.Queue[string(buffer[len(CMIT):len(CMIT)+UUIDLength])] != nil {

				log.Printf("Committed the request.\n")
				delete(currentXa.Queue, string(buffer)[len(CMIT):len(CMIT)+UUIDLength])
				currentXa.ID = ""

			} else if previousXa.ConnHandle != currentXa.ConnHandle && previousXa.ConnHandle != nil {

				log.Printf("Pass CMIT to previous syncpoint\n")
				Commit(previousXa, buffer[len(CMIT):len(CMIT)+UUIDLength])
			}
		} else if string(buffer)[:len(BACK)] == BACK {

			log.Printf("BACK id %s Received.\n", buffer[len(BACK):])

			if currentXa.Queue[string(buffer[len(BACK):len(BACK)+UUIDLength])] != nil {

				log.Printf("Rolledback and resend the request\n")
				Send(previousXa, currentXa, currentXa.Queue[string(buffer[len(BACK):len(BACK)+UUIDLength])], true)

			} else if previousXa.ConnHandle != currentXa.ConnHandle && previousXa.ConnHandle != nil {

				log.Printf("Pass BACK to previous syncpoint\n")
				Rollback(previousXa, buffer[len(BACK):len(BACK)+UUIDLength])

			}
		}
	}
}

func Send(previousXa XA, currentXa XA, buffer []byte, trans bool) (int, error) {

	if trans {

		if currentXa.ID == "" {

			id := uuid.New()
			currentXa.Queue[id] = buffer
			count, err := currentXa.ConnHandle.Write(append([]byte(UOW+id), buffer...))
			log.Printf("Sent %d bytes.\n", len(buffer))

			go Wait(previousXa, currentXa)

			return count - (len(UOW) + UUIDLength), err

		} else {

			count, err := currentXa.ConnHandle.Write(append([]byte(UOW+currentXa.ID), buffer...))
			log.Printf("Sent %d bytes.\n", len(buffer))

			return count - (len(UOW) + UUIDLength), err
		}

	} else {

		count, err := currentXa.ConnHandle.Write(buffer)
		log.Printf("Sent %d bytes.\n", len(buffer))

		return count, err
	}
}

func Receive(xa XA) ([]byte, []byte, int, error) {

	buffer := make([]byte, 1024)
	count, err := xa.ConnHandle.Read(buffer)
	log.Printf("Received %d bytes.\n", len(buffer))

	if string(buffer)[:len(UOW)] == UOW {

		return buffer[len(UOW) : len(UOW)+UUIDLength], buffer[len(UOW)+UUIDLength:], count - (len(UOW) + UUIDLength), err
	}

	return nil, buffer, count, err
}
