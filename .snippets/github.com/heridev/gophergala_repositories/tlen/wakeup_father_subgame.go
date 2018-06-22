package main

import (
	"log"
	"math"

	"github.com/seletskiy/go-android-rpc/android"
	"github.com/seletskiy/go-android-rpc/android/sdk"
	"github.com/zazab/zhash"
)

type WakeUpFatherSubgame struct {
	Subgame
	descView     sdk.TextView
	talkLocation Location

	NextLocation Location
}

func (subgame WakeUpFatherSubgame) GetButtonTitle() string {
	return "Talk to father"
}

func (subgame WakeUpFatherSubgame) GetLayoutName() string {
	return "main_layout"
}

func (subgame WakeUpFatherSubgame) Enter(state *State) {
	defer android.PanicHandler()
	subgame.Subgame.Enter(state)

	ignoreDescription := `
<p>Dad is watching TV.</p>

<p>"Go see your mother. She needs to check you after hibernation."</p>

<p>Presenter continues to talk about all the possibilities that awaits you and 3'000 settlers on new planet.</p>`

	locations := map[string]Location{
		"1": &BaseLocation{
			Description: `<p>Dad wakes up. He look at you with empty eyes.</p>`,
		},

		"2": &BaseLocation{
			ButtonTitle: `"Hey, dad! The monster was sooo scary!"`,
			Description: `<p>"Huh?"</p>

<p>Dad turns on TV.</p>
<p>You hear how presenter of the TV program tells you about ` +
				`new planet, detected at the edge of the ship's radars range.</p>`,
		},
		"3.1": &BaseLocation{
			ButtonTitle: `"It tried to hurt me! But I was brave."`,
			Description: ignoreDescription,
		},
		"3.2": &BaseLocation{
			ButtonTitle: `"I was so scared!"`,
			Description: ignoreDescription,
		},
		"4": &BaseLocation{
			ButtonTitle: `Ignore TV and dad`,
		},
		"5.1": &BaseLocation{
			ButtonTitle: `Caress ` + Cat,
		},
		"5.2": &BaseLocation{
			ButtonTitle: `Disregard ` + Cat,
		},
		"6": &FinalWakeUpLocation{
			BaseLocation{ButtonTitle: `Go outside`},
			&subgame,
		},
	}

	locations["1"].Link(locations["2"])

	locations["2"].Link(locations["3.1"])
	locations["2"].Link(locations["3.2"])

	locations["3.1"].Link(locations["4"])
	locations["3.2"].Link(locations["4"])

	locations["4"].Link(subgame.NextLocation)
	//locations["4"].Link(locations["5.2"])

	//locations["5.1"].Link(locations["6"])
	//locations["5.2"].Link(locations["6"])

	subgame.talkLocation = locations["1"]

	subgame.SetLayoutName("main_layout")
	subgame.SetLocation(&AwaitWakeUpLocation{
		game: &subgame,
	})

	subgame.Start()
}

func (subgame *WakeUpFatherSubgame) Talk() {
	defer android.PanicHandler()

	subgame.SetLocation(subgame.talkLocation)
	subgame.SwitchLocation()
}

type AwaitWakeUpLocation struct {
	BaseLocation
	game          *WakeUpFatherSubgame
	minX          float64
	minY          float64
	minZ          float64
	maxX          float64
	maxY          float64
	maxZ          float64
	tick          int64
	wakeUpSuccess bool
}

func (location AwaitWakeUpLocation) Enter(state *State) {
	log.Printf("%#v", "await wakeuasdasd")
	desc := location.game.CreateView("android.widget.TextView").(sdk.TextView)
	desc.SetText1s("YOU NEED SHAKE YOUR FATHER TO WAKE HIM UP")
	desc.SetTextSize(50.0)
	location.game.descView = desc
	location.game.AttachView(desc.View)

	sensors := zhash.HashFromMap(android.GetSensorsList())
	accelerometerId, err := sensors.GetString(
		"sensors", "TYPE_ACCELEROMETER",
	)
	if err != nil {
		panic(err)
	}

	android.SubscribeToSensorValues(
		accelerometerId,
		&location,
	)
}

type FinalWakeUpLocation struct {
	BaseLocation
	game *WakeUpFatherSubgame
}

func (location FinalWakeUpLocation) Enter(state *State) {
	log.Printf("%#v", "final wake up locatin")
	location.game.SetLocation(location.game.NextLocation)
	location.game.SwitchLocation()
}

func (location *AwaitWakeUpLocation) OnChange(values []float64) {
	if location.wakeUpSuccess {
		return
	}

	x := values[0]
	y := values[1]
	z := values[2]

	location.tick++

	if location.tick > 5 {
		location.maxX = 0
		location.maxY = 0
		location.maxZ = 0
		location.minX = 0
		location.maxY = 0
		location.minZ = 0

		location.tick = 0
	}

	if x < location.minX {
		location.minX = x
	}
	if y < location.minY {
		location.minY = y
	}
	if z < location.minZ {
		location.minZ = z
	}

	if x > location.maxX {
		location.maxX = x
	}
	if y > location.maxY {
		location.maxY = y
	}
	if z > location.maxZ {
		location.maxZ = z
	}

	log.Printf("wakeup_father_subgame.go:197 %#v", location)

	if math.Abs(location.maxX)+math.Abs(location.minX) > 20 ||
		math.Abs(location.maxY)+math.Abs(location.minY) > 20 ||
		math.Abs(location.maxZ)+math.Abs(location.minZ) > 20 {
		//log.Printf("%#v\n", "YOU WIN")
		location.wakeUpSuccess = true
		location.game.Talk()
	}

	//log.Printf("%s %s %s | %s %s %s\n",
	//    location.minX, location.minY, location.minZ,
	//    location.maxX, location.maxY, location.maxZ)
}

func (handler AwaitWakeUpLocation) OnAccuracyChange() {
	//noop
}
