// this was written in a few hours as fast as possible
// for the gopher gala.
// That should explain why it's really bad :)
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/sfreiberg/gotwilio"
)

var (
	listenAddr    = flag.String("listenAddr", "127.0.0.1:11111", "Web server listen address")
	wsListenAddr  = flag.String("wsListenAddr", ":11112", "Websocket server listen address")
	redisProtocol = flag.String("redisProtocol", "tcp", "Redis server protocol")
	redisAddress  = flag.String("redisAddress", "localhost:6379", "Redis server address")

	accountSid = flag.String("accountSID", "", "Twillio Account ID")
	authToken  = flag.String("authToken", "", "Twillio Auth Token")

	pool *redis.Pool

	wsClients = make(map[wsConn]int)
)

// websocket connection
type wsConn struct {
	websocket *websocket.Conn
	clientIP  string
}

// JSONUp represents one row of POSTed or collected JSON.
type JSONUp struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Value  string `json:"value"`
}

// jsonUpRecord is internal
type jsonUpRecord struct {
	User UpUser `json:"-"` // omit
	JSONUp
	ValueHistory []string `json:"sparkline"`
}

// UpUser is a record of a user of JSONUp
type UpUser struct {
	ID            string
	PhoneAreaCode string
	PhoneNumber   string
	VerifyCode    string
	Verified      bool
}

func (r jsonUpRecord) ID() string {
	return r.User.ID + ":" + r.Name
}

func redisSubscribeJSON(userID string) chan string {
	c := make(chan string)
	pubSubConn := redis.PubSubConn{Conn: pool.Get()}
	pubSubConn.Subscribe("user:" + userID)

	go func() {
		for {
			switch v := pubSubConn.Receive().(type) {
			case redis.Message:
				c <- string(v.Data)
			case redis.Subscription:
			case error:
				fmt.Printf("ERROR\n %s", v)
				return
			}
		}
	}()

	go publishUsersUpRecords(userID)

	return c
}

func wsServer(ws *websocket.Conn) {
	defer ws.Close()
	client := ws.Request().RemoteAddr

	sockID := wsConn{ws, client}
	wsClients[sockID] = 0

	log.Println("Client connected:", client)
	log.Println("for address:", ws.Request().URL.Path)
	log.Println("Websocket connections", len(wsClients))

	c := redisSubscribeJSON(ws.Request().URL.Path[1:])

	for {
		requestJSON := <-c
		err := websocket.Message.Send(ws, requestJSON)
		if err != nil {
			log.Println("Could not send to client, removing")
			delete(wsClients, sockID)
		}
	}
}

func startWsServer(port string) {
	log.Println("Websocket server Listening on " + port)
	err := http.ListenAndServeTLS(port, "server.crt", "server.key", websocket.Handler(wsServer))
	if err != nil {
		panic(err)
	}
}

func loadUser(userID string) (*UpUser, error) {
	conn := pool.Get()
	defer conn.Close()

	userjson, err := redis.String(conn.Do("GET", "user:"+userID))
	if err != nil {
		log.Printf("Redis GET Error %s", err)
		return nil, err
	} else {
		var upUser UpUser
		err = json.Unmarshal([]byte(userjson), &upUser)
		if err != nil {
			log.Printf("Marshal Error Error %s", err)
			panic(err)
		}
		return &upUser, nil
	}
}

func (u *UpUser) SaveUser() {
	userJson, err := json.Marshal(u)

	if err != nil {
		log.Printf("Marshal Error Error %s", err)
		panic(err)
	}

	conn := pool.Get()
	defer conn.Close()
	// TODO, maybe store the user as redis hash?
	_, err = conn.Do("SETNX", "user:"+u.ID, userJson)
	if err != nil {
		log.Printf("Redis SET Error %s", err)
		panic(err)
	}
}

func (u *UpUser) SendVerifyCode(ph_area string, ph_num string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnum := string(r.Int31())[:5]
	log.Printf(rnum)
	conn := pool.Get()
	defer conn.Close()

	message := "Your verification code for JSONUp is " + rnum
	twilio := gotwilio.NewTwilioClient(*accountSid, *authToken)
	twilio.SendSMS("+61421812892", "+"+ph_area+ph_num, message, "", "")

	return rnum
}

func saveUserEndpoint(w http.ResponseWriter, req *http.Request) {
	var user UpUser

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&user)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	idx := strings.LastIndex(req.URL.Path, "/")
	userID := req.URL.Path[idx+1:]

	if userID == user.ID {
		// save user
		user.SaveUser()
	} else {
		w.WriteHeader(422)
	}
}

