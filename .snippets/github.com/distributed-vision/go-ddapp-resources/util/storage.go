package util

import (
	"io/ioutil"
	"os"
	"path"
)

type Storage interface {
	Open(readOnOpen bool) ([]byte, error)
	Read(length, position uint) ([]byte, error)
	Write(buffer []byte, position uint) error
	Sync() error
	Close() error
	Delete() error
	IsOpen() bool
}

type FileStorage struct {
	path string
	flag int
	mode os.FileMode
	file *os.File
}

func NewFileStorage(path string, mode os.FileMode) *FileStorage {
	return &FileStorage{path: path, mode: mode}
}

func (this *FileStorage) Open(readOnOpen bool) ([]byte, error) {
	err := makeDirectoryPath(dirname(this.path), os.ModePerm)

	if err != nil {
		return nil, err
	}

	this.file, err = os.OpenFile(this.path, os.O_RDWR|os.O_CREATE, this.mode)

	if err != nil {
		return nil, err
	}

	if readOnOpen {
		return ioutil.ReadFile(this.path)
	}

	return nil, nil
}

func (this *FileStorage) IsOpen() bool {
	return this.file != nil
}

func (this *FileStorage) Read(length, position uint) ([]byte, error) {
	buffer := make([]byte, length)
	_, err := this.file.ReadAt(buffer, int64(position))

	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (this *FileStorage) Write(buffer []byte, position uint) error {
	_, err := this.file.WriteAt(buffer, int64(position))
	return err
}

func (this *FileStorage) Sync() error {
	return this.file.Sync()
}

func (this *FileStorage) Close() error {
	err := this.file.Close()
	this.file = nil
	return err
}

func (this *FileStorage) Delete() error {
	this.Close()
	return os.Remove(this.path)
}

func makeDirectoryPath(path string, mode os.FileMode) error {
	return os.MkdirAll(path, mode)
}

func dirname(p string) string {
	return path.Dir(p)
}
