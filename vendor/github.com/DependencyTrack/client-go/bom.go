package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type BOMService struct {
	client *Client
}

type BOMUploadRequest struct {
	ProjectUUID    *uuid.UUID `json:"project,omitempty"`
	ProjectName    string     `json:"projectName,omitempty"`
	ProjectVersion string     `json:"projectVersion,omitempty"`
	AutoCreate     bool       `json:"autoCreate"`
	BOM            string     `json:"bom"`
}

type bomUploadResponse struct {
	Token BOMUploadToken `json:"token"`
}

type BOMUploadToken string

type BOMFormat string

const (
	BOMFormatJSON BOMFormat = "JSON"
	BOMFormatXML  BOMFormat = "XML"
)

type BOMVariant string

const (
	BOMVariantInventory           BOMVariant = "inventory"
	BOMVariantWithVulnerabilities BOMVariant = "withVulnerabilities"
)

func (bs BOMService) ExportComponent(ctx context.Context, componentUUID uuid.UUID, format BOMFormat) (bom string, err error) {
	params := make(map[string]string)
	if format != "" {
		params["format"] = string(format)
	}

	req, err := bs.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/bom/cyclonedx/component/%s", componentUUID), withParams(params))
	if err != nil {
		return
	}

	req.Header.Set("Accept", "application/vnd.cyclonedx+json")

	_, err = bs.client.doRequest(req, &bom)
	return
}

func (bs BOMService) ExportProject(ctx context.Context, projectUUID uuid.UUID, format BOMFormat, variant BOMVariant) (bom string, err error) {
	params := make(map[string]string)
	if format != "" {
		params["format"] = string(format)
	}
	if variant != "" {
		params["variant"] = string(variant)
	}

	req, err := bs.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/bom/cyclonedx/project/%s", projectUUID), withParams(params))
	if err != nil {
		return
	}

	req.Header.Set("Accept", "application/vnd.cyclonedx+json")

	_, err = bs.client.doRequest(req, &bom)
	return
}

func (bs BOMService) Upload(ctx context.Context, uploadReq BOMUploadRequest) (token BOMUploadToken, err error) {
	req, err := bs.client.newRequest(ctx, http.MethodPut, "/api/v1/bom", withBody(uploadReq))
	if err != nil {
		return
	}

	var uploadRes bomUploadResponse
	_, err = bs.client.doRequest(req, &uploadRes)
	if err != nil {
		return
	}

	token = uploadRes.Token
	return
}

type bomProcessingResponse struct {
	Processing bool `json:"processing"`
}

func (bs BOMService) IsBeingProcessed(ctx context.Context, token BOMUploadToken) (bool, error) {
	req, err := bs.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/bom/token/%s", token))
	if err != nil {
		return false, err
	}

	var processingResponse bomProcessingResponse
	_, err = bs.client.doRequest(req, &processingResponse)
	if err != nil {
		return false, err
	}

	return processingResponse.Processing, nil
}
