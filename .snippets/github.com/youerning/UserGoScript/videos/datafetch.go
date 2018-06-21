package videos

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Fetch(num int, dbFile string) {
	DBInit(dbFile)
	db, err := gorm.Open("sqlite3", dbFile)
	CheckErr(err)
	defer db.Close()

	for i := 0; i < num; i++ {
		videoUrl := ApiSite2 + fmt.Sprintf(VideoSuffix, i+1)
		fmt.Println(videoUrl)
		store(db, videoUrl)
	}
}

func JsonOut(v VideoORM) string {
	out := VideoXml{}
	out.Last = v.Last
	out.VID = v.VID
	out.Tid = v.Tid
	out.Name = v.Name
	out.Type = v.Type
	out.Pic = v.Pic
	out.Lang = v.Lang
	out.Area = v.Area
	out.Year = v.Year
	out.State = v.State
	out.Note = v.Note
	out.Actor = v.Actor
	out.Director = v.Director
	tmpDl := DataList{}
	fmt.Println(v.DL)
	err := json.Unmarshal([]byte(v.DL), &tmpDl)
	CheckErr(err)
	out.DL = tmpDl
	out.Desc = v.Desc
	fmt.Println(out)
	ret, err := json.Marshal(out)
	return string(ret)
}

func V2VXML(v []VideoORM) []VideoXml {
	vLis := []VideoXml{}
	for _, v := range v {
		out := VideoXml{}
		out.Last = v.Last
		out.VID = v.ID
		out.Tid = v.Tid
		out.Name = v.Name
		out.Type = v.Type
		out.Pic = v.Pic
		out.Lang = v.Lang
		out.Area = v.Area
		out.Year = v.Year
		out.State = v.State
		out.Note = v.Note
		out.Actor = v.Actor
		out.Director = v.Director
		tmpDl := DataList{}
		// fmt.Println(v.DL)
		err := json.Unmarshal([]byte(v.DL), &tmpDl)
		CheckErr(err)
		out.DL = tmpDl
		out.Desc = v.Desc
		vLis = append(vLis, out)
	}

	return vLis
}

func CatalogStore(dbFile string) {
	var err error
	DBInit(dbFile)
	db, err := gorm.Open("sqlite3", dbFile)
	CheckErr(err)
	defer db.Close()
	catalogUrl := ApiSite2 + CatalogSuffix
	fmt.Println(catalogUrl)
	catalog := GetCatalog(catalogUrl)
	// fmt.Println(catalog.Class.TyRet)
	for _, ty := range catalog.Class.TyRet {
		// fmt.Println(ty)
		db.Create(&ty)
	}
}

func GetPageInfo(apiSite string) (int, int, int) {
	url := apiSite + fmt.Sprintf(VideoSuffix, 1)
	r := GetVideoList(url)
	return r.List.PageSize, r.List.Pagecount, r.List.Recordcount
}

func store(db *gorm.DB, url string) {
	var err error
	VideoRss := GetVideoList(url)
	var tmpDl []byte
	for _, v := range VideoRss.List.Videos {
		video := &VideoORM{}
		video.Last = v.Last
		video.VID = v.VID
		video.Tid = v.Tid
		video.Name = v.Name
		video.Type = v.Type
		video.Pic = v.Pic
		video.Lang = v.Lang
		video.Area = v.Area
		video.Year = v.Year
		video.State = v.State
		video.Note = v.Note
		video.Actor = v.Actor
		video.Director = v.Director
		tmpDl, err = json.Marshal(v.DL)
		CheckErr(err)
		video.DL = string(tmpDl)
		video.Desc = v.Desc
		// fmt.Println("video===>", v)
		db.Create(video)
	}
}

func GetCatalog(url string) *RssCatalog {
	var ret = RssCatalog{}
	resp, err := http.Get(url)
	CheckErr(err)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	err = xml.Unmarshal(data, &ret)
	CheckErr(err)
	return &ret
}

func GetVideoList(url string) *Rss {
	var ret = Rss{}
	resp, err := http.Get(url)
	CheckErr(err)
	if resp.StatusCode != 200 {
		data, err := ioutil.ReadAll(resp.Body)
		if len(data) > 1024 {
			fmt.Println(string(data)[:1024])
		} else {
			fmt.Println(string(data))
		}
		time.Sleep(10 * time.Second)
		resp, err = http.Get(url)
		CheckErr(err)

	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	err = xml.Unmarshal(data, &ret)
	CheckErr(err)
	return &ret
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
