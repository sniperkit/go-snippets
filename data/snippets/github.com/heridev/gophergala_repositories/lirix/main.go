package main

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"time"
	"fmt"
	"strconv"
	"html/template"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type RawCurrentData struct {
	LocationId int `json:"id"`
	UnixTime int `json:"dt"`
	LocationName string `json:"name"`

	System struct {
		Country string `json:"country"`
		Sunrise int `json:"sunrise"`
		Sunset int `json:"sunset"`
	} `json:"sys"`

	Main struct {
		Temperature float64 `json:"temp"`
		Humidity float64 `json:"humidity"`
		Pressure float64 `json:"pressure"`
	} `json:"main"`

	Wind struct {
		Speed float64 `json:"speed"`
		Direction float64 `json:"deg"`
		Gusts float64 `json:"gust"`
	} `json:"wind"`

	Clouds struct {
		Coverage float64 `json:"all"`
	} `json:"clouds"`

	Weather []struct {
		Id int `json:"id"`
		Type string `json:"main"`
		Description string `json:"description"`
		IconId string `json:"icon"`
	} `json:"weather"`

	Rain struct {
		Amount float64 `json:"3h"`
	} `json:"rain"`

	Snow struct {
		Amount float64 `json:"3h"`
	} `json:"snow"`
}

type CurrentData struct {
	DisplayName string
	LocationId string
	TargetTime string

	Sunrise string
	Sunset string

	WeatherDescription string
	WeatherIcon string

	Temperature string
	Humidity string
	Pressure string

	WindSpeed string
	WindDirection string
	WindGusts string

	CloudCoverage string

	RainAmount string
	SnowAmount string
}

type Location struct {
	DisplayName string
}

type RawDetailedData struct {
	City struct {
		Id int `json:"id"`
		Name string `json:"name"`
		Country string `json:"country"`
	} `json:"city"`

	Data []struct {
		UnixTime int `json:"dt"`

		Main struct {
			Temp float64 `json:"temp"`
			Humidity float64 `json:"humidity"`
			Pressure float64 `json:"pressure"`
		} `json:"main"`

		Wind struct {
			Speed float64 `json:"speed"`
			Direction float64 `json:"deg"`
			Gusts float64 `json:"gust"`
		} `json:"wind"`

		Clouds struct {
			Coverage float64 `json:"all"`
		} `json:"clouds"`

		Weather []struct {
			Id int `json:"id"`
			Type string `json:"main"`
			Description string `json:"description"`
			IconId string `json:"icon"`
		} `json:"weather"`

		Rain struct {
			Amount float64 `json:"3h"`
		} `json:"rain"`

		Snow struct {
			Amount float64 `json:"3h"`
		} `json:"snow"`
	} `json:"list"`
}

type DetailedData struct {
	DisplayName string
	LocationId string

	Forecast []struct {
		TargetTime string

		WeatherDescription string
		WeatherIcon string

		Temperature string
		Humidity string
		Pressure string

		WindSpeed string
		WindDirection string
		WindGusts string

		CloudCoverage string
	}
}

type RawSearchResultData struct {
	Data []struct {
		LocationId int `json:"id"`
		LocationName string `json:"name"`

		System struct {
			Country string `json:"country"`
		} `json:"sys"`
	} `json:"list"`
}

type ResultData struct {
	LocationId string
	LocationName string
}

func getCurrentData(id string) CurrentData {
	resp, _ := http.Get("http://api.openweathermap.org/data/2.5/weather?id=" + id)
	
	var raw_data RawCurrentData
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &raw_data)

	var ret CurrentData
	ret.LocationId = strconv.Itoa(raw_data.LocationId)
	ret.DisplayName = raw_data.LocationName + ", " + raw_data.System.Country

	var t time.Time
	t = time.Unix(int64(raw_data.UnixTime), 0)
	ret.TargetTime = t.Format("Monday, January 02, 2006 @ 03:04:05PM (MST)")

	t = time.Unix(int64(raw_data.System.Sunrise), 0)
	ret.Sunrise = t.Format("3:04pm (MST)")

	t = time.Unix(int64(raw_data.System.Sunset), 0)
	ret.Sunset = t.Format("3:04pm (MST)")

	ret.WeatherDescription = raw_data.Weather[0].Description
	ret.WeatherIcon = raw_data.Weather[0].IconId
	ret.Temperature = strconv.FormatFloat(raw_data.Main.Temperature - 273.15, 'f', 1, 32)
	ret.Humidity = strconv.FormatFloat(raw_data.Main.Humidity, 'f', 0, 32)
	ret.Pressure = strconv.FormatFloat(raw_data.Main.Pressure, 'f', 2, 32)
	ret.WindSpeed = strconv.FormatFloat(raw_data.Wind.Speed, 'f', 1, 32)
	ret.WindDirection = strconv.FormatFloat(raw_data.Wind.Direction, 'f', 0, 32)
	ret.WindGusts = strconv.FormatFloat(raw_data.Wind.Gusts, 'f', 1, 32)
	ret.CloudCoverage = strconv.FormatFloat(raw_data.Clouds.Coverage, 'f', 0, 32)
	ret.RainAmount = strconv.FormatFloat(raw_data.Rain.Amount, 'f', 0, 32)
	ret.SnowAmount = strconv.FormatFloat(raw_data.Snow.Amount, 'f', 0, 32)

	return ret
}

