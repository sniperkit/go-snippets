package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gophergala/edrans-smartcity/algorithm"
	"github.com/gophergala/edrans-smartcity/factory"
	"github.com/gophergala/edrans-smartcity/models"
	"github.com/gorilla/mux"
)

var sessions map[string]*models.City

type handler func(w http.ResponseWriter, r *http.Request, ctx *context) (int, interface{})

type context struct {
	Body   []byte
	CityID string
}

func main() {
	var port int
	var err error
	flag.IntVar(&port, "port", 2489, "port server will be launched")
	flag.Parse()

	sessions = make(map[string]*models.City)
	sessions["default"], err = factory.CreateRectangularCity(10, 10, "default")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	time.Sleep(2 * time.Second)

	muxRouter := mux.NewRouter()
	muxRouter.StrictSlash(false)

	muxRouter.Handle("/mobile/city/{cityID}", handler(getMobileCity)).Methods("GET")
	muxRouter.Handle("/sample-city", handler(postSampleCity)).Methods("POST")
	muxRouter.Handle("/sample-city", handler(getSampleCity)).Methods("GET")
	muxRouter.Handle("/emergency/{cityID}", handler(postEmergency)).Methods("POST")
	muxRouter.Handle("/city/{cityID}", handler(getIndex)).Methods("GET")
	muxRouter.HandleFunc("/", handleFile("main.html"))
	muxRouter.HandleFunc("/city/img/0.jpg", handleFile("img/0.jpg"))
	muxRouter.HandleFunc("/city/img/1.jpg", handleFile("img/1.jpg"))
	muxRouter.HandleFunc("/city/img/2.jpg", handleFile("img/2.jpg"))
	muxRouter.HandleFunc("/city/img/3.jpg", handleFile("img/3.jpg"))
	muxRouter.HandleFunc("/city/img/-1.jpg", handleFile("img/-1.jpg"))

	http.Handle("/", muxRouter)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Println("Cannot launch server:", err)
		os.Exit(2)
	}
	fmt.Printf("Listening on port %d...\n", port)
	http.Serve(listener, nil)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var ctx context
	var e error

	ctx.Body, e = ioutil.ReadAll(r.Body)
	if e != nil {
		fmt.Println("Error when reading body!")
	}

	vars := mux.Vars(r)
	ctx.CityID, _ = vars["cityID"]

	status, response := h(w, r, &ctx)
	if status == -1 {
		return
	}
	if status == 0 {
		status = 200
	}
	if response == nil {
		response = map[string]string{"status": "ok"}
	}
	if status < 200 || status >= 300 {
		response = map[string]interface{}{"error": response}
	}
	responseJSON, _ := json.Marshal(response)
	w.WriteHeader(status)
	w.Write(responseJSON)
}

type cityParams struct {
	SizeHorizontal int    `json:"size-horizontal"`
	SizeVertical   int    `json:"size-vertical"`
	Name           string `json:"name"`
}

type cityParams2 struct {
	SizeHorizontal string `json:"size-horizontal"`
	SizeVertical   string `json:"size-vertical"`
	Name           string `json:"name"`
}

func getSampleCity(w http.ResponseWriter, r *http.Request, ctx *context) (status int, response interface{}) {
	var sample = cityParams{SizeHorizontal: 10, SizeVertical: 10, Name: fmt.Sprintf("Sample-city-%d", len(sessions))}
	ctx.Body, _ = json.Marshal(sample)
	return postSampleCity(w, r, ctx)
}

func postSampleCity(w http.ResponseWriter, r *http.Request, ctx *context) (status int, response interface{}) {
	var in cityParams
	status = 302
	var url string
	err := json.Unmarshal(ctx.Body, &in)
	if in.Name == "" || err != nil {
		status = 400
		response = "Bad json"
		return
	}

	if sessions[in.Name] != nil {
		status = 403
		response = "city already exists"
		return
	}

	sessions[in.Name], err = factory.CreateRectangularCity(in.SizeHorizontal, in.SizeVertical, in.Name)
	if err != nil {
		status = 400
		response = err.Error()
		return
	}

	if status != 302 {
		url = "/error"
	} else {
		url = fmt.Sprintf("/city/%s", in.Name)
	}
	http.Redirect(w, r, url, status)
	return -1, nil
}

