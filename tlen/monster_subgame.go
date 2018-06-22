package main

import "github.com/seletskiy/go-android-rpc/android"

type MonsterSubgame struct {
	Subgame

	GlobalLocations map[string]Location
}

func (subgame MonsterSubgame) GetButtonTitle() string {
	return "Monster Title"
}

func (subgame MonsterSubgame) GetLayoutName() string {
	return "main_layout"
}

func (subgame MonsterSubgame) Enter(state *State) {
	defer android.PanicHandler()

	subgame.Subgame.Enter(state)

	monsterQuestion := `<p>"Alive? Energy? Energy."</p>
<p>You start to feel yourself more clearly.</p>
<p>Start to feel something around you. And you realize the something is alive.</p>
<p>You don't need another eternity to understand that the something is deadly strange, deadly dark and deadly dangerous. You feel it doesn't even know who are you, but even more you know it want you not to be.</p>`

	locations := map[string]Location{
		"1": &BaseLocation{
			Description: `<p><b>It is pitch black.</b></p>
			<p>You try to make a move to escape, but barely feel your body.</p>
			<p>Just powerful clammy darkness and strange feeling of existence and non-existence at the same time.</p>`,
		},

		"2": &BaseLocation{
			ButtonTitle: "Try to move",
			Description: `<p>You feel different. Time goes different.</p><p>You feel your mind wandering.</p>`,
		},

		"3": &BaseLocation{
			ButtonTitle: `"Who are you?"`,
			Description: `<p>"You? Things? Cold. Solid. Energy. Dark."</p>`,
		},

		"4.1": &BaseLocation{
			ButtonTitle: "Play.",
			Description: monsterQuestion,
		},

		"4.2": &BaseLocation{
			ButtonTitle: "Toys.",
			Description: monsterQuestion,
		},

		"4.3": &BaseLocation{
			ButtonTitle: "Mother.",
			Description: monsterQuestion,
		},

		"4.4": &BaseLocation{
			ButtonTitle: "Father.",
			Description: monsterQuestion,
		},

		"5": &BaseLocation{
			ButtonTitle: "Run away",
			Description: "<p>It is pitch black. There is nowhere to run.</p>",
		},

		"6.1": &BaseLocation{
			ButtonTitle: "Anger",
			Description: "<p>You feel dark presence goes away.</p>",
		},

		"6.2": &BaseLocation{
			ButtonTitle: "Panic",
			Description: `<p>You feel <font color="red">pain</font>.</p>`,
		},

		"7": &BaseLocation{
			ButtonTitle: "Anger",
			Description: `<p>You feel dark presence goes away.</p>
<p>You suddenly became aware of your body.</p>
<p>Rumbling around you. Low and pleasant rumbling. G-r-r-r-r. Like a tiger. Like a... an engines. Ginormous photon engines.</p>
<p>You're waking up.</p>
<p>You're seing you father. You're feeling bright <font color="red">pain</font>.</p>`,
		},

		"8": &BaseLocation{
			ButtonTitle: "Cry",
			Description: `<p>Your father is looking at you inquiringly.</p>
<p>"It's a just bad dream.", he says.</p>`,
		},

		"9": &BaseLocation{
			ButtonTitle: `I've fought with a monster!`,
			Description: `<p>Your father puts you on your feet.</p>`,
		},

		"10": &BaseLocation{
			ButtonTitle: `"I feel pain!"`,
			Description: `<p>"You'd better see your mom. She's our lovely doctor."</p>`,
		},

		"11": &BaseLocation{
			ButtonTitle: `"Are we <b>home</b>?"`,
			Description: `<p>"Not yet. But it looks like we are so close".</p>
<p>Grin touches your father's lips.</p>`,
		},

		"12": &BaseLocation{
			ButtonTitle: `End`,
			Description: `<p>Father takes your away from Imaginarium and locks the door on <b>magnet lock.</b></p>`,
		},
	}

	locations["1"].Link(locations["2"])

	locations["2"].Link(locations["3"])

	locations["3"].Link(locations["4.1"])
	locations["3"].Link(locations["4.2"])
	locations["3"].Link(locations["4.3"])
	locations["3"].Link(locations["4.4"])

	locations["4.1"].Link(locations["5"])
	locations["4.2"].Link(locations["5"])
	locations["4.3"].Link(locations["5"])
	locations["4.4"].Link(locations["5"])
	locations["4.4"].Link(locations["5"])

	locations["5"].Link(locations["6.1"])
	locations["5"].Link(locations["6.2"])

	locations["6.1"].Link(locations["7"])
	locations["6.2"].Link(locations["7"])

	locations["7"].Link(locations["8"])

	locations["8"].Link(locations["9"])

	locations["9"].Link(locations["10"])

	locations["10"].Link(locations["11"])

	locations["11"].Link(locations["12"])

	locations["12"].Link(subgame.GlobalLocations["bunk"])

	subgame.SetLayoutName("main_layout")
	subgame.SetLocation(locations["1"])

	subgame.Start()
}
