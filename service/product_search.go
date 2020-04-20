package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"time"
)

func HandleSearchProducts(c *gin.Context) {
	query := c.Query("q")
	filter := c.Query("filter")
	var filterStr string
	if filter != "" {
		filterStr = `"filter": {"term": { "productCategories.raw": "` + filter + `"}},`
	}
	indexName := c.Param("indexName")
	client := resty.New()
	response, err := client.R().
		SetBody([]byte(`
			  {
			  "query": {
			    "script_score": {
			      "query": {
					"bool": {
						`+filterStr+`
						"must": {
			        		"multi_match": {
			          			"query": "`+query+`",
					  			"fields": ["productName^100", "*"]	
			        		}
						}
					}	
			      },
			      "script": {
			        "source": "Math.max(_score, Math.max(decayDateGauss(params.date_origin, params.date_scale, params.date_offset, params.date_decay, doc['launchDate'].value) * 90, doc['discountPercent'].value * 80))",
			        "params": {
			          "date_origin": "`+time.Now().Format("2006-01-02T15:04:05Z")+`",
			          "date_scale": "365d",
			          "date_offset": "0",
			          "date_decay": 0.5
			        }
			      }
			    }
			  },
			  "aggs": {
			    "productCategories": {
			      "terms": {
			        "field": "productCategories.raw"
			      }
			    }
			  },
			  "explain": false
			}
		`)).
		SetHeader("Content-Type", "application/json").
		Post("http://elasticsearch:9200/" + indexName + "/_search")

	if err != nil {
		fmt.Printf("Error %s", err.Error())
		c.JSON(400, err.Error())
	} else {
		c.Data(200, "application/json", response.Body())
	}
}
