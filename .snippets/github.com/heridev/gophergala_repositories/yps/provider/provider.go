// Package provider contains the interface definition for what a provider should implement as functions
package provider

type (
	// DownloadFunc defines the logic of fetching the URL specific to a provider
	DownloadFunc func(string) ([]byte, error)

	// Provider interface determines the common functions across providers
	Provider interface {
		IsValidURL(string) bool
		IsVideo(string) bool
		IsPlaylist(string) bool
		URLToFile(string, DownloadFunc) (string, error)
	}
)
