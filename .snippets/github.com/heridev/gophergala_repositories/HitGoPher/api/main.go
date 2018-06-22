package main

import (
    "net/http"
    "golang.org/x/net/websocket"
    "time"
    "fmt"
    "math/rand"
)

// Declare global variable
var MapPlayer map[string] int = map[string] int {}
var channel []*websocket.Conn = make([]*websocket.Conn, 0)
var golbalNumber string = ""

type T struct {
    Action string
    Position string
    Player string
}



func randomGopher() (string) {
    var position = []string{"A1","A2","A3",
                            "B1","B2","B3",
                            "C1","C2","C3",}

    rand.Seed(time.Now().UnixNano())
    randNumber := rand.Intn(8)
    return position[randNumber]
}

func newPlayer(name string) (){
    MapPlayer[name] = 0
}

func winnerAddPoint(name string, point int) {
    MapPlayer[name] += point

}
func sentToClient(ws *websocket.Conn) {
    //var position = randomGopher()
    //var point = getCurrentPoint()
    //var data = prepairData(position, point)

    channel = append(channel, ws)
    var receiveData T
    for {
        err := websocket.JSON.Receive(ws,&receiveData)
        if err != nil {
            return
        }
        fmt.Println("======Received======")
        fmt.Println(receiveData)
        if receiveData.Action == "newPlayer" {
            newPlayer(receiveData.Player)
            websocket.JSON.Send(ws, map[string] interface{} {
                        "position" : golbalNumber,
                        "pointInfo" : MapPlayer,
                    });
                }

            if golbalNumber == receiveData.Position {
                winnerAddPoint(receiveData.Player,1)
                golbalNumber = randomGopher()
                for _, value := range channel {
                    websocket.JSON.Send(value, map[string] interface{} {
                        "position" : golbalNumber,
                        "pointInfo" : MapPlayer,
                    });
                }
                fmt.Println("Position : " + golbalNumber)
                fmt.Println(channel)
            }
        
    }
}

func main() {
    golbalNumber = randomGopher()
    fmt.Println(golbalNumber)
    http.Handle("/start", websocket.Handler(sentToClient))
    err := http.ListenAndServe(":12345", nil)
    if err != nil {
        panic("ListenAndServe: "+ err.Error())
    }
}