package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web" // ✅ Correct Beego v2 import
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"strings"
)

type SearchController struct {
	web.Controller // ✅ Use Beego v2's `web.Controller`
}

// SearchHandler handles search requests
func (c *SearchController) SearchHandler() {
	query := c.GetString("query") // Get search query from the request

	// Initialize the Elasticsearch client (Use the already running Elasticsearch instance)
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", // Your Elasticsearch instance URL
		},
	})
	if err != nil {
		log.Println("Error creating the client:", err)
		c.Data["json"] = map[string]string{"error": "Failed to connect to Elasticsearch"}
		c.ServeJSON()
		return
	}

	// Perform search using Elasticsearch's _search API
searchResult, err := esClient.Search(
    esClient.Search.WithIndex("kibana_sample_data_ecommerce"), // Use Kibana's sample data index
    esClient.Search.WithBody(strings.NewReader(`{
        "query": {
            "match": {
                "category.keyword": {
                    "query": "` + query + `",
                    "fuzziness": "AUTO"
                }
            }
        }
    }`)),
    esClient.Search.WithPretty(),
)

	if err != nil {
		log.Println("Error getting response:", err)
		c.Data["json"] = map[string]string{"error": "Search request failed"}
		c.ServeJSON()
		return
	}
	defer searchResult.Body.Close()

// Parse the result
var response map[string]interface{}
if err := json.NewDecoder(searchResult.Body).Decode(&response); err != nil {
    log.Fatalf("Error parsing the response body: %s", err)
}

// Print response for debugging
log.Printf("Elasticsearch response: %+v\n", response)

	// Extract hits safely
	hitsData, ok := response["hits"].(map[string]interface{})
	if !ok {
		log.Println("Error: hits field missing or invalid")
		c.Data["json"] = []string{}
		c.ServeJSON()
		return
	}

	hits, ok := hitsData["hits"].([]interface{})
	if !ok {
		log.Println("Error: hits.hits field missing or invalid")
		c.Data["json"] = []string{}
		c.ServeJSON()
		return
	}

	// Return top 20 product names as JSON
	var products []string
	for _, hit := range hits {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		name, ok := source["name"].(string)
		if !ok {
			continue
		}

		products = append(products, name)
	}

	// Send the result back to the frontend
	c.Data["json"] = products
	c.ServeJSON()
}
