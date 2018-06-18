package scripts

import (
	"database/sql"
	"fmt"
	"time"
	// import the Postgres lib
	_ "github.com/lib/pq"
	// import application related stuff
	"git.monkey-works.de/scripting/api"
	"monkey-works.de/model"
)

// The endpoint of our web service we want to connect to
// Change this to your endpoint
const URL string = "xxx.xxx.xxx.x" // ip address or name of the host
const USERNAME string = "user"
const PASSWORD string = "user_pw"
const DATABASE string = "database_name"
const PORT string = "5432"                         // default is 5432
const QUERY string = "Select * FROM \"variables\"" // your SQL Statement

// UPDATE_INTERVAL the interval at which the connector is updated.
// Time is defined in milliseconds.
const UPDATE_INTERVAL time.Duration = 1000

// Data is a global channel to consume received data.
var Data = make(chan *Response)

// Stop is a global channel to control periodic task.
var Stop = make(chan bool)

var db *sql.DB

var app model.Application

// status shows if the application is running smoothly or has any errors
var status model.StringDataItem

// DataItem symbolises a data structure of simple data item provided by REST interface.
// Adapt this to your data structure.
type DataItem struct {
	Name      string `json:name`
	Value     string `json:value`
	Timestamp string `json:timestamp`
}

// Response is the reponse of REST interface.
// Adapt this to your data structure.
type Response struct {
	DataItems []DataItem `json:"dataItems`
}

// @script
func InitializeScripting(application model.Application) {
	fmt.Println("Hello Scripting")
	var ok bool
	// get status item
	status, ok = application.ClientDataModel().FindDataItemByName("status").(model.StringDataItem)
	if !ok {
		panic("DataItem \"status\" not found")
	}
	app = application
	// listen to two buttons  and trigger start and stop if button get pressed
	registerButton()
	registerStopButton()
}

// processData is processing the received data
func processData(data *Response) {
	if data == nil {
		status.SetCurrentValue("Error")
		return
	}
	for i, value := range data.DataItems {
		//
		// Add your code here ...
		// We'll make an example and print the last Element of DataItems into the  3 Labels
		valueDataItem := app.ClientDataModel().FindDataItemByName("valueData").(model.StringDataItem)
		nameDataItem := app.ClientDataModel().FindDataItemByName("name").(model.StringDataItem)
		timestampDataItem := app.ClientDataModel().FindDataItemByName("time").(model.StringDataItem)
		if i == len(data.DataItems)-1 {
			valueDataItem.SetCurrentValue(value.Value)
			nameDataItem.SetCurrentValue(value.Name)
			timestampDataItem.SetCurrentValue(value.Timestamp)
			status.SetCurrentValue("Retrieved data successfully")
		}
	}
}

// getData starts every UPDATE_INTERVAL milliseconds. Sends a  request to get data from the data base
// the response is forwarded to the global channel if available
func getData() {
	var err error
	db, err = sql.Open("postgres", "host="+URL+" port="+PORT+" user="+USERNAME+" dbname="+DATABASE+" password="+PASSWORD+" sslmode=disable")
	if err != nil {
		status.SetCurrentValue("Connection error ")
		Stop <- true
		return
	}
	var resp Response
	// new ticker that ticks every UPDATE_INTERVAL milliseconds
	ticker := time.NewTicker(time.Millisecond * UPDATE_INTERVAL)
	defer ticker.Stop()
	for range ticker.C {
		resp.DataItems = nil
		rows, err := db.Query(QUERY)
		if err != nil {
			//some error happened
			status.SetCurrentValue("Could not reach Server, or SQL Statement is wrong. ")
			Data <- nil
			continue
		}

		for rows.Next() {
			var dataItem DataItem
			rows.Scan(&dataItem.Name, &dataItem.Value, &dataItem.Timestamp)
			resp.DataItems = append(resp.DataItems, dataItem)
		}

		Data <- &resp
	}
	Stop <- true
}

// StartSimulation is an endless loop that checks if we need to keep requesting data
// or stop the application.
func StartSimulation() {
	go getData()
	for {
		select {
		case val := <-Data:
			processData(val)
		case <-Stop:
			status.SetCurrentValue("Stop")
			return
		}
	}
}

// registerButton registers a start button that starts the connection to the database
func registerButton() {
	// get button reference
	btn := app.ClientDataModel().FindDataItemByName("refreshTriggered").(model.BooleanDataItem)
	btn.SetCurrentValue(false)

	// register adapter
	btn.AddAdapter(func(change api.FeatureChange) {
		val := btn.CurrentValue()
		// if button was pressed
		if val {
			// reset button state
			btn.SetCurrentValue(false)
			StartSimulation()
		}
	})
}

// registerStopButton registers a stop button that stops the connection to the database
func registerStopButton() {
	// get button reference
	btn := app.ClientDataModel().FindDataItemByName("stopTriggered").(model.BooleanDataItem)
	btn.SetCurrentValue(false)

	// register adapter
	btn.AddAdapter(func(change api.FeatureChange) {
		val := btn.CurrentValue()
		// if button was pressed
		if val {
			// reset button state
			btn.SetCurrentValue(false)
			Stop <- true
		}
	})
}
