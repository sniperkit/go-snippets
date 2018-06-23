##https://api.slack.com/methods/auth.test/test##
在这个页面上选择auth.test方法，然后点击TestMethod就可以看到自己的UserID了。

##其实这里还是通过自己的Token来获取自己用户信息
##API 就是https://slack.com/api/auth.test参数为用户的Token,返回的response为用户的个人信息

```go
package main

import (
	"fmt"
	"github.com/nlopes/slack"
)

var token = "YOUR-TOKEN"

func main() {
	api := slack.New(token) //Here use your own token
	groups, err := api.GetGroups(false)
	if err != nil {
		panic(err)
	}

	auth, err := api.AuthTest() //The reponse will include your user information
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", auth)
}
```
