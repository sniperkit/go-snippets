package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"git.andrewcsellers.com/acsellers/card_sharp/config"
	"git.andrewcsellers.com/acsellers/card_sharp/lobby"
	"github.com/acsellers/multitemplate"
	"github.com/acsellers/platform/controllers"
	"github.com/acsellers/platform/router"
	"github.com/acsellers/words"
	"golang.org/x/net/websocket"
)

func main() {
	go WebsocketServer()
	log.Fatal(http.ListenAndServe(config.WebPort(), Router()))
}

func Router() http.Handler {
	r := router.NewRouter()
	r.Many(PageCtrl{NewRenderableCtrl("desktop.html")})
	r.Many(FrontGameCtrl{NewRenderableCtrl("desktop.html"), nil})
	r.One(PlayerCtrl{RenderableCtrl: NewRenderableCtrl("player.html")})
	r.Mount(controllers.AssetModule{
		AssetLocation: "public",
	})

	return r
}

type RenderableCtrl struct {
	*router.BaseController
	Template, Layout string
}

func NewRenderableCtrl(layout string) RenderableCtrl {
	return RenderableCtrl{
		&router.BaseController{},
		"", layout,
	}
}

var m = sync.Mutex{}

func (rc RenderableCtrl) Render() router.Result {
	if *config.Dev {
		m.Lock()
		config.CompileTemplates()
	}

	ctx := &multitemplate.Context{
		Main:   rc.Template,
		Layout: rc.Layout,
		Dot:    rc.Context,
	}
	buf := &bytes.Buffer{}
	err := config.Tmpl.ExecuteContext(buf, ctx)
	rc.Log.Println(err)
	if *config.Dev {
		m.Unlock()
	}
	if err != nil {
		return router.InternalError{err}
	} else {
		return router.Rendered{Content: buf}
	}
}

type PageCtrl struct {
	RenderableCtrl
}

func (PageCtrl) Path() string {
	return ""
}

func (pc PageCtrl) Index() router.Result {
	pc.Template = "front.html"
	return pc.Render()
}

func (pc PageCtrl) Show() router.Result {
	pc.Template = "front.html"
	return pc.Render()
}

type FrontGameCtrl struct {
	RenderableCtrl
	*lobby.Lobby
}

func (FrontGameCtrl) Path() string {
	return "games"
}

func (fgc *FrontGameCtrl) PreItem() router.Result {
	var value string
	if cookie, err := fgc.Request.Cookie("party-lobby"); err != nil {
		fgc.Log.Printf("Could not retrieve Cookie: %s\n", err.Error())
		return router.NotAllowed{
			Request:  fgc.Request,
			Fallback: "/games/new",
		}
	} else {
		value = cookie.Value
	}
	var lobbyid string
	if err := config.Cookie.Decode("party-lobby", value, &lobbyid); err != nil {
		fgc.Log.Printf("Could not decode Cookie: %s\n", err.Error())
		return router.NotAllowed{
			Request:  fgc.Request,
			Fallback: "/games/new",
		}
	}
	l := lobby.Find(lobbyid)
	if l == nil {
		fgc.Log.Println("Could not retrieve lobby from cookie")
		return router.NotAllowed{
			Request:  fgc.Request,
			Fallback: "/games/new",
		}
		return nil
	}
	fgc.Lobby = l
	fgc.Context["Lobby"] = l
	return nil
}
func (fgc FrontGameCtrl) New() router.Result {
	fgc.Template = "new_game_lobby.html"
	return fgc.Render()

}

