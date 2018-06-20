package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

func main() {
	m, err := MD5All(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x  %s\n", m[path], path)
	}

}

func MD5All(root string) (map[string][]byte, error) {
	c := make(chan result)
	var wg sync.WaitGroup
	numDigesters := 10
	wg.Add(numDigesters)

	done := make(chan struct{})
	paths, errc := walkFiles(done, root)

	for i := 0; i < numDigesters; i++ {
		go func() {
			digester(done, paths, c)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	m := make(map[string][]byte)
	for r := range c {
		if r.err != nil {
			log.Println(r.err)
		}
		m[r.path] = r.sum
	}

	// Check whether the Walk failed.
	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
}

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(paths)

		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			case <-done:
				return errors.New("walk cancelled")
			}
			return nil
		})
	}()
	return paths, errc
}

type result struct {
	path string
	sum  []byte
	err  error
}

func digester(done <-chan struct{}, paths <-chan string, c chan<- result) {
	sumFile := func(file *os.File) ([]byte, error) {
		defer file.Close()
		h := md5.New()
		if _, err := io.Copy(h, file); err != nil {
			return nil, err
		}
		return h.Sum(nil), nil
	}

	for path := range paths {
		r := result{path: path}

		file, err := os.Open(path)
		var sum []byte
		if err != nil {
			r.err = err
		} else {
			sum, err = sumFile(file)
			if err != nil {
				r.err = err
			} else {
				r.sum = sum
			}
		}
		select {
		case c <- r:
		case <-done:
			return
		}
	}
}
