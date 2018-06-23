package theme

import (
	"encoding/json"
	"github.com/gophergala/cheppirc/message"
	"github.com/gophergala/cheppirc/user"
	"log"
	"sync"
	"strings"
)

type ThemeData struct {
	// Messages map[string][]message.Message
	Targets map[string]Target
	Users map[string]map[string]user.User
	Uuid string
	Nick string
	sync.RWMutex
}

type Target struct {
	Name string
	Type string
	Messages []message.Message
}

func (d *ThemeData) AddMessage(target, sender, text string, mtype string, updater chan []byte) {
	log.Println("ADDMESSAGE:", text, "DEBUG USERS:", d.Users)

	d.Lock()
	if _, ok := d.Targets[target]; !ok {
		log.Println("ADDMESSAGE: Target not found. Target:", target)
		d.Targets[target] = *NewTarget(target)
	}

	tempT := d.Targets[target]
	m := message.Message{sender, text, tempT.Name, mtype}
	tempT.AddMessage(m)
	d.Targets[target] = tempT
	d.Unlock()

	b, err := json.Marshal(m)
	if err != nil {
		log.Println("Error marshalling message:", err.Error())
	}
	log.Println(d.Targets)
	updater <- b
}

func (d *ThemeData) SetUsers(target, nick, info string) {
	d.Lock()
	if _, ok := d.Users[target]; !ok {
		log.Println("SETUSERS: Target not found. Target:", target)
		d.Users[target] = make(map[string]user.User)
	}
	d.Users[target][nick] = user.User{nick, info}
	d.Unlock()
}

func (t *Target) AddMessage(m message.Message) {
	t.Messages = append(t.Messages, m)
}

func NewThemeData() *ThemeData {
	d := &ThemeData{}
	//d.Messages = make(map[string][]message.Message)
	d.Targets = make(map[string]Target)
	d.Users = make(map[string]map[string]user.User)
	return d
}

func NewTarget(name string) *Target {
	var t Target
	if (name[0] == 35) {
		//If the target start with a # then it's a channel
		targetName := strings.Trim(name, "# ")
		t = Target{targetName, "channel", []message.Message{}}
	} else {
		t = Target{name, "other", []message.Message{}}
	}
	//t.Messages = []Message{}
	return &t
}
