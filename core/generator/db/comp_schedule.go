package generator

func init() {
	// Register with both lowercase and capitalized names
	RegisterComponent("schedule", func(ctx *WorkflowContext) error {
		// Schedule component is a marker for the scheduler
		// It doesn't perform actions during flow execution
		return nil
	})

	RegisterComponent("Scheduler", func(ctx *WorkflowContext) error {
		// Schedule component is a marker for the scheduler
		// It doesn't perform actions during flow execution
		return nil
	})
}
