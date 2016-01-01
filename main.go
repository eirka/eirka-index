package main

import (
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strings"
	"sync"

	"github.com/eirka/eirka-libs/db"

	local "eirka-index/config"
)

var (
	sitemap    map[string]*SiteData
	globaldata GlobalData
	mu         sync.RWMutex
)

func init() {
	// Database connection settings
	dbase := db.Database{
		User:           local.Settings.Database.User,
		Password:       local.Settings.Database.Password,
		Proto:          local.Settings.Database.Proto,
		Host:           local.Settings.Database.Host,
		Database:       local.Settings.Database.Database,
		MaxIdle:        local.Settings.Database.MaxIdle,
		MaxConnections: local.Settings.Database.MaxConnections,
	}

	// Set up DB connection
	dbase.NewDb()

	// map to hold site data so we dont hit the database every time
	sitemap = make(map[string]*SiteData)

	globaldata = GlobalData{
		Primcss: "/static/prim.css",
		Primjs:  "/static/prim.js",
		Imgsrv:  "images.eirka.com",
		Apisrv:  "api.trish.io",
	}

}

func main() {

	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}

	// parse our template
	t := template.Must(template.New("templates").Funcs(funcMap).Delims("[[", "]]").Parse(index))
	t = template.Must(t.Parse(head))
	t = template.Must(t.Parse(header))

	r := gin.Default()

	// load template into gin
	r.SetHTMLTemplate(t)
	// serve our assets
	r.Static("/static", "/data/prim/static")

	// use the details middleware
	r.Use(Details())

	r.GET("/", IndexController)
	r.GET("/page/:id", IndexController)
	r.GET("/thread/:id/:page", IndexController)
	r.GET("/directory", IndexController)
	r.GET("/image/:id", IndexController)
	r.GET("/tags/:page", IndexController)
	r.GET("/tags", IndexController)
	r.GET("/tag/:id/:page", IndexController)
	r.GET("/account", IndexController)
	r.GET("/trending", IndexController)
	r.GET("/favorites/:page", IndexController)
	r.GET("/favorites", IndexController)
	r.GET("/error", IndexController)

	r.NoRoute(ErrorController)

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", local.Settings.Index.Address, local.Settings.Index.Port),
		Handler: r,
	}

	gracehttp.Serve(s)

}

type GlobalData struct {
	Primcss string
	Primjs  string
	Imgsrv  string
	Apisrv  string
}

type SiteData struct {
	Ib    uint
	Title string
	Desc  string
	Nsfw  bool
	Style string
}

// gets the details from the request for the page handler variables
func Details() gin.HandlerFunc {
	return func(c *gin.Context) {

		host := c.Request.Host

		mu.RLock()
		site := sitemap[host]
		mu.RUnlock()

		if site == nil {

			sitedata := &SiteData{}

			// Get Database handle
			dbase, err := db.GetDb()
			if err != nil {
				c.Error(err)
				c.Abort()
				return
			}

			err = dbase.QueryRow(`SELECT ib_id,ib_title,ib_description,ib_nsfw FROM imageboards WHERE ib_domain = ?`, host).Scan(&sitedata.Ib, &sitedata.Title, &sitedata.Desc, &sitedata.Nsfw)
			if err != nil {
				c.Error(err)
				c.Abort()
				return
			}

			mu.Lock()
			sitemap[host] = sitedata
			mu.Unlock()

		}

		c.Next()

	}
}

// Handles index page generation
func IndexController(c *gin.Context) {

	host := c.Request.Host

	mu.RLock()
	site := sitemap[host]
	mu.RUnlock()

	c.HTML(http.StatusOK, "index", gin.H{
		"ib":      site.Ib,
		"title":   site.Title,
		"desc":    site.Desc,
		"nsfw":    site.Nsfw,
		"style":   site.Style,
		"primjs":  globaldata.Primjs,
		"primcss": globaldata.Primcss,
		"imgsrv":  globaldata.Imgsrv,
		"apisrv":  globaldata.Apisrv,
	})

	return

}

// Handles error messages for wrong routes
func ErrorController(c *gin.Context) {

	c.String(http.StatusNotFound, "Not Found")

	return

}
