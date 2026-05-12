// Package buildinfo holds version metadata injected at build time via -ldflags.
package buildinfo

import "fmt"

// These are set by the build script via -ldflags="-X ...".
// They default to "dev" values when building outside the release pipeline.
var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

// String returns a human-readable version string, e.g. "v1.0.0 (abc1234, 2026-05-12)".
func String() string {
	return fmt.Sprintf("%s (%s, %s)", Version, Commit, Date)
}
