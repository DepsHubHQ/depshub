package rules

import (
	"testing"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewRuleNoMultipleVersions(t *testing.T) {
	rule := NewRuleNoMultipleVersions()

	assert.Equal(t, "no-multiple-versions", rule.GetName())
	assert.Equal(t, LevelError, rule.GetLevel())
	assert.Equal(t, "Disallow the use of multiple versions of the same package", rule.GetMessage())
}

func TestRuleNoMultipleVersions_Check(t *testing.T) {
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
			name: "single manifest with no version conflicts",
			manifests: []types.Manifest{
				{
					Path: "manifest1",
					Dependencies: []types.Dependency{
						{
							Name:    "pkg1",
							Version: "1.0.0",
							Definition: types.Definition{
								RawLine: "pkg1@1.0.0",
								Line:    1,
							},
						},
						{
							Name:    "pkg2",
							Version: "2.0.0",
							Definition: types.Definition{
								RawLine: "pkg2@2.0.0",
								Line:    2,
							},
						},
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "single manifest with version conflict",
			manifests: []types.Manifest{
				{
					Path: "manifest1",
					Dependencies: []types.Dependency{
						{
							Name:    "pkg1",
							Version: "1.0.0",
							Definition: types.Definition{
								RawLine: "pkg1@1.0.0",
								Line:    1,
							},
						},
						{
							Name:    "pkg1",
							Version: "2.0.0",
							Definition: types.Definition{
								RawLine: "pkg1@2.0.0",
								Line:    2,
							},
						},
					},
				},
			},
			want: []Mistake{
				{
					Rule: NewRuleNoMultipleVersions(),
					Definitions: []types.Definition{
						{
							Path:    "manifest1",
							RawLine: "pkg1@1.0.0",
							Line:    1,
						},
						{
							Path:    "manifest1",
							RawLine: "pkg1@2.0.0",
							Line:    2,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple manifests with version conflicts",
			manifests: []types.Manifest{
				{
					Path: "manifest1",
					Dependencies: []types.Dependency{
						{
							Name:    "pkg1",
							Version: "1.0.0",
							Definition: types.Definition{
								RawLine: "pkg1@1.0.0",
								Line:    1,
							},
						},
					},
				},
				{
					Path: "manifest2",
					Dependencies: []types.Dependency{
						{
							Name:    "pkg1",
							Version: "2.0.0",
							Definition: types.Definition{
								RawLine: "pkg1@2.0.0",
								Line:    1,
							},
						},
					},
				},
			},
			want: []Mistake{
				{
					Rule: NewRuleNoMultipleVersions(),
					Definitions: []types.Definition{
						{
							Path:    "manifest1",
							RawLine: "pkg1@1.0.0",
							Line:    1,
						},
						{
							Path:    "manifest2",
							RawLine: "pkg1@2.0.0",
							Line:    1,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple manifests with multiple version conflicts",
			manifests: []types.Manifest{
				{
					Path: "manifest1",
					Dependencies: []types.Dependency{
						{
							Name:    "pkg1",
							Version: "1.0.0",
							Definition: types.Definition{
								RawLine: "pkg1@1.0.0",
								Line:    1,
							},
						},
						{
							Name:    "pkg2",
							Version: "1.0.0",
							Definition: types.Definition{
								RawLine: "pkg2@1.0.0",
								Line:    2,
							},
						},
					},
				},
				{
					Path: "manifest2",
					Dependencies: []types.Dependency{
						{
							Name:    "pkg1",
							Version: "2.0.0",
							Definition: types.Definition{
								RawLine: "pkg1@2.0.0",
								Line:    1,
							},
						},
						{
							Name:    "pkg2",
							Version: "2.0.0",
							Definition: types.Definition{
								RawLine: "pkg2@2.0.0",
								Line:    2,
							},
						},
					},
				},
			},
			want: []Mistake{
				{
					Rule: NewRuleNoMultipleVersions(),
					Definitions: []types.Definition{
						{
							Path:    "manifest1",
							RawLine: "pkg1@1.0.0",
							Line:    1,
						},
						{
							Path:    "manifest2",
							RawLine: "pkg1@2.0.0",
							Line:    1,
						},
					},
				},
				{
					Rule: NewRuleNoMultipleVersions(),
					Definitions: []types.Definition{
						{
							Path:    "manifest1",
							RawLine: "pkg2@1.0.0",
							Line:    2,
						},
						{
							Path:    "manifest2",
							RawLine: "pkg2@2.0.0",
							Line:    2,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "same version in different manifests",
			manifests: []types.Manifest{
				{
					Path: "manifest1",
					Dependencies: []types.Dependency{
						{
							Name:    "pkg1",
							Version: "1.0.0",
							Definition: types.Definition{
								RawLine: "pkg1@1.0.0",
								Line:    1,
							},
						},
					},
				},
				{
					Path: "manifest2",
					Dependencies: []types.Dependency{
						{
							Name:    "pkg1",
							Version: "1.0.0",
							Definition: types.Definition{
								RawLine: "pkg1@1.0.0",
								Line:    1,
							},
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
			rule := NewRuleNoMultipleVersions()
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
