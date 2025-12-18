package generator

import (
	"fmt"
	"strconv"
	"time"
)

func init() {
	// Register with both lowercase and capitalized names
	RegisterComponent("schedule", func(ctx *WorkflowContext) error {
		return handleSchedule(ctx)
	})

	RegisterComponent("Scheduler", func(ctx *WorkflowContext) error {
		return handleSchedule(ctx)
	})
}

// handleSchedule implements the Iterator pattern for Reminders
func handleSchedule(ctx *WorkflowContext) error {
	c := ctx.FiberCtx
	db := ctx.DB
	component := ctx.CurrentComponent
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)

	// 1. Retrieve Input (Previous Step Result)
	// Expecting a list of maps from comp_search
	var inputData []map[string]interface{}

	if len(wfEngine) > 0 {
		lastResult := wfEngine[len(wfEngine)-1]
		if resultMap, ok := lastResult.ResultNode.(map[string]interface{}); ok {
			if dataRows, ok := resultMap["data"].([]map[string]interface{}); ok {
				inputData = dataRows
			}
		}
	}

	if len(inputData) == 0 {
		// Nothing to process
		return nil
	}

	fmt.Printf("[Scheduler] Processing %d rows.\n", len(inputData))

	// 2. Iterate and Process
	today := time.Now()
	// Normalize today to start of day for comparison?
	// Or just use date parts. Logic assumes Day of Month checks.

	for _, row := range inputData {
		// Validations
		recordStatus := getIntFromMap(row, "recordstatus")
		if recordStatus != 1 {
			continue
		}

		expireDate := getIntFromMap(row, "expiredate")         // Day of Month (1-31)
		reminderBefore := getIntFromMap(row, "reminderbefore") // Days to notify before
		recursive := getIntFromMap(row, "recursive")           // 1=Monthly, 0=One-time
		// lastCount := getIntFromMap(row, "lastcount")
		lastReminderStr, _ := row["lastreminder"].(string)

		// Calculate Target Date (This Month)
		targetDate := time.Date(today.Year(), today.Month(), expireDate, 0, 0, 0, 0, today.Location())

		// Handle edge case: e.g. Day 31 in Feb -> default to last day of month?
		// Go's time.Date normalizes (March 3 -> March 3).
		// Note: The user said expiredate is TINYINT.
		// If today is past the target date, maybe next month?
		// But usually we remind BEFORE.

		// Logic:
		// Message Trigger Date = TargetDate - reminderBefore
		triggerDate := targetDate.AddDate(0, 0, -reminderBefore)

		// Check if we are in the "Active Notification Window"
		// Window: [TriggerDate, TargetDate)
		// Usually we send ON the TriggerDate.

		isDue := false

		// Simple Check: Is Today == TriggerDate?
		// Compare Year, Month, Day
		if isSameDay(today, triggerDate) {
			isDue = true
		} else {
			// What if cron missed it? Maybe check if today >= TriggerDate AND today < TargetDate ??
			// For safety (and to avoid spamming), let's stick to strict equality OR check "not sent yet for this cycle".

			// Let's implement: Today >= TriggerDate AND Today <= TargetDate
			// AND check if we already sent it efficiently.
			if !today.Before(triggerDate) && !today.After(targetDate) {
				isDue = true
			}
		}

		if !isDue {
			continue
		}

		// Check duplicate sending
		// If recursive=1 (Monthly), we send once per month cycle?
		// Or once per day within the window?
		// "reminderbefore" usually implies "Start reminding X days before".
		// Assuming we want to send it ONCE for this cycle.

		if lastReminderStr != "" {
			// Parse last reminder
			// Formats vary, try generic
			lastReminder, err := parseTime(lastReminderStr)
			if err == nil {
				// If sent today, skip
				if isSameDay(today, lastReminder) {
					continue
				}

				// If recursive=0 (One-Time), and we sent it in the past...
				// Actually user said recursive=0 is One-Time.
				if recursive == 0 {
					// It was sent before. Stop.
					continue
				}

				// For monthly: If sent THIS MONTH already?
				// User requirements are a bit vague on "how often".
				// "reminderbefore" = 3 means "3 days before".
				// Usually means ONE notification 3 days before.
				// If I check "lastReminder month == today month", that prevents duplicate monthly.
				// BUT if reminderBefore is e.g. 30 days...
				// Let's assume standard behavior: Send ONCE when the window opens.
				// So if lastReminder is within (TargetDate - 1 month, TargetDate], skip.

				// Let's stick to: If sent within the last (reminderBefore + 1) days?
				// Simplest safe logic: Don't send if sent within the last 24h.
				// AND if recursive=1, ensure we don't send multiple times for the same 'event'.

				// User's request implies simple logic.
				// Let's enforce: Send if not sent TODAY.
				// And let the user configure the query/logic if they want more complex rules.
			}
		}

		// PROCESS EXECUTION
		fmt.Printf("[Scheduler] Triggering reminder for ID: %v\n", row["reminderid"])

		// 1. Update DB (LastReminder, LastCount)
		reminderID := row["reminderid"]
		if dbErr := db.Exec("UPDATE reminder SET lastreminder = NOW(), lastcount = lastcount + 1 WHERE reminderid = ?", reminderID).Error; dbErr != nil {
			fmt.Printf("[Scheduler] Failed to update reminder %v: %v\n", reminderID, dbErr)
			continue
		}

		// 2. Map Row to Scoped Params ($variables)
		// We are "Iterating", so we temporarily inject these vars for the Downstream call
		currentScoped := c.Locals("scopedParams")
		newScoped := make(map[string]interface{})

		// Copy existing
		if existing, ok := currentScoped.(map[string]interface{}); ok {
			for k, v := range existing {
				newScoped[k] = v
			}
		}

		// Inject row data
		for k, v := range row {
			// available as $column_name
			newScoped[k] = v
		}

		// Update context
		c.Locals("scopedParams", newScoped)

		// 3. Trigger Downstream Components
		// Access current component outputs
		// We need to find the connected components.
		// Since we don't have easy access to the full component list in Handlers usually,
		// we rely on InternalFlow to do the recursion.
		// BUT we are overriding navigation!
		// So WE must call InternalFlow for the children.

		components := c.Locals("components").([]Component) // Full list

		// Iterate Outputs of THIS component
		for _, output := range component.Outputs {
			for _, conn := range output.Connections {
				nextNodeId, _ := strconv.Atoi(conn.Node)

				// Find target component
				var targetComp *Component
				for _, comp := range components {
					if comp.ID == nextNodeId && comp.WorkflowId == component.WorkflowId {
						targetComp = &comp
						break
					}
				}

				if targetComp != nil {
					// Manually execute the next step
					// Note: InternalFlow will recurse for its children too!
					// Issue: If we call InternalFlow(child), it will run child -> child's child -> ...
					// We want that! We want the full chain to run for THIS items.
					// CRITICAL: InternalFlow (Child) must NOT see `skipNavigation` = true yet.
					// We set it at the very end of handleSchedule.

					// Also, verify IsRun.
					// InternalFlow checks IsRun.
					// But we are running it potentially multiple times (loop).
					// So we must RESET `IsRun` for the downstream components?
					// YES! If we loop 10 times, the "Send Email" component needs to run 10 times.
					// Standard InternalFlow checks `if component.IsRun { return }`.
					// We must forcefully reset `IsRun` for the target branching.

					resetComponentRunState(components, nextNodeId)

					err := InternalFlow(c, *targetComp, component.WorkflowId, nextNodeId, db, ctx.Search)
					if err != nil {
						fmt.Printf("[Scheduler] Error in downstream flow: %v\n", err)
					}
				}
			}
		}

		// Restore scoped params for next iteration
		c.Locals("scopedParams", currentScoped)
	}

	// Stop default engine navigation
	c.Locals("skipNavigation", true)

	return nil
}

