package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

// TemplateConfig holds configuration for templates
type TemplateConfig struct {
	Delims []string
}

// NewTemplateConfig creates a new template configuration
func NewTemplateConfig() *TemplateConfig {
	return &TemplateConfig{
		Delims: []string{"{{", "}}"},
	}
}

// New creates a new template with the given name and config
func New(name string, config TemplateConfig) (*template.Template, error) {
	t := template.New(name)
	if len(config.Delims) >= 2 {
		t = t.Delims(config.Delims[0], config.Delims[1])
	}
	return t, nil
}

// Execute executes a template with the given context
func Execute(t *template.Template, context interface{}) (string, error) {
	var buf bytes.Buffer
	err := t.Execute(&buf, context)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

// ApplyTemplateFuncs applies custom functions to a template
func ApplyTemplateFuncs(t *template.Template, funcMap template.FuncMap) {
	if funcMap != nil {
		t.Funcs(funcMap)
	}
}