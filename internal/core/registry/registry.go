package registry

import (
	"context"

	"erp6-be-golang/internal/core/logger"
	"gorm.io/gorm"

	"github.com/redis/go-redis/v9"
)

// Module adalah interface dasar semua modul
type Module interface {
	Name() string
	Init(cfg interface{}, db *gorm.DB, rdb *redis.Client) error
}

// Registry menyimpan daftar module yang aktif
type Registry struct {
	modules []Module
}

// NewRegistry membuat instance registry baru
func NewRegistry() *Registry {
	return &Registry{
		modules: []Module{},
	}
}

// LoadModules: aktifkan modul berdasarkan daftar dari config (sementara dummy)
func (r *Registry) LoadModules(moduleNames []string) []Module {
	// sementara return kosong, nanti kita bisa mapping modul beneran
	return r.modules
}

// InitAll inisialisasi semua module yang ada di registry
func (r *Registry) InitAll(ctx context.Context, db *gorm.DB, rdb *redis.Client, log logger.Logger) error {
	for _, m := range r.modules {
		if err := m.Init(nil, db, rdb); err != nil {
			log.Error(ctx, "module init failed", map[string]interface{}{
				"module": m.Name(),
				"error":  err.Error(),
			})
			return err
		}
		log.Info(ctx, "module initialized", map[string]interface{}{"module": m.Name()})
	}
	return nil
}

// GetModules mengembalikan semua module aktif
func (r *Registry) GetModules() []Module {
	return r.modules
}
