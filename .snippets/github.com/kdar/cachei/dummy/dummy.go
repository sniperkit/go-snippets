package dummy

import (
	"errors"
	"github.com/kdar/cachei"
	"reflect"
)

func init() {
	cachei.Register("dummy", &Source{})
}

type Source struct{}

func (s *Source) Setup(config cachei.DataSource) error {
	return nil
}

func (s *Source) Open() error {
	return nil
}

func (s *Source) Close() error {
	return nil
}

func (s *Source) OutSetFn(key string, expires int, out interface{}, f cachei.CacheFunc) (error, error) {
	outval := reflect.ValueOf(out)
	if outval.Kind() != reflect.Ptr {
		return nil, errors.New("out must be a pointer")
	}

	fret, ferr, cerr := s.GetSetFn(key, expires, f)
	if fret != nil {
		fval := reflect.ValueOf(fret)
		outval.Elem().Set(fval)
	}

	return ferr, cerr
}

func (s *Source) GetSetFn(key string, expires int, f cachei.CacheFunc) (interface{}, error, error) {
	fret, ferr := f()
	return fret, ferr, nil
}

func (s *Source) Out(key string, out interface{}) error {
	return nil
}

func (s *Source) Get(key string) (interface{}, error) {
	return nil, errors.New("could not find cache")
}

func (s *Source) Set(key string, value interface{}, expires int) error {
	return nil
}

func (s *Source) Delete(key string) error {
	return nil
}
