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
	query := c.GetString("query")
	log.Println("Received autocomplete query:", query)

	// Initialize Elasticsearch client
	esClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Println("Error creating Elasticsearch client:", err)
		c.Data["json"] = map[string]string{"error": "Failed to connect to Elasticsearch"}
		c.ServeJSON()
		return
	}

	// Construct the search query for autocomplete
	searchQuery := fmt.Sprintf(`{
		"size": 10,
		"query": {
			"match_phrase_prefix": {
				"products.product_name": "%s"
			}
		}
	}`, query)
	

	// Use Source() to pass the query instead of BodyString()
	searchResult, err := esClient.Search().
		Index("kibana_sample_data_ecommerce").
		Source(strings.NewReader(searchQuery)).
		Pretty(true).
		Do(context.Background()) // Use context.Background() instead of c.Ctx

	if err != nil {
		log.Println("Error executing autocomplete query:", err)
		c.Data["json"] = map[string]string{"error": "Autocomplete request failed"}
		c.ServeJSON()
		return
	}

	// Log the raw response for debugging
	rawBody, err := json.Marshal(searchResult)
	if err != nil {
		log.Println("Error marshalling response body:", err)
		c.Data["json"] = map[string]string{"error": "Failed to process response"}
		c.ServeJSON()
		return
	}
	log.Println("Raw Elasticsearch response:", string(rawBody))

	// Parse the JSON response
	var response map[string]interface{}
	if err := json.Unmarshal(rawBody, &response); err != nil {
		log.Println("Error parsing response body:", err)
		c.Data["json"] = map[string]string{"error": "Invalid response format"}
		c.ServeJSON()
		return
	}

	// Extract hits safely
	hitsData, ok := response["hits"].(map[string]interface{})
	if !ok {
		log.Println("Error: 'hits' field missing or invalid")
		c.Data["json"] = []string{}
		c.ServeJSON()
		return
	}

	hits, ok := hitsData["hits"].([]interface{})
	if !ok {
		log.Println("Error: 'hits.hits' field missing or invalid")
		c.Data["json"] = []string{}
		c.ServeJSON()
		return
	}

	// Log the number of hits to verify we're getting results
	log.Printf("Found %d hits\n", len(hits))

	// Parse product names from the hits
	var suggestions []string
	for _, hit := range hits {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		// Log the full source to verify the structure
		log.Println("Source:", source)

		// Check if products is an array and extract product_name
		products, ok := source["products"].([]interface{})
		if !ok {
			continue
		}

		for _, product := range products {
			if productMap, ok := product.(map[string]interface{}); ok {
				if name, ok := productMap["product_name"].(string); ok {
					suggestions = append(suggestions, name)
				}
			}
		}
	}

	// Log suggestions for debugging
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
		Source(strings.NewReader(searchQuery)).
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



