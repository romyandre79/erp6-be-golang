package module

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
	"erp6-be-golang/internal/core/logger"
)

// Module interface untuk setiap modul
type Module interface {
	Name() string
	Init(ctx context.Context, db *gorm.DB, redis *redis.Client, log logger.Logger) error
	RegisterRoutes(r *gin.Engine)
	Migrate(db *gorm.DB) error
	Shutdown(ctx context.Context) error
}
