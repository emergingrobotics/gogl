// Package transport carries JSON-RPC calls to a GL.iNet router, owning session
// acquisition, renewal, and retry. Layers above it do not know a session exists.
package transport

import (
	"context"
	"fmt"
)

// Transport performs authenticated JSON-RPC calls against a GL.iNet router.
type Transport interface {
	// Call invokes method on group with args, decoding the result into out.
	// out may be nil to discard the result; args may be nil for no arguments.
	Call(ctx context.Context, group, method string, args, out any) error

	// Close stops the keepalive goroutine and releases resources. Safe to call
	// more than once.
	Close() error
}

// Error is a JSON-RPC error returned by the router, carrying the group and
// method that produced it so a failure is traceable to its call site.
type Error struct {
	Code    int
	Message string
	Group   string
	Method  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("gogl: rpc error %d on %s.%s: %s", e.Code, e.Group, e.Method, e.Message)
}
