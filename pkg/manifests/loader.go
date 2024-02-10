package manifests

import (
	"context"
)

type Loader interface {
	// Root returns the root location for this Loader.
	Path() string

	// Load returns the bytes read from the location or an error.
	Load(ctx context.Context) ([]byte, error)
}

func NewLoader(root, pathOrURI, targetFile string) (Loader, error) {
	if IsRemoteFile(pathOrURI) {
		return newHTTPFileLoader(pathOrURI)
	}

	return newLocalFileLoader(root, pathOrURI, targetFile)
}
