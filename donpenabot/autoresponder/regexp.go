package autoresponder

import (
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

var weekdays map[string]int

func init() {
	weekdays = map[string]int{
		"domingo":   0,
		"lunes":     1,
		"martes":    2,
		"miércoles": 3,
		"jueves":    4,
		"viernes":   5,
		"sábado":    6,
	}
}

func fecha(date time.Time, message string, api anaconda.TwitterApi) {
	for {
		currentdate := time.Now()
		if currentdate.UTC().Unix() >= date.UTC().Unix() {
			fmt.Printf(message)
			api.PostTweet(message, nil)
			return
		}
		time.Sleep(60 * time.Minute)
	}
}

// func main() {
// 	// anaconda.SetConsumerKey(os.Getenv("APIKEY"))
// 	// anaconda.SetConsumerSecret(os.Getenv("APISECRET"))
// 	// api := anaconda.NewTwitterApi(os.Getenv("ACCESSTOKEN"), os.Getenv("ACCESSTOKENSECRET"))
// 	r, _ := regexp.Compile("(@donpenabot\\s)+(recuerdame)+(\\s.{0,140})+(\\sel\\s)+([0-9]{2,4}-[0-9]{2,4}-[0-9]{2,4})")
// 	fmt.Println(r.FindAllString("blah blah @donpenabot recuerdame ir al banco @donpenabot recuerdame ir al banco el 11-02-16 ijoiijo", -1))
// 	//s := regexp.MustCompile("a*").Split("abaabaccadaaae", 5)
// 	//str, _ := regexp.MatchString("cantame una canci(o|ó)n", "cantame una cancion")
// 	//fmt.Println(str)
// 	//g, _ := regexp.Compile(".{1,117} ([0-9]{2}\\/[0-9]{2}\\/[0-9]{2,4}|[0-9]{2}\\-[0-9]{2}\\-[0-9]{2,4}|mañana|lunes|martes|mi(e|é)rcoles|jueves|viernes|s(a|á)bado|domingo)").Match()
// 	//fmt.Println(g.FindAllString("bjj cantame ijdoifjodis", -1))
// 	//array := regexp.MustCompile("@donpenabot recuerdame .{1,117} ([0-9]{2}\\/[0-9]{2}\\/[0-9]{2,4}|[0-9]{2}\\-[0-9]{2}\\-[0-9]{2,4}|mañana|lunes|martes|mi(e|é)rcoles|jueves|viernes|s(a|á)bado|domingo)").Split("@donpenabot recuerdame ir al banco mañana", -1)
// 	fecha1 := getDate("domingo")
// 	fmt.Println(fecha1.Format(time.ANSIC))
// 	// 	v := url.Values{}
// 	// 	v.Set("count", "10")
// 	// 	results1, _ := api.GetHomeTimeline(v)
// 	// 	time.Sleep(5 * time.Second)
// 	// 	results2, _ := api.GetHomeTimeline(v)
// 	//
// 	// 	ret := compare(results1, results2[0:7])
// 	// 	slice1 := []string{}
// 	// 	slice2 := []string{}
// 	//
// 	// 	for x := 0; x < len(results1); x++ {
// 	// 		slice1 = append(slice1, results1[x].IdStr)
// 	// 		slice2 = append(slice2, results2[x].IdStr)
// 	// 	}
// 	// 	fmt.Printf("%s", slice1)
// 	// 	fmt.Printf("%s", slice2)
// 	// 	fmt.Printf("%d", len(ret))
// }

func RespondTweet(tweet anaconda.Tweet, api *anaconda.TwitterApi) error {
	switch {
	case regexp.MustCompile("@donpenabot c(a|á)ntame una canci(o|ó)n").Match([]byte(tweet.Text)):
		str1 := "@" + tweet.User.ScreenName + " https://www.youtube.com/watch?v=RpCTu2ymqiM"
		api.PostTweet(str1, nil)
		fmt.Println("caso 1 jhbbkjkjn")
		fmt.Println(str1)
	case regexp.MustCompile("@donpenabot m(a|á)ndame besitos").Match([]byte(tweet.Text)):
		str2 := "@" + tweet.User.ScreenName + " mua mua mua"
		api.PostTweet(str2, nil)
		fmt.Println("caso 2")
	case regexp.MustCompile("@donpenabot chinga tu madre").Match([]byte(tweet.Text)):
		str3 := "@" + tweet.User.ScreenName + " luego por que los desaparezco"
		api.PostTweet(str3, nil)
		fmt.Println("caso 3")
	case regexp.MustCompile("@donpenabot h(a|á)blame en ingl(e|é)s").Match([]byte(tweet.Text)):
		str4 := "@" + tweet.User.ScreenName + " The steit mosbe de greit promotor an regulator of de efishien opereichion of de marquez in su Mary de..."
		api.PostTweet(str4, nil)
		fmt.Println("caso 4")
		fmt.Println(str4)
	case regexp.MustCompile("(@donpenabot\\s)+(recuerdame)+(\\s.{0,140})+(\\sel\\s)+([0-9]{2,4}-[0-9]{2,4}-[0-9]{2,4})").Match([]byte(tweet.Text)):
		str5 := "@" + tweet.User.ScreenName + " "
		array := regexp.MustCompile("(@donpenabot\\s)+(recuerdame)+(\\s.{0,140})+(\\sel\\s)+([0-9]{2,4}-[0-9]{2,4}-[0-9]{2,4})").FindAllString(tweet.Text, -1)
		fmt.Printf("%s", array)
		fmt.Printf("%d", len(array))
		//api.PostTweet(str5, nil)
		fmt.Println("caso 5")
		fmt.Println(str5)
		// case regexp.MustCompile("@donpenabot recu(e|é)rdame .{1,117} ([0-9]{2}\\/[0-9]{2}\\/[0-9]{2,4}|[0-9]{2}\\-[0-9]{2}\\-[0-9]{2,4}|mañana|lunes|martes|mi(e|é)rcoles|jueves|viernes|s(a|á)bado|domingo)").Match([]byte(tweet.Text)):
		// array := regexp.MustCompile("@donpenabot recuerdame .{1,117}+([0-9]{2}\\/[0-9]{2}\\/[0-9]{2,4}|[0-9]{2}\\-[0-9]{2}\\-[0-9]{2,4}|mañana|lunes|martes|mi(e|é)rcoles|jueves|viernes|s(a|á)bado|domingo)").Split(tweet.Text, -1)
		// //fecha := time.Now()
		// //str4 := tweet.User.ScreenName + ""
		// fmt.Printf("%s", array)
	}
	return nil
}

func getDate(dia string) time.Time {
	if dia == "mañana" {
		return time.Now().Add(24 * time.Hour)
	}
	currentDate := time.Now()
	if int(currentDate.Weekday()) < weekdays[dia] {
		daysToAdd := weekdays[dia] - int(currentDate.Weekday())
		daysToAdd = daysToAdd * 24
		return currentDate.Add(time.Duration(daysToAdd) * time.Hour)
	} else {
		daysToAdd := 6 - int(currentDate.Weekday()) + weekdays[dia]
		daysToAdd = daysToAdd * 24
		return currentDate.Add(time.Duration(daysToAdd) * time.Hour)
	}
}

func compare(x, y []anaconda.Tweet) []anaconda.Tweet {
	m := make(map[string]int)
	for _, tweet := range y {
		m[tweet.IdStr]++
	}
	var ret []anaconda.Tweet
	for _, tweet := range x {
		if m[tweet.IdStr] > 0 {
			m[tweet.IdStr]--
			continue
		}
		ret = append(ret, tweet)
	}
	return ret
}

func Respond(api *anaconda.TwitterApi) {
	v := url.Values{}
	v.Set("count", "10")
	results1, _ := api.GetMentionsTimeline(v)
	for {
		results2, _ := api.GetMentionsTimeline(v)
		comparacion := compare(results2, results1)
		for _, tweet := range comparacion {
			go RespondTweet(tweet, api)
		}
		results1 = results2
		time.Sleep(30 * time.Second)
	}

}
