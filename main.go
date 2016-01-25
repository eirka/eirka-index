package main

import (
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"

	"github.com/eirka/eirka-libs/config"
	"github.com/eirka/eirka-libs/csrf"
	"github.com/eirka/eirka-libs/db"

	local "github.com/eirka/eirka-index/config"
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

	// Get limits and stuff from database
	config.GetDatabaseSettings()

}

func main() {

	// parse our template
	t := template.Must(template.New("templates").Delims("[[", "]]").Parse(index))
	t = template.Must(t.Parse(head))
	t = template.Must(t.Parse(header))
	t = template.Must(t.Parse(navmenu))
	t = template.Must(t.Parse(angular))
	t = template.Must(t.ParseGlob(fmt.Sprintf("%s/includes/*.tmpl", local.Settings.Directories.AssetsDir)))

	r := gin.Default()

	// load template into gin
	r.SetHTMLTemplate(t)

	// use the details middleware
	r.Use(Details())
	// generates our csrf cookie
	r.Use(csrf.Cookie())

	r.GET("/", IndexController)
	r.GET("/page/:id", IndexController)
	r.GET("/thread/:id/:page", IndexController)
	r.GET("/directory", IndexController)
	r.GET("/directory/:page", IndexController)
	r.GET("/image/:id", IndexController)
	r.GET("/tags/:page", IndexController)
	r.GET("/tags", IndexController)
	r.GET("/tag/:id/:page", IndexController)
	r.GET("/account", IndexController)
	r.GET("/trending", IndexController)
	r.GET("/favorites/:page", IndexController)
	r.GET("/favorites", IndexController)
	r.GET("/admin", IndexController)
	r.GET("/error", ErrorController)

	r.NoRoute(ErrorController)

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", local.Settings.Index.Address, local.Settings.Index.Port),
		Handler: r,
	}

	gracehttp.Serve(s)

}
