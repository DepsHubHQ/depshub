package rules

import (
	"testing"

	"github.com/depshubhq/depshub/internal/config"
	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestRuleMinWeeklyDownloads(t *testing.T) {
	rule := NewRuleMinWeeklyDownloads()

	// Test rule metadata
	t.Run("metadata", func(t *testing.T) {
		assert.Equal(t, "min-weekly-downloads", rule.GetName())
		assert.Equal(t, types.LevelError, rule.GetLevel())
		assert.Equal(t, "Minimum weekly downloads not met", rule.GetMessage())
	})

	tests := []struct {
		name      string
		manifests []types.Manifest
		info      types.PackagesInfo
		want      []types.Mistake
		wantErr   bool
	}{
		{
			name: "package meets minimum downloads",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name: "popular-pkg",
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
				"popular-pkg": types.Package{
					Downloads: []types.Download{
						{Downloads: 600},
						{Downloads: 500}, // Total: 1100 > MinWeeklyDownloads
					},
				},
			},
			want:    []types.Mistake{},
			wantErr: false,
		},
		{
			name: "package below minimum downloads",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name: "unpopular-pkg",
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
				"unpopular-pkg": types.Package{
					Downloads: []types.Download{
						{Downloads: 400},
						{Downloads: 300}, // Total: 700 < MinWeeklyDownloads
					},
				},
			},
			want: []types.Mistake{
				{
					Rule: *rule,
					Definitions: []types.Definition{
						{
							Path:    "",
							RawLine: "",
							Line:    0,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple packages with mixed download counts",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name: "pkg1",
							Definition: types.Definition{
								Path:    "",
								RawLine: "",
								Line:    0,
							},
						},
						{
							Name: "pkg2",
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
				"pkg1": types.Package{
					Downloads: []types.Download{
						{Downloads: 800},
						{Downloads: 300}, // Total: 1100 > MinWeeklyDownloads
					},
				},
				"pkg2": types.Package{
					Downloads: []types.Download{
						{Downloads: 400},
						{Downloads: 200}, // Total: 600 < MinWeeklyDownloads
					},
				},
			},
			want: []types.Mistake{
				{
					Rule: *rule,
					Definitions: []types.Definition{
						{
							Path:    "",
							RawLine: "",
							Line:    0,
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
							Name: "missing-pkg",
							Definition: types.Definition{
								Path:    "",
								RawLine: "",
								Line:    0,
							},
						},
					},
				},
			},
			info:    types.PackagesInfo{},
			want:    []types.Mistake{},
			wantErr: false,
		},
		{
			name:      "empty manifests",
			manifests: []types.Manifest{},
			info:      types.PackagesInfo{},
			want:      []types.Mistake{},
			wantErr:   false,
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
