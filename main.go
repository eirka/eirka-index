package main

import (
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strings"
	"sync"

	"github.com/eirka/eirka-libs/config"
	"github.com/eirka/eirka-libs/db"

	local "eirka-index/config"
)

var (
	sitemap map[string]*SiteData
	mu      sync.RWMutex
)

func init() {

	// map to hold site data so we dont hit the database every time
	sitemap = make(map[string]*SiteData)

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

	// Get limits and stuff from database
	config.GetDatabaseSettings()

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
	t = template.Must(t.ParseGlob(fmt.Sprintf("%sincludes/*.tmpl", local.Settings.Directories.AssetsDir)))

	r := gin.Default()

	// load template into gin
	r.SetHTMLTemplate(t)

	// serve our assets
	r.Static("/assets", local.Settings.Directories.AssetsDir)

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
