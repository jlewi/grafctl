package grafana

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"testing"

	"github.com/jlewi/grafctl/api"
)

func Test_GetLogsLink(t *testing.T) {
	type testCase struct {
		Name      string
		PanesFile string
		BaseURL   string
		Expected  string
	}

	cases := []testCase{
		{
			Name:      "basic",
			PanesFile: filepath.Join("..", "..", "api", "test_data/pane.json"),
			BaseURL:   "https://grafana.acme.com",
			Expected:  "https://grafana.acme.com/explore?orgId=1&panes=%7B%22eja%22%3A%7B%22datasource%22%3A%22somesource%22%2C%22queries%22%3A%5B%7B%22builderOptions%22%3A%7B%22columns%22%3A%5B%7B%22hint%22%3A%22time%22%2C%22name%22%3A%22Timestamp%22%7D%2C%7B%22hint%22%3A%22log_level%22%2C%22name%22%3A%22SeverityText%22%7D%2C%7B%22hint%22%3A%22log_message%22%2C%22name%22%3A%22Body%22%7D%5D%2C%22database%22%3A%22views%22%2C%22limit%22%3A1000%2C%22mode%22%3A%22list%22%2C%22queryType%22%3A%22logs%22%2C%22simplelogQuery%22%3A%22cluster%3Aprod+AND+service%3Afoyle%22%2C%22table%22%3A%22logs%22%7D%2C%22datasource%22%3A%7B%22type%22%3A%22grafana-clickhouse-datasource%22%2C%22uid%22%3A%22someuid%22%7D%2C%22editorType%22%3A%22simplelog%22%2C%22format%22%3A2%2C%22pluginVersion%22%3A%224.5.0%22%2C%22queryType%22%3A%22logs%22%2C%22rawSql%22%3A%22SELECT+Timestamp+as+%5C%22timestamp%5C%22%2C+Body+as+%5C%22body%5C%22%2C+SeverityText+as+%5C%22level%5C%22+FROM+%5C%22views%5C%22.%5C%22logs%5C%22+LIMIT+1000+---+cluster%3Aprod+AND+service%3Afoyle%22%2C%22refId%22%3A%22A%22%7D%5D%2C%22range%22%3A%7B%22from%22%3A%22now-5m%22%2C%22to%22%3A%22now%22%7D%2C%22panelsState%22%3A%7B%22logs%22%3A%7B%22columns%22%3A%7B%220%22%3A%22timestamp%22%2C%221%22%3A%22body%22%7D%2C%22visualisationType%22%3A%22logs%22%7D%7D%7D%7D&schemaVersion=1",
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory")
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			tDir := filepath.Join(cwd, c.PanesFile)
			raw, err := os.ReadFile(tDir)
			if err != nil {
				t.Fatalf("Failed to read file %v: %v", tDir, err)
			}
			panes := &api.Panes{}
			if err := json.Unmarshal(raw, panes); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			actual, err := GetLogsLink(c.BaseURL, "1", *panes)
			if err != nil {
				t.Fatalf("Error calling GetLogsLink: %v", err)
			}

			if actual != c.Expected {
				t.Errorf("Got %v;\n Want %v", actual, c.Expected)
			}
		})
	}
}

func TestParseURL(t *testing.T) {
	type testCase struct {
		Name         string
		Input        string
		ExpectedBase string
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory")
	}

	tFile := filepath.Join(cwd, "test_data/testlink.yaml")

	raw, err := os.ReadFile(tFile)
	if err != nil {
		t.Fatalf("Failed to read file %v: %v", tFile, err)
	}

	link := &api.GrafanaLink{}
	if err := yaml.Unmarshal(raw, link); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	testUrl, err := LinkToURL(*link)
	if err != nil {
		t.Fatalf("Error calling LinkToURL: %v", err)
	}

	cases := []testCase{
		{
			Name:         "basic",
			Input:        testUrl,
			ExpectedBase: "https://grafana.acme.com",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			baseUrl, queryArgs, panes, err := ParseURL(c.Input)
			if err != nil {
				t.Fatalf("Error calling ParseURL: %v", err)
			}

			if len(queryArgs) != 2 {
				t.Fatalf("Expected 2 query args; got %d", len(queryArgs))
			}

			if len(panes) != 1 {
				t.Fatalf("Expected 1 pane; got %d", len(panes))
			}

			if baseUrl != c.ExpectedBase {
				t.Errorf("Got %v;\n Want %v", baseUrl, c.ExpectedBase)
			}
		})
	}
}
