package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Project struct {
	UUID               uuid.UUID         `json:"uuid,omitempty"`
	Author             string            `json:"author,omitempty"`
	Publisher          string            `json:"publisher,omitempty"`
	Group              string            `json:"group,omitempty"`
	Name               string            `json:"name,omitempty"`
	Description        string            `json:"description,omitempty"`
	Version            string            `json:"version,omitempty"`
	Classifier         string            `json:"classifier,omitempty"`
	CPE                string            `json:"cpe,omitempty"`
	PURL               string            `json:"purl,omitempty"`
	SWIDTagID          string            `json:"swidTagId,omitempty"`
	DirectDependencies string            `json:"directDependencies,omitempty"`
	Properties         []ProjectProperty `json:"properties,omitempty"`
	Tags               []Tag             `json:"tags,omitempty"`
	Active             bool              `json:"active"`
	Metrics            ProjectMetrics    `json:"metrics"`
	LastBOMImport      int               `json:"lastBomImport"`
}

type ProjectService struct {
	client *Client
}

func (ps ProjectService) Get(ctx context.Context, projectUUID uuid.UUID) (p Project, err error) {
	req, err := ps.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/project/%s", projectUUID))
	if err != nil {
		return
	}

	_, err = ps.client.doRequest(req, &p)
	return
}

func (ps ProjectService) GetAll(ctx context.Context, po PageOptions) (p Page[Project], err error) {
	req, err := ps.client.newRequest(ctx, http.MethodGet, "/api/v1/project", withPageOptions(po))
	if err != nil {
		return
	}

	res, err := ps.client.doRequest(req, &p.Items)
	if err != nil {
		return
	}

	p.TotalCount = res.TotalCount
	return
}

func (ps ProjectService) Create(ctx context.Context, project Project) (p Project, err error) {
	req, err := ps.client.newRequest(ctx, http.MethodPut, "/api/v1/project", withBody(project))
	if err != nil {
		return
	}

	_, err = ps.client.doRequest(req, &p)
	return
}

func (ps ProjectService) Patch(ctx context.Context, projectUUID uuid.UUID, project Project) (p Project, err error) {
	req, err := ps.client.newRequest(ctx, http.MethodPatch, fmt.Sprintf("/api/v1/project/%s", projectUUID), withBody(project))
	if err != nil {
		return
	}

	_, err = ps.client.doRequest(req, &p)
	return
}

func (ps ProjectService) Update(ctx context.Context, project Project) (p Project, err error) {
	req, err := ps.client.newRequest(ctx, http.MethodPost, "/api/v1/project", withBody(project))
	if err != nil {
		return
	}

	_, err = ps.client.doRequest(req, &p)
	return
}

func (ps ProjectService) Delete(ctx context.Context, projectUUID uuid.UUID) (err error) {
	req, err := ps.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/project/%s", projectUUID))
	if err != nil {
		return
	}

	_, err = ps.client.doRequest(req, nil)
	return
}

func (ps ProjectService) Lookup(ctx context.Context, name, version string) (p Project, err error) {
	params := map[string]string{
		"name":    name,
		"version": version,
	}

	req, err := ps.client.newRequest(ctx, http.MethodGet, "/api/v1/project/lookup", withParams(params))
	if err != nil {
		return
	}

	_, err = ps.client.doRequest(req, &p)
	return
}

type ProjectCloneRequest struct {
	ProjectUUID         uuid.UUID `json:"project"`
	Version             string    `json:"version"`
	IncludeAuditHistory bool      `json:"includeAuditHistory"`
	IncludeComponents   bool      `json:"includeComponents"`
	IncludeProperties   bool      `json:"includeProperties"`
	IncludeServices     bool      `json:"includeServices"`
	IncludeTags         bool      `json:"includeTags"`
}

func (ps ProjectService) Clone(ctx context.Context, cloneReq ProjectCloneRequest) (err error) {
	req, err := ps.client.newRequest(ctx, http.MethodPut, "/api/v1/project/clone", withBody(cloneReq))
	if err != nil {
		return
	}

	_, err = ps.client.doRequest(req, nil)
	return
}
