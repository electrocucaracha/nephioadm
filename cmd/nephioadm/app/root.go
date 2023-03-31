/*
Copyright © 2023

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"os"

	internal "github.com/electrocucaracha/nephioadm/internal/app"
	"github.com/electrocucaracha/nephioadm/internal/k8s"
	"github.com/electrocucaracha/nephioadm/internal/kpt"
	"github.com/spf13/cobra"
)

type GlobalOptions struct {
	basePath      string
	nephioRepoURI string
	gitServiceURI string
	debug         bool
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nephioadm",
		Short: "nephioadm: easily bootstrap Nephio cluster",
	}

	provider := internal.NewProvider(&kpt.CommandLine{}, k8s.ReadResourceFromFile,
		k8s.WriteResourceToFile)

	cmd.AddCommand(NewInitCommand(provider))
	cmd.AddCommand(NewJoinCommand(provider))

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

func GetCommandFlags(cmd *cobra.Command, opts *GlobalOptions) *cobra.Command {
	flags := cmd.Flags()

	flags.StringVar(&opts.basePath, "base-path", internal.DefaultBasePath,
		"The local directory to write the Nephio's packages to")
	flags.StringVar(&opts.nephioRepoURI, "nephio-repo", "https://github.com/nephio-project/nephio-packages.git",
		"URI of a git repository containing Nephio's packages (System, WebUI, ConfigSync) as subdirectories")
	flags.StringVar(&opts.gitServiceURI, "git-service", "https://github.com/nephio-test/",
		"URI of a Git Service")
	flags.BoolVar(&opts.debug, "debug", false, "Enable debug mode")

	return cmd
}
