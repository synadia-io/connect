// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package model

import "encoding/json"
import "fmt"

type SecretListRequest map[string]interface{}

type SecretListResponse struct {
	// Secrets corresponds to the JSON schema field "secrets".
	Secrets []SecretSummary `json:"secrets" yaml:"secrets" mapstructure:"secrets"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SecretListResponse) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["secrets"]; raw != nil && !ok {
		return fmt.Errorf("field secrets in SecretListResponse: required")
	}
	type Plain SecretListResponse
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = SecretListResponse(plain)
	return nil
}
