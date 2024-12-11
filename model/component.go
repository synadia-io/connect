package model

type ComponentStatus string

const (
	ComponentStatusActive       ComponentStatus = "stable"
	ComponentStatusPreview      ComponentStatus = "preview"
	ComponentStatusExperimental ComponentStatus = "experimental"
	ComponentStatusDeprecated   ComponentStatus = "deprecated"
)

type ComponentKind string

const (
	ComponentKindScanner ComponentKind = "scanner"
	ComponentKindSource  ComponentKind = "source"
	ComponentKindSink    ComponentKind = "sink"
)

type Component struct {
	RuntimeId   string          `yaml:"runtime_id" json:"runtime_id"`
	VersionId   string          `yaml:"version_id" json:"version_id"`
	Name        string          `yaml:"name" json:"name"`
	Label       string          `yaml:"label" json:"label"`
	Kind        ComponentKind   `yaml:"kind" json:"kind"`
	Status      ComponentStatus `yaml:"status" json:"status"`
	Description string          `yaml:"description" json:"description"`
	Fields      []Field         `yaml:"fields,omitempty" json:"fields,omitempty"`
}

type FieldKind string

const (
	FieldKindScalar FieldKind = "scalar"
	FieldKindMap    FieldKind = "map"
	FieldKindList   FieldKind = "list"
)

var fieldKinds = []FieldKind{FieldKindScalar, FieldKindMap, FieldKindList}

type FieldType string

const (
	FieldTypeBool       FieldType = "bool"
	FieldTypeInt        FieldType = "int"
	FieldTypeObject     FieldType = "object"
	FieldTypeScanner    FieldType = "scanner"
	FieldTypeString     FieldType = "string"
	FieldTypeExpression FieldType = "expression"
	FieldTypeCondition  FieldType = "condition"
)

var fieldTypes = []FieldType{
	FieldTypeBool,
	FieldTypeInt,
	FieldTypeObject,
	FieldTypeScanner,
	FieldTypeString,
	FieldTypeExpression,
	FieldTypeCondition,
}

type Field struct {
	Name        string       `yaml:"name" json:"name"`
	Label       string       `yaml:"label" json:"label"`
	Type        FieldType    `yaml:"type" json:"type"`
	Kind        FieldKind    `yaml:"kind" json:"kind"`
	Description string       `yaml:"description,omitempty" json:"description,omitempty"`
	Secret      bool         `yaml:"secret,omitempty" json:"secret,omitempty"`
	Default     any          `yaml:"default,omitempty" json:"default,omitempty"`
	Optional    bool         `yaml:"optional,omitempty" json:"optional,omitempty"`
	Examples    []any        `yaml:"examples,omitempty" json:"examples,omitempty"`
	Fields      []Field      `yaml:"fields,omitempty" json:"fields,omitempty"`
	Constraints []Constraint `yaml:"constraints,omitempty" json:"constraints,omitempty"`

	// RenderHint is a hint to the UI on how to render the field
	RenderHint string `yaml:"render_hint,omitempty" json:"render_hint,omitempty"`
}

type Constraint struct {
	Regex  string   `yaml:"regex,omitempty" json:"regex,omitempty"`
	Range  *Range   `yaml:"range,omitempty" json:"range,omitempty"`
	Enum   []string `yaml:"enum,omitempty" json:"enum,omitempty"`
	Preset string   `yaml:"preset,omitempty" json:"preset,omitempty"`
}

type Range struct {
	LesserThan       float64 `yaml:"lt,omitempty" json:"lt,omitempty"`
	LesserThanEqual  float64 `yaml:"lte,omitempty" json:"lte,omitempty"`
	GreaterThan      float64 `yaml:"gt,omitempty" json:"gt,omitempty"`
	GreaterThanEqual float64 `yaml:"gte,omitempty" json:"gte,omitempty"`
}
