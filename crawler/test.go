package videos

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	// "strconv"
	// "flag"
	"fmt"
	"net/http"
	"net/url"
	// "os"
	"time"
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
	jobs := make(chan int, 10)
	results := make(chan int, 10)

	for w := 0; w < 5; w++ {
		go worker(w, jobs, results)
	}

	time.Sleep(2)
	fmt.Println("push job")
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 0; a < 5; a++ {
		<-results
	}

}

func worker(w int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", w, "start job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", w, "finished job", j)
		results <- j * 2
	}

}

func checkProxy() {
	proxyUrl, err := url.Parse("http://42.245.252.35:80")
	proxyClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	res, err := proxyClient.Head("http://www.baidu.com")
	checkErr(err)
	fmt.Println(res.StatusCode)
}
func dbInit() {
	db, err := gorm.Open("sqlite3", "test.db")
	checkErr(err)
	defer db.Close()
	db.AutoMigrate(&Video{})
}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
