package main

import (
	"log"
	"runtime"
	"time"

	"net/http"
)
import "golang.org/x/net/websocket"

func main() {
	log.Println("Starting gogalamon server")

	runtime.GOMAXPROCS(runtime.NumCPU())

	http.Handle("/", http.FileServer(http.Dir("static/")))
	http.Handle("/sock/", websocket.Handler(wsHandler))

	go func() {
		var i int
		for {
			NextRenderId <- i
			i++
		}
	}()

	go mainLoop()
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}

var NewEntity = make(chan Entity)
var NewUser = make(chan *User)

const framesPerSecond = 30

func mainLoop() {
	var entities []Entity
	users := make(map[*User]struct{})
	overworld := NewOverworld()
	ticker := time.Tick(time.Second / framesPerSecond)
	ships := make(map[EntityShip]struct{})
	shipFrame := 0

	var planets []*Planet
	{
		go NewPlanet(2000, 4000, "planet_gas")
		go NewPlanet(6000, 0, "planet_gas")
		go NewPlanet(3000, -5000, "planet_gas")
		go NewPlanet(-7000, -1000, "planet_gas")
		go NewPlanet(-6000, -6000, "planet_gas")
	}

	for {
		select {
		case <-ticker:
			{
				var i, place int
				for i < len(entities) {
					entity := entities[i]
					if entity.update(overworld, planets) {
						entities[place] = entity
						place++
					} else {
						overworld.remove(entity)
						if ship, ok := entity.(EntityShip); ok {
							delete(ships, ship)
						}
					}
					i++
				}
				lastPlace := place
				for place < len(entities) {
					entities[place] = nil
					place++
				}
				entities = entities[:lastPlace]
			}
			{
				planetInfos := make([]PlanetInfo, len(planets))
				for i := range planetInfos {
					planetInfos[i] = planets[i].planetInfo()
				}
				shipFrame++
				var shipInfos []shipInfo
				if shipFrame > framesPerSecond*3 {
					shipInfos = make([]shipInfo, 0, len(ships))
					shipFrame = 0
					for ship := range ships {
						info := ship.shipInfo()
						if info.X*info.X+info.Y*info.Y <= 10000*10000 {
							shipInfos = append(shipInfos, info)
						}
					}
				}
				nextUsers := make(map[*User]struct{})
				wait := make(chan *User)
				for user := range users {
					go user.render(overworld, planetInfos, shipInfos, wait)

					if msg := user.GetChatMessage(); msg != nil {
						for other := range users {
							if other != user {
								other.Send(msg)
							}
						}
					}
				}
				for i := 0; i < len(users); i++ {
					user := <-wait
					if user != nil {
						nextUsers[user] = struct{}{}
					}
				}
				users = nextUsers
			}
			{
			}
		case entity := <-NewEntity:
			entities = append(entities, entity)
			if planet, ok := entity.(*Planet); ok {
				planets = append(planets, planet)
			}
			if ship, ok := entity.(EntityShip); ok {
				ships[ship] = struct{}{}
			}
		case user := <-NewUser:
			users[user] = struct{}{}
		}
	}
}

type Entity interface {
	update(overworld *Overworld, planets []*Planet) (alive bool)
	RenderInfo() RenderInfo
}

type team uint

type EntityTeam interface {
	team() team
}

const (
	TeamNone = team(iota)
	TeamPirates
	TeamGophers
	TeamPythons
	TeamMax
)

func (t team) String() string {
	switch t {
	case TeamPirates:
		return "pirate"
	case TeamGophers:
		return "gopher"
	case TeamPythons:
		return "python"
	}
	return "pirate"
}

type EntityDamage interface {
	Entity
	damage(damage int, teamSource team)
}

type EntityShip interface {
	shipInfo() shipInfo
}

type shipInfo struct {
	X, Y float32
	Team string
}

type V2 [2]float32

var NextRenderId = make(chan int)

type RenderInfo struct {
	I int
	X float32
	Y float32
	R float32
	N string //name
}

type EntitySound interface {
	Entity
	sound() string
}
