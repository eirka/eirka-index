package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

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
	API         string
	Img         string
	Title       string
	Desc        string
	Nsfw        bool
	Style       string
	Logo        string
	Base        string
	Discord     string
	Imageboards []Imageboard
}

// Imageboard holds an imageboards metadata
type Imageboard struct {
	Title   string
	Address string
}

// Details gets the imageboard settings from the request for the page handler variables
func Details() gin.HandlerFunc {
	return func(c *gin.Context) {

		host := c.Request.Host

		mu.RLock()
		// check the sitemap to see if its cached
		site := sitemap[host]
		mu.RUnlock()

		// if not query the database
		if site == nil {

			sitedata := &SiteData{}

			// Get Database handle
			dbase, err := db.GetDb()
			if err != nil {
				c.JSON(e.ErrorMessage(e.ErrInternalError))
				c.Error(err).SetMeta("Details.GetDb")
				c.Abort()
				return
			}

			// get the info about the imageboard
			err = dbase.QueryRow(`SELECT ib_id,ib_title,ib_description,ib_nsfw,ib_api,ib_img,ib_style,ib_logo,ib_discord FROM imageboards WHERE ib_domain = ?`, host).Scan(&sitedata.Ib, &sitedata.Title, &sitedata.Desc, &sitedata.Nsfw, &sitedata.API, &sitedata.Img, &sitedata.Style, &sitedata.Logo, &sitedata.Discord)
			if err == sql.ErrNoRows {
				c.JSON(e.ErrorMessage(e.ErrNotFound))
				c.Error(err).SetMeta("Details.QueryRow")
				c.Abort()
				return
			} else if err != nil {
				c.JSON(e.ErrorMessage(e.ErrInternalError))
				c.Error(err).SetMeta("Details.QueryRow")
				c.Abort()
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

				err = rows.Scan(&ib.Title, &ib.Address)
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

		c.Set("host", host)

		c.Next()

	}

}

// IndexController generates pages for angularjs frontend
func IndexController(c *gin.Context) {

	mu.RLock()
	site := sitemap[c.MustGet("host").(string)]
	mu.RUnlock()

	var discord string

	// add a cache breaker because their thing is dumb
	if site.Discord != "" {
		nonce := strconv.Itoa(int(site.Ib)) + strconv.Itoa(int(time.Now().Unix()))
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

// ErrorController generates pages and a 404 response
func ErrorController(c *gin.Context) {

	mu.RLock()
	site := sitemap[c.MustGet("host").(string)]
	mu.RUnlock()

	var discord string

	// add a cache breaker because their thing is dumb
	if site.Discord != "" {
		nonce := strconv.Itoa(int(site.Ib)) + strconv.Itoa(int(time.Now().Unix()))
		discord = strings.Join([]string{site.Discord, nonce}, "?")
	}

	c.HTML(http.StatusNotFound, "index", gin.H{
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
