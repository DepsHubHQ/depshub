package rules

import (
	"testing"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewRuleNoAnyTag(t *testing.T) {
	rule := NewRuleNoAnyTag()

	assert.Equal(t, "no-any-tag", rule.GetName())
	assert.Equal(t, LevelWarning, rule.GetLevel())
	assert.Equal(t, `Disallow the use of the "any" version tag`, rule.GetMessage())
}

func TestRuleNoAnyTag_Check(t *testing.T) {
	tests := []struct {
		name      string
		manifests []types.Manifest
		want      []Mistake
		wantErr   bool
	}{
		{
			name:      "no manifests",
			manifests: []types.Manifest{},
			want:      nil,
			wantErr:   false,
		},
		{
			name: "manifest with valid version tags",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Definition: types.Definition{Path: "dep1"},
							Version:    "1.0.0",
						},
						{
							Definition: types.Definition{Path: "dep2"},
							Version:    "^2.0.0",
						},
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "manifest with star version",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Definition: types.Definition{Path: "dep1"},
							Version:    "*",
						},
					},
				},
			},
			want: []Mistake{
				{
					Rule: NewRuleNoAnyTag(),
					Definitions: []types.Definition{
						{Path: "dep1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "manifest with latest version",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Definition: types.Definition{Path: "dep1"},
							Version:    "latest",
						},
					},
				},
			},
			want: []Mistake{
				{
					Rule: NewRuleNoAnyTag(),
					Definitions: []types.Definition{
						{Path: "dep1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "manifest with empty version",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Definition: types.Definition{Path: "dep1"},
							Version:    "",
						},
					},
				},
			},
			want: []Mistake{
				{
					Rule: NewRuleNoAnyTag(),
					Definitions: []types.Definition{
						{Path: "dep1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple manifests with mixed version tags",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Definition: types.Definition{Path: "dep1"},
							Version:    "*",
						},
						{
							Definition: types.Definition{Path: "dep2"},
							Version:    "1.0.0",
						},
					},
				},
				{
					Dependencies: []types.Dependency{
						{
							Definition: types.Definition{Path: "dep3"},
							Version:    "latest",
						},
						{
							Definition: types.Definition{Path: "dep4"},
							Version:    "",
						},
						{
							Definition: types.Definition{Path: "dep5"},
							Version:    "^2.0.0",
						},
					},
				},
			},
			want: []Mistake{
				{
					Rule: NewRuleNoAnyTag(),
					Definitions: []types.Definition{
						{Path: "dep1"},
					},
				},
				{
					Rule: NewRuleNoAnyTag(),
					Definitions: []types.Definition{
						{Path: "dep3"},
					},
				},
				{
					Rule: NewRuleNoAnyTag(),
					Definitions: []types.Definition{
						{Path: "dep4"},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewRuleNoAnyTag()
			got, err := rule.Check(tt.manifests)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
