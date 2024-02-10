package manifests

import (
	"context"
	"fmt"
	"os"
	"path"
)

var _ Loader = (*localFileLoader)(nil)

type localFileLoader struct {
	path string
}

func newLocalFileLoader(root, fileOrDir, targetFile string) (*localFileLoader, error) {
	if !path.IsAbs(fileOrDir) {
		fileOrDir = path.Join(root, fileOrDir)
	}

	info, err := os.Stat(fileOrDir)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		fileOrDir = path.Join(fileOrDir, targetFile)
	} else {
		if path.Base(fileOrDir) != targetFile {
			return nil, fmt.Errorf("path %s should be pointing to a %s", fileOrDir, targetFile)
		}
	}

	fileOrDir = path.Clean(fileOrDir)

	return &localFileLoader{path: fileOrDir}, nil
}

func (lfl *localFileLoader) Path() string {
	return lfl.path
}

func (lfl *localFileLoader) Load(_ context.Context) ([]byte, error) {
	return os.ReadFile(lfl.path)
}
