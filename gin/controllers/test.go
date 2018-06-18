package controllers

import (
    //"fmt"
    "net/http"

    "github.com/gin-gonic/gin"
)

type Test struct {
}

func (t *Test) Show(c *gin.Context) {
    //example := c.MustGet("example").(string)
    //fmt.Println(example)
    //c.Set("example2", "67890")
    c.JSON(http.StatusOK, gin.H{"message": "hello world", "status": http.StatusOK})
}
