# go-map-navigator

🗺️ A lightweight Go library for navigating nested data structures with dot notation, wildcards, and array indexing

## 🚀 Features

- **Dot Notation**: Access nested values with `"user.address.city"`
- **Wildcards**: Get all matching values with `"users.*.email"`
- **Array Indexing**: Access array elements with `"items[0].name"`
- **Type Safety**: Automatic type conversions with safety checks
- **Zero Dependencies**: Only uses Go standard library
- **High Performance**: Optimized for speed with minimal allocations
- **Create Mode**: Optionally create missing paths
- **Read/Write**: Support for both reading and modifying values
- **Node Modifiers**: Transform values during navigation
- **Template Support**: Dynamic value resolution with Go templates

## 📦 Installation

```bash
go get github.com/passionintellectual/go-map-navigator
```

## 🔧 Quick Start

```go
package main

import (
    "fmt"
    "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator"
)

func main() {
    // Sample nested data
    data := map[string]interface{}{
        "users": []interface{}{
            map[string]interface{}{
                "name": "John Doe",
                "age": 30,
                "address": map[string]interface{}{
                    "city": "New York",
                    "country": "USA",
                },
            },
            map[string]interface{}{
                "name": "Jane Smith",
                "age": 25,
                "address": map[string]interface{}{
                    "city": "London",
                    "country": "UK",
                },
            },
        },
        "settings": map[string]interface{}{
            "theme": "dark",
            "notifications": true,
        },
    }

    // Create a navigator
    nav := &mapnavigator.MapNavigator{
        ReadOnly: true, // Set to false to allow modifications
    }

    // Simple navigation
    city, err := nav.VisitMapStringNode(data, "users", "0", "address", "city")
    if err == nil {
        fmt.Println("City:", city) // Output: City: New York
    }

    // Using wildcards
    names, err := nav.VisitMapStringNode(data, "users", "*", "name")
    if err == nil {
        fmt.Println("Names:", names) // Output: Names: [John Doe Jane Smith]
    }
}
```

## 📖 Advanced Usage

### Node Modifiers

Transform values during navigation:

```go
import (
    "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-set-modifier"
)

// Create a modifier to set a value
setModifier := &map_node_set_modifier.MapNodeSetModifier{
    Value: "San Francisco",
}

nav := &mapnavigator.MapNavigator{
    NodeModifier: setModifier,
    ReadOnly:     false,
}

// This will set the city to "San Francisco"
nav.VisitMapStringNode(data, "users", "0", "address", "city")
```

### Creating Missing Paths

```go
nav := &mapnavigator.MapNavigator{
    CreateProperty: true,
    ReadOnly:       false,
}

// This will create the path if it doesn't exist
nav.VisitMapStringNode(data, "users", "0", "profile", "bio")
```

## 🛠️ API Reference

### MapNavigator

The main struct for navigating maps:

```go
type MapNavigator struct {
    NodeModifier   models.MapNodeModifier // Optional modifier for transforming values
    ReadOnly       bool                   // If true, prevents modifications
    CreateProperty bool                   // If true, creates missing paths
}
```

### Methods

- `VisitMapNode(mp map[interface{}]interface{}, ks ...string) (interface{}, error)` - Navigate a map with interface{} keys
- `VisitMapStringNode(mp map[string]interface{}, ks ...string) (interface{}, error)` - Navigate a map with string keys
- `VisitNode(node interface{}, ks ...string) (interface{}, error)` - Navigate any node (map, slice, or value)

### Special Keys

- `*` - Wildcard: matches all keys at the current level
- `-` - Apply node modifier to current value
- `[n]` - Array index notation (when used in path strings)

## 🧩 Node Modifiers

Available modifiers in the `map-node-modifiers` package:

- **SetModifier**: Set a value at the target location
- **DeleteModifier**: Delete a value at the target location
- **ReplaceModifier**: Replace values matching a pattern
- **ConditionalModifier**: Apply modifications based on conditions
- **CompositeModifier**: Combine multiple modifiers
- **ExpandCollectionModifier**: Expand collections with templates

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Originally extracted from `go-common-lib` to provide a standalone, reusable map navigation library for the Go community.