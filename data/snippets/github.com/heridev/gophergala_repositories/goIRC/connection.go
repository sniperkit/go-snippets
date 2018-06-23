package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type ConnectionStatus int

const (
	SocketConnected ConnectionStatus = iota
	UserPassSent
	UserNickSent
	UserUserInfoSent
	UserRegistered
)

type User struct {
	Nick     string
	Ident    string
	RealName string
	Conn     net.Conn
	Status   ConnectionStatus
	Host     string
}

func (u *User) getHead() string {
	//return fmt.Sprintf(":%s!%s@%s", u.Nick, u.Ident, u.Host)
	return fmt.Sprintf(":%s!%s@127.0.0.1", u.Nick, u.Ident)
}

func handleConnection(conn net.Conn, buses map[string]*EventBus) {
	client := User{Status: UserPassSent, Conn: conn}
	//myIP := net.Conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)

	commands := make(map[string]func(map[string]*EventBus, *User, string, string))
	commands["JOIN"] = handleJoin
	commands["TOPIC"] = handleTopic
	commands["PRIVMSG"] = handleMsg
	commands["NICK"] = handleNick
	commands["PART"] = handlePart
	commands["HELP"] = handleHelp
	commands["LIST"] = handleList
	commands["PING"] = handlePing
	commands["PONG"] = handlePong

	for {
		status, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		status = strings.TrimSpace(status)
		statLen := strings.Split(status, " ")

		// allows user to enter empty strings
		if len(status) == 0 {
			conn.Write([]byte(""))
			continue
		} else if len(statLen) < 2 {
			cmd := strings.SplitN(status, " ", 1)
			cmd[0] = strings.ToUpper(cmd[0])
			if _, ok := commands[cmd[0]]; ok {
				commands[cmd[0]](buses, &client, "", "")
			}
		} else {
			if client.Status < UserRegistered {
				regCmd := strings.SplitN(status, " ", 2)
				regCmd[0] = strings.ToUpper(regCmd[0])
				fmt.Println("-" + regCmd[0] + "-" + regCmd[1])
				switch regCmd[0] {
				case "NICK":
					//client.Nick = regCmd[1]
					handleNick(buses, &client, regCmd[1], "")
					client.Status = UserNickSent
				case "USER":
					fmt.Println("hit user case")
					var uname, hname, sname, rname string
					fmt.Sscanf(regCmd[1], "%s %s %s :%s", uname, hname, sname, rname) //TODO(jz) need to split on : in case real name has spaces
					fmt.Println(hname + uname)
					//client.Ident = uname
					client.RealName = rname
					client.Status = UserRegistered
					client.Ident = client.Nick
					fmt.Println("username:" + client.Ident)
					buses[client.Nick] = &EventBus{subscribers: make(map[EventType][]Subscriber), channel: nil}
					buses[client.Nick].Subscribe(PrivMsg, &client)
					sendWelcome(&client)
				case "PASS": //need to remove this at some point!
					client.Nick = regCmd[1]
					client.Ident = regCmd[1]
					client.RealName = regCmd[1]

					buses[client.Nick] = &EventBus{subscribers: make(map[EventType][]Subscriber), channel: nil}
					buses[client.Nick].Subscribe(PrivMsg, &client)
					client.Status = UserRegistered
					sendWelcome(&client)

				//conn.Write([]byte("Welcome " + regCmd[1] + ")

				default:
					client.Write("you must register first. try nick or user?")
				}

			} else {
				// split <command> <target> :<data>

				var cmd, target, data string
				s := strings.SplitN(status, ":", 2)
				if len(s) > 1 {
					data = s[1]
				}
				_, err = fmt.Sscanf(s[0], "%s %s", &cmd, &target)
				if err != nil {
					fmt.Println(err)
					client.Write("Invalid input! CHECK YOUR(self) SYNTAX")
					continue
				}
				cmd = strings.ToUpper(cmd)
				if _, ok := commands[cmd]; ok {
					commands[cmd](buses, &client, target, data)
				}
			}
		}
	}
}

func checkEventBus(buses map[string]*EventBus, client *User, target string) (*EventBus, bool) {
	b, ok := buses[target]
	if !ok {
		client.Write(fmt.Sprintf(canned_responses[ERR_NOSUCHCHANNEL], client.Nick))
	}
	return b, ok
}
func checkSubscribed(bus *EventBus, client *User, event_type EventType) bool {
	for _, v := range bus.subscribers[event_type] {
		if v == client {
			return true
		}
	}
	return false
}
func isChannel(target string) bool {
	return string(target[0]) == "#"
}
func handlePart(buses map[string]*EventBus, client *User, target string, data string) {
	message := fmt.Sprintf("%s parted %s!\n", client.Nick, target)
	b, ok := checkEventBus(buses, client, target)
	if !ok {
		return
	}
	if ok = checkSubscribed(b, client, UserPart); !ok {
		return
	}
	buses[target].Publish(&Event{event_type: UserPart, event_data: message})
	delete(b.channel.mode, client.Nick)

	buses[target].Unsubscribe(UserPart, client)
	buses[target].Unsubscribe(UserJoin, client)
	buses[target].Unsubscribe(Topic, client)
	buses[target].Unsubscribe(PrivMsg, client)
	// possibile race condition
	if len(b.channel.mode) == 0 {
		delete(buses, target)
		fmt.Println(target + " closed")
	}
}

