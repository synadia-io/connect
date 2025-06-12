package cli

import (
	"time"

	"github.com/synadia-io/connect/model"
)

// mockClient implements the client.Client interface for testing
type mockClient struct {
	// ConnectorClient methods
	connectors       []model.ConnectorSummary
	connector        *model.Connector
	connectorStatus  *model.ConnectorStatus
	connectorError   error
	instances        []model.Instance
	startOptions     *model.ConnectorStartOptions
	createCalled     bool
	deleteCalled     bool
	patchCalled      bool
	startCalled      bool
	stopCalled       bool
	
	// LibraryClient methods
	runtimes   []model.RuntimeSummary
	runtime    *model.Runtime
	components []model.ComponentSummary
	component  *model.Component
	
	// Client methods
	account string
}

func newMockClient() *mockClient {
	return &mockClient{
		account: "test-account",
		connectors: []model.ConnectorSummary{},
		instances: []model.Instance{},
		runtimes: []model.RuntimeSummary{},
		components: []model.ComponentSummary{},
	}
}

// Client interface
func (m *mockClient) Account() string {
	return m.account
}

func (m *mockClient) Close() {}

// ConnectorClient interface
func (m *mockClient) ListConnectors(timeout time.Duration) ([]model.ConnectorSummary, error) {
	if m.connectorError != nil {
		return nil, m.connectorError
	}
	return m.connectors, nil
}

func (m *mockClient) GetConnector(id string, timeout time.Duration) (*model.Connector, error) {
	if m.connectorError != nil {
		return nil, m.connectorError
	}
	return m.connector, nil
}

func (m *mockClient) GetConnectorStatus(id string, timeout time.Duration) (*model.ConnectorStatus, error) {
	if m.connectorError != nil {
		return nil, m.connectorError
	}
	return m.connectorStatus, nil
}

func (m *mockClient) CreateConnector(id, description, runtimeId string, steps model.Steps, timeout time.Duration) (*model.Connector, error) {
	m.createCalled = true
	if m.connectorError != nil {
		return nil, m.connectorError
	}
	return m.connector, nil
}

func (m *mockClient) PatchConnector(id string, patch string, timeout time.Duration) (*model.Connector, error) {
	m.patchCalled = true
	if m.connectorError != nil {
		return nil, m.connectorError
	}
	return m.connector, nil
}

func (m *mockClient) DeleteConnector(id string, timeout time.Duration) error {
	m.deleteCalled = true
	return m.connectorError
}

func (m *mockClient) ListConnectorInstances(id string, timeout time.Duration) ([]model.Instance, error) {
	if m.connectorError != nil {
		return nil, m.connectorError
	}
	return m.instances, nil
}

func (m *mockClient) StartConnector(id string, startOpts *model.ConnectorStartOptions, timeout time.Duration) ([]model.Instance, error) {
	m.startCalled = true
	m.startOptions = startOpts
	if m.connectorError != nil {
		return nil, m.connectorError
	}
	return m.instances, nil
}

func (m *mockClient) StopConnector(id string, timeout time.Duration) ([]model.Instance, error) {
	m.stopCalled = true
	if m.connectorError != nil {
		return nil, m.connectorError
	}
	return m.instances, nil
}

// LibraryClient interface
func (m *mockClient) ListRuntimes(timeout time.Duration) ([]model.RuntimeSummary, error) {
	return m.runtimes, nil
}

func (m *mockClient) GetRuntime(id string, timeout time.Duration) (*model.Runtime, error) {
	return m.runtime, nil
}

func (m *mockClient) SearchComponents(filter *model.ComponentSearchFilter, timeout time.Duration) ([]model.ComponentSummary, error) {
	return m.components, nil
}

func (m *mockClient) GetComponent(runtimeId string, kind model.ComponentKind, id string, timeout time.Duration) (*model.Component, error) {
	return m.component, nil
}

// Helper to create a mock AppContext
func newMockAppContext() (*AppContext, *mockClient) {
	mockCl := newMockClient()
	return &AppContext{
		Client:         mockCl,
		DefaultTimeout: 5 * time.Second,
	}, mockCl
}