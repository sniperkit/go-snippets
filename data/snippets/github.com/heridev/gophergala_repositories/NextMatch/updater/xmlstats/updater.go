package xmlstats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type (
	Events struct {
		Event      []Event `json:"event"`
		EventsDate string  `json:"events_date"`
	}

	Results struct {
		EventID                 string  `json:"event_id"`
		EventSeasonType         string  `json:"event_season_type"`
		EventStartDateTime      string  `json:"event_start_date_time"`
		EventStatus             string  `json:"event_status"`
		Opponent                Team    `json:"opponent"`
		OpponentEventsLost      float64 `json:"opponent_events_lost"`
		OpponentEventsWon       float64 `json:"opponent_events_won"`
		OpponentPointsScored    float64 `json:"opponent_points_scored"`
		Site                    Site    `json:"site"`
		Team                    Team    `json:"team"`
		TeamEventLocationType   string  `json:"team_event_location_type"`
		TeamEventNumberInSeason float64 `json:"team_event_number_in_season"`
		TeamEventResult         string  `json:"team_event_result"`
		TeamEventsLost          float64 `json:"team_events_lost"`
		TeamEventsWon           float64 `json:"team_events_won"`
		TeamPointsScored        float64 `json:"team_points_scored"`
	}

	Team struct {
		Abbreviation string `json:"abbreviation"`
		Active       bool   `json:"active"`
		City         string `json:"city"`
		Conference   string `json:"conference"`
		Division     string `json:"division"`
		FirstName    string `json:"first_name"`
		FullName     string `json:"full_name"`
		LastName     string `json:"last_name"`
		SiteName     string `json:"site_name"`
		State        string `json:"state"`
		TeamID       string `json:"team_id"`
		Logo         string
	}

	Site struct {
		Capacity float64 `json:"capacity"`
		City     string  `json:"city"`
		Name     string  `json:"name"`
		State    string  `json:"state"`
		Surface  string  `json:"surface"`
	}

	BoxScore struct {
		AwayPeriodScores []float64 `json:"away_period_scores"`
		AwayStats        []Stats   `json:"away_stats"`
		AwayTeam         Team      `json:"away_team"`
		AwayTotals       Stats     `json:"away_totals"`
		EventInformation struct {
			Attendance float64 `json:"attendance"`
			Duration   string  `json:"duration"`
			SeasonType string  `json:"season_type"`
			Site       struct {
				Capacity float64 `json:"capacity"`
				City     string  `json:"city"`
				Name     string  `json:"name"`
				State    string  `json:"state"`
				Surface  string  `json:"surface"`
			} `json:"site"`
			StartDateTime string  `json:"start_date_time"`
			Temperature   float64 `json:"temperature"`
		} `json:"event_information"`
		HomePeriodScores []float64 `json:"home_period_scores"`
		HomeStats        []Stats   `json:"home_stats"`
		HomeTeam         Team      `json:"home_team"`
		HomeTotals       Stats     `json:"home_totals"`
		Officials        []struct {
			FirstName string      `json:"first_name"`
			LastName  string      `json:"last_name"`
			Position  interface{} `json:"position"`
		} `json:"officials"`
	}

	Stats struct {
		Assists                       float64 `json:"assists"`
		Blocks                        float64 `json:"blocks"`
		DefensiveRebounds             float64 `json:"defensive_rebounds"`
		DisplayName                   string  `json:"display_name"`
		FieldGoalPercentage           float64 `json:"field_goal_percentage"`
		FieldGoalsAttempted           float64 `json:"field_goals_attempted"`
		FieldGoalsMade                float64 `json:"field_goals_made"`
		FirstName                     string  `json:"first_name"`
		FreeThrowPercentage           float64 `json:"free_throw_percentage"`
		FreeThrowsAttempted           float64 `json:"free_throws_attempted"`
		FreeThrowsMade                float64 `json:"free_throws_made"`
		IsStarter                     bool    `json:"is_starter"`
		LastName                      string  `json:"last_name"`
		Minutes                       float64 `json:"minutes"`
		OffensiveRebounds             float64 `json:"offensive_rebounds"`
		PersonalFouls                 float64 `json:"personal_fouls"`
		Points                        float64 `json:"points"`
		Position                      string  `json:"position"`
		Steals                        float64 `json:"steals"`
		TeamAbbreviation              string  `json:"team_abbreviation"`
		ThreePointFieldGoalsAttempted float64 `json:"three_point_field_goals_attempted"`
		ThreePointFieldGoalsMade      float64 `json:"three_point_field_goals_made"`
		ThreePointPercentage          float64 `json:"three_point_percentage"`
		Turnovers                     float64 `json:"turnovers"`
	}

	Event struct {
		EventID       string `json:"event_id"`
		EventStatus   string `json:"event_status"`
		AwayTeam      Team   `json:"away_team"`
		HomeTeam      Team   `json:"home_team"`
		SeasonType    string `json:"season_type"`
		Site          Site   `json:"site"`
		Sport         string `json:"sport"`
		StartDateTime string `json:"start_date_time"`
	}
)

