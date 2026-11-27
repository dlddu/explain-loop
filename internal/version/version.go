// Package version holds the service identity reported by /healthz and MCP initialize.
package version

// Name is the service name.
const Name = "explain-loop"

// Version is overridden at build time with -ldflags "-X .../version.Version=<sha>".
var Version = "dev"
