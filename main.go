package main

import (
	"github.com/DianaBurca/fetcher/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	driver := gin.Default()
	driver.GET("/fetch", utils.FetchHandler)
	driver.GET("/.well-known/live", utils.Health)
	driver.GET("/.well-known/ready", utils.Health)
	driver.Run()
}
