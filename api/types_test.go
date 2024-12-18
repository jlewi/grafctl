package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

var expected = Panes{
	"eja": PaneBody{
		Datasource: "somesource",
		Queries: []Query{
			{
				RefID: "A",
				Datasource: Datasource{
					Type: "grafana-clickhouse-datasource",
					UID:  "someuid",
				},
				EditorType: "simplelog",
				RawSQL:     "SELECT Timestamp as \"timestamp\", Body as \"body\", SeverityText as \"level\" FROM \"views\".\"logs\" LIMIT 1000 --- cluster:prod AND service:foyle",
				BuilderOptions: BuilderOptions{
					Database:  "views",
					Table:     "logs",
					QueryType: "logs",
					Mode:      "list",
					Columns: []Column{
						{Name: "Timestamp", Hint: "time"},
						{Name: "SeverityText", Hint: "log_level"},
						{Name: "Body", Hint: "log_message"},
					},
					Meta:           Meta{OtelEnabled: false},
					SimplelogQuery: "cluster:prod AND service:foyle",
					Limit:          1000,
				},
				AdditionalFields: map[string]any{},
				PluginVersion:    "4.5.0",
				Format:           2,
				QueryType:        "logs",
			},
		},
		Range: TimeRange{
			From: "now-5m",
			To:   "now",
		},
		PanelsState: PanelsState{
			Logs: LogsState{
				Columns: map[string]string{
					"0": "timestamp",
					"1": "body",
				},
				VisualisationType: "logs",
			},
		},
	},
}

func Test_Pane(t *testing.T) {
	type testCase struct {
		Name     string
		Input    string
		Expected *Panes
	}

	cases := []testCase{
		{
			Name:     "basic",
			Input:    "pane.json",
			Expected: &expected,
		},
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory")
	}
	tDir := filepath.Join(cwd, "test_data")

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			tFile := filepath.Join(tDir, c.Input)
			raw, err := os.ReadFile(tFile)
			if err != nil {
				t.Fatalf("Failed to read file %v: %v", tFile, err)
			}
			actual := &Panes{}
			if err := json.Unmarshal(raw, actual); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			if d := cmp.Diff(c.Expected, actual); d != "" {
				t.Fatalf("Unexpected diff:\n%+v", d)
			}
		})
	}
}

func Test_QueryUnmarshalYaml(t *testing.T) {
	type testCase struct {
		Name     string
		Input    map[string]any
		Expected *Query
	}

	cases := []testCase{
		{
			Name: "additionalfields",
			Input: map[string]any{
				"refId":      "A",
				"editorType": "sql",
				"customarg":  "customvalue",
			},
			Expected: &Query{
				RefID:      "A",
				EditorType: "sql",
				AdditionalFields: map[string]interface{}{
					"customarg": "customvalue",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			raw, err := yaml.Marshal(c.Input)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}
			actual := &Query{}

			if err := yaml.Unmarshal(raw, actual); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			if d := cmp.Diff(c.Expected, actual); d != "" {
				t.Fatalf("Unexpected diff:\n%+v", d)
			}
		})
	}
}
