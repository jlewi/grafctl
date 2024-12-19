package grafana

import (
	"encoding/json"
	"fmt"

	"github.com/jlewi/grafctl/api"
	"github.com/pkg/errors"

	// This is what inspired strategic patch merge in K8s
	// https://github.com/kubernetes/apimachinery/blob/8c60292e48e46c4faa1e92acb232ce6adb37512c/pkg/util/strategicpatch/patch.go#L38
	jsonpatch "github.com/evanphx/json-patch/v5"
)

type Patcher struct {
	Clock Clock
}

func NewPatcher(clock Clock) *Patcher {
	return &Patcher{Clock: clock}
}

// ApplyPatch applies the patch defined in the file to the app and returns the panes
//
// N.B. Bases ends up getting modified in place. We might want to make a copy so we don't modify the original.
func (a *Patcher) ApplyPatch(bases []*api.GrafanaLink, patch api.PanePatch) (*api.GrafanaLink, error) {
	if patch.Template == "" {
		return nil, errors.New("Template must be specified in the patch and should be the name of the template to apply")
	}

	if patch.Query == nil {
		return nil, errors.New("Query must be specified in the patch")
	}

	// Find the base
	var base *api.GrafanaLink
	baseNames := make([]string, 0, len(bases))
	for _, b := range bases {
		if b.Metadata.Name == patch.Template {
			base = b
			break
		}
		baseNames = append(baseNames, b.Metadata.Name)
	}

	if base == nil {
		return nil, errors.Errorf("Unable to apply the patch because there is no template %v in the links; add the template to the links in your configuration or select one of your existing links. The known bases are %v", patch.Template, baseNames)
	}

	if len(base.Panes) != 1 {
		// N.B. Right now we assume there is only 1 pane in the base resource and that's the one to apply the patch
		// to. Its TBD whether Grafana has resources where there is more than 1 pane and if there is how identify
		// the correct pane to apply the patch to. I guess we could just specify a parameter in the patch.
		return nil, errors.Errorf("Unable to apply patch to the GrafanaLink. GrafanaLink has %v panes; expected 1", len(base.Panes))
	}

	for k := range base.Panes {
		paneBody := base.Panes[k]
		if err := ApplyPatchToPane(&paneBody, patch); err != nil {
			return nil, errors.Wrapf(err, "Failed to apply patch to template %v", patch.Template)
		}

		if patch.FixTime == nil || *patch.FixTime {
			p := NewRelativeTimeParser()
			p.Clock = a.Clock
			from, err := p.ParseGrafanaRelativeTime(patch.Range.From)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse relative time in from field of value %v", patch.Range.From)
			}
			// From is unix epoch in milliseconds
			paneBody.Range.From = fmt.Sprintf("%d", from.Unix()*1000)

			to, err := p.ParseGrafanaRelativeTime(patch.Range.To)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse relative time in to field of value %v", patch.Range.To)
			}
			// From is unix epoch in milliseconds
			paneBody.Range.To = fmt.Sprintf("%d", to.Unix()*1000)
		}

		base.Panes[k] = paneBody
	}
	// n.b. is hardcoding the paneId here a problem? I think it just needs to be unique within the dictionary.
	return base, nil
}

// ApplyPatchToPane applies the patch to the pane.
func ApplyPatchToPane(pane *api.PaneBody, patch api.PanePatch) error {
	if len(pane.Queries) != 1 {
		return errors.Errorf("Unable to apply patch to the PaneBody. PaneBody has %v queries; expected 1", len(pane.Queries))
	}

	q := pane.Queries[0]

	if err := applyPatch(&q, patch.Query); err != nil {
		return errors.Wrapf(err, "Failed to patch query")
	}

	pane.Queries[0] = q
	return nil
}

// apply the patch. The patch is merged into base.
// base should be a pointer
//
// n.b applyPatch works by serializing the patch into a dictionary. If a key is present in the patch with an
// empty value it will override the field in the base of the same name. Therefore, you usually want your fields
// to be tagged with omit empty.
func applyPatch(base any, patch any) error {
	baseBytes, err := json.Marshal(base)
	if err != nil {
		return errors.Wrapf(err, "Error marshalling base")
	}

	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return errors.Wrapf(err, "Error marshalling patch")
	}
	mergedBytes, err := jsonpatch.MergePatch(baseBytes, patchBytes)
	if err != nil {
		return errors.Wrapf(err, "Error applying the merge patch")
	}

	return json.Unmarshal(mergedBytes, base)
}
