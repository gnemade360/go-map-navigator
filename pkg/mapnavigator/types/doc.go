// Package types provides common type definitions and utilities for the mapnavigator library.
//
// This package contains type constants, utility structures, and helper functions
// that are used throughout the mapnavigator library for type safety and consistency.
//
// Key Components:
//
// - Type constants for different data types
// - CollectionMapKey for handling collection keys with indices
// - MapEntry for representing key-value pairs
// - Conversion utilities between maps and structured data
//
// Usage:
//
//	// Create a collection key
//	key := types.NewCollectionMapKey("users", 0)
//
//	// Convert map to entries
//	entries := types.ConvertToMapEntries(myMap)
//
// Thread Safety:
//
// All utilities in this package are safe for concurrent use.
package types