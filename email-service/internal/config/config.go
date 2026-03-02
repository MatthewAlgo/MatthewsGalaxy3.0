package config

import "os"

// GetEnvOrDefault returns the value of an environment variable or a default value.
// This is the single source of truth for this utility — do not duplicate in other packages.
func GetEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
