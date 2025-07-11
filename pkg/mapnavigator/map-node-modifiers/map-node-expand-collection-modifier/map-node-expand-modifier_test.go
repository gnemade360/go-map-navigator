package map_node_expand_collection_modifier_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	map_nav_models "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-expand-collection-modifier"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/templates"
)

func TestMapNodeExpandModifier_ModifyNode(t *testing.T) {
	tests := []struct {
		name            string
		config          *models.MapNodeExpandModifierConfig
		templateContext map[string]interface{}
		input           interface{}
		expected        interface{}
	}{
		{
			name: "disabled modifier returns node unchanged",
			config: &models.MapNodeExpandModifierConfig{
				MapNodeModifierConfig: &models.MapNodeModifierConfig{
					Disabled: true,
				},
			},
			input:    []interface{}{"item1", "item2"},
			expected: []interface{}{"item1", "item2"},
		},
		{
			name:     "nil input returns nil",
			config:   &models.MapNodeExpandModifierConfig{},
			input:    nil,
			expected: nil,
		},
		{
			name:     "non-array input returns unchanged",
			config:   &models.MapNodeExpandModifierConfig{},
			input:    "not-an-array",
			expected: "not-an-array",
		},
		{
			name: "expand array with simple items",
			config: &models.MapNodeExpandModifierConfig{
				RepeatOn: []interface{}{
					map[string]interface{}{"name": "item3"},
					map[string]interface{}{"name": "item4"},
				},
				RepeatTemplate: models.RepeatTemplate{
					Template:     "name: ~(item.name)~",
					LeftLimiter:  "~(",
					RightLimiter: ")~",
				},
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{},
				},
			},
			input: []interface{}{"existing1", "existing2"},
			expected: []interface{}{
				"existing1",
				"existing2",
				map[string]interface{}{"name": "item3"},
				map[string]interface{}{"name": "item4"},
			},
		},
		{
			name: "expand with custom item field name",
			config: &models.MapNodeExpandModifierConfig{
				ItemField: "current",
				RepeatOn: []interface{}{
					map[string]interface{}{"id": 1},
				},
				RepeatTemplate: models.RepeatTemplate{
					Template:     "id: ~(current.id)~",
					LeftLimiter:  "~(",
					RightLimiter: ")~",
				},
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{},
				},
			},
			input:    []interface{}{},
			expected: []interface{}{map[string]interface{}{"id": 1}},
		},
		{
			name: "expand with template context",
			config: &models.MapNodeExpandModifierConfig{
				RepeatOn: []interface{}{
					map[string]interface{}{"name": "user1"},
				},
				RepeatTemplate: models.RepeatTemplate{
					Template:     "name: ~(item.name)~\nprefix: ~(prefix)~",
					LeftLimiter:  "~(",
					RightLimiter: ")~",
				},
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{},
				},
			},
			templateContext: map[string]interface{}{
				"prefix": "USER",
			},
			input: []interface{}{},
			expected: []interface{}{
				map[string]interface{}{
					"name":   "user1",
					"prefix": "USER",
				},
			},
		},
		{
			name: "expand with map template",
			config: &models.MapNodeExpandModifierConfig{
				RepeatOn: []interface{}{"value1", "value2"},
				RepeatTemplate: models.RepeatTemplate{
					Template: map[string]interface{}{
						"field": "~(item)~",
					},
					LeftLimiter:  "~(",
					RightLimiter: ")~",
				},
				TemplateConfigHolder: &models.TemplateConfigHolder{
					TemplateConfig: templates.TemplateConfig{},
				},
			},
			input:    []interface{}{},
			expected: []interface{}{map[string]interface{}{"field": "value1"}, map[string]interface{}{"field": "value2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifier := &map_node_expand_collection_modifier.MapNodeExpandModifier{
				Config:          tt.config,
				TemplateContext: tt.templateContext,
			}
			
			result := modifier.ModifyNode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapNodeExpandModifier_GetTemplatisedObject(t *testing.T) {
	tests := []struct {
		name     string
		modifier *map_node_expand_collection_modifier.MapNodeExpandModifier
		item     interface{}
		validate func(t *testing.T, result interface{})
	}{
		{
			name: "simple string templating",
			modifier: &map_node_expand_collection_modifier.MapNodeExpandModifier{
				Config: &models.MapNodeExpandModifierConfig{
					RepeatTemplate: models.RepeatTemplate{
						Template:     "value: ~(item)~",
						LeftLimiter:  "~(",
						RightLimiter: ")~",
					},
					TemplateConfigHolder: &models.TemplateConfigHolder{
						TemplateConfig: templates.TemplateConfig{},
					},
				},
				TemplateContext: map[string]interface{}{},
			},
			item: "test",
			validate: func(t *testing.T, result interface{}) {
				m, ok := result.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "test", m["value"])
			},
		},
		{
			name: "complex object templating",
			modifier: &map_node_expand_collection_modifier.MapNodeExpandModifier{
				Config: &models.MapNodeExpandModifierConfig{
					ItemField: "user",
					RepeatTemplate: models.RepeatTemplate{
						Template: `name: ~(user.name)~
email: ~(user.email)~
active: true`,
						LeftLimiter:  "~(",
						RightLimiter: ")~",
					},
					TemplateConfigHolder: &models.TemplateConfigHolder{
						TemplateConfig: templates.TemplateConfig{},
					},
				},
				TemplateContext: map[string]interface{}{},
			},
			item: map[string]interface{}{
				"name":  "John",
				"email": "john@example.com",
			},
			validate: func(t *testing.T, result interface{}) {
				m, ok := result.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "John", m["name"])
				assert.Equal(t, "john@example.com", m["email"])
				assert.Equal(t, true, m["active"])
			},
		},
		{
			name: "default delimiters",
			modifier: &map_node_expand_collection_modifier.MapNodeExpandModifier{
				Config: &models.MapNodeExpandModifierConfig{
					RepeatTemplate: models.RepeatTemplate{
						Template: "field: ~(item)~",
					},
					TemplateConfigHolder: &models.TemplateConfigHolder{
						TemplateConfig: templates.TemplateConfig{},
					},
				},
				TemplateContext: map[string]interface{}{},
			},
			item: "value",
			validate: func(t *testing.T, result interface{}) {
				m, ok := result.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "value", m["field"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.modifier.GetTemplatisedObject(tt.item)
			tt.validate(t, result)
		})
	}
}

func TestMapNodeExpandModifier_GetNewTemplateContextForItem(t *testing.T) {
	modifier := &map_node_expand_collection_modifier.MapNodeExpandModifier{
		Config: &models.MapNodeExpandModifierConfig{
			ItemField: "current",
		},
		TemplateContext: map[string]interface{}{
			"global": "value",
			"prefix": "PREFIX",
		},
	}
	
	item := map[string]interface{}{"id": 123}
	result := modifier.GetNewTemplateContextForItem(item)
	
	ctx, ok := result.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, item, ctx["current"])
	assert.Equal(t, "value", ctx["global"])
	assert.Equal(t, "PREFIX", ctx["prefix"])
}

