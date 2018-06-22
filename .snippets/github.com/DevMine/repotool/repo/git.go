// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repo

import (
	"archive/tar"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	g2g "github.com/libgit2/git2go"

	"github.com/DevMine/repotool/config"
	"github.com/DevMine/repotool/model"
)

// gitRepo is a repository with some things specific to git.
type gitRepo struct {
	model.Repository
	cfg    config.DataConfig
	r      *g2g.Repository
	tmpDir string
}

// New creates a new gitRepo object.
func newGitRepo(cfg config.DataConfig, repository model.Repository, gitDir string, useTmpDir bool) (*gitRepo, error) {
	r, err := g2g.OpenRepository(gitDir)
	if err != nil {
		return nil, err
	}

	var tmpDir string
	if useTmpDir {
		tmpDir = gitDir
	}

	return &gitRepo{Repository: repository, cfg: cfg, r: r, tmpDir: tmpDir}, nil
}

// FetchCommits fetches all commits from a Git repository and adds them to
// the list of commits of the repository object.
func (gr *gitRepo) FetchCommits() error {
	gr.Commits = make([]model.Commit, 0) // give number of commits

	rw, err := gr.r.Walk()
	if err != nil {
		return err
	}

	err = rw.PushHead()
	if err != nil {
		return err
	}

	err = rw.Iterate(gr.addCommit)
	if err != nil {
		return err
	}

	return nil
}

// GetRepository returns the repository structre contained in a git repository.
func (gr gitRepo) GetRepository() *model.Repository {
	return &gr.Repository
}

// GetName returns the name of a git repository.
func (gr gitRepo) GetName() string {
	return gr.Name
}

// GetVCS returns the VCS type (shall be "git").
func (gr gitRepo) GetVCS() string {
	return gr.VCS
}

// GetCloneURL returns the git repository clone URL.
func (gr gitRepo) GetCloneURL() string {
	return gr.CloneURL
}

// GetClonePath returns the clone path of a git repository.
func (gr gitRepo) GetClonePath() string {
	return gr.ClonePath
}

// GetDefaultBranch returns the git repository default branch.
func (gr gitRepo) GetDefaultBranch() string {
	return gr.DefaultBranch
}

// GetCommits returns the list of commits in the git repository.
// If the list is empty of nil, this probably means that a call to
// FetchCommits() is needed to populate the list.
func (gr gitRepo) GetCommits() []model.Commit {
	return gr.Commits
}

// Cleanup frees open repositories and removes temporary created files, if any.
func (gr gitRepo) Cleanup() error {
	if gr.r != nil {
		gr.r.Free()
	}

	if len(gr.tmpDir) > 0 {
		return os.RemoveAll(gr.tmpDir)
	}
	return nil
}

var deltaMap = map[g2g.Delta]*string{
	g2g.DeltaUnmodified: nil,
	g2g.DeltaAdded:      &model.StatusAdded,
	g2g.DeltaDeleted:    &model.StatusDeleted,
	g2g.DeltaModified:   &model.StatusModified,
	g2g.DeltaRenamed:    &model.StatusRenamed,
	g2g.DeltaCopied:     &model.StatusCopied,
	g2g.DeltaIgnored:    nil,
	g2g.DeltaUntracked:  nil,
	g2g.DeltaTypeChange: nil,
}

