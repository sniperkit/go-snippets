// main.go for figdash web tool
// as no API in Fig, built up Fig commands with commend line calls to Fig
//

package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/codegangsta/cli"
)

var project = ""

var headerTempl = template.Must(template.New("dash").Parse(headerStr))
var footerTempl = template.Must(template.New("dash").Parse(footerStr))

func fixProjectName(c *cli.Context) string {
	if c.String("projectname") == "" {
		var cwd, err = os.Getwd()
		if err == nil {
			var wd = path.Base(cwd)
			fmt.Printf("cwd = %s\n", cwd)
			fmt.Printf("wd = %s\n", wd)
			project = wd
			return wd
		}
		project = "unknown"
		return ""
	}
	project = c.String("projectname")
	return c.String("projectname")
}

func handler(w http.ResponseWriter, req *http.Request) {
	headerTempl.Execute(w, req.FormValue("s"))
	fmt.Fprintf(w, "<h3> for project %s </h3>", project)

	// now show each container as a box?
	cmd := exec.Command("fig", "ps")
	cmd.Stdin = os.Stdin
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	// start the command after having set up the pipe
	if err := cmd.Start(); err != nil {
		return
	}

	// read command's stdout line by line
	in := bufio.NewScanner(stdout)

	fmt.Fprintf(w, "<pre><samp>")
	for in.Scan() {
		fmt.Fprintf(w, in.Text()) // write each line to your log, or anything you need
		fmt.Fprintf(w, "<br/>")
	}
	if err := in.Err(); err != nil {
		log.Printf("error: %s", err)
	}
	cmd.Stderr = os.Stderr
	//cmd.Run()
	fmt.Fprintf(w, "</samp></pre>")

	footerTempl.Execute(w, req.FormValue("s"))
}

func killHandler(w http.ResponseWriter, req *http.Request) {
	headerTempl.Execute(w, req.FormValue("s"))
	fmt.Fprintf(w, "<h3> for project %s </h3>", project)
	// now kill all fig services
	fmt.Fprintf(w, "<h1>KILLing all services for %s</h1>", project)

	footerTempl.Execute(w, req.FormValue("s"))
}

func stopHandler(w http.ResponseWriter, req *http.Request) {
	headerTempl.Execute(w, req.FormValue("s"))
	fmt.Fprintf(w, "<h3> for project %s </h3>", project)
	// now stop each container
	fmt.Fprintf(w, "<h1>STOPing all services for %s</h1>", project)

	footerTempl.Execute(w, req.FormValue("s"))
}

func startHandler(w http.ResponseWriter, req *http.Request) {
	headerTempl.Execute(w, req.FormValue("s"))
	fmt.Fprintf(w, "<h3> for project %s </h3>", project)
	// now start each container
	fmt.Fprintf(w, "<h1>STARTing all services for %s</h1>", project)

	footerTempl.Execute(w, req.FormValue("s"))
}

func main() {
	app := cli.NewApp()
	app.Name = "figdash"
	app.Usage = "fig dashboard"
	app.Version = "0.0.1"
	app.Email = "mkobar@rkosecurity.com"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Show more output",
		},
		// version flag support is builtin
		//cli.BoolFlag{
		//	Name:  "version",
		//	Usage: "Print version and exit",
		//},
		cli.StringFlag{
			Name:   "file, f",
			Value:  "fig.yml",
			Usage:  "Specify an alternate fig file",
			EnvVar: "FIG_FILE",
		},
		cli.StringFlag{
			Name:   "projectname, p",
			Value:  "notset",
			Usage:  "Specify an alternate project name",
			EnvVar: "FIG_PROJECT_NAME",
		},
	}
	app.Commands = []cli.Command{
		// build - NOT supported yet
		// help - NOT supported yet
		{
			Name:  "kill",
			Usage: "Force stop service containers.",
			Action: func(c *cli.Context) {
				cmd := exec.Command("fig", "kill")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			},
		},
		// logs - not supported
		// port - NOT supported
		{
			Name:  "ps",
			Usage: "List containers",
			Action: func(c *cli.Context) {
				var pn = fixProjectName(c)
				if pn != "" {
					fmt.Printf("ProjectName: %s\n", pn)
				} else {
					fmt.Printf("ProjectName: %s\n", "unknown")
				}
				cmd := exec.Command("fig", "ps")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			},
		},
		// pull - NOT supported
		{
			Name:  "rm",
			Usage: "Remove stopped service containers.",
			Action: func(c *cli.Context) {
				cmd := exec.Command("fig", "rm")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			},
		},
		// run - NOT supported
		// scale - NOT supported
		{
			Name:  "start",
			Usage: "Start existing containers for a service.",
			Action: func(c *cli.Context) {
				cmd := exec.Command("fig", "start")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			},
		},
		{
			Name:  "stop",
			Usage: "Stop existing containers without removing them.",
			Action: func(c *cli.Context) {
				cmd := exec.Command("fig", "stop")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			},
		},
		// up - NOT supported - due to logs
		{
			Name:  "web",
			Usage: "Enable web monitoring on http://localhost/1984/PROJECTNAME.",
			Action: func(c *cli.Context) {
				var pn = fixProjectName(c)
				if pn == "" {
					pn = "unknown"
				}
				fmt.Printf("Starting web Fig Dashboard at http://localhost:1984/%s\n", pn)
				fmt.Printf("use Ctrl-C to exit\n")
				http.HandleFunc("/"+pn, handler) // default
				http.HandleFunc("/"+pn+"/kill", killHandler)
				http.HandleFunc("/"+pn+"/stop", stopHandler)
				http.HandleFunc("/"+pn+"/start", startHandler)
				err := http.ListenAndServe(":1984", nil)
				if err != nil {
					log.Fatal("ListenAndServe:", err)
				}
			},
		},
		// up - NOT supported - due to logs
	}
	app.Run(os.Args)
}

const headerStr = `
<html>
<head>
<title>Fig Dashboard</title>
</head>
<body>
<h2>Fig Dashboard v 0.0.1</h2>
`
const footerStr = `
</body>
</html>
`
