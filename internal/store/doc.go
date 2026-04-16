// Package store is the persistence layer: SQLite-backed users, households, categories,
// transactions, aggregates (aggregate_sum.go including SumIncomeExpenseCentsInRange, SumIncomeExpenseCentsInTwoRanges, SumRunningTotalAndIncomeExpenseInTwoRanges, aggregate_category.go including ListCategoryAmountsInRange + ErrInvalidCategoryAmountKind, aggregate_time.go; household SQL fragments in transaction_sql.go), and
// daily movement for the heatmap. Open via [moana/internal/dbutil.OpenStore].
// Household type in household.go; reads in household_query.go; writes in household_mutate.go (name, invites, detach; [ErrDuplicateUserEmail] from member insert).
// User types in user_types.go; scanning in user_scan.go; reads in user_get.go, user_list.go (preallocated slices), user_sql.go; profile/password/create in user_mutate.go ([CreateUser] trims email, returns [ErrDuplicateUserEmail] or [ErrInvalidUserEmail]).
// Transaction types in transaction_types.go; reads in transaction_get.go, transaction_list.go (+ sql/row helpers, LIKE escape in like_escape.go); writes in transaction_mutate.go (create/update + category validation helpers).
// SQLite TEXT timestamps use moana/internal/timeutil (ParseSQLiteTimestamp, FormatSQLiteUTC, NowSQLiteUTC).
// Category type in category_types.go (household-scoped); reads in category_query.go (+ category_sql.go); writes in category_mutate.go.
// Tests use store_test_helpers_test.go (shared in-memory store) and are split by concern (transactions, users/households, category flow).
package store
