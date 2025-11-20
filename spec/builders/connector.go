package builders

import "github.com/synadia-io/connect/v2/spec"

type ConnectorBuilder struct {
	s *spec.ConnectorSpec
}

func Connector() *ConnectorBuilder {
	return &ConnectorBuilder{
		s: &spec.ConnectorSpec{},
	}
}

func (b *ConnectorBuilder) Description(v string) *ConnectorBuilder {
	b.s.Description = v
	return b
}

func (b *ConnectorBuilder) RuntimeId(v string) *ConnectorBuilder {
	b.s.RuntimeId = v
	return b
}

func (b *ConnectorBuilder) Steps(sb *StepsBuilder) *ConnectorBuilder {
	b.s.Steps = sb.Build()
	return b
}

func (b *ConnectorBuilder) Build() spec.ConnectorSpec {
	return *b.s
}
