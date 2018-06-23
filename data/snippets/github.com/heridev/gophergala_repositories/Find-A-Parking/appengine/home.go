package home

import (
	"html/template"
	"net/http"
	"appengine"
    "appengine/datastore"
    "time"
    "fmt"
    "encoding/json"
    "appengine/channel"
    "strconv"
)


type Response map[string]interface{}

func (r Response) String() (s string) {
        b, err := json.Marshal(r)
        if err != nil {
                s = ""
                return
        }
        s = string(b)
        return
}

type Parking struct {
	Owner string
	Mail string
	Price float64
	Location appengine.GeoPoint
	Recorded time.Time
}

type Transactions struct {
	Park Parking
	Begin time.Time
	End time.Time
	Customer string
}

func init() {
	http.HandleFunc("/",index)
	http.HandleFunc("/home",home)
	http.HandleFunc("/rent",rent)
	http.HandleFunc("/createPark",createParking)
	http.HandleFunc("/getToken",getToken)
	http.HandleFunc("/getParkings",getParkings)
	http.HandleFunc("/parkAuto",parkAuto)
}

func parkAuto(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	price, _ := strconv.ParseFloat("50", 64)
	lat, _ := strconv.ParseFloat("25.670708", 64)
	lng, _ := strconv.ParseFloat("-100.308172", 64)
	geo := appengine.GeoPoint{Lat: lat,Lng: lng}
	parking := &Parking{
		Owner:	"Garces Parks",
		Mail:	"garces@parks.com",
		Price:	price,
		Location: geo,
		Recorded: time.Now(),
	}

	channel.SendJSON(c, "customer", parking)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"success": "ok"})
    return

}

func getParkings(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Parking")
	var parkings []Parking
	q.GetAll(c, &parkings)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"parkings": parkings})
	return
}

func getToken(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	token, err := channel.Create(c, "customer")
	if err != nil {
	    token = ""
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"token": token})
    return
}

func index(w http.ResponseWriter, r *http.Request) {
	d := map[string]interface{}{"Titulo": "Find A Park"}
	t, err := template.ParseFiles("templates/index.html", "templates/base.html")
  	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
  	} else{
    	err := t.ExecuteTemplate(w,"base", d)
	    if err != nil {
	      http.Error(w, err.Error(), http.StatusInternalServerError)
	    } 
  	}
}

func home(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	

	token, err := channel.Create(c, "customer")
	if err != nil {
	    token = ""
	}

	d := map[string]interface{}{"Titulo": "Home", "Token": token}
  	t, err := template.ParseFiles("templates/home.html", "templates/base.html")
  	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
  	} else{
	    err := t.ExecuteTemplate(w,"base", d)
	    if err != nil {
	      http.Error(w, err.Error(), http.StatusInternalServerError)
	    } 
  	}
}

func rent(w http.ResponseWriter, r *http.Request) {
	

	d := map[string]interface{}{"Titulo": "Rent"}


  	t, err := template.ParseFiles("templates/renta.html", "templates/base.html")
  	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
  	} else{
	    err := t.ExecuteTemplate(w,"base", d)
	    if err != nil {
	      http.Error(w, err.Error(), http.StatusInternalServerError)
	    } 
  	}	
}

func createParking(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	lat, _ := strconv.ParseFloat(r.FormValue("lat"), 64)
	lng, _ := strconv.ParseFloat(r.FormValue("lng"), 64)
	geo := appengine.GeoPoint{Lat: lat,Lng: lng}
	parking := &Parking{
		Owner:	r.FormValue("name"),
		Mail:	r.FormValue("email"),
		Price:	price,
		Location: geo,
		Recorded: time.Now(),
	}

	key := datastore.NewIncompleteKey(c, "Parking", nil);
	_, err := datastore.Put(c,key,parking)
	if err != nil {
	    http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	channel.SendJSON(c, "customer", parking)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"success": err == nil})
    return
}