func TestMapNodeExpandModifier_GetItemTemplate(t *testing.T) {
	tests := []struct {
		name             string
		modifier         *map_node_expand_collection_modifier.MapNodeExpandModifier
		context          interface{}
		expectedTemplate string
		expectedType     string
	}{
		{
			name: "string template",
			modifier: &map_node_expand_collection_modifier.MapNodeExpandModifier{
				Config: &models.MapNodeExpandModifierConfig{
					RepeatTemplate: models.RepeatTemplate{
						Template: "name: {{.item.name}}",
					},
				},
			},
			expectedTemplate: "name: {{.item.name}}",
			expectedType:     "nil",
		},
		{
			name: "map template",
			modifier: &map_node_expand_collection_modifier.MapNodeExpandModifier{
				Config: &models.MapNodeExpandModifierConfig{
					RepeatTemplate: models.RepeatTemplate{
						Template: map[string]interface{}{
							"field1": "value1",
							"field2": "value2",
						},
					},
				},
			},
			expectedTemplate: "field1: value1\nfield2: value2\n",
			expectedType:     "map",
		},
		{
			name: "empty template",
			modifier: &map_node_expand_collection_modifier.MapNodeExpandModifier{
				Config: &models.MapNodeExpandModifierConfig{
					RepeatTemplate: models.RepeatTemplate{
						Template: nil,
					},
				},
			},
			expectedTemplate: "",
			expectedType:     "nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, templateType := tt.modifier.GetItemTemplate(tt.context)
			assert.Equal(t, tt.expectedTemplate, template)
			if tt.expectedType == "nil" {
				assert.Nil(t, templateType)
			} else if tt.expectedType == "map" {
				assert.NotNil(t, templateType)
			}
		})
	}
}

