package types

type Manifest struct {
	Path         string
	Dependencies []Dependency
}

type Dependency struct {
	Name    string
	Version string
	Dev     bool
	Definition
}

type Definition struct {
	RawLine string
	Line    int
}
