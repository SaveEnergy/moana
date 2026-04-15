// Package app wires core dependencies into [moana/internal/handlers.App] (composition root for HTTP).
// app_test.go exercises [New] and [HTTPHandler] without importing testutil (avoids a package cycle).
package app
