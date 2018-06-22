package main

/* ====================================================================================================

cmdporter : a wifi intercom to talk to various devices

By Fred Ménez, Gaël Reyrol, Thierry Vo

==================================================================================================== */

/* TODO Serial

x looks for serial device depending on OS (Macos, Linux)
x discover serial device or read configuration
x load commands params from file

*/

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gophergala/cmdporter/vp/nec"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
)

var (
	SerialPortStatus bool = false
	g_Device         Device
)

type CmdRequest struct {
	Command string `json:"command"`
}

type CmdResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

func Render(w http.ResponseWriter, view string, content interface{}) {
	layout, err := ioutil.ReadFile(path.Join("views", "layout.html"))
	if err != nil {
		log.Fatal(err)
	}
	page, err := ioutil.ReadFile(path.Join("views", view))
	if err != nil {
		log.Fatal(err)
	}

	layoutTemplate := template.New("layout")
	pageTemplate := template.New("page")

	template.Must(layoutTemplate.Parse(string(layout)))
	template.Must(pageTemplate.Parse(string(page)))

	pageBuffer := new(bytes.Buffer)
	pageTemplate.Execute(pageBuffer, content)

	layoutContent := map[string]interface{}{"View": string(pageBuffer.Bytes())}
	layoutTemplate.Execute(w, layoutContent)
}

func ParseBody(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
	}

	return body
}

func main() {
	g_Device = nec.Nec_m271_m311
	LoadCommands(g_Device)

	// Start Http Server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		content := map[string]interface{}{
			"SerialPortStatus": SerialPortStatus,
			"Device":           g_Device.GetName(),
		}

		Render(w, "index.html", content)
	})

	http.HandleFunc("/cmd", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			req := CmdRequest{}
			res := CmdResponse{nil, nil}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			if err := json.Unmarshal(ParseBody(r), &req); err != nil {
				fmt.Println(err)
			}

			// Search for submited command on device
			commands := g_Device.GetCommandsList()
			if ok := commands[req.Command]; ok != nil {
				fmt.Printf("Found command : %s => %v\n", req.Command, commands[req.Command])
				g_Device.DoCmd(req.Command)
				res.Data = "Success"
				jsonRes, _ := json.Marshal(res)
				fmt.Fprintf(w, "%s", string(jsonRes))
				return
			}

			// Command not found in device commands list
			res.Error = "CommandNotFound"
			jsonRes, _ := json.Marshal(res)
			fmt.Fprintf(w, "%s", string(jsonRes))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Println("Running for device", g_Device.GetName())
	log.Println("Waiting for http connections on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func LoadCommands(d Device) {
	var err error

	// Load json file into string
	jsonbytes, err := ioutil.ReadFile(d.GetJsonPath())
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(-1)
	}
	json_string := string(jsonbytes)

	// Load it into our intermediate struct containing string encoded commands either in base 10 or hexa
	var IntermediateStruct = JSONCommands{}
	err = json.Unmarshal([]byte(json_string), &IntermediateStruct)
	if err != nil {
		fmt.Println("err :", err)
		os.Exit(-1)
	}

	// Convert these string encoded commands into bytes
	for key, value := range IntermediateStruct.Commands {
		command := value
		for _, cvalue := range command.StringCodedBytes {
			// TODO check whether string encoded commands actually begins with 0x, if not then it's base 10
			cmd_bytes, err := hex.DecodeString(cvalue[2:])
			if err != nil {
				fmt.Println("err :", err)
				os.Exit(-1)
			}
			// FIX this for commands containing more than one byte
			IntermediateStruct.Commands[key].Bytes = append(IntermediateStruct.Commands[key].Bytes, cmd_bytes[0])
		}
	}

	//CREATE A MAPPING FOR THE nec_m271_m311 COMMANDS
	for _, IntermediateCmd := range IntermediateStruct.Commands {
		d.RegisterCmd(IntermediateCmd.CommandName, IntermediateCmd.Bytes)
	}
	log.Println("test n :", IntermediateStruct.Name)
	d.SetName(IntermediateStruct.Name)

	log.Printf("Loaded %d commands for %s\n", d.GetNumCommands(), d.GetName())
}

type JSONCommands struct {
	Name     string
	Commands []JSONCommand
}

type JSONCommand struct {
	CommandName      string
	StringCodedBytes []string `json:"bytes"`
	Bytes            []byte
}
