package map_node_set_modifier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	map_nav_models "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-set-modifier"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/templates"
)

func TestMapNodeSetModifier_ModifyNode(t *testing.T) {
	tests := []struct {
		name            string
		config          *models.MapNodeSetModifierConfig
		templateContext map[string]interface{}
		input           interface{}
		expected        interface{}
	}{
		{
			name: "disabled modifier returns node unchanged",
			config: &models.MapNodeSetModifierConfig{
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					Disabled: true,
				},
			},
			input:    "test-value",
			expected: "test-value",
		},
		{
			name:     "nil input returns nil",
			config:   &models.MapNodeSetModifierConfig{},
			input:    nil,
			expected: nil,
		},
		{
			name: "simple string value set",
			config: &models.MapNodeSetModifierConfig{
				ValueToSet: "new-value",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{},
				},
			},
			input:    "old-value",
			expected: "new-value",
		},
		{
			name: "set value with template interpolation",
			config: &models.MapNodeSetModifierConfig{
				ValueToSet: "{{.Prefix}}-{{.NodeValue}}-{{.Suffix}}",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{
						LeftDelim:  "{{",
						RightDelim: "}}",
					},
				},
			},
			templateContext: map[string]interface{}{
				"Prefix": "START",
				"Suffix": "END",
			},
			input:    "middle",
			expected: "START-middle-END",
		},
		{
			name: "set value with vars interpolation",
			config: &models.MapNodeSetModifierConfig{
				ValueToSet: "{{.vars.type}}: {{.NodeValue}}",
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					Vars: map[string]interface{}{
						"type": "STRING",
					},
				},
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{
						LeftDelim:  "{{",
						RightDelim: "}}",
					},
				},
			},
			input:    "test",
			expected: "STRING: test",
		},
		{
			name: "set value for numeric input",
			config: &models.MapNodeSetModifierConfig{
				ValueToSet: "Number: {{.NodeValue}}",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{
						LeftDelim:  "{{",
						RightDelim: "}}",
					},
				},
			},
			input:    42,
			expected: "Number: 42",
		},
		{
			name: "set value in map with selector",
			config: &models.MapNodeSetModifierConfig{
				MapNodeModifierConfig: &models.MapNodeModifierConfig{},
				ValueToSet: "updated",
				Selector:   "field",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{},
				},
			},
			input: map[string]interface{}{
				"field": "old",
				"other": "keep",
			},
			expected: map[string]interface{}{
				"field": "updated",
				"other": "keep",
			},
		},
		{
			name: "set nested value in map",
			config: &models.MapNodeSetModifierConfig{
				MapNodeModifierConfig: &models.MapNodeModifierConfig{},
				ValueToSet: "new-email",
				Selector:   "user.email",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{},
				},
			},
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"name":  "John",
					"email": "old@example.com",
				},
			},
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name":  "John",
					"email": "new-email",
				},
			},
		},
		{
			name: "create property if absent",
			config: &models.MapNodeSetModifierConfig{
				ValueToSet: "created-value",
				Selector:   "newField",
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					CreatePropertyIfAbsent: true,
				},
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{},
				},
			},
			input: map[string]interface{}{
				"existing": "value",
			},
			expected: map[string]interface{}{
				"existing": "value",
				"newField": "created-value",
			},
		},
		{
			name: "boolean input conversion",
			config: &models.MapNodeSetModifierConfig{
				ValueToSet: "Boolean value was: {{.NodeValue}}",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{
						LeftDelim:  "{{",
						RightDelim: "}}",
					},
				},
			},
			input:    true,
			expected: "Boolean value was: true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifier := &map_node_set_modifier.MapNodeSetModifier{
				Config:          tt.config,
				TemplateContext: tt.templateContext,
			}
			
			result := modifier.ModifyNode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewMapNodeSetModifier(t *testing.T) {
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
			name: "with MapNodeModifierConfig and MapNodeSetModifierConfig",
			config: &models.MapNodeModifierConfig{
				Options: &models.MapNodeSetModifierConfig{
					ValueToSet: "test-value",
					Selector:   "field",
				},
			},
			mapContext: map[string]interface{}{"key": "value"},
			tmplConfig: tmplConfig,
			expectNil:  false,
			validateFn: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				assert.NotNil(t, modifier)
				sm, ok := modifier.(*map_node_set_modifier.MapNodeSetModifier)
				assert.True(t, ok)
				assert.Equal(t, "test-value", sm.Config.ValueToSet)
				assert.Equal(t, "field", sm.Config.Selector)
			},
		},
		{
			name: "with map options",
			config: &models.MapNodeModifierConfig{
				Options: map[string]interface{}{
					"valueToSet":             "new-value",
					"selector":               "path.to.field",
					"createPropertyIfAbsent": true,
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
			modifier := map_node_set_modifier.NewMapNodeSetModifier(
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

func TestSetModifier_ComplexScenarios(t *testing.T) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}

	t.Run("set value with complex template and vars", func(t *testing.T) {
		modifier := &map_node_set_modifier.MapNodeSetModifier{
			Config: &models.MapNodeSetModifierConfig{
				ValueToSet: "{{.vars.env}}_{{.NodeValue}}_{{.timestamp}}",
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					Vars: map[string]interface{}{
						"env": "{{.Environment}}",
					},
				},
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
			},
			TemplateContext: map[string]interface{}{
				"Environment": "PROD",
				"timestamp":   "2024-01-01",
			},
		}
		
		result := modifier.ModifyNode("service")
		assert.Equal(t, "PROD_service_2024-01-01", result)
	})

	t.Run("set multiple nested values", func(t *testing.T) {
		data := map[string]interface{}{
			"config": map[string]interface{}{
				"database": map[string]interface{}{
					"host": "localhost",
					"port": 5432,
				},
			},
		}
		
		// First modifier - set host
		modifier1 := &map_node_set_modifier.MapNodeSetModifier{
			Config: &models.MapNodeSetModifierConfig{
				ValueToSet: "db.example.com",
				Selector:   "config.database.host",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
			},
		}
		
		// Second modifier - set port
		modifier2 := &map_node_set_modifier.MapNodeSetModifier{
			Config: &models.MapNodeSetModifierConfig{
				ValueToSet: "3306",
				Selector:   "config.database.port",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
			},
		}
		
		result1 := modifier1.ModifyNode(data)
		result2 := modifier2.ModifyNode(result1)
		
		finalMap := result2.(map[string]interface{})
		dbConfig := finalMap["config"].(map[string]interface{})["database"].(map[string]interface{})
		assert.Equal(t, "db.example.com", dbConfig["host"])
		assert.Equal(t, "3306", dbConfig["port"])
	})

	t.Run("create nested property", func(t *testing.T) {
		modifier := &map_node_set_modifier.MapNodeSetModifier{
			Config: &models.MapNodeSetModifierConfig{
				ValueToSet: "new-value",
				Selector:   "new.nested.field",
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					CreatePropertyIfAbsent: true,
				},
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
			},
		}
		
		data := map[string]interface{}{
			"existing": "data",
		}
		
		result := modifier.ModifyNode(data)
		resultMap := result.(map[string]interface{})
		assert.Equal(t, "data", resultMap["existing"])
		// Verify nested structure was created
		if newMap, ok := resultMap["new"].(map[string]interface{}); ok {
			if nestedMap, ok := newMap["nested"].(map[string]interface{}); ok {
				assert.Equal(t, "new-value", nestedMap["field"])
			}
		}
	})
}

