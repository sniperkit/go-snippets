package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

var mutex = &sync.Mutex{}

// right now, it cannot handle spacebar and other languages, need to extend it later
var wordList = []string{"cat", "dog", "max", "delicious", "games", "ant", "min", "computer", "macbook", "microsoft", "apple", "pineapple",
	"genesis", "Thailand", "Google", "food", "Golang", "gopher", "httpserver", "cassandra", "zookeeper", "hadoop", "combination", "algorithm",
	"wonderland", "Game of Thrones"}

type Word struct {
	x, y         int
	str          string
	velo         int
	progress     int
	typedKeyChan chan rune
}

type Player struct {
	score int
	hp    int
	mp    int
}

type Screen struct {
	width  int
	height int
	top    int
	left   int
	right  int
	bottom int
}

type TypingGeeks struct {
	eventChan     chan termbox.Event
	exitChan      chan string
	wordMap       map[rune]Word
	wordChan      chan Word
	typedKeyChan  chan rune
	curWordKey    rune
	screen        Screen
	rowSize       int
	colSize       int
	wordVeloRange int
	wordVeloBase  int
	wordFps       int
	fps           int // frame per sec, must be between 1-60, default is 25
	player1       Player
}

func (t *TypingGeeks) Initialise() {
	// prepare rand engine
	rand.Seed(time.Now().UnixNano())

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	// initialise all go channels
	t.eventChan = make(chan termbox.Event)
	t.exitChan = make(chan string)
	t.wordMap = make(map[rune]Word)
	t.wordChan = make(chan Word)
	t.typedKeyChan = make(chan rune)
	// screen related
	t.screen = Screen{
		top:    2,
		left:   5,
		right:  5,
		bottom: 0,
	}
	t.screen.width, t.screen.height = termbox.Size()
	t.colSize = t.screen.width - (t.screen.left + t.screen.right)
	t.rowSize = t.screen.height - (t.screen.top + t.screen.bottom)
	t.fps = 25
	// initialise player
	t.player1 = Player{
		hp:    10,
		mp:    0,
		score: 0,
	}
	// initialise word velocities
	t.wordVeloBase = 1
	t.wordVeloRange = 1
	t.wordFps = 1
}

func (t *TypingGeeks) WaitExit() {
	<-t.exitChan
	termbox.Close()
}

func (t *TypingGeeks) GoWordFeeder() {
	/*	wordFile, err := os.Open("words.txt")
		if err != nil {
			panic(err)
		}
	*/
	counter := 0
	for {
		// set level (word velocities, pop-up frequency)
		t.wordVeloRange = t.player1.score/50*2 + 1
		t.wordVeloBase = t.player1.score/100 + 1
		t.wordFps = t.player1.score/50*2 + 1

		t.wordChan <- Word{
			y:            0,
			x:            rand.Intn(t.colSize),
			str:          wordList[counter],
			velo:         t.wordVeloBase + rand.Intn(t.wordVeloRange),
			progress:     0,
			typedKeyChan: make(chan rune),
		}
		wordFpsSleepTime := time.Duration(1000000/t.wordFps) * time.Microsecond
		time.Sleep(wordFpsSleepTime)
		counter++
		if counter >= len(wordList) {
			counter = 0
		}
	}
}

func (t *TypingGeeks) GoMainProcessor() {
	for {
		select {
		case newWord := <-t.wordChan: // receive new word and spawn a new go word routine
			// add newWord to wordPool for rendering
			key, _ := utf8.DecodeRuneInString(newWord.str)
			if _, exist := t.wordMap[key]; exist {
				continue
			}
			t.wordMap[key] = newWord
			// spawn go routine for each word to process itself (moving)
			go func(key rune) {
				veloSleepTime := time.Duration(2000000/newWord.velo) * time.Microsecond
				leftSleepTime := veloSleepTime
				startTime := time.Now()
				for {
					select {
					case typedKey := <-t.wordMap[key].typedKeyChan:
						curWord, exist := t.wordMap[key]
						if !exist {
							return
						}
						for pos, char := range curWord.str {
							if pos == curWord.progress {
								if char == typedKey {
									curWord.progress++
									if curWord.progress >= len(curWord.str) {
										// TODO: finish whole word, implement successful attempt effect
										t.player1.score += len(curWord.str)
										//  destroy word
										delete(t.wordMap, t.curWordKey)
										t.curWordKey = 0
										return
									}
									t.wordMap[t.curWordKey] = curWord
									break
								} else {
									// TODO: wrong key, implement fail attempt effect
									t.decreasePlayer1HP(1)
								}
							}
						}
						leftSleepTime = veloSleepTime - time.Now().Sub(startTime)
					case <-time.After(leftSleepTime):
						// due to issue #3117, we gotta assign value like this for map of struct for now.
						curWord, exist := t.wordMap[key]
						if !exist {
							return
						}
						curWord.y++
						t.wordMap[key] = curWord
						// delete word that goes out of windows in wordMap
						// TODO: need to watch out of race condition for map, too. see -> https://blog.golang.org/go-maps-in-action
						if t.wordMap[key].y > t.rowSize {
							// TODO: fail to finish word, implement fail attempt effect
							t.decreasePlayer1HP(1)
							// destroy word
							delete(t.wordMap, key)
							return
						}
						// set startTime over again
						startTime = time.Now()
						leftSleepTime = veloSleepTime
					}
				}
			}(key)
		case <-time.After(500 * time.Millisecond):
			//fmt.Println("timeout")
		}
	}
}

