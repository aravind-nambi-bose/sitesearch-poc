package main

import (
	"github.com/gin-gonic/gin"
	"sitesearch/service"
)

func main() {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, "all good")
	})

	router.PUT("/index/:indexName", service.HandlePutIndex)
	router.PUT("/index/:indexName/synonyms", service.HandlePutSynonyms)
	router.POST("/index/:indexName/products", service.HandlePostProducts)

	router.GET("/index/:indexName", service.HandleSearchProducts)

	router.Run(":8090")
}
