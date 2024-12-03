package rules

import (
	"testing"
	"time"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestRuleMaxPackageAge(t *testing.T) {
	rule := NewRuleMaxPackageAge()

	// Test rule metadata
	t.Run("metadata", func(t *testing.T) {
		assert.Equal(t, "max-package-age", rule.GetName())
		assert.Equal(t, LevelError, rule.GetLevel())
		assert.Equal(t, "Disallow the use of any package that is older than a certain age (in months).", rule.GetMessage())
	})

	now := time.Now()
	mistakes := make([]Mistake, 0)

	// Test scenarios
	tests := []struct {
		name      string
		manifests []types.Manifest
		info      PackagesInfo
		want      []Mistake
		wantErr   bool
	}{
		{
			name: "package older than max age",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "old-pkg",
							Version: "1.0.0",
							Definition: types.Definition{
								Line: 1,
							},
						},
					},
				},
			},
			info: PackagesInfo{
				"old-pkg": {
					Time: map[string]time.Time{
						"1.0.0": now.AddDate(0, -(MaxPackageAge + 6), 0), // 6 months older than max age
					},
				},
			},
			want: []Mistake{
				{
					Rule: rule,
					Definitions: []types.Definition{
						{
							Line: 1,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "package within max age",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "recent-pkg",
							Version: "1.0.0",
							Definition: types.Definition{
								Line: 1,
							},
						},
					},
				},
			},
			info: PackagesInfo{
				"recent-pkg": {
					Time: map[string]time.Time{
						"1.0.0": now.AddDate(0, -(MaxPackageAge - 1), 0), // 1 month newer than max age
					},
				},
			},
			want:    mistakes,
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
			info:    PackagesInfo{},
			want:    mistakes,
			wantErr: false,
		},
		{
			name: "version not found in package times",
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
			info: PackagesInfo{
				"test-pkg": {
					Time: map[string]time.Time{
						"1.0.0": now,
					},
				},
			},
			want:    mistakes,
			wantErr: false,
		},
		{
			name: "multiple manifests with mixed age status",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:    "old-pkg",
							Version: "1.0.0",
							Definition: types.Definition{
								Line: 1,
							},
						},
					},
				},
				{
					Dependencies: []types.Dependency{
						{
							Name:    "new-pkg",
							Version: "1.0.0",
							Definition: types.Definition{
								Line: 2,
							},
						},
					},
				},
			},
			info: PackagesInfo{
				"old-pkg": {
					Time: map[string]time.Time{
						"1.0.0": now.AddDate(0, -(MaxPackageAge + 12), 0), // 12 months older than max age
					},
				},
				"new-pkg": {
					Time: map[string]time.Time{
						"1.0.0": now.AddDate(0, -12, 0), // Only 12 months old
					},
				},
			},
			want: []Mistake{
				{
					Rule: rule,
					Definitions: []types.Definition{
						{
							Line: 1,
						},
					},
				},
			},
			wantErr: false,
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
			assert.Equal(t, tt.want, got)
		})
	}
}