func (t *TypingGeeks) decreasePlayer1HP(val int) {
	t.player1.hp -= val
	if t.player1.hp <= 0 {
		t.exitChan <- "Game Over"
	}
}

func (t *TypingGeeks) drawWord(x, y int, word string, fg, bg termbox.Attribute) {
	for pos, char := range word {
		termbox.SetCell(x+pos, y, char, fg, bg)
	}
}

func (t *TypingGeeks) drawTitle() {
	title := "|     TypingGeeks     |"
	t.drawWord(t.screen.width/2-len(title)/2-3, 0, title, termbox.ColorWhite, termbox.ColorBlack)
}

func (t *TypingGeeks) drawPlayerStatus() {
	t.drawWord(t.screen.width-30, 0, "Score: "+strconv.Itoa(t.player1.score), termbox.ColorWhite, termbox.ColorBlack)
	t.drawWord(15, 0, "HP: "+strconv.Itoa(t.player1.hp), termbox.ColorWhite, termbox.ColorBlack)
	t.drawWord(30, 0, "MP: "+strconv.Itoa(t.player1.mp), termbox.ColorWhite, termbox.ColorBlack)

}

func (t *TypingGeeks) drawNavBarLine(y int) {
	for x := 0; x < t.screen.width; x++ {
		termbox.SetCell(x, y, '-', termbox.ColorWhite, termbox.ColorBlack)
	}
}

func (t *TypingGeeks) GoRender() {
	for {
		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		// draw Title
		t.drawTitle()
		// render navbar showing score, hp, mp(future)
		t.drawPlayerStatus()
		t.drawNavBarLine(1)
		// render all words
		for _, word := range t.wordMap {
			for pos, char := range word.str {
				if pos >= word.progress {
					// render word by also having border from top and left
					termbox.SetCell(t.screen.left+word.x+pos, t.screen.top+word.y, char, termbox.ColorWhite, termbox.ColorBlack)
				}
			}
		}
		termbox.Flush()
		fpsSleepTime := time.Duration(1000000/t.fps) * time.Microsecond
		time.Sleep(fpsSleepTime)
	}

}

func (t *TypingGeeks) GoKeyAnalyzer() {
	for {
		select {
		case typedKey := <-t.typedKeyChan:
			if t.curWordKey != 0 {
				curWord := t.wordMap[t.curWordKey]
				curWord.typedKeyChan <- typedKey
			} else {
				// due to issue #3117, we gotta assign value like this for map of struct for now.
				if curWord, exist := t.wordMap[typedKey]; exist {
					t.curWordKey = typedKey
					curWord.typedKeyChan <- typedKey
				} else {
					// TODO: wrong key, no word exist, implement fail attempt effect
					t.decreasePlayer1HP(1)
				}
			}
		case <-time.After(500 * time.Millisecond):
		}
	}
}

func (t *TypingGeeks) GoEventTrigger() {
	for {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			switch event.Key {
			case termbox.KeyCtrlC, termbox.KeyCtrlX:
				// exit
				t.exitChan <- "end"
			case termbox.KeyEsc:
				// esc to cancel typing current word
				if t.curWordKey != 0 {
					tmp := t.wordMap[t.curWordKey]
					tmp.progress = 0
					t.wordMap[t.curWordKey] = tmp
					t.curWordKey = 0
				}
			}
			// TODO: this is blocking if you type fast, maybe use select? or use slice of channel?
			t.typedKeyChan <- event.Ch
		case termbox.EventResize:
			fmt.Println("Resize")
		case termbox.EventError:
			panic(event.Err)
		}
	}
}

func main() {
	t := new(TypingGeeks)
	t.Initialise()
	go t.GoEventTrigger()
	go t.GoKeyAnalyzer()
	go t.GoRender()
	go t.GoMainProcessor()
	go t.GoWordFeeder()
	t.WaitExit()
}
