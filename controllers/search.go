package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/olivere/elastic/v7" // Correct import for the olivere elastic client
	"log"
	"strings" // Import strings package
	"context" // Import context package
)

type SearchController struct {
	web.Controller
}



// AutocompleteHandler handles autocomplete suggestions
func (c *SearchController) AutocompleteHandler() {
	queryParams := c.GetString("query")
	log.Println("Received autocomplete query:", queryParams)

	// Initialize Elasticsearch client
	esClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Println("Error creating Elasticsearch client:", err)
		c.Data["json"] = map[string]string{"error": "Failed to connect to Elasticsearch"}
		c.ServeJSON()
		return
	}

	// Your existing query
	query := map[string]interface{}{
		"size": 20, // Number of results to return
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match_phrase_prefix": map[string]interface{}{
							"products.product_name": map[string]interface{}{
								"query": queryParams,
								"max_expansions": 50,
								"boost": 4,
							},
						},
					},
					{
						"match": map[string]interface{}{
							"products.product_name": map[string]interface{}{
								"query": queryParams,
								"operator": "or",
								"fuzziness": "AUTO",
								"prefix_length": 1,
								"boost": 2,
							},
						},
					},
					{
						"wildcard": map[string]interface{}{
							"products.product_name": map[string]interface{}{
								"value": fmt.Sprintf("*%s*", strings.ToLower(queryParams)),
								"boost": 1,
							},
						},
					},
				},
				"minimum_should_match": 1,
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"products.product_name": map[string]interface{}{},
			},
		},
	}

	// Convert query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		log.Println("Error marshaling query JSON:", err)
		c.Data["json"] = map[string]string{"error": "Failed to build query"}
		c.ServeJSON()
		return
	}

	log.Println("Constructed Elasticsearch Query:", string(queryJSON))

	// Execute search request
	searchResult, err := esClient.Search().
	Index("kibana_sample_data_ecommerce").
	Source(string(queryJSON)). // ✅ Use Source() with marshaled JSON
	Do(context.Background())


	if err != nil {
		log.Println("Error executing search query:", err)
		c.Data["json"] = map[string]string{"error": "Autocomplete request failed"}
		c.ServeJSON()
		return
	}

	// Process the search results
	var suggestions []string
	for _, hit := range searchResult.Hits.Hits {
		var hitSource map[string]interface{}
		if err := json.Unmarshal(hit.Source, &hitSource); err != nil {
			log.Println("Error unmarshaling hit source:", err)
			continue
		}

		// Extract product names from the products array
		if products, exists := hitSource["products"].([]interface{}); exists {
			for _, productItem := range products {
				if productInfo, ok := productItem.(map[string]interface{}); ok {
					if name, found := productInfo["product_name"].(string); found {
						suggestions = append(suggestions, name)
					}
				}
			}
		}
	}

	// Log and return the autocomplete suggestions
	log.Println("Returning autocomplete suggestions:", suggestions)
	c.Data["json"] = suggestions
	c.ServeJSON()
}






func (c *SearchController) SearchHandler() {
	query := c.GetString("query")
	log.Println("Received search query:", query)

	// Initialize Elasticsearch client
	esClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Println("Error creating Elasticsearch client:", err)
		c.Data["json"] = map[string]string{"error": "Failed to connect to Elasticsearch"}
		c.ServeJSON()
		return
	}

	// Construct search query using a simple match
	searchQuery := fmt.Sprintf(`{
		"query": {
			"match": {
				"products.product_name": {
					"query": "%s",
					"fuzziness": "AUTO"
				}
			}
		},
		"size": 5
	}`, query)
	log.Println("Constructed search query:", searchQuery)

	// Execute search request
	searchResult, err := esClient.Search().
		Index("kibana_sample_data_ecommerce").
		Source(searchQuery). // ✅ Use Source() directly
		Pretty(true).
		Do(context.Background())

	if err != nil {
		log.Println("Error executing search query:", err)
		c.Data["json"] = map[string]string{"error": "Search request failed"}
		c.ServeJSON()
		return
	}

	// Process search results
	for _, hit := range searchResult.Hits.Hits {
		var productData map[string]interface{}
		if err := json.Unmarshal(hit.Source, &productData); err != nil {
			log.Println("Error unmarshalling product details:", err)
			continue
		}

		// Extract product details from the array
		if products, ok := productData["products"].([]interface{}); ok {
			for _, productItem := range products {
				if productInfo, ok := productItem.(map[string]interface{}); ok {
					// Check if product name matches
					if name, exists := productInfo["product_name"].(string); exists && strings.EqualFold(name, query) {
						// Extract relevant details
						product := map[string]string{
							"name":         name,
							"description":  "No description available",
							"price":        "N/A",
							"manufacturer": "Unknown",
							"base_price":   "N/A",
						}

						// Extract base price
						if price, exists := productInfo["base_price"].(float64); exists {
							product["base_price"] = fmt.Sprintf("$%.2f", price)
						}

						// Extract manufacturer
						if manufacturer, exists := productInfo["manufacturer"].(string); exists {
							product["manufacturer"] = manufacturer
						}

						log.Println("Returning product details:", product)
						c.Data["json"] = product
						c.ServeJSON()
						return
					}
				}
			}
		}
	}

	// Return empty if no matching product found
	log.Println("No matching product found.")
	c.Data["json"] = map[string]string{"error": "No matching product found"}
	c.ServeJSON()
}



