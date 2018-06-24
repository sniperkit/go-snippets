package main

import (
	// "bufio"
	// "encoding/json"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"project_c9f7_i5l8_o0p4_p0j8/msg"
	"project_c9f7_i5l8_o0p4_p0j8/workerApp/worker"
	"runtime"

	"github.com/arcaneiceman/GoVector/govec"
)

// Logger is a GoVector logger for distributed logging
var Logger *govec.GoLog

func checkErr(err error) {
	_, _, line, _ := runtime.Caller(1)
	if err != nil {
		log.Fatal(fmt.Sprintf("Line %v -- %v", line, err))
	}
}

// HandleTasks receives tasks from the server and sends to the
// ServerMsgProcessor to be passed on to the Worker.
func handleTasks(conn net.Conn, msgProcessor *worker.ServerMsgProcessor) {
	var inBuf [512]byte
	var localAddr = conn.LocalAddr().String()
	log.Println("Waiting at addr: ", localAddr)

	var inMsg msg.FromServer
	// reader := bufio.NewReader(conn)
	for {
		var sizeBuf [2]byte
		_, sizeErr := conn.Read(sizeBuf[:])
		// log.Println("Size Buf contains:", sizeBuf)
		checkErr(sizeErr)
		msgSize := binary.LittleEndian.Uint16(sizeBuf[:])

		n, err := conn.Read(inBuf[0:msgSize])
		checkErr(err)
		Logger.UnpackReceive("Received-Message", inBuf[0:n], &inMsg)
		log.Printf("WConn: Received message %v\n", inMsg)
		msgProcessor.Process(inMsg)
	}

}

func runSender(conn net.Conn, outMsgChan chan msg.FromWorker) {
	for {
		select {
		case outMsg := <-outMsgChan:
			outBuf := Logger.PrepareSend(fmt.Sprintf("Sending-%v-Message", msg.TypeStr(outMsg.Type)), outMsg)
			var size uint16
			size = uint16(len(outBuf))
			sizeBuf := make([]byte, 2)
			binary.LittleEndian.PutUint16(sizeBuf, uint16(size))
			outBuf = append(sizeBuf, outBuf...)
			// log.Print("Adding sizeBuf to message with size: ", len(sizeBuf), "Msg size:", size, sizeBuf)
			n, err := conn.Write(outBuf)
			checkErr(err)
			log.Printf("WConn: Sent message %v, with %v bytes.\n", outMsg, n)
		}
	}
}

func main() {
	// Parse Arguments:
	serverAddr := os.Args[1]
	myAddr := os.Args[2]
	myID := string(os.Args[3])

	log.SetFlags(log.Lshortfile)

	// Initialize the Govec.
	Logger = govec.Initialize("worker"+myID, "workerlogfile"+myID)

	// Open connection to Server.
	conn, err := net.Dial("tcp", serverAddr)
	checkErr(err)
	var connMsg msg.WorkerConnectionMsg
	connMsg.WorkerId = msg.WorkerId(myID)
	connMsg.WorkerAddress = myAddr
	outBuf := Logger.PrepareSend("Worker-Connecting", connMsg)
	_, err = conn.Write(outBuf)
	checkErr(err)

	// Get the response from the server
	var inBuf [512]byte
	var response msg.WorkerConnectionResp
	_, err = conn.Read(inBuf[0:])
	checkErr(err)
	Logger.UnpackReceive("Received-Conn-Response", inBuf[0:], &response)
	log.Printf("Received response %v\n", response)

	if !response.IsAccepted {
		log.Println("Server refused connection.")
		os.Exit(1)
	} else {
		log.Println("Server accepted connection.")
	}

	outMsgChan := make(chan msg.FromWorker)
	smp := worker.NewServerMsgProcessor(msg.WorkerId(myID), outMsgChan)
	// Handle the tasks given to the
	go runSender(conn, outMsgChan)
	handleTasks(conn, smp)
}
