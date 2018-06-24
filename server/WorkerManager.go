// WorkerService
package main

import (
	// "bufio"
	// "encoding/json"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"project_c9f7_i5l8_o0p4_p0j8/msg"
	// "github.com/arcaneiceman/GoVector/govec"
)

//====================================================================
// Data Structures

type Channels struct {
	cIn  chan msg.FromWorker
	cOut chan msg.FromServer
}

// Represents the state of a worker in the system. Stores the connection state and the work state of the worker.
type WorkerMetadata struct {
	WorkerId  msg.WorkerId
	Conn      net.Conn
	RequestId int
	cMsgOut   chan msg.FromServer
	cMsgIn    chan msg.FromWorker
	cQuit     chan int
}

// Store all the data related to workers.
var WorkerData map[msg.WorkerId]WorkerMetadata
var LastSelectedWorkers map[msg.WorkerId]bool

//====================================================================
type WorkerManager int

func (wm *WorkerManager) Initialize(serviceAddr string) {
	// Initialize data structures.
	WorkerData = make(map[msg.WorkerId]WorkerMetadata)
	LastSelectedWorkers = make(map[msg.WorkerId]bool)

	// Create the worker service.
	workerListener, err := net.Listen("tcp", serviceAddr)
	checkErr(err)
	fmt.Println("Listening for Workers at %v", serviceAddr)
	go handleWorkers(workerListener.(*net.TCPListener))
}

// TODO
// While it's true that we could just reque our request, that would require
// reassigning all workers, and have them read from the database, after each
// checkpoint, even if the set of workers hasn't changed.
func (wm *WorkerManager) AnyJoined() bool {
	return true
}

func (wm *WorkerManager) SelectWorkers() []msg.WorkerId {
	// if (len(LastSelectedWorkers) != 0) {
	// log.Fatal("SelectWorkers() called with populated LastSelectedWorkers")
	// }

	var selectedWorkers []msg.WorkerId
	for key := range WorkerData {
		selectedWorkers = append(selectedWorkers, key)
		LastSelectedWorkers[key] = true
	}

	return selectedWorkers
}

func Writer(worker WorkerMetadata) {
	defer capturePanic(worker.WorkerId)

	for {
		select {
		case msg := <-worker.cMsgOut:
			SendMessage(worker.Conn, msg)
		case <-worker.cQuit:
			log.Printf("Writer %v terminating\n", worker.WorkerId)
			return
		}
	}
}

func Reader(workerId msg.WorkerId, conn *net.TCPConn) {
	defer capturePanic(workerId)

	log.Printf("Reader() %v started conn: %v\n", workerId, conn)
	// reader := bufio.NewReader(conn)
	for {
		msg := ReadMessage(conn)
		log.Printf("Reader() %v received message %v\n", workerId, msg)

		// LOCK
		for {
			workerData, ok := WorkerData[workerId]
			if !ok {
				log.Printf("Reader read for non-existent worker %v\n", workerId)
				return
			} else if workerData.cMsgIn == nil {
				log.Printf("Reader() %v, no cMsgIn\n", workerId)
			} else {
				workerData.cMsgIn <- msg
				break
			}
		}
	}
}

func (wm *WorkerManager) PrepareWorkers(workers []msg.WorkerId, cIn chan msg.FromWorker) {
	// TODO: acquire lock
	for _, worker := range workers {
		data, ok := WorkerData[worker]
		if !ok {
			log.Printf("Preparing a non-existent worker.\n")
		} else {
			data.cMsgIn = cIn
			data.cMsgOut = make(chan msg.FromServer, 1000)
			data.cQuit = make(chan int)
			WorkerData[worker] = data
			go Writer(data)
			log.Printf("Worker prepared %v\n", data)
		}
	}
}

// TODO:
func capturePanic(workerId msg.WorkerId) {
	if r := recover(); r != nil {
		log.Printf("Recovered %v\n", r)
		DeleteWorker(workerId)
		// var ok bool
		// _, ok = r.(error)
		// if !ok {
		//     log.Printf("Captured panic: %v\n", r)
		// }
	}
}

func (wm *WorkerManager) SendMessageToWorker(msg msg.FromServer) {
	defer capturePanic(msg.DstWorker)

	// TODO: acquire lock
	workerData, ok := WorkerData[msg.DstWorker]
	if ok == false {
		log.Printf("Sending message %v non-existent worker.\n", msg)
	} else {
		workerData.cMsgOut <- msg
	}
}

