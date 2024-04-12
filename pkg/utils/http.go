package utils

import "net/http"

// MockableRequestDoer is an interface that can be used wherever we have an
// HTTP client that needs to be mocked to inject custom responses and errors
type MockableRequestDoer interface {
	Do(*http.Request) (*http.Response, error)
}
