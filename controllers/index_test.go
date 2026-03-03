package controllers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eirka/eirka-index/templates"
	"github.com/eirka/eirka-libs/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	local "github.com/eirka/eirka-index/config"
)

func setupTemplateRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Parse templates the same way main.go does
	t := template.Must(template.New("templates").Delims("[[", "]]").Parse(templates.Index))
	t = template.Must(t.Parse(templates.Head))
	t = template.Must(t.Parse(templates.PrimConfig))

	// Create a dummy headinclude template that's used by the head template
	t = template.Must(t.Parse(`[[define "headinclude"]][[end]]`))

	// Set the template in the router
	r.SetHTMLTemplate(t)

	return r
}

func TestIndexControllerTemplateRendering(t *testing.T) {
	r := setupTemplateRouter()

	// Create a test site config
	testSite := &local.SiteData{
		Ib:    1,
		API:   "api.test.com",
		Img:   "img.test.com",
		Title: "Test Board",
		Desc:  "A test imageboard",
		Style: "test.css",
		Logo:  "logo.png",
	}

	// Configure test settings
	config.Settings = &config.Config{
		Prim: config.Prim{
			CSS: "test.css",
			JS:  "test.js",
		},
	}

	// Set up a route that will use our test site data
	r.GET("/", func(c *gin.Context) {
		c.Set("sitemap", testSite)
		c.Set("csrf_token", "test-csrf-token")
		IndexController(c)
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

	// Check that the HTML contains expected content
	html := w.Body.String()
	assert.Contains(t, html, "<div id=\"app\"></div>", "Should contain Vue app mount point")
	assert.Contains(t, html, "window.primConfig=", "Should contain Vue config injection")
	assert.Contains(t, html, "test-csrf-token", "Should contain CSRF token")
	assert.Contains(t, html, "Test Board", "Should contain board title")
	assert.Contains(t, html, "test.js", "Should contain JS reference")
	assert.Contains(t, html, "test.css", "Should contain CSS reference")
	assert.Contains(t, html, "//api.test.com", "Should contain API server with protocol-relative prefix")
	assert.Contains(t, html, "//img.test.com", "Should contain image server with protocol-relative prefix")
}

func TestIndexControllerQuotesInTitle(t *testing.T) {
	r := setupTemplateRouter()

	// Test that special characters in DB fields are safely JSON-encoded
	// and don't break the JavaScript or allow XSS injection
	testSite := &local.SiteData{
		Ib:    1,
		API:   "api.test.com",
		Img:   "img.test.com",
		Title: `O'Brien's <Board> & "More"`,
		Desc:  "A board with special chars",
		Style: "test.css",
		Logo:  "logo.png",
	}

	config.Settings = &config.Config{
		Prim: config.Prim{
			CSS: "test.css",
			JS:  "test.js",
		},
	}

	r.GET("/", func(c *gin.Context) {
		c.Set("sitemap", testSite)
		c.Set("csrf_token", "test-csrf-token")
		IndexController(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

	html := w.Body.String()
	assert.Contains(t, html, "window.primConfig=", "Should contain config injection")
	// json.Marshal escapes <, >, & to unicode escapes, preventing script injection
	assert.Contains(t, html, `\u003c`, "Angle brackets should be unicode-escaped")
	assert.Contains(t, html, `\u003e`, "Angle brackets should be unicode-escaped")
	assert.Contains(t, html, `\u0026`, "Ampersands should be unicode-escaped")
	// Single quotes are safe inside JSON double-quoted strings
	assert.Contains(t, html, "O'Brien's", "Single quotes should pass through in JSON strings")
	// The raw dangerous characters should NOT appear unescaped in the script block
	assert.NotContains(t, html, `<Board>`, "Raw angle brackets must not appear in script context")
}

func TestErrorControllerTemplateRendering(t *testing.T) {
	r := setupTemplateRouter()

	// Create a test site config
	testSite := &local.SiteData{
		Ib:    1,
		API:   "api.test.com",
		Img:   "img.test.com",
		Title: "Test Board",
		Desc:  "A test imageboard",
		Style: "test.css",
		Logo:  "logo.png",
	}

	// Configure test settings
	config.Settings = &config.Config{
		Prim: config.Prim{
			CSS: "test.css",
			JS:  "test.js",
		},
	}

	// Set up a route that will use our test site data
	r.GET("/error", func(c *gin.Context) {
		c.Set("sitemap", testSite)
		c.Set("csrf_token", "test-csrf-token")
		ErrorController(c)
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)
	r.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusNotFound, w.Code, "Status code should be 404")

	// Check that the HTML contains expected content
	html := w.Body.String()
	assert.Contains(t, html, "<div id=\"app\"></div>", "Should contain Vue app mount point")
	assert.Contains(t, html, "window.primConfig=", "Should contain Vue config injection")
	assert.Contains(t, html, "test-csrf-token", "Should contain CSRF token")
	assert.Contains(t, html, "Test Board", "Should contain board title")
}