func TestNewMapNodeExpandModifier(t *testing.T) {
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
			name: "with MapNodeModifierConfig and MapNodeExpandModifierConfig",
			config: &models.MapNodeModifierConfig{
				Options: &models.MapNodeExpandModifierConfig{
					ItemField: "item",
					RepeatOn:  []interface{}{"a", "b", "c"},
				},
			},
			mapContext: map[string]interface{}{"key": "value"},
			tmplConfig: tmplConfig,
			expectNil:  false,
			validateFn: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				assert.NotNil(t, modifier)
				em, ok := modifier.(*map_node_expand_collection_modifier.MapNodeExpandModifier)
				assert.True(t, ok)
				assert.Equal(t, "item", em.Config.ItemField)
				assert.Len(t, em.Config.RepeatOn, 3)
			},
		},
		{
			name: "with map options",
			config: &models.MapNodeModifierConfig{
				Options: map[string]interface{}{
					"itemField": "current",
					"repeatOn":  []interface{}{1, 2, 3},
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
			modifier := map_node_expand_collection_modifier.NewMapNodeExpandModifier(
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

func BenchmarkMapNodeExpandModifier_ModifyNode(b *testing.B) {
	modifier := &map_node_expand_collection_modifier.MapNodeExpandModifier{
		Config: &models.MapNodeExpandModifierConfig{
			RepeatOn: []interface{}{
				map[string]interface{}{"id": 1},
				map[string]interface{}{"id": 2},
				map[string]interface{}{"id": 3},
			},
			RepeatTemplate: models.RepeatTemplate{
				Template:     "id: ~(item.id)~",
				LeftLimiter:  "~(",
				RightLimiter: ")~",
			},
			TemplateConfigHolder: &models.TemplateConfigHolder{
				TemplateConfig: templates.TemplateConfig{},
			},
		},
	}
	
	input := []interface{}{}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = modifier.ModifyNode(input)
	}
}

func ExampleMapNodeExpandModifier_ModifyNode() {
	// Create an expand modifier to add items to a collection
	modifier := &map_node_expand_collection_modifier.MapNodeExpandModifier{
		Config: &models.MapNodeExpandModifierConfig{
			RepeatOn: []interface{}{
				map[string]interface{}{"name": "Alice", "role": "admin"},
				map[string]interface{}{"name": "Bob", "role": "user"},
			},
			RepeatTemplate: models.RepeatTemplate{
				Template: `name: ~(item.name)~
role: ~(item.role)~
created: ~(timestamp)~`,
				LeftLimiter:  "~(",
				RightLimiter: ")~",
			},
			TemplateConfigHolder: &models.TemplateConfigHolder{
				TemplateConfig: templates.TemplateConfig{},
			},
		},
		TemplateContext: map[string]interface{}{
			"timestamp": "2024-01-01",
		},
	}
	
	existingUsers := []interface{}{
		map[string]interface{}{"name": "Charlie", "role": "guest"},
	}
	
	result := modifier.ModifyNode(existingUsers)
	// Result will contain Charlie plus the two new expanded users
	println(result)
}