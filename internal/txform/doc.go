// Package txform parses HTTP form fields for transaction create/update into validated cents,
// UTC timestamps, optional category id, and trimmed description — no HTML or routing.
// [Parsed] lives in parsed.go; [Parse] in parse.go.
package txform
