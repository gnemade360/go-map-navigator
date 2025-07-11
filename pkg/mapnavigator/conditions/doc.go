// Package conditions provides conditional evaluation capabilities for the mapnavigator library.
//
// This package contains various condition types that can be used to conditionally
// apply modifiers or transformations to data based on specific criteria.
//
// Supported Conditions:
//
// - Value equality conditions
// - Type checking conditions
// - Existence conditions
// - Pattern matching conditions
// - Range conditions
// - Custom condition functions
//
// Usage with MapNavigator:
//
//	condition := conditions.NewEqualsCondition("active", true)
//	modifier := conditionalmodifier.NewConditionalModifier(condition, someModifier)
//	navigator := mapnavigator.NewMapNavigator(modifier)
//
// Thread Safety:
//
// Condition instances are generally safe for concurrent read operations,
// but should not be modified concurrently.
package conditions
