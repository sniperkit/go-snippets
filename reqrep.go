package main

import (
    "fmt"
    "bitbucket.org/gdamore/mangos"
    "bitbucket.org/gdamore/mangos/protocol/req"
    "bitbucket.org/gdamore/mangos/protocol/rep"
    "bitbucket.org/gdamore/mangos/transport/all"
)

func main() {
    url := "tcp://127.0.0.1:8000"

    // 1. Create Response Server
    responseServerReady := make(chan struct{})
    responseServer, err := rep.NewSocket()
    defer responseServer.Close()

    all.AddTransports(responseServer)

    // 2. Run Response Server
    go func() {
        var err error
        var serverMsg *mangos.Message

        if err = responseServer.Listen(url); err != nil {
            fmt.Printf("\nServer listen failed: %v", err)
            return
        }

        close(responseServerReady)

        for {
            if serverMsg, err = responseServer.RecvMsg(); err != nil {
                fmt.Printf("\nServer receive failed: %v", err)
            }

            fmt.Println("serverMsg: %v", serverMsg.Body)

            serverMsg.Body = append(serverMsg.Body, []byte(" World!")...)

            fmt.Println("new serverMsg: %v", serverMsg.Body)

            if err = responseServer.SendMsg(serverMsg); err != nil {
                fmt.Printf("\nServer send failed: %v", err)
                return
            }
        }
    }()

    // 3. Create Request Client
    requestSocket, err := req.NewSocket()
    defer requestSocket.Close()
    all.AddTransports(requestSocket)

    // 4. Client dials Server
    if err = requestSocket.Dial(url); err != nil {
        fmt.Printf("\nClient dial failed: %v", err)
        return
    }
    <-responseServerReady

    // 5. Client sends message
    clientMessageBytes := []byte("Hello!")

    if err = requestSocket.Send(clientMessageBytes); err != nil {
        fmt.Printf("\nClient send failed: %v", err)
        return
    }

    var clientMsg []byte

    if clientMsg, err = requestSocket.Recv(); err != nil {
        fmt.Printf("\nClient receive failed: %v", err)
        return
    }

    fmt.Printf("\nClient receive message: %s\n", clientMsg)
}
