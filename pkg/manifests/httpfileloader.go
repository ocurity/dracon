package manifests

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var _ Loader = (*httpFileLoader)(nil)

type httpFileLoader struct {
	uri    string
	client *http.Client
}

func newHTTPFileLoader(uri string) (*httpFileLoader, error) {
	return &httpFileLoader{uri: uri}, nil
}

func (hfl *httpFileLoader) Load(ctx context.Context) (content []byte, err error) {
	if hfl.client == nil {
		hfl.client = &http.Client{}
	}
	resp, err := hfl.client.Get(hfl.uri)
	if err != nil {
		return nil, fmt.Errorf("%s: could not create request: %w", hfl.uri, err)
	}
	defer func() {
		err = errors.Join(err, resp.Body.Close())
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("%s: could not fetch file: %d", hfl.uri, resp.StatusCode)
	}

	content, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: could not read body of response: %w", hfl.uri, err)
	}
	return content, nil
}

func (hfl *httpFileLoader) Path() string {
	return hfl.uri
}
