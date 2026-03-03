package middleware

import (
	"database/sql"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"

	local "github.com/eirka/eirka-index/config"
)

var (
	sitemap = make(map[string]*local.SiteData)
	mu      = new(sync.RWMutex)
)

// Details gets the imageboard settings from the request for the page handler variables
func Details() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get host and normalize it (strip port if present)
		host := c.Request.Host
		if hostParts := strings.Split(host, ":"); len(hostParts) > 1 {
			host = hostParts[0]
		}

		mu.RLock()
		// check the sitemap to see if its cached
		site := sitemap[host]
		mu.RUnlock()

		// if not query the database
		if site == nil {
			mu.Lock()
			// Double-check within the write lock to prevent race
			site = sitemap[host]
			if site == nil {
				sitedata := &local.SiteData{}

				// Get Database handle
				dbase, err := db.GetDb()
				if err != nil {
					mu.Unlock()
					c.JSON(e.ErrorMessage(e.ErrInternalError))
					c.Error(err).SetMeta("Details.GetDb")
					c.Abort()
					return
				}

				// get the info about the imageboard
				err = dbase.QueryRow(`SELECT ib_id,ib_title,ib_description,ib_nsfw,ib_api,ib_img,ib_style,ib_logo,ib_discord FROM imageboards WHERE ib_domain = ?`, host).Scan(&sitedata.Ib, &sitedata.Title, &sitedata.Desc, &sitedata.Nsfw, &sitedata.API, &sitedata.Img, &sitedata.Style, &sitedata.Logo, &sitedata.Discord)
				if err == sql.ErrNoRows {
					mu.Unlock()
					c.JSON(e.ErrorMessage(e.ErrNotFound))
					c.Error(err).SetMeta("Details.QueryRow")
					c.Abort()
					return
				} else if err != nil {
					mu.Unlock()
					c.JSON(e.ErrorMessage(e.ErrInternalError))
					c.Error(err).SetMeta("Details.QueryRow")
					c.Abort()
					return
				}

				sitemap[host] = sitedata
			}
			// Capture the value while still holding the lock
			site = sitemap[host]
			mu.Unlock()
		}

		c.Set("host", host)

		// set the site data for the request
		// this is used in the controllers
		c.Set("sitemap", site)

		c.Next()

	}

}
