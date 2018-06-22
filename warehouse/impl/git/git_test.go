package git

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/warpfork/go-errcat"
	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	riofs "go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/testutil"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// TODO: write tests across http/git protocol

var (
	RelPathBare  = riofs.MustRelPath("./bare")
	RelPathHash1 = riofs.MustRelPath("./hash1")
	RelPathHash2 = riofs.MustRelPath("./hash2")
	RelPathHash3 = riofs.MustRelPath("./hash3")
	RelPathHash4 = riofs.MustRelPath("./hash4")
	RelPathRepoA = riofs.MustRelPath("./repo-a")
	RelPathRepoB = riofs.MustRelPath("./repo-b")
	RelPathRepoC = riofs.MustRelPath("./repo-c")

	RelPathRepoA_DotGit = riofs.MustRelPath("./repo-a/.git")
	RelPathRepoB_DotGit = riofs.MustRelPath("./repo-b/.git")
	RelPathRepoC_DotGit = riofs.MustRelPath("./repo-c/.git")
)

/*
	tar -xzf warehouse/impl/git/fixtures/git_test_fixture.tar.gz \
	--to-stdout hash1 hash2 hash3 hash4 hashB hashC pwd
*/
const (
	hash1 = "f0e549c50372ac71af894309db05a63695e460a3"
	hash2 = "b64afb86af7150438beb62ac1b832e8a3ba831b9"
	hash3 = "86a7bda1e3a9b9ceceb7678aa77710db5c3f2b12"
	hash4 = "fed30e790df90ab04b767b727b6424f12f3077b9"
	hashB = "3501565627b7702806a9f8fd9ed892ee2c7f22c5"
	hashC = "3973bfaf4795e75924b03c25ae9e23b1374d3795"
	fixWd = "/tmp/tmp.TWgvjghRi7"
)

func TestSlugifyRemote(t *testing.T) {
	testItems := []struct {
		in  string
		out string
	}{
		{"http://codesource.google.com/project", "http%3A%2F%2Fcodesource.google.com%2Fproject"},
		{"https://go.polydawn.net/rio", "https%3A%2F%2Fgo.polydawn.net%2Frio"},
		{"https://go.polydawn.net/rio with spaces", "https%3A%2F%2Fgo.polydawn.net%2Frio+with+spaces"},
	}
	for _, remote := range testItems {
		t.Run(fmt.Sprintf("Slugify %s", remote.in), func(t *testing.T) {
			if result := SlugifyRemote(remote.in); result != remote.out {
				t.Errorf("expected %s but got %s", remote.out, result)
			}
		})
	}
}

func TestMustBeFullHash(t *testing.T) {
	incorrectSize := errcat.Errorf(rio.ErrUsage, "git commit hashes are 40 characters")
	incorrectEncoding := errcat.Errorf(rio.ErrUsage, "git commit hashes are hex strings")
	testItems := []struct {
		in  string
		out error
	}{
		{"", incorrectSize},
		{"1234567890123456789012345678901234567890", nil},
		{"0123456789abcdef0123456789abcdef01234567", nil},
		{"0123456789ABCDEF0123456789ABCDEF01234567", nil},
		{"0123456789ABCDEF0123456789ABCDEF012345678", incorrectSize},
		{"0123456789ABCDEF0123456789ABCDEF0123456", incorrectSize},
		{"0123456789ABCDEF0123456789ABCDEF0123456G", incorrectEncoding},
		{"0123456789ABCDEF0123456789ABCDEF0123456H", incorrectEncoding},
		{"0123456789ABCDEF0123456789ABCDEF0123456I", incorrectEncoding},
		{"0123456789ABCDEF0123456789ABCDEF0123456Z", incorrectEncoding},
		{"0123456789ABCDEF0123456789ABCDEF0123456z", incorrectEncoding},
		{"z0123456789ABCDEF0123456789ABCDEF0123456", incorrectEncoding},
	}
	for _, item := range testItems {
		t.Run(fmt.Sprintf("hash: %s", item.in), func(t *testing.T) {
			result := mustBeFullHash(item.in)
			resultString := fmt.Sprintf("%#v", result)
			expectedString := fmt.Sprintf("%#v", item.out)
			if resultString != expectedString {
				t.Errorf("expected %s but got %s", expectedString, resultString)
			}
		})
	}
}

