package model

import "bytes"

// https://stackoverflow.com/questions/55014001/check-if-json-is-object-or-array

func jsonIsArray(data []byte) bool {
	// Get slice of data with optional leading whitespace removed.
	// See RFC 7159, Section 2 for the definition of JSON whitespace.
	x := bytes.TrimLeft(data, " \t\r\n")
	return len(x) > 0 && x[0] == '['
}

func jsonIsObject(data []byte) bool {
	// Get slice of data with optional leading whitespace removed.
	// See RFC 7159, Section 2 for the definition of JSON whitespace.
	x := bytes.TrimLeft(data, " \t\r\n")
	return len(x) > 0 && x[0] == '{'
}
