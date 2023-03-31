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
	internal "github.com/electrocucaracha/nephioadm/internal/app"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewInitCommand(provider internal.Provider) *cobra.Command {
	var globalOpts GlobalOptions

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Run this command in order to set up the Nephio control plane",
		RunE: func(cmd *cobra.Command, args []string) error {
			backendBaseUrl, _ := cmd.Flags().GetString("backend-base-url")
			webUIClusterType, _ := cmd.Flags().GetString("webui-cluster-type")

			runnerOpts := &internal.NephioRunnerOptions{
				BasePath:         globalOpts.basePath,
				NephioRepoURI:    globalOpts.nephioRepoURI,
				GitServiceURI:    globalOpts.gitServiceURI,
				BackendBaseUrl:   backendBaseUrl,
				WebUIClusterType: webUIClusterType,
				Debug:            globalOpts.debug,
			}

			if err := provider.Init(runnerOpts); err != nil {
				return errors.Wrap(err, "failed to init nephio cluster plane")
			}

			return nil
		},
	}

	cmd.Flags().String("backend-base-url", "http://localhost:7007", "Nephio WebUI URL")
	cmd.Flags().String("webui-cluster-type", "NodePort", "Nephio WebUI Cluster Type")

	cmd = GetCommandFlags(cmd, &globalOpts)

	return cmd
}
