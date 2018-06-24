package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "sort"
  "time"
  "html/template"

  "github.com/qiniu/api.v7/auth/qbox"
  "github.com/qiniu/api.v7/kodo"
  "qiniupkg.com/api.v7/conf"
  "qiniupkg.com/api.v7/kodocli"

  "github.com/Piasy/ghrc/api"
  "github.com/Piasy/ghrc/models"
)

func check(err error, msg string) {
  if err != nil {
    fmt.Println(err)
    panic(msg)
  }
}

func main() {
  page := 1
  users := make([]*models.User, 0)
  for page <= 34 {
    onePage, _ := api.GetUsers(page, os.Args[1])
    users = append(users, onePage...)
    page++
    time.Sleep(time.Second)
  }
  sort.Sort(sort.Reverse(models.UserByDashboard(users)))
  users = append(users, &models.User{})
  now := time.Now()
  rank := models.Ranks{UpdatedAt: fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day()), Ranks: users}
  data, _ := json.Marshal(rank)
  err := ioutil.WriteFile("./ghrc.json", data, 0644)
  check(err, "WriteFile ./ghrc.json fail")

  tpl, err := ioutil.ReadFile("./index.tpl")
  check(err, "read index.tpl fail")
  t, err := template.New("index").Parse(string(tpl))
  check(err, "Parse index.tpl fail")
  tplData := struct {
    UpdateTime string
    Ranks      []*models.User
  }{
    UpdateTime: rank.UpdatedAt,
    Ranks: rank.Ranks,
  }
  out, err := os.Create("./index.html")
  check(err, "Create index.html fail")
  err = t.Execute(out, tplData)
  check(err, "instantiate tpl fail")

  //初始化AK，SK
  conf.ACCESS_KEY = os.Args[2]
  conf.SECRET_KEY = os.Args[3]

  //创建一个Client
  c := kodo.New(0, nil)

  //生成一个上传token
  token := c.MakeUptoken(&kodo.PutPolicy{
    Scope: "ghrc:ghrc.json",
    //设置Token过期时间
    Expires: 180,
  })

  //构建一个uploader
  zone := 0
  uploader := kodocli.NewUploader(zone, nil)

  var ret models.PutRet
  //调用PutFileWithoutKey方式上传，没有设置saveasKey以文件的hash命名
  res := uploader.PutFile(nil, &ret, token, "ghrc.json", "./ghrc.json", nil)
  //打印返回的信息
  fmt.Println(ret)
  //打印出错信息
  if res != nil {
    fmt.Println("io.Put failed:", res)
    return
  }
  //生成一个上传token
  token = c.MakeUptoken(&kodo.PutPolicy{
    Scope: "ghrc:index.html",
    //设置Token过期时间
    Expires: 180,
  })
  res = uploader.PutFile(nil, &ret, token, "index.html", "./index.html", nil)
  //打印返回的信息
  fmt.Println(ret)
  //打印出错信息
  if res != nil {
    fmt.Println("io.Put failed:", res)
    return
  }

  // 刷新CDN缓存
  mac := qbox.NewMac(os.Args[2], os.Args[3])
  req, _ := http.NewRequest("POST", "http://fusion.qiniuapi.com/refresh", bytes.NewBuffer([]byte("{\"urls\":[\"http://ghrc.babits.top/ghrc.json\", \"http://ghrc.babits.top/index.html\"]}")))
  qboxToken, _ := mac.SignRequest(req, false)
  req.Header.Add("Host", "fusion.qiniuapi.com")
  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("Authorization", "QBox " + qboxToken)
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    fmt.Println("Refresh cache fail:", err)
    return
  }
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println("Refresh cache succeed", string(body))
}
