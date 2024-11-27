package types

type Manifest struct {
	Path         string
	Dependencies []Dependency
	*Lockfile
}

type Lockfile struct {
	Path string
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
