package main

import (
	"github.com/seletskiy/go-android-rpc/android"

	// required for linking
	_ "github.com/seletskiy/go-android-rpc/android/modules"
)

const (
	ViewInvisible = 4
	ViewVisible   = 0
	ViewGone      = 8
)

const Cat = `<font color="green">春</font>`

func initLocations() map[string]Location {
	locations := make(map[string]Location)

	githubMessage := `<p>Visit our github repo to get the latest version.</p>
<p>http://github.com/seletskiy/tlen</p>`

	locations["github"] = &GithubSubgame{}

	locations["bunk"] = NewTestLocation(
		"Go to bunk",
		"<p>It's yours family bunk.</p>"+
		"<p>You see "+Cat+".</p>"+
			"<p>You see your dad watching TV.</p>"+
			"<p>There are door to ship lobby.</p>",
	)

	locations["lobby"] = NewTestLocation("Go to lobby", "You enter the lobby")
	locations["din"] = NewTestLocation("Go to dinnary", "Dinnay")
	locations["med"] = NewTestLocation("Go to medical", "Medical cabinet")
	locations["cap"] = NewTestLocation("Go to captain room", githubMessage)
	// "Captain's room")

	locations["stub"] = &StubSubgame{
		GlobalLocations: locations,
	}

	locations["woman1"] = &Woman1Subgame{
		GlobalLocations: locations,
	}

	locations["mother1"] = &Mother1Subgame{
		GlobalLocations: locations,
	}

	locations["captain1"] = &Captain1Subgame{
		GlobalLocations: locations,
	}

	locations["cook1"] = &Cook1Subgame{
		GlobalLocations: locations,
	}

	locations["monster"] = &MonsterSubgame{
		GlobalLocations: locations,
	}

	locations["caress_cat"] = &CaressCatSubgame{
		OriginLocation: locations["bunk"],
	}

	locations["locked"] = &LockedSubgame{
		GlobalLocations: locations,
	}

	locations["wakeup_father_subgame"] = &WakeUpFatherSubgame{
		NextLocation: locations["github"],
	}

	locations["bunk"].Link(locations["caress_cat"])
	locations["bunk"].Link(locations["wakeup_father_subgame"])
	locations["bunk"].Link(locations["lobby"])

	locations["lobby"].Link(locations["bunk"])
	//locations["lobby"].Link(locations["din"])
	//locations["lobby"].Link(locations["med"])
	locations["lobby"].Link(locations["locked"])

	//locations["med"].Link(locations["mother1"])
	//locations["med"].Link(locations["captain1"])
	//locations["med"].Link(locations["lobby"])

	//locations["din"].Link(actions["cook1"])
	//locations["din"].Link(locations["lobby"])

	locations["locked"].Link(locations["cap"])
	locations["cap"].Link(locations["lobby"])

	// go to lobby, if not to github
	locations["github"].Link(locations["lobby"])
	//locations["wakeup_father_subgame"].Link(locations["github"])
	locations["lobby"].Link(locations["github"])
	//locations["bunk"].Link(locations["father2"])

	//locations["lobby"].Link(locations["woman2"])

	//locations["med"].Link(locations["pick_lock"])

	//locations["din"].Link(locations["mother2"])
	//locations["din"].Link(locations["cook2"])
	//if !game.IsCaptainAbused()
	//locations["bunk"].Link(locations["father3"])
	//}

	//if game.ScalpelIsStolen {
	//locations["lobby"].Link(locations["woman3"])
	//}

	//locations["med"].Link(locations["mother3"])
	//locations["din"].Link(locations["cook3"])
	//locations["cap"].Link(locations["explore_captain_room"])

	return locations
}

func start() {
	locations := initLocations()

	game := NewGame(&State{})
	game.SetLocation(locations["monster"])
	game.SetLayoutName("main_layout")
	game.Start()
}

func main() {
	defer android.PanicHandler()

	android.OnStart(start)
	android.Enter()
}
