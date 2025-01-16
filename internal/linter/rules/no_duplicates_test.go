package rules

import (
	"testing"

	"github.com/depshubhq/depshub/internal/config"
	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewRuleNoDuplicates(t *testing.T) {
	rule := NewRuleNoDuplicates()

	assert.Equal(t, "no-duplicates", rule.GetName())
	assert.Equal(t, types.LevelError, rule.GetLevel())
	assert.Equal(t, "Disallow the same package to be listed multiple times", rule.GetMessage())
}

func TestRuleNoDuplicates_Check(t *testing.T) {
	tests := []struct {
		name      string
		manifests []types.Manifest
		want      []types.Mistake
		wantErr   bool
	}{
		{
			name:      "no manifests",
			manifests: []types.Manifest{},
			want:      nil,
			wantErr:   false,
		},
		{
			name: "manifest with no duplicates",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:       "pkg1",
							Definition: types.Definition{Path: "path/pkg1"},
						},
						{
							Name:       "pkg2",
							Definition: types.Definition{Path: "path/pkg2"},
						},
						{
							Name:       "pkg3",
							Definition: types.Definition{Path: "path/pkg3"},
						},
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "manifest with single duplicate",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:       "pkg1",
							Definition: types.Definition{Path: "path/pkg1"},
						},
						{
							Name:       "pkg2",
							Definition: types.Definition{Path: "path/pkg2"},
						},
						{
							Name:       "pkg1", // duplicate
							Definition: types.Definition{Path: "path/pkg1-duplicate"},
						},
					},
				},
			},
			want: []types.Mistake{
				{
					Rule: *NewRuleNoDuplicates(),
					Definitions: []types.Definition{
						{Path: "path/pkg1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "manifest with multiple duplicates",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:       "pkg1",
							Definition: types.Definition{Path: "path/pkg1"},
						},
						{
							Name:       "pkg2",
							Definition: types.Definition{Path: "path/pkg2"},
						},
						{
							Name:       "pkg1", // duplicate
							Definition: types.Definition{Path: "path/pkg1-duplicate"},
						},
						{
							Name:       "pkg2", // duplicate
							Definition: types.Definition{Path: "path/pkg2-duplicate"},
						},
					},
				},
			},
			want: []types.Mistake{
				{
					Rule: *NewRuleNoDuplicates(),
					Definitions: []types.Definition{
						{Path: "path/pkg1"},
					},
				},
				{
					Rule: *NewRuleNoDuplicates(),
					Definitions: []types.Definition{
						{Path: "path/pkg2"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple manifests with duplicates",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:       "pkg1",
							Definition: types.Definition{Path: "path1/pkg1"},
						},
						{
							Name:       "pkg1", // duplicate
							Definition: types.Definition{Path: "path1/pkg1-duplicate"},
						},
					},
				},
				{
					Dependencies: []types.Dependency{
						{
							Name:       "pkg2",
							Definition: types.Definition{Path: "path2/pkg2"},
						},
						{
							Name:       "pkg2", // duplicate
							Definition: types.Definition{Path: "path2/pkg2-duplicate"},
						},
					},
				},
			},
			want: []types.Mistake{
				{
					Rule: *NewRuleNoDuplicates(),
					Definitions: []types.Definition{
						{Path: "path1/pkg1"},
					},
				},
				{
					Rule: *NewRuleNoDuplicates(),
					Definitions: []types.Definition{
						{Path: "path2/pkg2"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "manifest with empty dependencies",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{},
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "manifest with single dependency",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name:       "pkg1",
							Definition: types.Definition{Path: "path/pkg1"},
						},
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewRuleNoDuplicates()
			got, err := rule.Check(tt.manifests, nil, config.Config{})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
