// Package store is the persistence layer: SQLite-backed users, households, categories,
// transactions, aggregates (aggregate_sum.go including SumIncomeExpenseCentsInRange, SumIncomeExpenseCentsInTwoRanges, SumRunningTotalAndIncomeExpenseInTwoRanges, aggregate_category.go including ListCategoryAmountsInRange + ErrInvalidCategoryAmountKind, aggregate_time.go; household SQL fragments in transaction_sql.go; income/expense AND filters in sql_amount_filter.go), and
// daily movement for the heatmap; per-user notifications (notification_*.go including [Store.MarkNotificationRead], migration v9). Open via [moana/internal/dbutil.OpenStore].
// Household type in household.go; reads in household_query.go (+ household_sql.go); writes in household_mutate.go (+ household_sql.go insert/update; name, invites, detach; [ErrDuplicateUserEmail] from member insert).
// User types in user_types.go; scanning in user_scan.go; reads in user_get.go, user_list.go (preallocated slices), user_sql.go (get-by-email/get-by-id + lists + profile/password/detach updates); profile/password/create in user_mutate.go ([CreateUser] trims email, returns [ErrDuplicateUserEmail] or [ErrInvalidUserEmail]).
// Transaction types in transaction_types.go; reads in transaction_get.go, transaction_list.go (+ transaction_sql.go list/get SQL + conditional insert/update, sql/row helpers, LIKE escape in like_escape.go); writes in transaction_mutate.go (category validation helpers).
// SQLite TEXT timestamps use moana/internal/timeutil (ParseSQLiteTimestamp, FormatSQLiteUTC, NowSQLiteUTC).
// Category type in category_types.go (household-scoped); reads in category_query.go (+ category_sql.go list/get SQL); writes in category_mutate.go (+ category_sql.go insert/update/delete).
// Tests use store_test_helpers_test.go (shared in-memory store) and are split by concern (transactions, users/households, category flow).
// Context deadlines and cancellation propagate to [database/sql]; tests use pre-cancelled contexts in store_test_helpers_test.go to lock in that contract on hot read and write paths.
// Updates/deletes that must affect exactly one row use [execExactlyOneRow] (exec_result.go) instead of duplicating RowsAffected plumbing.
package store