func TestStringToHash(t *testing.T) {
	incorrectSize := errcat.Errorf(rio.ErrUsage, "git commit hashes are 40 characters")
	incorrectEncoding := errcat.Errorf(rio.ErrUsage, "git commit hashes are hex strings")
	testItems := []struct {
		in  string
		out plumbing.Hash
		err error
	}{
		{"", plumbing.Hash{}, incorrectSize},
		{"1234567890123456789012345678901234567890", plumbing.NewReferenceFromStrings("", "1234567890123456789012345678901234567890").Hash(), nil},
		{"0123456789abcdef0123456789abcdef01234567", plumbing.NewReferenceFromStrings("", "0123456789abcdef0123456789abcdef01234567").Hash(), nil},
		{"0123456789ABCDEF0123456789ABCDEF01234567", plumbing.NewReferenceFromStrings("", "0123456789ABCDEF0123456789ABCDEF01234567").Hash(), nil},
		{"0123456789ABCDEF0123456789ABCDEF012345678", plumbing.Hash{}, incorrectSize},
		{"0123456789ABCDEF0123456789ABCDEF0123456", plumbing.Hash{}, incorrectSize},
		{"0123456789ABCDEF0123456789ABCDEF0123456G", plumbing.Hash{}, incorrectEncoding},
		{"0123456789ABCDEF0123456789ABCDEF0123456H", plumbing.Hash{}, incorrectEncoding},
		{"0123456789ABCDEF0123456789ABCDEF0123456I", plumbing.Hash{}, incorrectEncoding},
		{"0123456789ABCDEF0123456789ABCDEF0123456Z", plumbing.Hash{}, incorrectEncoding},
		{"0123456789ABCDEF0123456789ABCDEF0123456z", plumbing.Hash{}, incorrectEncoding},
		{"z0123456789ABCDEF0123456789ABCDEF0123456", plumbing.Hash{}, incorrectEncoding},
	}
	for _, item := range testItems {
		t.Run(fmt.Sprintf("hash: %s", item.in), func(t *testing.T) {
			result, err := StringToHash(item.in)
			resultString := fmt.Sprintf("%#v", result)
			errString := fmt.Sprintf("%#v", err)
			expectedString := fmt.Sprintf("%#v", item.out)
			expectedErrString := fmt.Sprintf("%#v", item.err)
			if resultString != expectedString {
				t.Errorf("expected %s but got %s", expectedString, resultString)
			}
			if errString != expectedErrString {
				t.Errorf("expected %s but got %s", expectedErrString, errString)
			}
		})
	}
}

func TestHasFoldedSuffix(t *testing.T) {
	testItems := []struct {
		in     string
		suffix string
		out    bool
	}{
		{"", "", true},
		{"foobargrill", "", true},
		{"foobargrill", "grill", true},
		{"foobargrill", "GRILL", true},
		{"fooBaRgRiLl", "BARgrill", true},
		{"fooBaRgRiLl", "x", false},
		{"fooBaRgRiLl", "foo", false},
		{"fooBaRgRiLl", "gril", false},
	}
	for _, item := range testItems {
		t.Run(fmt.Sprintf("string %s has suffix %s", item.in, item.suffix), func(t *testing.T) {
			if result := HasFoldedSuffix(item.in, item.suffix); result != item.out {
				t.Errorf("expected %v but got %v", item.out, result)
			}
		})
	}
}

