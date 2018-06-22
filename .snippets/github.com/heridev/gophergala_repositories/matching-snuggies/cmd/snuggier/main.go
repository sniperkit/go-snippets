/*
Command snuggier is a command line 3D slicing application that converts 3D
models to G-code for 3D printing using a snuggied server.

	snuggier -o model.gcode model.stl

Call snuggier with the -h flag to see available command line configuration.

	snuggier -h
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gophergala/matching-snuggies/slicerjob"
)

func main() {
	server := flag.String("server", "localhost:8888", "snuggied server address")
	slicerBackend := flag.String("backend", "slic3r", "backend slicer")
	slicerPreset := flag.String("preset", "hq", "specify a configuration preset for the backend")
	presets := flag.Bool("L", false, "get list of available configuration presets for Slic3r")
	gcodeDest := flag.String("o", "", "specify an output gcode filename")
	flag.Parse()

	client := &Client{
		ServerAddr: *server,
	}

	if *presets == true {
		presets, err := client.SlicerPresets()
		if err != nil {
			fmt.Errorf("something bad happened: %v", err)
		}
		for i := range presets {
			fmt.Println(presets[i])
		}
		return
	}

	if flag.NArg() < 1 {
		log.Fatalf("missing argument: mesh file")
	}
	meshpath := flag.Arg(0)

	// start intercepting signals from the operating system
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// send files to the slicer to be printed and poll the slicer until the job
	// has completed.
	log.Printf("sending file(s) to snuggied server at %v", *server)
	job, err := client.SliceFile(*slicerBackend, *slicerPreset, meshpath)
	if err != nil {
		log.Fatalf("sending files: %v", err)
	}

	// poll the server until the job has completed.  use exponential backoff to
	// reduce spam for slice slicing jobs.
	maxTick := time.Second * 5
	currentTick := 100 * time.Millisecond
	tick := time.After(currentTick)
	for job.Status != slicerjob.Complete {
		select {
		case s := <-sig:
			// stop intercepting signals. if the job cancellation is taking too
			// long let the future signals terminate the process naturally.
			signal.Stop(sig)
			log.Printf("signal: %v", s)
			err := client.Cancel(job)
			if err != nil {
				log.Printf("failed to cancel job: %v", err)
			}
			log.Printf("slicing job canceled")
			return
		case <-tick:
			job, err = client.SlicerStatus(job)
			if err != nil {
				// TODO:
				// detect potentially intermittent network failures and
				// continue polling up to some reasonable time limit.
				log.Fatalf("waiting: %v", err)
			}

			if currentTick < maxTick {
				currentTick *= 2
			}
			if currentTick > maxTick {
				currentTick = maxTick
			}
			tick = time.After(currentTick)
		}
	}
	// stop intercepting signals because it because much more difficult to stop
	// gracefully while reading gcode from the server.
	signal.Stop(sig)

	// download gcode from the slicer and write to the specified file.
	log.Printf("retreiving gcode file")
	r, err := client.GCode(job)
	if err != nil {
		log.Fatalf("gcode: %v", err)
	}
	defer r.Close()
	var f *os.File
	if *gcodeDest == "" {
		f = os.Stdout
	} else {
		f, err = os.Create(*gcodeDest)
		if err != nil {
			log.Panic(err)
		}
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Panic(err)
		}
	}()
	_, err = io.Copy(f, r)
	if err != nil {
		log.Panic(err)
	}
}

type Client struct {
	Client     *http.Client
	ServerAddr string
	HTTPS      bool
}

// SliceFiles tells the server to slice the specified paths.
func (c *Client) SliceFile(backend, preset string, path string) (*slicerjob.Job, error) {
	// check that a mesh file is given as the first argument and open it
	// so it may to encode in the form.
	if !IsMeshFile(path) {
		log.Fatalf("path is not a mesh file: %v", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// write the multipart form out to a temporary file.  the temporary
	// file is closed and unlinked when the function terminates.
	tmp, err := ioutil.TempFile("", "matching-snuggies-post-")
	if err != nil {
		return nil, fmt.Errorf("tempfile: %v", err)
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()
	bodyw := multipart.NewWriter(tmp)
	err = c.writeJobForm(bodyw, backend, preset, path, f)
	if err != nil {
		return nil, fmt.Errorf("tempfile: %v", err)
	}

	// seek back to the beginning of the form and POST it to the slicer
	// server.  decode a slicerjob.Job from successful responses.
	var job *slicerjob.Job
	_, err = tmp.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("tempfile: %v", err)
	}
	url := c.url("/slicer/jobs")
	log.Printf("POST %v", url)
	resp, err := c.client().Post(url, bodyw.FormDataContentType(), tmp)
	if err != nil {
		return nil, fmt.Errorf("POST /slicer/jobs: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return nil, httpStatusError(resp)
	}
	err = json.NewDecoder(resp.Body).Decode(&job)
	if err != nil {
		return nil, fmt.Errorf("response: %v", err)
	}
	return job, nil
}

func (c *Client) writeJobForm(w *multipart.Writer, backend, preset, filename string, r io.Reader) error {
	err := w.WriteField("slicer", backend)
	if err != nil {
		return err
	}
	err = w.WriteField("preset", preset)
	if err != nil {
		return err
	}
	file, err := w.CreateFormFile("meshfile", filepath.Base(filename))
	if err != nil {
		return err
	}
	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}
	return w.Close()
}

func (c *Client) Cancel(job *slicerjob.Job) error {
	if job.ID == "" {
		return fmt.Errorf("job missing id")
	}
	url := c.url("/slicer/jobs/" + job.ID)
	log.Printf("DELETE %v", url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("request: %v", err)
	}
	resp, err := c.client().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return httpStatusError(resp)
	}
	return nil
}

func (c *Client) SlicerPresets() ([]string, error) {
	url := c.url("/slicer/presets/slic3r")
	//log.Printf("GET %v", url)
	resp, err := c.client().Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET /slicer/presets/slic3r: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, httpStatusError(resp)
	}
	preset := new(slicerjob.SlicerPreset)
	err = json.NewDecoder(resp.Body).Decode(preset)
	if err != nil {
		return nil, fmt.Errorf("GET /slicer/presets/slic3r: %v", err)
	}

	return preset.Presets, nil
}

// SlicerStatus returns a current copy of the provided job.
func (c *Client) SlicerStatus(job *slicerjob.Job) (*slicerjob.Job, error) {
	if job.ID == "" {
		return nil, fmt.Errorf("job missing id")
	}
	var jobcurr *slicerjob.Job
	url := c.url("/slicer/jobs/" + job.ID)
	log.Printf("GET %v", url)
	resp, err := c.client().Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET /slicer/jobs/: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, httpStatusError(resp)
	}
	err = json.NewDecoder(resp.Body).Decode(&jobcurr)
	if err != nil {
		return nil, fmt.Errorf("response: %v", err)
	}
	return jobcurr, nil
}

// GCode requests the gcode for job.
func (c *Client) GCode(job *slicerjob.Job) (io.ReadCloser, error) {
	url := c.url("/slicer/gcodes/" + job.ID)
	log.Printf("GET %v", url)
	resp, err := c.client().Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET /slicer/codes/: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, httpStatusError(resp)
	}
	return resp.Body, nil
}

func (c *Client) client() *http.Client {
	if c.Client == nil {
		return http.DefaultClient
	}
	return c.Client
}

func (c *Client) url(pathquery string) string {
	pathquery = strings.TrimPrefix(pathquery, "/")
	scheme := "http"
	if c.HTTPS {
		scheme = "https"
	}
	return scheme + "://" + c.ServerAddr + "/" + pathquery
}

var meshExts = map[string]bool{
	".stl": true,
	".amf": true,
}

func IsMeshFile(path string) bool {
	return meshExts[filepath.Ext(path)]
}

func httpStatusError(resp *http.Response) error {
	p, _ := ioutil.ReadAll(io.LimitReader(resp.Body, 85))
	msg := trimMessage(string(p), 80)
	return fmt.Errorf("http %d %s: %q", resp.StatusCode, http.StatusText(resp.StatusCode), msg)
}

func trimMessage(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) < n {
		return s
	}
	var rs []rune
	var m int
	for _, c := range s {
		if m >= n {
			break
		}
		rs = append(rs, c)
		m++
	}
	return string(rs) + "..."
}
