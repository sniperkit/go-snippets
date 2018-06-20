package redis

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/kdar/cachei"
	//"github.com/vmihailenco/msgpack"

	"reflect"
)

func init() {
	cachei.Register("redis", &Source{})
}

type Source struct {
	conn   redis.Conn
	config cachei.DataSource
}

func (s *Source) Setup(config cachei.DataSource) error {
	s.config = config
	return s.Open()
}

func (s *Source) Open() (err error) {
	if s.config.Port == 0 {
		s.config.Port = 6379
	}

	s.conn, err = redis.Dial("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	return
}

func (s *Source) Close() error {
	return s.conn.Close()
}

func (s *Source) GetSetFn(key string, expires int, f cachei.CacheFunc) (interface{}, error, error) {
	var i interface{}
	ferr, cerr := s.OutSetFn(key, expires, &i, f)
	return i, ferr, cerr
}

func (s *Source) OutSetFn(key string, expires int, out interface{}, f cachei.CacheFunc) (error, error) {
	outval := reflect.ValueOf(out)
	if outval.Kind() != reflect.Ptr {
		return nil, errors.New("ret must be a pointer")
	}

	shouldSet := false
	v, err := redis.Bytes(s.conn.Do("GET", key))
	if err != redis.ErrNil || err == nil {
		err = s.config.Coder.Unmarshal(v, &out)
		if err != nil {
			fmt.Println(err)
			shouldSet = true
		} else {
			return nil, nil
		}
	}

	if shouldSet || err == redis.ErrNil {
		if f != nil {
			fret, rerr := f()
			fval := reflect.ValueOf(fret)
			outval.Elem().Set(fval)

			if fret != nil {
				err = s.Set(key, fret, expires)
				if err != nil {
					return rerr, err
				}
			}

			return rerr, rerr
		}
	}

	return nil, err
}

func (s *Source) Get(key string) (interface{}, error) {
	var i interface{}
	err := s.Out(key, &i)
	return i, err
}

func (s *Source) Out(key string, out interface{}) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return errors.New("out must be a pointer")
	}

	v, err := redis.Bytes(s.conn.Do("GET", key))
	if err == nil || err != redis.ErrNil {
		err = s.config.Coder.Unmarshal(v, &out)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("could not find cache")
}

func (s *Source) Set(key string, value interface{}, expires int) error {
	b, err := s.config.Coder.Marshal(value)
	if err != nil {
		return err
	}

	s.conn.Send("MULTI")
	s.conn.Send("SET", key, string(b))
	s.conn.Send("EXPIRE", key, expires)
	_, err = s.conn.Do("EXEC")

	return err
}

func (s *Source) Delete(key string) error {
	_, err := s.conn.Do("DEL", key)
	return err
}