func TestHasFoldedPrefix(t *testing.T) {
	testItems := []struct {
		in     string
		prefix string
		out    bool
	}{
		{"", "", true},
		{"foobargrill", "", true},
		{"foobargrill", "foo", true},
		{"foobargrill", "FOO", true},
		{"FOOBARGRILL", "foo", true},
		{"FOOBARGRILL", "FOO", true},
		{"FoObargrill", "foo", true},
		{"FoObargrill", "FOO", true},
		{"FoObargrill", "FoO", true},
		{"foobargrill", "x", false},
		{"foobargrill", "grill", false},
		{"foobargrill", "foox", false},
		{"foobargrill", "bar", false},
	}
	for _, item := range testItems {
		t.Run(fmt.Sprintf("string %s has prefix %s", item.in, item.prefix), func(t *testing.T) {
			if result := HasFoldedPrefix(item.in, item.prefix); result != item.out {
				t.Errorf("expected %v but got %v", item.out, result)
			}
		})
	}
}

func untar(t *testing.T, tarball, target string) error {
	t.Logf("untar: %s to %s\n", tarball, target)
	tarballFile, err := os.Open(tarball)
	if err != nil {
		t.Fatal(err)
	}
	defer tarballFile.Close()
	tarfile, err := gzip.NewReader(tarballFile)
	if err != nil {
		t.Fatal(err)
	}
	tarReader := tar.NewReader(tarfile)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				t.Fatal(err)
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			t.Fatal(err)
		}
	}
	return nil
}

// like WithTmpDir but unpacks a tarball into it.
func WithTarballTmpDir(t *testing.T, fn func(absPath riofs.AbsolutePath)) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf("unable to get working dir")
	}
	fixtureTarball := filepath.Join(pwd, "fixtures", "git_test_fixture.tar.gz")
	testutil.WithTmpdir(func(absPath riofs.AbsolutePath) {
		err := untar(t, fixtureTarball, absPath.String())
		if err != nil {
			t.Fatalf("unable to unpack fixture: %s", err.Error())
		}
		fn(absPath)
	})
}

func LogLs(t *testing.T, path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		t.Log(file.Name())
	}
}

// Run basic tests the fixture works and is valid
func TestFixtures(t *testing.T) {
	WithTarballTmpDir(t, func(absPath riofs.AbsolutePath) {
		LogLs(t, absPath.String())
		testItems := []struct {
			hash string
			path riofs.RelPath
		}{
			{hash1, RelPathHash1},
			{hash2, RelPathHash2},
			{hash3, RelPathHash3},
			{hash4, RelPathHash4},
		}

		for _, item := range testItems {
			t.Run(item.path.String(), func(t *testing.T) {
				hashFromFile, err := ioutil.ReadFile(absPath.Join(item.path).String())
				if err != nil {
					t.Errorf("unexpected error: \"%s\"", err)
				}
				hashString := strings.TrimSpace(string(hashFromFile))
				if hashString != item.hash {
					t.Errorf("expected \"%s\" but got \"%s\"", item.hash, hashString)
				}
			})
		}
	})
}

func TestNewControllerNonExistingWarehouse(t *testing.T) {
	controller, err := NewController(nil, api.WarehouseAddr("bogus"))
	if err == nil {
		t.Errorf("expected an error")
	}
	errc := err.(errcat.Error)
	if errc.Category() != rio.ErrWarehouseUnavailable {
		t.Errorf("expected error category \"%s\" but got \"%s\"", rio.ErrWarehouseUnavailable, errc.Category())
	}
	if controller != nil {
		t.Errorf("expected controller to be nil")
	}
}

func TestNewControllerWithFixture(t *testing.T) {
	testItems := []riofs.RelPath{
		RelPathBare,
		RelPathRepoA,
		RelPathRepoB,
		RelPathRepoC,
		RelPathRepoA_DotGit,
		RelPathRepoB_DotGit,
		RelPathRepoC_DotGit,
	}
	WithTarballTmpDir(t, func(absPath riofs.AbsolutePath) {
		for _, item := range testItems {
			t.Run(item.String(), func(t *testing.T) {
				wareAddr := api.WarehouseAddr(absPath.Join(item).String())
				t.Log(wareAddr)
				controller, err := NewController(nil, wareAddr)
				if err != nil {
					t.Errorf("failed to create controller: \"%s\"", err.Error())
					errc := err.(errcat.Error)
					t.Errorf("%s", errc.Category())
					t.Errorf("%s", errc.Details())
				}
				if controller == nil {
					t.Errorf("expected controller to not be nil")
				}
			})
		}
	})
}

