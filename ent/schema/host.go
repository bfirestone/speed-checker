package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Host holds the schema definition for the Host entity.
type Host struct {
	ent.Schema
}

// Fields of the Host.
func (Host) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("Friendly name for the host"),
		field.String("hostname").
			Comment("Hostname or IP address"),
		field.Int("port").
			Default(5201).
			Comment("iperf3 server port"),
		field.Enum("type").
			Values("lan", "vpn", "remote").
			Comment("Host type: lan, vpn, or remote"),
		field.Bool("active").
			Default(true).
			Comment("Whether this host should be included in tests"),
		field.String("description").
			Optional().
			Comment("Optional description of the host"),
	}
}

// Edges of the Host.
func (Host) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("iperf_tests", IperfTest.Type),
	}
}
