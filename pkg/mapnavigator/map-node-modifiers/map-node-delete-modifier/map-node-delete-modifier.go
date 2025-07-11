package map_node_delete_modifier

import (
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	"reflect"
	"strings"

	map_navigator "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator"
	map_nav_models "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/templates"
	"gopkg.in/yaml.v3"
)

type MapNodeDeleteModifier struct {
	TemplateContext map[string]interface{}
	Config          *models.MapNodeDeleteModifierConfig
	*models.TemplateConfigHolder
}

func (m *MapNodeDeleteModifier) ModifyNode(i interface{}) interface{} {

	if m.Config != nil && m.Config.Disabled {
		return i
	}
	if m.TemplateContext == nil {
		m.TemplateContext = map[string]interface{}{}
	}
	if i == nil {
		return i
	}
	// TODO delete the item
	if m.Config.Selector == "-" || m.Config.Selector == "*" {
		return &map_navigator.MapNavigatorDeleted{
			Original: i,
		}
	}
	if deleteFrom, er := GetValue(m.TemplateContext, m.Config.DeleteFrom); er == nil && deleteFrom != nil {
		k := reflect.TypeOf(deleteFrom).Kind()
		switch k {
		case reflect.Slice:
			//return i
			return &map_navigator.MapNavigatorDeleted{
				Original: i,
			}
		}

	}

	if k := reflect.TypeOf(i).Kind(); k == reflect.Map {
		accessors := getAccessors(m.Config.Selector)
		mn := &map_navigator.MapNavigator{
			NodeModifier: map_nav_models.MapNodeModifierFunc(func(v interface{}) interface{} {
				return &map_navigator.MapNavigatorDeleted{Original: v}
			}),
		}
		mn.VisitNode(i, accessors...)
		return i
	}

	return i
}

func getAccessors(accessor string) []string {
	accessors := []string{}
	if len(accessor) > 0 {
		temps := strings.Split(accessor, ".")
		for _, s := range temps {
			if len(s) > 0 {
				accessors = append(accessors, s)
			}
		}
	}
	return accessors
}

func GetValue(base interface{}, accessor string) (interface{}, error) {
	accessors := []string{}
	if len(accessor) > 0 {
		temps := strings.Split(accessor, ".")
		for _, s := range temps {
			if len(s) > 0 {
				accessors = append(accessors, s)
			}
		}

		mn := map_navigator.MapNavigator{
			NodeModifier: nil,
			ReadOnly:     false,
		}

		return mn.VisitNode(base, accessors...)
	}

	return nil, nil
}

func NewMapNodeDeleteModifier(config interface{}, mapContext map[string]interface{}, templateConfig templates.TemplateConfig) map_nav_models.MapNodeModifier {
	m := &MapNodeDeleteModifier{
		TemplateContext: mapContext,
	}
	if config != nil {
		if c, ok := config.(*models.MapNodeModifierConfig); ok {
			if cc, okk := c.Options.(*models.MapNodeDeleteModifierConfig); okk {
				m.Config = cc
			} else {
				optionsBytes, er := yaml.Marshal(c.Options)
				if er != nil {
					return nil
				}
				opts := &models.MapNodeDeleteModifierConfig{
					MapNodeModifierConfig: &models.MapNodeModifierConfig{},
				}
				yaml.Unmarshal(optionsBytes, opts)
				opts.TemplateConfigHolder = &models.TemplateConfigHolder{
					TemplateConfig: templateConfig,
				}
				m.Config = opts

			}
		} else if mapConfig, isMap := c.Options.(map[string]interface{}); isMap {
			optionsBytes, er := yaml.Marshal(mapConfig["options"])
			if er != nil {
				return nil
			}
			opts := &models.MapNodeDeleteModifierConfig{}
			yaml.Unmarshal(optionsBytes, opts)

			m.Config = opts
		}
	}
	return m
}
