package main

import "github.com/gin-gonic/gin"

func getStatus(c *gin.Context) {
	content := gin.H{"status": "ok"}
	c.IndentedJSON(200, content)
}
