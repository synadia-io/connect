package client

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/synadia-io/connect/model"
)

type connectorClient struct {
	t *Transport
}

func (c *connectorClient) subject(suffix string) string {
	return fmt.Sprintf("$CONSVC.%s.CONNECTORS.%s", c.t.Account(), suffix)
}

func (c *connectorClient) ListConnectors(timeout time.Duration) ([]model.ConnectorSummary, error) {
	req := model.ConnectorListRequest{}

	var resp model.ConnectorListResponse
	gotResponse, err := c.t.RequestJson(c.subject("LIST"), req, &resp, WithTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("unable to list connectors: %v", err)
	}

	if !gotResponse {
		return []model.ConnectorSummary{}, nil
	}

	result := resp.Connectors
	slices.SortFunc(result, func(a, b model.ConnectorSummary) int {
		return strings.Compare(a.ConnectorId, b.ConnectorId)
	})

	return result, nil
}

func (c *connectorClient) GetConnector(name string, timeout time.Duration) (*model.Connector, error) {
	req := model.ConnectorGetRequest{
		Id: name,
	}

	var resp model.ConnectorGetResponse
	gotResponse, err := c.t.RequestJson(c.subject("GET"), req, &resp, WithTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("unable to get connector: %v", err)
	}

	if !gotResponse {
		return nil, nil
	}

	return resp.Connector, nil
}

func (c *connectorClient) GetConnectorStatus(name string, timeout time.Duration) (*model.ConnectorStatus, error) {
	req := model.ConnectorStatusRequest{
		ConnectorId: name,
	}

	var resp model.ConnectorStatusResponse
	gotResponse, err := c.t.RequestJson(c.subject("STATUS"), req, &resp, WithTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("unable to get connector status: %v", err)
	}

	if !gotResponse {
		return nil, nil
	}

	return &resp.Status, nil
}

func (c *connectorClient) CreateConnector(id, description, runtimeId string, steps model.Steps, timeout time.Duration) (*model.Connector, error) {
	req := model.ConnectorCreateRequest{
		Id:          id,
		Description: description,
		RuntimeId:   runtimeId,
		Steps:       steps,
	}

	var resp model.ConnectorCreateResponse
	gotResponse, err := c.t.RequestJson(c.subject("CREATE"), req, &resp, WithTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("unable to create connector: %v", err)
	}

	if !gotResponse {
		return nil, nil
	}

	return &resp.Connector, nil
}

func (c *connectorClient) PatchConnector(id string, patch string, timeout time.Duration) (*model.Connector, error) {
	req := model.ConnectorPatchRequest{
		ConnectorId: id,
		Patch:       patch,
	}

	var resp model.ConnectorPatchResponse
	gotResponse, err := c.t.RequestJson(c.subject("PATCH"), req, &resp, WithTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("unable to patch connector: %v", err)
	}

	if !gotResponse {
		return nil, nil
	}

	return &resp.Connector, nil
}

func (c *connectorClient) DeleteConnector(id string, timeout time.Duration) error {
	req := model.ConnectorDeleteRequest{
		Id: id,
	}

	var resp model.ConnectorDeleteResponse
	_, err := c.t.RequestJson(c.subject("DELETE"), req, &resp, WithTimeout(timeout))
	if err != nil {
		return fmt.Errorf("unable to delete connector: %v", err)
	}

	return nil
}

func (c *connectorClient) ListConnectorInstances(id string, timeout time.Duration) ([]model.Instance, error) {
	req := model.ConnectorInstancesRequest{
		ConnectorId: &id,
	}

	var resp model.ConnectorInstancesResponse
	gotResponse, err := c.t.RequestJson(c.subject("INSTANCES"), req, &resp, WithTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("unable to list connector instances: %v", err)
	}

	if !gotResponse {
		return []model.Instance{}, nil
	}

	return resp.Instances, nil
}

func (c *connectorClient) StartConnector(id string, startOpts *model.ConnectorStartOptions, timeout time.Duration) ([]model.Instance, error) {
	req := model.ConnectorStartRequest{
		ConnectorId: id,
		Options:     startOpts,
	}

	var resp model.ConnectorStartResponse
	hasResponded, err := c.t.RequestJson(c.subject("START"), req, &resp, WithTimeout(timeout))
	if err != nil {
		return nil, err
	}

	if !hasResponded {
		return nil, nil
	}

	return resp.Instances, nil
}

func (c *connectorClient) StopConnector(id string, timeout time.Duration) ([]model.Instance, error) {
	req := model.ConnectorStopRequest{
		ConnectorId: id,
	}

	var resp model.ConnectorStopResponse
	hasResponded, err := c.t.RequestJson(c.subject("STOP"), req, &resp, WithTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("unable to stop connector: %v", err)
	}

	if !hasResponded {
		return nil, nil
	}

	return resp.Instances, nil
}