// addCommit is conform to the g2g.RevWalIterator type in order to be used
// by the g2g.Iterate() function to iterate over all commits of a Git repository.
func (gr *gitRepo) addCommit(c *g2g.Commit) bool {
	if c == nil {
		return false
	}

	var commit model.Commit

	oID := c.Id()
	if oID == nil {
		return false
	}
	commit.VCSID = oID.String()

	commit.Message = c.Message()

	var author model.Developer
	author.Name = c.Author().Name
	author.Email = c.Author().Email
	commit.Author = author

	var committer model.Developer
	committer.Name = c.Committer().Name
	committer.Email = c.Committer().Email
	commit.Committer = committer

	commit.CommitDate = c.Committer().When
	commit.AuthorDate = c.Author().When

	parentC := c.Parent(0)
	if parentC == nil {
		return false
	}

	parentTree, err := parentC.Tree()
	if err != nil {
		return false
	}

	cTree, err := c.Tree()
	if err != nil {
		return false
	}

	diffOpts, err := g2g.DefaultDiffOptions()
	if err != nil {
		return false
	}

	diff, err := gr.r.DiffTreeToTree(parentTree, cTree, &diffOpts)
	if err != nil {
		return false
	}

	stats, err := diff.Stats()
	if err != nil {
		return false
	}

	if gr.cfg.CommitDeltas {
		nDeltas, err := diff.NumDeltas()
		if err != nil {
			return false
		}

		for d := 0; d < nDeltas; d++ {
			var cdd model.DiffDelta

			if gr.cfg.CommitPatches {
				patch, err := diff.Patch(d)
				if err != nil || patch == nil {
					return false
				}
				p, err := patch.String()
				if err != nil {
					return false
				}
				cdd.Patch = &p
			}

			diffDelta, err := diff.GetDelta(d)
			if err != nil {
				return false
			}
			cdd.Status = deltaMap[diffDelta.Status]

			var isBin bool
			if (diffDelta.Flags & g2g.DiffFlagBinary) > 0 {
				isBin = true
			}
			cdd.Binary = &isBin

			// TODO compute similarity to add to cdd.Similarity

			cdd.OldFilePath = &diffDelta.OldFile.Path
			cdd.NewFilePath = &diffDelta.NewFile.Path

			commit.DiffDelta = append(commit.DiffDelta, cdd)
		}
	}

	commit.FileChangedCount = stats.FilesChanged()
	commit.InsertionsCount = stats.Insertions()
	commit.DeletionsCount = stats.Deletions()

	if ok := isCommitValid(commit); !ok {
		return false
	}
	gr.Commits = append(gr.Commits, commit)

	return true
}

// extractGitURL returns a git repository clone URL as a string, given the
// path to its location on disk.
func extractGitURL(path string) (*string, error) {
	f, err := os.Open(filepath.Join(path, ".git", "config"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("url ?= ?(.+)")
	match := re.FindStringSubmatch(string(bs))
	if len(match) != 2 {
		return nil, errors.New("invalid git config file")
	}

	if len(match[1]) == 0 {
		return nil, errors.New("cannot extract git clone url")
	}

	return &match[1], nil
}

// extractGitDefaultBranch returns the branch to which HEAD of a git
// repository is pointing at.
func extractGitDefaultBranch(path string) (*string, error) {
	f, err := os.Open(filepath.Join(path, ".git", "HEAD"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("ref: refs\\/heads\\/[a-zA-Z0-9].+")
	match := re.FindString(string(bs))
	if len(match) == 0 {
		return nil, errors.New("no branch (detached HEAD state)")
	}

	index := strings.LastIndex(match, "/")
	branch := string(match[index+1:])
	if len(branch) == 0 {
		// we shall not be able to reach here but just in case
		return nil, errors.New("no branch (detached HEAD state)")
	}

	return &branch, nil
}

// untarGitFolder extracts the root's .git directory contained in a tar archive
// of a git repository into destPath.
func untarGitFolder(destPath, archivePath string) error {
	var err error

	archiveFile, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	// make sure to create dest path
	if err = os.MkdirAll(destPath, os.ModePerm); err != nil {
		return err
	}

	tr := tar.NewReader(archiveFile)

	// make sure we keep the trailing /
	// FIXME not compatible with Windows
	basePath := filepath.Base(strings.TrimSuffix(archivePath, ".tar"))
	dotGitDirPath := filepath.Join(basePath, ".git") + "/"
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// we only want to extract the .git/ subtree and skip the rest
		if strings.HasPrefix(hdr.Name, dotGitDirPath) {
			hdr.Name = strings.TrimPrefix(hdr.Name, basePath)
			mode := hdr.FileInfo().Mode()
			switch {
			case mode&os.ModeDir != 0:
				if err := os.Mkdir(filepath.Join(destPath, hdr.Name), mode); err != nil {
					return err
				}
			case mode&os.ModeSymlink != 0:
				os.Symlink(hdr.Linkname, filepath.Join(destPath, hdr.Name))
			default: // consider it a regular file
				createFile := func() error {
					f, err := os.Create(filepath.Join(destPath, hdr.Name))
					if err != nil {
						return err
					}
					defer f.Close()

					if _, err := io.Copy(f, tr); err != nil {
						return err
					}
					return nil
				}

				if err = createFile(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
