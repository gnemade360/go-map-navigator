package map_node_replace_modifier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	map_nav_models "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-replace-modifier"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/templates"
)

func TestMapNodeReplaceModifier_ModifyNode(t *testing.T) {
	tests := []struct {
		name            string
		config          *models.MapNodeReplaceModifierConfig
		templateContext map[string]interface{}
		input           interface{}
		expected        interface{}
	}{
		{
			name: "disabled modifier returns node unchanged",
			config: &models.MapNodeReplaceModifierConfig{
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					Disabled: true,
				},
			},
			input:    "test-value",
			expected: "test-value",
		},
		{
			name: "nil input with find not null returns nil",
			config: &models.MapNodeReplaceModifierConfig{
				Find: "test",
			},
			input:    nil,
			expected: nil,
		},
		{
			name: "nil input with find null converts to string null",
			config: &models.MapNodeReplaceModifierConfig{
				Find: "null",
			},
			input:    nil,
			expected: "null",
		},
		{
			name: "simple string replacement",
			config: &models.MapNodeReplaceModifierConfig{
				Find:        "old",
				ReplaceWith: "new",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{},
				},
			},
			input:    "old-value-old",
			expected: "new-value-new",
		},
		{
			name: "replace first occurrence only",
			config: &models.MapNodeReplaceModifierConfig{
				Find:             "test",
				ReplaceWith:      "replaced",
				ReplaceFirstOnly: true,
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{},
				},
			},
			input:    "test-test-test",
			expected: "replaced-test-test",
		},
		{
			name: "replace with template interpolation",
			config: &models.MapNodeReplaceModifierConfig{
				Find:        "{{.SearchFor}}",
				ReplaceWith: "{{.ReplaceValue}}",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{
						LeftDelim:  "{{",
						RightDelim: "}}",
					},
				},
			},
			templateContext: map[string]interface{}{
				"SearchFor":    "hello",
				"ReplaceValue": "world",
			},
			input:    "hello there",
			expected: "world there",
		},
		{
			name: "replace with vars interpolation",
			config: &models.MapNodeReplaceModifierConfig{
				Find:        "PLACEHOLDER",
				ReplaceWith: "{{.vars.prefix}}-{{.NodeValue}}",
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					Vars: map[string]interface{}{
						"prefix": "PREFIX",
					},
				},
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{
						LeftDelim:  "{{",
						RightDelim: "}}",
					},
				},
			},
			input:    "PLACEHOLDER",
			expected: "PREFIX-PLACEHOLDER",
		},
		{
			name: "empty find string returns unchanged",
			config: &models.MapNodeReplaceModifierConfig{
				Find:        "",
				ReplaceWith: "new",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{},
				},
			},
			input:    "unchanged",
			expected: "unchanged",
		},
		{
			name: "replace in map with selector",
			config: &models.MapNodeReplaceModifierConfig{
				Find:        "old",
				ReplaceWith: "new",
				Selector:    "field",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{},
				},
			},
			input: map[string]interface{}{
				"field": "old-value",
				"other": "old-keep",
			},
			expected: map[string]interface{}{
				"field": "new-value",
				"other": "old-keep",
			},
		},
		{
			name: "non-string input returns unchanged",
			config: &models.MapNodeReplaceModifierConfig{
				Find:        "test",
				ReplaceWith: "new",
			},
			input:    42,
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifier := &map_node_replace_modifier.MapNodeReplaceModifier{
				Config:          tt.config,
				TemplateContext: tt.templateContext,
			}
			
			result := modifier.ModifyNode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewMapNodeReplaceModifier(t *testing.T) {
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
			name: "with MapNodeModifierConfig and MapNodeReplaceModifierConfig",
			config: &models.MapNodeModifierConfig{
				Caption: "Replace Test",
				Options: &models.MapNodeReplaceModifierConfig{
					Find:        "search",
					ReplaceWith: "replace",
				},
			},
			mapContext: map[string]interface{}{"key": "value"},
			tmplConfig: tmplConfig,
			expectNil:  false,
			validateFn: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				assert.NotNil(t, modifier)
				rm, ok := modifier.(*map_node_replace_modifier.MapNodeReplaceModifier)
				assert.True(t, ok)
				assert.Equal(t, "search", rm.Config.Find)
				assert.Equal(t, "replace", rm.Config.ReplaceWith)
				assert.Equal(t, "Replace Test", rm.Config.Caption)
			},
		},
		{
			name: "with map options",
			config: &models.MapNodeModifierConfig{
				Options: map[string]interface{}{
					"find":             "pattern",
					"replaceWith":      "replacement",
					"replaceFirstOnly": true,
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
			modifier := map_node_replace_modifier.NewMapNodeReplaceModifier(
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

func TestReplaceModifier_ComplexScenarios(t *testing.T) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}

	t.Run("multiple replacements with template context", func(t *testing.T) {
		modifier := &map_node_replace_modifier.MapNodeReplaceModifier{
			Config: &models.MapNodeReplaceModifierConfig{
				Find:        "{{.SearchPattern}}",
				ReplaceWith: "{{.vars.prefix}}-{{.findResult}}-{{.vars.suffix}}",
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					Vars: map[string]interface{}{
						"prefix": "START",
						"suffix": "END",
					},
				},
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
			},
			TemplateContext: map[string]interface{}{
				"SearchPattern": "MARKER",
			},
		}
		
		result := modifier.ModifyNode("This has MARKER and another MARKER here")
		expected := "This has START-MARKER-END and another START-MARKER-END here"
		assert.Equal(t, expected, result)
	})

	t.Run("replace with array delimiter", func(t *testing.T) {
		modifier := &map_node_replace_modifier.MapNodeReplaceModifier{
			Config: &models.MapNodeReplaceModifierConfig{
				Find:        "a" + templates.ArrDelimiter + "b",
				ReplaceWith: "replaced",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
			},
		}
		
		result := modifier.ModifyNode("a b c")
		assert.Equal(t, "replaced replaced c", result)
	})

	t.Run("nested map modification", func(t *testing.T) {
		modifier := &map_node_replace_modifier.MapNodeReplaceModifier{
			Config: &models.MapNodeReplaceModifierConfig{
				Find:        "secret",
				ReplaceWith: "***",
				Selector:    "user.password",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
			},
		}
		
		data := map[string]interface{}{
			"user": map[string]interface{}{
				"name":     "john",
				"password": "secret123",
			},
		}
		
		result := modifier.ModifyNode(data)
		resultMap := result.(map[string]interface{})
		userMap := resultMap["user"].(map[string]interface{})
		assert.Equal(t, "***123", userMap["password"])
		assert.Equal(t, "john", userMap["name"])
	})
}

func TestReplaceModifier_EdgeCases(t *testing.T) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}

	t.Run("replace with special characters", func(t *testing.T) {
		modifier := &map_node_replace_modifier.MapNodeReplaceModifier{
			Config: &models.MapNodeReplaceModifierConfig{
				Find:        "\\n",
				ReplaceWith: " ",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
			},
		}
		
		result := modifier.ModifyNode("line1\\nline2\\nline3")
		assert.Equal(t, "line1 line2 line3", result)
	})

	t.Run("replace with empty string", func(t *testing.T) {
		modifier := &map_node_replace_modifier.MapNodeReplaceModifier{
			Config: &models.MapNodeReplaceModifierConfig{
				Find:        "remove",
				ReplaceWith: "",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
			},
		}
		
		result := modifier.ModifyNode("remove-this-remove")
		assert.Equal(t, "-this-", result)
	})

	t.Run("case sensitive replacement", func(t *testing.T) {
		modifier := &map_node_replace_modifier.MapNodeReplaceModifier{
			Config: &models.MapNodeReplaceModifierConfig{
				Find:        "Test",
				ReplaceWith: "Exam",
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
			},
		}
		
		result := modifier.ModifyNode("Test test TEST")
		assert.Equal(t, "Exam test TEST", result)
	})
}

func BenchmarkMapNodeReplaceModifier_ModifyNode(b *testing.B) {
	modifier := &map_node_replace_modifier.MapNodeReplaceModifier{
		Config: &models.MapNodeReplaceModifierConfig{
			Find:        "old",
			ReplaceWith: "new",
			TemplateConfigHolder: &models.TemplateConfigHolder{
				TemplateConfig: templates.TemplateConfig{},
			},
		},
	}
	
	input := "old-value-old-text-old"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modifier.ModifyNode(input)
	}
}

func ExampleMapNodeReplaceModifier_ModifyNode() {
	// Create a replace modifier to mask sensitive data
	modifier := &map_node_replace_modifier.MapNodeReplaceModifier{
		Config: &models.MapNodeReplaceModifierConfig{
			Find:        "[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{4}",
			ReplaceWith: "****-****-****-****",
			TemplateConfigHolder: &models.TemplateConfigHolder{
				TemplateConfig: templates.TemplateConfig{},
			},
		},
	}
	
	data := "Credit card: 1234-5678-9012-3456"
	result := modifier.ModifyNode(data)
	// Result will be: "Credit card: ****-****-****-****"
	println(result)
}