func getMobileCity(w http.ResponseWriter, r *http.Request, ctx *context) (status int, response interface{}) {
	if ctx.CityID == "" {
		status = 404
		response = "City doesn't exist"
		return
	}
	response = *sessions[ctx.CityID]
	return
}

type emergencyRequest struct {
	Service      string `json:"service"`
	WhereRequest string `json:"where"`
	Where        int
}

func postEmergency(w http.ResponseWriter, r *http.Request, ctx *context) (status int, response interface{}) {
	var emergency emergencyRequest
	e := json.Unmarshal(ctx.Body, &emergency)
	if e != nil {
		status = 400
		response = e.Error()
		return
	}

	emergency.Where, e = strconv.Atoi(emergency.WhereRequest)
	if e != nil {
		status = 400
		response = e.Error()
		return
	}

	city := sessions[ctx.CityID]
	if city == nil {
		status = 404
		response = "city not found"
		return
	}

	defer city.CleanError()
	vehicle, e := city.CallService(emergency.Service)
	if e != nil {
		status = 400
		response = e.Error()
		return
	}

	paths, e := algorithm.GetPaths(city, vehicle.Position.ID, emergency.Where)
	if e != nil {
		status = 400
		response = e.Error()
		return
	}

	paths = algorithm.CalcEstimatesForVehicle(vehicle, paths)
	if len(paths) == 0 {
		status = 400
		response = "Oh no... there is no way to go"
		return
	}

	toRun1 := algorithm.ChooseBest(paths)
	paths, _ = algorithm.GetPaths(city, emergency.Where, vehicle.BasePosition.ID)
	paths = algorithm.CalcEstimatesForVehicle(vehicle, paths)

	if len(paths) == 0 {
		status = 400
		response = "Oh no... there is no way to come back"
		return
	}

	vehicle.Alert <- toRun1
	vehicle.Alert <- algorithm.ChooseBest(paths)
	response = fmt.Sprintf("%s on the way to %d. It is %d blocks away", emergency.Service, emergency.Where, len(toRun1.Links))
	return
}

func getIndex(w http.ResponseWriter, r *http.Request, ctx *context) (status int, response interface{}) {
	var index = make([]string, 0)
	file, e := ioutil.ReadFile("index.html")
	if e != nil {
		return 503, "index not found"
	}
	fileLines := strings.Split(string(file), "\n")
	for i := 0; i < len(fileLines); i++ {
		index = append(index, fileLines[i])
		if strings.Contains(fileLines[i], "<table") {
			if sessions[ctx.CityID] == nil {
				status = 404
				response = "city does not exist"
				return
			}
			table := createTable(ctx.CityID)
			for j := 0; j < len(table); j++ {
				index = append(index, table[j])
			}
		}
	}
	http.ServeContent(w, r, "city", time.Now(), bytes.NewReader([]byte(strings.Join(index, "\n"))))
	status = -1
	response = strings.Join(index, "\n")
	return
}

func createTable(cityID string) []string {
	var table = make([]string, 0)
	city := sessions[cityID]
	locations := city.GetLocations()
	for i := 0; i < city.Size[0]; i++ {
		table = append(table, "<tr>")
		myNodes := getNodes(locations, -i)
		for j := 0; j < len(myNodes); j++ {
			var color string
			switch myNodes[j].Vehicle {
			case 0:
				color = fmt.Sprintf(`bgcolor="#0000FF"`)
			case 1:
				color = fmt.Sprintf(`bgcolor="#FF0000"`)
			case 2:
				color = fmt.Sprintf(`bgcolor="#228B22"`)
			}
			insert := fmt.Sprintf(`<img src="img/%d.jpg" height="20" width="20" /> %d`, myNodes[j].Input, myNodes[j].Weight)
			table = append(table, fmt.Sprintf(`<td style="width:100px" %s> %s </td>`, color, insert))
		}
		table = append(table, "</tr>")
	}
	return table
}

func getNodes(locations []models.Location, nodes int) []models.Location {
	var local = make([]models.Location, 0)
	for i := 0; i < len(locations); i++ {
		if locations[i].Long == nodes {
			local = append(local, locations[i])
		}
	}
	return local
}

func handleFile(path string) http.HandlerFunc {
	path = filepath.Join("", path)
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}
