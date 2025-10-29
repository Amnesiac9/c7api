package c7api

import (
	"encoding/json"
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