func pushEndpoint(w http.ResponseWriter, req *http.Request) {
	var jsonCollection []JSONUp

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&jsonCollection)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	idx := strings.LastIndex(req.URL.Path, "/")
	userID := req.URL.Path[idx+1:]

	user, err := loadUser(userID)
	if err != nil {
		user = &UpUser{ID: userID}
		user.SaveUser()
	}

	for _, jsonRecord := range jsonCollection {
		r := jsonUpRecord{JSONUp: jsonRecord, User: *user}
		go pushToRedis(&r)
	}

	w.WriteHeader(200)
}

func pushToRedis(up *jsonUpRecord) (err error) {
	ju := up.JSONUp

	conn := pool.Get()
	defer conn.Close()

	key := "user:" + up.ID()

	_, err = conn.Do("SETEX", key+":status", 60, ju.Status)
	if err != nil {
		log.Printf("Redis Push Error %s", err)
		return
	}

	_, err = conn.Do("LPUSH", key+":values", ju.Value)
	if err != nil {
		log.Printf("Redis Push Error %s", err)
		return
	}

	// limit array to 20 entries
	_, err = conn.Do("LTRIM", key+":values", 0, 20)
	if err != nil {
		log.Printf("Redis Push Error %s", err)
		return
	}

	// add to users set of ups.
	_, err = conn.Do("SADD", "user:"+up.User.ID+":ups", ju.Name)
	if err != nil {
		log.Printf("Redis SADD Error %s", err)
		return
	}

	// Get sparkline data
	values, err := redis.Strings(conn.Do("LRANGE", key+":values", 0, 20))
	if err != nil {
		return
	}

	log.Printf("%s", values)
	up.ValueHistory = values

	// Publish Web Event.
	data, _ := json.Marshal(up)
	_, err = conn.Do("PUBLISH", "user:"+up.User.ID, data)
	if err != nil {
		log.Printf("Redis Push Error %s", err)
		return
	}

	return
}

func publishUsersUpRecords(userID string) {
	conn := pool.Get()
	defer conn.Close()

	_, err := loadUser(userID)
	if err != nil {
		return
	}

	iter := "0"
	for {
		result, err := conn.Do("SSCAN", "user:"+userID+":ups", iter, "COUNT", 20)
		if err != nil {
			log.Fatalln(err)
		}

		scanResults, ok := result.([]interface{})
		if !ok {
			log.Fatalln("Cannot cast scan results")
		}

		//log.Printf("%s", scanResults)
		// borrowed from https://github.com/garyburd/redigo/issues/81
		keys, ok := scanResults[1].([]interface{})
		if !ok {
			log.Fatalln("Cannot cast scan results")
		}

		for _, key := range keys {
			name, ok := key.([]byte)
			if !ok {
				log.Fatalln("Cannot cast key")
			}
			sendUpRecord(userID, string(name))
		}

		if b, ok := scanResults[0].([]byte); !ok || string(b) == "0" {
			//log.Printf("ok: %s", ok)
			//log.Println("Done fetching keys")
			break
		} else {
			iter = string(b)
		}

	}
}

func sendUpRecord(userID string, name string) {
	//log.Printf("name: %s", name)
	conn := pool.Get()
	defer conn.Close()

	status, err := redis.String(conn.Do("GET", "user:"+userID+":"+name+":status"))

	if status == "" {
		status = "GONE"
	}

	jsonUp := JSONUp{Name: name, Status: status}
	up := jsonUpRecord{JSONUp: jsonUp}

	// Get sparkline data
	sparkline, err := redis.Strings(conn.Do("LRANGE", "user:"+userID+":"+name+":values", 0, 20))
	if err != nil {
		return
	}

	//log.Printf("sparkline: %s", sparkline)
	up.ValueHistory = sparkline

	// Publish Web Event.
	data, _ := json.Marshal(up)
	_, err = conn.Do("PUBLISH", "user:"+userID, data)
	if err != nil {
		log.Printf("Redis Push Error %s", err)
		return
	}
}

func main() {
	flag.Parse()

	pool = newPool()

	router := mux.NewRouter()

	// Push endpoint
	router.HandleFunc("/push/{userId}", pushEndpoint).Methods("POST")

	// Save User
	router.HandleFunc("/saveUser/{userId}", saveUserEndpoint).Methods("POST")

	// Static public files
	publicFiles := http.FileServer(http.Dir("public"))
	router.Handle("/", publicFiles)

	// This is really dumb. #TODO, use strip prefix or something.
	router.Handle("/js/app.js", publicFiles)
	router.Handle("/css/app.css", publicFiles)

	// Start Web Server
	http.Handle("/", router)

	go startWsServer(*wsListenAddr)

	log.Println("Web server Listening on " + *listenAddr)
	err := http.ListenAndServe(*listenAddr, nil)
	if err != nil {
		panic(err)
	}
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(*redisProtocol, *redisAddress)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
