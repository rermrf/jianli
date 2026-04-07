package main

import (
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"jianli/internal/config"
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

	resumeStore := store.NewResumeStore(db)
	draftStore := store.NewResumeDraftStore(db, resumeStore)
	visitorStore := store.NewVisitorStore(db)
	resumeHandler := handler.NewResumeHandler(resumeStore, pdf.NewExporter())
	draftHandler := handler.NewResumeDraftHandler(draftStore)
	visitorHandler := handler.NewVisitorHandler(visitorStore)
	uploadHandler := handler.NewUploadHandler(filepath.Clean("./data/uploads/avatars"))

	router.POST("/api/auth/verify", handler.VerifyAuth(cfg.AuthKey))
	router.GET("/api/resume", resumeHandler.Get)
	router.GET("/api/resume/pdf", resumeHandler.ExportPDF)
	router.POST("/api/visitors", visitorHandler.Create)
	router.PATCH("/api/visitors/:id", visitorHandler.UpdateDuration)
	router.Static("/uploads", filepath.Clean("./data/uploads"))

	protected := router.Group("/api")
	protected.Use(middleware.Auth(cfg.AuthKey))
	protected.PUT("/resume", resumeHandler.Update)
	protected.POST("/resume/drafts", draftHandler.Create)
	protected.GET("/resume/drafts", draftHandler.List)
	protected.GET("/resume/drafts/:id", draftHandler.Get)
	protected.PUT("/resume/drafts/:id/publish", draftHandler.Publish)
	protected.DELETE("/resume/drafts/:id", draftHandler.Delete)
	protected.GET("/visitors", visitorHandler.List)
	protected.GET("/visitors/stats", visitorHandler.Stats)
	protected.POST("/upload/avatar", uploadHandler.UploadAvatar)

	return router
}
