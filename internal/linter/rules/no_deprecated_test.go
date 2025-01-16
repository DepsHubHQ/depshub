package rules

import (
	"testing"

	"github.com/depshubhq/depshub/internal/config"
	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestRuleNoDeprecated(t *testing.T) {
	rule := NewRuleNoDeprecated()

	// Test rule metadata
	t.Run("metadata", func(t *testing.T) {
		assert.Equal(t, "no-deprecated", rule.GetName())
		assert.Equal(t, types.LevelError, rule.GetLevel())
		assert.Equal(t, "Disallow the use of deprecated package versions", rule.GetMessage())
	})

	// Test scenarios
	tests := []struct {
		name      string
		manifests []types.Manifest
		info      types.PackagesInfo
		want      []types.Mistake
		wantErr   bool
	}{
		{
			name: "deprecated package version",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "test-pkg",
							Version: "1.0.0",
							Definition: types.Definition{
								Path:    "",
								RawLine: "",
								Line:    0,
							},
						},
					},
				},
			},
			info: types.PackagesInfo{
				"test-pkg": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {
							Version:    "1.0.0",
							Deprecated: "This version is deprecated",
						},
					},
				},
			},
			want: []types.Mistake{
				{
					Rule: *rule,
					Definitions: []types.Definition{{
						Path:    "",
						RawLine: "",
						Line:    0,
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "non-deprecated package version",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "test-pkg",
							Version: "1.0.0",
						},
					},
				},
			},
			info: types.PackagesInfo{
				"test-pkg": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {
							Version:    "1.0.0",
							Deprecated: "",
						},
					},
				},
			},
			want:    []types.Mistake{},
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
			want:    []types.Mistake{},
			wantErr: false,
		},
		{
			name: "version not found in package versions",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "test-pkg",
							Version: "2.0.0",
						},
					},
				},
			},
			info: types.PackagesInfo{
				"test-pkg": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {
							Version:    "1.0.0",
							Deprecated: "",
						},
					},
				},
			},
			want:    []types.Mistake{},
			wantErr: false,
		},
		{
			name: "multiple manifests with mixed deprecated status",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:       "pkg1",
							Version:    "1.0.0",
							Definition: types.Definition{},
						},
					},
				},
				{
					Dependencies: []types.Dependency{
						{
							Name:       "pkg2",
							Version:    "2.0.0",
							Definition: types.Definition{},
						},
					},
				},
			},
			info: types.PackagesInfo{
				"pkg1": {
					Versions: map[string]types.PackageVersion{
						"1.0.0": {
							Version:    "1.0.0",
							Deprecated: "Deprecated version",
						},
					},
				},
				"pkg2": {
					Versions: map[string]types.PackageVersion{
						"2.0.0": {
							Version:    "2.0.0",
							Deprecated: "",
						},
					},
				},
			},
			want: []types.Mistake{
				{
					Rule: *rule,
					Definitions: []types.Definition{{
						Path:    "",
						RawLine: "",
						Line:    0,
					}},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rule.Check(tt.manifests, tt.info, config.Config{})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
