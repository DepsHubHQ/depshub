package types

type Manifest struct {
	Path         string
	Dependencies []Dependency
}

type Dependency struct {
	Name    string
	Version string
}
