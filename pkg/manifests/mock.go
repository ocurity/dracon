package manifests

import "fmt"

type FakeLoader struct {
	path string
}

func (l FakeLoader) Root() string {
	return ""
}

func (l FakeLoader) New(newRoot string) (Loader, error) {
	if newRoot == l.path {
		return nil, nil
	}
	return nil, fmt.Errorf("%s not exist", newRoot)
}
func (l FakeLoader) Load(location string) ([]byte, error) {
	return nil, nil
}
func (l FakeLoader) Cleanup() error {
	return nil
}
