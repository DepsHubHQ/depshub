package cargo

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/depshubhq/depshub/pkg/types"
)

func TestCargo_Managed(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "valid cargo.toml",
			path:     "cargo.toml",
			expected: true,
		},
		{
			name:     "invalid filename",
			path:     "not-cargo.toml",
			expected: false,
		},
		{
			name:     "cargo.toml in nested path",
			path:     filepath.Join("some", "nested", "path", "cargo.toml"),
			expected: true,
		},
	}

	cargo := Cargo{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cargo.Managed(tt.path); got != tt.expected {
				t.Errorf("Cargo.Managed() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCargo_Dependencies(t *testing.T) {
	cargo := Cargo{}
	deps, err := cargo.Dependencies("testdata/Cargo.toml")
	if err != nil {
		t.Fatalf("Failed to parse dependencies: %v", err)
	}

	// Test specific dependencies with Definition fields
	tests := []struct {
		name          string
		dep           string
		version       string
		isDev         bool
		line          int
		containsInRaw string // String that should be present in RawLine
	}{
		{
			name:          "dev dependency tokio",
			dep:           "tokio",
			line:          9,
			version:       "1.0.0",
			isDev:         true,
			containsInRaw: `tokio = { version = "1.0.0"`,
		},
		{
			name:          "regular dependency serde",
			dep:           "serde",
			line:          19,
			version:       "1.0",
			isDev:         false,
			containsInRaw: `serde = "1.0"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, dep := range deps {
				if dep.Name == tt.dep {
					found = true
					if dep.Line != tt.line {
						t.Errorf("Expected line %d for %s, got %d", tt.line, tt.dep, dep.Line)
					}
					if dep.Version != tt.version {
						t.Errorf("Expected version %s for %s, got %s", tt.version, tt.dep, dep.Version)
					}
					if dep.Dev != tt.isDev {
						t.Errorf("Expected Dev=%v for %s, got %v", tt.isDev, tt.dep, dep.Dev)
					}
					if !strings.Contains(dep.RawLine, tt.containsInRaw) {
						t.Errorf("Expected RawLine to contain '%s', got '%s'", tt.containsInRaw, dep.RawLine)
					}
				}
			}
			if !found {
				t.Errorf("Dependency %s not found", tt.dep)
			}
		})
	}
}

func TestFindLineInfo(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		key          string
		expectedLine int
		expectedRaw  string
		shouldFind   bool
	}{
		{
			name: "simple dependency",
			content: `[dependencies]
serde = "1.0"
tokio = "1.0"`,
			key:          "serde",
			expectedRaw:  `serde = "1.0"`,
			expectedLine: 2,
			shouldFind:   true,
		},
		{
			name: "dependency with features",
			content: `[dependencies]
tokio = { version = "1.0.0", features = ["full"] }`,
			key:          "tokio",
			expectedRaw:  `tokio = { version = "1.0.0", features = ["full"] }`,
			expectedLine: 2,
			shouldFind:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line, rawLine := findLineInfo([]byte(tt.content), tt.key)
			if line != tt.expectedLine {
				t.Errorf("Expected line %d, got %d", tt.expectedLine, line)
			}
			if rawLine != tt.expectedRaw {
				t.Errorf("Expected raw line '%s', got '%s'", tt.expectedRaw, rawLine)
			}
		})
	}
}

func TestCargo_GetType(t *testing.T) {
	cargo := Cargo{}
	if got := cargo.GetType(); got != types.Cargo {
		t.Errorf("Cargo.GetType() = %v, want %v", got, types.Cargo)
	}
}

func TestCargo_LockfilePath(t *testing.T) {
	cargo := Cargo{}
	path := "testdata/Cargo.toml"
	_, err := cargo.LockfilePath(path)
	if err == nil {
		t.Error("Expected error for non-existent lockfile, got nil")
	}
}

func TestCleanVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{"caret version", "^1.0.0", "1.0.0"},
		{"tilde version", "~1.0.0", "1.0.0"},
		{"exact version", "1.0.0", "1.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanVersion(tt.version); got != tt.expected {
				t.Errorf("cleanVersion() = %v, want %v", got, tt.expected)
			}
		})
	}
}
