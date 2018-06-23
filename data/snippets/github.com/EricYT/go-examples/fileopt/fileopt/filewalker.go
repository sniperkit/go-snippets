package parallel

import (
	"os"
	"path/filepath"

	"github.com/EricYT/go-examples/fileopt/worker"

	tomb "gopkg.in/tomb.v1"
)

type FilterFun func(path string, info os.FileInfo, err error) (worker.Jobber, error)

type FileGenerator struct {
	tomb *tomb.Tomb

	root      string
	filterFun FilterFun

	jobCh chan worker.Jobber
}

func NewFileGenerator(root string, fun FilterFun) *FileGenerator {
	f := &FileGenerator{
		tomb:      new(tomb.Tomb),
		root:      root,
		filterFun: fun,
		jobCh:     make(chan worker.Jobber, 10),
	}

	go func() {
		defer f.tomb.Done()
		f.tomb.Kill(f.runLoop())
	}()

	return f
}

func (f *FileGenerator) Wait() error {
	return f.tomb.Wait()
}

func (f *FileGenerator) Kill(reason error) {
	f.tomb.Kill(reason)
}

func (f *FileGenerator) Generate() <-chan worker.Jobber {
	return f.jobCh
}

func (f *FileGenerator) runLoop() error {
	//FIXME: For very large directories Walk can be inefficient.
	return filepath.Walk(f.root, f.Walk)
}

func (f *FileGenerator) Walk(path string, info os.FileInfo, err error) error {
	// skip the root path
	if filepath.Clean(f.root) == filepath.Clean(path) {
		return nil
	}

	select {
	case <-f.tomb.Dying():
		return nil
	default:
	}

	// filter a specifial file
	j, e := f.filterFun(path, info, err)
	if e != nil {
		return e
	}

	if j != nil {
		f.jobCh <- j
	}

	return nil
}
