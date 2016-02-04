package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"sync"

	"github.com/eirka/eirka-libs/config"
	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
)

var (
	sitemap = make(map[string]*SiteData)
	mu      = new(sync.RWMutex)
)

// SiteData holds imageboard settings
type SiteData struct {
	Ib          uint
	Api         string
	Img         string
	Title       string
	Desc        string
	Nsfw        bool
	Style       string
	Logo        string
	Base        string
	Imageboards []Imageboard
}

type Imageboard struct {
	Title   string
	Address string
}

// Details gets the imageboard settings from the request for the page handler variables
func Details() gin.HandlerFunc {
	return func(c *gin.Context) {

		var host, base string

		host = c.Request.Host

		// figure out our path and host
		if path.Dir(host) == "." {
			// we're on the base in this case
			host = path.Base(host)
		} else {
			host = path.Dir(host)
			base = fmt.Sprintf("%s/", path.Base(host))
		}

		fmt.Println(host, base)

		mu.RLock()
		// check the sitemap to see if its cached
		site := sitemap[host]
		mu.RUnlock()

		// if not query the database
		if site == nil {

			sitedata := &SiteData{}

			// set the base for angularjs
			sitedata.Base = base

			// Get Database handle
			dbase, err := db.GetDb()
			if err != nil {
				c.JSON(e.ErrorMessage(e.ErrInternalError))
				c.Error(err).SetMeta("Details.GetDb")
				c.Abort()
				return
			}

			// get the info about the imageboard
			err = dbase.QueryRow(`SELECT ib_id,ib_title,ib_description,ib_nsfw,ib_api,ib_img,ib_style,ib_logo FROM imageboards WHERE ib_domain = ?`, host).Scan(&sitedata.Ib, &sitedata.Title, &sitedata.Desc, &sitedata.Nsfw, &sitedata.Api, &sitedata.Img, &sitedata.Style, &sitedata.Logo)
			if err == sql.ErrNoRows {
				c.JSON(e.ErrorMessage(e.ErrNotFound))
				c.Error(err).SetMeta("Details.QueryRow")
				return
			} else if err != nil {
				c.JSON(e.ErrorMessage(e.ErrInternalError))
				c.Error(err).SetMeta("Details.QueryRow")
				return
			}

			// collect the links to the other imageboards for nav menu
			rows, err := dbase.Query(`SELECT ib_title,ib_domain FROM imageboards WHERE ib_id != ?`, sitedata.Ib)
			if err != nil {
				c.JSON(e.ErrorMessage(e.ErrInternalError))
				c.Error(err).SetMeta("Details.Query")
				c.Abort()
				return
			}
			defer rows.Close()

			for rows.Next() {

				ib := Imageboard{}

				err := rows.Scan(&ib.Title, &ib.Address)
				if err != nil {
					return
				}

				sitedata.Imageboards = append(sitedata.Imageboards, ib)
			}
			if rows.Err() != nil {
				c.JSON(e.ErrorMessage(e.ErrInternalError))
				c.Error(err).SetMeta("Details.Query")
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

// IndexController generates pages for angularjs frontend
func IndexController(c *gin.Context) {

	// Get parameters from csrf middleware
	csrf_token := c.MustGet("csrf_token").(string)

	host := c.Request.Host

	mu.RLock()
	site := sitemap[host]
	mu.RUnlock()

	c.HTML(http.StatusOK, "index", gin.H{
		"primjs":      config.Settings.Prim.Js,
		"primcss":     config.Settings.Prim.Css,
		"ib":          site.Ib,
		"base":        site.Base,
		"apisrv":      site.Api,
		"imgsrv":      site.Img,
		"title":       site.Title,
		"desc":        site.Desc,
		"nsfw":        site.Nsfw,
		"style":       site.Style,
		"logo":        site.Logo,
		"imageboards": site.Imageboards,
		"csrf":        csrf_token,
	})

	return

}

// ErrorController generates pages and a 404 response
func ErrorController(c *gin.Context) {

	// Get parameters from csrf middleware
	csrf_token := c.MustGet("csrf_token").(string)

	host := c.Request.Host

	mu.RLock()
	site := sitemap[host]
	mu.RUnlock()

	c.HTML(http.StatusNotFound, "index", gin.H{
		"primjs":      config.Settings.Prim.Js,
		"primcss":     config.Settings.Prim.Css,
		"ib":          site.Ib,
		"base":        site.Base,
		"apisrv":      site.Api,
		"imgsrv":      site.Img,
		"title":       site.Title,
		"desc":        site.Desc,
		"nsfw":        site.Nsfw,
		"style":       site.Style,
		"logo":        site.Logo,
		"imageboards": site.Imageboards,
		"csrf":        csrf_token,
	})

	return

}
