package rules

import (
	"testing"

	"github.com/depshubhq/depshub/internal/config"
	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestRuleNoUnstable(t *testing.T) {
	rule := NewRuleNoUnstable()

	// Test rule metadata
	t.Run("metadata", func(t *testing.T) {
		assert.Equal(t, "no-unstable", rule.GetName())
		assert.Equal(t, types.LevelError, rule.GetLevel())
		assert.Equal(t, "Disallow the use of unstable versions (< 1.0.0)", rule.GetMessage())
	})

	tests := []struct {
		name      string
		manifests []types.Manifest
		want      int
		wantErr   bool
	}{
		{
			name: "should detect unstable version",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Version: "0.1.0",
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/test-package v0.1.0`,
								Line:    1,
							},
						},
					},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "should allow stable version",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Version: "1.0.0",
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/test-package v1.0.0`,
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
			name: "should handle multiple dependencies",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Version: "0.9.9",
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/unstable-package v0.9.9`,
								Line:    1,
							},
						},
						{
							Version: "1.0.0",
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/stable-package v1.0.0`,
								Line:    2,
							},
						},
						{
							Version: "0.1.0",
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/another-unstable v0.1.0`,
								Line:    3,
							},
						},
					},
				},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "should handle version with prefix",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Version: "v0.1.0",
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/test-package v0.1.0`,
								Line:    1,
							},
						},
					},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "should handle invalid version format",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Version: "invalid-version",
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/test-package invalid-version`,
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
			name: "should handle empty dependencies",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{},
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name:      "should handle empty manifests",
			manifests: []types.Manifest{},
			want:      0,
			wantErr:   false,
		},
		{
			name: "should handle complex version strings",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Version: "^0.1.0-beta.1",
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/test-package ^0.1.0-beta.1`,
								Line:    1,
							},
						},
						{
							Version: "~1.2.3-alpha+001",
							Definition: types.Definition{
								Path:    "go.mod",
								RawLine: `require github.com/another-package ~1.2.3-alpha+001`,
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
			name: "should handle multiple manifest files",
			manifests: []types.Manifest{
				{
					Dependencies: []types.Dependency{
						{
							Version: "0.9.9",
							Definition: types.Definition{
								Path:    "project1/go.mod",
								RawLine: `require github.com/unstable-package v0.9.9`,
								Line:    1,
							},
						},
					},
				},
				{
					Dependencies: []types.Dependency{
						{
							Version: "0.1.0",
							Definition: types.Definition{
								Path:    "project2/go.mod",
								RawLine: `require github.com/another-unstable v0.1.0`,
								Line:    1,
							},
						},
					},
				},
			},
			want:    2,
			wantErr: false,
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
