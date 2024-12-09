package api

import (
	"encoding/json"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	queryKnownFields = []string{"refId", "datasource", "editorType", "rawSql", "builderOptions", "pluginVersion", "format", "queryType"}
)

// N.B. Merging the datastructures requires omitempty tags to be added to the fields

// Panes is a map from the ID of the pane
// to the body of the pane.
type Panes map[string]PaneBody

// PaneBody represents a pane in the log explorer
// viewer:
// https://grafana.com/docs/grafana/latest/explore/get-started-with-explore/#generate-explore-urls-from-external-tools
type PaneBody struct {
	Datasource  string      `json:"datasource,omitempty" yaml:"datasource,omitempty"`
	Queries     []Query     `json:"queries,omitempty" yaml:"queries,omitempty"`
	Range       TimeRange   `json:"range,omitempty" yaml:"range,omitempty"`
	PanelsState PanelsState `json:"panelsState,omitempty" yaml:"panelsState,omitempty"`
}

// Query represents a query in the log explorer
// https://grafana.com/docs/grafana/latest/explore/get-started-with-explore/#generate-explore-urls-from-external-tools
// Most of these fields are not required by Grafana but our fields used by the logs datasource.
type Query struct {
	RefID      string     `json:"refId,omitempty" yaml:"refId,omitempty"`
	Datasource Datasource `json:"datasource,omitempty" yaml:"datasource,omitempty"`
	EditorType string     `json:"editorType,omitempty" yaml:"editorType,omitempty"`
	RawSQL     string     `json:"rawSql,omitempty" yaml:"rawSql,omitempty"`
	// TODO(jeremy): BuilderOptions is not a standard field in Grafana. Should we treat it as an AdditionalField?
	BuilderOptions   BuilderOptions         `json:"builderOptions,omitempty" yaml:"builderOptions,omitempty"`
	PluginVersion    string                 `json:"pluginVersion,omitempty" yaml:"pluginVersion,omitempty"`
	Format           int                    `json:"format,omitempty" yaml:"format,omitempty"`
	QueryType        string                 `json:"queryType,omitempty" yaml:"queryType,omitempty"`
	AdditionalFields map[string]interface{} `json:"-" yaml:"-"`
}

type QueryKnownFields struct {
	RefID          string         `json:"refId,omitempty" yaml:"refId,omitempty"`
	Datasource     Datasource     `json:"datasource,omitempty" yaml:"datasource,omitempty"`
	EditorType     string         `json:"editorType,omitempty" yaml:"editorType,omitempty"`
	RawSQL         string         `json:"rawSql,omitempty" yaml:"rawSql,omitempty"`
	BuilderOptions BuilderOptions `json:"builderOptions,omitempty" yaml:"builderOptions,omitempty"`
	PluginVersion  string         `json:"pluginVersion,omitempty" yaml:"pluginVersion,omitempty"`
	Format         int            `json:"format,omitempty" yaml:"format,omitempty"`
	QueryType      string         `json:"queryType,omitempty" yaml:"queryType,omitempty"`
}

// UnmarshalJSON method custom unmarshal function to deal with additional fields
func (q *Query) UnmarshalJSON(data []byte) error {

	n := &yaml.Node{}
	if err := yaml.Unmarshal(data, n); err != nil {
		return err
	}

	return q.UnmarshalYAML(n)
}

// UnmarshalYAML method custom unmarshal function to deal with additional fields
func (q *Query) UnmarshalYAML(value *yaml.Node) error {
	// See: https://choly.ca/post/go-json-marshalling/

	// Define an Alias type for the query type to avoid recursion.
	// This gives us all the fields but none of the methods
	type Alias Query

	// Create an alias of the query type that points to q so that
	// we end up decoding into q.
	aux := (*Alias)(q)
	// Unmarshal into the alias
	if err := value.Decode(aux); err != nil {
		return err
	}

	// Decode additional fields
	q.AdditionalFields = make(map[string]any)
	if err := value.Decode(&q.AdditionalFields); err != nil {
		return err
	}

	// Remove known fields from AdditionalFields
	for _, f := range queryKnownFields {
		delete(q.AdditionalFields, f)
	}

	return nil
}

// MarshalYAML custom marshall function
func (c *Query) MarshalYAML() ([]byte, error) {
	// Create a map to hold all fields
	// Need to keep this in sync with the fields
	data := map[string]interface{}{
		"refId":          c.RefID,
		"datasource":     c.Datasource,
		"editorType":     c.EditorType,
		"rawSql":         c.RawSQL,
		"builderOptions": c.BuilderOptions,
		"pluginVersion":  c.PluginVersion,
		"format":         c.Format,
		"queryType":      c.QueryType,
	}

	// Add all additional fields to the map
	for key, value := range c.AdditionalFields {
		data[key] = value
	}

	// Marshal the combined map to JSON
	return yaml.Marshal(data)
}

func (c *Query) MarshalJSON() ([]byte, error) {
	yData, err := yaml.Marshal(c)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to marshal query to yaml")
	}

	m := make(map[string]interface{})
	if err := yaml.Unmarshal(yData, &m); err != nil {
		return nil, errors.Wrapf(err, "Failed to unmarshal yaml to map")
	}

	return json.Marshal(m)
}

type Datasource struct {
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	UID  string `json:"uid,omitempty" yaml:"uid,omitempty"`
}

// BuilderOptions is the options for the builderOptions panel.
// TODO(jeremy): Should we add support for additionalFields as in Query with custom unmarshal and marshal functions
// So that we properly handle additional fields that could be added later?
type BuilderOptions struct {
	Database       string   `json:"database,omitempty" yaml:"database,omitempty"`
	Table          string   `json:"table,omitempty" yaml:"table,omitempty"`
	QueryType      string   `json:"queryType,omitempty" yaml:"queryType,omitempty"`
	Mode           string   `json:"mode,omitempty" yaml:"mode,omitempty"`
	Columns        []Column `json:"columns,omitempty" yaml:"columns,omitempty"`
	Meta           Meta     `json:"meta,omitempty" yaml:"meta,omitempty"`
	Limit          int      `json:"limit,omitempty" yaml:"limit,omitempty"`
	SimplelogQuery string   `json:"simplelogQuery,omitempty" yaml:"simplelogQuery,omitempty"`
}

type Column struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Hint string `json:"hint,omitempty" yaml:"hint,omitempty"`
}

type Meta struct {
	OtelEnabled bool `json:"otelEnabled,omitempty" yaml:"otelEnabled,omitempty"`
}

type TimeRange struct {
	From string `json:"from,omitempty" yaml:"from,omitempty"`
	To   string `json:"to,omitempty" yaml:"to,omitempty"`
}

type PanelsState struct {
	Logs LogsState `json:"logs,omitempty" yaml:"logs,omitempty"`
}

type LogsState struct {
	Columns           map[string]string `json:"columns,omitempty" yaml:"columns,omitempty"`
	VisualisationType string            `json:"visualisationType,omitempty" yaml:"visualisationType,omitempty"`
}
