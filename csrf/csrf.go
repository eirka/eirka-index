package csrf

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	// the name of CSRF header
	HeaderName = "X-XSRF-TOKEN"
	// the name of the form field
	FormFieldName = "csrf_token"
	// the name of CSRF cookie
	CookieName = "csrf_token"
	// the name of the session cookie for angularjs
	SessionName = "XSRF-TOKEN"
)

// generates two cookies: a long term csrf token for a user, and a masked session token to verify against
func Cookie() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Vary", "Cookie")

		// the token from the users cookie
		var csrfToken []byte

		tokenCookie, err := c.Request.Cookie(CookieName)
		if err == nil {
			csrfToken = b64decode(tokenCookie.Value)
		}

		// if the user doesnt have a csrf token create one
		if len(csrfToken) != tokenLength {
			// creates a 32 bit token
			csrfToken = generateToken()

			// set the users csrf token tookie
			csrfCookie := &http.Cookie{
				Name:     CookieName,
				Value:    b64encode(csrfToken),
				Expires:  time.Now().Add(356 * 24 * time.Hour),
				Path:     "/",
				HttpOnly: true,
			}

			http.SetCookie(c.Writer, csrfCookie)

		}

		// set the users csrf token tookie
		sessionCookie := &http.Cookie{
			Name:  SessionName,
			Value: b64encode(maskToken(csrfToken)),
			Path:  "/",
		}

		http.SetCookie(c.Writer, sessionCookie)

		c.Next()

	}
}
