package persistence

import (
	"github.com/jinzhu/gorm"
	"goworkshop/model"
)

type GormDataStore struct {
	DBInstance *gorm.DB
}

type DataStore interface {
	//books operations
	CreateBook(*model.Book) error
	GetBook(string) (model.Book, error)
	GetBooks() ([]model.Book, error)
	UpdateBook(string, *model.Book) error
	DeleteBook(string) error

	//author operations
	CreateAuthor(author *model.Author) error
	GetAuthor(string) (model.Author, error)
	GetAuthors() ([]model.Author, error)
	UpdateAuthor(string, *model.Author) error
	DeleteAuthor(string) error
}

func isEmptyUUID(uuid string) bool {
	if len(uuid) == 0 {
		return true
	} else {
		return false
	}
}
