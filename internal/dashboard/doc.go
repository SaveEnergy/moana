// Package dashboard holds dashboard presentation logic: period stats, heatmap cells,
// outflow donut gradient, and [BuildPageData] for the overview template.
// [PageData] and related types live in page_data_types.go; outflow assembly in
// page_data_outflow.go; rolling heatmap window in page_data_heatmap.go; orchestration in page_data.go.
// It does not import net/http; handlers render templates using [PageData].
package dashboard
