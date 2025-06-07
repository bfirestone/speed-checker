package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// IperfTest holds the schema definition for the IperfTest entity.
type IperfTest struct {
	ent.Schema
}

// Fields of the IperfTest.
func (IperfTest) Fields() []ent.Field {
	return []ent.Field{
		field.Time("timestamp").
			Default(time.Now),
		field.Float("sent_mbps").
			Comment("Upload speed in Mbps (sent to server)"),
		field.Float("received_mbps").
			Comment("Download speed in Mbps (received from server)"),
		field.Float("retransmits").
			Optional().
			Comment("Number of retransmits"),
		field.Float("mean_rtt_ms").
			Optional().
			Comment("Mean round-trip time in milliseconds"),
		field.Int("duration_seconds").
			Default(10).
			Comment("Test duration in seconds"),
		field.String("protocol").
			Default("TCP").
			Comment("Protocol used (TCP/UDP)"),
		field.Bool("success").
			Default(true).
			Comment("Whether the test completed successfully"),
		field.String("error_message").
			Optional().
			Comment("Error message if test failed"),
		field.String("daemon_id").
			Optional().
			Comment("Identifier of the daemon that performed the test"),
	}
}

// Edges of the IperfTest.
func (IperfTest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("host", Host.Type).
			Ref("iperf_tests").
			Unique(),
	}
}
