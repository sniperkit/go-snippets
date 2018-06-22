package wpdb

import "html/template"

// TODO: configure types correctly and better relationships
type Option struct {
	Id       int64  `gorm:"column:option_id; primary_key:yes"`
	Name     string `gorm:"column:option_name" sql:"size:64"`
	Value    string `gorm:"column:option_value" sql:"type:longtext"`
	Autoload string `gorm:"column:autoload" sql:"size:20"`
}

func (o Option) TableName() string {
	return "wp_options"
}

type User struct {
	Id               int64  `gorm:"column:ID; primary_key:yes"`
	HumanizedName    string `gorm:"column:user_nicename" sql:"size:50"`
	Username         string `gorm:"column:display_name" sql:"size:250"`
	LoginName        string `gorm:"column:user_login" sql:"size:60"`
	Password         string `gorm:"column:user_pass" sql:"size:64"`
	Mail             string `gorm:"column:user_email" sql:"size:100"`
	Homepage         string `gorm:"column:user_url" sql:"size:100"`
	RegistrationDate string `gorm:"column:user_registered"`
	ActivationKey    string `gorm:"column:user_activation_key" sql:"size:60"`
	Online           bool   `gorm:"column:user_status"`
}

func (u User) TableName() string {
	return "wp_users"
}

type Post struct {
	Id       int64  `gorm:"column:ID; primary_key:yes"`
	UserId   int64  `gorm:"column:post_author"`
	PostDate string `gorm:"column:post_date"`
	Content  string `gorm:"column:post_content" sql:"type:longtext" json:",omitempty"` // The comma for JSON is needed.
	Title    string `gorm:"column:post_title" sql:"type:text"`
	Name     string `gorm:"column:post_name" sql:"size:200"`
}

func (p Post) TableName() string {
	return "wp_posts"
}

func (p Post) ContentAsHTML() template.HTML {
	return template.HTML(p.Content)
}

type Comment struct {
	Id int64 `gorm:"column:comment_ID; primary_key:yes" json:"id"`
	// Get rid of 'Author' and instead let the User variable point to AnonymousUser
	Author         string `gorm:"column:comment_author" sql:"type:tinytext" json:"author_name"`
	AuthorMail     string `gorm:"column:comment_author_email" sql:"size:100" json:"author_mail"`
	AuthorHomepage string `gorm:"column:comment_author_url" sql:"size:100" json:"author_url"`
	Content        string `gorm:"column:comment_content" json:"content"`
	UserId         int64  `gorm:"column:user_id"`
	PostId         int64  `gorm:"column:comment_post_ID"`
}

func (c Comment) TableName() string {
	return "wp_comments"
}
