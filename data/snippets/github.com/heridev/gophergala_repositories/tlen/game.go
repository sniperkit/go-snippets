package main

import (
	"log"
	"strconv"

	"github.com/seletskiy/go-android-rpc/android"
	"github.com/seletskiy/go-android-rpc/android/sdk"
)

type State struct {
	Cat      int
	Progress int

	ViewId      int
	Location    Location
	LayoutId    string
	LayoutName  string
	MoveCounter int
}

type Game struct {
	state      *State
	headerView sdk.TextView
	descView   sdk.TextView

	//viewId      int
	viewObjects map[string][]sdk.View

	LocationOnClickHandler GameLocationOnClickHandler
}

func NewGame(state *State) *Game {
	return &Game{
		state: state,
	}
}

func (game *Game) IncrementMoveCounter() {
	game.state.MoveCounter += 1
}

func (game *Game) SetLocation(location Location) {
	if location == nil {
		log.Printf("%#v", "!!!!LOCATION IS NIL")
	}
	game.state.Location = location
}

func (game *Game) SetLayoutName(layoutName string) {
	layoutResponse := android.GetLayoutById(layoutName)
	game.state.LayoutId = layoutResponse["layout_id"].(string)
	game.state.LayoutName = layoutName
}

func (game *Game) SwitchLayout() {
	log.Printf("%#v", game.state.LayoutName)
	android.ChangeLayout(game.state.LayoutName)
}

func (game *Game) Start() {
	game.headerView = android.GetViewById("header_text").(sdk.TextView)
	game.descView = android.GetViewById("desc_text").(sdk.TextView)
	game.viewObjects = make(map[string][]sdk.View)

	game.SwitchLocation()
}

func (game *Game) ClearViews() {

	for _, view := range game.viewObjects[game.state.LayoutName] {
		android.RemoveView(view, game.state.LayoutId)
	}

	game.viewObjects[game.state.LayoutName] = []sdk.View{}
}

func (game *Game) SwitchLocation() {
	game.ClearViews()

	location := game.state.Location

	android.SetTextFromHtml(game.headerView, location.GetHeader())
	android.SetTextFromHtml(game.descView, location.GetDescription())

	//game.SetLocationArterfactsVisibility(true)

	linkedLocations := location.GetLinkedLocations()
	for _, linkedLocation := range linkedLocations {
		button := game.CreateView("android.widget.Button").(sdk.Button)
		android.SetTextFromHtml(button, linkedLocation.GetButtonTitle())

		android.OnClick(button,
			GameLocationOnClickHandler{
				game, button, linkedLocation,
			},
		)

		game.AttachView(button.View)
	}

	location.Enter(game.state)
}

func (game *Game) SetLocationArterfactsVisibility(should bool) {
	headerVisible, _ := game.headerView.IsShown()
	descVisible, _ := game.descView.IsShown()

	if should {
		if !headerVisible {
			game.headerView.SetVisibility(ViewVisible)
		}
		if !descVisible {
			game.descView.SetVisibility(ViewVisible)
		}

	} else {
		if headerVisible {
			game.headerView.SetVisibility(ViewGone)
		}

		if descVisible {
			game.descView.SetVisibility(ViewGone)
		}
	}
}

func (game *Game) LocationOnClick(button sdk.Button, location Location) {
	game.IncrementMoveCounter()
	game.SetLocation(location)
	game.SwitchLocation()
}

func (game *Game) CreateView(viewName string) interface{} {
	game.state.ViewId++
	id := strconv.Itoa(game.state.ViewId)
	created := android.CreateView(id, viewName)
	return created
}

func (game *Game) AttachView(view sdk.View) {
	game.viewObjects[game.state.LayoutName] = append(
		game.viewObjects[game.state.LayoutName], view)

	android.AttachView(view, game.state.LayoutId)
}

type GameLocationOnClickHandler struct {
	game     *Game
	button   sdk.Button
	location Location
}

func (handler GameLocationOnClickHandler) OnClick() {
	handler.game.LocationOnClick(handler.button, handler.location)
}
