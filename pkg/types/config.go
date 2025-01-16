package types

type Config interface {
	Apply(manifestPath string, packageName string, rule Rule) error
}