func getDetailedData(id string) DetailedData {
	resp, _ := http.Get("http://api.openweathermap.org/data/2.5/forecast?id=" + id)

	var raw_data RawDetailedData
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &raw_data)

	var ret DetailedData
	ret.LocationId = strconv.Itoa(raw_data.City.Id)

	ret.Forecast = make([]struct {
		TargetTime string
		WeatherDescription string
		WeatherIcon string
		Temperature string
		Humidity string
		Pressure string
		WindSpeed string
		WindDirection string
		WindGusts string
		CloudCoverage string
	}, len(raw_data.Data))

	for i := 0; i < len(raw_data.Data); i++ {
		var t time.Time
		t = time.Unix(int64(raw_data.Data[i].UnixTime), 0)
		ret.Forecast[i].TargetTime = t.Format("Mon 3:04pm (MST)")

		ret.Forecast[i].WeatherDescription = raw_data.Data[i].Weather[0].Description
		ret.Forecast[i].WeatherIcon = raw_data.Data[i].Weather[0].IconId
		ret.Forecast[i].Temperature = strconv.FormatFloat(raw_data.Data[i].Main.Temp - 273.15, 'f', 1, 32)
		ret.Forecast[i].Humidity = strconv.FormatFloat(raw_data.Data[i].Main.Humidity, 'f', 0, 32)
		ret.Forecast[i].Pressure = strconv.FormatFloat(raw_data.Data[i].Main.Pressure, 'f', 2, 32)
		ret.Forecast[i].WindSpeed = strconv.FormatFloat(raw_data.Data[i].Wind.Speed, 'f', 1, 32)
		ret.Forecast[i].WindDirection = strconv.FormatFloat(raw_data.Data[i].Wind.Direction, 'f', 0, 32)
		ret.Forecast[i].WindGusts = strconv.FormatFloat(raw_data.Data[i].Wind.Gusts, 'f', 1, 32)
		ret.Forecast[i].CloudCoverage = strconv.FormatFloat(raw_data.Data[i].Clouds.Coverage, 'f', 0, 32)
	}

	return ret
}

func getSearchResultData(query string) []ResultData {
	resp, _ := http.Get("http://api.openweathermap.org/data/2.5/find?q=" + query + "&type=accurate")

	var raw_data RawSearchResultData
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &raw_data)

	ret := make([]ResultData, len(raw_data.Data))

	for i := 0; i < len(raw_data.Data); i++ {
		ret[i].LocationId = strconv.Itoa(raw_data.Data[i].LocationId)
		ret[i].LocationName = raw_data.Data[i].LocationName + ", " + raw_data.Data[i].System.Country
	}

	return ret
}

func main() {
	fmt.Println("Go to http://localhost:3000 to view Lirix.")

	db, _ := sql.Open("sqlite3", "./data.db")

	var templates = template.Must(template.ParseGlob("views/*"))

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		rows, _ := db.Query("SELECT * FROM locations;", nil)
		cnt, _ := db.Query("SELECT COUNT(*) FROM locations;", nil)
		cnt.Next()
		var count int
		cnt.Scan(&count)

		weather := make([]CurrentData, count)

		i := 0
		for rows.Next() {
			var id string
			var display_name string
			rows.Scan(&id, &display_name)

			weather[i] = getCurrentData(id)
			weather[i].DisplayName = display_name
			i++
		}

		tmpl_data := map[string]interface{} {
			"title": "Lirix | Overview",
			"data": weather,
		}

		templates.ExecuteTemplate(res, "index", tmpl_data)
	})

	http.HandleFunc("/about", func(res http.ResponseWriter, req *http.Request) {
		templates.ExecuteTemplate(res, "about", nil)
	})

	/*http.HandleFunc("/search_results", func(res http.ResponseWriter, req *http.Request) {
		results := getSearchResultData(req.URL.Query().Get("q"))

		tmpl_data := map[string]interface{} {
			"title": "Lirix | Search Results",
			"query": req.URL.Query().Get("q"),
			"data": results,
		}

		templates.ExecuteTemplate(res, "search_results", tmpl_data)
	})*/

	http.HandleFunc("/detail", func(res http.ResponseWriter, req *http.Request) {
		result := getDetailedData(req.URL.Query().Get("id"))

		row, _ := db.Query("SELECT display_name FROM locations WHERE id = \"" + req.URL.Query().Get("id") + "\" LIMIT 1", nil)
		row.Next()
		var display_name string
		row.Scan(&display_name)

		result.DisplayName = display_name

		tmpl_data := map[string]interface{} {
			"data": result,
		}

		templates.ExecuteTemplate(res, "detail", tmpl_data)
	})

	/*http.HandleFunc("/add", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("Attempting to add new location...")

		data := getCurrentData(req.URL.Query().Get("id"))

		_, err := db.Exec("INSERT INTO locations (id, display_name) VALUES('" + req.URL.Query().Get("id") + "', '" + data.DisplayName + "');", nil)

		fmt.Println(err)

		http.Redirect(res, req, "/", http.StatusFound)
	})

	http.HandleFunc("/delete", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("Attempting to remove location...")
		fmt.Println(req.URL.Query().Get("id"))

		db.Exec("DELETE FROM locations WHERE id = " + req.URL.Query().Get("id") + " LIMIT 1;", nil)

		http.Redirect(res, req, "/", http.StatusFound)
	})*/

	http.ListenAndServe(":3000", nil)
}