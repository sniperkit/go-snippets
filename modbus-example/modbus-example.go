package scripts

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
	// modbus library
	"github.com/goburrow/modbus"
	// app related stuff
	"git.monkey-works.de/scripting/api"
	"monkey-works.de/model"
)

// UPDATE_INTERVAL is the interval at which the connector is updated.
// Time is defined in milliseconds.
const UPDATE_INTERVAL time.Duration = 1000

// Data is a global channel to consume received data
var Data = make(chan *Response)

// client is a modbus TCP Client
var client = modbus.TCPClient("127.0.0.1:502") // (client_ip:port)

// Stop is a global channel to control periodic tasks
var Stop = make(chan bool)

// Stop_getData is a global channel that makes sure our Goroutine stops
var Stop_getData = make(chan bool)

// address is the address you would like to address on your modbus TCP Client
const address uint16 = 40000

var app model.Application

var status model.StringDataItem

// DataItem receives the data from the modbus TCP Client and writes it into an array of bytes
type DataItem struct {
	data []byte
}

// Response is used to handle the DataItems
type Response struct {
	DataItems []DataItem
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

// processData modifies the received data
func processData(data *Response) {
	if data == nil || len(data.DataItems) == 0 {
		status.SetCurrentValue("Empty response.")
		return
	}
	fmt.Println(data.DataItems[0].data)
	// transform the response to an int number
	var value uint16
	slice := []byte{data.DataItems[0].data[1], data.DataItems[0].data[0]}
	value = binary.LittleEndian.Uint16(slice)

	// Print data
	valueDataItem := app.ClientDataModel().FindDataItemByName("valueData").(model.IntegerDataItem)
	valueDataItem.SetCurrentValue(int32(value))

	// ADD YOUR CODE HERE
	// example: increment the value by 1 and write it to a holding register
	if value >= math.MaxUint16-1 {
		status.SetCurrentValue("Integer overflow detected! Value will be reset to 0.")
		value = 0
	}
	value++
	_, err := client.WriteSingleRegister(address, value)
	if err != nil {
		status.SetCurrentValue("Can't write to register.")
	}
}

// getData starts every UPDATE_INTERVAL milliseconds. Sends a  request to get data from the modbus TCP client
// the response is forwarded to the global channel if available
func getData() {
	var resp Response
	var dataItem DataItem
	// new ticker that ticks every UPDATE_INTERVAL milliseconds
	ticker := time.NewTicker(time.Millisecond * UPDATE_INTERVAL)
	defer ticker.Stop()
	for range ticker.C {
		resp.DataItems = nil
		// Read 10 values from address 0
		results, err := client.ReadHoldingRegisters(address, 1)
		if err != nil {
			//some error happened
			status.SetCurrentValue("Connection error")
			Data <- nil
			continue
		}
		dataItem.data = results
		resp.DataItems = append(resp.DataItems, dataItem)
		Data <- &resp
		// stop Goroutine if stop was pressed
		select {
		case <-Stop_getData:
			status.SetCurrentValue("It's true!")
			fmt.Println("Goroutine Stops.")
			return
		default:
			status.SetCurrentValue("Keep running")
		}
	}
}

// StartSimulation constantly checks the data from the channels
// Starts the getData() function
func StartSimulation() {
	go getData()
	status.SetCurrentValue("Connection started")
	for {
		select {
		case val := <-Data:
			processData(val)
		case <-Stop:
			status.SetCurrentValue("Stop")
			Stop_getData <- true
			return
		}
	}
}

func registerButton() {
	// get button reference
	btn := app.ClientDataModel().FindDataItemByName("refreshTriggered").(model.BooleanDataItem)
	btn.SetCurrentValue(false)

	// registerButton creates an action listener for the Start button, it calls the StartSimulation function
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

// registerStopButton creates an action listener for the Stop button and stops the connections to the modbus TCP Server
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
