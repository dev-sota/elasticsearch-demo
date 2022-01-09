package main

import (
	"context"
	"fmt"
	"log"

	"github.com/olivere/elastic/v7"
)

const (
	url     = "http://127.0.0.1:9200"
	index   = "twitter"
	mapping = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"doc":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
                "retweets":{
                    "type":"long"
                },
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}
`
)

func main() {
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		log.Fatalln(err)
	}

	info, code, err := client.Ping(url).Do(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	exists, err := client.IndexExists(index).Do(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	if !exists {
		fmt.Println("Create a new Index.")
		res, err := client.CreateIndex(index).Body(mapping).IncludeTypeName(true).Do(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		if !res.Acknowledged {
			fmt.Println("Not acknowledged")
		}
	} else {
		fmt.Println("Exist")
	}
}
