package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
)

type LogEvent struct {
	Message  string
	Xp       int
	Life     int
	Speed    int
	Power    int
	Ancestry int
}

var gameLogMem = make([]string, 0)

type GameLog struct {
	recvEvents <-chan LogEvent
}

func (s *GameLog) InitLogEventStream(logEvents <-chan LogEvent) {
	s.recvEvents = logEvents

	go func() {
		for logEvent := range logEvents {
			s.storeLogEvent(logEvent)
		}
	}()
}

func NewGameLog() (*GameLog, error) {
	gameLog := &GameLog{
		recvEvents: make(chan LogEvent),
	}

	return gameLog, nil
}

func (s *GameLog) storeLogEvent(logEvent LogEvent) {

	var logString = fmt.Sprintf("%+v", logEvent)

	gameLogMem = append(gameLogMem, logString)
}

func (*GameLog) PrintGameLog() {
	ct.ChangeColor(ct.Yellow, true, ct.None, false)
	fmt.Println()
	fmt.Println("Game Log")
	fmt.Println("--------")
	for _, gameLog := range gameLogMem {
		fmt.Println(gameLog)
	}
	fmt.Println()
	ct.ResetColor()
}
