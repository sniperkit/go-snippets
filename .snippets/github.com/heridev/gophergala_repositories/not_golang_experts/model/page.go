package model

import (
	"strconv"
	"time"
)

type Page struct {
	Id            int64
	Url           string
	LastCheckedAt time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	HtmlString    string
}

type UserSubscribed struct {
	Email string
}

func PagesToCheck() []*Page {
	DB.LogMode(false)

	var pages []*Page
	DB.Where("last_checked_at < ?", time.Now().Add(-time.Minute*2)).Find(&pages)
	return pages
}

func FindOrCreatePageByUrl(url string) Page {
	page := Page{LastCheckedAt: time.Now()}
	DB.Where(Page{Url: url}).FirstOrCreate(&page)
	return page
}

func (p *Page) SubscribedUsersEmails() []string {
	var results []UserSubscribed
	querystring := "join subscriptions on subscriptions.user_id = users.id AND subscriptions.page_id = " + strconv.FormatInt(p.Id, 10)
	DB.Table("users").Select("users.email").Joins(querystring).Scan(&results)

	var emails = make([]string, len(results))
	for i, result := range results {
		emails[i] = result.Email
	}
	return emails
}

func (p Page) Save() {
	DB.Save(&p)
}