func (wm *WorkerManager) CloseWorkers(workers []msg.WorkerId) {
	// Since the superstep is done, we want to kill the listener processes we created above.
	log.Println("work(): deleting workers.")
	for _, worker := range workers {
		data, ok := WorkerData[worker]
		if !ok {
			log.Println("work(): deleting worker doesn't exist.")
		} else {
			data.cQuit <- 1
		}
	}

	// Create a new map here,
	LastSelectedWorkers = make(map[msg.WorkerId]bool)
}

//============================================================
func handleWorkers(l *net.TCPListener) {
	for {
		conn, err := l.Accept()
		checkErr(err)
		go handleWorkerConn(conn)
	}
}

//====================================================================
func handleWorkerConn(conn net.Conn) {
	var wcm msg.WorkerConnectionMsg
	var resp msg.WorkerConnectionResp

	var inBuf [512]byte
	n, err := conn.Read(inBuf[0:])
	checkErr(err)
	Logger.UnpackReceive("Worker-Conn", inBuf[0:n], &wcm)
	if err == nil {
		// Refuse a connection if this worker in the current selection.
		_, ok := LastSelectedWorkers[wcm.WorkerId] // TODO: lock
		if ok {
			log.Printf("Refusing conn from Worker %v. Already in current selection.\n", wcm.WorkerId)
			resp.WorkerId = wcm.WorkerId
			resp.IsAccepted = false
		} else {
			// Since this worker is not in the current selection, we can add it to map.
			var workerData WorkerMetadata
			workerData.WorkerId = wcm.WorkerId
			workerData.Conn = conn
			WorkerData[wcm.WorkerId] = workerData
			resp.IsAccepted = true
			tcpconn := conn.(*net.TCPConn)
			go Reader(wcm.WorkerId, tcpconn)
			log.Printf("Added worker %v.\n", workerData.WorkerId)
		}
	} else {
		log.Println("Worker did not send a WorkerConnectionMsg as its first msg. Refusing and closing Connection.")
		resp.WorkerId = "Badconnectionparam"
		resp.IsAccepted = false
	}

	outBuf := Logger.PrepareSend("Sending-Resp-To-Worker", resp)
	n, err = conn.Write(outBuf)
	checkErr(err)

	if n != len(outBuf) || err != nil {
		DeleteWorker(wcm.WorkerId)
	}
	return
}

func DeleteWorker(workerId msg.WorkerId) {
	// TODO: lock
	log.Printf("Deleting worker %v\n", workerId)
	delete(WorkerData, workerId)
}

//====================================================================
// Worker TCP Service
// TODO: We should be resilient to connection failures in Send() and Read(). At the moment we fail straight away.
func SendMessage(conn net.Conn, message msg.FromServer) {
	log.Printf("SendMessage() to addr: %v Worker:%v\n", conn.RemoteAddr().String(), message.DstWorker)
	outBuf := Logger.PrepareSend(fmt.Sprintf("Sending-%v-Message", msg.TypeStr(message.Type)), message)
	// log.Printf("%v\n", outBuf)
	var msgSize uint16
	msgSize = uint16(len(outBuf))
	sizeBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(sizeBuf, uint16(msgSize))
	outBuf = append(sizeBuf, outBuf...)
	size, err := conn.Write(outBuf)
	checkErr(err)
	if size != len(outBuf) {
		log.Fatal("size of message out != size written")
	}
}

// func ReadMessage(reader *bufio.Reader) msg.FromWorker {
func ReadMessage(conn *net.TCPConn) msg.FromWorker {
	var msg msg.FromWorker
	var inBuf [512]byte

	var sizeBuf [2]byte
	_, sizeErr := conn.Read(sizeBuf[:])
	// log.Println("Size Buf contains:", sizeBuf)
	checkErr(sizeErr)
	msgSize := binary.LittleEndian.Uint16(sizeBuf[:])

	n, err := conn.Read(inBuf[0:msgSize])
	// log.Println("WMAN: Reading", n, "bytes from connection.")
	checkErr(err)
	Logger.UnpackReceive("Reading-Message", inBuf[0:n], &msg)
	return msg
}
