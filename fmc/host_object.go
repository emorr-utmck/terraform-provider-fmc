package fmc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HostObjectUpdateInput struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Overridable bool   `json:"overridable"`
	Description string `json:"description"`
	Type        string `json:"type"`
	ID          string `json:"id"`
}

type HostObject struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Overridable bool   `json:"overridable"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type HostObjectResponse struct {
	Links struct {
		Self   string `json:"self"`
		Parent string `json:"parent"`
	} `json:"links"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Overridable bool   `json:"overridable"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Metadata    struct {
		Lastuser struct {
			Name string `json:"name"`
		} `json:"lastUser"`
		Domain struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"domain"`
		Iptype     string `json:"ipType"`
		Parenttype string `json:"parentType"`
	} `json:"metadata"`
}

type HostObjectsResponse struct {
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
	Items []struct {
		Links struct {
			Self   string `json:"self"`
			Parent string `json:"parent"`
		} `json:"links"`
		Type string `json:"type"`
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"items"`
	Paging struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
		Count  int `json:"count"`
		Pages  int `json:"pages"`
	} `json:"paging"`
}

func (v *Client) GetHostObjectByNameOrValue(ctx context.Context, nameOrValue string) (*HostObjectResponse, error) {
	url := fmt.Sprintf("%s/object/hosts?expanded=false&filter=nameOrValue:%s", v.domainBaseURL, nameOrValue)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("getting host object by name/value: %s - %s", url, err.Error())
	}
	resp := &HostObjectsResponse{}
	err = v.DoRequest(req, resp, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("getting host object by name/value: %s - %s", url, err.Error())
	}
	switch l := len(resp.Items); {
	case l > 1:
		return nil, fmt.Errorf("duplicates found, length of response is: %d, expected 1, please search using a unique id, name or value", l)
	case l == 0:
		return nil, fmt.Errorf("no host objects found, length of response is: %d, expected 1, please check your filter", l)
	}
	return v.GetHostObject(ctx, resp.Items[0].ID)
}

// /fmc_config/v1/domain/DomainUUID/object/hosts?bulk=true ( Bulk POST operation on host objects. )

func (v *Client) CreateHostObject(ctx context.Context, object *HostObject) (*HostObjectResponse, error) {
	url := fmt.Sprintf("%s/object/hosts", v.domainBaseURL)
	body, err := json.Marshal(&object)
	if err != nil {
		return nil, fmt.Errorf("creating host objects: %s - %s", url, err.Error())
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("creating host objects: %s - %s", url, err.Error())
	}
	item := &HostObjectResponse{}
	err = v.DoRequest(req, item, http.StatusCreated)
	if err != nil {
		return nil, fmt.Errorf("getting host objects: %s - %s", url, err.Error())
	}
	return item, nil
}

func (v *Client) GetHostObject(ctx context.Context, id string) (*HostObjectResponse, error) {
	url := fmt.Sprintf("%s/object/hosts/%s", v.domainBaseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("getting host objects: %s - %s", url, err.Error())
	}
	item := &HostObjectResponse{}
	err = v.DoRequest(req, item, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("getting host objects: %s - %s", url, err.Error())
	}
	return item, nil
}

func (v *Client) UpdateHostObject(ctx context.Context, id string, object *HostObjectUpdateInput) (*HostObjectResponse, error) {
	url := fmt.Sprintf("%s/object/hosts/%s", v.domainBaseURL, id)
	body, err := json.Marshal(&object)
	if err != nil {
		return nil, fmt.Errorf("updating host objects: %s - %s", url, err.Error())
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("updating host objects: %s - %s", url, err.Error())
	}
	item := &HostObjectResponse{}
	err = v.DoRequest(req, item, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("getting host objects: %s - %s", url, err.Error())
	}
	return item, nil
}

func (v *Client) DeleteHostObject(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/object/hosts/%s", v.domainBaseURL, id)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("deleting host objects: %s - %s", url, err.Error())
	}
	err = v.DoRequest(req, nil, http.StatusOK)
	return err
}
