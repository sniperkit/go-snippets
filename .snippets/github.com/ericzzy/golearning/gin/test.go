package main

import "github.com/gin-gonic/gin"
import "net/http"
import "fmt"
import "os"

func main() {
    fd, err := os.OpenFile("/var/log/gin.log", os.O_RDWR | os.O_APPEND | os.O_CREATE , 0666)
    defer fd.Close()
    if err != nil {
        fmt.Println("Log file could not be opened/created")
        os.Exit(1)
    }
    gin.DefaultWriter = fd
    router := gin.Default()

    // This handler will match /user/john but will not match neither /user/ or /user
    router.GET("/user/:name", func(c *gin.Context) {
        name := c.Param("name")
        c.String(http.StatusOK, "Hello %s", name)
    })

    // However, this one will match /user/john/ and also /user/john/send
    // If no other routers match /user/john, it will redirect to /user/join/
    router.GET("/user/:name/*action", func(c *gin.Context) {
        name := c.Param("name")
        action := c.Param("action")
        message := name + " is " + action
        c.String(http.StatusOK, message)
    })

    router.POST("/form_post", func(c *gin.Context) {
        message := c.PostForm("message")
        nick := c.DefaultPostForm("nick", "anonymous")

        c.JSON(http.StatusOK, gin.H{
            "status": "posted",
            "message": message,
            "nick": nick,
        })
    })

    router.POST("/post", func(c *gin.Context) {
        id := c.Query("id")
        page := c.DefaultQuery("page", "0")
        name := c.PostForm("name")
        message := c.PostForm("message")

        fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
    })

    router.Run(":8090")
}
