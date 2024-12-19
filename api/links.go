package api

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	LinkGVK = schema.FromAPIVersionAndKind(Group+"/"+Version, "GrafanaLink")
)

// GrafanaLink represents a link to a Grafana dashboard
type GrafanaLink struct {
	APIVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind" yaml:"kind"`
	Metadata   Metadata `json:"metadata" yaml:"metadata"`

	// BaseURL is the base URL for links generated from this template
	BaseURL string `json:"baseURL" yaml:"baseURL"`
	// Panes is a map from the ID of the pane to the body of the pane
	Panes Panes `json:"panes" yaml:"panes"`
}
