package debtool

// PackageSourceProvider provides a list of package repositories.
type PackageSourceProvider interface {
	GetPackageArchives() []PackageArchive
}
