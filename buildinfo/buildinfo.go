// Package buildinfo provides build information for the executable. Its variables are mostly injected at build time.
package buildinfo

import "time"

//goland:noinspection ALL
var (
	// Version is the short commit hash or version tag of the executable.
	Version = ""
	// Revision is the short commit hash of the executable.
	Revision = ""
	// CompiledAt is the RFC3339 formatted time indicating that time of the compilation.
	CompiledAt = ""

	defaultCompileTime = time.Now()
)

//goland:noinspection ALL
func CompiledAtTime() time.Time {
	if len(CompiledAt) == 0 {
		return defaultCompileTime
	}

	t, err := time.Parse(time.RFC3339, CompiledAt)
	if err != nil {
		return defaultCompileTime
	}

	return t
}
