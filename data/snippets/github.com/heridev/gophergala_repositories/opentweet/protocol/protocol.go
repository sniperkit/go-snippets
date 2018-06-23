package protocol

import (
	"net"
	"fmt"
	"time"
	"bufio"
	"strconv"
	"log"
)

type Tweet struct {
	Time time.Time
	Text string
}

type Server struct {
	port int
	fn Callback
}

type Callback func(name string, from, to time.Time) ([]Tweet, error)

var port = 12315

func GetTweets(server, name string, from, to time.Time) ([]Tweet, error) {
	serverString := fmt.Sprintf("%s:%v", server, port)
	conn, err := net.Dial("tcp", serverString)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	fmt.Fprintf(conn, "OT v1\r\n")
	fmt.Fprintf(conn, "%s\r\n", name)
	fmt.Fprintf(conn, "%v\r\n", from.Unix())
	fmt.Fprintf(conn, "%v\r\n", to.Unix())
	// if to == 0 {
	// 	fmt.Fprintf(conn, "now\r\n")
	// } else {
	// 	fmt.Fprintf(conn, "%v\r\n", to.Unix())
	// }
	
	reader := bufio.NewReader(conn)
	tweets := make([]Tweet, 0)
	for {
		timeStamp, err := readStringCRLF(reader)
		if err != nil {
			break
		}
		timeUnix, err := strconv.ParseInt(timeStamp, 10, 64)
		if err != nil {
			return nil, err
		}
		tweetText, err := readStringCRLF(reader)
		if err != nil {
			return nil, err
		}
		tweet := Tweet {time.Unix(timeUnix, 0), tweetText}
		tweets = append(tweets, tweet)
		log.Printf("Got a tweet from server %v %v\n", timeUnix, tweetText)
	}
	return tweets, nil
}

func readStringCRLF(reader *bufio.Reader) (string, error) {
	str, err := reader.ReadString('\r')
	if err != nil {
		return "", err
	}
	// remove \r
	str = str[:len(str)-1]
	// assume next char is \n
	_, err = reader.ReadByte()
	if err != nil {
		return "", err
	}
	return str, nil
}

func NewServer() (Server) {
	server := Server{port, nil}
	return server
}

func (server *Server) Register(fn Callback) error {
	server.fn = fn
	return nil
}

func (server *Server) ListenAndServe() error {
	lStr := fmt.Sprintf(":%v", server.port)
	listen, err := net.Listen("tcp", lStr)
	if err != nil {
		return err
	}
	log.Printf("Listening\n")
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}
		log.Printf("Got a Connection\n")
		go server.connection(conn)
	}
	return nil
}
	
func (server *Server) connection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	magic, err := readStringCRLF(reader)
	if err != nil {
		log.Printf("Could not read magic string: %v", err)
		return
	}
	if magic != "OT v1" {
		log.Printf("Invalid magic string: %s", magic)
		return
	}
	name, err := readStringCRLF(reader)
	if err != nil {
		log.Printf("Could not read name: %v", err)
		return
	}
	from, err := readStringCRLF(reader)
	if err != nil {
		log.Printf("Could not read from time: %v", err)
		return
	}
	fromUnix, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		log.Printf("Could not convert from time: %v", err)
		return
	}
	to, err := readStringCRLF(reader)
	if err != nil {
		log.Printf("Could not read to time: %v", err)
		return
	}
	toUnix, err := strconv.ParseInt(to, 10, 64)
	if err != nil {
		log.Printf("Could not convert to time: %v", err)
		return
	}
	
	log.Printf("Calling callback %v %v %v\n", name, fromUnix, toUnix)
	tweets, err := server.fn(name, time.Unix(fromUnix, 0), time.Unix(toUnix, 0))
	if err != nil {
		log.Printf("Could not retrieve tweets: %v", err)
		return
	}

	log.Printf("Got %v tweets from backend\n", len(tweets))
	for _, tweet := range tweets {
		fmt.Fprintf(conn, "%v\r\n", tweet.Time.Unix())
		fmt.Fprintf(conn, "%v\r\n", tweet.Text)
	}
	conn.Close()
}
