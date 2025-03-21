// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package model

import "encoding/json"
import "fmt"

type SecretDeleteRequest struct {
	// The id of the secret
	Id string `json:"id" yaml:"id" mapstructure:"id"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SecretDeleteRequest) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["id"]; raw != nil && !ok {
		return fmt.Errorf("field id in SecretDeleteRequest: required")
	}
	type Plain SecretDeleteRequest
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = SecretDeleteRequest(plain)
	return nil
}

type SecretDeleteResponse struct {
	// Whether the secret existed
	Existed bool `json:"existed" yaml:"existed" mapstructure:"existed"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SecretDeleteResponse) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["existed"]; raw != nil && !ok {
		return fmt.Errorf("field existed in SecretDeleteResponse: required")
	}
	type Plain SecretDeleteResponse
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = SecretDeleteResponse(plain)
	return nil
}
