package rules

import (
	"testing"

	"github.com/depshubhq/depshub/internal/config"
	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestRuleAllowedLicenses(t *testing.T) {
	rule := NewRuleAllowedLicenses()

	// Test rule metadata
	t.Run("metadata", func(t *testing.T) {
		assert.Equal(t, "allowed-licenses", rule.GetName())
		assert.Equal(t, types.LevelError, rule.GetLevel())
		assert.Equal(t, "The license of the package is not allowed.", rule.GetMessage())
	})

	// Test cases for Check method
	testCases := []struct {
		name      string
		manifests []types.Manifest
		info      types.PackagesInfo
		expected  []types.Mistake
	}{
		{
			name: "all licenses allowed",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:       "pkg1",
							Definition: types.Definition{},
						},
						{
							Name:       "pkg2",
							Definition: types.Definition{},
						},
					},
				},
			},
			info: types.PackagesInfo{
				"pkg1": {License: "MIT"},
				"pkg2": {License: "Apache-2.0"},
			},
			expected: []types.Mistake{},
		},
		{
			name: "empty license allowed",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:       "pkg1",
							Definition: types.Definition{},
						},
					},
				},
			},
			info: types.PackagesInfo{
				"pkg1": {License: ""},
			},
			expected: []types.Mistake{},
		},
		{
			name: "disallowed license",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:       "pkg1",
							Definition: types.Definition{},
						},
					},
				},
			},
			info: types.PackagesInfo{
				"pkg1": {License: "GPL-3.0"},
			},
			expected: []types.Mistake{
				{
					Rule: *rule,
					Definitions: []types.Definition{{
						Path:    "",
						Line:    0,
						RawLine: "",
					}},
				},
			},
		},
		{
			name: "multiple manifests with mixed licenses",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:       "pkg1",
							Definition: types.Definition{},
						},
					},
				},
				{
					Dependencies: []types.Dependency{
						{
							Name:       "pkg2",
							Definition: types.Definition{},
						},
					},
				},
			},
			info: types.PackagesInfo{
				"pkg1": {License: "MIT"},
				"pkg2": {License: "GPL-3.0"},
			},
			expected: []types.Mistake{
				{
					Rule: *rule,
					Definitions: []types.Definition{{
						Path:    "",
						Line:    0,
						RawLine: "",
					}},
				},
			},
		},
		{
			name: "package not in info",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:       "unknown-pkg",
							Definition: types.Definition{},
						},
					},
				},
			},
			info:     types.PackagesInfo{},
			expected: []types.Mistake{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mistakes, err := rule.Check(tc.manifests, tc.info, config.Config{})
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, mistakes)
		})
	}
}
