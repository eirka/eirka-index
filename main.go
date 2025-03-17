package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/facebookgo/pidfile"
	"github.com/gin-gonic/gin"

	"github.com/eirka/eirka-libs/config"
	"github.com/eirka/eirka-libs/csrf"
	"github.com/eirka/eirka-libs/db"

	local "github.com/eirka/eirka-index/config"
	c "github.com/eirka/eirka-index/controllers"
	m "github.com/eirka/eirka-index/middleware"
	"github.com/eirka/eirka-index/templates"
)

func init() {

	// Database connection settings
	dbase := db.Database{
		User:           local.Settings.Database.User,
		Password:       local.Settings.Database.Password,
		Proto:          local.Settings.Database.Protocol,
		Host:           local.Settings.Database.Host,
		Database:       local.Settings.Database.Database,
		MaxIdle:        local.Settings.Index.DatabaseMaxIdle,
		MaxConnections: local.Settings.Index.DatabaseMaxConnections,
	}

	// Set up DB connection
	dbase.NewDb()

	// Get limits and stuff from database
	config.GetDatabaseSettings()

}

func main() {

	// create pid file
	pidfile.SetPidfilePath("/run/eirka/eirka-index.pid")

	err := pidfile.Write()
	if err != nil {
		panic("Could not write pid file")
	}

	// parse our template
	t := template.Must(template.New("templates").Delims("[[", "]]").Parse(templates.Index))
	t = template.Must(t.Parse(templates.Head))
	t = template.Must(t.Parse(templates.Header))
	t = template.Must(t.Parse(templates.Navmenu))
	t = template.Must(t.Parse(templates.Angular))
	t = template.Must(t.Parse(templates.HeadInclude)) // Add empty templates for includes
	t = template.Must(t.Parse(templates.NavMenuInclude))
	t = template.Must(t.ParseGlob(fmt.Sprintf("%s/includes/*.tmpl", local.Settings.Directories.AssetsDir)))

	r := gin.Default()

	// load template into gin
	r.SetHTMLTemplate(t)

	// use the details middleware
	r.Use(m.Details())
	// generates our csrf cookie
	r.Use(csrf.Cookie())

	// these routes are handled by angularjs
	r.GET("/", c.IndexController)
	r.GET("/page/:id", c.IndexController)
	r.GET("/thread/:id/:page", c.IndexController)
	r.GET("/directory", c.IndexController)
	r.GET("/directory/:page", c.IndexController)
	r.GET("/image/:id", c.IndexController)
	r.GET("/tags/:page", c.IndexController)
	r.GET("/tags", c.IndexController)
	r.GET("/tag/:id/:page", c.IndexController)
	r.GET("/account", c.IndexController)
	r.GET("/trending", c.IndexController)
	r.GET("/favorites/:page", c.IndexController)
	r.GET("/favorites", c.IndexController)
	r.GET("/admin", c.IndexController)
	r.GET("/error", c.ErrorController)

	// if nothing matches
	r.NoRoute(c.ErrorController)

	s := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", local.Settings.Index.Host, local.Settings.Index.Port),
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           r,
	}

	err = gracehttp.Serve(s)
	if err != nil {
		panic("Could not start server")
	}

}
