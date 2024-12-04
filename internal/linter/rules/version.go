package rules

import (
	"regexp"
	"strconv"
	"strings"
)

func parseVersion(version string) (major, minor, patch int) {
	// Remove all characters except numbers and dots
	reg := regexp.MustCompile(`[^\d.]`)
	cleaned := reg.ReplaceAllString(version, "")

	// Split version string by dots
	parts := strings.Split(cleaned, ".")

	// Handle incomplete version strings
	if len(parts) < 1 {
		return 0, 0, 0
	}

	// Parse major version
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		major = 0
	}

	// Parse minor version if present
	if len(parts) > 1 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil {
			minor = 0
		}
	}

	// Parse patch version if present
	if len(parts) > 2 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			patch = 0
		}
	}

	return major, minor, patch
}

