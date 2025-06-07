package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// SpeedTest holds the schema definition for the SpeedTest entity.
type SpeedTest struct {
	ent.Schema
}

// Fields of the SpeedTest.
func (SpeedTest) Fields() []ent.Field {
	return []ent.Field{
		field.Time("timestamp").
			Default(time.Now),
		field.Float("download_mbps").
			Comment("Download speed in Mbps"),
		field.Float("upload_mbps").
			Comment("Upload speed in Mbps"),
		field.Float("ping_ms").
			Comment("Ping latency in milliseconds"),
		field.Float("jitter_ms").
			Optional().
			Comment("Jitter in milliseconds"),
		field.String("server_name").
			Optional().
			Comment("Speed test server name"),
		field.String("server_id").
			Optional().
			Comment("Speed test server ID"),
		field.String("isp").
			Optional().
			Comment("Internet Service Provider"),
		field.String("external_ip").
			Optional().
			Comment("External IP address"),
		field.String("result_url").
			Optional().
			Comment("URL to full test results"),
		field.String("daemon_id").
			Optional().
			Comment("Identifier of the daemon that performed the test"),
	}
}

// Edges of the SpeedTest.
func (SpeedTest) Edges() []ent.Edge {
	return nil
}
