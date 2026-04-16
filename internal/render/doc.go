// Package render executes HTML templates (app shell and standalone pages) without HTTP handler logic.
// [LayoutData] is in layout_data.go; [Engine] in engine.go; [Engine.Shell] in shell.go (buffers full layout before writing); [Engine.Simple] in simple.go (simple_test.go).
package render
