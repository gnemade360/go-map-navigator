package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator"
	models "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-nav-models"
)

func main() {
	// Example 1: JSON data extraction
	fmt.Println("=== JSON Data Extraction ===")
	jsonExtractionExample()

	// Example 2: JSON data transformation
	fmt.Println("\n=== JSON Data Transformation ===")
	jsonTransformationExample()

	// Example 3: JSON validation and cleaning
	fmt.Println("\n=== JSON Validation and Cleaning ===")
	jsonValidationExample()

	// Example 4: Complex JSON navigation
	fmt.Println("\n=== Complex JSON Navigation ===")
	complexNavigationExample()
}

func jsonExtractionExample() {
	// Sample JSON data representing an API response
	jsonData := `{
		"data": {
			"user": {
				"id": 12345,
				"profile": {
					"name": "John Doe",
					"email": "john@example.com",
					"preferences": {
						"language": "en",
						"timezone": "UTC"
					}
				},
				"orders": [
					{
						"id": "order-1",
						"total": 99.99,
						"status": "completed",
						"items": [
							{"name": "Product A", "quantity": 2},
							{"name": "Product B", "quantity": 1}
						]
					},
					{
						"id": "order-2",
						"total": 149.99,
						"status": "pending",
						"items": [
							{"name": "Product C", "quantity": 1}
						]
					}
				]
			}
		}
	}`

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly: true,
	}

	// Extract user information
	if result, err := navigator.VisitNode(data, "data", "user", "profile", "name"); err == nil {
		fmt.Printf("User name: %v\n", result)
	}

	if result, err := navigator.VisitNode(data, "data", "user", "profile", "email"); err == nil {
		fmt.Printf("User email: %v\n", result)
	}

	// Extract all order totals
	if result, err := navigator.VisitNode(data, "data", "user", "orders", "*", "total"); err == nil {
		fmt.Printf("Order totals: %v\n", result)
	}

	// Extract all item names from all orders
	if result, err := navigator.VisitNode(data, "data", "user", "orders", "*", "items", "*", "name"); err == nil {
		fmt.Printf("All item names: %v\n", result)
	}

	// Extract specific order status
	if result, err := navigator.VisitNode(data, "data", "user", "orders", "1", "status"); err == nil {
		fmt.Printf("Second order status: %v\n", result)
	}
}

func jsonTransformationExample() {
	// Sample product catalog JSON
	jsonData := `{
		"catalog": {
			"products": [
				{
					"id": "prod-1",
					"name": "Laptop",
					"price": 999.99,
					"category": "electronics",
					"available": true
				},
				{
					"id": "prod-2",
					"name": "Phone",
					"price": 599.99,
					"category": "electronics",
					"available": false
				},
				{
					"id": "prod-3",
					"name": "Book",
					"price": 19.99,
					"category": "books",
					"available": true
				}
			]
		}
	}`

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	// Transform prices to include tax (10%)
	taxModifier := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if price, ok := o.(float64); ok {
			return price * 1.10 // Add 10% tax
		}
		return o
	})

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: taxModifier,
		ReadOnly:     false,
	}

	// Apply tax to all prices
	if result, err := navigator.VisitNode(data, "catalog", "products", "*", "price"); err == nil {
		fmt.Printf("Prices with tax: %v\n", result)
	}

	// Print the transformed data
	if jsonOutput, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Transformed catalog:\n%s\n", jsonOutput)
	}
}

func jsonValidationExample() {
	// Sample JSON with some invalid/missing data
	jsonData := `{
		"users": [
			{
				"id": 1,
				"name": "John Doe",
				"email": "john@example.com",
				"age": 30
			},
			{
				"id": 2,
				"name": "",
				"email": "invalid-email",
				"age": -5
			},
			{
				"id": 3,
				"name": "Jane Smith",
				"email": "jane@example.com"
			}
		]
	}`

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	// Validation and cleaning modifier
	cleanModifier := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		switch v := o.(type) {
		case string:
			if v == "" {
				return "N/A" // Replace empty strings
			}
			return v
		case float64:
			if v < 0 {
				return 0 // Replace negative numbers
			}
			return v
		default:
			return o
		}
	})

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: cleanModifier,
		ReadOnly:     false,
	}

	// Clean all names
	if result, err := navigator.VisitNode(data, "users", "*", "name"); err == nil {
		fmt.Printf("Cleaned names: %v\n", result)
	}

	// Clean all ages
	if result, err := navigator.VisitNode(data, "users", "*", "age"); err == nil {
		fmt.Printf("Cleaned ages: %v\n", result)
	}

	// Print the cleaned data
	if jsonOutput, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Cleaned data:\n%s\n", jsonOutput)
	}
}

func complexNavigationExample() {
	// Complex nested JSON structure
	jsonData := `{
		"company": {
			"name": "Tech Corp",
			"departments": [
				{
					"name": "Engineering",
					"teams": [
						{
							"name": "Backend",
							"members": [
								{"name": "Alice", "role": "Senior", "salary": 80000},
								{"name": "Bob", "role": "Junior", "salary": 50000}
							]
						},
						{
							"name": "Frontend",
							"members": [
								{"name": "Charlie", "role": "Senior", "salary": 75000},
								{"name": "Diana", "role": "Mid", "salary": 60000}
							]
						}
					]
				},
				{
					"name": "Sales",
					"teams": [
						{
							"name": "Enterprise",
							"members": [
								{"name": "Eve", "role": "Senior", "salary": 70000}
							]
						}
					]
				}
			]
		}
	}`

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly: true,
	}

	// Get all employee names across all departments and teams
	if result, err := navigator.VisitNode(data, "company", "departments", "*", "teams", "*", "members", "*", "name"); err == nil {
		fmt.Printf("All employee names: %v\n", result)
	}

	// Get all salaries for Senior roles
	if result, err := navigator.VisitNode(data, "company", "departments", "*", "teams", "*", "members", "*", "salary"); err == nil {
		fmt.Printf("All salaries: %v\n", result)
	}

	// Get specific team information
	if result, err := navigator.VisitNode(data, "company", "departments", "0", "teams", "0", "name"); err == nil {
		fmt.Printf("First team name: %v\n", result)
	}

	// Get Engineering department's first team's members
	if result, err := navigator.VisitNode(data, "company", "departments", "0", "teams", "0", "members", "*", "name"); err == nil {
		fmt.Printf("Backend team members: %v\n", result)
	}

	// Calculate total salary for a specific team
	totalSalaryModifier := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if salaries, ok := o.([]interface{}); ok {
			total := 0.0
			for _, salary := range salaries {
				if s, ok := salary.(float64); ok {
					total += s
				}
			}
			return total
		}
		return o
	})

	salaryNavigator := &mapnavigator.MapNavigator{
		NodeModifier: totalSalaryModifier,
		ReadOnly:     false,
	}

	// Get total salary for Backend team
	if result, err := salaryNavigator.VisitNode(data, "company", "departments", "0", "teams", "0", "members", "*", "salary"); err == nil {
		fmt.Printf("Backend team total salary: %v\n", result)
	}
}