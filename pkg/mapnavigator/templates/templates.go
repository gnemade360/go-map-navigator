package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

// ArrDelimiter is the delimiter used for array strings
const ArrDelimiter = "|||"

// TemplateConfig holds configuration for templates
type TemplateConfig struct {
	LeftDelim     string
	RightDelim    string
	TemplateFuncs template.FuncMap
}

// NewTemplateConfig creates a new template configuration
func NewTemplateConfig() *TemplateConfig {
	return &TemplateConfig{
		LeftDelim:  "{{",
		RightDelim: "}}",
	}
}

// New creates a new template with the given name and config
func New(name string, config TemplateConfig) (*template.Template, error) {
	t := template.New(name)
	if config.LeftDelim != "" && config.RightDelim != "" {
		t = t.Delims(config.LeftDelim, config.RightDelim)
	}
	if config.TemplateFuncs != nil {
		t = t.Funcs(config.TemplateFuncs)
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

// Interpolate interpolates a template string with the given context
func Interpolate(templateStr string, context map[string]interface{}, config TemplateConfig) (string, error) {
	if templateStr == "" {
		return "", nil
	}
	
	tmpl, err := New("interpolate", config)
	if err != nil {
		return "", err
	}
	
	tmpl, err = tmpl.Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	
	return Execute(tmpl, context)
}

// Parse parses text as a template body for t
func Parse(tmpl *template.Template, text string) (*template.Template, error) {
	return tmpl.Parse(text)
}