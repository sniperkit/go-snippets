package main

import (
	"fmt"
	"strings"

	"encoding/json"
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
	"io/ioutil"
	"net/http"
	"time"
	"vueutil"
)

type GitCommit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			Sha string `json:"sha"`
			URL string `json:"url"`
		} `json:"tree"`
		URL          string `json:"url"`
		CommentCount int    `json:"comment_count"`
	} `json:"commit"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
	Author      struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Committer struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"committer"`
	Parents []struct {
		Sha     string `json:"sha"`
		URL     string `json:"url"`
		HTMLURL string `json:"html_url"`
	} `json:"parents"`
}

type Commit struct {
	*js.Object
	HtmlUrl          string `js:"html_url"`
	Sha              string `js:"sha"`
	CommitMessage    string `js:"commit_message"`
	AuthorHtmlUrl    string `js:"author_html_url"`
	CommitAuthorName string `js:"commit_author_name"`
	CommitAuthorDate string `js:"commit_author_date"`
}

type AppData struct {
	*js.Object
	Branches      []string  `js:"branches"`
	CurrentBranch string    `js:"currentBranch"`
	Commits       []*Commit `js:"commits"`
}

type AppProps struct {
	*js.Object
}

type App struct {
	Data  *AppData
	Props *AppProps
}

func NewAppData() interface{} {
	ad := &AppData{
		Object: js.Global.Get("Object").New(),
	}
	ad.Branches = []string{"master", "dev"}
	ad.CurrentBranch = "master"
	ad.Commits = []*Commit{}
	return ad
}

func NewApp(vm *vue.ViewModel) *App {
	return &App{
		Data: &AppData{
			Object: vm.Data,
		},
		Props: &AppProps{
			Object: vueutil.PropsData(vm),
		},
	}
}

func (app *App) SyncViewModel(vm *vue.ViewModel) {
	// [Vue warn]: Avoid replacing instance root $data. Use nested data properties instead.
	//vm.Data = t.Date.Object
	keys := js.Keys(app.Data.Object)
	for _, v := range keys {
		vm.Data.Set(v, app.Data.Get(v))
	}
	vm.Get("$options").Set("propsData", app.Props.Object)
}

func (app *App) FetchData() {
	url := `https://api.github.com/repos/vuejs/vue/commits?per_page=3&sha=`
	fmt.Println(url, app.Data.CurrentBranch)

	resp, err := http.Get(url + app.Data.CurrentBranch)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	gitcommits := []GitCommit{}
	err = json.Unmarshal(b, &gitcommits)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	Commits := []*Commit{}
	for _, g := range gitcommits {
		C := &Commit{
			Object: js.Global.Get("Object").New(),
		}
		C.HtmlUrl = g.HTMLURL
		C.Sha = g.Sha[0:7]
		C.CommitMessage = g.Commit.Message
		C.AuthorHtmlUrl = g.Author.HTMLURL
		C.CommitAuthorName = g.Commit.Author.Name
		C.CommitAuthorDate = formatDate(g.Commit.Author.Date)
		Commits = append(Commits, C)
	}
	app.Data.Commits = Commits
}

func formatDate(Date time.Time) string {
	return Date.Format("2006-01-02 15:04:05")
}

func main() {
	RegisterFilter()

	o := vue.NewOption()
	o.Data = NewAppData()
	o.OnLifeCycleEvent(vue.EvtCreated, func(vm *vue.ViewModel) {
		go func() {
			println("OnLifeCycleEvent", "EvtCreated")
			app := NewApp(vm)
			app.FetchData()
			app.SyncViewModel(vm)
		}()
	})
	o = vueutil.AddWatch(o, "currentBranch", func(vm *vue.ViewModel, newVal *js.Object, oldVal *js.Object) {
		go func() {
			app := NewApp(vm)
			app.FetchData()
			app.SyncViewModel(vm)
		}()
	})

	v := o.NewViewModel()
	v.Mount("#demo")
}

func RegisterFilter() {
	vue.NewFilter(func(v *js.Object) interface{} {
		return strings.Split(v.String(), "\n")[0]
	}).Register("truncate")
}
