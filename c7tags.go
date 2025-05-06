package c7api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Tag struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	ObjectType string `json:"objectType"`
}

type TagPayload_Post struct {
	ObjectId string `json:"objectId"`
	TagId    string `json:"tagId"`
}

type TagPayload_Get struct {
	Tags  []Tag `json:"tags"`
	Total int   `json:"total"`
}

func (tagPayload *TagPayload_Post) ToString() string {
	return fmt.Sprintf("ObjectId: %v | TagId: %v", tagPayload.ObjectId, tagPayload.TagId)
}

// allowed object types: "order" | "customer"
// Pass in
func AddTagById(tenant string, auth string, tagId string, targetId string, targetObjType string, retryCount int, rl genericRateLimiter) error {

	targetObjType = strings.ToLower(targetObjType)
	if targetObjType != "order" && targetObjType != "customer" {
		return errors.New("invalid object type for add tag. Must be either: order || customer")
	}

	if tagId == "" {
		return fmt.Errorf("no tagId provided")
	}

	tagPayload := TagPayload_Post{
		ObjectId: targetId,
		TagId:    tagId,
	}

	tagPayloadBytes, err := json.Marshal(tagPayload)
	if err != nil {
		return fmt.Errorf("while unmarshalling payload: %w | Tag Payload: %s", err, tagPayload.ToString())
	}

	// Does this url work with orders? Yes...
	urlt := strings.Replace(Endpoints.TagXObject, "{:obj}", targetObjType, 1)
	_, err = RequestWithRetryAndRead("POST", urlt, nil, &tagPayloadBytes, tenant, auth, retryCount, rl)
	if err != nil {
		return fmt.Errorf("while posting tag: %w | Tag Payload: %s", err, tagPayload.ToString())
	}

	return nil
}

// allowed object types: "order" | "customer"
// Pass in rl or nil as required.
func RemoveTagById(tenant string, auth string, tagId string, targetId string, targetObjType string, retryCount int, rl genericRateLimiter) error {
	if targetObjType != "order" && targetObjType != "customer" {
		return errors.New("invalid object type for add tag. Must be either: order || customer")
	}

	if tagId == "" {
		return fmt.Errorf("no tagId provided")
	}

	if targetId == "" {
		return fmt.Errorf("no targetId provided")
	}

	//https://api.commerce7.com/v1/tag-x-object/customer/{tagid}/{orderid}
	//https://api.commerce7.com/v1/tag-x-object/customer/0f464186-4985-4737-bcc5-f5c33be0a591/d23cb84a-31c7-4a94-83c6-c6086fc48984
	urlt := strings.Replace(Endpoints.TagXObject, "{:obj}", targetObjType, 1) + "/" + targetId
	_, err := RequestWithRetryAndRead("DELETE", urlt, nil, nil, tenant, auth, retryCount, rl)
	if err != nil {
		return fmt.Errorf("while posting tag: %w | Tag Payload: %s", err, urlt)
	}

	return nil
}

// allowed object types: "order" | "customer"
//
// Pass in raw search string
func GetTags(tenant string, auth string, objectType string, query string, rl genericRateLimiter) (*TagPayload_Get, error) {
	// Lowercase and validate
	objectType = strings.ToLower(objectType)
	if objectType != "order" && objectType != "customer" {
		return nil, errors.New("invalid object type for add tag. Must be either: order || customer")
	}

	// Create url and request
	escapedQuery := url.QueryEscape(query)
	urlt := fmt.Sprintf("%s/%s?q=%s", Endpoints.Tag, objectType, escapedQuery)
	tagsResp, err := RequestWithRetryAndRead("GET", urlt, nil, nil, tenant, auth, 2, rl)
	if err != nil {
		return nil, fmt.Errorf("while getting tags: %w", err)
	}

	// unmarshall the tags
	tags := TagPayload_Get{}
	err = json.Unmarshal(*tagsResp, &tags)
	if err != nil {
		return nil, fmt.Errorf("while unmarshalling tags payload: %w", err)
	}

	return &tags, nil
}