func (fgc FrontGameCtrl) Create() router.Result {
	fgc.Request.ParseForm()
	gid := fgc.Request.Form.Get("game_id")
	if gid == "" {
		return router.Redirect{
			Request: fgc.Request,
			URL:     "/games/new",
		}
	}

	id, err := strconv.Atoi(gid)
	if err != nil {
		return router.Redirect{
			Request: fgc.Request,
			URL:     "/games/new",
		}
	}

	g, err := config.Conn.Deck.Find(id)
	if err != nil {
		return router.Redirect{
			Request: fgc.Request,
			URL:     "/games/new",
		}
	}
	l := lobby.Create(g)
	en, err := config.Cookie.Encode("party-lobby", l.ID)
	if err != nil {
		return router.Redirect{
			Request: fgc.Request,
			URL:     "/games/new",
		}
	}

	http.SetCookie(fgc.Out, &http.Cookie{Name: "party-lobby", Value: en, Path: "/"})

	fgc.SetDefaultPlayer()
	return router.Redirect{
		Request: fgc.Request,
		URL:     "/games/" + l.ID,
	}
}

func (rc RenderableCtrl) SetDefaultPlayer() {
	p := lobby.CreatePlayer(words.QuickDouble())
	en, _ := config.Cookie.Encode("player-info", p.ID)
	http.SetCookie(rc.Out, &http.Cookie{Name: "player-info", Value: en, Path: "/"})
}

func (fgc FrontGameCtrl) Show() router.Result {
	fgc.Log.Println(fgc.Lobby)
	for _, p := range fgc.Lobby.Players {
		fgc.Log.Println(p)
	}
	fgc.Layout = "lobby.html"
	fgc.Template = "game_lobby.html"
	return fgc.Render()

}
func (fgc FrontGameCtrl) Edit() router.Result {
	fgc.Template = "edit_game_lobby.html"
	return fgc.Render()
}

func (fgc FrontGameCtrl) Update() router.Result {
	return router.Redirect{
		Request: fgc.Request,
		URL:     "/",
	}
}

func (fgc FrontGameCtrl) Join() router.Result {
	fgc.Template = "join_game.html"
	return fgc.Render()
}

func (fgc FrontGameCtrl) DoJoin() router.Result {
	fgc.Template = "join_game.html"
	fgc.Request.ParseForm()
	gid := fgc.Request.Form.Get("game_code")
	l := lobby.Find(gid)
	if l != nil {
		en, err := config.Cookie.Encode("party-lobby", l.ID)
		if err != nil {
			return fgc.Render()
		}

		http.SetCookie(fgc.Out, &http.Cookie{Name: "party-lobby", Value: en, Path: "/"})
		fgc.SetDefaultPlayer()
		return router.Redirect{
			Request: fgc.Request,
			URL:     "/players/edit",
		}
	}

	return fgc.Render()
}

func (fgc FrontGameCtrl) OtherBase(sr *router.SubRoute) {
	sr.Get("join").Action("Join")
	sr.Post("join").Action("DoJoin")
}

type PlayerCtrl struct {
	RenderableCtrl
	Game   *lobby.Lobby
	Player *lobby.Player
}

func (PlayerCtrl) Path() string {
	return "players"
}

func (pc *PlayerCtrl) PreFilter() router.Result {
	var value string
	if cookie, err := pc.Request.Cookie("party-lobby"); err != nil {
		pc.Log.Printf("Could not retrieve Cookie: %s\n", err.Error())
		return router.NotAllowed{
			Request:  pc.Request,
			Fallback: "/games/join",
		}
	} else {
		value = cookie.Value
	}
	var lobbyid string
	if err := config.Cookie.Decode("party-lobby", value, &lobbyid); err != nil {
		pc.Log.Printf("Could not decode Cookie: %s\n", err.Error())
		return router.NotAllowed{
			Request:  pc.Request,
			Fallback: "/games/join",
		}
	}
	l := lobby.Find(lobbyid)
	if l == nil {
		pc.Log.Println("Could not retrieve lobby from cookie")
		return router.NotAllowed{
			Request:  pc.Request,
			Fallback: "/games/join",
		}
		return nil
	}
	pc.Game = l
	pc.Context["Lobby"] = l

	if cookie, err := pc.Request.Cookie("player-info"); err != nil {
		pc.Log.Printf("Could not retrieve Cookie: %s\n", err.Error())
		pc.SetDefaultPlayer()
	} else {
		value = cookie.Value
	}
	var playerid string
	if err := config.Cookie.Decode("player-info", value, &playerid); err != nil {
		pc.Log.Printf("Could not decode Cookie: %s\n", err.Error())
		pc.SetDefaultPlayer()
	}
	p := lobby.FindPlayer(playerid)
	if p == nil {
		pc.Log.Println("Could not retrieve lobby from cookie")
		pc.SetDefaultPlayer()
	}
	pc.Player = p
	pc.Context["Player"] = p

	return nil
}

