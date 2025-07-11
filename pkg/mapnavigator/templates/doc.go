// Package templates provides template processing capabilities for the mapnavigator library.
//
// This package allows you to define templates that can be used to generate
// dynamic content or transformations based on data structure patterns.
//
// Templates can be used to:
//
// - Generate dynamic content based on data patterns
// - Apply consistent transformations across similar data structures
// - Create reusable data processing patterns
// - Implement complex data mapping scenarios
//
// Usage:
//
//	template := templates.NewTemplate(pattern)
//	result := template.Apply(data)
//
// Thread Safety:
//
// Template instances should be created once and can be safely used
// concurrently for read operations.
package templates