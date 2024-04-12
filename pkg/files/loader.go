package files

import (
	"context"
	"strings"

	"github.com/go-errors/errors"
)

// Loader is an interface implemented by a struct that can resolve a path or a
// URI and return the contents as a byte slice.
type Loader interface {
	// Root returns the root location for this Loader.
	Path() string

	// Load returns the bytes read from the location or an error.
	Load(ctx context.Context) ([]byte, error)
}

// ErrBadTargetFile is returned when the target file is not a hint for a file
// but a directory.
var ErrBadTargetFile = errors.New("target file should be a filename not a directory")

// NewLoader takes as an argumet a directory containing some configuration, a
// path or a URI and the expected file name that we want to fetch. If the
func NewLoader(configurationDir, pathOrURI, targetFile string) (Loader, error) {
	if strings.Contains(targetFile, "/") {
		return nil, errors.Errorf("%s: %w", targetFile, ErrBadTargetFile)
	}

	if IsRemoteFile(pathOrURI) {
		return newHTTPFileLoader(pathOrURI, targetFile)
	}

	return newLocalFileLoader(configurationDir, pathOrURI, targetFile)
}