func TestGetCommit(t *testing.T) {
	WithTarballTmpDir(t, func(absPath riofs.AbsolutePath) {
		wareAddr := api.WarehouseAddr(absPath.Join(RelPathBare).String())
		t.Log(wareAddr)
		controller := mustNewController(t, nil, wareAddr)
		// test hapy path conditions
		testItems := []struct {
			hash    string
			message string
		}{
			{hash1, "initial commit"},
			{hash2, "commit file1"},
			{hash3, "commit file2"},
			{hash4, "add submodule"},
		}
		for _, item := range testItems {
			t.Run(item.hash, func(t *testing.T) {
				commit, err := controller.GetCommit(item.hash)
				if err != nil {
					t.Errorf("unexpected error: \"%s\"", err)
				}
				if commit == nil {
					t.Fatal("commit should not be nil")
				}
				if string(commit.Hash.String()) != item.hash {
					t.Errorf("expected hash \"%s\" but got \"%s\"", item.hash, commit.Hash)
				}
				if strings.TrimSpace(commit.Message) != item.message {
					t.Errorf("expected message \"%s\" but got \"%s\"", item.message, commit.Message)
				}
			})
		}
		// test error conditions
		t.Run("test short commit hash", func(t *testing.T) {
			commit, err := controller.GetCommit("too-short")
			if commit != nil {
				t.Errorf("expected nil commit object")
			}
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			expectedMessage := "git commit hashes are 40 characters"
			if err.Error() != expectedMessage {
				t.Errorf("expected \"%s\" but got \"%s\"", expectedMessage, err)
			}
			errc := err.(errcat.Error)
			expectedCategory := rio.ErrUsage
			if errc.Category() != expectedCategory {
				t.Errorf("expected \"%s\" but got \"%s\"", expectedCategory, errc.Category())
			}
		})
		t.Run("test missing hash", func(t *testing.T) {
			commit, err := controller.GetCommit("ffffffffffffffffffffffffffffffffffffffff")
			if commit != nil {
				t.Errorf("expected nil commit object")
			}
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			expectedMessage := "commit not found"
			if err.Error() != expectedMessage {
				t.Errorf("expected \"%s\" but got \"%s\"", expectedMessage, err)
			}
			errc := err.(errcat.Error)
			expectedCategory := rio.ErrWareNotFound
			if errc.Category() != expectedCategory {
				t.Errorf("expected \"%s\" but got \"%s\"", expectedCategory, errc.Category())
			}
		})
		t.Run("test invalid hash", func(t *testing.T) {
			commit, err := controller.GetCommit("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")
			if commit != nil {
				t.Errorf("expected nil commit object")
			}
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			expectedMessage := "git commit hashes are hex strings"
			if err.Error() != expectedMessage {
				t.Errorf("expected \"%s\" but got \"%s\"", expectedMessage, err)
			}
			errc := err.(errcat.Error)
			expectedCategory := rio.ErrUsage
			if errc.Category() != expectedCategory {
				t.Errorf("expected \"%s\" but got \"%s\"", expectedCategory, errc.Category())
			}
		})
		t.Run("test corrupt", func(t *testing.T) {
			t.Skip("FIXME: I'm not sure how to do this properly.")
			// needs to somehow create a corrupt object
			enc := controller.store.NewEncodedObject()
			blob := &object.Blob{}
			err := blob.Encode(enc)
			if err != nil {
				t.Errorf("unexpected error: \"%s\"", err)
			}
			hash, err := controller.store.SetEncodedObject(enc)
			if err != nil {
				t.Errorf("unexpected error: \"%s\"", err)
			}
			commit, err := controller.GetCommit(hash.String())
			if commit != nil {
				t.Errorf("expected nil commit object")
			}
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			expectedMessage := "git commit hashes are hex strings"
			if err.Error() != expectedMessage {
				t.Errorf("expected \"%s\" but got \"%s\"", expectedMessage, err)
			}
			errc := err.(errcat.Error)
			expectedCategory := rio.ErrWareCorrupt
			if errc.Category() != expectedCategory {
				t.Errorf("expected \"%s\" but got \"%s\"", expectedCategory, errc.Category())
			}
		})
	})
}

