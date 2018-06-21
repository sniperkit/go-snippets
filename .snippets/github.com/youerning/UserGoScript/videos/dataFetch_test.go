package videos

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/jinzhu/gorm"
	// "os"
	"testing"
)

const dbDestTest = "video-resource-test.db"

func Test_const(t *testing.T) {
	fmt.Println(dbDest)
}

func Test_Fetch(t *testing.T) {
	t.Log("fetch start")
	Fetch(10, dbDestTest)
	t.Log("fetch done")
}

func Test_CatalogStore(t *testing.T) {
	var err error
	var ty Ty
	DBInit(dbDestTest)
	db, err := gorm.Open("sqlite3", dbDestTest)
	CheckErr(err)
	defer db.Close()
	catalogUrl := apiSite2 + catalogSuffix
	t.Log(catalogUrl)
	catalog := GetCatalog(catalogUrl)
	for _, ty = range catalog.Class.TyRet {
		db.Create(&ty)
	}
	if ty.ID > 1 {
		t.Log("catalog store success!")
	} else {
		t.Error("catalog store failed!")
	}
}

func Test_Api(t *testing.T) {
	var err error
	var retJson []byte
	catalogUrl := apiSite2 + catalogSuffix
	t.Log(catalogUrl)
	catalog := GetCatalog(catalogUrl)
	retJson, err = json.MarshalIndent(catalog, " ", " ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(retJson))

	videoUrl := apiSite2 + fmt.Sprintf(videoSuffix, 1)
	t.Log(videoUrl)
	videos := GetVideoList(videoUrl)
	retJson, err = json.MarshalIndent(videos.List.Videos[0], " ", " ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(retJson))
}

func Test_XmlParse(t *testing.T) {
	var err error
	var retJson []byte
	rss := Rss{}
	err = xml.Unmarshal([]byte(TestXMLRss), &rss)
	if err != nil {
		t.Error(err)
	}
	t.Log("marshal rss")
	// fmt.Println(rss)
	retJson, err = json.MarshalIndent(rss, " ", " ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(retJson))

	rsscatalog := RssCatalog{}
	err = xml.Unmarshal([]byte(TestXMLRssCatalog), &rsscatalog)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("marshal rsscatalog")
	// fmt.Println(rss)
	retJson, err = json.MarshalIndent(rsscatalog, " ", " ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(retJson))
}
