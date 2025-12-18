package generator

import (
	"mime/multipart"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// WorkflowContext holds the context for a component execution
type WorkflowContext struct {
	FiberCtx         *fiber.Ctx
	DB               *gorm.DB
	Params           []WorkflowDetailResult
	Search           bool
	FileHeader       *multipart.FileHeader
	Extras           map[string]interface{}
	DecisionResult   bool
	CurrentComponent Component
}

// ComponentHandler is the interface that all workflow components must implement
type ComponentHandler interface {
	Execute(ctx *WorkflowContext) error
}

// ComponentHandlerFunc allows using a function as a ComponentHandler
type ComponentHandlerFunc func(ctx *WorkflowContext) error

func (f ComponentHandlerFunc) Execute(ctx *WorkflowContext) error {
	return f(ctx)
}

// Registry holds all registered workflow components
type Registry struct {
	mu         sync.RWMutex
	components map[string]ComponentHandler
	dynamic    map[string]string // For Yaegi-based plugins (source code)
	external   map[string]string // For external executable plugins (command path)
}

// NewRegistry creates a new component registry
func NewRegistry() *Registry {
	return &Registry{
		components: make(map[string]ComponentHandler),
		dynamic:    make(map[string]string),
		external:   make(map[string]string),
	}
}

// Register adds a component to the registry
func (r *Registry) Register(name string, handler ComponentHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.components[strings.ToLower(name)] = handler
}

// RegisterDynamic registers a dynamic (Yaegi) component
func (r *Registry) RegisterDynamic(name string, sourceCode string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.dynamic[strings.ToLower(name)] = sourceCode
}

// RegisterExternal registers an external executable component
func (r *Registry) RegisterExternal(name string, commandPath string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.external[strings.ToLower(name)] = commandPath
}

// Get retrieves a component from the registry
func (r *Registry) Get(name string) (ComponentHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	name = strings.ToLower(name)

	// Check external plugins first
	if cmdPath, ok := r.external[name]; ok {
		return ComponentHandlerFunc(func(ctx *WorkflowContext) error {
			return RunExternalPlugin(cmdPath, ctx)
		}), true
	}

	// Check dynamic plugins
	if sourceCode, ok := r.dynamic[name]; ok {
		return ComponentHandlerFunc(func(ctx *WorkflowContext) error {
			runner := NewDynamicRunner()
			return runner.ExecuteHandler(sourceCode, ctx)
		}), true
	}

	// Check static components
	handler, ok := r.components[name]
	return handler, ok
}

// GlobalRegistry is the global instance
var GlobalRegistry = NewRegistry()

// RegisterComponent is a helper to register to the global registry
func RegisterComponent(name string, handler ComponentHandlerFunc) {
	GlobalRegistry.Register(name, handler)
}

// GetComponent is a helper to get from the global registry
func GetComponent(name string) (ComponentHandler, bool) {
	return GlobalRegistry.Get(name)
}
