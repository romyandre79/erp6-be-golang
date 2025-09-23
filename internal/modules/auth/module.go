package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"erp6-be-golang/internal/core/logger"
	"erp6-be-golang/internal/core/registry"
)

type AuthModule struct{}

func New() registry.Module { return &AuthModule{} }

func (m *AuthModule) Name() string { return "auth" }

func (m *AuthModule) Init(ctx context.Context, db *gorm.DB, rdb *redis.Client, log logger.Logger) error {
	// prepare user table etc later
	return nil
}

func (m *AuthModule) Migrate(db *gorm.DB) error {
	// auto-migrate user model later
	return nil
}

func (m *AuthModule) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/auth")
	g.POST("/register", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "register endpoint (stub)"})
	})
	g.POST("/login", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "login endpoint (stub)"})
	})
}

func (m *AuthModule) Shutdown(ctx context.Context) error { return nil }

func init() {
	registry.RegisterFactory("auth", func() registry.Module { return New() })
}