func TestContains(t *testing.T) {
	WithTarballTmpDir(t, func(absPath riofs.AbsolutePath) {
		wareAddr := api.WarehouseAddr(absPath.Join(RelPathBare).String())
		controller := mustNewController(t, nil, wareAddr)
		t.Run("invalid commit", func(t *testing.T) {
			if controller.Contains("bogus") {
				t.Errorf("expected %v", false)
			}
		})
		t.Run("bad commit", func(t *testing.T) {
			if controller.Contains(hashB) {
				t.Errorf("expected %v", false)
			}
		})
		t.Run("ok commit", func(t *testing.T) {
			if !controller.Contains(hash2) {
				t.Errorf("expected %v", true)
			}
		})
	})
}

func TestGetTree(t *testing.T) {
	WithTarballTmpDir(t, func(absPath riofs.AbsolutePath) {
		wareAddr := api.WarehouseAddr(absPath.Join(RelPathBare).String())
		controller := mustNewController(t, nil, wareAddr)
		t.Run("invalid commit", func(t *testing.T) {
			tree, err := controller.GetTree("bogus")
			if err == nil {
				t.Errorf("expected error to not be nil")
			}
			if tree != nil {
				t.Errorf("expected tree object to be nil")
			}
		})
		t.Run("bad commit", func(t *testing.T) {
			tree, err := controller.GetTree(hashB)
			if err == nil {
				t.Errorf("expected error to not be nil")
			}
			if tree != nil {
				t.Errorf("expected tree object to be nil")
			}
		})
		t.Run("ok commit", func(t *testing.T) {
			tree, err := controller.GetTree(hash2)
			if err != nil {
				t.Errorf("expected nil but got %s", err)
			}
			if tree == nil {
				t.Errorf("expected tree object to not be nil")
			}
			// TODO: Maybe write more tests here
		})
	})
}

// TODO
func TestSetCacheStorage(t *testing.T) {
	t.Log("TODO")
	t.SkipNow()
}

