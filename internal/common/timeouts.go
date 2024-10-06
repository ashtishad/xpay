package common

import "time"

// ServiceTimeouts defines default timeouts for different services
type ServiceTimeouts struct {
	Read    time.Duration
	Write   time.Duration
	Startup time.Duration
}

// Timeouts contains timeout configurations for different services
var Timeouts = struct {
	Auth    ServiceTimeouts
	User    ServiceTimeouts
	Wallet  ServiceTimeouts
	Card    ServiceTimeouts
	Server  ServiceTimeouts
	Default ServiceTimeouts
}{
	Auth: ServiceTimeouts{
		Read:  300 * time.Millisecond,
		Write: 500 * time.Millisecond,
	},
	User: ServiceTimeouts{
		Read:  300 * time.Millisecond,
		Write: 500 * time.Millisecond,
	},
	Wallet: ServiceTimeouts{
		Read:  300 * time.Millisecond,
		Write: 500 * time.Millisecond,
	},
	Card: ServiceTimeouts{
		Read:  300 * time.Millisecond,
		Write: 500 * time.Millisecond,
	},
	Server: ServiceTimeouts{
		Read:    5 * time.Second,
		Write:   10 * time.Second,
		Startup: 30 * time.Second,
	},
	Default: ServiceTimeouts{
		Read:    300 * time.Millisecond,
		Write:   500 * time.Millisecond,
		Startup: 5 * time.Second,
	},
}

// // Specific overrides for endpoints that need different timeouts
// const (
// 	TimeoutRegister         = 600 * time.Millisecond
// 	TimeoutComplexOperation = 1 * time.Second
// )
