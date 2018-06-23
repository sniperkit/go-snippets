package gogurt

type SourceArchive interface {
	Name() string
	URL(version string) string
}
