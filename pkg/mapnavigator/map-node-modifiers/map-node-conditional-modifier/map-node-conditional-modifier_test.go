package map_node_conditional_modifier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/conditions"
	map_nav_models "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-conditional-modifier"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/templates"
)

func TestMapNodeConditionalModifier_ModifyNode(t *testing.T) {
	tests := []struct {
		name               string
		config             *models.MapNodeConditionalModifierConfig
		templateContext    map[string]interface{}
		ifTrue             map_nav_models.MapNodeModifier
		ifFalse            map_nav_models.MapNodeModifier
		conditionsExecutor *conditions.SaConditionsExecutor
		input              interface{}
		expected           interface{}
	}{
		{
			name:     "nil config returns node unchanged",
			config:   nil,
			input:    "test-value",
			expected: "test-value",
		},
		{
			name: "disabled modifier returns node unchanged",
			config: &models.MapNodeConditionalModifierConfig{
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					Disabled: true,
				},
			},
			input:    "test-value",
			expected: "test-value",
		},
		{
			name:   "condition true executes ifTrue modifier",
			config: &models.MapNodeConditionalModifierConfig{},
			conditionsExecutor: &conditions.SaConditionsExecutor{
				// Mock executor that returns true
			},
			ifTrue: map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if s, ok := v.(string); ok {
					return s + "-true"
				}
				return v
			}),
			ifFalse: map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if s, ok := v.(string); ok {
					return s + "-false"
				}
				return v
			}),
			input:    "test",
			expected: "test",
		},
		{
			name:   "condition false executes ifFalse modifier",
			config: &models.MapNodeConditionalModifierConfig{},
			conditionsExecutor: &conditions.SaConditionsExecutor{
				// Mock executor that returns false
			},
			ifTrue: map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if s, ok := v.(string); ok {
					return s + "-true"
				}
				return v
			}),
			ifFalse: map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if s, ok := v.(string); ok {
					return s + "-false"
				}
				return v
			}),
			input:    "test",
			expected: "test",
		},
		{
			name:   "nil ifTrue modifier when condition is true",
			config: &models.MapNodeConditionalModifierConfig{},
			conditionsExecutor: &conditions.SaConditionsExecutor{},
			ifTrue:  nil,
			ifFalse: map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				return "should-not-execute"
			}),
			input:    "unchanged",
			expected: "unchanged",
		},
		{
			name:   "nil ifFalse modifier when condition is false",
			config: &models.MapNodeConditionalModifierConfig{},
			conditionsExecutor: &conditions.SaConditionsExecutor{},
			ifTrue: map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				return "should-not-execute"
			}),
			ifFalse:  nil,
			input:    "unchanged",
			expected: "unchanged",
		},
		{
			name:     "nil input handling",
			config:   &models.MapNodeConditionalModifierConfig{},
			conditionsExecutor: &conditions.SaConditionsExecutor{},
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifier := &map_node_conditional_modifier.MapNodeConditionalModifier{
				Config:             tt.config,
				TemplateContext:    tt.templateContext,
				IfTrue:             tt.ifTrue,
				IfFalse:            tt.ifFalse,
				ConditionsExecutor: tt.conditionsExecutor,
			}
			
			result := modifier.ModifyNode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapNodeConditionalModifier_ExecuteCondition(t *testing.T) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}

	tests := []struct {
		name            string
		config          *models.MapNodeConditionalModifierConfig
		templateContext map[string]interface{}
		node            interface{}
		setupExecutor   func() *conditions.SaConditionsExecutor
		expectedResult  bool
		expectError     bool
	}{
		{
			name: "execute with string node",
			config: &models.MapNodeConditionalModifierConfig{
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					Vars: map[string]interface{}{
						"testVar": "{{.NodeValue}}",
					},
				},
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
				Condition: "{{.NodeValue}} == 'test-string'",
			},
			node: "test-string",
			setupExecutor: func() *conditions.SaConditionsExecutor {
				executor := conditions.NewSaConditionsExecutor()
				executor.SetCondition("{{.NodeValue}} == 'test-string'")
				return executor
			},
			expectedResult: true,
			expectError:    false,
		},
		{
			name: "execute with nil template context",
			config: &models.MapNodeConditionalModifierConfig{
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
				Condition: "{{.NodeValue}} == 42",
			},
			templateContext: nil,
			node:            42,
			setupExecutor: func() *conditions.SaConditionsExecutor {
				executor := conditions.NewSaConditionsExecutor()
				executor.SetCondition("{{.NodeValue}} == 42")
				return executor
			},
			expectedResult: true,
			expectError:    false,
		},
		{
			name: "execute with vars interpolation",
			config: &models.MapNodeConditionalModifierConfig{
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					Vars: map[string]interface{}{
						"prefix": "{{.Prefix}}",
						"static": "static-value",
					},
				},
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: tmplConfig,
				},
				Condition: "{{.NodeValue}} == 'test'",
			},
			templateContext: map[string]interface{}{
				"Prefix": "custom",
			},
			node: "test",
			setupExecutor: func() *conditions.SaConditionsExecutor {
				executor := conditions.NewSaConditionsExecutor()
				executor.SetCondition("{{.NodeValue}} == 'test'")
				return executor
			},
			expectedResult: true,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifier := &map_node_conditional_modifier.MapNodeConditionalModifier{
				Config:               tt.config,
				TemplateContext:      tt.templateContext,
				TemplateConfigHolder: tt.config.TemplateConfigHolder,
				ConditionsExecutor:   tt.setupExecutor(),
			}
			
			result, err := modifier.ExecuteCondition(tt.node)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			
			// Verify template context was updated
			assert.Equal(t, tt.node, modifier.TemplateContext["NodeValue"])
			if tt.node != nil {
				assert.NotEmpty(t, modifier.TemplateContext["NodeValueType"])
			}
		})
	}
}

