/*
	The git warehouse is a repository style warehouse.
	This means that the hashes given represent a given state of the repository
	instead of describing the state of a fileset.
	This git repository has several properties that are unique.
	 - You cannot write to it idempotently.
	 - Repositories are not interchangeable.
	 - Different hashes can result in equivalent filesets
	 - We cannot detect different forks of a repository.

	This warehouse requires git-upload-pack to be on the path.
	But otherwise is much less dependent on the host's git implementation
*/
package git

/*
TODO:
	- how to export files for transmat?
	- Figure out the correct interface for a repository controller
		- particularly string vs api.WareID for incoming hashes
*/

import (
	"context"
	"encoding/hex"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	. "github.com/warpfork/go-errcat"
	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	riofs "go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/warehouse"

	srcd_osfs "gopkg.in/src-d/go-billy.v4/osfs"
	srcd_git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/filemode"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
	"gopkg.in/src-d/go-git.v4/storage"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var (
	_ warehouse.RepositoryController = &Controller{}
	_ io.Reader                      = &Reader{}
)

const githubHostname = "github.com"
const gitmodulesFile = ".gitmodules"

// protocols
const (
	protocolHTTP  = "http"
	protocolHTTPS = "https"
	protocolFile  = "file"
)

/*
	Combo of remote repo interactions and a local (possibly on-disk) cache;
	these responsibilities are combined because in the case of "remote" repos
	that are actually also local, we'll use their object stores directly
	rather than replicating the entire object store again to our cache areas.
*/
type Controller struct {
	// user's address retained for messages (minus leading/trailing whitespace)
	addr string
	// Address that we will actually use to perform remote operations
	sanitizedAddr string
	// The detected protocol based on the address given
	protocol string
	// slugified address used for cache location
	slugifiedAddr string
	// filesystem in which to search for and place cached repositories, will use memory if nil
	workingDirectory    riofs.FS
	transportAuthMethod transport.AuthMethod
	store               storage.Storer // git object storage
	allowClone          bool           // clones are disabled for local files
	allowFetch          bool           // fetch is disabled a clone is fresh (vs using a cached clone)
	newClone            bool           // whether this controller created the clone in use
	repo                *srcd_git.Repository
}

/*
	Initialize a new warehouse controller that operates on a local filesystem.

	May return errors of category:

	  - `rio.ErrUsage` -- for unsupported addressses
	  - `rio.ErrWarehouseUnavailable` -- if the warehouse doesn't exist
	  - `rio.ErrLocalCacheProblem` -- if unable to create a cache directory
*/
func NewController(workingDirectory riofs.FS, addr api.WarehouseAddr) (*Controller, error) {
	var err error
	// Verify that the addr is sensible up front.
	sanitizedAddr, err := SanitizeRemote(string(addr))
	if err != nil {
		return nil, err
	}
	endpoint, err := transport.NewEndpoint(sanitizedAddr)
	if err != nil {
		return nil, Errorf(rio.ErrUsage, "failed to create endpoint after sanitization")
	}
	// Stamp out a warehouse handle.
	//  More values will be accumulated in shortly.
	whCtrl := &Controller{
		addr:             string(addr),
		sanitizedAddr:    sanitizedAddr,
		workingDirectory: workingDirectory,
		protocol:         endpoint.Protocol,
	}

	// ping the remote and see if it responds
	_, err = whCtrl.lsRemote()
	if err != nil {
		if whCtrl.protocol == protocolFile {
			// Unlike remote repositories, an error from a local repository pretty much means it doesn't exist
			return nil, ErrorDetailed(rio.ErrWarehouseUnavailable, "warehouse does not exist", map[string]string{
				"cause":     err.Error(),
				"warehouse": sanitizedAddr,
			})
		}
		return nil, Errorf(rio.ErrWarehouseUnavailable, "warehouse unavailable: %s", err)
	}
	err = whCtrl.setCacheStorage()
	return whCtrl, err
}

/*
	Returns the commit for the given hash
*/
func (c *Controller) GetCommit(hash string) (*object.Commit, error) {
	commitHash, err := StringToHash(hash)
	if err != nil {
		return nil, err
	}
	commit, err := object.GetCommit(c.store, commitHash)
	if err == plumbing.ErrObjectNotFound {
		return nil, Errorf(rio.ErrWareNotFound, "commit not found")
	} else if err != nil {
		return nil, Errorf(rio.ErrWareCorrupt, "failed to get commit: %s", err)
	}
	return commit, nil
}

/*
	Returns true if the repository contains the commit for the given hash
*/
func (c *Controller) Contains(hash string) bool {
	commitHash, err := StringToHash(hash)
	if err != nil {
		return false
	}
	commit, err := object.GetCommit(c.store, commitHash)
	if err != nil {
		return false
	}
	return commit != nil
}

/*
	Returns the tree for a given commit hash
*/
func (c *Controller) GetTree(hash string) (*object.Tree, error) {
	commit, err := c.GetCommit(hash)
	if err != nil {
		return nil, err
	}
	tree, err := commit.Tree()
	if err != nil {
		return nil, Errorf(rio.ErrWareCorrupt, "commit missing tree: %s", err)
	}
	return tree, nil
}

/*
	This sets the cache storage for the controller.
	If the repository is local then we can just open the repository. No cache is needed.
	If a working directory is set, then we can clone repositories there and use them again later.
	If a working directory is _not_ set, then we must clone into memory.
*/
func (c *Controller) setCacheStorage() error {
	var err error
	if c.store != nil {
		return nil
	}
	if c.protocol == protocolFile {
		// if we are pulling from a local repository then no cache is needed at all! Just use the repo itself.
		if !filepath.IsAbs(c.sanitizedAddr) {
			return Errorf(rio.ErrUsage, "remote is not an absolute path")
		}
		c.store, err = filesystem.NewStorage(srcd_osfs.New(c.sanitizedAddr))
		if err != nil {
			return Errorf(rio.ErrWareCorrupt, "Could not open repository directory: %s", err)
		}
		return nil
	}
	c.allowClone = true // non-local repositories are allowed to clone
	if c.workingDirectory == nil {
		// if no working directory is given then use memory to clone
		c.store = memory.NewStorage()
		return nil
	}
	c.allowFetch = true // cached repositories may be updated
	workingDir := c.workingDirectory.BasePath()
	remotePath := SlugifyRemote(c.sanitizedAddr)
	path := workingDir.Join(riofs.MustRelPath(remotePath))
	srcd_fs := srcd_osfs.New(path.String())
	c.store, err = filesystem.NewStorage(srcd_fs)
	if err != nil {
		return Errorf(rio.ErrLocalCacheProblem, "Could not create repository cache: %s", err)
	}
	return nil
}

func (c *Controller) Clone(ctx context.Context) error {
	var err error
	c.repo, err = c.open(ctx, c.store, c.allowClone)
	return err
}

func (c *Controller) Update(ctx context.Context) error {
	if c.repo == nil {
		panic("cannot update repository before opening")
	}
	if c.allowFetch && !c.newClone {
		return c.repo.FetchContext(ctx, &srcd_git.FetchOptions{
			// Auth credentials, if required, to use with the remote repository.
			Auth: c.transportAuthMethod,
		})
	}
	return nil
}

/*
	Opens the repository or clones it if it is does not exist
	Will not clone if allowClone is false.
	If a clone occurs then controller.newClone will be set to true.

	I'm not a huge fan of the way this function does things.
	It's too complex and behaves sort of as a method and sort of as a function.
	It makes some tests easier but not in any way that you can't do the same
	from the calling function. I think it would be better to move all this to Clone
	and delete it or make it a function entirely independent of controller.
*/
func (c *Controller) open(ctx context.Context, store storage.Storer, allowClone bool) (*srcd_git.Repository, error) {
	repo, err := srcd_git.Open(store, nil)
	if err == srcd_git.ErrRepositoryNotExists && allowClone {
		// checkout dir is nil for a bare clone
		//FIXME: this should use the endpoint that we created. It has _special_ behaviour
		repo, err = srcd_git.CloneContext(ctx, store, nil, &srcd_git.CloneOptions{
			// The (possibly remote) repository URL to clone from.
			URL: c.sanitizedAddr,
			// Auth credentials, if required, to use with the remote repository.
			Auth: c.transportAuthMethod,
			// we handle submodule recursion ourselves.
			RecurseSubmodules: srcd_git.NoRecurseSubmodules,
		})
		if ctx.Err() != nil {
			return nil, Errorf(rio.ErrCancelled, "cancelled: %s", err)
		} else if err != nil {
			return nil, Errorf(rio.ErrWareCorrupt, "unable to clone repository: %s", err)
		}
		c.newClone = true
		return repo, nil
	} else if err != nil {
		return nil, Errorf(rio.ErrLocalCacheProblem, "unable to open cache repository: %s", err)
	}
	return repo, nil
}

/*
	Pretty straight forward `git ls-remote` implementation
	Returns the list of references available on the remote
*/
func (c *Controller) lsRemote() (memory.ReferenceStorage, error) {
	endpoint, err := transport.NewEndpoint(c.sanitizedAddr)
	if err != nil {
		return nil, err
	}
	gitClient, err := client.NewClient(endpoint)
	if err != nil {
		return nil, err
	}
	gitSession, err := gitClient.NewUploadPackSession(endpoint, nil)
	if err != nil {
		return nil, err
	}
	advertisedRefs, err := gitSession.AdvertisedReferences()
	if err != nil {
		return nil, err
	}
	refs, err := advertisedRefs.AllReferences()
	if err != nil {
		return nil, err
	}
	err = gitSession.Close()
	if err != nil {
		return nil, err
	}
	return refs, nil
}

/*
	We force https when talking to github because github may refuse to respond to http urls.
	Otherwise this will return an endpoint that will use the correct protocol based on the remote.
	We could improve behavior by overriding installing new protocols.
	Perhaps using a particular git user agent would avoid the strange github behavior.
*/
func SanitizeRemote(remote string) (string, error) {
	var err error
	remote = strings.TrimSpace(remote)
	if remote == "" {
		return "", Errorf(rio.ErrUsage, "empty git remote")
	}

	endpoint, err := transport.NewEndpoint(remote)
	if err != nil {
		return "", Errorf(rio.ErrUsage, "failed to parse URI: %s", err)
	}

	protocol := endpoint.Protocol
	if protocol == protocolFile {
		// absolutize paths
		if HasFoldedPrefix(remote, "file://") {
			if len(remote) <= 7 {
				return "", Errorf(rio.ErrUsage, "empty git remote")
			}
			remote = remote[7:]
		}
		if !filepath.IsAbs(remote) {
			remote, err = filepath.Abs(remote)
			if err != nil {
				return "", Errorf(rio.ErrUsage, "failed handling local path")
			}
			return SanitizeRemote(remote)
		}
	} else if protocol == protocolHTTP {
		// github will not send back a response over http so we force https in this case
		if HasFoldedSuffix(endpoint.Host, githubHostname) {
			parsedUrl, err := url.Parse(remote)
			if err != nil {
				return "", Errorf(rio.ErrUsage, "failed to parse URI: %s", err)
			}
			parsedUrl.Scheme = protocolHTTPS
			return SanitizeRemote(parsedUrl.String())
		}
	}
	// using the git library may help mitigate things like ssh commands being executed
	// but I haven't tested it and it's easy enough to say no
	pathString := endpoint.Host + endpoint.Path
	if pathString == "" {
		return "", Errorf(rio.ErrUsage, "warehouse has empty path: %s", endpoint.String())
	} else if pathString[0] == '-' {
		return "", Errorf(rio.ErrUsage, "warehouse host cannot start with '-'")
	}
	return remote, nil
}

/*
	Combination of strings.EqualFold and strings.HasSuffix
*/
func HasFoldedSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && strings.EqualFold(s[len(s)-len(suffix):], suffix)
}

