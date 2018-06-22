package main


import (
    "github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	SecretKey = "WangruiAlert"
)

func fatal (err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type UserCredentials struct {
	Username  string  `json:"username"`
	Password  string  `json:"password"`
}

type User struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Username string     `json:"username"`
	Password string     `json:"password"`
}

type Response struct {
	Data   string   `json:"data"`
}

type Token struct {
	Token  string   `json:"token"`
}

func StartServer() {
	http.HandleFunc("/login", LoginHandler)

	// TODO: 这里需了解如果使用negroni作为中间件
	//
	//
	http.Handle("/resource", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	))

	log.Println("Now listening...")
	http.ListenAndServe(":8080", nil)
}

func main() {
	StartServer()
}



func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user UserCredentials
        fmt.Println(r.Body)

	err := json.NewDecoder(r.Body).Decode(&user)
	// 将request中的账号信息传入创建的user结构体中
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	// 判断传进服务器的账号密码是否可用
	if strings.ToLower(user.Username) != "wangrui" {
		if user.Password != "123456" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("User info not validated")
			fmt.Fprint(w, "Invalid Credentials")
			return
		}
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error extracting the key")
		fatal(err)
	}

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}

	response := Token{tokenString}
	JsonResponse(response, w)
}


// 验证Token的中间件函数
func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error){
			return []byte(SecretKey), nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource ")
	}
}

func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
