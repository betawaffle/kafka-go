// Package kafka is a pure Go client library for dealing with Apache Kafka (versions 2.0 and later).
package kafka

// Logger is the instance of a StdLogger interface that this library writes connection
// management events to. By default it is set to discard all log messages, but you can
// set it to redirect wherever you want.
var Logger StdLogger = nopLogger{}

type nopLogger struct {
}

func (nopLogger) Print(v ...interface{}) {
}

func (nopLogger) Printf(f string, v ...interface{}) {
}

func (nopLogger) Println(v ...interface{}) {
}

// StdLogger is used to log error messages.
type StdLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// PanicHandler is called for recovering from panics spawned internally to the library (and thus
// not recoverable by the caller's goroutine). Defaults to nil, which means panics are not recovered.
var PanicHandler func(interface{})

// MaxRequestSize is the maximum size (in bytes) of any request that this library will attempt to send. Trying
// to send a request larger than this will result in an PacketEncodingError. The default of 100 MiB is aligned
// with Kafka's default `socket.request.max.bytes`, which is the largest request the broker will attempt
// to process.
var MaxRequestSize int32 = 100 * 1024 * 1024

// MaxResponseSize is the maximum size (in bytes) of any response that this library will attempt to parse. If
// a broker returns a response message larger than this value, this library will return a PacketDecodingError to
// protect the client from running out of memory. Please note that brokers do not have any natural limit on
// the size of responses they send. In particular, they can send arbitrarily large fetch responses to consumers
// (see https://issues.apache.org/jira/browse/KAFKA-2063).
var MaxResponseSize int32 = 100 * 1024 * 1024
