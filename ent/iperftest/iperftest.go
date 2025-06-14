// Code generated by ent, DO NOT EDIT.

package iperftest

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the iperftest type in the database.
	Label = "iperf_test"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTimestamp holds the string denoting the timestamp field in the database.
	FieldTimestamp = "timestamp"
	// FieldSentMbps holds the string denoting the sent_mbps field in the database.
	FieldSentMbps = "sent_mbps"
	// FieldReceivedMbps holds the string denoting the received_mbps field in the database.
	FieldReceivedMbps = "received_mbps"
	// FieldRetransmits holds the string denoting the retransmits field in the database.
	FieldRetransmits = "retransmits"
	// FieldMeanRttMs holds the string denoting the mean_rtt_ms field in the database.
	FieldMeanRttMs = "mean_rtt_ms"
	// FieldDurationSeconds holds the string denoting the duration_seconds field in the database.
	FieldDurationSeconds = "duration_seconds"
	// FieldProtocol holds the string denoting the protocol field in the database.
	FieldProtocol = "protocol"
	// FieldSuccess holds the string denoting the success field in the database.
	FieldSuccess = "success"
	// FieldErrorMessage holds the string denoting the error_message field in the database.
	FieldErrorMessage = "error_message"
	// FieldDaemonID holds the string denoting the daemon_id field in the database.
	FieldDaemonID = "daemon_id"
	// EdgeHost holds the string denoting the host edge name in mutations.
	EdgeHost = "host"
	// Table holds the table name of the iperftest in the database.
	Table = "iperf_tests"
	// HostTable is the table that holds the host relation/edge.
	HostTable = "iperf_tests"
	// HostInverseTable is the table name for the Host entity.
	// It exists in this package in order to avoid circular dependency with the "host" package.
	HostInverseTable = "hosts"
	// HostColumn is the table column denoting the host relation/edge.
	HostColumn = "host_iperf_tests"
)

// Columns holds all SQL columns for iperftest fields.
var Columns = []string{
	FieldID,
	FieldTimestamp,
	FieldSentMbps,
	FieldReceivedMbps,
	FieldRetransmits,
	FieldMeanRttMs,
	FieldDurationSeconds,
	FieldProtocol,
	FieldSuccess,
	FieldErrorMessage,
	FieldDaemonID,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "iperf_tests"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"host_iperf_tests",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultTimestamp holds the default value on creation for the "timestamp" field.
	DefaultTimestamp func() time.Time
	// DefaultDurationSeconds holds the default value on creation for the "duration_seconds" field.
	DefaultDurationSeconds int
	// DefaultProtocol holds the default value on creation for the "protocol" field.
	DefaultProtocol string
	// DefaultSuccess holds the default value on creation for the "success" field.
	DefaultSuccess bool
)

// OrderOption defines the ordering options for the IperfTest queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByTimestamp orders the results by the timestamp field.
func ByTimestamp(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTimestamp, opts...).ToFunc()
}

// BySentMbps orders the results by the sent_mbps field.
func BySentMbps(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSentMbps, opts...).ToFunc()
}

// ByReceivedMbps orders the results by the received_mbps field.
func ByReceivedMbps(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldReceivedMbps, opts...).ToFunc()
}

// ByRetransmits orders the results by the retransmits field.
func ByRetransmits(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRetransmits, opts...).ToFunc()
}

// ByMeanRttMs orders the results by the mean_rtt_ms field.
func ByMeanRttMs(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMeanRttMs, opts...).ToFunc()
}

// ByDurationSeconds orders the results by the duration_seconds field.
func ByDurationSeconds(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDurationSeconds, opts...).ToFunc()
}

// ByProtocol orders the results by the protocol field.
func ByProtocol(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldProtocol, opts...).ToFunc()
}

// BySuccess orders the results by the success field.
func BySuccess(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSuccess, opts...).ToFunc()
}

// ByErrorMessage orders the results by the error_message field.
func ByErrorMessage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldErrorMessage, opts...).ToFunc()
}

// ByDaemonID orders the results by the daemon_id field.
func ByDaemonID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDaemonID, opts...).ToFunc()
}

// ByHostField orders the results by host field.
func ByHostField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newHostStep(), sql.OrderByField(field, opts...))
	}
}
func newHostStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(HostInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, HostTable, HostColumn),
	)
}
