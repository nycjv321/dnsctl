package dns

// Client defines the interface for DNS management operations.
// Implementations provide platform-specific DNS configuration.
type Client interface {
	// ListNetworkServices returns all available network services/interfaces.
	ListNetworkServices() ([]string, error)

	// GetDNSServers returns the current DNS servers for a network service.
	GetDNSServers(service string) ([]string, error)

	// SetDNSServers sets the DNS servers for a network service.
	SetDNSServers(service string, servers []string) error

	// ClearDNSServers clears DNS servers, reverting to DHCP defaults.
	ClearDNSServers(service string) error

	// FlushCache flushes the DNS cache.
	FlushCache() error

	// Name returns the backend name for display purposes.
	Name() string
}
