package api

import (
	_ "embed"
	"crypto/subtle"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	FlareData "github.com/soulteary/flare/config/data"
	FlareDefine "github.com/soulteary/flare/config/define"
	FlareModel "github.com/soulteary/flare/config/model"
	FlareMDI "github.com/soulteary/flare/internal/resources/mdi"
)

//go:embed llm.txt
var llmGuide string

func RegisterRouting(router *gin.Engine) {
	// plain-text agent guide always available (no auth)
	router.GET("/llm.txt", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(llmGuide))
	})

	if !FlareDefine.AppFlags.EnableAPI {
		return
	}

	g := router.Group("/api/v1")
	g.Use(apiAuth)
	{
		g.GET("", index)
		g.GET("/apps", getApps)
		g.PUT("/apps", putApps)
		g.GET("/bookmarks", getBookmarks)
		g.PUT("/bookmarks", putBookmarks)
		g.GET("/settings", getSettings)
		g.PUT("/settings", putSettings)
		g.GET("/icons", getIcons)
	}
}

func apiAuth(c *gin.Context) {
	key := FlareDefine.AppFlags.APIKey
	if key == "" {
		c.Next()
		return
	}
	if apiKeyMatch(c.GetHeader("Authorization"), c.GetHeader("X-API-Key"), key) {
		c.Next()
		return
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "unauthorized"})
}

func apiKeyMatch(authorization, xAPIKey, key string) bool {
	const prefix = "Bearer "
	if strings.HasPrefix(authorization, prefix) {
		got := authorization[len(prefix):]
		if subtle.ConstantTimeCompare([]byte(got), []byte(key)) == 1 {
			return true
		}
	}
	return subtle.ConstantTimeCompare([]byte(xAPIKey), []byte(key)) == 1
}

func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"ok": true, "data": data})
}

func fail(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"ok": false, "error": msg})
}

func index(c *gin.Context) {
	ok(c, gin.H{
		"endpoints": []string{
			"GET  /api/v1",
			"GET  /api/v1/apps",
			"PUT  /api/v1/apps",
			"GET  /api/v1/bookmarks",
			"PUT  /api/v1/bookmarks",
			"GET  /api/v1/settings",
			"PUT  /api/v1/settings",
			"GET  /api/v1/icons?q=&limit=20",
		},
		"auth": "optional Bearer / X-API-Key when --api_key set",
		"notes": []string{
			"apps = apps.yml (favorite links on home)",
			"bookmarks = bookmarks.yml",
			"PUT replaces the whole document",
			"bookmark fields: name, link, icon, desc, category, private",
			"icon names: Remix Icon kebab-case, search via /icons",
			"human/agent guide: GET /llm.txt",
		},
	})
}

func getApps(c *gin.Context) {
	ok(c, FlareData.LoadFavoriteBookmarks())
}

func putApps(c *gin.Context) {
	var body FlareModel.Bookmarks
	if err := c.ShouldBindJSON(&body); err != nil {
		fail(c, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}
	if err := validateBookmarks(body); err != "" {
		fail(c, http.StatusBadRequest, err)
		return
	}
	if !FlareData.SaveFavoriteBookmarks(body) {
		fail(c, http.StatusInternalServerError, "save failed")
		return
	}
	ok(c, body)
}

func getBookmarks(c *gin.Context) {
	ok(c, FlareData.LoadNormalBookmarks())
}

func putBookmarks(c *gin.Context) {
	var body FlareModel.Bookmarks
	if err := c.ShouldBindJSON(&body); err != nil {
		fail(c, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}
	if err := validateBookmarks(body); err != "" {
		fail(c, http.StatusBadRequest, err)
		return
	}
	if !FlareData.SaveNormalBookmarks(body) {
		fail(c, http.StatusInternalServerError, "save failed")
		return
	}
	ok(c, body)
}

func getSettings(c *gin.Context) {
	ok(c, FlareData.GetAllSettingsOptions())
}

func putSettings(c *gin.Context) {
	var body FlareModel.Application
	if err := c.ShouldBindJSON(&body); err != nil {
		fail(c, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}
	if !FlareData.SaveAllSettingsOptions(body) {
		fail(c, http.StatusInternalServerError, "save failed")
		return
	}
	ok(c, body)
}

func getIcons(c *gin.Context) {
	q := c.Query("q")
	limit := 20
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	names := FlareMDI.SearchIcons(q, limit)
	ok(c, gin.H{"q": q, "count": len(names), "icons": names})
}

func validateBookmarks(body FlareModel.Bookmarks) string {
	for i, item := range body.Items {
		if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.URL) == "" {
			return "links[" + strconv.Itoa(i) + "]: name and link required"
		}
		if item.Icon != "" && !FlareMDI.IconExists(item.Icon) && !strings.HasPrefix(item.Icon, "http://") && !strings.HasPrefix(item.Icon, "https://") {
			return "links[" + strconv.Itoa(i) + "]: unknown icon \"" + item.Icon + "\" (use GET /api/v1/icons?q=...)"
		}
	}
	return ""
}
