package map_node_composite_modifier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	map_nav_models "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-composite-modifier"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/templates"
)

func TestMapNodeCompositeModifier_ModifyNode(t *testing.T) {
	tests := []struct {
		name            string
		config          *models.MapNodeCompositeModifierConfig
		templateContext map[string]interface{}
		nodeActions     []map_nav_models.MapNodeModifier
		input           interface{}
		expected        interface{}
	}{
		{
			name: "disabled modifier returns node unchanged",
			config: &models.MapNodeCompositeModifierConfig{
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					Disabled: true,
				},
			},
			input:    "test-value",
			expected: "test-value",
		},
		{
			name:     "nil config processes node",
			config:   nil,
			input:    "test-value",
			expected: "test-value",
		},
		{
			name: "single action modifier",
			config: &models.MapNodeCompositeModifierConfig{},
			nodeActions: []map_nav_models.MapNodeModifier{
				map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
					if s, ok := v.(string); ok {
						return s + "-modified"
					}
					return v
				}),
			},
			input:    "test",
			expected: "test-modified",
		},
		{
			name: "multiple action modifiers in sequence",
			config: &models.MapNodeCompositeModifierConfig{},
			nodeActions: []map_nav_models.MapNodeModifier{
				map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
					if s, ok := v.(string); ok {
						return s + "-first"
					}
					return v
				}),
				map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
					if s, ok := v.(string); ok {
						return s + "-second"
					}
					return v
				}),
			},
			input:    "test",
			expected: "test-first-second",
		},
		{
			name:     "nil input handling",
			config:   &models.MapNodeCompositeModifierConfig{},
			input:    nil,
			expected: nil,
		},
		{
			name: "with template context",
			config: &models.MapNodeCompositeModifierConfig{
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					Vars: map[string]interface{}{
						"prefix": "{{.Prefix}}",
					},
				},
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{
						LeftDelim:  "{{",
						RightDelim: "}}",
					},
				},
			},
			templateContext: map[string]interface{}{
				"Prefix": "custom",
			},
			nodeActions: []map_nav_models.MapNodeModifier{
				map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
					return v
				}),
			},
			input:    "test",
			expected: "test",
		},
		{
			name: "empty node actions",
			config: &models.MapNodeCompositeModifierConfig{},
			nodeActions: []map_nav_models.MapNodeModifier{},
			input:    "unchanged",
			expected: "unchanged",
		},
		{
			name: "complex object modification",
			config: &models.MapNodeCompositeModifierConfig{},
			nodeActions: []map_nav_models.MapNodeModifier{
				map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
					if m, ok := v.(map[string]interface{}); ok {
						m["modified"] = true
						return m
					}
					return v
				}),
			},
			input:    map[string]interface{}{"key": "value"},
			expected: map[string]interface{}{"key": "value", "modified": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifier := &map_node_composite_modifier.MapNodeCompositeModifier{
				Config:          tt.config,
				TemplateContext: tt.templateContext,
				NodeActions:     tt.nodeActions,
			}
			
			result := modifier.ModifyNode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewMapNodeCompositeModifier(t *testing.T) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}

	tests := []struct {
		name        string
		config      interface{}
		mapContext  map[string]interface{}
		tmplConfig  templates.TemplateConfig
		expectNil   bool
		validateFn  func(t *testing.T, modifier map_nav_models.MapNodeModifier)
	}{
		{
			name: "with MapNodeModifierConfig",
			config: &models.MapNodeModifierConfig{
				Caption: "test-caption",
				Options: &models.MapNodeCompositeModifierConfig{
					MapNodeModifierConfig: &models.MapNodeModifierConfig{
						Disabled: false,
					},
				},
			},
			mapContext: map[string]interface{}{"key": "value"},
			tmplConfig: tmplConfig,
			expectNil:  false,
			validateFn: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				assert.NotNil(t, modifier)
				cm, ok := modifier.(*map_node_composite_modifier.MapNodeCompositeModifier)
				assert.True(t, ok)
				assert.Equal(t, "test-caption", cm.Config.Caption)
			},
		},
		{
			name: "with map options",
			config: &models.MapNodeModifierConfig{
				Caption: "map-caption",
				Options: map[string]interface{}{
					"disabled": false,
					"vars": map[string]interface{}{
						"key": "value",
					},
				},
			},
			mapContext: map[string]interface{}{},
			tmplConfig: tmplConfig,
			expectNil:  false,
			validateFn: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				assert.NotNil(t, modifier)
				cm, ok := modifier.(*map_node_composite_modifier.MapNodeCompositeModifier)
				assert.True(t, ok)
				assert.Equal(t, "map-caption", cm.Config.Caption)
			},
		},
		{
			name: "with direct map config",
			config: map[string]interface{}{
				"caption": "direct-caption",
				"disabled": false,
			},
			mapContext: nil,
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
			modifier := map_node_composite_modifier.NewMapNodeCompositeModifier(
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

func TestCompositeModifier_TemplateContextHandling(t *testing.T) {
	modifier := &map_node_composite_modifier.MapNodeCompositeModifier{
		Config: &models.MapNodeCompositeModifierConfig{},
		TemplateContext: nil,
	}
	
	// Test that nil template context is initialized
	result := modifier.ModifyNode("test")
	assert.NotNil(t, modifier.TemplateContext)
	assert.Equal(t, "test", modifier.TemplateContext["NodeValue"])
	assert.Equal(t, "string", modifier.TemplateContext["NodeValueType"])
	assert.Equal(t, "test", result)
}

func TestCompositeModifier_WithVars(t *testing.T) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}
	
	modifier := &map_node_composite_modifier.MapNodeCompositeModifier{
		Config: &models.MapNodeCompositeModifierConfig{
			MapNodeModifierConfig: &models.MapNodeModifierConfig{
				Vars: map[string]interface{}{
					"greeting": "Hello {{.Name}}",
					"static":   "static-value",
					"number":   42,
				},
			},
			TemplateConfigHolder: &models.TemplateConfigHolder{
				TemplateConfig: tmplConfig,
			},
		},
		TemplateContext: map[string]interface{}{
			"Name": "World",
		},
	}
	
	result := modifier.ModifyNode("test")
	
	// Check that vars were processed
	vars, ok := modifier.TemplateContext["vars"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Hello World", vars["greeting"])
	assert.Equal(t, "static-value", vars["static"])
	assert.Equal(t, 42, vars["number"])
	assert.Equal(t, "test", result)
}

func BenchmarkMapNodeCompositeModifier_SingleAction(b *testing.B) {
	modifier := &map_node_composite_modifier.MapNodeCompositeModifier{
		Config: &models.MapNodeCompositeModifierConfig{},
		NodeActions: []map_nav_models.MapNodeModifier{
			map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if s, ok := v.(string); ok {
					return s + "-modified"
				}
				return v
			}),
		},
	}
	
	input := "benchmark-test"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modifier.ModifyNode(input)
	}
}

func BenchmarkMapNodeCompositeModifier_MultipleActions(b *testing.B) {
	modifier := &map_node_composite_modifier.MapNodeCompositeModifier{
		Config: &models.MapNodeCompositeModifierConfig{},
		NodeActions: []map_nav_models.MapNodeModifier{
			map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				return v
			}),
			map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				return v
			}),
			map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				return v
			}),
		},
	}
	
	input := map[string]interface{}{"key": "value"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modifier.ModifyNode(input)
	}
}

func ExampleMapNodeCompositeModifier_ModifyNode() {
	// Create a composite modifier with multiple actions
	modifier := &map_node_composite_modifier.MapNodeCompositeModifier{
		Config: &models.MapNodeCompositeModifierConfig{},
		NodeActions: []map_nav_models.MapNodeModifier{
			map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				// First modifier: add prefix
				if s, ok := v.(string); ok {
					return "prefix-" + s
				}
				return v
			}),
			map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				// Second modifier: add suffix
				if s, ok := v.(string); ok {
					return s + "-suffix"
				}
				return v
			}),
		},
	}
	
	result := modifier.ModifyNode("test")
	println(result) // "prefix-test-suffix"
}