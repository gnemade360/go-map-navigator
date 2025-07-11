package conditions

import (
	"fmt"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/templates"
)

// SaConditions represents a condition that can be evaluated
type SaConditions struct {
	conditions []interface{}
}

// NewSaConditions creates a new SaConditions instance
func NewSaConditions(configMap map[string]interface{}, key string) *SaConditions {
	if configMap == nil {
		return &SaConditions{conditions: []interface{}{}}
	}
	
	conditions, ok := configMap[key].([]interface{})
	if !ok {
		// Try to get it as a single condition
		if cond, exists := configMap[key]; exists {
			conditions = []interface{}{cond}
		} else {
			conditions = []interface{}{}
		}
	}
	
	return &SaConditions{
		conditions: conditions,
	}
}

// Evaluate evaluates the conditions against the provided data
func (c *SaConditions) Evaluate(data interface{}) (bool, error) {
	// For now, return true as a stub implementation
	// This would normally evaluate the conditions
	return true, nil
}

// SaConditionsExecutor executes conditions with a template configuration
type SaConditionsExecutor struct {
	conditions     *SaConditions
	templateConfig templates.TemplateConfig
}

// NewSaConditionsExecutor creates a new executor with options
func NewSaConditionsExecutor(opts ...ExecutorOption) *SaConditionsExecutor {
	executor := &SaConditionsExecutor{
		templateConfig: templates.TemplateConfig{
			LeftDelim:  "{{",
			RightDelim: "}}",
		},
	}
	
	for _, opt := range opts {
		opt(executor)
	}
	
	return executor
}

// ExecutorOption is a function that configures an executor
type ExecutorOption func(*SaConditionsExecutor)

// WithConditions sets the conditions for the executor
func WithConditions(conditions interface{}) ExecutorOption {
	return func(e *SaConditionsExecutor) {
		switch c := conditions.(type) {
		case *SaConditions:
			e.conditions = c
		case SaConditions:
			e.conditions = &c
		default:
			e.conditions = &SaConditions{conditions: []interface{}{conditions}}
		}
	}
}

// WithTemplateConfig sets the template configuration for the executor
func WithTemplateConfig(config templates.TemplateConfig) ExecutorOption {
	return func(e *SaConditionsExecutor) {
		e.templateConfig = config
	}
}

// ExecuteConditions evaluates the conditions with the given context
func (e *SaConditionsExecutor) ExecuteConditions(context map[string]interface{}) (bool, error) {
	if e.conditions == nil {
		return false, fmt.Errorf("no conditions set")
	}
	return e.conditions.Evaluate(context)
}

// SetCondition sets a single condition string
func (e *SaConditionsExecutor) SetCondition(condition string) {
	if e.conditions == nil {
		e.conditions = &SaConditions{conditions: []interface{}{}}
	}
	e.conditions.conditions = []interface{}{condition}
}