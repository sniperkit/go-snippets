package model

type Book struct {
	Model

	Name string `gorm:"unique_index;size:50;not null;default:''"`
}

type BookItem struct {
	Model

	Name string `gorm:"index;not null;default:50;default:''"`

	BookId uint `gorm:"index;not null;default:0"`
	Book   Book
}

type BookPage struct {
	Model

	BookItem   BookItem
	BookItemId uint `gorm:"unique_index;not null;default:0"`

	Txt string
}
