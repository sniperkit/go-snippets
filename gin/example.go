package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Default With the Logger and Recovery middleware already attached
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := DB[user]
		if ok {
			c.JSON(200, gin.H{"user": user, "value": value})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	})

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			DB[user] = json.Value
			c.JSON(200, gin.H{"status": "ok"})
		}
	})

	// Default route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hi": "jiangew",
		})
	})

	// Params in path
	// This handler will match /user/john but will not match neither /user/ or /user
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hi, %s", name)
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " to " + action
		c.String(http.StatusOK, message)
	})

	// QueryString Params
	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?fname=Erwei&lname=Jiangew
	r.GET("/welcome", func(c *gin.Context) {
		fname := c.DefaultQuery("fname", "Guest")
		// lname := c.Request.URL.Query().Get("lname")
		lname := c.Query("lname")
		c.String(http.StatusOK, "Hi, %s %s", fname, lname)
	})

	// Multipart / Urlencoded Form
	r.POST("/form_post", func(c *gin.Context) {
		message := C.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")
		c.JSON(http.StatusOK, gin.H{
			"status": "posted",
			"message": message,
			"nick": nick,
		})
	})

	// Query + Post Form
	// POST /post?id=1234&page=1 HTTP/1.1 
	// Content-Type: application/x-www-form-urlencoded
	// name=manu&message=this_is_great
	r.POST("/post", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")

		fmt.Printf("id: %s, page: %s, name: %s, message: %s", id, page, name, message)
	})

	// Upload file
	// How to curl:
	// curl -X POST http://localhost:8080/upload \
	// 	-F "file=@/Users/appleboy/test.zip" \
  	// 	-H "Content-Type: multipart/form-data"
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Static("/", "./public")
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		if err := c.SaveUploadedFile(file, file.Filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("Uploaded successfully fileName: %s", file.Filename))
	})

	// Multiple files
	// How to curl:
	// curl -X POST http://localhost:8080/upload \
  	// 	-F "upload[]=@/Users/appleboy/test1.zip" \
  	// 	-F "upload[]=@/Users/appleboy/test2.zip" \
	// 	-H "Content-Type: multipart/form-data"
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Static("/", "./public") 
	r.POST("/upload", func(c *gin.Context) {
		// Multipart form
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		files := form.File["files"]
		for _, file := range files {
			if err := c.SaveUploadedFile(file, file.Filename); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}
		}

		c.String(http.StatusOK, fmt.Sprintf("Uploaded successfully fileCount: %s", len(files)))
	})

	// Grouping routes
	// Simple group: v1
	v1 := r.Group("/v1") 
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}
	// Simple group: v2
	v2 := r.Group("/v2") 
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
		v2.POST("/read", readEndpoint)
	}

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
