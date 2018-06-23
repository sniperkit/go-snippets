package github

//import "reflect"
import "fmt"
import "strconv"
import "bytes"
import "encoding/json"
import "github.com/pravj/hackman/utils/request"

const (
	TEAM_CREATION_ENDPOINT   string = "https://api.github.com/orgs/"
	TEAM_MEMBERSHIP_ENDPOINT string = "https://api.github.com/teams/"
        HOOK_CREATION_ENDPOINT   string = "https://api.github.com/orgs/"
)

type User struct {
	UserName string
}

type Team struct {
	Id    int64
	Name  string
	Users []User
}

type TeamParameter struct {
	Name       string `json:"name"`
	Permission string `json:"permission"`
}

type WebhookParameter struct {
  Name string `json:"name"`
  Config WebhookConfig `json:"config"`
}

type WebhookConfig struct {
  Url string `json:"url"`
  ContentType string `json:"content_type"`
}

type TeamResponse struct {
	Id int `json:"id"`
}

type Organization struct {
	Name  string
	Teams []Team
}

func CreateTeams(org *Organization, accessToken string) {
	for _, team := range org.Teams {
		payloadJson, _ := json.Marshal(TeamParameter{Name: team.Name, Permission: "push"})
		payloadReader := bytes.NewReader(payloadJson)

		body := request.Post(TEAM_CREATION_ENDPOINT+org.Name+"/teams", accessToken, payloadReader, false)
		fmt.Println(string(body))

		var tr TeamResponse
		json.Unmarshal(body, &tr)

		AddTeamMembers(team.Users, tr.Id, accessToken)
	}
}

func AddTeamMembers(users []User, id int, accessToken string) {
	for _, user := range users {
		url := TEAM_MEMBERSHIP_ENDPOINT + strconv.Itoa(id) + "/memberships/" + user.UserName
		body := request.Put(url, accessToken, nil)
		fmt.Println(string(body))
	}
}

func AddOrgWebhook(org *Organization, accessToken string) {
  webhookConfig := WebhookConfig{Url: "http://requestb.in/uxp2gxux", ContentType: "json"}
  payloadJson, _ := json.Marshal(WebhookParameter{Name: "web", Config: webhookConfig})
  payloadReader := bytes.NewReader(payloadJson)

  body := request.Post(HOOK_CREATION_ENDPOINT + org.Name + "/hooks", accessToken, payloadReader, true)
  fmt.Println(string(body))
}

func main() {
	//user1 := User{UserName: "pravj"}
	//user2 := User{UserName: "iMshyam"}

	//fmt.Println(reflect.TypeOf(user1))

	//team1 := Team{Id: 0, Name: "testing", Users: []User{user1}}
	//team2 := Team{Id: 0, Name: "shyam", Users: []User{user2}}
	//fmt.Println(reflect.TypeOf(team1))
	org := Organization{Name: "mockers", Teams: []Team{}}
        AddOrgWebhook(&org, "xxxxxxxxx")

	//CreateTeams(&org, "xxxxxxxxx")
}
