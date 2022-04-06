package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	g.GET("/hello/:name", func(c *gin.Context) {
		c.String(200, "Hello %s", c.Param("name"))
	})
	g.RunTLS(":3000", "./certs/server.crt", "./certs/server.key")
}
