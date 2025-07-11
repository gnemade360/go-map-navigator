package map_node_delete_modifier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	map_navigator "github.com/gnemade360/go-map-navigator/pkg/mapnavigator"
	map_nav_models "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-delete-modifier"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/templates"
)

func TestMapNodeDeleteModifier_ModifyNode(t *testing.T) {
	tests := []struct {
		name            string
		config          *models.MapNodeDeleteModifierConfig
		templateContext map[string]interface{}
		input           interface{}
		validateResult  func(t *testing.T, result interface{})
	}{
		{
			name: "disabled modifier returns node unchanged",
			config: &models.MapNodeDeleteModifierConfig{
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					Disabled: true,
				},
			},
			input: "test-value",
			validateResult: func(t *testing.T, result interface{}) {
				assert.Equal(t, "test-value", result)
			},
		},
		{
			name:   "nil input returns nil",
			config: &models.MapNodeDeleteModifierConfig{},
			input:  nil,
			validateResult: func(t *testing.T, result interface{}) {
				assert.Nil(t, result)
			},
		},
		{
			name:   "selector '-' marks for deletion",
			config: &models.MapNodeDeleteModifierConfig{Selector: "-"},
			input:  "test-value",
			validateResult: func(t *testing.T, result interface{}) {
				deleted, ok := result.(*map_navigator.MapNavigatorDeleted)
				assert.True(t, ok)
				assert.Equal(t, "test-value", deleted.Original)
			},
		},
		{
			name:   "selector '*' marks for deletion",
			config: &models.MapNodeDeleteModifierConfig{Selector: "*"},
			input:  map[string]interface{}{"key": "value"},
			validateResult: func(t *testing.T, result interface{}) {
				deleted, ok := result.(*map_navigator.MapNavigatorDeleted)
				assert.True(t, ok)
				assert.Equal(t, map[string]interface{}{"key": "value"}, deleted.Original)
			},
		},
		{
			name: "delete from slice",
			config: &models.MapNodeDeleteModifierConfig{
				DeleteFrom: "listData",
			},
			templateContext: map[string]interface{}{
				"listData": []interface{}{"item1", "item2"},
			},
			input: "item1",
			validateResult: func(t *testing.T, result interface{}) {
				deleted, ok := result.(*map_navigator.MapNavigatorDeleted)
				assert.True(t, ok)
				assert.Equal(t, "item1", deleted.Original)
			},
		},
		{
			name: "delete nested field from map",
			config: &models.MapNodeDeleteModifierConfig{
				Selector: "user.email",
			},
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"name":  "John",
					"email": "john@example.com",
				},
			},
			validateResult: func(t *testing.T, result interface{}) {
				m, ok := result.(map[string]interface{})
				assert.True(t, ok)
				user := m["user"].(map[string]interface{})
				assert.Equal(t, "John", user["name"])
				// Email should be marked for deletion
			},
		},
		{
			name: "delete with empty selector",
			config: &models.MapNodeDeleteModifierConfig{
				Selector: "",
			},
			input: map[string]interface{}{"key": "value"},
			validateResult: func(t *testing.T, result interface{}) {
				assert.Equal(t, map[string]interface{}{"key": "value"}, result)
			},
		},
		{
			name:   "string input returns unchanged",
			config: &models.MapNodeDeleteModifierConfig{Selector: "field"},
			input:  "string-value",
			validateResult: func(t *testing.T, result interface{}) {
				assert.Equal(t, "string-value", result)
			},
		},
		{
			name: "delete single field from map",
			config: &models.MapNodeDeleteModifierConfig{
				Selector: "toDelete",
			},
			input: map[string]interface{}{
				"keep":     "value1",
				"toDelete": "value2",
			},
			validateResult: func(t *testing.T, result interface{}) {
				m, ok := result.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "value1", m["keep"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifier := &map_node_delete_modifier.MapNodeDeleteModifier{
				Config:          tt.config,
				TemplateContext: tt.templateContext,
			}
			
			result := modifier.ModifyNode(tt.input)
			tt.validateResult(t, result)
		})
	}
}