func handlePing(buses map[string]*EventBus, client *User, target string, data string) {
	client.Write("PONG " + target)
}

func handlePong(buses map[string]*EventBus, client *User, target string, data string) {
	//no op for fun
}

func handleJoin(buses map[string]*EventBus, client *User, target string, data string) {
	fmt.Println("!!!!!!!!! JOIN")
	if !isChannel(target) {
		return
	}
	b, ok := checkEventBus(buses, client, target)
	if !ok {
		newChannel := Channel{name: target, topic: "gogo new channel!", mode: make(map[string]Mode)}
		buses[newChannel.name] = &EventBus{subscribers: make(map[EventType][]Subscriber), channel: &newChannel}
		b = buses[newChannel.name]
	}
	if ok = checkSubscribed(b, client, UserJoin); !ok {
		b.channel.mode[client.Nick] = Voice
		b.Subscribe(UserJoin, client)
		b.Subscribe(UserPart, client)
		b.Subscribe(PrivMsg, client)
		b.Subscribe(Topic, client)
		//message := fmt.Sprintf("%s joined %s!\n", client.Nick, target)
		message := fmt.Sprintf("%s JOIN %s\n", client.getHead(), target)
		//send names
		var names string
		for _, val := range buses[target].subscribers[PrivMsg] {
			names = names + " " + val.GetInfo()
		}
		client.Write(":" + HOST_STRING + " 332 " + client.Nick + " " + target + ":no topic set")
		client.Write(":" + HOST_STRING + " 333 " + client.Nick + " " + target + " admin!admin@localhost 1419044230")
		client.Write(":" + HOST_STRING + " 353 " + client.Nick + " " + target + " :" + names)
		client.Write(":" + HOST_STRING + " 366 " + client.Nick + " * :END of /NAMES list.")
		///end send names
		b.Publish(&Event{UserJoin, message})
	}

}

func handleTopic(buses map[string]*EventBus, client *User, target string, data string) {
	b, ok := checkEventBus(buses, client, target)
	if !ok {
		return
	}
	if ok = checkSubscribed(b, client, Topic); !ok {
		return
	}

	if len(data) > 0 {
		b.channel.topic = data
		message := fmt.Sprintf("%s changed the channel topic to %s", client.Nick, data)
		b.Publish(&Event{Topic, message})
	} else {
		message := fmt.Sprintf("%s\n", b.channel.topic)
		client.Write(message)
	}
}

func handleNick(buses map[string]*EventBus, client *User, target string, data string) {
	client.Nick = target
	client.Write("nick set to:" + client.Nick)
}

func handleMsg(buses map[string]*EventBus, client *User, target string, data string) {
	b, ok := checkEventBus(buses, client, target)
	if !ok {
		return
	}
	if !isChannel(target) {
		message := fmt.Sprintf("%s PRIVMSG %s: %s\n", client.getHead(), target, data)
		b.Publish(&Event{event_type: PrivMsg, event_data: message})
		buses[client.Nick].Publish(&Event{event_type: PrivMsg, event_data: message})
	}
	if ok = checkSubscribed(b, client, PrivMsg); !ok {
		return
	}
	message := fmt.Sprintf("%s PRIVMSG %s: %s\n", client.getHead(), target, data)
	b.Publish(&Event{event_type: PrivMsg, event_data: message})
}

//func handleList(conn net.Conn, buses map[string]*EventBus) {
func handleList(buses map[string]*EventBus, client *User, target string, data string) {
	if len(buses) == 0 {
		client.Write("No Channels Exist")
	} else {
		client.Write("Channels")
		for k, _ := range buses {
			if k[:1] == "#" {
				client.Conn.Write([]byte(k + "\n"))
			}
		}
		client.Conn.Write([]byte("End of List\n"))
		//client.Write(k)
	}
}

func handleHelp(buses map[string]*EventBus, client *User, target string, data string) {
	k, ok := Help[target]
	if !ok {
		client.Write("\nAvailable Commands: (Enter HELP <command> for further details")
		for h := range Help {
			client.Write(h)
		}
	} else {
		client.Write("Summary: " + k.Summary + "\r\nUsage: " + k.Syntax)
	}
}
