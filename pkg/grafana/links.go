package grafana

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/go-logr/zapr"
	"github.com/jlewi/grafctl/api"
	"github.com/jlewi/monogo/yamlfiles"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func LinkToURL(link api.GrafanaLink) (string, error) {
	// TODO(jeremy): Rather than hard coding orgid should we add it to GrafanaLink and include it in the URL?
	return GetLogsLink(link.BaseURL, "1", link.Panes)
}

// GetLogsLink returns a link to the Datadog logs matching the given query.
func GetLogsLink(baseUrl string, orgId string, panes api.Panes) (string, error) {
	// Create a new url.Values object
	queryParams := url.Values{}

	queryParams.Add("orgId", orgId)
	queryParams.Add("schemaVersion", "1")

	panesData, err := json.Marshal(panes)

	if err != nil {
		return "", errors.Wrapf(err, "Error marshalling panes data")
	}

	queryParams.Add("panes", string(panesData))

	// Encode the values into a query string
	encodedQuery := queryParams.Encode()
	u := fmt.Sprintf("%s/explore?%s", baseUrl, encodedQuery)
	return u, nil
}

// ParseURL parses the input URL and returns
// baseUrl - The base URL
// a map of query parameters other than the panes object
// The panes object.
func ParseURL(inputURL string) (string, map[string][]string, []*api.Panes, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", nil, nil, errors.Wrapf(err, "failed to parse URL: %v", inputURL)
	}

	values := parsedURL.Query()

	panesJson := []string{}

	queryArgs := map[string][]string{}
	for key, value := range values {
		if key == "panes" {
			panesJson = append(panesJson, value...)
		} else {
			queryArgs[key] = value
		}
	}

	panes := make([]*api.Panes, 0)
	for _, paneJson := range panesJson {
		pane := &api.Panes{}
		err := json.Unmarshal([]byte(paneJson), pane)
		if err != nil {
			return "", nil, nil, errors.Wrapf(err, "Error unmarshalling panes")
		}
		panes = append(panes, pane)
	}

	// Get only the scheme and host
	baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	// Clean the path to keep only the base path without the last segment
	cleanedPath := path.Clean("/")
	baseURL = baseURL + cleanedPath

	baseURL = strings.TrimSuffix(baseURL, "/")
	return baseURL, queryArgs, panes, nil
}

// URLToLink converts a URL to a GrafanaLink
func URLToLink(logUrl string) (*api.GrafanaLink, error) {
	baseUrl, queryParams, panes, err := ParseURL(logUrl)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stdout, "Query Arguments:\n%v\n", queryParams)

	if len(panes) == 0 {
		return nil, errors.New("No panes found in URL")
	}
	if len(panes) > 1 {
		// This means the panes argument is repeated. What should we do in that case
		return nil, errors.New("Multiple panes found in URL")
	}

	link := &api.GrafanaLink{
		APIVersion: api.LinkGVK.GroupVersion().String(),
		Kind:       api.LinkGVK.Kind,
		BaseURL:    baseUrl,
		Panes:      *panes[0],
	}
	return link, nil
}

// LoadGrafanaLinksInDir looks for YAML files in the given directory containing GrafanaLink resources
func LoadGrafanaLinksInDir(dir string) ([]*api.GrafanaLink, error) {
	log := zapr.NewLogger(zap.L())
	files, err := yamlfiles.Find(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "Error finding files in %v", dir)
	}

	links := make([]*api.GrafanaLink, 0)
	for _, f := range files {
		nodes, err := yamlfiles.Read(f)
		if err != nil {
			log.Error(err, "Error reading file", "file", f)
			continue
		}

		for _, node := range nodes {
			if node.GetKind() != api.LinkGVK.Kind {
				continue
			}

			link := &api.GrafanaLink{}
			if err := node.YNode().Decode(link); err != nil {
				log.Error(err, "Failed to decode GrafanaLink", "file", "name", node.GetName())
				continue
			}
			links = append(links, link)
		}
	}

	return links, nil
}
