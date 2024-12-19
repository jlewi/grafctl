package api

const (
	Group   = "grafctl.foyle.io"
	Version = "v1alpha1"
)

// N.B. We need to redefine Metadata and not reuse the version in the K8s libraries
// because we want it to have yaml tags so we can serialize with the YAML library.

// Metadata holds an optional name of the project.
type Metadata struct {
	Name        string            `yaml:"name,omitempty"`
	Namespace   string            `yaml:"namespace,omitempty"`
	Labels      map[string]string `yaml:"labels"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}
