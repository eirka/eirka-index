package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/eirka/eirka-libs/config"
	"github.com/gin-gonic/gin"

	local "github.com/eirka/eirka-index/config"
)

// renderSPA builds the template data for the Vue SPA and renders the page.
func renderSPA(c *gin.Context, status int) {
	site := c.MustGet("sitemap").(*local.SiteData)

	// Build the primConfig object for the Vue app
	primConfig := map[string]any{
		"ib_id":      site.Ib,
		"title":      site.Title,
		"img_srv":    "//" + site.Img,
		"api_srv":    "//" + site.API,
		"csrf_token": c.MustGet("csrf_token").(string),
		"logo":       site.Logo,
	}

	if site.Discord != "" {
		// Add a cache-busting query param to the discord widget URL
		if u, err := url.Parse(site.Discord); err == nil {
			q := u.Query()
			q.Set("v", strconv.FormatInt(time.Now().Unix(), 10))
			u.RawQuery = q.Encode()
			primConfig["discord_widget"] = u.String()
		}
	}

	cfgJSON, err := json.Marshal(primConfig)
	if err != nil {
		c.String(http.StatusInternalServerError, "config error")
		c.Abort()
		return
	}

	c.HTML(status, "index", gin.H{
		"primjs":  config.Settings.Prim.JS,
		"primcss": config.Settings.Prim.CSS,
		"title":   site.Title,
		"desc":    site.Desc,
		"nsfw":    site.Nsfw,
		"style":   site.Style,
		"config":  template.JS(cfgJSON),
	})
}

// IndexController generates pages for the Vue frontend
func IndexController(c *gin.Context) {
	renderSPA(c, http.StatusOK)
}