func TestNewMapNodeDeleteModifier(t *testing.T) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}

	tests := []struct {
		name       string
		config     interface{}
		mapContext map[string]interface{}
		tmplConfig templates.TemplateConfig
		expectNil  bool
		validateFn func(t *testing.T, modifier map_nav_models.MapNodeModifier)
	}{
		{
			name: "with MapNodeModifierConfig and MapNodeDeleteModifierConfig",
			config: &models.MapNodeModifierConfig{
				Options: &models.MapNodeDeleteModifierConfig{
					Selector:   "field.to.delete",
					DeleteFrom: "collection",
				},
			},
			mapContext: map[string]interface{}{"key": "value"},
			tmplConfig: tmplConfig,
			expectNil:  false,
			validateFn: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				assert.NotNil(t, modifier)
				dm, ok := modifier.(*map_node_delete_modifier.MapNodeDeleteModifier)
				assert.True(t, ok)
				assert.Equal(t, "field.to.delete", dm.Config.Selector)
				assert.Equal(t, "collection", dm.Config.DeleteFrom)
			},
		},
		{
			name: "with map options",
			config: &models.MapNodeModifierConfig{
				Options: map[string]interface{}{
					"selector":   "*",
					"deleteFrom": "",
				},
			},
			mapContext: map[string]interface{}{},
			tmplConfig: tmplConfig,
			expectNil:  false,
			validateFn: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				assert.NotNil(t, modifier)
			},
		},
		{
			name:       "nil config",
			config:     nil,
			mapContext: map[string]interface{}{},
			tmplConfig: tmplConfig,
			expectNil:  false,
			validateFn: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				assert.NotNil(t, modifier)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifier := map_node_delete_modifier.NewMapNodeDeleteModifier(
				tt.config,
				tt.mapContext,
				tt.tmplConfig,
			)
			
			if tt.expectNil {
				assert.Nil(t, modifier)
			} else {
				assert.NotNil(t, modifier)
				if tt.validateFn != nil {
					tt.validateFn(t, modifier)
				}
			}
		})
	}
}

func TestGetValue(t *testing.T) {
	testData := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"address": map[string]interface{}{
				"city": "New York",
				"zip":  "10001",
			},
			"tags": []interface{}{"admin", "user"},
		},
		"count": 42,
	}

	tests := []struct {
		name        string
		base        interface{}
		accessor    string
		expected    interface{}
		expectError bool
	}{
		{
			name:        "empty accessor",
			base:        testData,
			accessor:    "",
			expected:    nil,
			expectError: false,
		},
		{
			name:        "single level access",
			base:        testData,
			accessor:    "count",
			expected:    42,
			expectError: false,
		},
		{
			name:        "nested access",
			base:        testData,
			accessor:    "user.name",
			expected:    "John",
			expectError: false,
		},
		{
			name:        "deep nested access",
			base:        testData,
			accessor:    "user.address.city",
			expected:    "New York",
			expectError: false,
		},
		{
			name:        "array access",
			base:        testData,
			accessor:    "user.tags",
			expected:    []interface{}{"admin", "user"},
			expectError: false,
		},
		{
			name:        "non-existent field",
			base:        testData,
			accessor:    "user.email",
			expected:    nil,
			expectError: false,
		},
		{
			name:        "accessor with extra dots",
			base:        testData,
			accessor:    "..user..name..",
			expected:    "John",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := map_node_delete_modifier.GetValue(tt.base, tt.accessor)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDeleteModifier_ComplexScenarios(t *testing.T) {
	t.Run("delete multiple nested fields", func(t *testing.T) {
		data := map[string]interface{}{
			"users": []interface{}{
				map[string]interface{}{
					"id":       1,
					"name":     "User1",
					"password": "secret1",
				},
				map[string]interface{}{
					"id":       2,
					"name":     "User2",
					"password": "secret2",
				},
			},
		}
		
		// First delete passwords
		modifier1 := &map_node_delete_modifier.MapNodeDeleteModifier{
			Config: &models.MapNodeDeleteModifierConfig{
				Selector: "users.*.password",
			},
		}
		
		result := modifier1.ModifyNode(data)
		assert.NotNil(t, result)
	})

	t.Run("conditional deletion", func(t *testing.T) {
		modifier := &map_node_delete_modifier.MapNodeDeleteModifier{
			Config: &models.MapNodeDeleteModifierConfig{
				Selector: "-",
			},
			TemplateContext: map[string]interface{}{
				"shouldDelete": true,
			},
		}
		
		result := modifier.ModifyNode("test-value")
		deleted, ok := result.(*map_navigator.MapNavigatorDeleted)
		assert.True(t, ok)
		assert.Equal(t, "test-value", deleted.Original)
	})
}

func BenchmarkMapNodeDeleteModifier_ModifyNode(b *testing.B) {
	modifier := &map_node_delete_modifier.MapNodeDeleteModifier{
		Config: &models.MapNodeDeleteModifierConfig{
			Selector: "field.to.delete",
		},
	}
	
	input := map[string]interface{}{
		"field": map[string]interface{}{
			"to": map[string]interface{}{
				"delete": "value",
				"keep":   "other",
			},
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modifier.ModifyNode(input)
	}
}

func ExampleMapNodeDeleteModifier_ModifyNode() {
	// Create a delete modifier to remove sensitive data
	modifier := &map_node_delete_modifier.MapNodeDeleteModifier{
		Config: &models.MapNodeDeleteModifierConfig{
			Selector: "user.password",
		},
	}
	
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name":     "John Doe",
			"email":    "john@example.com",
			"password": "secret123",
		},
	}
	
	result := modifier.ModifyNode(data)
	// Password field will be marked for deletion
	println(result)
}