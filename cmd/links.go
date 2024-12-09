package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jlewi/monogo/helpers"

	"github.com/jlewi/grafctl/api"
	"github.com/jlewi/grafctl/pkg/application"
	"github.com/jlewi/grafctl/pkg/config"
	"github.com/jlewi/grafctl/pkg/grafana"
	"github.com/jlewi/grafctl/pkg/version"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewExploreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "links",
	}
	cmd.AddCommand(NewExploreToURL())
	cmd.AddCommand(NewParseURL())
	return cmd
}

// NewExploreToURL creates a command to turn queries into URLs
func NewExploreToURL() *cobra.Command {
	var patchFile string
	var baseURL string
	var open bool
	cmd := &cobra.Command{
		Use: "build",
		Run: func(cmd *cobra.Command, args []string) {
			err := func() error {
				app := application.NewApp()
				if err := app.LoadConfig(cmd); err != nil {
					return err
				}
				if err := app.SetupLogging(); err != nil {
					return err
				}

				version.LogVersion()

				configDir := app.Config.GetConfigDir()

				bases, err := grafana.LoadGrafanaLinksInDir(configDir)
				if err != nil {
					return errors.Wrapf(err, "Error loading Grafana links from %v", configDir)
				}

				patch := &api.PanePatch{}

				patchData, err := os.ReadFile(patchFile)
				if err != nil {
					return errors.Wrapf(err, "Error reading patch file %v", patchFile)
				}

				if err := yaml.Unmarshal(patchData, patch); err != nil {
					return errors.Wrapf(err, "Couldn't unmarshal the patch in file %v", patchFile)
				}

				patcher := grafana.NewPatcher(grafana.RealClock{})
				link, err := patcher.ApplyPatch(bases, *patch)
				if err != nil {
					return errors.Wrapf(err, "Error applying patch")
				}
				u, err := grafana.LinkToURL(*link)
				fmt.Printf("Grafana URL:\n%v\n", u)
				if open {
					if err := browser.OpenURL(u); err != nil {
						return errors.Wrapf(err, "Error opening URL %v", u)
					}
				}
				return nil
			}()

			if err != nil {
				fmt.Printf("Error running request;\n %+v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&patchFile, "patch-file", "p", "", "A file containing the JSON object containing a map of pane IDs to panes")
	cmd.Flags().StringVarP(&baseURL, config.BaseURLFlagName, "", "", "The base URL for your grafana URLs.")
	cmd.Flags().BoolVarP(&open, "open", "", false, "Open the URL in a browser")
	return cmd
}

// NewParseURL creates a command to parse URLs
func NewParseURL() *cobra.Command {
	var panesFile string
	var logUrl string
	var name string
	cmd := &cobra.Command{
		Use: "parse",
		Run: func(cmd *cobra.Command, args []string) {
			err := func() error {
				app := application.NewApp()
				if err := app.LoadConfig(cmd); err != nil {
					return err
				}
				if err := app.SetupLogging(); err != nil {
					return err
				}

				version.LogVersion()

				var o io.Writer

				if panesFile != "" {
					f, err := os.Create(panesFile)
					if err != nil {
						return errors.Wrapf(err, "Error creating file %v", panesFile)
					}
					defer f.Close()

					if name == "" {
						// Default to the name of the file
						filename := filepath.Base(panesFile)

						// Strip the suffix (file extension)
						name = filename[:len(filename)-len(filepath.Ext(filename))]
					}

					o = f
				} else {
					o = os.Stdout
				}

				link, err := grafana.URLToLink(logUrl)
				if err != nil {
					return errors.Wrapf(err, "Error parsing URL")
				}
				link.Metadata.Name = name

				// Pretty print the json of the panes to the file
				encoder := yaml.NewEncoder(o)
				encoder.SetIndent(2)
				if err := encoder.Encode(link); err != nil {
					return errors.Wrapf(err, "Error writing panes to file")
				}
				return nil
			}()

			if err != nil {
				fmt.Printf("Error running request;\n %+v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&panesFile, "link-file", "o", "", "File to write the panes to. If not specified the panes will be written to stdout.")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name to give the resource when saving to a file")
	cmd.Flags().StringVarP(&logUrl, "url", "u", "", "The URL to parse")
	helpers.IgnoreError(cmd.MarkFlagRequired("url"))
	return cmd
}
