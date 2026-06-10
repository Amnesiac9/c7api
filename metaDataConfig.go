package c7api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type MetaDataConfig struct {
	Id         string   `json:"id"`
	Title      string   `json:"title"`
	ObjectType string   `json:"objectType"`
	Code       string   `json:"code"`
	DataType   string   `json:"dataType"`
	IsRequired bool     `json:"isRequired"`
	Options    []string `json:"options"`
	SortOrder  int      `json:"sortOrder"`
	CreatedAt  string   `json:"createdAt"`
	UpdatedAt  string   `json:"updatedAt"`
}

type MetaDataConfigPayload struct {
	MetaDataConfigs []MetaDataConfig `json:"metaDataConfigs"`
	Total           int              `json:"total"`
}

// New meta data / custom field
type MetaDataConfigPost struct {
	Title      string   `json:"title"`
	Code       string   `json:"code"`
	DataType   string   `json:"dataType"`
	IsRequired bool     `json:"isRequired"`
	SortOrder  int      `json:"sortOrder"`
	Options    []string `json:"options"`
}

func GetMetaDataConfigs(tenant, obj, q, c7appAuthEncoded string, retryCount int, rl genericRateLimiter) (*MetaDataConfigPayload, error) {
	reqUrl := Endpoints.MetaDataConfig + url.PathEscape(obj)
	resp, err := RequestWithRetryAndRead(http.MethodGet, reqUrl, nil, nil, tenant, c7appAuthEncoded, retryCount, rl)
	if err != nil {
		return nil, err
	}

	var metaDataConfigPayload MetaDataConfigPayload
	err = json.Unmarshal(*resp, &metaDataConfigPayload)
	if err != nil {
		return nil, err
	}
	return &metaDataConfigPayload, nil
}

func GetMetaDataConfigById(metadataId, objectType, tenant, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*MetaDataConfig, error) {

	if !IsValidMetaDataConfigObjectType(objectType) {
		return nil, fmt.Errorf("not a valid object type for the metadata endpoint")
	}

	reqUrl := Endpoints.MetaDataConfig + objectType + "/" + metadataId

	resp, err := RequestWithRetryAndRead(http.MethodGet, reqUrl, nil, nil, tenant, c7AppAuthEncoded, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	var c7MetaData MetaDataConfig
	if err := json.Unmarshal(*resp, &c7MetaData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &c7MetaData, nil
}

// Request
//
// {"title":"d2","code":"d2","dataType":"Select","isRequired":false,"sortOrder":1,"options":["21e2","2323"]}
//
// # Response on success
//
// {"id":"05c64236-d697-42dc-a3d7-bdb96774e4a2","title":"d2","objectType":"Customer","code":"d2","dataType":"Select","isRequired":false,"options":["21e2","2323"],"sortOrder":1,"createdAt":"2026-06-10T04:37:01.830Z","updatedAt":"2026-06-10T04:37:01.830Z"}
func PostMetaDataConfig(objectPayload *MetaDataConfigPost, objectType, objectId, tenant, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*MetaDataConfig, error) {

	if !IsValidMetaDataConfigObjectType(objectType) {
		return nil, fmt.Errorf("not a valid object type for the metadata endpoint: %s", objectType)
	}

	objectBytes, err := json.Marshal(objectPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata object payload: %w", err)
	}

	reqUrl := Endpoints.MetaDataConfig + objectType + "/" + objectId

	resp, err := RequestWithRetryAndRead(http.MethodPost, reqUrl, nil, &objectBytes, tenant, c7AppAuthEncoded, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	var c7MetaData MetaDataConfig
	if err := json.Unmarshal(*resp, &c7MetaData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata after attempted post: %w", err)
	}

	return &c7MetaData, nil
}

func DeleteMetaDataConfigById(metadataId, objectType, tenant, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) error {

	if !IsValidMetaDataConfigObjectType(objectType) {
		return fmt.Errorf("not a valid object type for the metadata endpoint")
	}

	reqUrl := Endpoints.MetaDataConfig + objectType + "/" + metadataId

	_, err := RequestWithRetryAndRead(http.MethodDelete, reqUrl, nil, nil, tenant, c7AppAuthEncoded, retryCount, rl)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	return nil
}
