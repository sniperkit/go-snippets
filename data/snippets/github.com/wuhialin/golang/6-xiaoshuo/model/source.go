package model

type Source struct {
	Id                int64
	Name, Domain, Url string
	MainAction        string `mapstructure:"main_action"`
	State             int
	CreateAt          int64 `mapstructure:"create_at"`
}
