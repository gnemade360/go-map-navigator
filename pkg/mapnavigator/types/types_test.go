package types_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/types"
)

func TestNewCollectionMapKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		index    int
		expected types.CollectionMapKey
	}{
		{
			name:  "simple key with positive index",
			key:   "test-key",
			index: 5,
			expected: types.CollectionMapKey{
				Key:   "test-key",
				Index: 5,
			},
		},
		{
			name:  "empty key with zero index",
			key:   "",
			index: 0,
			expected: types.CollectionMapKey{
				Key:   "",
				Index: 0,
			},
		},
		{
			name:  "key with negative index",
			key:   "negative",
			index: -1,
			expected: types.CollectionMapKey{
				Key:   "negative",
				Index: -1,
			},
		},
		{
			name:  "key with large index",
			key:   "large",
			index: 999999,
			expected: types.CollectionMapKey{
				Key:   "large",
				Index: 999999,
			},
		},
		{
			name:  "special characters in key",
			key:   "special-@#$%-key",
			index: 42,
			expected: types.CollectionMapKey{
				Key:   "special-@#$%-key",
				Index: 42,
			},
		},
		{
			name:  "unicode key",
			key:   "unicode-键-🔑",
			index: 10,
			expected: types.CollectionMapKey{
				Key:   "unicode-键-🔑",
				Index: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.NewCollectionMapKey(tt.key, tt.index)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertToMapEntries(t *testing.T) {
	tests := []struct {
		name     string
		input    map[interface{}]interface{}
		validate func(t *testing.T, entries []types.MapEntry)
	}{
		{
			name:  "empty map",
			input: map[interface{}]interface{}{},
			validate: func(t *testing.T, entries []types.MapEntry) {
				assert.Empty(t, entries)
			},
		},
		{
			name: "string keys and values",
			input: map[interface{}]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			validate: func(t *testing.T, entries []types.MapEntry) {
				assert.Len(t, entries, 3)
				// Create a map for easy validation
				resultMap := make(map[interface{}]interface{})
				for _, entry := range entries {
					resultMap[entry.Key] = entry.Value
				}
				assert.Equal(t, "value1", resultMap["key1"])
				assert.Equal(t, "value2", resultMap["key2"])
				assert.Equal(t, "value3", resultMap["key3"])
			},
		},
		{
			name: "mixed key types",
			input: map[interface{}]interface{}{
				"string-key": "string-value",
				42:           "int-key",
				true:         "bool-key",
				3.14:         "float-key",
			},
			validate: func(t *testing.T, entries []types.MapEntry) {
				assert.Len(t, entries, 4)
				// Verify all entries are present
				resultMap := make(map[interface{}]interface{})
				for _, entry := range entries {
					resultMap[entry.Key] = entry.Value
				}
				assert.Equal(t, "string-value", resultMap["string-key"])
				assert.Equal(t, "int-key", resultMap[42])
				assert.Equal(t, "bool-key", resultMap[true])
				assert.Equal(t, "float-key", resultMap[3.14])
			},
		},
		{
			name: "mixed value types",
			input: map[interface{}]interface{}{
				"int":    123,
				"float":  45.67,
				"bool":   false,
				"nil":    nil,
				"slice":  []int{1, 2, 3},
				"map":    map[string]int{"a": 1},
			},
			validate: func(t *testing.T, entries []types.MapEntry) {
				assert.Len(t, entries, 6)
				resultMap := make(map[interface{}]interface{})
				for _, entry := range entries {
					resultMap[entry.Key] = entry.Value
				}
				assert.Equal(t, 123, resultMap["int"])
				assert.Equal(t, 45.67, resultMap["float"])
				assert.Equal(t, false, resultMap["bool"])
				assert.Nil(t, resultMap["nil"])
				assert.Equal(t, []int{1, 2, 3}, resultMap["slice"])
			},
		},
		{
			name: "single entry",
			input: map[interface{}]interface{}{
				"only": "one",
			},
			validate: func(t *testing.T, entries []types.MapEntry) {
				assert.Len(t, entries, 1)
				assert.Equal(t, "only", entries[0].Key)
				assert.Equal(t, "one", entries[0].Value)
			},
		},
		{
			name:  "nil map",
			input: nil,
			validate: func(t *testing.T, entries []types.MapEntry) {
				assert.Empty(t, entries)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.ConvertToMapEntries(tt.input)
			tt.validate(t, result)
		})
	}
}

func TestConvertFromMapEntries(t *testing.T) {
	tests := []struct {
		name     string
		input    []types.MapEntry
		expected map[interface{}]interface{}
	}{
		{
			name:     "empty slice",
			input:    []types.MapEntry{},
			expected: map[interface{}]interface{}{},
		},
		{
			name:     "nil slice",
			input:    nil,
			expected: map[interface{}]interface{}{},
		},
		{
			name: "string entries",
			input: []types.MapEntry{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
				{Key: "key3", Value: "value3"},
			},
			expected: map[interface{}]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		{
			name: "mixed type entries",
			input: []types.MapEntry{
				{Key: "string", Value: "text"},
				{Key: 42, Value: "number"},
				{Key: true, Value: "boolean"},
				{Key: 3.14, Value: "pi"},
			},
			expected: map[interface{}]interface{}{
				"string": "text",
				42:       "number",
				true:     "boolean",
				3.14:     "pi",
			},
		},
		{
			name: "entries with nil values",
			input: []types.MapEntry{
				{Key: "exists", Value: "value"},
				{Key: "nil-value", Value: nil},
			},
			expected: map[interface{}]interface{}{
				"exists":    "value",
				"nil-value": nil,
			},
		},
		{
			name: "duplicate keys (last wins)",
			input: []types.MapEntry{
				{Key: "dup", Value: "first"},
				{Key: "other", Value: "middle"},
				{Key: "dup", Value: "last"},
			},
			expected: map[interface{}]interface{}{
				"dup":   "last",
				"other": "middle",
			},
		},
		{
			name: "single entry",
			input: []types.MapEntry{
				{Key: "single", Value: "entry"},
			},
			expected: map[interface{}]interface{}{
				"single": "entry",
			},
		},
		{
			name: "complex value types",
			input: []types.MapEntry{
				{Key: "slice", Value: []int{1, 2, 3}},
				{Key: "map", Value: map[string]string{"nested": "map"}},
				{Key: "struct", Value: struct{ Name string }{Name: "test"}},
			},
			expected: map[interface{}]interface{}{
				"slice":  []int{1, 2, 3},
				"map":    map[string]string{"nested": "map"},
				"struct": struct{ Name string }{Name: "test"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.ConvertFromMapEntries(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	// Test that converting to entries and back preserves the data
	testCases := []struct {
		name     string
		original map[interface{}]interface{}
	}{
		{
			name:     "empty map",
			original: map[interface{}]interface{}{},
		},
		{
			name: "simple strings",
			original: map[interface{}]interface{}{
				"a": "1",
				"b": "2",
				"c": "3",
			},
		},
		{
			name: "mixed types",
			original: map[interface{}]interface{}{
				"string": "value",
				123:      456,
				true:     false,
				3.14:     2.71,
				"nil":    nil,
			},
		},
		{
			name: "nested structures",
			original: map[interface{}]interface{}{
				"nested": map[interface{}]interface{}{
					"inner": "value",
				},
				"array": []interface{}{1, "two", 3.0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert to entries
			entries := types.ConvertToMapEntries(tc.original)
			
			// Convert back to map
			result := types.ConvertFromMapEntries(entries)
			
			// Should be equal
			assert.Equal(t, tc.original, result)
		})
	}
}

func TestMapEntry_Fields(t *testing.T) {
	// Test that MapEntry fields are accessible
	entry := types.MapEntry{
		Key:   "test-key",
		Value: "test-value",
	}
	
	assert.Equal(t, "test-key", entry.Key)
	assert.Equal(t, "test-value", entry.Value)
	
	// Test with different types
	entry2 := types.MapEntry{
		Key:   42,
		Value: []int{1, 2, 3},
	}
	
	assert.Equal(t, 42, entry2.Key)
	assert.Equal(t, []int{1, 2, 3}, entry2.Value)
}

func TestCollectionMapKey_Fields(t *testing.T) {
	// Test that CollectionMapKey fields are accessible
	cmk := types.CollectionMapKey{
		Key:   "collection",
		Index: 10,
	}
	
	assert.Equal(t, "collection", cmk.Key)
	assert.Equal(t, 10, cmk.Index)
}

func BenchmarkConvertToMapEntries(b *testing.B) {
	// Create a test map
	testMap := make(map[interface{}]interface{})
	for i := 0; i < 100; i++ {
		testMap[i] = i * 2
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = types.ConvertToMapEntries(testMap)
	}
}

func BenchmarkConvertFromMapEntries(b *testing.B) {
	// Create test entries
	entries := make([]types.MapEntry, 100)
	for i := 0; i < 100; i++ {
		entries[i] = types.MapEntry{Key: i, Value: i * 2}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = types.ConvertFromMapEntries(entries)
	}
}

func BenchmarkRoundTripConversion(b *testing.B) {
	// Create a test map
	testMap := make(map[interface{}]interface{})
	for i := 0; i < 50; i++ {
		testMap[i] = i * 2
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entries := types.ConvertToMapEntries(testMap)
		_ = types.ConvertFromMapEntries(entries)
	}
}

func ExampleNewCollectionMapKey() {
	// Create a new collection map key
	key := types.NewCollectionMapKey("users", 5)
	
	// Use the key
	println(key.Key)   // "users"
	println(key.Index) // 5
}

func ExampleConvertToMapEntries() {
	// Create a map with mixed types
	originalMap := map[interface{}]interface{}{
		"name":  "John",
		"age":   30,
		"active": true,
	}
	
	// Convert to entries
	entries := types.ConvertToMapEntries(originalMap)
	
	// Process entries
	for _, entry := range entries {
		println(entry.Key, ":", entry.Value)
	}
}

func ExampleConvertFromMapEntries() {
	// Create entries
	entries := []types.MapEntry{
		{Key: "name", Value: "Alice"},
		{Key: "score", Value: 95},
		{Key: "passed", Value: true},
	}
	
	// Convert to map
	resultMap := types.ConvertFromMapEntries(entries)
	
	// Use the map
	println(resultMap["name"])   // "Alice"
	println(resultMap["score"])  // 95
	println(resultMap["passed"]) // true
}

// Test for ensuring proper type handling
func TestTypeHandling(t *testing.T) {
	// Test that the conversion functions handle Go's type system correctly
	t.Run("interface{} preservation", func(t *testing.T) {
		original := map[interface{}]interface{}{
			"key": interface{}("value"),
		}
		
		entries := types.ConvertToMapEntries(original)
		result := types.ConvertFromMapEntries(entries)
		
		assert.Equal(t, original, result)
		assert.Equal(t, reflect.TypeOf(original["key"]), reflect.TypeOf(result["key"]))
	})
}