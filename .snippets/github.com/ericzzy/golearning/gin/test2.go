package main

import (
    //"time"
     "fmt"
     "os"
     "runtime"

    "github.com/gin-gonic/gin"
    "github.com/ericzzy/golearning/gin/controllers"
)


//func Logger() gin.HandlerFunc {
//    return func(c *gin.Context) {
//        t := time.Now()
//        c.Set("example", "12345")
//        c.Next()
//        latency := time.Since(t)
//        fmt.Println(latency)
//        status := c.Writer.Status()
//        fmt.Println("status is", status)
//        if example2, exists := c.Get("example2"); exists {
//            fmt.Println("example2 is", example2)
//        }
//    }
//}


func main() {
    runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)
    fd, err := os.OpenFile("/var/log/gin.log", os.O_RDWR | os.O_APPEND | os.O_CREATE , 0666)
    defer fd.Close()
    if err != nil {
        fmt.Println("Log file could not be opened/created")
        os.Exit(1)
    }
    gin.DefaultWriter = fd
    gin.SetMode(gin.ReleaseMode)
    //r := gin.New()
    r := gin.Default()
 
    test := &controllers.Test{}
    r.GET("/test", test.Show) 

    r.Run(":8090")
}
