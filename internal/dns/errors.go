package dns

import "errors"

// ErrNoDNSBackend is returned when no supported DNS management system is detected.
var ErrNoDNSBackend = errors.New("no supported DNS management system detected")
