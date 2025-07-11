package types

// Type constants for different data types used throughout the library.
const (
	// Map represents the map data type
	Map = "map"
)

// CollectionMapKey represents a key in a collection map with an associated index.
// This is useful for handling collections where you need to track both the key and position.
type CollectionMapKey struct {
	// Key is the string key identifier
	Key   string
	// Index is the numeric position within the collection
	Index int
}

// NewCollectionMapKey creates a new CollectionMapKey with the specified key and index.
// This is the preferred way to create CollectionMapKey instances.
func NewCollectionMapKey(key string, index int) CollectionMapKey {
	return CollectionMapKey{
		Key:   key,
		Index: index,
	}
}

// MapEntry represents a key-value pair from a map.
// This is useful for converting maps to slices or for iteration purposes.
type MapEntry struct {
	// Key is the map key (can be any type)
	Key   interface{}
	// Value is the map value (can be any type)
	Value interface{}
}

// ConvertToMapEntries converts a map[interface{}]interface{} to a slice of MapEntry.
// This is useful for scenarios where you need to process map data as a slice.
func ConvertToMapEntries(m map[interface{}]interface{}) []MapEntry {
	entries := make([]MapEntry, 0, len(m))
	for k, v := range m {
		entries = append(entries, MapEntry{Key: k, Value: v})
	}
	return entries
}

// ConvertFromMapEntries converts a slice of MapEntry back to a map[interface{}]interface{}.
// This is the inverse operation of ConvertToMapEntries.
func ConvertFromMapEntries(entries []MapEntry) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	for _, entry := range entries {
		m[entry.Key] = entry.Value
	}
	return m
}