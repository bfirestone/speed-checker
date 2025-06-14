// Code generated by ent, DO NOT EDIT.

package host

import (
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the host type in the database.
	Label = "host"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldHostname holds the string denoting the hostname field in the database.
	FieldHostname = "hostname"
	// FieldPort holds the string denoting the port field in the database.
	FieldPort = "port"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldActive holds the string denoting the active field in the database.
	FieldActive = "active"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// EdgeIperfTests holds the string denoting the iperf_tests edge name in mutations.
	EdgeIperfTests = "iperf_tests"
	// Table holds the table name of the host in the database.
	Table = "hosts"
	// IperfTestsTable is the table that holds the iperf_tests relation/edge.
	IperfTestsTable = "iperf_tests"
	// IperfTestsInverseTable is the table name for the IperfTest entity.
	// It exists in this package in order to avoid circular dependency with the "iperftest" package.
	IperfTestsInverseTable = "iperf_tests"
	// IperfTestsColumn is the table column denoting the iperf_tests relation/edge.
	IperfTestsColumn = "host_iperf_tests"
)

// Columns holds all SQL columns for host fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldHostname,
	FieldPort,
	FieldType,
	FieldActive,
	FieldDescription,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultPort holds the default value on creation for the "port" field.
	DefaultPort int
	// DefaultActive holds the default value on creation for the "active" field.
	DefaultActive bool
)

// Type defines the type for the "type" enum field.
type Type string

// Type values.
const (
	TypeLan    Type = "lan"
	TypeVpn    Type = "vpn"
	TypeRemote Type = "remote"
)

func (_type Type) String() string {
	return string(_type)
}

// TypeValidator is a validator for the "type" field enum values. It is called by the builders before save.
func TypeValidator(_type Type) error {
	switch _type {
	case TypeLan, TypeVpn, TypeRemote:
		return nil
	default:
		return fmt.Errorf("host: invalid enum value for type field: %q", _type)
	}
}

// OrderOption defines the ordering options for the Host queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByHostname orders the results by the hostname field.
func ByHostname(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHostname, opts...).ToFunc()
}

// ByPort orders the results by the port field.
func ByPort(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPort, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByActive orders the results by the active field.
func ByActive(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldActive, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByIperfTestsCount orders the results by iperf_tests count.
func ByIperfTestsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newIperfTestsStep(), opts...)
	}
}

// ByIperfTests orders the results by iperf_tests terms.
func ByIperfTests(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newIperfTestsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newIperfTestsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(IperfTestsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, IperfTestsTable, IperfTestsColumn),
	)
}
