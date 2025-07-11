package types

// CollectionMapKey represents a key in a collection map
type CollectionMapKey struct {
	Key   string
	Index int
}

// NewCollectionMapKey creates a new CollectionMapKey
func NewCollectionMapKey(key string, index int) CollectionMapKey {
	return CollectionMapKey{
		Key:   key,
		Index: index,
	}
}

// MapEntry represents an entry in a map
type MapEntry struct {
	Key   interface{}
	Value interface{}
}

// ConvertToMapEntries converts a map to a slice of MapEntry
func ConvertToMapEntries(m map[interface{}]interface{}) []MapEntry {
	entries := make([]MapEntry, 0, len(m))
	for k, v := range m {
		entries = append(entries, MapEntry{Key: k, Value: v})
	}
	return entries
}

// ConvertFromMapEntries converts a slice of MapEntry back to a map
func ConvertFromMapEntries(entries []MapEntry) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	for _, entry := range entries {
		m[entry.Key] = entry.Value
	}
	return m
}