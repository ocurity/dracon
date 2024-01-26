package dtrack

import (
	"context"
	"net/http"
)

type RepositoryMetaComponent struct {
	LatestVersion string `json:"latestVersion"`
}

type RepositoryService struct {
	client *Client
}

func (rs RepositoryService) GetMetaComponent(ctx context.Context, purl string) (r RepositoryMetaComponent, err error) {
	params := map[string]string{
		"purl": purl,
	}

	req, err := rs.client.newRequest(ctx, http.MethodGet, "/api/v1/repository/latest", withParams(params))
	if err != nil {
		return
	}

	_, err = rs.client.doRequest(req, &r)
	return
}
