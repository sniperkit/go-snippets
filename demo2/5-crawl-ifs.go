/*
* @Author: wuhailin
* @Date:   2018-01-04 10:17:59
* @Last Modified by:   wuhailin
* @Last Modified time: 2018-01-04 10:28:19
 */
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

func main() {
	f, _ := os.OpenFile("D:/beta.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	defer f.Close()
	csvWrite := csv.NewWriter(f)
	defer csvWrite.Flush()
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Parallelism: 2})
	c.SetRequestTimeout(30 * time.Minute)
	u, _ := url.Parse("http://pms.beta.gw-ec.com/order/ordersearch/details?purchase_sn=&search_type=sku&keyword=&order_type=0&purchase_type=0&purchaser=&input_purchaser_select=&input_purchaser=&purchaser_group=&level_purchaser_group=&select_level_user=&provider_type=provider_sn&provider_keyword=&purchase_deliver_sn=&cancel_min=&cancel_max=&create_time_start=2018-01-01&create_time_end=2018-01-02&over_start=&over_end=&last_start=&last_end=&min_stockin=&max_stockin=&provider_deliver_sn=&stock_id=&transfer_stock_sn=&is_new=-1&is_tax=&latest_warehouse_date_start=&latest_warehouse_date_end=&reply_arrival_date_start=&reply_arrival_date_end=&receipt_status=&time_date_type=delivery&time_date_start=&time_date_end=&butname=%E6%90%9C%E7%B4%A2&get_total=1")
	c.OnError(func(r *colly.Response, e error) {
		log.Println(e)
	})
	var stats map[string]interface{}
	c.OnResponse(func(r *colly.Response) {
		log.Println(r.Request.URL.String())
		if err := json.Unmarshal(r.Body, &stats); err != nil {
			panic(err)
		}
		var data []string
		data = append(data, r.Request.URL.Query().Get("page"))
		data = append(data, fmt.Sprintf(`%s`, stats["total_fine_str"]))
		data = append(data, fmt.Sprintf(`%s`, stats["cur_fine_total_str"]))
		csvWrite.Write(data)
		log.Println(data)
	})
	c.OnRequest(func(r *colly.Request) {
	})
	c.SetCookieJar(getCookieJar(u))
	for i := 1; i <= 27; i++ {
		c.Visit(u.String() + `&page=` + fmt.Sprintf(`%d`, i))
	}
	c.Wait()
}

func getCookieJar(u *url.URL) *cookiejar.Jar {
	cookie := &http.Cookie{}
	cookie.Name = "PHPSESSID"
	cookie.Value = "oigupt17o4hh3actjf4jj4sep3"
	cookie.Path = "/"
	cookie.Domain = u.Hostname()
	cookie.Expires = time.Now().Add(time.Hour * 24 * 365)
	c2 := cookie
	c2.Name = "pms[page_size]"
	c2.Value = "500"
	c2.Path = "/order/ordersearch"
	o := cookiejar.Options{}
	jar, err := cookiejar.New(&o)
	if err != nil {
		log.Fatalln(err)
	}
	jar.SetCookies(u, []*http.Cookie{cookie, c2})
	return jar
}
