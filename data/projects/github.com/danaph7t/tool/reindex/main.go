package main

import (
	"encoding/json"
	"log"
	// "sync/atomic"
	"time"

	"github.com/btlike/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/olivere/elastic.v3"
)

var (
	db *gorm.DB
	es *elastic.Client
)

//重建索引，在某些情况下，es与数据库中数据不一致，需要重建全部索引
func main() {
	initial()
	// reindex()
	bulkreindex()
}

func initial() {
	var (
		err    error
		conn   = "root:password@tcp(127.0.0.1:3306)/torrent?charset=utf8&parseTime=True&loc=Local"
		esaddr = "http://127.0.0.1:9200"
	)
	db, err = gorm.Open("mysql", conn)
	if err != nil {
		log.Fatal(err)
	}
	es, err = elastic.NewClient(elastic.SetURL(esaddr))
	if err != nil {
		log.Fatal(err)
	}
	es.CreateIndex("torrent").Do()
}

type torrent struct {
	ID         int64     `gorm:"column:id"`
	Infohash   string    `gorm:"column:infohash"`
	Data       string    `gorm:"column:data"`
	CreateTime time.Time `gorm:"column:create_time"`
}

type torrentSearch struct {
	Name       string
	Length     int64
	Heat       int64
	CreateTime time.Time
}

type tableManager struct {
	name  string
	index int
}

func bulkreindex() {
	var tables []tableManager
	tables = append(tables, tableManager{name: "0", index: 0})
	tables = append(tables, tableManager{name: "1", index: 0})
	tables = append(tables, tableManager{name: "2", index: 0})
	tables = append(tables, tableManager{name: "3", index: 0})
	tables = append(tables, tableManager{name: "4", index: 0})
	tables = append(tables, tableManager{name: "5", index: 0})
	tables = append(tables, tableManager{name: "6", index: 0})
	tables = append(tables, tableManager{name: "7", index: 0})
	tables = append(tables, tableManager{name: "8", index: 0})
	tables = append(tables, tableManager{name: "9", index: 0})
	tables = append(tables, tableManager{name: "a", index: 0})
	tables = append(tables, tableManager{name: "b", index: 0})
	tables = append(tables, tableManager{name: "c", index: 0})
	tables = append(tables, tableManager{name: "d", index: 0})
	tables = append(tables, tableManager{name: "e", index: 0})
	tables = append(tables, tableManager{name: "f", index: 0})

	for _, manager := range tables {
		table := "torrent" + manager.name
		log.Println("begin table ", table)

		var data []torrent
		for i := manager.index; i < 20; i++ {
			log.Println("begin row ", i)
			err := db.Table(table).Limit(50000).Offset(i * 50000).Find(&data).Error
			if err != nil {
				log.Println(err)
				return
			}
			if len(data) == 0 {
				break
			}

			var beforeRequests int64
			var befores int64
			var afters int64
			var failures int64
			var afterRequests int64
			/*
				beforeFn := func(executionId int64, requests []elastic.BulkableRequest) {
					atomic.AddInt64(&beforeRequests, int64(len(requests)))
					atomic.AddInt64(&befores, 1)
				}*/
			afterFn := func(executionId int64, requests []elastic.BulkableRequest, response *elastic.BulkResponse, err error) {
				// atomic.AddInt64(&afters, 1)
				if err != nil {
					log.Println(err)
					// atomic.AddInt64(&failures, 1)
				}
				// atomic.AddInt64(&afterRequests, int64(len(requests)))
			}

			// p, err := es.BulkProcessor().BulkSize(-1).BulkActions(-1).Before(beforeFn).After(afterFn).Do()
			p, err := es.BulkProcessor().BulkSize(-1).BulkActions(-1).After(afterFn).Do()
			if err != nil {
				log.Fatal(err)
			}

			var search = torrentSearch{}
			var t repository.Torrent
			for _, v := range data {
				err = json.Unmarshal([]byte(v.Data), &t)
				if err != nil {
					// log.Println(err, t, v.ID)
					// db.Table(table).Where("infohash=?", v.Infohash).Delete(&t)
					continue
				}
				search = torrentSearch{
					Name:       t.Name,
					Length:     t.Length,
					CreateTime: t.CreateTime,
					Heat:       1,
				}
				request := elastic.NewBulkIndexRequest().Index("torrent").Type(manager.name).Id(v.Infohash).Doc(search)
				p.Add(request)
			}
			log.Println("begin flush")
			p.Flush()
			log.Println("finish flush")

			err = p.Close()
			if err != nil {
				log.Fatal(err)
			}
			if failures != 0 {
				log.Println(befores, afters, beforeRequests, afterRequests, "fail", failures)
			}
			data = make([]torrent, 0)
		}
	}
}

func reindex() {
	tables := []string{"c", "d", "e", "f"}
	for _, v := range tables {
		table := "torrent" + v
		log.Println("begin table ", table)
		for i := 12; i < 100; i++ {
			log.Println("begin row ", i)
			var data []torrent
			err := db.Table(table).Limit(10000).Offset(i * 10000).Find(&data).Error
			if err != nil {
				log.Println(err)
				return
			}
			if len(data) == 0 {
				break
			}

			for _, v := range data {
				var t repository.Torrent
				err = json.Unmarshal([]byte(v.Data), &t)
				if err != nil {
					log.Println(err, t, v.ID)
					// db.Table(table).Where("infohash=?", v.Infohash).Delete(&t)
					continue
				}
				search := torrentSearch{
					Name:       t.Name,
					Length:     t.Length,
					CreateTime: t.CreateTime,
					Heat:       1,
				}
				_, err = es.Index().
					Index("torrent").
					Type("infohash").
					Id(v.Infohash).
					BodyJson(search).
					Refresh(false).
					Do()
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
