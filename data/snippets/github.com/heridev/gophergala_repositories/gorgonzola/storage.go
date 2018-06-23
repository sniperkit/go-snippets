package gorgonzola

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Link contains Json-job link
type Link struct {
	URL     string
	Created time.Time
	Fetched time.Time
}

// Job contains single flattened job offer
type Job struct {
	LinkKey               string
	Hash                  string
	CompanyName           string
	CompanyURL            string
	CompanyRemoteFriendly bool
	CompanyMarket         string
	CompanySize           string
	Position              string
	Title                 string
	Description           string
	URL                   string
	Type                  string
	Posted                string
	Location              string
	Skills                []string
	SalaryRangeFrom       int
	SalaryRangeTo         int
	SalaryRangeCurrency   string
	EquityFrom            float32
	EquityTo              float32
	Perks                 []string
	Apply                 string
	Active                bool
	Created               time.Time
}

// Storage interface to handle data persistence
type Storage interface {
	AddURL(r *http.Request, url string) error
	GetJobs(r *http.Request, limit int) ([]Job, error)
	GetJob(r *http.Request, hash string) (*Job, error)
	Update(r *http.Request) error
}

func (j *Job) getHash() string {
	h := md5.New()
	io.WriteString(h, j.CompanyName)
	io.WriteString(h, j.Position)
	io.WriteString(h, j.Title)
	io.WriteString(h, j.URL)
	io.WriteString(h, j.Posted)
	return fmt.Sprintf("%x", h.Sum(nil))
}
