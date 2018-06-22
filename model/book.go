package model

import "fmt"

//Book - The DTO used to access books
type Book struct {
	Entity
	Title       string `json:"title"`
	NoPages     int    `json:"noPages"`
	ReleaseDate string `json:"releaseDate"`
	Author      Author `json:"author" gorm:"foreignkey:AuthorId"`
	AuthorId    int    `json:"-"`
}

func (book Book) String() string {
	return fmt.Sprintf("Book{UUID='%s', Title='%s', NoPages=%d, ReleaseDate='%s',Author=%s}",
		book.UUID, book.Title, book.NoPages, book.ReleaseDate, book.Author)
}

//BooksList - A list of Book objects
type BooksList []Book

func (b *BooksList) Get(bookUUID string) (Book, error) {
	err := fmt.Errorf("could not find book by uuid %s", bookUUID)
	for _, book := range *b {
		if book.UUID == bookUUID {
			return book, nil
		}
	}
	return Book{}, err
}

func (b *BooksList) Delete(bookUUID string) error {
	var err = fmt.Errorf("could not find book by uuid %s", bookUUID)
	var updatedBooks BooksList
	for _, book := range *b {
		if book.UUID == bookUUID {
			err = nil
		} else {
			updatedBooks = append(updatedBooks, book)
		}
	}
	if err == nil {
		*b = updatedBooks
	}
	return err
}

func (b *BooksList) Update(updatedBook Book) (Book, error) {
	var err = fmt.Errorf("could not find book by uuid %s", updatedBook.UUID)
	var newBooks BooksList
	for _, book := range *b {
		if book.UUID == updatedBook.UUID {
			newBooks = append(newBooks, updatedBook)
			err = nil
		} else {
			newBooks = append(newBooks, book)
		}
	}
	if err == nil {
		*b = newBooks
	}
	return updatedBook, err
}

//Adds the given book into the list
func (b *BooksList) Add(book Book) {
	*b = append(*b, book)
}
