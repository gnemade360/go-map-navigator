package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-nav-models"
)

func TestMapNodeModifierFunc_ModifyNode(t *testing.T) {
	tests := []struct {
		name     string
		modifier models.MapNodeModifierFunc
		input    interface{}
		expected interface{}
	}{
		{
			name: "string modifier",
			modifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if s, ok := v.(string); ok {
					return s + "-modified"
				}
				return v
			}),
			input:    "test",
			expected: "test-modified",
		},
		{
			name: "int modifier",
			modifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if i, ok := v.(int); ok {
					return i * 2
				}
				return v
			}),
			input:    5,
			expected: 10,
		},
		{
			name: "nil handler",
			modifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if v == nil {
					return "nil-value"
				}
				return v
			}),
			input:    nil,
			expected: "nil-value",
		},
		{
			name: "map modifier",
			modifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if m, ok := v.(map[string]interface{}); ok {
					m["modified"] = true
					return m
				}
				return v
			}),
			input:    map[string]interface{}{"key": "value"},
			expected: map[string]interface{}{"key": "value", "modified": true},
		},
		{
			name: "slice modifier",
			modifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if s, ok := v.([]interface{}); ok {
					return append(s, "new-item")
				}
				return v
			}),
			input:    []interface{}{"item1", "item2"},
			expected: []interface{}{"item1", "item2", "new-item"},
		},
		{
			name: "identity modifier",
			modifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
				return v
			}),
			input:    "unchanged",
			expected: "unchanged",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.modifier.ModifyNode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapNodeModifier_Interface(t *testing.T) {
	// Test that MapNodeModifierFunc implements MapNodeModifier interface
	var _ models.MapNodeModifier = models.MapNodeModifierFunc(func(v interface{}) interface{} {
		return v
	})
	
	// Test polymorphic behavior
	modifiers := []models.MapNodeModifier{
		models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return "modifier1"
		}),
		models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return "modifier2"
		}),
	}
	
	for i, mod := range modifiers {
		result := mod.ModifyNode("test")
		assert.Equal(t, "modifier"+string(rune('1'+i)), result)
	}
}

func TestMapNodeModifierFunc_ChainedModifiers(t *testing.T) {
	// Test chaining multiple modifiers
	modifier1 := models.MapNodeModifierFunc(func(v interface{}) interface{} {
		if s, ok := v.(string); ok {
			return s + "-first"
		}
		return v
	})
	
	modifier2 := models.MapNodeModifierFunc(func(v interface{}) interface{} {
		if s, ok := v.(string); ok {
			return s + "-second"
		}
		return v
	})
	
	input := "test"
	result1 := modifier1.ModifyNode(input)
	result2 := modifier2.ModifyNode(result1)
	
	assert.Equal(t, "test-first-second", result2)
}

func TestMapNodeModifierFunc_ComplexTypes(t *testing.T) {
	type CustomStruct struct {
		Name  string
		Value int
	}
	
	modifier := models.MapNodeModifierFunc(func(v interface{}) interface{} {
		if cs, ok := v.(CustomStruct); ok {
			cs.Value = cs.Value * 2
			return cs
		}
		return v
	})
	
	input := CustomStruct{Name: "test", Value: 10}
	result := modifier.ModifyNode(input)
	
	expected := CustomStruct{Name: "test", Value: 20}
	assert.Equal(t, expected, result)
}

func BenchmarkMapNodeModifierFunc_ModifyNode(b *testing.B) {
	modifier := models.MapNodeModifierFunc(func(v interface{}) interface{} {
		if s, ok := v.(string); ok {
			return s + "-modified"
		}
		return v
	})
	
	input := "benchmark-test"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modifier.ModifyNode(input)
	}
}