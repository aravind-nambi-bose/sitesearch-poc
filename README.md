# Sitesearch PoC

## Setup

1. Start up the containers: (Open Distro ElasticSearch & Golang based Microservice)

`docker-compose up --build`

2. Initialize the ElasticSearch index:

`sh ./init/init.sh`

3. Query samples:

Get all -

http://localhost:8080/index/product_catalog_search_en_us

Query -

http://localhost:8080/index/product_catalog_search_en_us?q=wireless

Query with category filter -

http://localhost:8080/index/product_catalog_search_en_us?q=wireless&filter=Banded%20Wireless

Note:

Datamodel: [./model/product.go](./model/product.go)

Script score for relevancy ranking: [/service/product_search.go#L36](/service/product_search.go#L36)


