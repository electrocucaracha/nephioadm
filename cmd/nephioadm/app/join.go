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

func NewJoinCommand(provider internal.Provider) *cobra.Command {
	var opts GlobalOptions

	cmd := &cobra.Command{
		Use:   "join",
		Short: "Run this command in order to join a Cluster to the existing Nephio control plane",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := &internal.NephioRunnerOptions{
				BasePath:      opts.basePath,
				NephioRepoURI: opts.nephioRepoURI,
				GitServiceURI: opts.gitServiceURI,
				Debug:         opts.debug,
			}

			if err := provider.Join(opts); err != nil {
				return errors.Wrap(err, "failed to join to the nephio cluster plane")
			}

			return nil
		},
	}

	cmd = GetCommandFlags(cmd, &opts)

	return cmd
}
