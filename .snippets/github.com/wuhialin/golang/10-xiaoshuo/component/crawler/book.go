package crawler

import (
	"../../model"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/jaytaylor/html2text"
	"log"
	"strings"
)

func bookItem(e *colly.HTMLElement) {
	var err error
	pages := e.DOM.Find("#chapterlist > li > a")
	if e.DOM.Find("h1.page-title").Length() == 0 || pages.Length() == 0 {
		return
	}
	title := strings.TrimSpace(e.ChildText("h1.page-title"))
	if title == "" {
		return
	}
	book := model.Book{Name: title}
	if err = db.FirstOrCreate(&book, "name = ?", title).Error; err != nil {
		log.Println(err)
		return
	}
	pages.Each(func(_ int, selection *goquery.Selection) {
		name := strings.TrimSpace(selection.Text())
		if name == "" {
			return
		}
		bookItem := model.BookItem{
			BookId: book.ID,
			Name:   name,
		}
		db.FirstOrCreate(&bookItem, "book_id = ? AND name = ?", book.ID, name)
	})
}

func bookPage(e *colly.HTMLElement) {
	txt := e.DOM.Find("div#zcontent")
	if txt.Length() == 0 {
		return
	}
	bookName := e.DOM.Find("div.mod.mod-back.breadcrumb > div.bd > a:nth-of-type(3)")
	if bookName.Length() == 0 {
		return
	}
	bookItemName := e.DOM.Find("h1.page-title")
	if bookItemName.Length() == 0 {
		return
	}
	name := strings.TrimSpace(bookName.Text())
	if name == "" {
		return
	}
	book := model.Book{Name: name}
	db.Find(&book, "name = ?", name)
	if book.ID == 0 {
		return
	}
	name = strings.TrimSpace(bookItemName.Text())
	if name == "" {
		return
	}
	bookItem := model.BookItem{
		BookId: book.ID,
		Name:   name,
	}
	db.Find(&bookItem, "name = ? AND book_id = ?", name, book.ID)
	if bookItem.ID == 0 {
		return
	}
	ret, err := txt.Html()
	if err != nil {
		log.Println("get page txt error:", err)
		return
	}
	text, err := html2text.FromString(ret, html2text.Options{PrettyTables: true})
	if err != nil {
		log.Println(err)
		return
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	page := model.BookPage{
		BookItemId: bookItem.ID,
		Txt:        text,
	}
	db.FirstOrCreate(&page, "book_item_id = ?", bookItem.ID)
}
