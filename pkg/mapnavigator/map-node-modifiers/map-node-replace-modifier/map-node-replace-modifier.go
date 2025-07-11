package map_node_replace_modifier

import (
	"fmt"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	"reflect"
	"strings"

	map_navigator "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator"
	map_nav_models "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/templates"
	"gopkg.in/yaml.v3"
)

type MapNodeReplaceModifier struct {
	TemplateContext map[string]interface{}
	Config          *models.MapNodeReplaceModifierConfig
}

func (m *MapNodeReplaceModifier) ModifyNode(i interface{}) interface{} {
	if m.Config != nil && m.Config.Disabled {
		return i
	}
	//if conf, ok := m.TemplateContext["Config"]; ok {
	//    //do something here
	//    if appConf, valid := conf.(*action_config.AppConfig); valid && appConf != nil {
	//        name, _ := appConf.UniqueNameProvider.ProvideUniqueName(m.TemplateContext["Result"], appConf.Resource)
	//        fmt.Printf("\nInside ModifyNode of replace-name: %v\n", name)
	//        fmt.Printf("\ncaption-: %v\n", m.Config.Caption)
	//    }
	//}
	if m.TemplateContext == nil {
		m.TemplateContext = map[string]interface{}{}
	}

	if i == nil {
		if m.Config.Find != "null" {
			return i
		} else {
			i = "null"
		}
	}
	kind := reflect.TypeOf(i).Kind()
	if nodestr, isString := i.(string); isString {

		fmt.Printf("\nReplace rule: replacing text-nodestr: %v\n", nodestr)

		m.TemplateContext["NodeValue"] = nodestr
		m.TemplateContext["NodeValueType"] = reflect.TypeOf(i).String()
		getVars(m.Config.Vars, m.TemplateContext, m.Config.TemplateConfig)

		findinterpolated, err := templates.Interpolate(m.Config.Find, m.TemplateContext, m.Config.TemplateConfig)
		if len(findinterpolated) == 0 || findinterpolated == "[]" || err != nil {
			return i
		}
		findArr := strings.Split(findinterpolated, templates.ArrDelimiter)
		for _, f := range findArr {
			m.TemplateContext["findResult"] = f
			replaceWith := ""
			fmt.Printf("\nReplace Rule-findinterpolated: %v\n", findinterpolated)
			if tmpl, err := templates.New("m.Config.ReplaceWith", m.Config.TemplateConfig); err == nil {
				tmpl = tmpl.Funcs(map[string]interface{}{})

				tmpl, err = templates.Parse(tmpl, m.Config.ReplaceWith)
				if err != nil {
					fmt.Printf("Error: %v", err)
				}
				if ff, er := templates.Execute(tmpl, m.TemplateContext); er == nil {
					replaceWith = ff
					fmt.Printf("\nReplace Rule-replaceWith : %v\n", findinterpolated)

				} else {
					fmt.Printf("\nReplace Rule Error while executing template-er: %v\n", er)
				}
			} else {
				return i
			}
			if m.Config.ReplaceFirstOnly {
				i = strings.Replace(i.(string), f, replaceWith, 1)
			} else {
				fmt.Printf("\nReplace Rule: Final replace-: \n")
				i = strings.ReplaceAll(i.(string), f, replaceWith)
			}
		}
	} else if kind == reflect.Map || kind == reflect.Slice {
		mn := map_navigator.NewMapNavigator(m)
		if len(strings.Trim(m.Config.Selector, " ")) == 0 {
			m.Config.Selector = "-"
		}
		ks := strings.Split(m.Config.Selector, ".")
		_, er := mn.VisitNode(i, ks...)
		if er != nil {
			fmt.Printf("\nWarn: %v\n", er)
		}
		return i
	}

	return i
}

func NewMapNodeReplaceModifier(config interface{}, mapContext map[string]interface{}, templateConfig templates.TemplateConfig) map_nav_models.MapNodeModifier {

	m := &MapNodeReplaceModifier{
		TemplateContext: mapContext,
	}
	if config != nil {
		if c, ok := config.(*models.MapNodeModifierConfig); ok {
			if cc, okk := c.Options.(*models.MapNodeReplaceModifierConfig); okk {
				m.Config = cc
			} else {
				optionsBytes, er := yaml.Marshal(c.Options)
				if er != nil {
					return nil
				}
				ropts := getOptionsFromByts(optionsBytes, templateConfig)
				m.Config = ropts
				m.Config.Caption = c.Caption
			}
		} else if mapConfig, isMap := c.Options.(map[string]interface{}); isMap {
			optionsBytes, er := yaml.Marshal(mapConfig["options"])
			if er != nil {
				return nil
			}
			opts := getOptionsFromByts(optionsBytes, templateConfig)
			m.Config = opts
			m.Config.Caption = c.Caption
		}
	}
	return m
}

func getOptionsFromByts(optionsBytes []byte, templateConfig templates.TemplateConfig) *models.MapNodeReplaceModifierConfig {
	sopts := &models.MapNodeReplaceModifierConfig{}
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
