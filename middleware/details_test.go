package middleware

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eirka/eirka-libs/config"
	"github.com/eirka/eirka-libs/csrf"
	"github.com/eirka/eirka-libs/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	local "github.com/eirka/eirka-index/config"
)

// SimpleHandler is a simple handler for testing that doesn't require templates
func SimpleHandler(c *gin.Context) {
	// Get sitemap from middleware
	site := c.MustGet("sitemap").(*local.SiteData)

	// Return JSON instead of HTML
	c.JSON(http.StatusOK, gin.H{
		"ib_id":             site.Ib,
		"title":             site.Title,
		"imageboards_count": len(site.Imageboards),
	})
}

func performHTMLRequest(r http.Handler, method, path string, host string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(nil))
	req.Header.Set("Content-Type", "text/html")
	req.Header.Set("X-Real-IP", "127.0.0.1")
	req.Host = host // Set the host dynamically
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// setupRouter creates a gin router with middleware for testing
// clearSiteCache clears the sitemap cache between tests
func clearSiteCache() {
	mu.Lock()
	defer mu.Unlock()
	// Reset the sitemap to an empty map
	sitemap = make(map[string]*local.SiteData)
}

func setupRouter() (*gin.Engine, sqlmock.Sqlmock, error) {
	// Clear the sitemap cache for each test
	clearSiteCache()

	config.Settings = &config.Config{
		Prim: config.Prim{
			CSS: "test.css",
			JS:  "test.js",
		},
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(Details())
	router.Use(csrf.Cookie())

	// Use our simple handler instead of IndexController to avoid template issues
	router.GET("/", SimpleHandler)

	// Setup test database with mock
	mock, err := db.NewTestDb()
	if err != nil {
		return nil, nil, err
	}

	return router, mock, nil
}

func TestDetailsSQLSuccess(t *testing.T) {
	var err error

	router, mock, err := setupRouter()
	assert.NoError(t, err, "Setup should not error")
	defer db.CloseDb()

	// Mock first query for imageboard settings
	ibrows := sqlmock.NewRows([]string{"ib_id", "ib_title", "ib_description", "ib_nsfw", "ib_api", "ib_img", "ib_style", "ib_logo", "ib_discord"}).
		AddRow(1, "test board", "a test board", false, "http://test.board/api", "http://test.board/images", "style.css", "logo.png", "http://test.board/discord.json")
	mock.ExpectQuery(`SELECT ib_id,ib_title,ib_description,ib_nsfw,ib_api,ib_img,ib_style,ib_logo,ib_discord FROM imageboards WHERE ib_domain = \?`).
		WithArgs("test.board").
		WillReturnRows(ibrows)

	// Mock second query for other imageboards
	otheribrows := sqlmock.NewRows([]string{"ib_title", "ib_domain"}).
		AddRow("other board", "http://other.board").
		AddRow("another board", "http://another.board")
	mock.ExpectQuery(`SELECT ib_title,ib_domain FROM imageboards WHERE ib_id != \?`).
		WithArgs(1).
		WillReturnRows(otheribrows)

	resp := performHTMLRequest(router, "GET", "/", "test.board")

	// Check response
	assert.Equal(t, 200, resp.Code, "HTTP request code should match")
	assert.Contains(t, resp.Body.String(), "\"ib_id\":1", "Response should contain ib_id:1")
	assert.Contains(t, resp.Body.String(), "\"imageboards_count\":2", "Response should have 2 imageboards")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")
}

func TestDetailsNonExistentDomain(t *testing.T) {
	router, mock, err := setupRouter()
	assert.NoError(t, err, "Setup should not error")
	defer db.CloseDb()

	// Mock database returning no rows for unknown domain
	mock.ExpectQuery(`SELECT ib_id,ib_title,ib_description,ib_nsfw,ib_api,ib_img,ib_style,ib_logo,ib_discord FROM imageboards WHERE ib_domain = \?`).
		WithArgs("test.board").
		WillReturnError(sql.ErrNoRows)

	resp := performHTMLRequest(router, "GET", "/", "test.board")

	// Should return 404 when domain not found
	assert.Equal(t, 404, resp.Code, "Should return 404 for non-existent domain")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")
}

func TestDetailsCachedLookup(t *testing.T) {
	router, mock, err := setupRouter()
	assert.NoError(t, err, "Setup should not error")
	defer db.CloseDb()

	// First request - sets up the cache
	ibrows := sqlmock.NewRows([]string{"ib_id", "ib_title", "ib_description", "ib_nsfw", "ib_api", "ib_img", "ib_style", "ib_logo", "ib_discord"}).
		AddRow(1, "test board", "a test board", false, "http://test.board/api", "http://test.board/images", "style.css", "logo.png", "http://test.board/discord.json")
	mock.ExpectQuery(`SELECT ib_id,ib_title,ib_description,ib_nsfw,ib_api,ib_img,ib_style,ib_logo,ib_discord FROM imageboards WHERE ib_domain = \?`).
		WithArgs("test.board").
		WillReturnRows(ibrows)

	otheribrows := sqlmock.NewRows([]string{"ib_title", "ib_domain"}).
		AddRow("other board", "http://other.board")
	mock.ExpectQuery(`SELECT ib_title,ib_domain FROM imageboards WHERE ib_id != \?`).
		WithArgs(1).
		WillReturnRows(otheribrows)

	// First request - populates cache
	resp1 := performHTMLRequest(router, "GET", "/", "test.board")
	assert.Equal(t, 200, resp1.Code, "First request should succeed")

	// Second request - should use cached data and not query the database
	resp2 := performHTMLRequest(router, "GET", "/", "test.board")
	assert.Equal(t, 200, resp2.Code, "Second request should succeed")

	// Body content should be the same for both requests
	assert.Equal(t, resp1.Body.String(), resp2.Body.String(), "Response should be identical")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")
}

func TestDetailsDatabaseErrors(t *testing.T) {
	router, mock, err := setupRouter()
	assert.NoError(t, err, "Setup should not error")
	defer db.CloseDb()

	// Test case 1: First query fails
	mock.ExpectQuery(`SELECT ib_id,ib_title,ib_description,ib_nsfw,ib_api,ib_img,ib_style,ib_logo,ib_discord FROM imageboards WHERE ib_domain = \?`).
		WithArgs("test.board").
		WillReturnError(fmt.Errorf("database error"))

	resp1 := performHTMLRequest(router, "GET", "/", "test.board")
	assert.Equal(t, 500, resp1.Code, "Should return 500 on database error")

	// Test case 2: Second query fails
	ibrows := sqlmock.NewRows([]string{"ib_id", "ib_title", "ib_description", "ib_nsfw", "ib_api", "ib_img", "ib_style", "ib_logo", "ib_discord"}).
		AddRow(2, "another board", "description", false, "http://api", "http://img", "style.css", "logo.png", "discord.json")
	mock.ExpectQuery(`SELECT ib_id,ib_title,ib_description,ib_nsfw,ib_api,ib_img,ib_style,ib_logo,ib_discord FROM imageboards WHERE ib_domain = \?`).
		WithArgs("test.board").
		WillReturnRows(ibrows)

	mock.ExpectQuery(`SELECT ib_title,ib_domain FROM imageboards WHERE ib_id != \?`).
		WithArgs(2).
		WillReturnError(fmt.Errorf("database error"))

	resp2 := performHTMLRequest(router, "GET", "/", "test.board")
	assert.Equal(t, 500, resp2.Code, "Should return 500 on database error")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")
}

func TestHostPortStripping(t *testing.T) {
	router, mock, err := setupRouter()
	assert.NoError(t, err, "Setup should not error")
	defer db.CloseDb()

	// The host will include a port, but the middleware should strip it before the database query
	// The query should use just "test.board", not "test.board:8080"
	mock.ExpectQuery(`SELECT ib_id,ib_title,ib_description,ib_nsfw,ib_api,ib_img,ib_style,ib_logo,ib_discord FROM imageboards WHERE ib_domain = \?`).
		WithArgs("test.board").
		WillReturnRows(sqlmock.NewRows([]string{"ib_id", "ib_title", "ib_description", "ib_nsfw", "ib_api", "ib_img", "ib_style", "ib_logo", "ib_discord"}).
			AddRow(1, "test board", "a test board", false, "http://test.board/api", "http://test.board/images", "style.css", "logo.png", "http://test.board/discord.json"))

	mock.ExpectQuery(`SELECT ib_title,ib_domain FROM imageboards WHERE ib_id != \?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"ib_title", "ib_domain"}).
			AddRow("other board", "http://other.board"))

	// Use a host with a port
	resp := performHTMLRequest(router, "GET", "/", "test.board:8080")
	assert.Equal(t, 200, resp.Code, "Request with port in host should succeed")

	// Verify the host was properly normalized (port was stripped)
	assert.Contains(t, resp.Body.String(), "\"ib_id\":1", "Response should contain correct data")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")
}
