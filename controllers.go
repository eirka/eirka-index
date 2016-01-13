package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"

	"github.com/eirka/eirka-libs/config"
	"github.com/eirka/eirka-libs/db"
)

var (
	sitemap = make(map[string]*SiteData)
	mu      = new(sync.RWMutex)
)

type SiteData struct {
	Ib          uint
	Api         string
	Img         string
	Title       string
	Desc        string
	Nsfw        bool
	Style       string
	Logo        string
	Imageboards []Imageboard
}

type Imageboard struct {
	Title   string
	Address string
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

			err = dbase.QueryRow(`SELECT ib_id,ib_title,ib_description,ib_nsfw,ib_api,ib_img,ib_style,ib_logo FROM imageboards WHERE ib_domain = ?`, host).Scan(&sitedata.Ib, &sitedata.Title, &sitedata.Desc, &sitedata.Nsfw, &sitedata.Api, &sitedata.Img, &sitedata.Style, &sitedata.Logo)
			if err != nil {
				c.Error(err)
				c.Abort()
				return
			}

			rows, err := dbase.Query(`SELECT ib_title,ib_domain FROM imageboards WHERE ib_id != ?`, sitedata.Ib)
			if err != nil {
				c.Error(err)
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
			err = rows.Err()
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

	// Get parameters from validate middleware
	csrf_token := c.MustGet("csrf_token").(string)

	host := c.Request.Host

	mu.RLock()
	site := sitemap[host]
	mu.RUnlock()

	c.HTML(http.StatusOK, "index", gin.H{
		"primjs":      config.Settings.Prim.Js,
		"primcss":     config.Settings.Prim.Css,
		"ib":          site.Ib,
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

// Handles error messages for wrong routes
func ErrorController(c *gin.Context) {

	c.String(http.StatusNotFound, "Not Found")

	return

}