func TestClone(t *testing.T) {
	WithTarballTmpDir(t, func(absPath riofs.AbsolutePath) {
		wareAddr := api.WarehouseAddr(absPath.Join(RelPathBare).String())
		t.Run("clone disabled", func(t *testing.T) {
			controller := mustNewController(t, nil, wareAddr)
			if controller.repo != nil {
				t.Error("expected nil repository")
			}
			if controller.allowClone {
				t.Error("expected controller to have clone disabled")
			}
			ctx := context.Background()
			err := controller.Clone(ctx)
			if err != nil {
				t.Error(err)
			}
			if controller.repo == nil {
				t.Errorf("expected non-nil repository")
			}
		})
		t.Run("clone enabled", func(t *testing.T) {
			controller := mustNewController(t, nil, wareAddr)
			if controller.repo != nil {
				t.Error("expected nil repository")
			}
			controller.allowClone = true
			ctx := context.Background()
			err := controller.Clone(ctx)
			if err != nil {
				t.Error(err)
			}
			if controller.repo == nil {
				t.Errorf("expected non-nil repository")
			}
		})
		t.Run("clone cancelled default storage", func(t *testing.T) {
			// The test is against a locally stored filesystem.
			// By default it will simply open the repository in place
			controller := mustNewController(t, nil, wareAddr)
			if controller.repo != nil {
				t.Error("expected nil repository")
			}
			controller.allowClone = true
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			cancel()
			err := controller.Clone(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if controller.repo == nil {
				t.Errorf("expected non-nil repository")
			}
		})
		t.Run("clone cancelled with different storage", func(t *testing.T) {
			// replacing the storage will cause the system
			// to copy the repo into the new storage.
			controller := mustNewController(t, nil, wareAddr)
			controller.allowClone = true
			if controller.repo != nil {
				t.Error("expected nil repository")
			}
			controller.store = memory.NewStorage()
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			cancel()
			err := controller.Clone(ctx)
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			errc := err.(errcat.Error)
			if errc.Category() != rio.ErrCancelled {
				t.Errorf("expected error category %s but got %s", rio.ErrCancelled, errc.Category())
			}
			if controller.repo != nil {
				t.Errorf("expected nil repository")
			}
		})
	})
}

func TestUpdate(t *testing.T) {
	WithTarballTmpDir(t, func(absPath riofs.AbsolutePath) {
		wareAddr := api.WarehouseAddr(absPath.Join(RelPathBare).String())
		testItems := []struct {
			fetchAllowed bool // updates can be disabled for some conditions
			isNewClone   bool // a fresh clone is fully up-to-date
			doesUpdate   bool //
		}{
			{false, false, false},
			{false, true, false},
			{true, false, true},
			{true, true, false},
		}
		for _, item := range testItems {
			testName := fmt.Sprintf("update after open{fetch=%v,fresh=%v}", item.fetchAllowed, item.isNewClone)
			t.Run(testName, func(t *testing.T) {
				store := memory.NewStorage()
				controller := mustNewController(t, nil, wareAddr)
				if controller.repo != nil {
					t.Error("expected nil repository")
				}
				controller.store = store
				controller.allowClone = true
				ctx := context.Background()
				err := controller.Clone(ctx)
				if err != nil {
					t.Fatal(err)
				}
				// move master back to hash3
				ref := plumbing.NewReferenceFromStrings("refs/heads/master", hash3)
				store.SetReference(ref)
				// remove hash4 from storage
				plumbHash4 := plumbing.NewHash(hash4)
				delete(store.ObjectStorage.Commits, plumbHash4)
				delete(store.ObjectStorage.Objects, plumbHash4)
				obj, err := store.EncodedObject(plumbing.AnyObject, plumbHash4)
				if err != plumbing.ErrObjectNotFound {
					t.Fatal(err)
				}
				if obj != nil {
					t.Fatalf("object hash %s should be missing", hash4)
				}
				controller.allowFetch = item.fetchAllowed
				controller.newClone = item.isNewClone
				err = controller.Update(ctx)
				if err == transport.ErrEmptyUploadPackRequest {
					t.Skip("FIXME: Can't get fetch to actually work.")
				}
				if err != nil {
					t.Fatal(err)
				}
				if controller.repo == nil {
					t.Errorf("expected non-nil repository")
				}
				checkRef, err := controller.store.Reference("refs/heads/master")
				if err != nil {
					t.Fatal(err)
				}
				updated := checkRef.Hash() == plumbHash4
				if item.doesUpdate != updated {
					t.Logf("reference: %v", checkRef)
					t.Errorf("expected hashes (%s == %s) to be %v", checkRef.Hash(), plumbHash4, item.doesUpdate)
				}
			})
		}
	})
}

func TestOpen(t *testing.T) {
	WithTarballTmpDir(t, func(absPath riofs.AbsolutePath) {
		wareAddr := api.WarehouseAddr(absPath.Join(RelPathBare).String())
		t.Run("clone non-existent", func(t *testing.T) {
			controller := mustNewController(t, nil, wareAddr)
			// Basically this means that the repo became unavailable _after_
			// creation of the controller but before we tried to pull the repo.
			controller.sanitizedAddr = absPath.String()
			allowClone := true
			ctx := context.Background()
			store := memory.NewStorage()
			repo, err := controller.open(ctx, store, allowClone)
			if err == nil {
				t.Errorf("expected an error but got nil")
			}
			errc := err.(errcat.Error)
			if errc.Category() != rio.ErrWareCorrupt {
				t.Errorf("expected error category %s but got %s", rio.ErrWareCorrupt, errc.Category())
			}
			if errc.Message() != "unable to clone repository: repository not found" {
				t.Errorf("unexpected error message: %s", errc.Message())
			}
			if repo != nil {
				t.Errorf("expected repository to be nil")
			}
		})
		t.Run("clone enabled memory", func(t *testing.T) {
			controller := mustNewController(t, nil, wareAddr)
			allowClone := true
			ctx := context.Background()
			store := memory.NewStorage()
			repo, err := controller.open(ctx, store, allowClone)
			if err != nil {
				t.Errorf("expected an nil but got %s", err)
			}
			if repo == nil {
				t.Errorf("expected repository to not be nil")
			}
		})
		t.Run("clone disabled", func(t *testing.T) {
			controller := mustNewController(t, nil, wareAddr)
			allowClone := false
			ctx := context.Background()
			store := memory.NewStorage()
			repo, err := controller.open(ctx, store, allowClone)
			if err == nil {
				t.Errorf("expected an error but got nil")
			}
			errc := err.(errcat.Error)
			if errc.Category() != rio.ErrLocalCacheProblem {
				t.Errorf("expected error category %s but got %s", rio.ErrLocalCacheProblem, errc.Category())
			}
			if errc.Message() != "unable to open cache repository: repository does not exist" {
				t.Errorf("unexpected error message: %s", errc.Message())
			}
			if repo != nil {
				t.Errorf("expected repository to be nil")
			}
		})
		t.Run("clone canceled memory", func(t *testing.T) {
			controller := mustNewController(t, nil, wareAddr)
			allowClone := true
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			cancel()
			store := memory.NewStorage()
			repo, err := controller.open(ctx, store, allowClone)
			errc := err.(errcat.Error)
			if errc.Category() != rio.ErrCancelled {
				t.Errorf("expected error category %s but got %s", rio.ErrCancelled, errc.Category())
			}
			msg := "cancelled: sending upload-req message: encoding first want line: context canceled"
			if errc.Message() != msg {
				t.Errorf("unexpected error message: %s", errc.Message())
			}
			if repo != nil {
				t.Errorf("expected repository to be nil")
			}
		})
		// TODO: set up a server and do non-file tests
	})
}

func TestLsRemote(t *testing.T) {
	WithTarballTmpDir(t, func(absPath riofs.AbsolutePath) {
		wareAddr := api.WarehouseAddr(absPath.Join(RelPathBare).String())
		controller := mustNewController(t, nil, wareAddr)
		result, err := controller.lsRemote()
		if err != nil {
			t.Fatal(err)
		}
		iter, err := result.IterReferences()
		if err != nil {
			t.Fatal(err)
		}
		counter := 0
		expectedRefs := []*plumbing.Reference{
			plumbing.NewReferenceFromStrings("HEAD", "ref: refs/heads/master"), // symbolic reference target must be prefixed by "ref: "
			plumbing.NewReferenceFromStrings("refs/heads/master", hash4),
		}
		iter.ForEach(func(ref *plumbing.Reference) error {
			found := false
			for _, expRef := range expectedRefs {
				if reflect.DeepEqual(ref, expRef) {
					found = true
					break
				}
			}
			counter++
			if found != true {
				t.Errorf("unexpected reference: %+v", ref)
			}
			return nil
		})
		if counter != len(expectedRefs) {
			t.Errorf("expected %d references but got %d", len(expectedRefs), counter)
		}
	})
}

func TestSubmodules(t *testing.T) {
	WithTarballTmpDir(t, func(absPath riofs.AbsolutePath) {
		wareAddr := api.WarehouseAddr(absPath.Join(RelPathBare).String())
		t.Log(wareAddr)
		controller := mustNewController(t, nil, wareAddr)
		for _, testItem := range []struct {
			hash string
		}{
			{hash1},
			{hash2},
			{hash3},
		} {
			t.Run(testItem.hash, func(t *testing.T) {
				// we expect that these hashes do _not_ have submodules
				result, err := controller.Submodules(testItem.hash)
				if err != nil {
					t.Fatal(err)
				}
				if len(result) != 0 {
					t.Fatal("expected empty submodule list")
				}
			})
		}

		result, err := controller.Submodules(hash4)
		if err != nil {
			t.Fatal(err)
		}
		expectedSubmodulesLength := 1
		if len(result) != expectedSubmodulesLength {
			t.Fatalf("expected %d submodules but got %d", expectedSubmodulesLength, len(result))
		}
		submoduleAbsPath := riofs.MustAbsolutePath(fixWd)
		for _, value := range result {
			assertSubmoduleEqual(t, value, Submodule{
				Path:   "repo-b",
				Name:   "repo-b",
				URL:    submoduleAbsPath.Join(RelPathRepoB).String(),
				Branch: "",
				Hash:   hashB,
			})
		}
	})
}

// Example of using the reader
func TestReader(t *testing.T) {
	WithTarballTmpDir(t, func(absPath riofs.AbsolutePath) {
		// The bare repository is a mirror of repo A
		// we will compare the files extracted from bare to the files in repo A
		repoPath := absPath.Join(RelPathRepoA)
		warePath := absPath.Join(RelPathBare)
		wareAddr := api.WarehouseAddr(warePath.String())
		controller := mustNewController(t, nil, wareAddr)
		t.Run("test reader bad hash", func(t *testing.T) {
			wareReader, err := controller.NewReader("bogus")
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			errc := err.(errcat.Error)
			if errc.Category() != rio.ErrUsage {
				t.Fatalf("expected category %s but got %s", rio.ErrUsage, errc.Category())
			}
			if wareReader != nil {
				t.Fatal("expected nil reader")
			}
		})
		t.Run("create reader extract files", func(t *testing.T) {
			wareReader, err := controller.NewReader(hash4)
			if err != nil {
				t.Fatal(err)
			}
			defer wareReader.Close()
			// reading without calling next once is not allowed
			func() {
				defer func() {
					fmt.Print("recover")
					recover()
				}()
				buf := []byte{}
				wareReader.Read(buf)
				t.Fatal("expected panic")
			}()
			for {
				header, err := wareReader.Next()
				if err == io.EOF {
					break
				} else if err != nil {
					t.Fatal(err)
					break
				}
				t.Logf("%+v", header)
				buf := make([]byte, header.Size)
				n, err := wareReader.Read(buf)
				if err != nil {
					t.Fatal(err)
				}
				if n == 0 {
					t.Fatal(err)
				}
				relpath := riofs.MustRelPath(header.Name)
				filepath := repoPath.Join(relpath).String()
				expected, err := ioutil.ReadFile(filepath)
				if err != nil {
					t.Fatalf("error reading file %s: %s", filepath, err)
				}
				if !bytes.Equal(expected, buf) {
					t.Errorf("expected file contents\n%s\nbut got\n%s", expected, buf)
				}
			}
		})
	})
}

func assertTreeEntryEqual(t *testing.T, a, b *object.TreeEntry) {
	if a.Name != b.Name {
		t.Errorf("Expected name %v but got %v", b.Name, a.Name)
	}
	if a.Mode != b.Mode {
		t.Errorf("Expected filemode %v but got %v", b.Mode, a.Mode)
	}
	if a.Hash != b.Hash {
		t.Errorf("Expected hash %v but got %v", b.Hash, a.Hash)
	}
}

func assertSubmoduleEqual(t *testing.T, a, b Submodule) {
	t.Helper()
	if a.Name != b.Name {
		t.Errorf("Name %s not equal to Name %s", a.Name, b.Name)
	}
	if a.Path != b.Path {
		t.Errorf("Path %s not equal to Path %s", a.Path, b.Path)
	}
	if a.URL != b.URL {
		t.Errorf("URL %s not equal to URL %s", a.URL, b.URL)
	}
	if a.Branch != b.Branch {
		t.Errorf("Branch %s not equal to Branch %s", a.Branch, b.Branch)
	}
	if a.Hash != b.Hash {
		t.Errorf("Hash %s not equal to Hash %s", a.Hash, b.Hash)
	}
}

func mustNewController(t *testing.T, workingDir riofs.FS, wareAddr api.WarehouseAddr) *Controller {
	t.Helper()
	controller, err := NewController(workingDir, wareAddr)
	if err != nil {
		t.Errorf("failed to create controller: \"%s\"", err.Error())
		errc := err.(errcat.Error)
		t.Errorf("%s", errc.Category())
		t.Errorf("%s", errc.Details())
		t.FailNow()
	}
	return controller
}
