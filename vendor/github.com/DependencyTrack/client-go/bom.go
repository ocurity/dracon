package dtrack

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type BOMService struct {
	client *Client
}

type BOMUploadRequest struct {
	ProjectUUID    *uuid.UUID `json:"project,omitempty"`
	ProjectName    string     `json:"projectName,omitempty"`
	ProjectVersion string     `json:"projectVersion,omitempty"`
	ParentUUID     *uuid.UUID `json:"parentUUID,omitempty"`    // Since v4.8.0
	ParentName     string     `json:"parentName,omitempty"`    // Since v4.8.0
	ParentVersion  string     `json:"parentVersion,omitempty"` // Since v4.8.0
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
	BOMVariantVDR                 BOMVariant = "vdr" // Since v4.7.0
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

func (bs BOMService) PostBom(ctx context.Context, uploadReq BOMUploadRequest) (token BOMUploadToken, err error) {
	params := make(url.Values)
	if uploadReq.ProjectUUID != nil {
		params["project"] = append(params["project"], uploadReq.ProjectUUID.String())
	}
	if uploadReq.AutoCreate {
		params["autoCreate"] = append(params["autoCreate"], "true")
	}
	if uploadReq.ProjectName != "" {
		params["projectName"] = append(params["projectName"], uploadReq.ProjectName)
	}
	if uploadReq.ProjectVersion != "" {
		params["projectVersion"] = append(params["projectVersion"], uploadReq.ProjectVersion)
	}
	if uploadReq.ParentUUID != nil {
		params["parentUUID"] = append(params["parentUUID"], uploadReq.ParentUUID.String())
	}
	if uploadReq.ParentName != "" {
		params["parentName"] = append(params["parentName"], uploadReq.ParentName)
	}
	if uploadReq.ParentVersion != "" {
		params["parentVersion"] = append(params["parentVersion"], uploadReq.ParentVersion)
	}
	if uploadReq.BOM != "" {
		params["bom"] = append(params["bom"], uploadReq.BOM)
	}

	req, err := bs.client.newRequest(ctx, http.MethodPost, "/api/v1/bom", withMultiPart(params))
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
