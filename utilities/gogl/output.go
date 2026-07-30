package main

import (
	"encoding/json"
	"io"
)

// writeJSON emits v as indented JSON.
//
// Indented rather than compact: this output is read by people at least as often as by
// jq, and a wall of one-line JSON is hostile to the first case for no gain in the second.
func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
