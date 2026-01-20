//go:build linux

package config

func defaultServiceName() string {
	// On Linux, we auto-detect the first interface
	return ""
}
