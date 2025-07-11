package main

import (
	"encoding/json"
	"fmt"

	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator"
	models "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-nav-models"
)

func main() {
	// Example 1: Basic map navigation
	fmt.Println("=== Basic Map Navigation ===")
	basicNavigation()

	// Example 2: Array navigation
	fmt.Println("\n=== Array Navigation ===")
	arrayNavigation()

	// Example 3: Read-only navigation
	fmt.Println("\n=== Read-only Navigation ===")
	readOnlyNavigation()

	// Example 4: Property creation
	fmt.Println("\n=== Property Creation ===")
	propertyCreation()
}

func basicNavigation() {
	// Create a sample map
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"id":   123,
			"name": "John Doe",
			"profile": map[string]interface{}{
				"email":    "john@example.com",
				"location": "New York",
			},
		},
		"settings": map[string]interface{}{
			"theme": "dark",
			"lang":  "en",
		},
	}

	// Create a navigator with a simple modifier
	navigator := &mapnavigator.MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o // Return as-is for basic navigation
		}),
		ReadOnly: true,
	}

	// Navigate to nested properties
	if result, err := navigator.VisitNode(data, "user", "name"); err == nil {
		fmt.Printf("User name: %v\n", result)
	}

	if result, err := navigator.VisitNode(data, "user", "profile", "email"); err == nil {
		fmt.Printf("User email: %v\n", result)
	}

	if result, err := navigator.VisitNode(data, "settings", "theme"); err == nil {
		fmt.Printf("Theme: %v\n", result)
	}
}

func arrayNavigation() {
	// Create a sample map with arrays
	data := map[string]interface{}{
		"products": []interface{}{
			map[string]interface{}{
				"id":    1,
				"name":  "Laptop",
				"price": 999.99,
				"tags":  []interface{}{"electronics", "computer"},
			},
			map[string]interface{}{
				"id":    2,
				"name":  "Phone",
				"price": 599.99,
				"tags":  []interface{}{"electronics", "mobile"},
			},
			map[string]interface{}{
				"id":    3,
				"name":  "Book",
				"price": 19.99,
				"tags":  []interface{}{"education", "reading"},
			},
		},
	}

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly: true,
	}

	// Get all product names using wildcard
	if result, err := navigator.VisitNode(data, "products", "*", "name"); err == nil {
		fmt.Printf("All product names: %v\n", result)
	}

	// Get specific product by index
	if result, err := navigator.VisitNode(data, "products", "0", "name"); err == nil {
		fmt.Printf("First product name: %v\n", result)
	}

	if result, err := navigator.VisitNode(data, "products", "1", "price"); err == nil {
		fmt.Printf("Second product price: %v\n", result)
	}

	// Get all prices
	if result, err := navigator.VisitNode(data, "products", "*", "price"); err == nil {
		fmt.Printf("All product prices: %v\n", result)
	}
}

func readOnlyNavigation() {
	data := map[string]interface{}{
		"counter": 10,
		"status":  "active",
	}

	// Create a read-only navigator
	navigator := &mapnavigator.MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			if val, ok := o.(int); ok {
				return val + 1 // This modification won't be applied due to ReadOnly
			}
			return o
		}),
		ReadOnly: true,
	}

	// Navigate and attempt to modify
	if result, err := navigator.VisitNode(data, "counter"); err == nil {
		fmt.Printf("Counter value (read-only): %v\n", result)
	}

	// Original data remains unchanged
	fmt.Printf("Original counter: %v\n", data["counter"])
}

func propertyCreation() {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
		},
	}

	// Create a navigator with property creation enabled
	navigator := &mapnavigator.MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly:       false,
		CreateProperty: true,
	}

	// Try to access a non-existent property
	if result, err := navigator.VisitNode(data, "user", "email"); err == nil {
		fmt.Printf("Created property 'email': %v\n", result)
	}

	// Print the modified data
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Modified data:\n%s\n", jsonData)
	}
}
