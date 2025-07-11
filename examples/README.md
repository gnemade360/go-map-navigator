# Go Map Navigator Examples

This directory contains comprehensive examples demonstrating the various features and use cases of the go-map-navigator library.

## Examples Overview

### 1. Basic Usage (`basic/`)
- **File**: `main.go`
- **Description**: Demonstrates fundamental map navigation operations
- **Features**:
  - Basic map navigation with nested properties
  - Array navigation with index and wildcard access
  - Read-only navigation
  - Property creation with `CreateProperty` flag

**Key Concepts**:
- Creating MapNavigator instances
- Using `VisitNode` method for navigation
- Working with nested map structures
- Accessing array elements by index

### 2. Modifiers (`modifiers/`)
- **File**: `main.go`
- **Description**: Shows how to use various node modifiers to transform data
- **Features**:
  - Replace modifier for string substitution
  - Set modifier for adding new properties
  - Custom modifier functions
  - Chaining multiple modifiers

**Key Concepts**:
- MapNodeModifier interface
- Built-in modifier implementations
- Custom modifier creation
- Sequential modifier application

### 3. JSON Processing (`json-processing/`)
- **File**: `main.go`
- **Description**: Real-world JSON data processing scenarios
- **Features**:
  - JSON data extraction from complex structures
  - JSON data transformation with business logic
  - JSON validation and cleaning
  - Complex nested navigation patterns

**Key Concepts**:
- Working with JSON-unmarshaled data
- Deep navigation in complex structures
- Data transformation patterns
- Practical API response processing

### 4. Conditional Processing (`conditional-processing/`)
- **File**: `main.go`
- **Description**: Advanced conditional logic and data processing
- **Features**:
  - Basic conditional modifiers
  - Complex multi-condition logic
  - Conditional data transformation
  - Template-based conditions

**Key Concepts**:
- Conditional modifier usage
- Complex business rule implementation
- Template-based conditional logic
- Multi-step conditional processing

### 5. Data Validation (`data-validation/`)
- **File**: `main.go`
- **Description**: Comprehensive data validation and cleaning examples
- **Features**:
  - Basic field validation
  - Data format normalization
  - Batch validation with error reporting
  - Format validation and standardization

**Key Concepts**:
- Input validation patterns
- Data cleaning and normalization
- Error collection and reporting
- Format standardization

## Running the Examples

To run any example, navigate to the example directory and execute:

```bash
cd examples/basic
go run main.go
```

Or run from the project root:

```bash
go run examples/basic/main.go
```

## Common Patterns

### 1. Creating a Basic Navigator

```go
navigator := &mapnavigator.MapNavigator{
    NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
        return o // No modification
    }),
    ReadOnly: true,
}
```

### 2. Navigating to Nested Properties

```go
// Navigate to user.profile.email
result, err := navigator.VisitNode(data, "user", "profile", "email")
```

### 3. Using Wildcards for Array Processing

```go
// Get all product names
result, err := navigator.VisitNode(data, "products", "*", "name")
```

### 4. Creating Custom Modifiers

```go
customModifier := models.MapNodeModifierFunc(func(o interface{}) interface{} {
    // Apply custom transformation logic
    return transformedValue
})
```

### 5. Conditional Processing

```go
conditionalModifier := models.MapNodeModifierFunc(func(o interface{}) interface{} {
    if condition {
        return modifiedValue
    }
    return o
})
```

## Best Practices

1. **Error Handling**: Always check for errors returned by `VisitNode`
2. **ReadOnly Mode**: Use `ReadOnly: true` when you only need to read data
3. **Modifier Reuse**: Create reusable modifiers for common transformations
4. **Type Safety**: Always perform type assertions when working with interface{} values
5. **Performance**: Consider the impact of deep navigation on large data structures

## Advanced Features

### Template-Based Conditions
The library supports template-based conditional logic using Go's text/template syntax:

```go
config := &models.MapNodeConditionalModifierConfig{
    Condition: "{{.NodeValue}} > 100",
    // ... other configuration
}
```

### Composite Modifiers
Multiple modifiers can be chained together for complex transformations:

```go
// First apply validation, then transformation
validator := createValidator()
transformer := createTransformer()

// Apply in sequence
navigator1 := &mapnavigator.MapNavigator{NodeModifier: validator}
navigator2 := &mapnavigator.MapNavigator{NodeModifier: transformer}
```

## Testing Your Code

The examples include comprehensive test patterns that you can adapt for your own use cases:

1. **Unit Testing**: Test individual modifiers in isolation
2. **Integration Testing**: Test complete navigation chains
3. **Edge Case Testing**: Test with nil values, empty structures, and invalid data
4. **Performance Testing**: Benchmark critical navigation paths

## Getting Help

- Check the main README.md for API documentation
- Look at the test files for additional usage patterns
- Examine the source code for detailed implementation details
- Review the examples for common use cases

## Contributing

If you have additional examples or improvements to existing examples, please:

1. Follow the existing code style and patterns
2. Include comprehensive comments explaining the logic
3. Add error handling for robust examples
4. Test your examples thoroughly
5. Update this README with any new examples