/*
	Combination of strings.EqualFold and strings.HasSuffix
*/
func HasFoldedPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}

/*
	Transform the rio commit ID to a git hash.
	Performs some basic checks on inputs.
*/
func StringToHash(hash string) (plumbing.Hash, error) {
	if err := mustBeFullHash(hash); err != nil {
		return plumbing.Hash{}, err
	}
	return plumbing.NewHash(hash), nil
}

/*
	A git hash must be exactly 40 hex characters
*/
func mustBeFullHash(hash string) error {
	if len(hash) != 40 {
		return Errorf(rio.ErrUsage, "git commit hashes are 40 characters")
	}
	if _, err := hex.DecodeString(hash); err != nil {
		return Errorf(rio.ErrUsage, "git commit hashes are hex strings")
	}
	return nil
}

/*
	Return a string that's safe to use as a dir name.

	Uses URL query escaping so it remains roughly readable.
	Does not attempt any URL normalization.
*/
func SlugifyRemote(remoteURL string) string {
	return url.QueryEscape(remoteURL)
}

type Submodule struct {
	// Name module name
	Name string
	// Path defines the path, relative to the top-level directory of the Git
	// working tree.
	Path string
	// URL defines a URL from which the submodule repository can be cloned.
	URL string
	// Branch is a remote branch name for tracking updates in the upstream
	// submodule. Optional value.
	Branch string
	// Ye olde hash!
	Hash string
}

