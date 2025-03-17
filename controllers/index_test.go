package controllers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/eirka/eirka-index/config"
	"github.com/eirka/eirka-index/templates"
)

func setupTemplateRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Parse templates the same way main.go does
	t := template.Must(template.New("templates").Delims("[[", "]]").Parse(templates.Index))
	t = template.Must(t.Parse(templates.Head))
	t = template.Must(t.Parse(templates.Header))
	t = template.Must(t.Parse(templates.Navmenu))
	t = template.Must(t.Parse(templates.Angular))

	// Create a dummy headinclude template that's used by the head template
	t = template.Must(t.Parse(`[[define "headinclude"]][[end]]`))

	// Create a dummy navmenuinclude template that's used by the header template
	t = template.Must(t.Parse(`[[define "navmenuinclude"]][[end]]`))

	// Set the template in the router
	r.SetHTMLTemplate(t)

	return r
}

func TestIndexControllerTemplateRendering(t *testing.T) {
	r := setupTemplateRouter()

	// Create a test site config
	testSite := &config.SiteData{
		Ib:    1,
		API:   "api.test.com",
		Img:   "img.test.com",
		Title: "Test Board",
		Desc:  "A test imageboard",
		Style: "test.css",
		Logo:  "logo.png",
		Base:  "",
		Imageboards: []config.Imageboard{
			{Title: "Other Board", Address: "other.board"},
		},
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
		// Set the site data and csrf token in the context
		c.Set("sitemap", testSite)
		c.Set("csrf_token", "test-csrf-token")

		// Call the actual controller
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
	assert.Contains(t, html, "ng-app=\"prim\"", "Should contain Angular app directive")
	assert.Contains(t, html, "test-csrf-token", "Should contain CSRF token")
	assert.Contains(t, html, "Test Board", "Should contain board title")
	assert.Contains(t, html, "test.js", "Should contain JS reference")
	assert.Contains(t, html, "test.css", "Should contain CSS reference")
}

func TestErrorControllerTemplateRendering(t *testing.T) {
	r := setupTemplateRouter()

	// Create a test site config
	testSite := &config.SiteData{
		Ib:    1,
		API:   "api.test.com",
		Img:   "img.test.com",
		Title: "Test Board",
		Desc:  "A test imageboard",
		Style: "test.css",
		Logo:  "logo.png",
		Base:  "",
		Imageboards: []config.Imageboard{
			{Title: "Other Board", Address: "other.board"},
		},
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
		// Set the site data and csrf token in the context
		c.Set("sitemap", testSite)
		c.Set("csrf_token", "test-csrf-token")

		// Call the actual controller
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
	assert.Contains(t, html, "ng-app=\"prim\"", "Should contain Angular app directive")
	assert.Contains(t, html, "test-csrf-token", "Should contain CSRF token")
	assert.Contains(t, html, "Test Board", "Should contain board title")
}
