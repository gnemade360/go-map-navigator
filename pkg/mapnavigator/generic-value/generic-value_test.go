package generic_value_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	genericvalue "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/generic-value"
)

func TestToGenericValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: nil,
		},
		{
			name:     "string value",
			input:    "test",
			expected: "test",
		},
		{
			name:     "int value",
			input:    42,
			expected: 42,
		},
		{
			name:     "float value",
			input:    3.14,
			expected: 3.14,
		},
		{
			name:     "bool value",
			input:    true,
			expected: true,
		},
		{
			name:     "slice value",
			input:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "map value",
			input:    map[string]int{"a": 1},
			expected: map[string]int{"a": 1},
		},
		{
			name:     "struct value",
			input:    struct{ Name string }{Name: "test"},
			expected: struct{ Name string }{Name: "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := genericvalue.ToGenericValue(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: "",
		},
		{
			name:     "string value",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "int value",
			input:    42,
			expected: "42",
		},
		{
			name:     "int8 value",
			input:    int8(127),
			expected: "127",
		},
		{
			name:     "int16 value",
			input:    int16(32767),
			expected: "32767",
		},
		{
			name:     "int32 value",
			input:    int32(2147483647),
			expected: "2147483647",
		},
		{
			name:     "int64 value",
			input:    int64(9223372036854775807),
			expected: "9223372036854775807",
		},
		{
			name:     "uint value",
			input:    uint(42),
			expected: "42",
		},
		{
			name:     "uint8 value",
			input:    uint8(255),
			expected: "255",
		},
		{
			name:     "uint16 value",
			input:    uint16(65535),
			expected: "65535",
		},
		{
			name:     "uint32 value",
			input:    uint32(4294967295),
			expected: "4294967295",
		},
		{
			name:     "uint64 value",
			input:    uint64(18446744073709551615),
			expected: "18446744073709551615",
		},
		{
			name:     "float32 value",
			input:    float32(3.14),
			expected: "3.14",
		},
		{
			name:     "float64 value",
			input:    float64(3.14159),
			expected: "3.14159",
		},
		{
			name:     "bool true",
			input:    true,
			expected: "true",
		},
		{
			name:     "bool false",
			input:    false,
			expected: "false",
		},
		{
			name:     "slice value",
			input:    []int{1, 2, 3},
			expected: "[1 2 3]",
		},
		{
			name:     "map value",
			input:    map[string]int{"a": 1},
			expected: "map[a:1]",
		},
		{
			name:     "struct value",
			input:    struct{ Name string }{Name: "test"},
			expected: "{test}",
		},
		{
			name:     "negative int",
			input:    -42,
			expected: "-42",
		},
		{
			name:     "zero value",
			input:    0,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := genericvalue.ToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    int
		expectError bool
	}{
		{
			name:        "int value",
			input:       42,
			expected:    42,
			expectError: false,
		},
		{
			name:        "int8 value",
			input:       int8(127),
			expected:    127,
			expectError: false,
		},
		{
			name:        "int16 value",
			input:       int16(32767),
			expected:    32767,
			expectError: false,
		},
		{
			name:        "int32 value",
			input:       int32(2147483647),
			expected:    2147483647,
			expectError: false,
		},
		{
			name:        "int64 value within int range",
			input:       int64(42),
			expected:    42,
			expectError: false,
		},
		{
			name:        "uint value",
			input:       uint(42),
			expected:    42,
			expectError: false,
		},
		{
			name:        "uint8 value",
			input:       uint8(255),
			expected:    255,
			expectError: false,
		},
		{
			name:        "uint16 value",
			input:       uint16(65535),
			expected:    65535,
			expectError: false,
		},
		{
			name:        "uint32 value",
			input:       uint32(42),
			expected:    42,
			expectError: false,
		},
		{
			name:        "uint64 value within int range",
			input:       uint64(42),
			expected:    42,
			expectError: false,
		},
		{
			name:        "float32 value",
			input:       float32(42.0),
			expected:    42,
			expectError: false,
		},
		{
			name:        "float64 value",
			input:       float64(42.0),
			expected:    42,
			expectError: false,
		},
		{
			name:        "float with decimal truncated",
			input:       float64(42.7),
			expected:    42,
			expectError: false,
		},
		{
			name:        "string valid int",
			input:       "42",
			expected:    42,
			expectError: false,
		},
		{
			name:        "string negative int",
			input:       "-42",
			expected:    -42,
			expectError: false,
		},
		{
			name:        "string with spaces",
			input:       " 42 ",
			expected:    42,
			expectError: false,
		},
		{
			name:        "string invalid int",
			input:       "not a number",
			expected:    0,
			expectError: true,
		},
		{
			name:        "string empty",
			input:       "",
			expected:    0,
			expectError: true,
		},
		{
			name:        "bool value",
			input:       true,
			expected:    0,
			expectError: true,
		},
		{
			name:        "nil value",
			input:       nil,
			expected:    0,
			expectError: true,
		},
		{
			name:        "slice value",
			input:       []int{1, 2, 3},
			expected:    0,
			expectError: true,
		},
		{
			name:        "map value",
			input:       map[string]int{"a": 1},
			expected:    0,
			expectError: true,
		},
		{
			name:        "zero value",
			input:       0,
			expected:    0,
			expectError: false,
		},
		{
			name:        "negative value",
			input:       -100,
			expected:    -100,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := genericvalue.ToInt(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGetReflectValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		validate func(t *testing.T, v reflect.Value)
	}{
		{
			name:  "nil value",
			input: nil,
			validate: func(t *testing.T, v reflect.Value) {
				assert.False(t, v.IsValid())
			},
		},
		{
			name:  "string value",
			input: "test",
			validate: func(t *testing.T, v reflect.Value) {
				assert.True(t, v.IsValid())
				assert.Equal(t, reflect.String, v.Kind())
				assert.Equal(t, "test", v.String())
			},
		},
		{
			name:  "int value",
			input: 42,
			validate: func(t *testing.T, v reflect.Value) {
				assert.True(t, v.IsValid())
				assert.Equal(t, reflect.Int, v.Kind())
				assert.Equal(t, int64(42), v.Int())
			},
		},
		{
			name:  "pointer value",
			input: new(int),
			validate: func(t *testing.T, v reflect.Value) {
				assert.True(t, v.IsValid())
				assert.Equal(t, reflect.Ptr, v.Kind())
			},
		},
		{
			name:  "slice value",
			input: []int{1, 2, 3},
			validate: func(t *testing.T, v reflect.Value) {
				assert.True(t, v.IsValid())
				assert.Equal(t, reflect.Slice, v.Kind())
				assert.Equal(t, 3, v.Len())
			},
		},
		{
			name:  "map value",
			input: map[string]int{"a": 1, "b": 2},
			validate: func(t *testing.T, v reflect.Value) {
				assert.True(t, v.IsValid())
				assert.Equal(t, reflect.Map, v.Kind())
				assert.Equal(t, 2, v.Len())
			},
		},
		{
			name:  "struct value",
			input: struct{ Name string }{Name: "test"},
			validate: func(t *testing.T, v reflect.Value) {
				assert.True(t, v.IsValid())
				assert.Equal(t, reflect.Struct, v.Kind())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := genericvalue.GetReflectValue(tt.input)
			tt.validate(t, result)
		})
	}
}

func TestToInt_EdgeCases(t *testing.T) {
	// Test overflow scenarios
	t.Run("int64 overflow", func(t *testing.T) {
		// This test depends on the system's int size
		if reflect.TypeOf(int(0)).Size() == 4 {
			// 32-bit system
			_, err := genericvalue.ToInt(int64(math.MaxInt32 + 1))
			assert.Error(t, err, "should error on overflow")
		}
	})

	t.Run("uint64 large value", func(t *testing.T) {
		// Test with a uint64 value that might overflow int
		if reflect.TypeOf(int(0)).Size() == 4 {
			_, err := genericvalue.ToInt(uint64(math.MaxUint32))
			assert.Error(t, err, "should error on overflow")
		}
	})

	t.Run("string with leading zeros", func(t *testing.T) {
		result, err := genericvalue.ToInt("0042")
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("string with plus sign", func(t *testing.T) {
		result, err := genericvalue.ToInt("+42")
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})
}

func BenchmarkToString(b *testing.B) {
	testCases := []struct {
		name  string
		value interface{}
	}{
		{"nil", nil},
		{"string", "test string"},
		{"int", 42},
		{"float", 3.14159},
		{"bool", true},
		{"slice", []int{1, 2, 3, 4, 5}},
		{"map", map[string]int{"a": 1, "b": 2, "c": 3}},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = genericvalue.ToString(tc.value)
			}
		})
	}
}

func BenchmarkToInt(b *testing.B) {
	testCases := []struct {
		name  string
		value interface{}
	}{
		{"int", 42},
		{"string", "42"},
		{"float", 42.0},
		{"int64", int64(42)},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = genericvalue.ToInt(tc.value)
			}
		})
	}
}

func ExampleToString() {
	fmt.Println(genericvalue.ToString(42))
	fmt.Println(genericvalue.ToString("hello"))
	fmt.Println(genericvalue.ToString(true))
	fmt.Println(genericvalue.ToString(nil))
	// Output:
	// 42
	// hello
	// true
	// 
}

func ExampleToInt() {
	val1, _ := genericvalue.ToInt(42)
	fmt.Println(val1)
	
	val2, _ := genericvalue.ToInt("123")
	fmt.Println(val2)
	
	val3, _ := genericvalue.ToInt(3.14)
	fmt.Println(val3)
	// Output:
	// 42
	// 123
	// 3
}