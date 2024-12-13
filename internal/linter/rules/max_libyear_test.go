package rules

import (
	"fmt"
	"testing"
	"time"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestRuleMaxLibyear(t *testing.T) {
	rule := NewRuleMaxLibyear()

	// Test rule metadata
	t.Run("metadata", func(t *testing.T) {
		assert.Equal(t, "max-libyear", rule.GetName())
		assert.Equal(t, LevelError, rule.GetLevel())
		assert.Equal(t, "The total libyear of all dependencies is too high", rule.GetMessage())
	})

	// Use a fixed time for test data setup
	baseTime := time.Now()

	tests := []struct {
		name      string
		manifests []types.Manifest
		info      types.PackagesInfo
		want      []Mistake
		wantErr   bool
	}{
		{
			name: "within libyear limit",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "pkg1",
							Version: "1.0.0",
						},
						{
							Name:    "pkg2",
							Version: "2.0.0",
						},
					},
				},
			},
			info: types.PackagesInfo{
				"pkg1": types.Package{
					Time: map[string]time.Time{
						"1.0.0": baseTime.AddDate(0, -6, 0),
					},
				},
				"pkg2": types.Package{
					Time: map[string]time.Time{
						"2.0.0": baseTime.AddDate(-1, 0, 0),
					},
				},
			},
			want:    []Mistake{}, // Total < MaxLibyear
			wantErr: false,
		},
		{
			name: "exceeds libyear limit",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "old-pkg",
							Version: "1.0.0",
						},
					},
				},
			},
			info: types.PackagesInfo{
				"old-pkg": types.Package{
					Time: map[string]time.Time{
						"1.0.0": baseTime.AddDate(-31, 0, 0),
					},
				},
			},
			want: []Mistake{
				{
					Rule: rule,
					Definitions: []types.Definition{
						{
							// Don't hardcode the exact value since it will change based on when the test runs
							// Instead verify the format and approximate range in the test
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple packages exceeding limit",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "pkg1",
							Version: "1.0.0",
						},
						{
							Name:    "pkg2",
							Version: "2.0.0",
						},
					},
				},
			},
			info: types.PackagesInfo{
				"pkg1": types.Package{
					Time: map[string]time.Time{
						"1.0.0": baseTime.AddDate(-16, 0, 0),
					},
				},
				"pkg2": types.Package{
					Time: map[string]time.Time{
						"2.0.0": baseTime.AddDate(-16, 0, 0),
					},
				},
			},
			want: []Mistake{
				{
					Rule: rule,
					Definitions: []types.Definition{
						{
							// Don't hardcode the exact value since it will change based on when the test runs
							// Instead verify the format and approximate range in the test
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "package not found in info",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "missing-pkg",
							Version: "1.0.0",
						},
					},
				},
			},
			info:    types.PackagesInfo{},
			want:    []Mistake{},
			wantErr: false,
		},
		{
			name: "version not found in package time map",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "pkg",
							Version: "1.0.0",
						},
					},
				},
			},
			info: types.PackagesInfo{
				"pkg": types.Package{
					Time: map[string]time.Time{},
				},
			},
			want:    []Mistake{},
			wantErr: false,
		},
		{
			name:      "empty manifests",
			manifests: []types.Manifest{},
			info:      types.PackagesInfo{},
			want:      []Mistake{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rule.Check(tt.manifests, tt.info)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// For cases where we expect mistakes, verify the structure and ranges rather than exact values
			if len(tt.want) > 0 {
				assert.Len(t, got, 1)
				assert.Equal(t, rule, got[0].Rule)
				assert.Len(t, got[0].Definitions, 1)

				// Verify the Path format and that it contains expected parts
				path := got[0].Definitions[0].Path
				assert.Contains(t, path, "Allowed libyear: 25.00")
				assert.Contains(t, path, "Total libyear:")

				// Parse the total libyear value to verify it's in the expected range
				var allowedLibyear, totalLibyear float64
				_, err := fmt.Sscanf(path, "Allowed libyear: %f. Total libyear: %f", &allowedLibyear, &totalLibyear)
				assert.NoError(t, err)
				assert.Greater(t, totalLibyear, DefaultMaxLibyear)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
