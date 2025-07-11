package map_node_action_repository

import (
	"fmt"
	mnmodels "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	map_node_composite_modifier "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-composite-modifier"
	map_node_conditional_modifier "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-conditional-modifier"
	map_node_delete_modifier "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-delete-modifier"
	map_node_expand_collection_modifier "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-expand-collection-modifier"
	map_node_replace_modifier "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-replace-modifier"
	map_node_set_modifier "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-modifiers/map-node-set-modifier"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/templates"
	"sync"

	"gopkg.in/yaml.v3"
)

type MapNodeModifierRepository struct {
	*sync.RWMutex
	MapNodeModifierCollection MapNodeModifierCollectionType
	TemplateConfig            templates.TemplateConfig
}
type MapNodeModifierFunc func(config interface{}, mapContext map[string]interface{}, templateConfig templates.TemplateConfig) mnmodels.MapNodeModifier
type MapNodeModifierCollectionType map[string]MapNodeModifierFunc

var Instance *MapNodeModifierRepository = nil //NewMapNodeModifierRepository()

func GetMapNodeModifierRepository(templateConfig templates.TemplateConfig) *MapNodeModifierRepository {

	instance := NewMapNodeModifierRepository(templateConfig)
	return instance
}

func NewMapNodeModifierRepository(templateConfig templates.TemplateConfig, params ...interface{}) *MapNodeModifierRepository {

	// templateConfig.AddFuncMap(map[string]interface{}{})

	instance := &MapNodeModifierRepository{
		RWMutex:                   &sync.RWMutex{},
		MapNodeModifierCollection: make(MapNodeModifierCollectionType, 0),
	}

	instance.TemplateConfig = templateConfig
	instance.Register("set", map_node_set_modifier.NewMapNodeSetModifier)
	instance.Register("replace", map_node_replace_modifier.NewMapNodeReplaceModifier)
	instance.Register("conditional", func(config interface{}, mapContext map[string]interface{}, templateConfig templates.TemplateConfig) mnmodels.MapNodeModifier {
		cond := map_node_conditional_modifier.NewMapNodeConditionalModifier(config, mapContext, templateConfig)
		c := cond.(*map_node_conditional_modifier.MapNodeConditionalModifier)
		if c.Config != nil && c.Config.IfTrue != nil {
			c.IfTrue, _ = instance.Get(c.Config.IfTrue.NodeAction, c.Config.IfTrue, mapContext)

		}
		if c.Config != nil && c.Config.IfFalse != nil {
			c.IfFalse, _ = instance.Get(c.Config.IfFalse.NodeAction, c.Config.IfFalse, mapContext)
		}

		return cond
	})
	instance.Register("composite", func(config interface{}, mapContext map[string]interface{}, templateConfig templates.TemplateConfig) mnmodels.MapNodeModifier {
		compositeModifier := map_node_composite_modifier.NewMapNodeCompositeModifier(config, mapContext, templateConfig)
		if c, ok := compositeModifier.(*map_node_composite_modifier.MapNodeCompositeModifier); ok {
			for _, childNodeAction := range c.Config.NodeActions {
				fmt.Printf("\nchildNodeAction-childNodeAction: %v\n", childNodeAction)
				childMod, _ := instance.Get(childNodeAction.NodeAction, childNodeAction, mapContext)
				c.NodeActions = append(c.NodeActions, childMod)
				//
				//instance.get
			}
		}

		return compositeModifier
	})
	instance.Register("expand", func(config interface{}, mapContext map[string]interface{}, templateConfig templates.TemplateConfig) mnmodels.MapNodeModifier {
		newmodifier := map_node_expand_collection_modifier.NewMapNodeExpandModifier(config, mapContext, templateConfig)

		if c, ok := newmodifier.(*map_node_expand_collection_modifier.MapNodeExpandModifier); ok {

			var _ = c

		}

		return newmodifier
	})

	instance.Register("delete", map_node_delete_modifier.NewMapNodeDeleteModifier)

	return instance
}

func (mmr *MapNodeModifierRepository) Register(e string, ser MapNodeModifierFunc) *MapNodeModifierRepository {
	mmr.Lock()
	defer func() {
		mmr.Unlock()
	}()
	if mmr.MapNodeModifierCollection == nil {
		mmr.MapNodeModifierCollection = make(MapNodeModifierCollectionType, 0)
	}
	mmr.MapNodeModifierCollection[e] = ser
	return mmr
}

func (mmr *MapNodeModifierRepository) Get(e string, config interface{}, mapContext map[string]interface{}) (mnmodels.MapNodeModifier, bool) {
	mmr.RLock()
	defer mmr.RUnlock()
	var conf interface{}
	if newConfigByts, err := yaml.Marshal(config); err == nil {
		mp := &models.MapNodeModifierConfig{}
		yaml.Unmarshal(newConfigByts, &mp)
		conf = mp
	} else {
		conf = config
	}
	toReturn, exists := mmr.MapNodeModifierCollection[e]
	modifier := toReturn(conf, mapContext, mmr.TemplateConfig)
	return modifier, exists
}
