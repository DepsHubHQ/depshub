package rules

import (
	"testing"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestRuleMaxMinorUpdates(t *testing.T) {
	tests := []struct {
		name           string
		manifests      []types.Manifest
		packagesInfo   PackagesInfo
		expectedLength int
		expectError    bool
	}{
		{
			name: "no dependencies",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{},
				},
			},
			packagesInfo:   PackagesInfo{},
			expectedLength: 0,
			expectError:    false,
		},
		{
			name: "below threshold - single dependency no update",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "test-pkg",
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "test/path",
							},
						},
					},
				},
			},
			packagesInfo: PackagesInfo{
				"test-pkg": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
					},
				},
			},
			expectedLength: 0,
			expectError:    false,
		},
		{
			name: "below threshold - minor update available",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "test-pkg",
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "test/path",
							},
						},
					},
				},
			},
			packagesInfo: PackagesInfo{
				"test-pkg": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
						"1.1.0": {},
					},
				},
			},
			expectedLength: 1,
			expectError:    false,
		},
		{
			name: "above threshold - too many minor updates",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "pkg1",
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "path1",
							},
						},
						{
							Name:    "pkg2",
							Version: "2.0.0",
							Definition: types.Definition{
								Path: "path2",
							},
						},
					},
				},
			},
			packagesInfo: PackagesInfo{
				"pkg1": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
						"1.1.0": {},
					},
				},
				"pkg2": {
					Versions: map[string]types.PackageVersion{
						"2.0.0": {},
						"2.1.0": {},
					},
				},
			},
			expectedLength: 1,
			expectError:    false,
		},
		{
			name: "ignore patch version updates",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "test-pkg",
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "test/path",
							},
						},
					},
				},
			},
			packagesInfo: PackagesInfo{
				"test-pkg": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
						"1.0.1": {},
					},
				},
			},
			expectedLength: 0,
			expectError:    false,
		},
		{
			name: "ignore major version updates",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "test-pkg",
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "test/path",
							},
						},
					},
				},
			},
			packagesInfo: PackagesInfo{
				"test-pkg": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
						"2.0.0": {},
					},
				},
			},
			expectedLength: 0,
			expectError:    false,
		},
		{
			name: "multiple manifests",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "pkg1",
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "path1",
							},
						},
					},
				},
				{
					Dependencies: []types.Dependency{
						{
							Name:    "pkg2",
							Version: "2.0.0",
							Definition: types.Definition{
								Path: "path2",
							},
						},
					},
				},
			},
			packagesInfo: PackagesInfo{
				"pkg1": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
						"1.1.0": {},
					},
				},
				"pkg2": {
					Versions: map[string]types.PackageVersion{
						"2.0.0": {},
						"2.1.0": {},
					},
				},
			},
			expectedLength: 1,
			expectError:    false,
		},
		{
			name: "package not in info",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "missing-pkg",
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "test/path",
							},
						},
					},
				},
			},
			packagesInfo:   PackagesInfo{},
			expectedLength: 0,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewRuleMaxMinorUpdates()

			// Test rule metadata
			assert.Equal(t, "max-minor-updates", rule.GetName())
			assert.Equal(t, LevelError, rule.GetLevel())
			assert.Equal(t, "The total number of minor updates is too high", rule.GetMessage())

			// Test rule check
			mistakes, err := rule.Check(tt.manifests, tt.packagesInfo)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedLength, len(mistakes))

			// Additional checks for when mistakes are found
			if tt.expectedLength > 0 {
				assert.Equal(t, rule, mistakes[0].Rule)
				assert.NotEmpty(t, mistakes[0].Definitions)
			}
		})
	}
}