/*
	Returns the submodule map for a given commit hash.

	Keys in the map are submodules "names" as given in the committed
	'.gitmodules' file; values contain the repeated name, the path,
	the url from the '.gitmodules' file, and the hash as a string.

	If a path specified in '.gitmodules' does not have a gitlink (i.e., no
	real hash is there, or, some other kind of object is in the repo at that
	path), we'll return an error.
	It's also quite possible that there will be *more* gitlinks in the tree;
	we won't notice this here, because this function does not do a full treewalk.
*/
func (c *Controller) Submodules(commitHash string) (map[string]Submodule, error) {
	result := map[string]Submodule{}

	commit, err := c.GetCommit(commitHash)
	if err != nil {
		return result, err
	}

	// If there is no gitmodules file then we can return nothing.
	// FIXME: I'm not sure what the chance of getting _other_ errors here is.
	// As far as I can tell it's unlikely but idk. This may need better handling.
	fd, err := commit.File(gitmodulesFile)
	if err != nil {
		return result, nil
	}

	data, err := fd.Contents()
	if err != nil {
		return result, Errorf(rio.ErrWareCorrupt, "found but could not read %s", gitmodulesFile)
	}

	cfgModules := config.NewModules()
	err = cfgModules.Unmarshal([]byte(data))
	if err != nil {
		return result, Errorf(rio.ErrWareCorrupt, "found but could not parse %s", gitmodulesFile)
	}
	if cfgModules == nil {
		return result, nil
	}

	/*
		Due diligence for tracking submodules.
		We expect the git submodule file to correspond to valid submodules entries in the commit tree.
		We are going to return this mapping of file entries to tree entries.
	*/
	tree, err := commit.Tree()
	if err != nil {
		return result, Errorf(rio.ErrWareCorrupt, "commit missing tree")
	}
	for name, submodule := range cfgModules.Submodules {
		if submodule == nil {
			// Should protect against incomplete entries
			return result, Errorf(rio.ErrWareCorrupt, "nil submodule")
		}

		entry, err := tree.FindEntry(submodule.Path)
		if err != nil {
			return result, Errorf(rio.ErrWareCorrupt, "submodule entry missing matching tree entry")
		}

		if entry.Mode != filemode.Submodule {
			return result, Errorf(rio.ErrWareCorrupt, "gitmodule entry is not a submodule")
		}

		// Mostly a translation struct, but adds the hash.
		result[name] = Submodule{
			Name:   submodule.Name,
			Path:   submodule.Path,
			URL:    submodule.URL,
			Branch: submodule.Branch,
			Hash:   entry.Hash.String(),
		}
	}
	return result, nil
}

