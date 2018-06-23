
type ViewFunc func(http.ResponseWriter, *http.Request)

func BasicAuth(f ViewFunc) ViewFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := database.DB()
		CheckErr(err, w)
		defer db.Close()
		timestamp := r.URL.Query().Get("timestamp")
		t := time.Now()
		fmt.Println(t)
		tm, _ := time.Parse("20060102150405", timestamp)
		fmt.Println(tm)
		x := t.Sub(tm).Minutes()
		fmt.Println(x)
		basicAuthPrefix := "Basic "
		auth := r.Header.Get("Authorization") //获取 request header
		if strings.HasPrefix(auth, basicAuthPrefix) {
			payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):]) // 解码认证信息
			if err == nil {
				account := model.Account{}
				pair := bytes.SplitN(payload, []byte(":"), 2)
				appid := string(pair[0])
				db.Joins("JOIN apps ON apps.account_id = accounts.id").Where("apps.sid = ?", appid).Find(&account)
				sign_str := appid + account.AuthToken + timestamp
				sha1Encoder := sha1.New()
				sha1Encoder.Write([]byte(sign_str))
				signResult := sha1Encoder.Sum(nil)
				sign := fmt.Sprintf("%x", signResult)
				fmt.Println(sign)
				if len(pair) == 2 && bytes.Equal(pair[0], []byte(appid)) && bytes.Equal(pair[1], []byte(sign)) {
					f(w, r)
					return
				}
			}
		}
		// 认证失败，提示 401 Unauthorized
		// Restricted 可以改成其他的值
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func CheckErr(err error, res http.ResponseWriter) {
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}
