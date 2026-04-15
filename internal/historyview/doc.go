// Package historyview builds the history ledger page: template types ([PageData], [Nav], [DayGroup]) in page_types.go;
// filter nav URLs in nav.go ([BuildNav]); day grouping in groups.go ([GroupByDay], groups_test.go); [BuildPage] in page.go;
// query normalization ([ParseHistoryURL], [HistoryURLParams]) in query_parse.go; bounded transaction fetch in limit.go.
// It does not import net/http.
package historyview