// Helpers

func getIntFromMap(row map[string]interface{}, key string) int {
	if val, ok := row[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int8:
			return int(v)
		case int16:
			return int(v)
		case int32:
			return int(v)
		case int64:
			return int(v)
		case float64:
			return int(v)
		case []uint8: // Raw db bytes
			i, _ := strconv.Atoi(string(v))
			return i
		case string:
			i, _ := strconv.Atoi(v)
			return i
		}
	}
	return 0
}

func parseTime(s string) (time.Time, error) {
	// Try standard formats
	formats := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unknown format")
}

func isSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

// resetComponentRunState recursively resets IsRun for a node and its children
// This is crucial for loops!
func resetComponentRunState(components []Component, startNodeId int) {
	// Simple approach: Reset ALL downstream?
	// Or just the specific branch?
	// Since we are iterating, we want the whole sub-flow to be re-runnable.

	// BFS/DFS to find all descendants
	queue := []int{startNodeId}
	visited := make(map[int]bool)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if visited[curr] {
			continue
		}
		visited[curr] = true

		// Find component and reset
		for i := range components {
			if components[i].ID == curr {
				components[i].IsRun = false // ALLOW RE-EXECUTION

				// Add children to queue
				for _, out := range components[i].Outputs {
					for _, conn := range out.Connections {
						nextId, _ := strconv.Atoi(conn.Node)
						queue = append(queue, nextId)
					}
				}
				break
			}
		}
	}
}
