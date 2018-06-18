package scripts

import (
	"database/sql"
	"fmt"
	"time"
	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	// import of application related stuff
	"git.monkey-works.de/scripting/api"
	"monkey-works.de/model"
)

// The endpoint of our web service we want to connect to
// Change this to your endpoint
const URL string = "tcp(127.0.0.1:3306)/schema" // tcp(yourHost_IP:Port)/Schema_name
const USERNAME string = "user"
const PASSWORD string = "1234"
const QUERY string = "SELECT * FROM variables WHERE id = ?" // your SQL Statement, "?" is a placeholder

// UPDATE_INTERVAL is the interval in which we refresh the Data from the lables
// the number is the interval in milliseconds.
const UPDATE_INTERVAL time.Duration = 2000

var db *sql.DB //you database

var id int

var app model.Application // app stuff, we'll need it for the DataItems

var status model.StringDataItem // Status Item to show errors and progress

// ok and Stop handle errors and stop the application if necessary
var ok bool
var Stop bool

/*
 _________________________________
|  id  |    name    |    value    |
|______|____________|_____________|
|  1   | DateItem1  |  HelloWorld |
|______|____________|_____________|
*/
// SQLData helps us to process the data that we get from the database.
type SQLData struct {
	id    int
	name  string
	value string
}

// @script
// InitializeScripting functions as the main method.
// model.Application allows us to access data from the Workbench and change it.
func InitializeScripting(application model.Application) {
	fmt.Println("Hello Scripting")
	// get status item
	status, ok = application.ClientDataModel().FindDataItemByName("status").(model.StringDataItem)
	if !ok {
		panic("DataItem \"status\" not found")
	}
	app = application
	// create start and stop button
	getStartButton()
	getStopButton()
}

// CheckConnectionToDataBase checks if there is a connection to the database.
func CheckConnectionToDataBase() bool {
	//ping the database to check if there is a connection
	err := db.Ping()
	if err != nil {
		status.SetCurrentValue("Failed to connect to database!")
		return false
	}
	status.SetCurrentValue("Connection ready to go")
	return true
}

// GetData connects to the database and checks if there is a connection to the database.
// It sends SQL query to the database and receives data from the server.
func GetData() {
	//infinite loop that reads the same data every single time (updates every UPDATE_INTERVAL milliseconds)
	// new ticker that ticks every UPDATE_INTERVAL milliseconds
	ticker := time.NewTicker(time.Millisecond * UPDATE_INTERVAL)
	defer ticker.Stop()
	for range ticker.C {
		// error for error handling
		var err error
		// Database connection
		// sql.Open("mysql","user:passwort@schema")
		db, err = sql.Open("mysql", USERNAME+":"+PASSWORD+"@"+URL)
		defer db.Close()
		if err != nil {
			status.SetCurrentValue("Wrong URL!")
			return
		}
		// check connection
		if !CheckConnectionToDataBase() {
			status.SetCurrentValue("Failed to communicate with the Database!")
			return
		}
		// if something went wrong or the Stop button was triggered we exit this method
		if Stop == true {
			return
		}
		// Execute MySQL Query
		sendSQLQueryToDataBaseAndReceiveData(1)
	}
}

// sendSQLQueryToDataBaseAndReceiveData sends SQL query to the Server and receives data.
func sendSQLQueryToDataBaseAndReceiveData(i int) {
	// struct to get data from SQL QUERY
	var sqlData SQLData
	// SQL Statement that gets executed
	stmtOut, err := db.Prepare(QUERY)
	if err != nil {
		status.SetCurrentValue("SQL Query failed, please check the SQL statement.")
		Stop = true
	}
	defer stmtOut.Close()
	// the real execution for the query
	err = stmtOut.QueryRow(i).Scan(&sqlData.id, &sqlData.name, &sqlData.value) // WHERE id=i
	if err != nil {
		status.SetCurrentValue("SQL Query failed, please check the SQL statement.")
		Stop = true
	} else {
		// Print the Data to show the user the result of the Query
		PrintData(sqlData.id, sqlData.name, sqlData.value)
	}
}

// PrintData shows the result of the query
// and prints Data into the labels.
func PrintData(id int, name string, value string) {
	// Print data into the DataItems
	// get DataItems
	valueDataItem := app.ClientDataModel().FindDataItemByName("yourDataItemName").(model.StringDataItem) // we have a DataItem for each id
	nameDataItem := app.ClientDataModel().FindDataItemByName("yourDataItemName").(model.StringDataItem)

	// set new data item values
	valueDataItem.SetCurrentValue(value)
	nameDataItem.SetCurrentValue(name)
	status.SetCurrentValue("Retrieved data successfully")
}

// StartSimulation is an endless loop.
// We start getting data the database and
// stop the application when "Stop" is true.
func StartSimulation() {
	// Start to get Data
	Stop = false
	go GetData()
	// infinite loop that checks the status of Stop, and keeps running as long as stop in not false
	for {
		switch {
		case Stop == false:
			// all good!
			status.SetCurrentValue("Receiving Data...")
		case Stop == true:
			// stop the application
			status.SetCurrentValue("Retrieved data successfully, stopped the connection to the database.")
		}
	}
}

// getStartButton action listener that starts the application when it gets triggered.
func getStartButton() {
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

// getStopButton is an action listener that stops the application when it gets triggered.
func getStopButton() {
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
			Stop = true
		}
	})
}
