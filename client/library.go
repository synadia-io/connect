package client

import (
    "fmt"
    "github.com/synadia-io/connect/model"
    "time"
)

type libraryEntityKind string

const (
    runtimes   libraryEntityKind = "RUNTIMES"
    components libraryEntityKind = "COMPONENTS"
)

type libraryClient struct {
    t *Transport
}

func (c *libraryClient) subject(kind libraryEntityKind, suffix string) string {
    return fmt.Sprintf("$CONLIB.%s.%s", kind, suffix)
}

func (c *libraryClient) ListRuntimes(timeout time.Duration) ([]model.RuntimeSummary, error) {
    req := model.RuntimeListRequest{}
    var resp model.RuntimeListResponse
    gotResponse, err := c.t.RequestJson(c.subject(runtimes, "LIST"), req, &resp, WithTimeout(timeout))
    if err != nil {
        return nil, fmt.Errorf("unable to list runtimes: %v", err)
    }

    if !gotResponse {
        return []model.RuntimeSummary{}, nil
    }

    return resp.Runtimes, nil
}

func (c *libraryClient) GetRuntime(id string, timeout time.Duration) (*model.Runtime, error) {
    req := model.RuntimeGetRequest{
        Name: id,
    }
    var resp model.RuntimeGetResponse
    gotResponse, err := c.t.RequestJson(c.subject(runtimes, "GET"), req, &resp, WithTimeout(timeout))
    if err != nil {
        return nil, fmt.Errorf("unable to get runtime: %v", err)
    }

    if !gotResponse {
        return nil, nil
    }

    return &resp.Runtime, nil
}

func (c *libraryClient) SearchComponents(filter *model.ComponentSearchFilter, timeout time.Duration) ([]model.ComponentSummary, error) {
    req := model.ComponentSearchRequest{
        Filter: filter,
    }

    var resp model.ComponentSearchResponse
    gotResponse, err := c.t.RequestJson(c.subject(components, "SEARCH"), req, &resp, WithTimeout(timeout))
    if err != nil {
        return nil, fmt.Errorf("unable to search components: %v", err)
    }

    if !gotResponse {
        return []model.ComponentSummary{}, nil
    }

    return resp.Components, nil
}

func (c *libraryClient) GetComponent(runtimeId string, kind model.ComponentKind, id string, timeout time.Duration) (*model.Component, error) {
    req := model.ComponentGetRequest{
        RuntimeId: runtimeId,
        Kind:      kind,
        Name:      id,
    }

    var resp model.ComponentGetResponse
    gotResponse, err := c.t.RequestJson(c.subject(components, "GET"), req, &resp, WithTimeout(timeout))
    if err != nil {
        return nil, fmt.Errorf("unable to get component: %v", err)
    }

    if !gotResponse {
        return nil, nil
    }

    return &resp.Component, nil
}
