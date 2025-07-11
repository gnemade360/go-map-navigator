// Package mapnavigator provides a flexible library for navigating, querying, and modifying nested map and slice data structures.
//
// This library allows you to traverse deeply nested data structures using simple string paths,
// apply modifiers to transform data, and handle complex data manipulation scenarios with ease.
//
// Basic Usage:
//
//	// Create a navigator with a modifier
//	navigator := mapnavigator.NewMapNavigator(someModifier)
//
//	// Navigate to a nested value
//	result, err := navigator.VisitNode(data, "users", "0", "name")
//
// Key Features:
//
// - Navigate nested maps and slices using dot notation or array indices
// - Apply modifiers to transform data during navigation
// - Support for wildcard operations with "*" to process all elements
// - Conditional processing with various condition types
// - Read-only mode for safe data querying
// - Property creation mode for dynamic data structure building
//
// Path Syntax:
//
// - "key" - access map key
// - "0", "1", "2" - access array element by index
// - "*" - apply operation to all elements in map or array
// - "-" - apply modifier to current node
//
// The library supports various modifiers including:
//
// - Set operations to modify values
// - Delete operations to remove elements
// - Replace operations for value substitution
// - Conditional operations based on various criteria
// - Composite operations for complex transformations
// - Collection expansion for flattening nested structures
//
// Thread Safety:
//
// MapNavigator instances are not thread-safe. Create separate instances
// for concurrent operations or use appropriate synchronization mechanisms.
//
// Author: Ganesh Nemade
package mapnavigator