const (
	shortf    = "20060102"
	eventURI  = "https://erikberg.com/events.json?sport=%s&date=%s" // league[nba,nfl] date
	resultURI = "https://erikberg.com/%s/results/%s.json"           // team_id
	scoreURI  = "https://erikberg.com/%s/boxscore/%s.json"          // team_id
	userAgent = "nextmatch/0.1 (https://twitter.com/oscarryz)"
	auth      = "Bearer %s"
)

var (
	Token string
)

var cache = make(map[string]interface{})

func doRequest(uri string, result interface{}) error {

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf(auth, Token))
	req.Header.Add("User-agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	log.Printf("%s for URI %s", resp.Status, uri)

	err = decode(resp, result)
	return err
}

// BySport returns events in certain sport [with an optional date]
func BySport(sport string, date ...string) (ev Events, err error) {

	if len(date) < 1 {
		date = append(date, time.Now().Format(shortf))
	}
	uri := fmt.Sprintf(eventURI, sport, date[0])
	if cache[uri] != nil {
		log.Printf("cache  for URI %s", uri)

		return cache[uri].(Events), nil
	}

	err = doRequest(uri, &ev)
	cache[uri] = ev

	for i, _ := range ev.Event {
		ev.Event[i].AwayTeam.Logo = teamLogos[ev.Event[i].AwayTeam.TeamID]
		ev.Event[i].HomeTeam.Logo = teamLogos[ev.Event[i].HomeTeam.TeamID]
		//ev.AwayTeam.Logo = teamLogos[ev.AwayTeam.TeamID]
		//ev.HomeTeam.Logo = teamLogos[ev.HomeTeam.TeamID]
	}

	return ev, err
}

// Result
func Result(sport, teamId string) (results Results, err error) {
	uri := fmt.Sprintf(resultURI, sport, teamId)
	if cache[uri] != nil {
		return cache[uri].(Results), nil
	}
	err = doRequest(uri, &results)
	cache[uri] = results
	return results, err
}

// Score
func Score(sport, eventId string) (results BoxScore, err error) {
	uri := fmt.Sprintf(scoreURI, sport, eventId)
	if cache[uri] != nil {
		return cache[uri].(BoxScore), nil
	}
	err = doRequest(uri, &results)

	results.AwayTeam.Logo = teamLogos[results.AwayTeam.TeamID]
	results.HomeTeam.Logo = teamLogos[results.HomeTeam.TeamID]

	cache[uri] = results
	return results, err
}

func decode(resp *http.Response, d interface{}) error {
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(d)
}

func Unmarshal(b string, d interface{}) error {
	data := []byte(b)
	return json.Unmarshal(data, d)
}

func (e *Events) ById(id string) *Event {
	for _, v := range e.Event {
		if v.EventID == id {
			return &v
		}
	}

	return nil
}

