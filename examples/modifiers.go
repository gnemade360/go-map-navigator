package main

import (
	"fmt"
	"log"

	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-set-modifier"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-delete-modifier"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-replace-modifier"
)

func main() {
	// Create a sample data structure
	data := map[interface{}]interface{}{
		"config": map[interface{}]interface{}{
			"database": map[interface{}]interface{}{
				"host":     "localhost",
				"port":     5432,
				"username": "admin",
				"password": "secret123",
			},
			"cache": map[interface{}]interface{}{
				"enabled": true,
				"ttl":     3600,
			},
		},
	}

	fmt.Println("Original data:")
	printData(data)

	// Example 1: Set Modifier - Update a value
	fmt.Println("\n--- Set Modifier Example ---")
	setModifier := &map_node_set_modifier.MapNodeSetModifier{
		Value: "db.example.com",
	}
	
	nav := &mapnavigator.MapNavigator{
		NodeModifier: setModifier,
		ReadOnly:     false,
	}
	
	_, err := nav.VisitNode(data, "config", "database", "host")
	if err != nil {
		log.Printf("Error setting value: %v\n", err)
	} else {
		fmt.Println("Updated database host")
		printData(data)
	}

	// Example 2: Delete Modifier - Remove a field
	fmt.Println("\n--- Delete Modifier Example ---")
	deleteModifier := &map_node_delete_modifier.MapNodeDeleteModifier{}
	
	nav.NodeModifier = deleteModifier
	_, err = nav.VisitNode(data, "config", "database", "password")
	if err != nil {
		log.Printf("Error deleting value: %v\n", err)
	} else {
		fmt.Println("Deleted password field")
		printData(data)
	}

	// Example 3: Replace Modifier - Replace values
	fmt.Println("\n--- Replace Modifier Example ---")
	replaceModifier := &map_node_replace_modifier.MapNodeReplaceModifier{
		OldValue: 5432,
		NewValue: 5433,
	}
	
	nav.NodeModifier = replaceModifier
	_, err = nav.VisitNode(data, "config", "database", "port")
	if err != nil {
		log.Printf("Error replacing value: %v\n", err)
	} else {
		fmt.Println("Replaced port value")
		printData(data)
	}

	// Example 4: Using CreateProperty to add new fields
	fmt.Println("\n--- Create Property Example ---")
	setModifier2 := &map_node_set_modifier.MapNodeSetModifier{
		Value: "production",
	}
	
	nav2 := &mapnavigator.MapNavigator{
		NodeModifier:   setModifier2,
		ReadOnly:       false,
		CreateProperty: true, // Enable property creation
	}
	
	_, err = nav2.VisitNode(data, "config", "environment")
	if err != nil {
		log.Printf("Error creating property: %v\n", err)
	} else {
		fmt.Println("Created new environment property")
		printData(data)
	}
}

func printData(data interface{}) {
	fmt.Printf("%+v\n", data)
}