/*
	Creates a reader/iterator for extracting files from a particular commit.
*/
func (c *Controller) NewReader(commitHash string) (*Reader, error) {
	tree, err := c.GetTree(commitHash)
	if err != nil {
		return nil, err
	}
	return &Reader{tree.Files(), nil}, nil
}

func (c *Controller) NewTreeWalker(commitHash string) (*object.TreeWalker, error) {
	tree, err := c.GetTree(commitHash)
	if err != nil {
		return nil, err
	}
	return object.NewTreeWalker(tree, true, nil), nil
}

/*
	Reader is a utility for extracting files from the repository.
	It is loosely based on the `archive/tar` package.
*/
type Reader struct {
	iter    *object.FileIter
	current *object.File
}

/*
	Closes the iterator
*/
func (r *Reader) Close() {
	r.iter.Close()
}

/*
	Moves the reader to the next file and returns the files header information
*/
func (r *Reader) Next() (*Header, error) {
	var err error
	r.current, err = r.iter.Next()
	if err != nil {
		return nil, err
	}
	isBinary, err := r.current.IsBinary()
	if err != nil {
		return nil, err
	}
	header := Header{
		Name:   r.current.Name,
		Mode:   int64(r.current.Mode),
		Size:   r.current.Size,
		Binary: isBinary,
	}
	return &header, err
}

/*
	Reads the current file contents.
*/
func (r *Reader) Read(b []byte) (int, error) {
	reader, err := r.current.Reader()
	if err != nil {
		return 0, err
	}
	return reader.Read(b)
}

/*
	File header information
*/
type Header struct {
	Name   string // name of header file entry
	Mode   int64  // permission and mode bits
	Size   int64  // length in bytes
	Binary bool
}
