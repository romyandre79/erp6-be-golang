package server

import (
	"erp6-be-golang/internal/core/registry"

	"github.com/gin-gonic/gin"
)

// NewRouter membuat router baru dengan semua route dari module
func NewRouter(reg *registry.Registry) *gin.Engine {
	r := gin.Default()

	// Loop semua module aktif
	for _, m := range reg.GetModules() {
		// cek apakah module punya router
		if routerProvider, ok := m.(interface {
			RegisterRoutes(*gin.Engine)
		}); ok {
			routerProvider.RegisterRoutes(r)
		}
	}

	return r
}
