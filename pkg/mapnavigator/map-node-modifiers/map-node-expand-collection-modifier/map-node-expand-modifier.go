package map_node_expand_collection_modifier

import (
	"fmt"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-node-action-repository/models"
	"reflect"
	"strings"

	map_nav_models "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/templates"
	SaTypes "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/types"
	"gopkg.in/yaml.v3"
)

type MapNodeExpandModifier struct {
	TemplateContext    map[string]interface{}
	Config             *models.MapNodeExpandModifierConfig
	LeftTemplateDelim  string
	RightTemplateDelim string
	ChildNodeAction    map_nav_models.MapNodeModifier
}

func (m *MapNodeExpandModifier) ModifyNode(i interface{}) interface{} {
	if m.TemplateContext == nil {
		m.TemplateContext = map[string]interface{}{}
	}

	if m.Config != nil && m.Config.Disabled {
		return i
	}
	if i == nil {
		return i
	}

	// get array
	repeatOnCollection := m.Config.RepeatOn
	//get expansion path object/slice
	if arr, isArray := i.([]interface{}); isArray && arr != nil {
		for _, itm := range repeatOnCollection {
			var newItem = m.GetTemplatisedObject(itm)
			if i != nil {
				arr = append(arr, newItem)
			}
		}
		return arr
	}
	// get template

	// append to the slice

	return i
}

func (m *MapNodeExpandModifier) GetTemplatisedObject(itm interface{}) interface{} {

	newTemplateContext := m.GetNewTemplateContextForItem(itm)
	itmTemplate, itmType := m.GetItemTemplate(newTemplateContext)
	if len(itmTemplate) == 0 {
		return nil
	}
	if itmType == nil {
		itmType = SaTypes.Map
	}
	var newItm interface{}
	rlimiter := m.Config.RepeatTemplate.RightLimiter
	llimiter := m.Config.RepeatTemplate.LeftLimiter
	if len(rlimiter) == 0 {
		rlimiter = ")~"
	}
	if len(llimiter) == 0 {
		llimiter = "~("
	}
	tmplConfig := templates.TemplateConfig{
		RightDelim:    rlimiter,
		LeftDelim:     llimiter,
		TemplateFuncs: m.Config.TemplateConfig.TemplateFuncs,
	}

	if nis, err := templates.Interpolate(itmTemplate, newTemplateContext, tmplConfig); err == nil {
		byt := []byte(nis)
		mp := map[string]interface{}{}
		if er := yaml.Unmarshal(byt, &mp); er != nil {
			fmt.Printf("\nyaml unmarshalling error while getting new repeated item-: %v\n", er)
		}
		return mp
	} else {
		fmt.Printf("\nError while templating the repeating item-err: %v\n", err)
	}
	return newItm

}

func (m *MapNodeExpandModifier) GetNewTemplateContextForItem(itm interface{}) interface{} {
	nc := map[string]interface{}{}
	itemFieldName := m.Config.ItemField
	if len(itemFieldName) == 0 {
		itemFieldName = "item"
	}
	nc[itemFieldName] = itm

	for k, v := range m.TemplateContext {
		nc[k] = v
	}
	return nc
}

func (m *MapNodeExpandModifier) GetItemTemplate(context interface{}) (string, reflect.Type) {
	templateContext := context
	if templateContext == nil {
		templateContext = m.TemplateContext
	}
	tmpl := m.Config.RepeatTemplate.Template
	if tTextStr, ok := tmpl.(string); ok {
		return tTextStr, nil
	} else if tMap, isMap := tmpl.(map[string]interface{}); isMap {
		if byts, err := yaml.Marshal(tMap); err == nil {
			return string(byts), SaTypes.Map
		}
	}

	return "", nil
}

func getAccessor(tTextStr string) []string {
	acs := strings.Split(tTextStr, ".")
	accessors := []string{}
	for _, ac := range acs {
		if len(ac) > 0 {
			accessors = append(accessors, ac)
		}
	}
	return accessors
}

func NewMapNodeExpandModifier(config interface{}, mapContext map[string]interface{}, templateConfig templates.TemplateConfig) map_nav_models.MapNodeModifier {
	m := &MapNodeExpandModifier{
		TemplateContext: mapContext,
	}
	if config != nil {
		if c, ok := config.(*models.MapNodeModifierConfig); ok {
			if cc, okk := c.Options.(*models.MapNodeExpandModifierConfig); okk {
				m.Config = cc
			} else if ccMap, isMap := c.Options.(map[string]interface{}); isMap {
				optionsBytes, er := yaml.Marshal(ccMap)
				if er != nil {
					return nil
				}
				opts := &models.MapNodeExpandModifierConfig{
					MapNodeModifierConfig: &models.MapNodeModifierConfig{},
					TemplateConfigHolder:  &models.TemplateConfigHolder{TemplateConfig: templateConfig},
				}
				yaml.Unmarshal(optionsBytes, opts)
				opts.MapNodeModifierConfig = c
				m.Config = opts
			}
		} else {
			mp := config
			optionsBytes, er := yaml.Marshal(mp)
			if er != nil {
				return nil
			}
			opts := &models.MapNodeExpandModifierConfig{
				MapNodeModifierConfig: &models.MapNodeModifierConfig{},
			}
			yaml.Unmarshal(optionsBytes, opts)
			opts.MapNodeModifierConfig = c
			m.Config = opts
		}
	}

	return m
}
