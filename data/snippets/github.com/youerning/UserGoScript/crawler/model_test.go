package videos

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"strconv"
)

type Video struct {
	gorm.Model
	Url        string `gorm:"index`
	ImageSrc   string
	Name       string `gorm:"index"`
	Alias      string
	Desc       string
	VideoSrc   string
	ThunderSrc string
}

func main() {
	db, err := gorm.Open("sqlite3", "80.db")
	checkErr(err)
	defer db.Close()
	db.AutoMigrate(&Video{})
	// var v *Video

	// for i := 1; i <= 10; i++ {
	// 	url := strconv.Itoa(i)
	// 	v = &Video{Url: url}
	// 	db.Create(v)
	// }

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
