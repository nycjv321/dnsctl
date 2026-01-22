package tui

import "strings"

// Service type constants for ASCII art lookup.
const (
	ServiceTypeWifi        = "wifi"
	ServiceTypeEthernet    = "ethernet"
	ServiceTypeThunderbolt = "thunderbolt"
	ServiceTypeUSB         = "usb"
)

// ASCII art for each service type (compact icons).
const (
	wifiArt = `
 (((•)))`

	ethernetArt = `
  [==]`

	thunderboltArt = `
    ⚡`

	usbArt = `
  [⊏=`
)

// ServiceArt maps service types to their ASCII art.
var ServiceArt = map[string]string{
	ServiceTypeWifi:        wifiArt,
	ServiceTypeEthernet:    ethernetArt,
	ServiceTypeThunderbolt: thunderboltArt,
	ServiceTypeUSB:         usbArt,
}

// GetServiceType returns the art key for a service name.
// Maps macOS service names to internal service type constants.
// Returns empty string for unknown services.
// Note: USB is checked before Ethernet because "USB Ethernet Adapter"
// should be classified as USB, not Ethernet.
func GetServiceType(serviceName string) string {
	lower := strings.ToLower(serviceName)
	switch {
	case strings.Contains(lower, "wi-fi"):
		return ServiceTypeWifi
	case strings.Contains(lower, "thunderbolt"):
		return ServiceTypeThunderbolt
	case strings.Contains(lower, "usb"):
		return ServiceTypeUSB
	case strings.Contains(lower, "ethernet"):
		return ServiceTypeEthernet
	default:
		return ""
	}
}

// GetServiceArt returns the ASCII art for a service name.
// Returns empty string for unknown services.
func GetServiceArt(serviceName string) string {
	serviceType := GetServiceType(serviceName)
	if art, ok := ServiceArt[serviceType]; ok {
		return art
	}
	return ""
}

// GetServiceLabel returns a short label for the service type.
// Returns empty string for unknown services.
func GetServiceLabel(serviceName string) string {
	serviceType := GetServiceType(serviceName)
	switch serviceType {
	case ServiceTypeWifi:
		return "Wi-Fi"
	case ServiceTypeEthernet:
		return "Ethernet"
	case ServiceTypeThunderbolt:
		return "Thunderbolt"
	case ServiceTypeUSB:
		return "USB"
	default:
		return ""
	}
}

// GetServiceIcon returns a compact inline icon for a service name.
// Returns empty string for unknown services.
func GetServiceIcon(serviceName string) string {
	serviceType := GetServiceType(serviceName)
	switch serviceType {
	case ServiceTypeWifi:
		return "(((•)))"
	case ServiceTypeEthernet:
		return "[==]"
	case ServiceTypeThunderbolt:
		return "⚡"
	case ServiceTypeUSB:
		return "[⊏="
	default:
		return ""
	}
}
