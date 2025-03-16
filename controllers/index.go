package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/eirka/eirka-index/config"
	"github.com/gin-gonic/gin"
)

// IndexController generates pages for angularjs frontend
func IndexController(c *gin.Context) {

	// get sitemap from session middleware
	site := c.MustGet("sitemap").(*config.SiteData)

	var discord string

	// add a cache breaker because their thing is dumb
	if site.Discord != "" {
		nonce := strconv.FormatUint(uint64(site.Ib), 10) + strconv.FormatInt(time.Now().Unix(), 10)
		discord = strings.Join([]string{site.Discord, nonce}, "?")
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"primjs":      config.Settings.Prim.JS,
		"primcss":     config.Settings.Prim.CSS,
		"ib":          site.Ib,
		"base":        site.Base,
		"apisrv":      site.API,
		"imgsrv":      site.Img,
		"title":       site.Title,
		"desc":        site.Desc,
		"nsfw":        site.Nsfw,
		"style":       site.Style,
		"logo":        site.Logo,
		"discord":     discord,
		"imageboards": site.Imageboards,
		"csrf":        c.MustGet("csrf_token").(string),
	})

}
