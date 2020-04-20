package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"sitesearch/model"
	"strings"
)

func upsert(product *model.Product, indexName string) (*resty.Response, error) {
	client := resty.New()

	upsertStmt := map[string]interface{}{}
	upsertStmt["doc"] = product
	upsertStmt["doc_as_upsert"] = true
	upsertJson, _ := json.Marshal(upsertStmt)
	fmt.Println(string(upsertJson))

	return client.R().
		SetBody(upsertJson).
		SetHeader("Content-Type", "application/json").
		Post("http://elasticsearch:9200/" + indexName + "/_doc/" + product.BaseProduct + "/_update")
}

func HandlePutSynonyms(c *gin.Context) {
	indexName := c.Param("indexName")
	var synonymsJson = map[string][]string{}
	c.ShouldBind(&synonymsJson)

	client := resty.New()
	response, err := client.R().
		SetBody([]byte(`
				{
		    		"analysis" : {
						"filter" : {
							"synonym" : {
								"type" : "synonym",
								"synonyms" : ["`+strings.Join(synonymsJson["synonyms"], "\",\"")+`"]
							}
						}
					}
				}
		`)).
		SetHeader("Content-Type", "application/json").
		Put("http://elasticsearch:9200/" + indexName + "/_settings") //Example: product_catalog_search_en_us

	if err != nil {
		fmt.Printf("Error %s", err.Error())
		c.JSON(400, err.Error())
	} else {
		c.Data(200, "application/json", response.Body())
	}
}

func HandlePutIndex(c *gin.Context) {
	indexName := c.Param("indexName")
	client := resty.New()
	response, err := client.R().
		SetBody([]byte(`
		{
			"settings" : {
				"number_of_shards" : 1,
				"index" : {
					"analysis" : {
						"analyzer" : {
							"synonym_analyzer" : {
								"tokenizer" : "standard",
								"filter" : ["synonym"]
							}
						},
						"filter" : {
							"synonym" : {
								"type" : "synonym",
								"lenient": true,
								"updateable": true,
								"synonyms" : ["glasses, goggles => frames"]
							}
						}
					}
				}
			},
			"mappings" : {
				"properties" : {
					"baseProduct" : { "type" : "keyword" },
					"name" : { "type" : "text", "analyzer": "standard", "search_analyzer": "synonym_analyzer" },
					"productType" : { "type" : "keyword" },
					"category" : { "type" : "text" },
					"productCategories" : { "type" : "text", "fields": { "raw": { "type": "keyword"} } },
					"price" : { "type" : "float" },
					"discountPrice" : { "type" : "float" },
					"discountPercent" : { "type" : "float" },
					"promote" : { "type" : "boolean" },
					"reviewsTotal" : { "type" : "long" },
					"averageRating" : { "type" : "float" },
					"launchDate" : { "type" : "date" }
				}
			}
		}
       `)).
		SetHeader("Content-Type", "application/json").
		Put("http://elasticsearch:9200/" + indexName) //Example: product_catalog_search_en_us

	if err != nil {
		fmt.Printf("Error %s", err.Error())
		c.JSON(400, err.Error())
	} else {
		c.Data(200, "application/json", response.Body())
	}
}

func HandlePostProducts(c *gin.Context) {
	indexName := c.Param("indexName")
	status := make(map[string]int)
	var products []model.Product
	if err := c.ShouldBind(&products); err != nil {
		c.JSON(400, err)
		return
	}
	for _, product := range products {
		if _, err := upsert(&product, indexName); err != nil {
			fmt.Println(err)
			status[product.BaseProduct] = 400
		} else {
			status[product.BaseProduct] = 200
		}
	}
	c.JSON(207, status)
}
