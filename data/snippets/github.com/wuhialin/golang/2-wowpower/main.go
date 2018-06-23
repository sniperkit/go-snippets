package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/scorredoira/email"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var write *csv.Writer

func main() {
	dir := os.TempDir()
	f, err := ioutil.TempFile(dir, "")
	if err != nil {
		log.Fatalln(err)
	}
	w := csv.NewWriter(f)
	filename := f.Name()
	defer func() {
		write.Flush()
		subject := "test"
		body := "test"
		to := []string{"wuhialin@vip.qq.com"}
		tmpFile := new(email.Attachment)
		tmpFile.Filename = "games.csv"
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalln(err)
		}
		tmpFile.Data = content
		attach := map[string]*email.Attachment{}
		attach["game.csv"] = tmpFile
		sendMail(to, subject, body, attach)
		err = os.Remove(filename)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	defer f.Close()
	w.Write([]string{"名称", "价格", "时间", "利率", "万元获得"})
	write = w
	query()
}

var createDate = time.Now().Format("2006-01-02")

func formatter(body []byte) {
	var games []map[string]string
	body = regexp.MustCompile(`\n`).ReplaceAll(body, []byte(""))
	body = regexp.MustCompile(`\s+`).ReplaceAll(body, []byte(" "))

	nameReg := regexp.MustCompile(`<dd\s+class="game-name">(.+?)</dd>`)
	for _, q := range nameReg.FindAllSubmatch(body, -1) {
		tmpMap := map[string]string{}
		tmpMap["name"] = string(q[1])
		games = append(games, tmpMap)
	}

	priceReg := regexp.MustCompile(`<span>押金：(\d+)纳币</span>`)
	for k, q := range priceReg.FindAllSubmatch(body, -1) {
		games[k]["price"] = string(q[1])
	}

	dayReg := regexp.MustCompile(`<span>周期：(\d+)天</span>`)
	for k, q := range dayReg.FindAllSubmatch(body, -1) {
		games[k]["day"] = string(q[1])
	}

	gainProfitReg := regexp.MustCompile(`<span>返利：(\d+(\.\d+)?)纳币</span>`)
	for k, q := range gainProfitReg.FindAllSubmatch(body, -1) {
		games[k]["gain_profit"] = string(q[1])
	}

	updateTime := time.Now().Unix()
	for _, row := range games {
		gainProfit, _ := strconv.ParseFloat(row["gain_profit"], 32)
		day, _ := strconv.ParseFloat(row["day"], 32)
		price, _ := strconv.ParseFloat(row["price"], 32)
		rage := gainProfit / day / price * 10000
		write.Write([]string{row["name"], row["price"], row["day"], row["gain_profit"], strconv.FormatFloat(float64(rage), 'f', 2, 64)})
		if hasByName(row["name"]) > 0 {
			continue
		}
		sqlString := `INSERT INTO wowpower_game (name, price, day, gain_profit, create_date, update_time) VALUE (?, ?, ?, ?, ?, ?)`
		_, err := db().Exec(sqlString, row["name"], row["price"], row["day"], row["gain_profit"], createDate, updateTime)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func query() {
	client := &http.Client{}
	client.Timeout = 3 * time.Second
	baseUrl := "http://www.wowpower.com/showNewGame?page="
	i := 1
	maxPage := 1 << 10
	flag := true
	for {
		response, err := client.Get(baseUrl + strconv.Itoa(i))
		i++
		if err != nil {
			continue
		}
		if response.StatusCode == http.StatusOK {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				continue
			}
			if flag {
				flag = false
				reg := regexp.MustCompile(`<span\s+id="gamesPage".*?>\s*(\d+)\s*</span>`)
				q := reg.FindSubmatch(body)
				if q != nil {
					maxPage, _ = strconv.Atoi(string(q[1]))
				}
			}
			formatter(body)
		}
		if i > maxPage {
			break
		}
	}
}

var connDb *sql.DB

func db() *sql.DB {
	if connDb != nil {
		return connDb
	}
	name := "mysql"
	driver := "root:123456@tcp(127.0.0.1)/test?charset=utf8"
	db, err := sql.Open(name, driver)
	if err != nil {
		log.Fatalln(err)
	}
	connDb = db
	return connDb
}

func hasByName(name string) (id int) {
	queryString := `SELECT id FROM wowpower_game WHERE name = ? AND create_date = ? LIMIT 1`
	rows, err := db().Query(queryString, name, createDate)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			log.Fatalln(err)
		}
	}
	return
}

func config() map[string]map[string]string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	content, err := ioutil.ReadFile(filepath.Join(path, "config.json"))
	if err != nil {
		log.Fatalln(err)
	}
	config := make(map[string]map[string]string)
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalln(err)
	}
	return config
}

func sendMail(to []string, subject, body string, attach map[string]*email.Attachment) {
	config := config()
	mail := config["mail"]
	_, err := smtp.Dial(mail["host"])
	if err != nil {
		log.Fatalln(err)
	}
	ht := strings.Split(mail["host"], ":")
	auth := smtp.PlainAuth("", mail["user"], mail["pass"], ht[0])
	m := email.NewMessage(subject, body)
	m.From.Address = mail["user"]
	m.To = to
	if len(attach) > 0 {
		m.Attachments = attach
	}
	if err := email.Send(mail["host"], auth, m); err != nil {
		log.Fatalln(err)
	}
}
