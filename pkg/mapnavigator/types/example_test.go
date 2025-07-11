package types

import (
	"fmt"
)

// ExampleNewCollectionMapKey demonstrates creating a new CollectionMapKey.
func ExampleNewCollectionMapKey() {
	// Create a collection map key
	key := NewCollectionMapKey("users", 0)
	
	fmt.Printf("Key: %s, Index: %d\n", key.Key, key.Index)
	
	// Output:
	// Key: users, Index: 0
}

// ExampleConvertToMapEntries demonstrates converting a map to MapEntry slice.
func ExampleConvertToMapEntries() {
	// Create a sample map
	data := map[interface{}]interface{}{
		"name":  "John Doe",
		"age":   30,
		"email": "john@example.com",
	}
	
	// Convert to entries
	entries := ConvertToMapEntries(data)
	
	fmt.Printf("Number of entries: %d\n", len(entries))
	for _, entry := range entries {
		fmt.Printf("Key: %v, Value: %v\n", entry.Key, entry.Value)
	}
	
	// Output:
	// Number of entries: 3
	// Key: name, Value: John Doe
	// Key: age, Value: 30
	// Key: email, Value: john@example.com
}

// ExampleConvertFromMapEntries demonstrates converting MapEntry slice back to map.
func ExampleConvertFromMapEntries() {
	// Create entries
	entries := []MapEntry{
		{Key: "name", Value: "Jane Smith"},
		{Key: "age", Value: 25},
		{Key: "active", Value: true},
	}
	
	// Convert back to map
	data := ConvertFromMapEntries(entries)
	
	fmt.Printf("Converted map: %v\n", data)
	
	// Output:
	// Converted map: map[active:true age:25 name:Jane Smith]
}

// ExampleMapEntry demonstrates working with MapEntry directly.
func ExampleMapEntry() {
	// Create a MapEntry
	entry := MapEntry{
		Key:   "config",
		Value: map[string]interface{}{"debug": true, "port": 8080},
	}
	
	fmt.Printf("Entry key: %v\n", entry.Key)
	fmt.Printf("Entry value: %v\n", entry.Value)
	
	// Output:
	// Entry key: config
	// Entry value: map[debug:true port:8080]
}