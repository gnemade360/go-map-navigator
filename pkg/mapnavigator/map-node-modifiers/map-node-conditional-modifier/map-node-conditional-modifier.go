package map_node_conditional_modifier

import (
	"fmt"
	"reflect"

	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"

	map_nav_models "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/conditions"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/templates"
	"gopkg.in/yaml.v3"
)

type MapNodeConditionalModifier struct {
	TemplateContext map[string]interface{}
	IfTrue          map_nav_models.MapNodeModifier
	IfFalse         map_nav_models.MapNodeModifier
	Config          *models.MapNodeConditionalModifierConfig
	*models.TemplateConfigHolder
	ConditionsExecutor *conditions.SaConditionsExecutor
}

func (m *MapNodeConditionalModifier) ModifyNode(node interface{}) interface{} {
	if m.Config == nil {
		return node
	}

	if m.Config != nil && m.Config.MapNodeModifierConfig != nil && m.Config.MapNodeModifierConfig.Disabled {
		return node
	}

	if conditionResult, err := m.ExecuteCondition(node); err != nil {
		fmt.Printf("\nError while executing condition-m.Condition: %v\n", m.Config.Condition)
	} else {
		if conditionResult && m.IfTrue != nil {
			//fmt.Printf("\nCondition result in the replacement rule is true-m.Config: %v\n", m.Config)
			//fmt.Printf("\nExecuting IFTRUE node modification-node: %v\n", node)
			return m.IfTrue.ModifyNode(node)
		} else if !conditionResult && m.IfFalse != nil {
			//fmt.Printf("\nCondition result in the replacement rule is false-m.Config: %v\n", m.Config)
			//fmt.Printf("\nExecuting IFFALSE node modification-node: %v\n", node)
			return m.IfFalse.ModifyNode(node)
		}
	}
	return node
}

func (m *MapNodeConditionalModifier) ExecuteCondition(node interface{}) (bool, error) {

	if m.TemplateContext == nil {
		m.TemplateContext = make(map[string]interface{})
	}
	m.TemplateContext["NodeValue"] = node
	if node != nil {
		m.TemplateContext["NodeValueType"] = reflect.TypeOf(node).String()
	}
	if m.Config != nil && m.Config.MapNodeModifierConfig != nil {
		tmplConfig := templates.TemplateConfig{}
		if m.TemplateConfigHolder != nil {
			tmplConfig = m.TemplateConfigHolder.TemplateConfig
		}
		getVars(m.Config.MapNodeModifierConfig.Vars, m.TemplateContext, tmplConfig)
	}

	//conditions := m.ConditionsExecutor.Conditions()

	if m.ConditionsExecutor == nil {
		fmt.Printf("\nConditionsExecutor is nil-: %v\n", m.ConditionsExecutor)
	}

	return m.ConditionsExecutor.ExecuteConditions(m.TemplateContext)

	//if finalStr, er := templates.Interpolate(m.Config.Condition, m.TemplateContext, m.TemplateConfig); er == nil {
	//    finalStr = strings.TrimSpace(finalStr)
	//    return len(finalStr) > 0 && finalStr != "false" && finalStr != "0", nil
	//} else {
	//    fmt.Println(er)
	//    return false, er
	//}
}

func NewMapNodeConditionalModifier(config interface{}, mapContext map[string]interface{}, tmplConfig templates.TemplateConfig) map_nav_models.MapNodeModifier {
	m := &MapNodeConditionalModifier{
		TemplateContext:      mapContext,
		TemplateConfigHolder: &models.TemplateConfigHolder{TemplateConfig: tmplConfig},
	}

	if config != nil {
		if c, ok := config.(*models.MapNodeModifierConfig); ok {
			if cc, okk := c.Options.(*models.MapNodeConditionalModifierConfig); okk {
				m.Config = cc
				m.Config.TemplateConfigHolder = &models.TemplateConfigHolder{TemplateConfig: tmplConfig}
				m.ConditionsExecutor = conditions.NewSaConditionsExecutor(
					conditions.WithConditions(m.Config.Conditions),
					conditions.WithTemplateConfig(tmplConfig),
				)
			} else if ccMap, isMap := c.Options.(map[string]interface{}); isMap {

				conditionsObj := conditions.NewSaConditions(ccMap, "conditions")
				delete(ccMap, "conditions")
				optionsBytes, er := yaml.Marshal(ccMap)
				if er != nil {
					return nil
				}

				copts := getOptionsFromByts(optionsBytes, tmplConfig)
				if ccMap != nil {
					ccMap["conditions"] = conditionsObj
				}
				// setting conditions now

				m.ConditionsExecutor = conditions.NewSaConditionsExecutor(
					conditions.WithConditions(conditionsObj),
					conditions.WithTemplateConfig(tmplConfig),
				)
				if len(copts.Condition) > 0 {
					m.ConditionsExecutor.SetCondition(copts.Condition)
				}
				copts.Conditions = *conditionsObj
				m.Config = copts
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
	if m.ConditionsExecutor == nil {
		fmt.Printf("\nm.ConditionsExecutor is nil-: %v\n", m.ConditionsExecutor)
	}
	return m
}

func getOptionsFromByts(optionsBytes []byte, templateConfig templates.TemplateConfig) *models.MapNodeConditionalModifierConfig {
	sopts := &models.MapNodeConditionalModifierConfig{}
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
			} else {
				fmt.Printf("\ntemplate interpolate in getVars function of conditional modifier er-er: %v\n", er)
			}
		} else {
			toReturn[key] = value
		}
	}
	context["vars"] = toReturn
	return toReturn
}
