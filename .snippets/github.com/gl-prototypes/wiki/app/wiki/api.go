package wiki

import "time"

// Store is an interface for storage backends
type Store interface {
	// Page gets a wiki page
	Page(title string) (*Page, error)

	// PutPage replaces a wiki page, renaming if necessary
	PutPage(title string, page *Page) error

	// Index returns a list of page titles
	Index() []string
}

// Page is a wiki page
type Page struct {
	Title   string
	Content string
	ModUser string
	ModTime time.Time
}

// Touch will update the last modified fields
// using the given user and the current time
func (p *Page) Touch(user string) {
	p.ModTime = time.Now()
	p.ModUser = user
}
