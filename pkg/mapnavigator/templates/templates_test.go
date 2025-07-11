package templates_test

import (
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator/templates"
)

func TestNewTemplateConfig(t *testing.T) {
	config := templates.NewTemplateConfig()
	
	assert.NotNil(t, config)
	assert.Equal(t, "{{", config.LeftDelim)
	assert.Equal(t, "}}", config.RightDelim)
}

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		templateName string
		config      templates.TemplateConfig
		expectError bool
		validate    func(t *testing.T, tmpl *template.Template)
	}{
		{
			name:         "default delimiters",
			templateName: "test-template",
			config: templates.TemplateConfig{
				LeftDelim:  "{{",
				RightDelim: "}}",
			},
			expectError: false,
			validate: func(t *testing.T, tmpl *template.Template) {
				assert.NotNil(t, tmpl)
				assert.Equal(t, "test-template", tmpl.Name())
			},
		},
		{
			name:         "custom delimiters",
			templateName: "custom-delim",
			config: templates.TemplateConfig{
				LeftDelim:  "<<",
				RightDelim: ">>",
			},
			expectError: false,
			validate: func(t *testing.T, tmpl *template.Template) {
				assert.NotNil(t, tmpl)
				assert.Equal(t, "custom-delim", tmpl.Name())
			},
		},
		{
			name:         "empty template name",
			templateName: "",
			config: templates.TemplateConfig{
				LeftDelim:  "{{",
				RightDelim: "}}",
			},
			expectError: false,
			validate: func(t *testing.T, tmpl *template.Template) {
				assert.NotNil(t, tmpl)
				assert.Equal(t, "", tmpl.Name())
			},
		},
		{
			name:         "special characters in delimiters",
			templateName: "special-chars",
			config: templates.TemplateConfig{
				LeftDelim:  "[[",
				RightDelim: "]]",
			},
			expectError: false,
			validate: func(t *testing.T, tmpl *template.Template) {
				assert.NotNil(t, tmpl)
			},
		},
		{
			name:         "single character delimiters",
			templateName: "single-char",
			config: templates.TemplateConfig{
				LeftDelim:  "<",
				RightDelim: ">",
			},
			expectError: false,
			validate: func(t *testing.T, tmpl *template.Template) {
				assert.NotNil(t, tmpl)
			},
		},
		{
			name:         "empty delimiters use defaults",
			templateName: "empty-delim",
			config: templates.TemplateConfig{
				LeftDelim:  "",
				RightDelim: "",
			},
			expectError: false,
			validate: func(t *testing.T, tmpl *template.Template) {
				assert.NotNil(t, tmpl)
				// Template should still be created with default delimiters
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := templates.New(tt.templateName, tt.config)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, tmpl)
				}
			}
		})
	}
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name         string
		templateText string
		context      interface{}
		expected     string
		expectError  bool
		setupFunc    func(*template.Template)
	}{
		{
			name:         "simple string template",
			templateText: "Hello, {{.Name}}!",
			context:      map[string]string{"Name": "World"},
			expected:     "Hello, World!",
			expectError:  false,
		},
		{
			name:         "template with multiple variables",
			templateText: "{{.First}} {{.Last}} is {{.Age}} years old",
			context: map[string]interface{}{
				"First": "John",
				"Last":  "Doe",
				"Age":   30,
			},
			expected:    "John Doe is 30 years old",
			expectError: false,
		},
		{
			name:         "template with conditionals",
			templateText: "{{if .IsActive}}Active{{else}}Inactive{{end}}",
			context:      map[string]bool{"IsActive": true},
			expected:     "Active",
			expectError:  false,
		},
		{
			name:         "template with range",
			templateText: "{{range .Items}}{{.}} {{end}}",
			context:      map[string]interface{}{"Items": []string{"a", "b", "c"}},
			expected:     "a b c ",
			expectError:  false,
		},
		{
			name:         "template with nested data",
			templateText: "{{.Person.Name}} lives in {{.Person.City}}",
			context: map[string]interface{}{
				"Person": map[string]string{
					"Name": "Alice",
					"City": "New York",
				},
			},
			expected:    "Alice lives in New York",
			expectError: false,
		},
		{
			name:         "empty template",
			templateText: "",
			context:      map[string]string{"Name": "Test"},
			expected:     "",
			expectError:  false,
		},
		{
			name:         "template with missing field",
			templateText: "Hello, {{.MissingField}}!",
			context:      map[string]string{"Name": "World"},
			expected:     "Hello, <no value>!",
			expectError:  false,
		},
		{
			name:         "nil context",
			templateText: "Static text",
			context:      nil,
			expected:     "Static text",
			expectError:  false,
		},
		{
			name:         "template with custom function",
			templateText: "{{upper .Name}}",
			context:      map[string]string{"Name": "world"},
			expected:     "WORLD",
			expectError:  false,
			setupFunc: func(tmpl *template.Template) {
				funcMap := template.FuncMap{
					"upper": strings.ToUpper,
				}
				tmpl.Funcs(funcMap)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := templates.NewTemplateConfig()
			tmpl, err := templates.New("test", *config)
			assert.NoError(t, err)
			
			if tt.setupFunc != nil {
				tt.setupFunc(tmpl)
			}
			
			_, err = tmpl.Parse(tt.templateText)
			assert.NoError(t, err)
			
			result, err := templates.Execute(tmpl, tt.context)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestExecute_CustomDelimiters(t *testing.T) {
	config := templates.TemplateConfig{
		LeftDelim:  "<<",
		RightDelim: ">>",
	}
	
	tmpl, err := templates.New("custom", config)
	assert.NoError(t, err)
	
	_, err = tmpl.Parse("Hello, <<.Name>>!")
	assert.NoError(t, err)
	
	result, err := templates.Execute(tmpl, map[string]string{"Name": "World"})
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", result)
}

func TestApplyTemplateFuncs(t *testing.T) {
	tests := []struct {
		name         string
		funcMap      template.FuncMap
		templateText string
		context      interface{}
		expected     string
	}{
		{
			name: "single function",
			funcMap: template.FuncMap{
				"double": func(n int) int { return n * 2 },
			},
			templateText: "{{double .Value}}",
			context:      map[string]int{"Value": 21},
			expected:     "42",
		},
		{
			name: "multiple functions",
			funcMap: template.FuncMap{
				"add": func(a, b int) int { return a + b },
				"mul": func(a, b int) int { return a * b },
			},
			templateText: "{{add .A .B}} {{mul .A .B}}",
			context:      map[string]int{"A": 3, "B": 4},
			expected:     "7 12",
		},
		{
			name:         "nil funcMap",
			funcMap:      nil,
			templateText: "{{.Value}}",
			context:      map[string]string{"Value": "test"},
			expected:     "test",
		},
		{
			name:         "empty funcMap",
			funcMap:      template.FuncMap{},
			templateText: "{{.Value}}",
			context:      map[string]string{"Value": "test"},
			expected:     "test",
		},
		{
			name: "function with string manipulation",
			funcMap: template.FuncMap{
				"repeat": func(s string, n int) string {
					return strings.Repeat(s, n)
				},
			},
			templateText: "{{repeat .Text .Count}}",
			context:      map[string]interface{}{"Text": "Ha", "Count": 3},
			expected:     "HaHaHa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := templates.NewTemplateConfig()
			tmpl, err := templates.New("test", *config)
			assert.NoError(t, err)
			
			templates.ApplyTemplateFuncs(tmpl, tt.funcMap)
			
			_, err = tmpl.Parse(tt.templateText)
			assert.NoError(t, err)
			
			result, err := templates.Execute(tmpl, tt.context)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTemplateIntegration(t *testing.T) {
	// Test a complete workflow
	config := templates.NewTemplateConfig()
	config.LeftDelim = "[["
	config.RightDelim = "]]"
	
	tmpl, err := templates.New("integration-test", *config)
	assert.NoError(t, err)
	
	// Apply custom functions
	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
	}
	templates.ApplyTemplateFuncs(tmpl, funcMap)
	
	// Parse template
	templateText := `
[[range .Users -]]
Name: [[upper .Name]]
Email: [[lower .Email]]
Title: [[title .Role]]
[[end]]`
	
	_, err = tmpl.Parse(templateText)
	assert.NoError(t, err)
	
	// Execute with context
	context := map[string]interface{}{
		"Users": []map[string]string{
			{"Name": "john doe", "Email": "JOHN@EXAMPLE.COM", "Role": "software engineer"},
			{"Name": "jane smith", "Email": "JANE@EXAMPLE.COM", "Role": "product manager"},
		},
	}
	
	result, err := templates.Execute(tmpl, context)
	assert.NoError(t, err)
	
	expected := `
Name: JOHN DOE
Email: john@example.com
Title: Software Engineer
Name: JANE SMITH
Email: jane@example.com
Title: Product Manager
`
	assert.Equal(t, expected, result)
}

func BenchmarkTemplateExecution(b *testing.B) {
	config := templates.NewTemplateConfig()
	tmpl, _ := templates.New("bench", *config)
	tmpl.Parse("Hello, {{.Name}}! You have {{.Count}} messages.")
	
	context := map[string]interface{}{
		"Name":  "User",
		"Count": 42,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = templates.Execute(tmpl, context)
	}
}

func BenchmarkTemplateWithFunctions(b *testing.B) {
	config := templates.NewTemplateConfig()
	tmpl, _ := templates.New("bench-func", *config)
	
	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
		"add":   func(a, b int) int { return a + b },
	}
	templates.ApplyTemplateFuncs(tmpl, funcMap)
	
	tmpl.Parse("{{upper .Name}} has {{add .Base .Bonus}} points")
	
	context := map[string]interface{}{
		"Name":  "player",
		"Base":  100,
		"Bonus": 50,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = templates.Execute(tmpl, context)
	}
}

func ExampleNew() {
	// Create a template with custom delimiters
	config := templates.TemplateConfig{
		LeftDelim:  "[[",
		RightDelim: "]]",
	}
	
	tmpl, _ := templates.New("example", config)
	tmpl.Parse("Hello, [[.Name]]!")
	
	result, _ := templates.Execute(tmpl, map[string]string{"Name": "World"})
	fmt.Println(result)
	// Output: Hello, World!
}

func ExampleApplyTemplateFuncs() {
	// Create a template with custom functions
	config := templates.NewTemplateConfig()
	tmpl, _ := templates.New("example", *config)
	
	// Add custom functions
	funcMap := template.FuncMap{
		"double": func(n int) int { return n * 2 },
		"upper":  strings.ToUpper,
	}
	templates.ApplyTemplateFuncs(tmpl, funcMap)
	
	tmpl.Parse("{{upper .Name}} scored {{double .Score}} points!")
	
	context := map[string]interface{}{
		"Name":  "alice",
		"Score": 21,
	}
	
	result, _ := templates.Execute(tmpl, context)
	fmt.Println(result)
	// Output: ALICE scored 42 points!
}