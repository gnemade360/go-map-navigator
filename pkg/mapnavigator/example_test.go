package mapnavigator

import (
	"fmt"
	"log"

	models "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-nav-models"
)

// ExampleMapNavigator demonstrates basic usage of MapNavigator for navigating nested data structures.
func ExampleMapNavigator() {
	// Create sample data structure
	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{
				"name":  "John Doe",
				"email": "john@example.com",
				"age":   30,
			},
			map[string]interface{}{
				"name":  "Jane Smith",
				"email": "jane@example.com",
				"age":   25,
			},
		},
		"config": map[string]interface{}{
			"debug": true,
			"port":  8080,
		},
	}

	// Create a navigator without any modifiers
	navigator := NewMapNavigator(nil)

	// Navigate to a specific user's name
	result, err := navigator.VisitNode(data, "users", "0", "name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("First user name: %v\n", result)

	// Navigate to config value
	result, err = navigator.VisitNode(data, "config", "port")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Port: %v\n", result)

	// Output:
	// First user name: John Doe
	// Port: 8080
}

// ExampleMapNavigator_wildcard demonstrates using wildcard (*) to process all elements.
func ExampleMapNavigator_wildcard() {
	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "John", "active": true},
			map[string]interface{}{"name": "Jane", "active": false},
		},
	}

	navigator := NewMapNavigator(nil)

	// Get all user names using wildcard
	result, err := navigator.VisitNode(data, "users", "*", "name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("All names: %v\n", result)

	// Output:
	// All names: [John Jane]
}

// ExampleMapNavigator_withModifier demonstrates using a navigator with a custom modifier.
func ExampleMapNavigator_withModifier() {
	data := map[string]interface{}{
		"message": "hello world",
	}

	// Create a modifier that converts strings to uppercase
	uppercaseModifier := models.MapNodeModifierFunc(func(node interface{}) interface{} {
		if str, ok := node.(string); ok {
			return fmt.Sprintf("MODIFIED: %s", str)
		}
		return node
	})

	navigator := NewMapNavigator(uppercaseModifier)

	// Navigate to the message and let the navigator apply the modifier
	result, err := navigator.VisitNode(data, "message")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Modified message: %v\n", result)

	// Output:
	// Modified message: MODIFIED: hello world
}

// ExampleMapNavigator_readOnly demonstrates read-only navigation.
func ExampleMapNavigator_readOnly() {
	data := map[string]interface{}{
		"config": map[string]interface{}{
			"debug": false,
		},
	}

	// Create a modifier that would change the value
	toggleModifier := models.MapNodeModifierFunc(func(node interface{}) interface{} {
		if b, ok := node.(bool); ok {
			return !b
		}
		return node
	})

	navigator := NewMapNavigator(toggleModifier)
	navigator.ReadOnly = true

	// Navigate to the value and apply modifier
	result, err := navigator.VisitNode(data, "config", "debug")
	if err != nil {
		log.Fatal(err)
	}

	// Apply modifier manually to demonstrate read-only behavior
	if navigator.NodeModifier != nil {
		modifiedResult := navigator.NodeModifier.ModifyNode(result)
		fmt.Printf("Result: %v\n", modifiedResult)
	}
	fmt.Printf("Original unchanged: %v\n", data["config"].(map[string]interface{})["debug"])

	// Output:
	// Result: true
	// Original unchanged: false
}

// ExampleNewMapNavigatorFunc demonstrates creating a navigator with a function-based modifier.
func ExampleNewMapNavigatorFunc() {
	data := map[string]interface{}{
		"numbers": []interface{}{1, 2, 3, 4, 5},
	}

	// Create navigator with a function that doubles numbers
	navigator := NewMapNavigatorFunc(func(node interface{}) interface{} {
		if num, ok := node.(int); ok {
			return num * 2
		}
		return node
	})

	// Apply modifier to all numbers using wildcard
	result, err := navigator.VisitNode(data, "numbers", "*")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Doubled numbers: %v\n", result)

	// Output:
	// Doubled numbers: [2 4 6 8 10]
}

// ExampleMapNavigator_createProperty demonstrates dynamic property creation.
func ExampleMapNavigator_createProperty() {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
		},
	}

	navigator := NewMapNavigator(nil)
	navigator.CreateProperty = true

	// Try to access a non-existent property
	result, err := navigator.VisitNode(data, "user", "email")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created property:%v\n", result)
	fmt.Printf("Data now contains: %v\n", data["user"])

	// Output:
	// Created property:
	// Data now contains: map[email: name:John]
}
