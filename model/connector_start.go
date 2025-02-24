// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package model

import "encoding/json"
import "fmt"

type ConnectorStartRequest struct {
	// The id of the connector
	ConnectorId string `json:"connector_id" yaml:"connector_id" mapstructure:"connector_id"`

	// Options corresponds to the JSON schema field "options".
	Options *ConnectorStartOptions `json:"options,omitempty" yaml:"options,omitempty" mapstructure:"options,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ConnectorStartRequest) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["connector_id"]; raw != nil && !ok {
		return fmt.Errorf("field connector_id in ConnectorStartRequest: required")
	}
	type Plain ConnectorStartRequest
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = ConnectorStartRequest(plain)
	return nil
}

type ConnectorStartResponse struct {
	// The instances of the connector that were started
	Instances []Instance `json:"instances" yaml:"instances" mapstructure:"instances"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ConnectorStartResponse) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["instances"]; raw != nil && !ok {
		return fmt.Errorf("field instances in ConnectorStartResponse: required")
	}
	type Plain ConnectorStartResponse
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = ConnectorStartResponse(plain)
	return nil
}
