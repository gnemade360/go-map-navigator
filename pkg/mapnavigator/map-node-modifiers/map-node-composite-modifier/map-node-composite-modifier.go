package map_node_composite_modifier

import (
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	"reflect"

	map_nav_models "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/templates"
	"gopkg.in/yaml.v3"
)

type MapNodeCompositeModifier struct {
	TemplateContext map[string]interface{}
	Config          *models.MapNodeCompositeModifierConfig
	NodeActions     []map_nav_models.MapNodeModifier
}

func (m *MapNodeCompositeModifier) ModifyNode(node interface{}) interface{} {

	//if conf, ok := m.TemplateContext["Config"]; ok {
	//    //do something here
	//    if appConf, valid := conf.(*action_config.AppConfig); valid && appConf != nil {
	//        name, _ := appConf.UniqueNameProvider.ProvideUniqueName(m.TemplateContext["Result"], appConf.Resource)
	//        fmt.Printf("\nInside ModifyNode of Composte-name: %v\n", name)
	//        fmt.Printf("\ncaption-: %v\n", m.Config.Caption)
	//
	//    }
	//}

	if m.Config != nil && m.Config.MapNodeModifierConfig != nil && m.Config.MapNodeModifierConfig.Disabled {
		return node
	}
	if m.TemplateContext == nil {
		m.TemplateContext = make(map[string]interface{})
	}
	m.TemplateContext["NodeValue"] = node
	if node != nil {
		m.TemplateContext["NodeValueType"] = reflect.TypeOf(node).String()
	}
	if m.Config != nil && m.Config.MapNodeModifierConfig != nil {
		getVars(m.Config.MapNodeModifierConfig.Vars, m.TemplateContext, m.Config.TemplateConfig)
	}
	if m.NodeActions != nil && len(m.NodeActions) > 0 {
		for _, nodeModifier := range m.NodeActions {
			node = nodeModifier.ModifyNode(node)
		}
	}
	return node
}

func NewMapNodeCompositeModifier(config interface{}, mapContext map[string]interface{}, tmplConfig templates.TemplateConfig) map_nav_models.MapNodeModifier {
	m := &MapNodeCompositeModifier{
		NodeActions:     []map_nav_models.MapNodeModifier{},
		TemplateContext: mapContext,
	}

	if config != nil {
		if c, ok := config.(*models.MapNodeModifierConfig); ok {
			if cc, okk := c.Options.(*models.MapNodeCompositeModifierConfig); okk {
				m.Config = cc
				m.Config.Caption = c.Caption
			} else if ccMap, isMap := c.Options.(map[string]interface{}); isMap {

				optionsBytes, er := yaml.Marshal(ccMap)
				if er != nil {
					return nil
				}

				copts := getOptionsFromByts(optionsBytes, tmplConfig)
				m.Config = copts
				m.Config.Caption = c.Caption
			}
		} else {
			mp := config
			optionsBytes, er := yaml.Marshal(mp)
			if er != nil {
				return nil
			}
			opts := getOptionsFromByts(optionsBytes, tmplConfig)

			m.Config = opts
		}
	}

	return m

}

func getOptionsFromByts(optionsBytes []byte, templateConfig templates.TemplateConfig) *models.MapNodeCompositeModifierConfig {

	sopts := &models.MapNodeCompositeModifierConfig{}
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
