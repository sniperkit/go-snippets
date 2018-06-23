package api

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

/*
Formats:

param key (value)
*/

type Params struct {
	URL     string
	Method  string
	Query   map[string][]string
	Headers map[string][]string
}

func NewParams(url string, buf io.Reader) (Params, error) {
	s := bufio.NewScanner(buf)
	p := Params{
		URL:     url,
		Query:   map[string][]string{},
		Headers: map[string][]string{},
	}

	for s.Scan() {
		line := s.Text()
		log.Printf("NewParams processing %#v", line)
		args := strings.Split(line, " ")
		p.parseLine(args)
	}

	err := s.Err()

	return p, err
}

func (p *Params) checkNumArgs(args []string, min int, max int) error {
	if len(args) < min || len(args) > max {
		return fmt.Errorf("expected %d or %d args, got %d", min, max, len(args))
	}

	return nil
}

func (p *Params) parseLine(args []string) error {
	if len(args) < 2 || len(args) > 3 {
		return fmt.Errorf("malformed argument %#v", args)
	}

	switch args[0] {
	case "method":
		if err := p.checkNumArgs(args, 2, 2); err != nil {
			return err
		}

		p.Method = args[1]
	case "query":
		if err := p.checkNumArgs(args, 2, 3); err != nil {
			return err
		}

		key := args[1]
		val := ""
		if len(args) == 3 {
			val = args[2]
		}

		p.Query[key] = append(p.Query[key], val)
		log.Printf("query: %#v", p.Query)

	case "header":
		if err := p.checkNumArgs(args, 2, 3); err != nil {
			return err
		}

		key := args[1]
		val := ""
		if len(args) == 3 {
			val = args[2]
		}

		p.Query[key] = append(p.Headers[key], val)
		log.Printf("query: %#v", p.Headers)

	}

	if p.Method == "" {
		p.Method = "GET"
	}

	return nil
}

func buildRequest(p Params) (*http.Request, error) {
	r, err := http.NewRequest(p.Method, p.URL, nil)
	if err != nil {
		return nil, err
	}

	for k, vs := range p.Headers {
		for _, v := range vs {
			r.Header.Add(k, v)
		}
	}

	q := url.Values{}
	for k, vs := range p.Query {
		for _, v := range vs {
			log.Printf("buildRequest adding %s %s", k, v)
			q.Add(k, v)
		}
	}

	r.URL.RawQuery = q.Encode()

	log.Printf("buildRequest final query: %#v", r.URL.Query())

	return r, nil
}

func DoRequest(p Params) (io.ReadCloser, error) {
	log.Printf("HTTP %#v %#v %#v", p.Method, p.URL, p.Query)
	req, err := buildRequest(p)
	if err != nil {
		return nil, err
	}

	var c http.Client

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
