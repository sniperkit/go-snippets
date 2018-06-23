package code

import (
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

var (
	// AppFS represents the filesystem of the app. It is exported to be used as a
	// test helper.
	AppFS afero.Fs

	// ErrCodePathEmpty is returned if Code.Path is empty or invalid
	ErrCodePathEmpty = errors.New("code path is empty or does not exist")

	// ErrProfileNoFound is returned if the profile was not found
	ErrProfileNoFound = errors.New("profile not found")

	// ErrProjectNotFound is returned if the project is not found
	ErrProjectNotFound = errors.New("project not found")

	// ErrStoryNotFound is returned if the story is not found
	ErrStoryNotFound = errors.New("story not found")

	// ErrInvalidURL is returned by AddProject if the URL given is not valid
	ErrInvalidURL = errors.New("invalid URL given")

	// ErrProjectAlreadyExists is returned if the project already exists
	ErrProjectAlreadyExists = errors.New("project already exists")

	// ErrCoderNotScanned is returned if San() was never called
	ErrCoderNotScanned = errors.New("code was not scanned")

	// ErrPathIsInvalid is returned if the given absolute path is invalid (it
	// does not make up a full coder path + profile + story + project).
	ErrPathIsInvalid = errors.New("path is invalid")
)

func init() {
	// initialize AppFs to use the OS filesystem
	AppFS = afero.NewOsFs()
}

// code implements the coder interface
type code struct {
	// path is the base path of this profile
	path string

	// excludePattern is a list of patterns to ignore
	excludePattern *regexp.Regexp

	mu       sync.RWMutex
	profiles map[string]*profile
}

// New returns a new empty Code, caller must call Load to load from cache or
// scan the code directory
func New(p string, ignore *regexp.Regexp) Coder {
	return &code{
		path:           path.Clean(p),
		excludePattern: ignore,
		profiles:       make(map[string]*profile),
	}
}

// Path returns the absolute path of this coder
func (c *code) Path() string { return c.path }

// Profile returns the profile given it's name or an error if no profile with
// this name was found
func (c *code) Profile(name string) (Profile, error) { return c.getProfile(name) }

// Scan loads the code from the cache (if it exists), otherwise it will
// initiate a full scan and save it in cache.
func (c *code) Scan() error {
	// validate the Code, we cannot load an invalid Code
	if err := c.validate(); err != nil {
		return err
	}
	c.scan()

	return nil
}

// ProjectByAbsolutePath returns the project corresponding to the absolute
// path.
func (c *code) ProjectByAbsolutePath(p string) (Project, error) {
	// clean the path
	p = path.Clean(p)
	// trim the coder path from the path we are looking for (along with the pathSeparator)
	p = strings.TrimPrefix(p, c.Path()+string(os.PathSeparator))
	// split the path by the pathSeparator now
	parts := strings.Split(p, string(os.PathSeparator))
	if len(parts) < 5 { // 5 because the template is story/stories/STORY-NAME/src/project
		return nil, ErrPathIsInvalid
	}
	// get the profile
	profile, err := c.getProfile(parts[0])
	if err != nil {
		return nil, err
	}
	// get the story
	var story Story
	if parts[1] == baseStoryName {
		story = profile.Base()
		parts = parts[3:]
	} else if parts[1] == storiesDirName {
		story = profile.Story(parts[2])
		if !story.Exists() {
			return nil, ErrStoryNotFound
		}
		parts = parts[4:]
	}

	for i := len(parts); i > 0; i-- {
		var prj Project
		prj, err = story.Project(path.Join(parts[:i]...))
		if err == nil {
			return prj, nil
		}
	}

	return nil, ErrProjectNotFound
}

// getProfile return the profile identified by name
func (c *code) getProfile(name string) (*profile, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// make sure we scanned already
	if len(c.profiles) == 0 {
		return nil, ErrCoderNotScanned
	}
	// get the profile
	p, ok := c.profiles[name]
	if !ok {
		return nil, ErrProfileNoFound
	}
	return p, nil
}

// addProfile adds the profile to the list of profiles
func (c *code) addProfile(name string) *profile {
	// if the profile already exists, return it
	if p, err := c.getProfile(name); err == nil {
		return p
	}
	// create a new profile and make sure it has a base story
	p := newProfile(c, name)
	_, err := AppFS.Stat(p.Base().GoPath())
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	// add it to the map and return it
	c.mu.Lock()
	c.profiles[name] = p
	c.mu.Unlock()

	return p
}

// scan scans the entire profile to build the workspaces
func (c *code) scan() {
	// initialize the variables
	var wg sync.WaitGroup
	// read the profile and scan all profiles
	entries, err := afero.ReadDir(AppFS, c.path)
	if err != nil {
		log.Error().Str("path", c.path).Msgf("error reading the directory: %s", err)
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			// does this folder match the exclude pattern?
			if c.excludePattern != nil && c.excludePattern.MatchString(entry.Name()) {
				continue
			}
			// create the profile
			log.Debug().Msgf("found profile: %s", entry.Name())
			wg.Add(1)
			go func(name string) {
				if p := c.addProfile(name); p != nil {
					p.scan()
				}
				wg.Done()
			}(entry.Name())
		}
	}
	wg.Wait()
}

func (c *code) validate() error {
	if c.path == "" {
		return ErrCodePathEmpty
	}
	if _, err := AppFS.Stat(c.path); err != nil {
		return ErrCodePathEmpty
	}

	return nil
}
