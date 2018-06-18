package main

import (
    "fmt"
    "bitbucket.org/gdamore/mangos"
    "bitbucket.org/gdamore/mangos/protocol/push"
    "bitbucket.org/gdamore/mangos/protocol/pull"
    "bitbucket.org/gdamore/mangos/transport/all"
)

func main() {
    url := "tcp://127.0.0.1:8000"

    // 1. Create Pull Server
    pullServerReady := make(chan struct{})
    pullServer, err := pull.NewSocket()
    defer pullServer.Close()

    all.AddTransports(pullServer)

    // 2. Run Pull Server
    go func() {
        var err error
        var serverMsg *mangos.Message

        if err = pullServer.Listen(url); err != nil {
            fmt.Printf("\nServer listen failed: %v", err)
            return
        }

        close(pullServerReady)

        for {
            if serverMsg, err = pullServer.RecvMsg(); err != nil {
                fmt.Printf("\nServer receive failed: %v", err)
            }

            fmt.Println("Server received: ", string(serverMsg.Body))

            if err = pullServer.SendMsg(serverMsg); err != nil {
                fmt.Printf("\nServer send failed: %v", err)
                return
            }
        }
    }()

    // 3. Create Push Client
    pushSocket, err := push.NewSocket()
    defer pushSocket.Close()
    all.AddTransports(pushSocket)

    // 4. Client dials Server
    if err = pushSocket.Dial(url); err != nil {
        fmt.Printf("\nClient dial failed: %v", err)
        return
    }
    <-pullServerReady

    // 5. Client sends message
    clientMessageBytes := []byte("Hello!")

    if err = pushSocket.Send(clientMessageBytes); err != nil {
        fmt.Printf("\nClient send failed: %v", err)
        return
    }

    fmt.Printf("\nClient sending: %s\n", clientMessageBytes)
}
