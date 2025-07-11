package models

import (
	. "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/conditions"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/templates"
)

type MapNodeModifierConfigHolder struct {
	Config []*MapNodeModifierConfig `json:"config" yaml:"config"`
}
type MapNodeModifierConfig struct {
	Disabled               bool                   `json:"disabled" yaml:"disabled"`
	Path                   string                 `json:"path" yaml:"path"`
	NodeAction             string                 `json:"nodeAction" yaml:"nodeAction"`
	Caption                string                 `json:"caption" yaml:"caption"`
	Order                  int                    `json:"order" yaml:"order"`
	Priority               int                    `json:"priority" yaml:"priority"`
	Options                interface{}            `json:"options,omitempty" yaml:"options,omitempty"`
	IterateOn              string                 `json:"iterateOn" yaml:"iterateOn"`
	IndexItemKey           string                 `json:"indexItemKey" yaml:"indexItemKey"`
	CreatePropertyIfAbsent bool                   `json:"createPropertyIfAbsent" yaml:"createPropertyIfAbsent"`
	Vars                   map[string]interface{} `json:"vars" yaml:"vars"`
}

type MapNodeConditionalModifierConfig struct {
	Condition string                 `json:"condition" yaml:"condition"`
	IfTrue    *MapNodeModifierConfig `json:"ifTrue" yaml:"ifTrue"`
	IfFalse   *MapNodeModifierConfig `json:"ifFalse" yaml:"ifFalse"`
	*TemplateConfigHolder
	Conditions SaConditions `json:"conditions" yaml:"conditions"`
	*MapNodeModifierConfig
}

type MapNodeCompositeModifierConfig struct {
	NodeActions []*MapNodeModifierConfig `json:"nodeActions" yaml:"nodeActions"`
	*TemplateConfigHolder
	*MapNodeModifierConfig
}

type MapNodeSetModifierConfig struct {
	Selector   string `json:"selector" yaml:"selector"`
	ValueToSet string `json:"valueToSet" yaml:"valueToSet"`
	*TemplateConfigHolder
	*MapNodeModifierConfig
}
type MapNodeDeleteModifierConfig struct {
	Selector   string `json:"selector" yaml:"selector" `
	DeleteFrom string `json:"deleteFrom" yaml:"deleteFrom"`
	*TemplateConfigHolder
	*MapNodeModifierConfig
}

type MapNodeReplaceModifierConfig struct {
	Find             string `json:"find" yaml:"find"`
	ReplaceWith      string `json:"replaceWith" yaml:"replaceWith"`
	ReplaceFirstOnly bool   `json:"replaceFirst" yaml:"replaceFirst"`
	Selector         string `json:"selector" yaml:"selector"`
	*TemplateConfigHolder
	*MapNodeModifierConfig
}

type TemplateConfigHolder struct {
	TemplateConfig templates.TemplateConfig
}

type MapNodeExpandModifierConfig struct {
	RepeatLocation string         `json:"repeatLocation" yaml:"repeatLocation"`
	RepeatOn       []interface{}  `json:"repeatOn" yaml:"repeatOn"`
	ItemField      string         `json:"itemField", yaml:"itemField"`
	RepeatTemplate RepeatTemplate `json:"repeatTemplate" yaml:"repeatTemplate"`
	*TemplateConfigHolder
	*MapNodeModifierConfig
}

type RepeatTemplate struct {
	Template     interface{} `json:"template" yaml:"template"`
	LeftLimiter  string      `json:"leftLimiter" yaml:"leftLimiter"`
	RightLimiter string      `json:"rightLimiter" yaml:"rightLimiter"`
}
