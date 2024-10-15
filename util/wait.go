package util

import (
	"fmt"
	"time"
)

// WaitFor function accepts a checking function, interval, timeout, and maxRetries to wait for a resource to be in a ready state.
func WaitFor(checkFunc func() (bool, error), interval, timeout time.Duration, maxRetries int) error {
	// Validate that at least one stopping condition is valid.
	if timeout <= 0 && maxRetries <= 0 {
		return fmt.Errorf("invalid parameters: both timeout and maxRetries cannot be less than or equal to 0")
	}

	// Initialize timeout and retry counters
	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	retries := 0

	for {
		select {
		case <-timeoutChan:
			if timeout > 0 {
				return fmt.Errorf("timeout reached while waiting")
			}
		case <-ticker.C:
			// Check the condition function
			ready, err := checkFunc()
			if err != nil {
				LogError("Error during wait: %v", err) // Log the error but don't return it immediately
				if ready {
					return err
				}
			}

			// If the check succeeds, exit the function.
			if ready {
				return nil
			}

			// Increment retry count only when retries are applicable
			if maxRetries > 0 {
				retries++
				LogInfo("Retrying... Attempt %d/%d", retries, maxRetries)
				if retries >= maxRetries {
					return fmt.Errorf("max retries reached (%d/%d)", retries, maxRetries)
				}
			} else {
				LogInfo("Retrying...")
			}
		}
	}
}

