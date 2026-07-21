package util

import (
	"os"
)

// PlatformVersion returns the application version from APP_VERSION environment variable.
// If APP_VERSION is not set, returns a default version string.
func PlatformVersion() string {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		return "1.0.0" // Default version if not set
	}
	return version
}
