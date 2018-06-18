package main

import (
	"fmt"
	"net/url"
	"log"
	"strings"
)

func main() {
	uri := "merchantId=687&orderId=1&email=af%40suretly.com&orderActualTill=2017-12-27%2006%3A34%3A19Z&price=1.02&callbackUrl=https%3A%2F%2Fdemo.suretly.io%3A3015%2Ftransaction_callback&returnUrl=https%3A%2F%2Fdemo.suretly.io%3A3015&action=pay&customer_fullName=%20%20&customer_phone=%2B%2B79133819546&customer_email=af%40suretly.com&transaction=37307991d4ee425092648aad15c9036e&object_type=transaction&status=success&payment_system=mandarinpayv1&sandbox=true&cb_processed_at=2017-12-25T06%3A34%3A23.2703470&831366e5-9b1e-4aeb-9ec6-3a76e3e09708=7b8e55e2-2cd3-4aa1-9041-8f1b4ecbb0c4&sign=265142320da58f14fbfdc4debd4985b6ac5766972f09446fed2c6b70a09c2134"
	uri, err := url.QueryUnescape(uri)
	fmt.Println(uri)

	u, err := url.Parse("http://bing.com/search/add?q=dotnet&a=1")
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "https"
	u.Host = "google.com"
	q := u.Query()
	q.Set("q", "golang")
	u.RawQuery = q.Encode()
	fmt.Println(u)

	a := strings.Split(u.Path, "/")
	if len(a) > 0 {
		fmt.Println("a:",a[len(a)-1])
	}

	fmt.Println(u.Path)
}