func TestSetModifier_EdgeCases(t *testing.T) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}

	t.Run("set empty value", func(t *testing.T) {
		modifier := &map_node_set_modifier.MapNodeSetModifier{
			Config: &models.MapNodeSetModifierConfig{
				ValueToSet: "",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
			},
		}
		
		result := modifier.ModifyNode("non-empty")
		assert.Equal(t, "", result)
	})

	t.Run("array input", func(t *testing.T) {
		modifier := &map_node_set_modifier.MapNodeSetModifier{
			Config: &models.MapNodeSetModifierConfig{
				ValueToSet: "modified",
				Selector:   "0",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
			},
		}
		
		input := []interface{}{"first", "second", "third"}
		result := modifier.ModifyNode(input)
		resultArr := result.([]interface{})
		assert.Equal(t, "modified", resultArr[0])
		assert.Equal(t, "second", resultArr[1])
		assert.Equal(t, "third", resultArr[2])
	})

	t.Run("nil value handling", func(t *testing.T) {
		modifier := &map_node_set_modifier.MapNodeSetModifier{
			Config: &models.MapNodeSetModifierConfig{
				ValueToSet: "{{.NodeValue}}",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
			},
		}
		
		result := modifier.ModifyNode((*string)(nil))
		assert.Equal(t, "<nil>", result)
	})
}

func BenchmarkMapNodeSetModifier_ModifyNode(b *testing.B) {
	modifier := &map_node_set_modifier.MapNodeSetModifier{
		Config: &models.MapNodeSetModifierConfig{
			ValueToSet: "{{.vars.prefix}}-{{.NodeValue}}-{{.vars.suffix}}",
			MapNodeModifierConfig: &models.MapNodeModifierConfig{
				Vars: map[string]interface{}{
					"prefix": "START",
					"suffix": "END",
				},
			},
			TemplateConfigHolder: &models.TemplateConfigHolder{
				TemplateConfig: templates.TemplateConfig{
					LeftDelim:  "{{",
					RightDelim: "}}",
				},
			},
		},
	}
	
	input := "benchmark"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modifier.ModifyNode(input)
	}
}

func ExampleMapNodeSetModifier_ModifyNode() {
	// Create a set modifier to transform values
	modifier := &map_node_set_modifier.MapNodeSetModifier{
		Config: &models.MapNodeSetModifierConfig{
			ValueToSet: "Processed: {{.NodeValue}} (Type: {{.NodeValueType}})",
			TemplateConfigHolder: &models.TemplateConfigHolder{
				TemplateConfig: templates.TemplateConfig{
					LeftDelim:  "{{",
					RightDelim: "}}",
				},
			},
		},
	}
	
	result := modifier.ModifyNode("input-data")
	// Result will be: "Processed: input-data (Type: string)"
	println(result)
}