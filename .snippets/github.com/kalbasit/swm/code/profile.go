package code

import (
	"path"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

const storiesDirName = "stories"

// profile represents the profile
type profile struct {
	// code links back to parent coder
	code *code

	// name is the name of the profile
	name string

	// stories is a list of workspaces
	mu      sync.RWMutex
	stories map[string]*story
}

func newProfile(c *code, name string) *profile {
	return &profile{
		name:    name,
		code:    c,
		stories: make(map[string]*story),
	}
}

// Coder returns the coder under which this exists
func (p *profile) Coder() Coder { return p.code }

// Name returns the name of the profile
func (p *profile) Name() string { return p.name }

// Base returns the base story
func (p *profile) Base() Story { return p.Story(baseStoryName) }

// Path returns the absolute path to this profile
func (p *profile) Path() string { return path.Join(p.code.Path(), p.name) }

// Story returns the story given it's name or an error if no story with this
// name was found
func (p *profile) Story(name string) Story { return p.getStory(name) }

func (p *profile) getStory(name string) *story {
	p.mu.RLock()
	s, ok := p.stories[name]
	p.mu.RUnlock()
	if !ok {
		return p.addStory(name)
	}
	return s
}

// addStory adds the story to the list of stories
func (p *profile) addStory(name string) *story {
	// if the story already exists, return it
	p.mu.RLock()
	s, ok := p.stories[name]
	p.mu.RUnlock()
	if ok {
		return s
	}
	// otherwise add it to the map
	s = newStory(p, name)
	p.mu.Lock()
	p.stories[name] = s
	p.mu.Unlock()

	return s
}

// scan scans the entire profile to build the workspaces
func (p *profile) scan() {
	// initialize the variables
	var wg sync.WaitGroup
	// scan the base story
	p.scanBaseStory(&wg)
	// scan the stories
	p.scanStories(&wg)

	wg.Wait()
}

func (p *profile) scanBaseStory(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		s := p.addStory(baseStoryName)
		s.scan()
		wg.Done()
	}()
}

func (p *profile) scanStories(wg *sync.WaitGroup) {
	// make sure the stories folder exists
	_, err := AppFS.Stat(path.Join(p.Path(), storiesDirName))
	if err != nil {
		return
	}
	// read the profile and scan all stories
	entries, err := afero.ReadDir(AppFS, path.Join(p.Path(), storiesDirName))
	if err != nil {
		log.Error().Str("path", path.Join(p.Path(), storiesDirName)).Msgf("error reading the directory: %s", err)
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			// create the story
			log.Debug().Str("profile", p.name).Msgf("found story: %s", entry.Name())
			wg.Add(1)
			go func(name string) {
				s := p.addStory(name)
				s.scan()
				wg.Done()
			}(entry.Name())
		}
	}
}
