package main

import (
	"fmt"
	"log"

	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator"
)

func main() {
	// Create a sample nested map structure
	data := map[interface{}]interface{}{
		"company": map[interface{}]interface{}{
			"name": "Tech Corp",
			"employees": map[interface{}]interface{}{
				"john": map[interface{}]interface{}{
					"age":        30,
					"department": "Engineering",
					"skills":     []string{"Go", "Python", "JavaScript"},
				},
				"jane": map[interface{}]interface{}{
					"age":        28,
					"department": "Marketing",
					"skills":     []string{"SEO", "Content Marketing", "Analytics"},
				},
			},
			"departments": []string{"Engineering", "Marketing", "Sales"},
		},
	}

	// Create a navigator instance
	nav := &mapnavigator.MapNavigator{
		ReadOnly: true, // Set to true for read-only operations
	}

	// Example 1: Navigate to a specific value
	fmt.Println("Example 1: Basic Navigation")
	value, err := nav.VisitNode(data, "company", "name")
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Company name: %v\n", value)
	}

	// Example 2: Navigate to nested values
	fmt.Println("\nExample 2: Nested Navigation")
	johnAge, err := nav.VisitNode(data, "company", "employees", "john", "age")
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("John's age: %v\n", johnAge)
	}

	// Example 3: Access array elements
	fmt.Println("\nExample 3: Array Access")
	firstDept, err := nav.VisitNode(data, "company", "departments", "0")
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("First department: %v\n", firstDept)
	}

	// Example 4: Use wildcards to get all employee names
	fmt.Println("\nExample 4: Wildcard Navigation")
	allEmployees, err := nav.VisitNode(data, "company", "employees", "*", "department")
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("All employee departments: %v\n", allEmployees)
	}
}