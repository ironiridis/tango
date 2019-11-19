package rpc

import "encoding/json"

// Message is a container for JSON signals exchanged with subprocesses.
type Message struct {
	// T is the type of message in the payload; 'set', 'get', 'status', 'event', etc
	T string

	// Optional specifies that if the message is not understood by the recipient or
	// otherwise fails to parse, to not treat as a failure condition. Zero value is
	// to assume the message is critical.
	Optional bool

	// M is the payload of the message, in the format specified by T.
	M json.RawMessage
}
