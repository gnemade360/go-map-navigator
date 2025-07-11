# Go Map Navigator

A flexible Go library for navigating, querying, and modifying nested map and slice data structures using simple string paths.

[![Go Reference](https://pkg.go.dev/badge/github.com/gnemade360/go-map-navigator.svg)](https://pkg.go.dev/github.com/gnemade360/go-map-navigator)
[![Go Report Card](https://goreportcard.com/badge/github.com/gnemade360/go-map-navigator)](https://goreportcard.com/report/github.com/gnemade360/go-map-navigator)

## Features

- **Simple Navigation**: Navigate nested data structures using string paths
- **Wildcard Operations**: Process all elements in maps or arrays using the `*` wildcard
- **Flexible Modifiers**: Apply custom transformations during navigation
- **Read-Only Mode**: Safely query data without modifications
- **Property Creation**: Dynamically create missing properties during navigation
- **Type Safety**: Comprehensive type checking and error handling
- **Extensible**: Easy to implement custom modifiers and conditions

## Installation

```bash
go get github.com/gnemade360/go-map-navigator
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/gnemade360/go-map-navigator/pkg/mapnavigator"
)

func main() {
    // Sample data structure
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

    // Create navigator
    navigator := mapnavigator.NewMapNavigator(nil)

    // Navigate to specific values
    result, err := navigator.VisitNode(data, "users", "0", "name")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("First user: %v\n", result) // Output: John Doe

    // Navigate to config
    result, err = navigator.VisitNode(data, "config", "port")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Port: %v\n", result) // Output: 8080
}
```

### Wildcard Operations

```go
// Get all user names
result, err := navigator.VisitNode(data, "users", "*", "name")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("All names: %v\n", result) // Output: [John Doe Jane Smith]
```

### Using Modifiers

```go
package main

import (
    "fmt"
    "strings"
    
    "github.com/gnemade360/go-map-navigator/pkg/mapnavigator"
    "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-nav-models"
)

func main() {
    data := map[string]interface{}{
        "message": "hello world",
    }

    // Create a modifier that converts strings to uppercase
    uppercaseModifier := models.MapNodeModifierFunc(func(node interface{}) interface{} {
        if str, ok := node.(string); ok {
            return strings.ToUpper(str)
        }
        return node
    })

    navigator := mapnavigator.NewMapNavigator(uppercaseModifier)

    // Apply modifier using "-" key
    result, err := navigator.VisitNode(data, "message", "-")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Modified: %v\n", result) // Output: HELLO WORLD
}
```

## Path Syntax

The library uses a simple string-based path syntax:

- `"key"` - Access map key
- `"0"`, `"1"`, `"2"` - Access array element by index
- `"*"` - Apply operation to all elements in map or array
- `"-"` - Apply modifier to current node

## Navigation Examples

### Maps
```go
// Navigate to nested map value
result, err := navigator.VisitNode(data, "config", "database", "host")
```

### Arrays
```go
// Navigate to array element
result, err := navigator.VisitNode(data, "users", "0", "name")
```

### Wildcards
```go
// Process all elements
result, err := navigator.VisitNode(data, "users", "*", "email")
```

### Modifiers
```go
// Apply modifier to current node
result, err := navigator.VisitNode(data, "users", "0", "name", "-")
```

## Configuration Options

### Read-Only Mode
```go
navigator := mapnavigator.NewMapNavigator(nil)
navigator.ReadOnly = true
```

### Property Creation
```go
navigator := mapnavigator.NewMapNavigator(nil)
navigator.CreateProperty = true
```

## Available Modifiers

The library includes several built-in modifiers:

- **Set Modifier**: Set values at specific paths
- **Delete Modifier**: Remove elements from data structures
- **Replace Modifier**: Replace values based on conditions
- **Conditional Modifier**: Apply modifiers based on conditions
- **Composite Modifier**: Chain multiple modifiers together
- **Expand Collection Modifier**: Flatten nested collections

## Custom Modifiers

You can create custom modifiers by implementing the `MapNodeModifier` interface:

```go
type CustomModifier struct{}

func (c CustomModifier) ModifyNode(node interface{}) interface{} {
    // Your custom logic here
    return node
}
```

Or use the function-based approach:

```go
modifier := models.MapNodeModifierFunc(func(node interface{}) interface{} {
    // Your custom logic here
    return node
})
```

## Error Handling

The library provides comprehensive error handling:

```go
result, err := navigator.VisitNode(data, "nonexistent", "key")
if err != nil {
    fmt.Printf("Navigation error: %v\n", err)
}
```

## Thread Safety

`MapNavigator` instances are not thread-safe. Create separate instances for concurrent operations or use appropriate synchronization mechanisms.

## Documentation

For complete API documentation, visit [pkg.go.dev](https://pkg.go.dev/github.com/gnemade360/go-map-navigator).

## Examples

More examples can be found in the `examples/` directory and in the test files throughout the codebase.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

If you encounter any issues or have questions, please open an issue on GitHub.