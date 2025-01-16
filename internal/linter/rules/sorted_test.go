package rules

import (
	"testing"

	"github.com/depshubhq/depshub/internal/config"
	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestRuleSorted(t *testing.T) {
	rule := NewRuleSorted()

	// Test rule metadata
	t.Run("metadata", func(t *testing.T) {
		assert.Equal(t, "sorted", rule.GetName())
		assert.Equal(t, types.LevelError, rule.GetLevel())
		assert.Equal(t, "All the dependencies should be ordered alphabetically", rule.GetMessage())
	})

	tests := []struct {
		name      string
		manifests []types.Manifest
		want      int
		wantErr   bool
	}{
		{
			name: "properly sorted dependencies",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name: "alpha",
							Dev:  false,
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/alpha v1.0.0`,
								Line:    1,
							},
						},
						{
							Name: "beta",
							Dev:  false,
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/beta v1.0.0`,
								Line:    2,
							},
						},
						{
							Name: "gamma",
							Dev:  false,
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/gamma v1.0.0`,
								Line:    3,
							},
						},
					},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "unsorted dependencies",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name: "beta",
							Dev:  false,
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/beta v1.0.0`,
								Line:    1,
							},
						},
						{
							Name: "alpha",
							Dev:  false,
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/alpha v1.0.0`,
								Line:    2,
							},
						},
					},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "mixed dev and non-dev dependencies",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name: "zeta",
							Dev:  false,
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/zeta v1.0.0`,
								Line:    1,
							},
						},
						{
							Name: "alpha",
							Dev:  true,
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/alpha v1.0.0 // dev`,
								Line:    2,
							},
						},
					},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "unsorted dev dependencies",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name: "charlie",
							Dev:  true,
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/charlie v1.0.0 // dev`,
								Line:    1,
							},
						},
						{
							Name: "alpha",
							Dev:  true,
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/alpha v1.0.0 // dev`,
								Line:    2,
							},
						},
					},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "multiple manifest files",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name: "beta",
							Dev:  false,
							Definition: types.Definition{
								Path:    "project1/go.mod",
								RawLine: `require github.com/beta v1.0.0`,
								Line:    1,
							},
						},
						{
							Name: "alpha",
							Dev:  false,
							Definition: types.Definition{
								Path:    "project1/go.mod",
								RawLine: `require github.com/alpha v1.0.0`,
								Line:    2,
							},
						},
					},
				},
				{
					Dependencies: []types.Dependency{
						{
							Name: "delta",
							Dev:  false,
							Definition: types.Definition{
								Path:    "project2/go.mod",
								RawLine: `require github.com/delta v1.0.0`,
								Line:    1,
							},
						},
						{
							Name: "gamma",
							Dev:  false,
							Definition: types.Definition{
								Path:    "project2/go.mod",
								RawLine: `require github.com/gamma v1.0.0`,
								Line:    2,
							},
						},
					},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "single dependency",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Name: "alpha",
							Dev:  false,
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/alpha v1.0.0`,
								Line:    1,
							},
						},
					},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "empty dependencies",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name:      "empty manifests",
			manifests: []types.Manifest{},
			want:      0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mistakes, err := rule.Check(tt.manifests, nil, config.Config{})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, mistakes, tt.want)

			// Verify mistake details when mistakes are expected
			if tt.want > 0 {
				for _, mistake := range mistakes {
					assert.NotEmpty(t, mistake.Definitions)

					// Verify Definition fields
					for _, def := range mistake.Definitions {
						assert.NotEmpty(t, def.Path)
						assert.NotEmpty(t, def.RawLine)
						assert.Greater(t, def.Line, 0)
					}
				}
			}
		})
	}
}
