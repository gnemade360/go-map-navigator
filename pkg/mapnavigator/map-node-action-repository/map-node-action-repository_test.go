package map_node_action_repository_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	map_nav_models "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-action-repository"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-composite-modifier"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-conditional-modifier"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-delete-modifier"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-expand-collection-modifier"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-replace-modifier"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-set-modifier"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/templates"
)

func TestGetMapNodeModifierRepository(t *testing.T) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}

	repo := map_node_action_repository.GetMapNodeModifierRepository(tmplConfig)
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.RWMutex)
	assert.NotNil(t, repo.MapNodeModifierCollection)
	assert.Equal(t, tmplConfig, repo.TemplateConfig)
}

func TestNewMapNodeModifierRepository(t *testing.T) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}

	repo := map_node_action_repository.NewMapNodeModifierRepository(tmplConfig)
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.RWMutex)
	assert.NotNil(t, repo.MapNodeModifierCollection)
	
	// Verify default modifiers are registered
	expectedModifiers := []string{"set", "replace", "conditional", "composite", "expand", "delete"}
	for _, name := range expectedModifiers {
		_, exists := repo.MapNodeModifierCollection[name]
		assert.True(t, exists, "Modifier '%s' should be registered", name)
	}
}

func TestMapNodeModifierRepository_Register(t *testing.T) {
	tmplConfig := templates.TemplateConfig{}
	repo := map_node_action_repository.NewMapNodeModifierRepository(tmplConfig)
	
	// Custom modifier function
	customModifier := func(config interface{}, mapContext map[string]interface{}, templateConfig templates.TemplateConfig) map_nav_models.MapNodeModifier {
		return map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return "custom-" + v.(string)
		})
	}
	
	// Register custom modifier
	result := repo.Register("custom", customModifier)
	assert.Equal(t, repo, result) // Should return self for chaining
	
	// Verify registration
	_, exists := repo.MapNodeModifierCollection["custom"]
	assert.True(t, exists)
}

func TestMapNodeModifierRepository_Get(t *testing.T) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}
	repo := map_node_action_repository.NewMapNodeModifierRepository(tmplConfig)
	
	tests := []struct {
		name         string
		modifierName string
		config       interface{}
		mapContext   map[string]interface{}
		validateType func(t *testing.T, modifier map_nav_models.MapNodeModifier)
	}{
		{
			name:         "get set modifier",
			modifierName: "set",
			config: &models.MapNodeModifierConfig{
				Options: &models.MapNodeSetModifierConfig{
					ValueToSet: "test",
				},
			},
			mapContext: map[string]interface{}{"key": "value"},
			validateType: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				_, ok := modifier.(*map_node_set_modifier.MapNodeSetModifier)
				assert.True(t, ok)
			},
		},
		{
			name:         "get replace modifier",
			modifierName: "replace",
			config: &models.MapNodeModifierConfig{
				Options: &models.MapNodeReplaceModifierConfig{
					Find:        "old",
					ReplaceWith: "new",
				},
			},
			mapContext: map[string]interface{}{},
			validateType: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				_, ok := modifier.(*map_node_replace_modifier.MapNodeReplaceModifier)
				assert.True(t, ok)
			},
		},
		{
			name:         "get conditional modifier",
			modifierName: "conditional",
			config: &models.MapNodeModifierConfig{
				Options: &models.MapNodeConditionalModifierConfig{},
			},
			mapContext: map[string]interface{}{},
			validateType: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				_, ok := modifier.(*map_node_conditional_modifier.MapNodeConditionalModifier)
				assert.True(t, ok)
			},
		},
		{
			name:         "get composite modifier",
			modifierName: "composite",
			config: &models.MapNodeModifierConfig{
				Options: &models.MapNodeCompositeModifierConfig{},
			},
			mapContext: map[string]interface{}{},
			validateType: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				_, ok := modifier.(*map_node_composite_modifier.MapNodeCompositeModifier)
				assert.True(t, ok)
			},
		},
		{
			name:         "get expand modifier",
			modifierName: "expand",
			config: &models.MapNodeModifierConfig{
				Options: &models.MapNodeExpandModifierConfig{},
			},
			mapContext: map[string]interface{}{},
			validateType: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				_, ok := modifier.(*map_node_expand_collection_modifier.MapNodeExpandModifier)
				assert.True(t, ok)
			},
		},
		{
			name:         "get delete modifier",
			modifierName: "delete",
			config: &models.MapNodeModifierConfig{
				Options: &models.MapNodeDeleteModifierConfig{},
			},
			mapContext: map[string]interface{}{},
			validateType: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				_, ok := modifier.(*map_node_delete_modifier.MapNodeDeleteModifier)
				assert.True(t, ok)
			},
		},
		{
			name:         "get with map config",
			modifierName: "set",
			config: map[string]interface{}{
				"options": map[string]interface{}{
					"valueToSet": "value",
				},
			},
			mapContext: map[string]interface{}{},
			validateType: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				assert.NotNil(t, modifier)
			},
		},
		{
			name:         "get non-existent modifier",
			modifierName: "non-existent",
			config:       nil,
			mapContext:   map[string]interface{}{},
			validateType: func(t *testing.T, modifier map_nav_models.MapNodeModifier) {
				// Should return nil or panic
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.modifierName == "non-existent" {
				// Test non-existent modifier
				defer func() {
					if r := recover(); r != nil {
						// Expected panic for non-existent modifier
					}
				}()
			}
			
			modifier, exists := repo.Get(tt.modifierName, tt.config, tt.mapContext)
			
			if tt.modifierName != "non-existent" {
				assert.True(t, exists)
				assert.NotNil(t, modifier)
				tt.validateType(t, modifier)
			}
		})
	}
}

