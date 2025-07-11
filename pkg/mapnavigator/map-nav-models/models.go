package models

// MapNodeModifierFunc is a function type that implements the MapNodeModifier interface.
// It provides a convenient way to create modifiers from simple functions.
type MapNodeModifierFunc func(interface{}) interface{}

// ModifyNode implements the MapNodeModifier interface by calling the underlying function.
// This allows functions to be used anywhere a MapNodeModifier is expected.
func (m MapNodeModifierFunc) ModifyNode(i interface{}) interface{} {
	return m(i)
}

// MapNodeModifier defines the interface for modifying nodes during navigation.
// Implementations of this interface can be used to transform data as it's being navigated.
type MapNodeModifier interface {
	// ModifyNode takes an input node and returns a modified version.
	// The modification can be any transformation, including deletion (return MapNavigatorDeleted).
	ModifyNode(interface{}) interface{}
}
