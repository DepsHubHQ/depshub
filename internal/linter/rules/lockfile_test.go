package rules

import (
	"testing"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewRuleLockfile(t *testing.T) {
	rule := NewRuleLockfile()

	assert.Equal(t, "lockfile", rule.GetName())
	assert.Equal(t, LevelError, rule.GetLevel())
	assert.Equal(t, "The lockfile should be always present", rule.GetMessage())
}

func TestRuleLockfile_Check(t *testing.T) {
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
			name: "manifest with lockfile",
			manifests: []types.Manifest{
				{
					Path: "path/to/manifest",
					Lockfile: &types.Lockfile{
						Path: "path/to/lockfile",
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "manifest without lockfile",
			manifests: []types.Manifest{
				{
					Path:     "path/to/manifest",
					Lockfile: nil,
				},
			},
			want: []Mistake{
				{
					Rule: NewRuleLockfile(),
					Definitions: []types.Definition{
						{
							Path: "path/to/manifest",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple manifests with mixed lockfile presence",
			manifests: []types.Manifest{
				{
					Path:     "path/to/manifest1",
					Lockfile: nil,
				},
				{
					Path: "path/to/manifest2",
					Lockfile: &types.Lockfile{
						Path: "path/to/lockfile2",
					},
				},
				{
					Path:     "path/to/manifest3",
					Lockfile: nil,
				},
			},
			want: []Mistake{
				{
					Rule: NewRuleLockfile(),
					Definitions: []types.Definition{
						{
							Path: "path/to/manifest1",
						},
					},
				},
				{
					Rule: NewRuleLockfile(),
					Definitions: []types.Definition{
						{
							Path: "path/to/manifest3",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewRuleLockfile()
			got, err := rule.Check(tt.manifests, nil)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
