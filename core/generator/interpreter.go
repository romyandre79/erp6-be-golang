package generator

import (
	"reflect"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// DynamicRunner handles execution of interpreted Go code
type DynamicRunner struct {
	interpreter *interp.Interpreter
}

// NewDynamicRunner creates a new Yaegi interpreter
func NewDynamicRunner() *DynamicRunner {
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	return &DynamicRunner{interpreter: i}
}

// Eval evaluates Go source code and returns the result
func (d *DynamicRunner) Eval(src string) (reflect.Value, error) {
	return d.interpreter.Eval(src)
}

// ExecuteHandler runs a specific handler function from the interpreted code
// The script must define a function: func Handle(ctx *WorkflowContext) error
func (d *DynamicRunner) ExecuteHandler(src string, ctx *WorkflowContext) error {
	// 1. Evaluate the source code to define the function
	_, err := d.interpreter.Eval(src)
	if err != nil {
		return err
	}

	// 2. Get the Handle function
	v, err := d.interpreter.Eval("Handle")
	if err != nil {
		return err
	}

	// 3. Call the function
	// Note: Passing *WorkflowContext directly might fail if symbols aren't exported.
	// We might need to use a simpler interface or map.
	// For now, let's try passing the struct.
	args := []reflect.Value{reflect.ValueOf(ctx)}
	results := v.Call(args)

	// 4. Check error return
	if len(results) > 0 && !results[0].IsNil() {
		return results[0].Interface().(error)
	}

	return nil
}
