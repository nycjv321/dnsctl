//go:build linux

package config

func defaultServiceName() string {
	// On Linux, interfaces are auto-detected at runtime.
	// Return "auto" as a placeholder to indicate auto-detection.
	return "auto"
}
