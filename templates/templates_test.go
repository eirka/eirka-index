package templates

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplatesParsing(t *testing.T) {
	// Test that all templates can be parsed without errors
	tmpl := template.New("templates").Delims("[[", "]]")

	// Parse each template and check for errors
	var err error

	tmpl, err = tmpl.Parse(Index)
	assert.NoError(t, err, "Index template should parse without errors")

	tmpl, err = tmpl.Parse(Head)
	assert.NoError(t, err, "Head template should parse without errors")

	tmpl, err = tmpl.Parse(PrimConfig)
	assert.NoError(t, err, "PrimConfig template should parse without errors")

	// Create dummy template for includes
	tmpl, err = tmpl.Parse(`[[define "headinclude"]][[end]]`)
	assert.NoError(t, err, "Dummy headinclude template should parse")

	// Test that the important template definitions exist
	for _, name := range []string{"index", "head", "primconfig"} {
		assert.NotNil(t, tmpl.Lookup(name), "Template '%s' should be defined", name)
	}
}