func TestMapNodeModifierRepository_ConditionalWithChildModifiers(t *testing.T) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}
	repo := map_node_action_repository.NewMapNodeModifierRepository(tmplConfig)
	
	// Create conditional config with child modifiers
	config := &models.MapNodeModifierConfig{
		Options: &models.MapNodeConditionalModifierConfig{
			IfTrue: &models.MapNodeModifierConfig{
				NodeAction: "set",
				Options: &models.MapNodeSetModifierConfig{
					ValueToSet: "true-value",
				},
			},
			IfFalse: &models.MapNodeModifierConfig{
				NodeAction: "replace",
				Options: &models.MapNodeReplaceModifierConfig{
					Find:        "old",
					ReplaceWith: "new",
				},
			},
		},
	}
	
	modifier, exists := repo.Get("conditional", config, map[string]interface{}{})
	assert.True(t, exists)
	
	condModifier, ok := modifier.(*map_node_conditional_modifier.MapNodeConditionalModifier)
	assert.True(t, ok)
	assert.NotNil(t, condModifier.IfTrue)
	assert.NotNil(t, condModifier.IfFalse)
}

func TestMapNodeModifierRepository_CompositeWithChildModifiers(t *testing.T) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}
	repo := map_node_action_repository.NewMapNodeModifierRepository(tmplConfig)
	
	// Create composite config with child modifiers
	config := &models.MapNodeModifierConfig{
		Options: &models.MapNodeCompositeModifierConfig{
			NodeActions: []*models.MapNodeModifierConfig{
				{
					NodeAction: "set",
					Options: &models.MapNodeSetModifierConfig{
						ValueToSet: "value1",
					},
				},
				{
					NodeAction: "replace",
					Options: &models.MapNodeReplaceModifierConfig{
						Find:        "pattern",
						ReplaceWith: "replacement",
					},
				},
			},
		},
	}
	
	modifier, exists := repo.Get("composite", config, map[string]interface{}{})
	assert.True(t, exists)
	
	compModifier, ok := modifier.(*map_node_composite_modifier.MapNodeCompositeModifier)
	assert.True(t, ok)
	assert.Len(t, compModifier.NodeActions, 2)
}

func TestMapNodeModifierRepository_ConcurrentAccess(t *testing.T) {
	tmplConfig := templates.TemplateConfig{}
	repo := map_node_action_repository.NewMapNodeModifierRepository(tmplConfig)
	
	// Test concurrent reads and writes
	var wg sync.WaitGroup
	numGoroutines := 10
	
	// Concurrent registrations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			name := fmt.Sprintf("custom%d", index)
			repo.Register(name, func(config interface{}, mapContext map[string]interface{}, templateConfig templates.TemplateConfig) map_nav_models.MapNodeModifier {
				return map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
					return v
				})
			})
		}(i)
	}
	
	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			config := &models.MapNodeModifierConfig{
				Options: &models.MapNodeSetModifierConfig{
					ValueToSet: "test",
				},
			}
			_, _ = repo.Get("set", config, map[string]interface{}{})
		}()
	}
	
	wg.Wait()
	
	// Verify all custom modifiers were registered
	for i := 0; i < numGoroutines; i++ {
		name := fmt.Sprintf("custom%d", i)
		_, exists := repo.MapNodeModifierCollection[name]
		assert.True(t, exists)
	}
}

func BenchmarkMapNodeModifierRepository_Get(b *testing.B) {
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}
	repo := map_node_action_repository.NewMapNodeModifierRepository(tmplConfig)
	
	config := &models.MapNodeModifierConfig{
		Options: &models.MapNodeSetModifierConfig{
			ValueToSet: "test-value",
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.Get("set", config, map[string]interface{}{})
	}
}

func ExampleMapNodeModifierRepository_Register() {
	// Create a new repository
	tmplConfig := templates.TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}
	repo := map_node_action_repository.NewMapNodeModifierRepository(tmplConfig)
	
	// Register a custom modifier
	repo.Register("uppercase", func(config interface{}, mapContext map[string]interface{}, templateConfig templates.TemplateConfig) map_nav_models.MapNodeModifier {
		return map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
			if s, ok := v.(string); ok {
				return strings.ToUpper(s)
			}
			return v
		})
	})
	
	// Use the custom modifier
	modifier, _ := repo.Get("uppercase", nil, map[string]interface{}{})
	result := modifier.ModifyNode("hello")
	println(result) // Output: HELLO
}