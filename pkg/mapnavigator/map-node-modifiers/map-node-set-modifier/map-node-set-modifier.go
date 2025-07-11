package map_node_set_modifier

import (
	"fmt"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	"reflect"
	"strings"

	generic_value "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/generic-value"
	map_navigator "github.com/gnemade360/go-map-navigator/pkg/mapnavigator"
	map_nav_models "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/templates"
	"gopkg.in/yaml.v3"
)

type MapNodeSetModifier struct {
	TemplateContext map[string]interface{}
	Config          *models.MapNodeSetModifierConfig
}

func (m *MapNodeSetModifier) ModifyNode(i interface{}) interface{} {

	if m.TemplateContext == nil {
		m.TemplateContext = map[string]interface{}{}
	}
	if i == nil {
		return i
	}

	if m.Config != nil && m.Config.MapNodeModifierConfig != nil && m.Config.MapNodeModifierConfig.Disabled {
		return i
	}
	typ := reflect.TypeOf(i)

	kind := typ.Kind()
	if nodestr, isString := i.(string); isString {

		m.TemplateContext["NodeValue"] = nodestr
		m.TemplateContext["NodeValueType"] = reflect.TypeOf(nodestr).String()
		if m.Config != nil {
			if m.Config.MapNodeModifierConfig != nil && m.Config.TemplateConfigHolder != nil {
				getVars(m.Config.MapNodeModifierConfig.Vars, m.TemplateContext, m.Config.TemplateConfigHolder.TemplateConfig)
			}

			tmplConfig := templates.TemplateConfig{}
			if m.Config.TemplateConfigHolder != nil {
				tmplConfig = m.Config.TemplateConfigHolder.TemplateConfig
			}
			if finalStr, er := templates.Interpolate(m.Config.ValueToSet, m.TemplateContext, tmplConfig); er == nil {
				return finalStr
			} else {
				fmt.Printf("\nwarn-er: %v\n", er)
			}
		}
	} else if kind != reflect.Map && kind != reflect.Slice {
		str := (&generic_value.SedulousTypeConverter{}).ConvertToString(i)

		m.TemplateContext["NodeValue"] = str
		m.TemplateContext["NodeValueType"] = typ.String()
		if m.Config != nil {
			if m.Config.MapNodeModifierConfig != nil && m.Config.TemplateConfigHolder != nil {
				getVars(m.Config.MapNodeModifierConfig.Vars, m.TemplateContext, m.Config.TemplateConfigHolder.TemplateConfig)
			}

			tmplConfig := templates.TemplateConfig{}
			if m.Config.TemplateConfigHolder != nil {
				tmplConfig = m.Config.TemplateConfigHolder.TemplateConfig
			}
			//m.TemplateContext["NodeValue"] = i
			if finalStr, er := templates.Interpolate(m.Config.ValueToSet, m.TemplateContext, tmplConfig); er == nil {
			//v := reflect.ValueOf(finalStr)
			//v.Convert(typ)
			//vi := v.Interface()
				return finalStr
			} else {
				fmt.Printf("\nwarn-er: %v\n", er)
			}
		}
	} else {
		mn := map_navigator.NewMapNavigator(m)
		// TODO: we need to take this from options
		var ks []string
		if m.Config != nil {
			if m.Config.MapNodeModifierConfig != nil {
				mn.CreateProperty = m.Config.MapNodeModifierConfig.CreatePropertyIfAbsent
			}
			ks = strings.Split(m.Config.Selector, ".")
		} else {
			ks = []string{}
		}
		mn.VisitNode(i, ks...)
		return i
	}

	return i
}

func NewMapNodeSetModifier(config interface{}, mapContext map[string]interface{}, templateConfig templates.TemplateConfig) map_nav_models.MapNodeModifier {
	m := &MapNodeSetModifier{
		TemplateContext: mapContext,
	}
	if config != nil {
		if c, ok := config.(*models.MapNodeModifierConfig); ok {
			if cc, okk := c.Options.(*models.MapNodeSetModifierConfig); okk {
				m.Config = cc
			} else {
				optionsBytes, er := yaml.Marshal(c.Options)
				if er != nil {
					return nil
				}

				sopts := getOptionsFromByts(optionsBytes, templateConfig)
				m.Config = sopts

			}
		} else if mapConfig, isMap := config.(map[string]interface{}); isMap {
			optionsBytes, er := yaml.Marshal(mapConfig)
			if er != nil {
				return nil
			}
			opts := getOptionsFromByts(optionsBytes, templateConfig)
			m.Config = opts
		}
	}
	return m
}

func getOptionsFromByts(optionsBytes []byte, templateConfig templates.TemplateConfig) *models.MapNodeSetModifierConfig {
	sopts := &models.MapNodeSetModifierConfig{
		Selector:             "",
		ValueToSet:           "",
		TemplateConfigHolder: nil,
	}
	opts := &models.MapNodeModifierConfig{}
	yaml.Unmarshal(optionsBytes, sopts)
	yaml.Unmarshal(optionsBytes, opts)
	sopts.MapNodeModifierConfig = opts
	sopts.TemplateConfigHolder = &models.TemplateConfigHolder{
		TemplateConfig: templateConfig,
	}
	sopts.MapNodeModifierConfig = opts
	return sopts
}

func getVars(obj, context map[string]interface{}, templateConfig templates.TemplateConfig) map[string]interface{} {
	if obj == nil {
		return nil
	}
	toReturn := map[string]interface{}{}
	if context == nil {
		context = map[string]interface{}{}
	}
	for key, value := range obj {
		if valStr, isString := value.(string); isString {

			if str, er := templates.Interpolate(valStr, context, templateConfig); er == nil {
				toReturn[key] = str
			}
		} else {
			toReturn[key] = value
		}
	}
	context["vars"] = toReturn
	return toReturn
}
