// Package store is the persistence layer: SQLite-backed users, households, categories,
// transactions, aggregates (aggregate_sum.go, aggregate_category.go, aggregate_time.go), and
// daily movement for the heatmap. Open via [moana/internal/dbutil.OpenStore].
// Household type in household.go; reads in household_query.go; writes in household_mutate.go (name, invites, detach).
// User types in user_types.go; scanning in user_scan.go; reads in user_get.go, user_list.go, user_sql.go; profile/password/create in user_mutate.go.
// Transaction types in transaction_types.go; reads in transaction_get.go, transaction_list.go (+ sql/row helpers); writes in transaction_mutate.go (create/update + category validation helpers).
// SQLite TEXT timestamps use moana/internal/timeutil (ParseSQLiteTimestamp, FormatSQLiteUTC, NowSQLiteUTC).
// Category type in category_types.go (household-scoped); reads in category_query.go (+ category_sql.go); writes in category_mutate.go.
package store
