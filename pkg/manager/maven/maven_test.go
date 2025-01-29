package maven

import (
	"github.com/depshubhq/depshub/pkg/types"
	"path/filepath"
	"strings"
	"testing"
)

func TestMaven_Managed(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "valid pom.xml path",
			path:     "pom.xml",
			expected: true,
		},
		{
			name:     "invalid filename",
			path:     "not-pom.xml",
			expected: false,
		},
		{
			name:     "pom.xml in nested path",
			path:     filepath.Join("some", "nested", "path", "pom.xml"),
			expected: true,
		},
	}

	maven := Maven{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maven.Managed(tt.path); got != tt.expected {
				t.Errorf("Maven.Managed() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMaven_Dependencies(t *testing.T) {
	cargo := Maven{}
	deps, err := cargo.Dependencies("testdata/pom.xml")
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
			name:          "regular dependency log4j-api",
			dep:           "org.apache.logging.log4j:log4j-api",
			version:       "2.17.2",
			line:          16,
			isDev:         false,
			containsInRaw: `<artifactId>log4j-api</artifactId>`,
		},
		{
			name:          "regular dependency spring-context",
			dep:           "org.springframework:spring-context",
			version:       "5.3.29",
			line:          24,
			isDev:         false,
			containsInRaw: `<artifactId>spring-context</artifactId>`,
		},
		{
			name:          "dependency with explicit version override log4j-core",
			dep:           "org.apache.logging.log4j:log4j-core",
			version:       "2.16.0",
			line:          40,
			isDev:         false,
			containsInRaw: `<artifactId>log4j-core</artifactId>`,
		},
		{
			name:          "runtime-scoped dependency mysql-connector-java",
			dep:           "mysql:mysql-connector-java",
			version:       "",
			isDev:         false,
			containsInRaw: `<artifactId>mysql-connector-java</artifactId>`,
		},
		{
			name:          "runtime-scoped dependency mysql-connector-java",
			dep:           "mysql:mysql-connector-java",
			version:       "",
			line:          60,
			isDev:         false,
			containsInRaw: `<artifactId>mysql-connector-java</artifactId>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, dep := range deps {
				if dep.Name == tt.dep {
					found = true
					if tt.line != 0 && dep.Line != tt.line {
						t.Errorf("Expected line %d for %s, got %d", tt.line, tt.dep, dep.Line)
					}
					if tt.version != "" && dep.Version != tt.version {
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

func TestMaven_GetType(t *testing.T) {
	maven := Maven{}
	if got := maven.GetType(); got != types.Maven {
		t.Errorf("Maven.GetType() = %v, want %v", got, types.Maven)
	}
}
