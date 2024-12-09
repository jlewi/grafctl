package api

// PanePatch represents the patch to be applied to one of your pane templates.
// This corresponds to the YAML that is passed on the command line
type PanePatch struct {
	APIVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind" yaml:"kind"`
	Metadata   Metadata `json:"metadata" yaml:"metadata"`

	// Template is the name of the template to apply the patch to
	Template string `json:"template" yaml:"template"`
	// Query is a patch to be applied to the first query in the pane.
	Query map[string]interface{} `json:"query,omitempty" yaml:"query,omitempty"`
	// Range is the time range for the query.
	// Uses the syntax supported by grafana for relative times and units
	// https://grafana.com/docs/grafana/latest/dashboards/use-dashboards/#time-units-and-relative-ranges
	Range TimeRange `json:"range,omitempty" yaml:"range,omitempty"`
	// FixTime is a flag to indicate whether to fix the time range to absolute time or use relative time
	// in the link. Default is true.
	FixTime *bool `json:"fixTime,omitempty" yaml:"fixTime,omitempty"`
}
