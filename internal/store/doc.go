// Package store is the persistence layer: SQLite-backed users, households, categories,
// transactions, aggregates, and daily movement for the heatmap. Open via [moana/internal/dbutil.OpenStore].
// User types and scanning live in user.go; reads in user_query.go; writes in user_mutate.go.
// Transaction types in transaction.go; reads in transaction_query.go; writes in transaction_mutate.go.
// Category in category.go; reads in category_query.go; writes in category_mutate.go.
package store
