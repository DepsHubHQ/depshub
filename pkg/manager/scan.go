package manager

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/depshubhq/depshub/internal/config"
	"github.com/depshubhq/depshub/pkg/manager/cargo"
	"github.com/depshubhq/depshub/pkg/manager/gem"
	gomanager "github.com/depshubhq/depshub/pkg/manager/go"
	"github.com/depshubhq/depshub/pkg/manager/hex"
	"github.com/depshubhq/depshub/pkg/manager/maven"
	"github.com/depshubhq/depshub/pkg/manager/npm"
	"github.com/depshubhq/depshub/pkg/manager/pip"
	"github.com/depshubhq/depshub/pkg/manager/pyproject"
	"github.com/depshubhq/depshub/pkg/types"
	ignore "github.com/sabhiram/go-gitignore"
)

type scanner struct {
	gitignore *ignore.GitIgnore
	config    config.Config
	managers  []Manager
}

func New(config config.Config) scanner {
	return scanner{
		config: config,
		managers: []Manager{
			npm.Npm{},
			gomanager.Go{},
			cargo.Cargo{},
			pip.Pip{},
			hex.Hex{},
			pyproject.Pyproject{},
			maven.Maven{},
			gem.Gem{},
		},
	}
}

func (s scanner) Scan(pathToScan string) ([]types.Manifest, error) {
	var manifests []types.Manifest

	log.Println("Scanning path:", pathToScan)

	// Check if there is a .gitignore file in the root directory
	gitignorePath := filepath.Join(pathToScan, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		s.loadGitignore(gitignorePath)
	}

	err := filepath.WalkDir(pathToScan, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip files matched by .gitignore
		if s.gitignore != nil && s.gitignore.MatchesPath(path) {
			return filepath.SkipDir
		}

		// Check if the path is ignored by the config
		ignored, err := s.config.Ignored(path)

		if err != nil {
			return err
		}

		if ignored {
			return filepath.SkipDir
		}

		dependencies, managerType, err := s.dependencies(path)

		if err != nil {
			return err
		}

		var lockfile *types.Lockfile
		lockfilePath, err := s.lockfilePath(path)

		if err == nil {
			lockfile = &types.Lockfile{
				Path: lockfilePath,
			}
		}

		if len(dependencies) != 0 {
			manifests = append(manifests, types.Manifest{
				Manager:      managerType,
				Path:         path,
				Dependencies: dependencies,
				Lockfile:     lockfile,
			})
		}

		return nil
	})

	return manifests, err
}

func (s scanner) UniqueDependencies(manifests []types.Manifest) (result []types.Dependency) {
	uniqueDependencies := make(map[string]types.Dependency)

	for _, manifest := range manifests {
		for _, dep := range manifest.Dependencies {
			uniqueDependencies[dep.Name] = dep
		}
	}

	for _, dep := range uniqueDependencies {
		result = append(result, dep)
	}

	return result
}

func (s scanner) dependencies(path string) ([]types.Dependency, types.ManagerType, error) {
	for _, m := range s.managers {
		if !m.Managed(path) {
			continue
		}

		dependencies, err := m.Dependencies(path)
		if err != nil {
			return nil, 0, err
		}

		return dependencies, m.GetType(), nil
	}

	return nil, 0, nil
}

func (s scanner) lockfilePath(path string) (string, error) {
	for _, m := range s.managers {
		if m.Managed(path) {
			return m.LockfilePath(path)
		}
	}
	return "", nil
}

func (s *scanner) loadGitignore(path string) error {
	// Ignore if gitignore is already loaded
	if s.gitignore != nil {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	lines = append(lines, ".git", "node_modules", "deps", "_build", "tmp")
	s.gitignore = ignore.CompileIgnoreLines(lines...)
	return nil
}
