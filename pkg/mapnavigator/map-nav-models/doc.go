// Package models defines the core interfaces and types for the mapnavigator library.
//
// This package contains the fundamental interfaces that enable the modular
// architecture of the mapnavigator library, allowing different components
// to work together seamlessly.
//
// Key Interfaces:
//
// - MapNodeModifier: Interface for implementing node modification logic
// - MapNodeModifierFunc: Function type that implements MapNodeModifier
//
// Usage:
//
//	// Implement a custom modifier
//	type MyModifier struct{}
//	func (m MyModifier) ModifyNode(node interface{}) interface{} {
//		// Custom modification logic
//		return modifiedNode
//	}
//
//	// Use function as modifier
//	modifierFunc := models.MapNodeModifierFunc(func(node interface{}) interface{} {
//		return processedNode
//	})
//
// Thread Safety:
//
// Modifier implementations should be designed to be thread-safe if they
// will be used concurrently. The interfaces themselves are thread-safe.
package models