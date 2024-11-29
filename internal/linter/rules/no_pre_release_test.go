package rules

import (
	"testing"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewRuleNoPreRelease(t *testing.T) {
	rule := NewRuleNoPreRelease()

	assert.Equal(t, "no-pre-release", rule.GetName())
	assert.Equal(t, LevelError, rule.GetLevel())
	assert.Equal(t, `Disallow the use of "alpha", "beta", "rc", etc. version tags`, rule.GetMessage())
}

func TestRuleNoPreRelease_Check(t *testing.T) {
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
			name: "manifest with stable versions",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Definition: types.Definition{Path: "pkg1"},
							Version:    "1.0.0",
						},
						{
							Definition: types.Definition{Path: "pkg2"},
							Version:    "^2.0.0",
						},
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "manifest with alpha version",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Definition: types.Definition{Path: "pkg1"},
							Version:    "1.0.0-alpha",
						},
					},
				},
			},
			want: []Mistake{
				{
					Rule: NewRuleNoPreRelease(),
					Definitions: []types.Definition{
						{Path: "pkg1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "manifest with beta version",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Definition: types.Definition{Path: "pkg1"},
							Version:    "2.0.0-beta.1",
						},
					},
				},
			},
			want: []Mistake{
				{
					Rule: NewRuleNoPreRelease(),
					Definitions: []types.Definition{
						{Path: "pkg1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "manifest with rc version",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Definition: types.Definition{Path: "pkg1"},
							Version:    "1.0.0-rc.2",
						},
					},
				},
			},
			want: []Mistake{
				{
					Rule: NewRuleNoPreRelease(),
					Definitions: []types.Definition{
						{Path: "pkg1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple manifests with mixed versions",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Definition: types.Definition{Path: "pkg1"},
							Version:    "1.0.0-alpha",
						},
						{
							Definition: types.Definition{Path: "pkg2"},
							Version:    "2.0.0",
						},
					},
				},
				{
					Dependencies: []types.Dependency{
						{
							Definition: types.Definition{Path: "pkg3"},
							Version:    "3.0.0-beta",
						},
						{
							Definition: types.Definition{Path: "pkg4"},
							Version:    "4.0.0-rc.1",
						},
						{
							Definition: types.Definition{Path: "pkg5"},
							Version:    "5.0.0",
						},
					},
				},
			},
			want: []Mistake{
				{
					Rule: NewRuleNoPreRelease(),
					Definitions: []types.Definition{
						{Path: "pkg1"},
					},
				},
				{
					Rule: NewRuleNoPreRelease(),
					Definitions: []types.Definition{
						{Path: "pkg3"},
					},
				},
				{
					Rule: NewRuleNoPreRelease(),
					Definitions: []types.Definition{
						{Path: "pkg4"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "version containing pre-release strings in package name",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Definition: types.Definition{Path: "alpha-pkg"},
							Version:    "1.0.0",
						},
						{
							Definition: types.Definition{Path: "beta-pkg"},
							Version:    "2.0.0",
						},
						{
							Definition: types.Definition{Path: "rc-pkg"},
							Version:    "3.0.0",
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
			rule := NewRuleNoPreRelease()
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
