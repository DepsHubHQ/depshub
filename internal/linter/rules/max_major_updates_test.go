package rules

import (
	"testing"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestRuleMaxMajorUpdates(t *testing.T) {
	tests := []struct {
		name           string
		manifests      []types.Manifest
		packagesInfo   types.PackagesInfo
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
			packagesInfo:   types.PackagesInfo{},
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
			packagesInfo: types.PackagesInfo{
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
			name: "below threshold - one major update out of five dependencies",
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
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "path2",
							},
						},
						{
							Name:    "pkg3",
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "path3",
							},
						},
						{
							Name:    "pkg4",
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "path4",
							},
						},
						{
							Name:    "pkg5",
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "path5",
							},
						},
					},
				},
			},
			packagesInfo: types.PackagesInfo{
				"pkg1": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
						"2.0.0": {},
					},
				},
				"pkg2": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
					},
				},
				"pkg3": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
					},
				},
				"pkg4": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
					},
				},
				"pkg5": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
					},
				},
			},
			expectedLength: 0,
			expectError:    false,
		},
		{
			name: "above threshold - two major updates out of five dependencies (40%)",
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
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "path2",
							},
						},
						{
							Name:    "pkg3",
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "path3",
							},
						},
						{
							Name:    "pkg4",
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "path4",
							},
						},
						{
							Name:    "pkg5",
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "path5",
							},
						},
					},
				},
			},
			packagesInfo: types.PackagesInfo{
				"pkg1": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
						"2.0.0": {},
					},
				},
				"pkg2": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
						"2.0.0": {},
					},
				},
				"pkg3": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
					},
				},
				"pkg4": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
					},
				},
				"pkg5": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
					},
				},
			},
			expectedLength: 1,
			expectError:    false,
		},
		{
			name: "ignore minor version updates",
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
			packagesInfo: types.PackagesInfo{
				"test-pkg": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
						"1.1.0": {},
					},
				},
			},
			expectedLength: 0,
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
			packagesInfo: types.PackagesInfo{
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
							Version: "1.0.0",
							Definition: types.Definition{
								Path: "path2",
							},
						},
					},
				},
			},
			packagesInfo: types.PackagesInfo{
				"pkg1": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
						"2.0.0": {},
					},
				},
				"pkg2": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
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
			packagesInfo:   types.PackagesInfo{},
			expectedLength: 0,
			expectError:    false,
		},
		{
			name: "check exact version match requirement",
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
			packagesInfo: types.PackagesInfo{
				"test-pkg": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {},
						"2.0.1": {}, // Different patch version
						"2.1.0": {}, // Different minor version
					},
				},
			},
			expectedLength: 0,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewRuleMaxMajorUpdates()

			// Test rule metadata
			assert.Equal(t, "max-major-updates", rule.GetName())
			assert.Equal(t, LevelError, rule.GetLevel())
			assert.Equal(t, "The total number of major updates is too high", rule.GetMessage())

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
