package memory

import (
	"sync"

	"github.com/gl-prototypes/wiki/app/wiki"
	"github.com/gliderlabs/com"
)

func init() {
	com.Register(&Component{}, "")
}

type Component struct {
	pages map[string]*wiki.Page
	mu    sync.Mutex
}

func (c *Component) Page(title string) (*wiki.Page, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.pages[title], nil
}

func (c *Component) PutPage(title string, page *wiki.Page) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.pages == nil {
		c.pages = make(map[string]*wiki.Page)
	}
	if title != page.Title {
		delete(c.pages, title)
	}
	c.pages[page.Title] = page
	return nil
}

func (c *Component) Index() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	var titles []string
	for k := range c.pages {
		titles = append(titles, k)
	}
	return titles
}
