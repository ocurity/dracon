package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Team struct {
	UUID             uuid.UUID     `json:"uuid,omitempty"`
	Name             string        `json:"name,omitempty"`
	APIKeys          []APIKey      `json:"apiKeys,omitempty"`
	Permissions      []Permission  `json:"permissions,omitempty"`
	MappedOIDCGroups []OIDCMapping `json:"mappedOidcGroups,omitempty"`
}

type APIKey struct {
	Key string `json:"key"`
}

type TeamService struct {
	client *Client
}

func (ts TeamService) Get(ctx context.Context, teamUUID uuid.UUID) (t Team, err error) {
	req, err := ts.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/team/%s", teamUUID))
	if err != nil {
		return
	}

	_, err = ts.client.doRequest(req, &t)
	return
}

func (ts TeamService) GetAll(ctx context.Context, po PageOptions) (p Page[Team], err error) {
	req, err := ts.client.newRequest(ctx, http.MethodGet, "/api/v1/team", withPageOptions(po))
	if err != nil {
		return
	}

	res, err := ts.client.doRequest(req, &p.Items)
	if err != nil {
		return
	}

	p.TotalCount = res.TotalCount
	return
}

func (ts TeamService) GenerateAPIKey(ctx context.Context, teamUUID uuid.UUID) (key string, err error) {
	req, err := ts.client.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/team/%s/key", teamUUID))
	if err != nil {
		return
	}

	var apiKey APIKey
	_, err = ts.client.doRequest(req, &apiKey)
	key = apiKey.Key
	return
}

func (ts TeamService) Create(ctx context.Context, team Team) (t Team, err error) {
	req, err := ts.client.newRequest(ctx, http.MethodPut, "/api/v1/team", withBody(team))
	if err != nil {
		return
	}

	_, err = ts.client.doRequest(req, &t)
	return
}

func (ts TeamService) Update(ctx context.Context, team Team) (t Team, err error) {
	req, err := ts.client.newRequest(ctx, http.MethodPost, "/api/v1/team", withBody(team))
	if err != nil {
		return
	}

	_, err = ts.client.doRequest(req, &t)
	return
}

func (ts TeamService) Delete(ctx context.Context, team Team) (err error) {
	req, err := ts.client.newRequest(ctx, http.MethodDelete, "/api/v1/team", withBody(team))
	if err != nil {
		return
	}

	_, err = ts.client.doRequest(req, nil)
	return
}
