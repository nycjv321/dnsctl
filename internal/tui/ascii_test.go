package tui

import "testing"

func TestGetServiceType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Wi-Fi", ServiceTypeWifi},
		{"wi-fi", ServiceTypeWifi},
		{"WI-FI", ServiceTypeWifi},
		{"Thunderbolt Bridge", ServiceTypeThunderbolt},
		{"Thunderbolt Ethernet Slot 1", ServiceTypeThunderbolt},
		{"USB 10/100/1000 LAN", ServiceTypeUSB},
		{"USB Ethernet Adapter", ServiceTypeUSB},
		{"Ethernet", ServiceTypeEthernet},
		{"Ethernet Adaptor (en0)", ServiceTypeEthernet},
		{"Unknown Service", ""},
		{"Bluetooth PAN", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := GetServiceType(tt.input)
			if got != tt.expected {
				t.Errorf("GetServiceType(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGetServiceArt(t *testing.T) {
	tests := []struct {
		serviceName string
		wantContain string
		wantEmpty   bool
	}{
		{"Wi-Fi", "Wi-Fi", false},
		{"Thunderbolt Bridge", "Thunderbolt", false},
		{"USB 10/100/1000 LAN", "USB", false},
		{"Ethernet", "Ethernet", false},
		{"Unknown Service", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.serviceName, func(t *testing.T) {
			art := GetServiceArt(tt.serviceName)
			if tt.wantEmpty {
				if art != "" {
					t.Errorf("GetServiceArt(%q) = %q, want empty string", tt.serviceName, art)
				}
			} else {
				if art == "" {
					t.Errorf("GetServiceArt(%q) returned empty string", tt.serviceName)
				}
			}
		})
	}
}

func TestGetServiceArt_ReturnsEmptyForUnknown(t *testing.T) {
	art := GetServiceArt("Some Unknown Interface")
	if art != "" {
		t.Errorf("GetServiceArt for unknown service should return empty string, got %q", art)
	}
}

func TestGetServiceLabel(t *testing.T) {
	tests := []struct {
		serviceName string
		expected    string
	}{
		{"Wi-Fi", "Wi-Fi"},
		{"Thunderbolt Bridge", "Thunderbolt"},
		{"USB 10/100/1000 LAN", "USB"},
		{"Ethernet", "Ethernet"},
		{"Unknown Service", ""},
	}

	for _, tt := range tests {
		t.Run(tt.serviceName, func(t *testing.T) {
			got := GetServiceLabel(tt.serviceName)
			if got != tt.expected {
				t.Errorf("GetServiceLabel(%q) = %q, want %q", tt.serviceName, got, tt.expected)
			}
		})
	}
}

func TestServiceArt_AllTypesHaveArt(t *testing.T) {
	types := []string{
		ServiceTypeWifi,
		ServiceTypeEthernet,
		ServiceTypeThunderbolt,
		ServiceTypeUSB,
	}

	for _, serviceType := range types {
		t.Run(serviceType, func(t *testing.T) {
			art, ok := ServiceArt[serviceType]
			if !ok {
				t.Errorf("ServiceArt missing entry for %q", serviceType)
			}
			if art == "" {
				t.Errorf("ServiceArt[%q] is empty", serviceType)
			}
		})
	}
}

func TestGetServiceIcon(t *testing.T) {
	tests := []struct {
		serviceName string
		expected    string
	}{
		{"Wi-Fi", "(((•)))"},
		{"wi-fi", "(((•)))"},
		{"Ethernet", "[==]"},
		{"Ethernet Adaptor (en0)", "[==]"},
		{"Thunderbolt Bridge", "⚡"},
		{"USB 10/100/1000 LAN", "[⊏="},
		{"USB Ethernet Adapter", "[⊏="},
		{"Unknown Service", ""},
		{"Bluetooth PAN", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.serviceName, func(t *testing.T) {
			got := GetServiceIcon(tt.serviceName)
			if got != tt.expected {
				t.Errorf("GetServiceIcon(%q) = %q, want %q", tt.serviceName, got, tt.expected)
			}
		})
	}
}

func TestGetServiceIcon_ReturnsEmptyForUnknown(t *testing.T) {
	icon := GetServiceIcon("Some Unknown Interface")
	if icon != "" {
		t.Errorf("GetServiceIcon for unknown service should return empty string, got %q", icon)
	}
}
