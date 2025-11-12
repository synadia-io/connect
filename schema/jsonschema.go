package schema

import (
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/synadia-io/connect/model"
)

func ToJsonSchema(runtime string, version string, kind string, cmp model.Component) (*jsonschema.Schema, error) {
	s := &jsonschema.Schema{
		ID:    fmt.Sprintf("%s.v%s.%s.%s", runtime, version, kind, cmp.Name),
		Title: cmp.Label,
		Type:  "object",
		Extra: map[string]any{
			"status": string(cmp.Status),
		},
		Properties: map[string]*jsonschema.Schema{},
		Required:   []string{},
	}

	if cmp.Description != nil {
		s.Description = *cmp.Description
	}

	if cmp.Icon != nil {
		s.Extra["icon"] = *cmp.Icon
	}

	for _, field := range cmp.Fields {
		fs, err := fieldSchema(field.Kind, field.Type, field)
		if err != nil {
			return nil, err
		}

		s.Properties[field.Name] = fs

		if field.Optional != nil && !*field.Optional {
			s.Required = append(s.Required, field.Name)
		}
	}

	return s, nil
}

func fieldSchema(fk model.ComponentFieldKind, ft model.ComponentFieldType, fld model.ComponentField) (*jsonschema.Schema, error) {
	switch fk {
	case model.ComponentFieldKindScalar:
		switch ft {
		case model.ComponentFieldTypeString, model.ComponentFieldTypeExpression, model.ComponentFieldTypeCondition:
			return stringSchema(fld)
		case model.ComponentFieldTypeInt:
			return integerSchema(fld)
		case model.ComponentFieldTypeBool:
			return booleanSchema(fld)
		case model.ComponentFieldTypeObject:
			return objectSchema(fld)
		case model.ComponentFieldTypeScanner:
			return scannerSchema(fld)
		default:
			return nil, fmt.Errorf("unsupported field type: %s", fld.Type)
		}
	case model.ComponentFieldKindList:
		itemSchema, err := fieldSchema(model.ComponentFieldKindScalar, ft, model.ComponentField{
			Type:   ft,
			Fields: fld.Fields,
		})
		if err != nil {
			return nil, err
		}

		result, err := commonSchema(fld)
		if err != nil {
			return nil, err
		}

		result.Type = "array"
		result.Items = itemSchema
		return result, nil

	case model.ComponentFieldKindMap:
		itemSchema, err := fieldSchema(model.ComponentFieldKindScalar, ft, model.ComponentField{
			Type:   ft,
			Fields: fld.Fields,
		})
		if err != nil {
			return nil, err
		}

		result, err := commonSchema(fld)
		if err != nil {
			return nil, err
		}

		result.Type = "object"
		result.AdditionalProperties = itemSchema
		return result, nil

	default:
		return nil, fmt.Errorf("unsupported field kind: %s", fk)
	}
}

func commonSchema(fld model.ComponentField) (*jsonschema.Schema, error) {
	result := &jsonschema.Schema{
		Title: fld.Label,
		Extra: map[string]any{},
	}

	if fld.Description != nil {
		result.Description = *fld.Description
	}

	if fld.Default != nil {
		b, err := json.Marshal(fld.Default)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal default value: %w", err)
		}

		result.Default = b
	}

	if fld.Secret != nil {
		result.Extra["secret"] = *fld.Secret
	}

	if fld.RenderHint != nil {
		result.Extra["preset"] = *fld.RenderHint
	}

	if fld.Examples != nil {
		result.Examples = fld.Examples
	}

	return result, nil
}

func stringSchema(fld model.ComponentField) (*jsonschema.Schema, error) {
	s, err := commonSchema(fld)
	if err != nil {
		return nil, err
	}

	s.Type = "string"

	switch fld.Type {
	case model.ComponentFieldTypeExpression:
		s.Extra["render_hint"] = "expression"
	case model.ComponentFieldTypeCondition:
		s.Extra["render_hint"] = "condition"
	}

	for _, constraint := range fld.Constraints {
		// -- add enum values
		if constraint.Enum != nil {
			for _, v := range constraint.Enum {
				s.Enum = append(s.Enum, v)
			}
		}

		if constraint.Preset != nil {
			s.Extra["preset"] = constraint.Preset
		}
		// todo: add other constraints
	}

	return s, nil
}

func booleanSchema(fld model.ComponentField) (*jsonschema.Schema, error) {
	s, err := commonSchema(fld)
	if err != nil {
		return nil, err
	}

	s.Type = "boolean"

	return s, nil
}

func integerSchema(fld model.ComponentField) (*jsonschema.Schema, error) {
	s, err := commonSchema(fld)
	if err != nil {
		return nil, err
	}

	s.Type = "integer"

	// todo: add range constraints

	return s, nil
}

func objectSchema(fld model.ComponentField) (*jsonschema.Schema, error) {
	s, err := commonSchema(fld)
	if err != nil {
		return nil, err
	}

	s.Type = "object"
	s.Properties = map[string]*jsonschema.Schema{}

	// add the child fields
	for _, childFld := range fld.Fields {
		childSchema, err := fieldSchema(childFld.Kind, childFld.Type, *childFld)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %s to jsonschema: %w", childFld.Name, err)
		}

		s.Properties[childFld.Name] = childSchema

		if childFld.Optional != nil && !*childFld.Optional {
			s.Required = append(s.Required, childFld.Name)
		}
	}

	return s, nil
}

func scannerSchema(fld model.ComponentField) (*jsonschema.Schema, error) {
	s, err := commonSchema(fld)
	if err != nil {
		return nil, err
	}

	s.Type = "object"
	s.AdditionalProperties = &jsonschema.Schema{
		Type: "object",
	}

	return s, nil
}
