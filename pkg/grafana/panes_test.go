package grafana

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jlewi/grafctl/api"
)

func Test_ApplyPatchToLink(t *testing.T) {
	type testCase struct {
		name     string
		bases    []*api.GrafanaLink
		patch    api.PanePatch
		expected *api.GrafanaLink
	}

	cases := []testCase{
		{
			name: "basic",
			bases: []*api.GrafanaLink{
				{
					Metadata: api.Metadata{
						Name: "test",
					},
					Panes: api.Panes{
						"eja": api.PaneBody{
							Queries: []api.Query{
								{
									BuilderOptions: api.BuilderOptions{
										Database: "somedatabase",
										Table:    "sometable",
									},
								},
							},
						},
					},
				},
			},
			patch: api.PanePatch{
				Template: "test",
				Query: map[string]any{
					"builderOptions": map[string]any{
						"simplelogQuery": "service:foo",
					},
					// This is a custom query argument
					"customarg": "customvalue",
				},
				Range: api.TimeRange{
					From: "now-1h",
					To:   "now",
				},
			},
			expected: &api.GrafanaLink{
				Metadata: api.Metadata{
					Name: "test",
				},
				Panes: api.Panes{
					"eja": api.PaneBody{
						Queries: []api.Query{
							{
								BuilderOptions: api.BuilderOptions{
									Database:       "somedatabase",
									Table:          "sometable",
									SimplelogQuery: "service:foo",
								},
								AdditionalFields: map[string]interface{}{
									"customarg": "customvalue",
								},
							},
						},
						Range: api.TimeRange{
							From: "1708863900000",
							To:   "1708867500000",
						},
					},
				},
			},
		},
	}

	applier := NewPatcher(FakeClock{})
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := applier.ApplyPatch(c.bases, c.patch)
			if err != nil {
				t.Fatalf("Error applying patch: %v", err)
			}
			if d := cmp.Diff(c.expected, actual); d != "" {
				t.Fatalf("Unexpected diff:\n%v", d)
			}
		})
	}
}

func Test_applyPatch(t *testing.T) {
	type testCase struct {
		name     string
		base     any
		patch    any
		expected any
	}

	cases := []testCase{
		{
			name: "basic",
			base: &api.BuilderOptions{
				Database: "somedatabase",
				Table:    "sometable",
			},
			patch: api.BuilderOptions{
				SimplelogQuery: "select *",
			},
			expected: &api.BuilderOptions{
				Database:       "somedatabase",
				Table:          "sometable",
				SimplelogQuery: "select *",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := applyPatch(c.base, c.patch); err != nil {
				t.Fatalf("Error applying patch: %v", err)
			}
			if d := cmp.Diff(c.expected, c.base); d != "" {
				t.Fatalf("Unexpected diff:\n%v", d)
			}
		})
	}
}
