package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/olivere/elastic/v7"
)

const (
	url       = "http://127.0.0.1:9200"
	indexName = "twitter"
	mapping   = `
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

type Tweet struct {
	User     string
	Message  string
	Retweets int
}

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

	exists, err := client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	if !exists {
		fmt.Println("Create a new Index.")
		res, err := client.CreateIndex(indexName).Body(mapping).IncludeTypeName(true).Do(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		if !res.Acknowledged {
			fmt.Println("Not acknowledged")
		}
	} else {
		fmt.Println("Exist")
	}

	id := "1"
	t := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
	put, err := client.Index().
		Index(indexName).
		Id(id).
		BodyJson(t).
		Do(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put.Id, put.Index, put.Type)

	get, err := client.Get().
		Index(indexName).
		Id(id).
		Do(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Got document %s in version %d from index %s, type %s\n", get.Id, get.Version, get.Index, get.Type)

	// Refresh before search
	_, err = client.Refresh().Index(indexName).Do(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	res, err := client.Search().
		Index(indexName).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Query took %d milliseconds\n", res.TookInMillis)
	fmt.Printf("Found a total of %d tweets\n", res.TotalHits())

	var ttyp Tweet
	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		t := item.(Tweet)
		fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
	}

	if res.TotalHits() > 0 {
		for _, hit := range res.Hits.Hits {
			var t Tweet
			err := json.Unmarshal(hit.Source, &t)
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		fmt.Print("Found no tweets\n")
	}
}