var teamLogos = map[string]string{
	"atlanta-hawks":          "http://content.sportslogos.net/logos/6/220/thumbs/5mdhgjh3aa92kih09pgi.gif",
	"boston-celtics":         "http://content.sportslogos.net/logos/6/213/thumbs/slhg02hbef3j1ov4lsnwyol5o.gif",
	"brooklyn-nets":          "http://content.sportslogos.net/logos/6/3786/thumbs/hsuff5m3dgiv20kovde422r1f.gif",
	"charlotte-hornets":      "http://content.sportslogos.net/logos/6/5120/thumbs/512019262015.gif",
	"chicago-bulls":          "http://content.sportslogos.net/logos/6/221/thumbs/hj3gmh82w9hffmeh3fjm5h874.gif",
	"cleveland-cavaliers":    "http://content.sportslogos.net/logos/6/222/thumbs/e4701g88mmn7ehz2baynbs6e0.gif",
	"dallas-mavericks":       "http://content.sportslogos.net/logos/6/228/thumbs/ifk08eam05rwxr3yhol3whdcm.gif",
	"denver-nuggets":         "http://content.sportslogos.net/logos/6/229/thumbs/xeti0fjbyzmcffue57vz5o1gl.gif",
	"detroit-pistons":        "http://content.sportslogos.net/logos/6/223/thumbs/3079.gif",
	"golden-state-warriors":  "http://content.sportslogos.net/logos/6/235/thumbs/qhhir6fj8zp30f33s7sfb4yw0.gif",
	"houston-rockets":        "http://content.sportslogos.net/logos/6/230/thumbs/8xe4813lzybfhfl14axgzzqeq.gif",
	"indiana-pacers":         "http://content.sportslogos.net/logos/6/224/thumbs/3083.gif",
	"los-angeles-clippers":   "http://content.sportslogos.net/logos/6/236/thumbs/bvv028jd1hhr8ee8ii7a0fg4i.gif",
	"los-angeles-lakers":     "http://content.sportslogos.net/logos/6/237/thumbs/uig7aiht8jnpl1szbi57zzlsh.gif",
	"memphis-grizzlies":      "http://content.sportslogos.net/logos/6/231/thumbs/793.gif",
	"miami-heat":             "http://content.sportslogos.net/logos/6/214/thumbs/burm5gh2wvjti3xhei5h16k8e.gif",
	"milwaukee-bucks":        "http://content.sportslogos.net/logos/6/225/thumbs/0295onf2c4xsbfsxye6i.gif",
	"minnesota-timberwolves": "http://content.sportslogos.net/logos/6/232/thumbs/zq8qkfni1g087f4245egc32po.gif",
	"new-orleans-pelicans":   "http://content.sportslogos.net/logos/6/4962/thumbs/496226812014.gif",
	"new-york-knicks":        "http://content.sportslogos.net/logos/6/216/thumbs/2nn48xofg0hms8k326cqdmuis.gif",
	"oklahoma-city-thunder":  "http://content.sportslogos.net/logos/6/2687/thumbs/khmovcnezy06c3nm05ccn0oj2.gif",
	"orlando-magic":          "http://content.sportslogos.net/logos/6/217/thumbs/wd9ic7qafgfb0yxs7tem7n5g4.gif",
	"philadelphia-76ers":     "http://content.sportslogos.net/logos/6/218/thumbs/qlpk0etqwelv8artgc7tvqefu.gif",
	"phoenix-suns":           "http://content.sportslogos.net/logos/6/238/thumbs/23843702014.gif",
	"portland-trail-blazers": "http://content.sportslogos.net/logos/6/239/thumbs/bahmh46cyy6eod2jez4g21buk.gif",
	"sacramento-kings":       "http://content.sportslogos.net/logos/6/240/thumbs/832.gif",
	"san-antonio-spurs":      "http://content.sportslogos.net/logos/6/233/thumbs/827.gif",
	"toronto-raptors":        "http://content.sportslogos.net/logos/6/227/thumbs/yfypcwqog6qx8658sn5w65huh.gif",
	"utah-jazz":              "http://content.sportslogos.net/logos/6/234/thumbs/m2leygieeoy40t46n1qqv0550.gif",
	"washington-wizards":     "http://content.sportslogos.net/logos/6/219/thumbs/b3619brnphtx65s2th4p9eggf.gif",
	// NFL
	"arizona-cardinals":    "http://content.sportslogos.net/logos/7/177/thumbs/kwth8f1cfa2sch5xhjjfaof90.gif",
	"atlanta-falcons":      "http://content.sportslogos.net/logos/7/173/thumbs/299.gif",
	"baltimore-ravens":     "http://content.sportslogos.net/logos/7/153/thumbs/318.gif",
	"buffalo-bills":        "http://content.sportslogos.net/logos/7/149/thumbs/n0fd1z6xmhigb0eej3323ebwq.gif",
	"carolina-panthers":    "http://content.sportslogos.net/logos/7/174/thumbs/f1wggq2k8ql88fe33jzhw641u.gif",
	"chicago-bears":        "http://content.sportslogos.net/logos/7/169/thumbs/364.gif",
	"cincinnati-bengals":   "http://content.sportslogos.net/logos/7/154/thumbs/403.gif",
	"cleveland-browns":     "http://content.sportslogos.net/logos/7/155/thumbs/2ioheczrkmc2ibc42c9r.gif",
	"dallas-cowboys":       "http://content.sportslogos.net/logos/7/165/thumbs/406.gif",
	"denver-broncos":       "http://content.sportslogos.net/logos/7/161/thumbs/9ebzja2zfeigaziee8y605aqp.gif",
	"detroit-lions":        "http://content.sportslogos.net/logos/7/170/thumbs/cwuyv0w15ruuk34j9qnfuoif9.gif",
	"green-bay-packers":    "http://content.sportslogos.net/logos/7/171/thumbs/dcy03myfhffbki5d7il3.gif",
	"houston-texans":       "http://content.sportslogos.net/logos/7/157/thumbs/570.gif",
	"indianapolis-colts":   "http://content.sportslogos.net/logos/7/158/thumbs/593.gif",
	"jacksonville-jaguars": "http://content.sportslogos.net/logos/7/159/thumbs/15988562013.gif",
	"kansas-city-chiefs":   "http://content.sportslogos.net/logos/7/162/thumbs/857.gif",
	"miami-dolphins":       "http://content.sportslogos.net/logos/7/150/thumbs/15041052013.gif",
	"minnesota-vikings":    "http://content.sportslogos.net/logos/7/172/thumbs/17227042013.gif",
	"new-england-patriots": "http://content.sportslogos.net/logos/7/151/thumbs/y71myf8mlwlk8lbgagh3fd5e0.gif",
	"new-orleans-saints":   "http://content.sportslogos.net/logos/7/175/thumbs/907.gif",
	"new-york-giants":      "http://content.sportslogos.net/logos/7/166/thumbs/919.gif",
	"new-york-jets":        "http://content.sportslogos.net/logos/7/152/thumbs/v7tehkwthrwefgounvi7znf5k.gif",
	"oakland-raiders":      "http://content.sportslogos.net/logos/7/163/thumbs/g9mgk6x3ge26t44cccm9oq1vl.gif",
	"philadelphia-eagles":  "http://content.sportslogos.net/logos/7/167/thumbs/960.gif",
	"pittsburgh-steelers":  "http://content.sportslogos.net/logos/7/156/thumbs/970.gif",
	"san-diego-chargers":   "http://content.sportslogos.net/logos/7/164/thumbs/8e1jhgblydtow4m3okwzxh67k.gif",
	"san-francisco-49ers":  "http://content.sportslogos.net/logos/7/179/thumbs/17994552009.gif",
	"seattle-seahawks":     "http://content.sportslogos.net/logos/7/180/thumbs/pfiobtreaq7j0pzvadktsc6jv.gif",
	"st.-louis-rams":       "http://content.sportslogos.net/logos/7/178/thumbs/1029.gif",
	"tampa-bay-buccaneers": "http://content.sportslogos.net/logos/7/176/thumbs/17636702014.gif",
	"tennessee-titans":     "http://content.sportslogos.net/logos/7/160/thumbs/1053.gif",
	"washington-redskins":  "http://content.sportslogos.net/logos/7/168/thumbs/im5xz2q9bjbg44xep08bf5czq.gif",
}
