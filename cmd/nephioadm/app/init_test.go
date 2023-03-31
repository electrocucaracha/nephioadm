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

package app_test

import (
	"github.com/electrocucaracha/nephioadm/cmd/nephioadm/app"
	internal "github.com/electrocucaracha/nephioadm/internal/app"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

type mock struct {
	Opts *internal.NephioRunnerOptions
}

func (m *mock) Init(opts *internal.NephioRunnerOptions) error {
	m.Opts = opts

	return nil
}

var _ = Describe("Init Command", func() {
	var provider mock
	var cmd *cobra.Command
	testData := &internal.NephioRunnerOptions{
		BasePath:         "/tmp",
		NephioRepoURI:    "http://gitea:3000/playground/test.git",
		GitServiceURI:    "http://gitea:3000/nephio-test",
		BackendBaseUrl:   "https://codespace-7007.preview.app.github.dev",
		WebUIClusterType: "LoadBalancer",
		Debug:            true,
	}

	BeforeEach(func() {
		provider = mock{}
		cmd = app.NewInitCommand(&provider)
	})

	DescribeTable("initialization execution process", func(shouldSucceed bool, args ...string) {
		cmd.SetArgs(args)
		err := cmd.Execute()

		if shouldSucceed {
			Expect(err).NotTo(HaveOccurred())
			if len(args) > 0 {
				Expect(testData).To(Equal(provider.Opts))
			}
		} else {
			Expect(err).To(HaveOccurred())
		}
	},
		Entry("when the default options are provided", true),
		Entry("when all options are defined", true,
			"--base-path", testData.BasePath,
			"--nephio-repo", testData.NephioRepoURI,
			"--git-service", testData.GitServiceURI,
			"--backend-base-url", testData.BackendBaseUrl,
			"--webui-cluster-type", testData.WebUIClusterType,
			"--debug"),
		Entry("when invalid option is provided", false, "--invalid"),
	)
})
