//Package survey supports the implementation of a survey server.
package survey

import (
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/surveyor"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
	"log"
	"time"
)

type Server struct {
	url        string
	sock       mangos.Socket
	surveySent chan string
}

type ResponseHandler func([]byte)

//Starts a Survey on the specified url for a specified duration for clients to repsond.  A set of workers can run to handle more traffic.
func (self *Server) Listen(url string, ms time.Duration, workers int, handler ResponseHandler) error {

	self.url = url
	self.surveySent = make(chan string)

	var err error

	if self.sock, err = surveyor.NewSocket(); err != nil {
		return err
	}
	self.sock.AddTransport(ipc.NewTransport())
	self.sock.AddTransport(tcp.NewTransport())

	if err = self.sock.Listen(url); err != nil {
		return err
	}
	err = self.sock.SetOption(mangos.OptionSurveyTime, time.Millisecond*ms)
	if err != nil {
		return err
	}
	for id := 0; id < workers; id++ {
		go self.processData(handler)
	}

	return nil

}

//Send the survey question to clients and set a channel of slice messages to process.
func (self *Server) Send(payload []byte) error {

	var err error

	if err = self.sock.Send(payload); err != nil {
		return err
	}

	self.surveySent <- "Survey Sent"

	return nil
}

//Handles the survey responses.
func (self *Server) processData(handler ResponseHandler) {

	var msg []byte
	var err error

	//Wait for a survey to be sent
	s := <-self.surveySent
	log.Println(s)

	for {
		if msg, err = self.sock.Recv(); err != nil {
			continue
		}
		go handler(msg)
	}
}

// var src = rand.NewSource(time.Now().UnixNano())

// func RandStringBytesMaskImprSrc(n int) string {
// 	b := make([]byte, n)
// 	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
// 	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
// 		if remain == 0 {
// 			cache, remain = src.Int63(), letterIdxMax
// 		}
// 		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
// 			b[i] = letterBytes[idx]
// 			i--
// 		}
// 		cache >>= letterIdxBits
// 		remain--
// 	}

// 	return string(b)
// }
