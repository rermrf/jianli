package main

import (
	"log"
	"net/http"
	"os"
	pathpkg "path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"jianli/internal/config"
	"jianli/internal/geoip"
	"jianli/internal/handler"
	"jianli/internal/middleware"
	"jianli/internal/pdf"
	"jianli/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := store.Open(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := newRouter(cfg, db).Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

func newRouter(cfg config.Config, db *gorm.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(cfg.FrontendOrigin))

	// Trust only loopback proxies so c.ClientIP() picks up X-Forwarded-For
	// when nginx/Caddy fronts the Go server. Non-loopback traffic falls
	// back to RemoteAddr which is also fine.
	if err := router.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		log.Printf("warn: SetTrustedProxies: %v", err)
	}

	// Resolver always non-nil: Lookup if the xdb file is loadable, else a
	// loopback-only fallback so visits from localhost / LAN still get a
	// sensible city ("本地网络") rather than "未知".
	var geoResolver geoip.Resolver
	if l, err := geoip.NewLookup(cfg.IP2RegionPath); err == nil {
		geoResolver = l
	} else {
		log.Printf("warn: ip2region disabled (%v); using loopback-only resolver", err)
		geoResolver = geoip.NewLoopbackResolver()
	}

	resumeStore := store.NewResumeStore(db)
	draftStore := store.NewResumeDraftStore(db, resumeStore)
	visitorStore := store.NewVisitorStore(db)
	siteSettingsStore := store.NewSiteSettingsStore(db)
	resumeHandler := handler.NewResumeHandler(resumeStore, siteSettingsStore, pdf.NewExporter(cfg.BrowserPath, cfg.Port))
	draftHandler := handler.NewResumeDraftHandler(draftStore)
	visitorHandler := handler.NewVisitorHandler(visitorStore, geoResolver)
	settingsHandler := handler.NewSettingsHandler(siteSettingsStore)
	uploadHandler := handler.NewUploadHandler(filepath.Clean("./data/uploads/avatars"))

	router.POST("/api/auth/verify", handler.VerifyAuth(cfg.AuthKey))
	router.GET("/api/resume", resumeHandler.Get)
	router.GET("/api/resume/pdf", resumeHandler.ExportPDF)
	router.POST("/api/visitors", visitorHandler.Create)
	router.POST("/api/visitors/:id", visitorHandler.UpdateDuration)
	router.PATCH("/api/visitors/:id", visitorHandler.UpdateDuration)
	router.POST("/api/visitors/:id/pdf-export", visitorHandler.MarkPDFExported)
	router.Static("/uploads", filepath.Clean("./data/uploads"))

	protected := router.Group("/api")
	protected.Use(middleware.Auth(cfg.AuthKey))
	protected.PUT("/resume", resumeHandler.Update)
	protected.GET("/settings", settingsHandler.Get)
	protected.PUT("/settings", settingsHandler.Update)
	protected.POST("/resume/drafts", draftHandler.Create)
	protected.GET("/resume/drafts", draftHandler.List)
	protected.GET("/resume/drafts/:id", draftHandler.Get)
	protected.PUT("/resume/drafts/:id/publish", draftHandler.Publish)
	protected.DELETE("/resume/drafts/:id", draftHandler.Delete)
	protected.GET("/visitors", visitorHandler.List)
	protected.GET("/visitors/stats", visitorHandler.Stats)
	protected.GET("/visitors/trend", visitorHandler.Trend)
	protected.POST("/upload/avatar", uploadHandler.UploadAvatar)

	router.NoRoute(serveFrontend(filepath.Clean("./web/dist")))

	return router
}

func serveFrontend(distDir string) gin.HandlerFunc {
	indexPath := filepath.Join(distDir, "index.html")

	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusNotFound)
			return
		}

		requestPath := c.Request.URL.Path
		if requestPath == "/api" || strings.HasPrefix(requestPath, "/api/") || requestPath == "/uploads" || strings.HasPrefix(requestPath, "/uploads/") {
			c.Status(http.StatusNotFound)
			return
		}

		cleanPath := pathpkg.Clean("/" + strings.TrimPrefix(requestPath, "/"))
		if cleanPath != "/" {
			candidatePath := filepath.Join(distDir, filepath.FromSlash(strings.TrimPrefix(cleanPath, "/")))
			if info, err := os.Stat(candidatePath); err == nil && !info.IsDir() {
				c.File(candidatePath)
				return
			}
		}

		if info, err := os.Stat(indexPath); err == nil && !info.IsDir() {
			c.File(indexPath)
			return
		}

		c.Status(http.StatusNotFound)
	}
}
