package events

import "sync"

type HandlerFunc func(data interface{}) error

var (
	handlers = make(map[string][]HandlerFunc)
	mu       sync.RWMutex
)

// Register listener
func Register(name string, handler HandlerFunc) {
	mu.Lock()
	defer mu.Unlock()
	handlers[name] = append(handlers[name], handler)
}

// Trigger event, hentikan jika ada error
func Trigger(name string, data interface{}) error {
	mu.RLock()
	defer mu.RUnlock()

	if list, ok := handlers[name]; ok {
		for _, h := range list {
			if err := h(data); err != nil {
				return err // hentikan chain jika ada error
			}
		}
	}
	return nil
}
