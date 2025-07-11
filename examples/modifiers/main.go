package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator"
	models "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-replace-modifier"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-set-modifier"
)

func main() {
	// Example 1: Using Replace Modifier
	fmt.Println("=== Replace Modifier Example ===")
	replaceModifierExample()

	// Example 2: Using Set Modifier
	fmt.Println("\n=== Set Modifier Example ===")
	setModifierExample()

	// Example 3: Custom Modifier
	fmt.Println("\n=== Custom Modifier Example ===")
	customModifierExample()

	// Example 4: Chaining Modifiers
	fmt.Println("\n=== Chaining Modifiers Example ===")
	chainingModifiersExample()
}

func replaceModifierExample() {
	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{
				"id":     1,
				"name":   "John Doe",
				"status": "inactive",
			},
			map[string]interface{}{
				"id":     2,
				"name":   "Jane Smith",
				"status": "active",
			},
		},
	}

	// Create a replace modifier to change "inactive" to "active"
	replaceModifier := map_node_replace_modifier.NewMapNodeReplaceModifier(
		"inactive", "active", false,
	)

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: replaceModifier,
		ReadOnly:     false,
	}

	// Apply the modifier to all user statuses
	if result, err := navigator.VisitNode(data, "users", "*", "status"); err == nil {
		fmt.Printf("Modified statuses: %v\n", result)
	}

	// Print the modified data
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Modified data:\n%s\n", jsonData)
	}
}

func setModifierExample() {
	data := map[string]interface{}{
		"products": []interface{}{
			map[string]interface{}{
				"id":    1,
				"name":  "Laptop",
				"price": 999.99,
			},
			map[string]interface{}{
				"id":    2,
				"name":  "Phone",
				"price": 599.99,
			},
		},
	}

	// Create a set modifier to add a discount flag
	setModifier := map_node_set_modifier.NewMapNodeSetModifier(
		"on_sale", true,
	)

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: setModifier,
		ReadOnly:     false,
	}

	// Apply the modifier to all products
	if result, err := navigator.VisitNode(data, "products", "*", "-"); err == nil {
		fmt.Printf("Modified products: %v\n", len(result.([]interface{})))
	}

	// Print the modified data
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Modified data:\n%s\n", jsonData)
	}
}

func customModifierExample() {
	data := map[string]interface{}{
		"employees": []interface{}{
			map[string]interface{}{
				"name":   "Alice",
				"salary": 50000,
				"role":   "developer",
			},
			map[string]interface{}{
				"name":   "Bob",
				"salary": 60000,
				"role":   "designer",
			},
		},
	}

	// Create a custom modifier that applies salary bonus
	bonusModifier := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if salary, ok := o.(int); ok {
			return salary + 5000 // Add $5000 bonus
		}
		if salary, ok := o.(float64); ok {
			return salary + 5000.0 // Add $5000 bonus
		}
		return o
	})

	navigator := &mapnavigator.MapNavigator{
		NodeModifier: bonusModifier,
		ReadOnly:     false,
	}

	// Apply bonus to all salaries
	if result, err := navigator.VisitNode(data, "employees", "*", "salary"); err == nil {
		fmt.Printf("Updated salaries: %v\n", result)
	}

	// Print the modified data
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Modified data:\n%s\n", jsonData)
	}
}

func chainingModifiersExample() {
	data := map[string]interface{}{
		"config": map[string]interface{}{
			"database": map[string]interface{}{
				"host":     "localhost",
				"port":     5432,
				"username": "admin",
				"password": "secret",
			},
		},
	}

	// First, mask sensitive data
	maskModifier := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if str, ok := o.(string); ok && str == "secret" {
			return "***masked***"
		}
		return o
	})

	navigator1 := &mapnavigator.MapNavigator{
		NodeModifier: maskModifier,
		ReadOnly:     false,
	}

	// Apply masking to password
	if _, err := navigator1.VisitNode(data, "config", "database", "password"); err != nil {
		log.Printf("Error masking password: %v", err)
	}

	// Then, add a timestamp
	timestampModifier := models.MapNodeModifierFunc(func(o interface{}) interface{} {
		if dbConfig, ok := o.(map[string]interface{}); ok {
			dbConfig["last_updated"] = "2023-01-01T00:00:00Z"
			return dbConfig
		}
		return o
	})

	navigator2 := &mapnavigator.MapNavigator{
		NodeModifier: timestampModifier,
		ReadOnly:     false,
	}

	// Apply timestamp to database config
	if _, err := navigator2.VisitNode(data, "config", "database", "-"); err != nil {
		log.Printf("Error adding timestamp: %v", err)
	}

	// Print the final modified data
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("Final modified data:\n%s\n", jsonData)
	}
}