func (pc PlayerCtrl) Edit() router.Result {
	pc.Template = "player_edit.html"
	return pc.Render()
}

func (pc PlayerCtrl) Update() router.Result {
	pc.Request.ParseForm()
	if pc.Request.Form.Get("name") != pc.Player.Name {
		pc.Player.Name = pc.Request.Form.Get("name")
	}
	pc.Game.Add(pc.Player)

	return router.Redirect{
		Request: pc.Request,
		URL:     "/players",
	}
}

func (pc PlayerCtrl) Show() router.Result {
	pc.Template = "player_show.html"
	if !strings.Contains(pc.Player.Status, "waiting") {
		pi := pc.Game.Instance.Players[pc.Player.ID]
		pc.Context["Hand"] = pc.Game.Instance.Hands[pi]
	}
	if pc.Player.Status == "judge" {
		pc.Context["Judging"] = pc.Game.Instance.CurrentPlays
	}

	return pc.Render()
}

func (pc PlayerCtrl) OtherItem(sr *router.SubRoute) {
	sr.Post("start_game").Action("StartGame")
	sr.Post("make_move").Action("MakeMove")
	sr.Post("pick_card").Action("PickCard")
}

func (pc PlayerCtrl) PickCard() router.Result {
	if pc.Game.Players[pc.Game.Instance.CurrentJudge] != pc.Player {
		return nil
	}
	pc.Request.ParseForm()
	pc.Game.PickCard(pc.Request.Form.Get("pid"))
	return router.String{Content: "ok"}
}
func (pc PlayerCtrl) StartGame() router.Result {
	if pc.Game.Czar == pc.Player {
		pc.Game.GameStart()
	}
	return router.String{Content: "ok"}
}
func (pc PlayerCtrl) MakeMove() router.Result {
	pc.Request.ParseForm()
	fmt.Println(pc.Request.Form)
	card := pc.Request.Form.Get("card")
	if card == "" {
		pc.Game.Sync(pc.Player)
		return router.String{Content: "resync"}
	}

	pc.Game.PlayCard(pc.Player, strings.TrimSpace(card))

	return router.String{Content: "ok"}
}

func PresenterServer(ws *websocket.Conn) {
	p := strings.Split(ws.LocalAddr().String(), "/")
	l := p[len(p)-1]
	wl := lobby.Find(l)
	fmt.Println(wl)
	msg := make(chan []byte)
	if wl.Presenter != nil {
		close(wl.Presenter)
	}
	wl.Presenter = msg
	for m := range msg {
		fmt.Println(ws.Write(m))
	}
}

func PlayerServer(ws *websocket.Conn) {
	p := strings.Split(ws.LocalAddr().String(), "/")
	l := p[len(p)-2]
	wl := lobby.Find(l)
	y := p[len(p)-1]
	wp := lobby.FindPlayer(y)
	msg := make(chan []byte)
	if wl.PlayerChan[wp.ID] != nil {
		close(wl.PlayerChan[wp.ID])
	}
	wl.PlayerChan[wp.ID] = msg
	for m := range msg {
		fmt.Println(ws.Write(m))
	}

}

func WebsocketServer() {
	http.Handle("/player/", websocket.Handler(PlayerServer))
	http.Handle("/presenter/", websocket.Handler(PresenterServer))
	log.Fatal(http.ListenAndServe(":8101", nil))
}
