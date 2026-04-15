// Package dashboard holds dashboard presentation logic: period stats, heatmap cells,
// outflow donut gradient, and [BuildPageData] for the overview template.
// [PageData] and related types live in page_data_types.go; period query normalization in page_data_period.go;
// outflow assembly in page_data_outflow.go; rolling heatmap window in page_data_heatmap.go; orchestration in page_data.go.
// It does not import net/http; handlers render templates using [PageData].
// build_page_data_test.go exercises [BuildPageData] against a real in-memory store.
package dashboard
