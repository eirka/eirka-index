package main

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/eirka/eirka-libs/db"
)

type SiteData struct {
	Ib    uint
	Api   string
	Img   string
	Title string
	Desc  string
	Nsfw  bool
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

			err = dbase.QueryRow(`SELECT ib_id,ib_title,ib_description,ib_nsfw,ib_api,ib_img FROM imageboards WHERE ib_domain = ?`, host).Scan(&sitedata.Ib, &sitedata.Title, &sitedata.Desc, &sitedata.Nsfw, &sitedata.Api, &sitedata.Img)
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
		"ib":     site.Ib,
		"apisrv": site.Api,
		"imgsrv": site.Img,
		"title":  site.Title,
		"desc":   site.Desc,
		"nsfw":   site.Nsfw,
	})

	return

}

// Handles error messages for wrong routes
func ErrorController(c *gin.Context) {

	c.String(http.StatusNotFound, "Not Found")

	return

}
