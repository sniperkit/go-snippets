package videos

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

type Rss struct {
	// XMLName xml.Name  `xml:"rss"`
	Version string    `xml:"version,attr"`
	List    VideoList `xml:"list"`
}

type VideoList struct {
	// XMLName     xml.Name `xml:"list"`
	Page        int        `xml:"page,attr"`
	Pagecount   int        `xml:"pagecount,attr"`
	PageSize    int        `xml:"pagesize,attr"`
	Recordcount int        `xml:"recordcount,attr"`
	Videos      []VideoXml `xml:"video"`
}

type VideoXml struct {
	// XMLName  xml.Name `xml:"video"`
	Last     string   `xml:"last" json:"last"`
	VID      uint     `xml:"id" json:"vid" gorm:"index"`
	Tid      int      `xml:"tid" json:"tid" gorm:"index"`
	Name     string   `xml:"name" json:"name"`
	Type     string   `xml:"type" json:"type"`
	Pic      string   `xml:"pic" json:"pic"`
	Lang     string   `xml:"lang" json:"lang"`
	Area     string   `xml:"area" json:"area"`
	Year     string   `xml:"year" json:"year"`
	State    string   `xml:"state" json:"state"`
	Note     string   `xml:"note" json:"note"`
	Actor    string   `xml:"actor" json:"actor"`
	Director string   `xml:"director" json:"director"`
	DL       DataList `xml:"dl" json:"datalist"`
	Desc     string   `xml:"des" json:"des"`
}

type VideoORM struct {
	gorm.Model
	VideoXml
	// XMLName  xml.Name `xml:"video"`
	DL string `xml:"dl json:"dl"`
}

func (VideoORM) TableName() string {
	return "videos"
}

type DataList struct {
	// XMLName xml.Name  `xml:"dl" json:"dl"`
	Data []DataRet `xml:"dd" json:"data"`
}

type DataRet struct {
	// XMLName xml.Name `xml:"dd"`
	Flag    string `xml:"flag,attr" json:"flag"`
	Content string `xml:",cdata" json:"content"`
}

type RssCatalog struct {
	// XMLName xml.Name  `xml:"rss"`
	Version string      `xml:"version,attr"`
	List    CatalogList `xml:"list"`
	Class   TyList      `xml:"class"`
}

type TyList struct {
	// XMLName xml.Name `xml:"ty"`
	TyRet []Ty `xml:"ty"`
}

type Ty struct {
	// gorm.Model
	// ID        uint
	CreatedAt time.Time
	ID        uint   `xml:"id,attr"`
	Name      string `xml:",chardata"`
}

func (Ty) TableName() string {
	return "catalog"
}

type CatalogList struct {
	// XMLName     xml.Name `xml:"list"`
	Page        int `xml:"page,attr"`
	Pagecount   int `xml:"pagecount,attr"`
	PageSize    int `xml:"pagesize,attr"`
	Recordcount int `xml:"recordcount,attr"`
}

func DBInit(dbFIle string) {
	db, err := gorm.Open("sqlite3", dbFIle)
	CheckErr(err)
	defer db.Close()
	db.AutoMigrate(&VideoORM{})
	db.AutoMigrate(&Ty{})
}