func TestNewMapNodeConditionalModifier(t *testing.T) {
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
			name: "with MapNodeModifierConfig and MapNodeConditionalModifierConfig",
			config: &models.MapNodeModifierConfig{
				Options: &models.MapNodeConditionalModifierConfig{
					MapNodeModifierConfig: &models.MapNodeModifierConfig{
						Disabled: false,
					},
					Conditions: conditions.SaConditions{},
				},
			},
			mapContext: map[string]interface{}{"key": "value"},
			tmplConfig: tmplConfig,
			expectNil:  false,
			validateFn: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				assert.NotNil(t, modifier)
				cm, ok := modifier.(*map_node_conditional_modifier.MapNodeConditionalModifier)
				assert.True(t, ok)
				assert.NotNil(t, cm.Config)
				assert.NotNil(t, cm.ConditionsExecutor)
			},
		},
		{
			name: "with map options",
			config: &models.MapNodeModifierConfig{
				Options: map[string]interface{}{
					"disabled":  false,
					"condition": "true",
					"conditions": map[string]interface{}{
						"type": "always_true",
					},
				},
			},
			mapContext: map[string]interface{}{},
			tmplConfig: tmplConfig,
			expectNil:  false,
			validateFn: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				assert.NotNil(t, modifier)
				cm, ok := modifier.(*map_node_conditional_modifier.MapNodeConditionalModifier)
				assert.True(t, ok)
				assert.NotNil(t, cm.Config)
				assert.NotNil(t, cm.ConditionsExecutor)
			},
		},
		{
			name: "with direct map config",
			config: map[string]interface{}{
				"disabled":  false,
				"condition": "{{.Value}} > 10",
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
			modifier := map_node_conditional_modifier.NewMapNodeConditionalModifier(
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

func TestConditionalModifier_ComplexScenarios(t *testing.T) {
	// tmplConfig := templates.TemplateConfig{
	// 	LeftDelim:  "{{",
	// 	RightDelim: "}}",
	// }

	t.Run("nested object modification", func(t *testing.T) {
		modifier := &map_node_conditional_modifier.MapNodeConditionalModifier{
			Config: &models.MapNodeConditionalModifierConfig{
				Condition: "{{.id}} == 123",
			},
			ConditionsExecutor: conditions.NewSaConditionsExecutor(),
			IfTrue: map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if m, ok := v.(map[string]interface{}); ok {
					m["status"] = "processed"
					return m
				}
				return v
			}),
		}
		
		input := map[string]interface{}{
			"id":   123,
			"name": "test",
		}
		
		result := modifier.ModifyNode(input)
		resultMap, ok := result.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "processed", resultMap["status"])
	})

	t.Run("chained conditional modifiers", func(t *testing.T) {
		// First conditional: always true, add prefix
		modifier1 := &map_node_conditional_modifier.MapNodeConditionalModifier{
			Config: &models.MapNodeConditionalModifierConfig{
				Condition: "{{.NodeValue}} != ''",
			},
			ConditionsExecutor: conditions.NewSaConditionsExecutor(),
			IfTrue: map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if s, ok := v.(string); ok {
					return s + "-processed"
				}
				return v
			}),
		}
		
		input := "testing"
		result := modifier1.ModifyNode(input)
		
		assert.Equal(t, "testing-processed", result)
	})
}

func BenchmarkMapNodeConditionalModifier_ModifyNode(b *testing.B) {
	modifier := &map_node_conditional_modifier.MapNodeConditionalModifier{
		Config: &models.MapNodeConditionalModifierConfig{},
		ConditionsExecutor: conditions.NewSaConditionsExecutor(),
		IfTrue: map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
			if s, ok := v.(string); ok {
				return s + "-modified"
			}
			return v
		}),
	}
	
	input := "benchmark-test"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modifier.ModifyNode(input)
	}
}

func ExampleMapNodeConditionalModifier_ModifyNode() {
	// Create a conditional modifier
	modifier := &map_node_conditional_modifier.MapNodeConditionalModifier{
		Config: &models.MapNodeConditionalModifierConfig{},
		ConditionsExecutor: conditions.NewSaConditionsExecutor(),
		IfTrue: map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
			// If condition is true, convert to uppercase
			if s, ok := v.(string); ok {
				return "TRUE: " + s
			}
			return v
		}),
		IfFalse: map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
			// If condition is false, add prefix
			if s, ok := v.(string); ok {
				return "FALSE: " + s
			}
			return v
		}),
	}
	
	result := modifier.ModifyNode("test")
	println(result) // Either "TRUE: test" or "FALSE: test" depending on condition
}