// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/bfirestone/speed-checker/ent/host"
	"github.com/bfirestone/speed-checker/ent/iperftest"
	"github.com/bfirestone/speed-checker/ent/schema"
	"github.com/bfirestone/speed-checker/ent/speedtest"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	hostFields := schema.Host{}.Fields()
	_ = hostFields
	// hostDescPort is the schema descriptor for port field.
	hostDescPort := hostFields[2].Descriptor()
	// host.DefaultPort holds the default value on creation for the port field.
	host.DefaultPort = hostDescPort.Default.(int)
	// hostDescActive is the schema descriptor for active field.
	hostDescActive := hostFields[4].Descriptor()
	// host.DefaultActive holds the default value on creation for the active field.
	host.DefaultActive = hostDescActive.Default.(bool)
	iperftestFields := schema.IperfTest{}.Fields()
	_ = iperftestFields
	// iperftestDescTimestamp is the schema descriptor for timestamp field.
	iperftestDescTimestamp := iperftestFields[0].Descriptor()
	// iperftest.DefaultTimestamp holds the default value on creation for the timestamp field.
	iperftest.DefaultTimestamp = iperftestDescTimestamp.Default.(func() time.Time)
	// iperftestDescDurationSeconds is the schema descriptor for duration_seconds field.
	iperftestDescDurationSeconds := iperftestFields[5].Descriptor()
	// iperftest.DefaultDurationSeconds holds the default value on creation for the duration_seconds field.
	iperftest.DefaultDurationSeconds = iperftestDescDurationSeconds.Default.(int)
	// iperftestDescProtocol is the schema descriptor for protocol field.
	iperftestDescProtocol := iperftestFields[6].Descriptor()
	// iperftest.DefaultProtocol holds the default value on creation for the protocol field.
	iperftest.DefaultProtocol = iperftestDescProtocol.Default.(string)
	// iperftestDescSuccess is the schema descriptor for success field.
	iperftestDescSuccess := iperftestFields[7].Descriptor()
	// iperftest.DefaultSuccess holds the default value on creation for the success field.
	iperftest.DefaultSuccess = iperftestDescSuccess.Default.(bool)
	speedtestFields := schema.SpeedTest{}.Fields()
	_ = speedtestFields
	// speedtestDescTimestamp is the schema descriptor for timestamp field.
	speedtestDescTimestamp := speedtestFields[0].Descriptor()
	// speedtest.DefaultTimestamp holds the default value on creation for the timestamp field.
	speedtest.DefaultTimestamp = speedtestDescTimestamp.Default.(func() time.Time)
}
