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
	"github.com/electrocucaracha/nephioadm/internal/app"
	"github.com/electrocucaracha/nephioadm/internal/kpt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockClient struct {
	SetLocalPathCallerCount int
	PkgGetCallerCount       int
	PkgTreeCallerCount      int
	PkgDiffCallerCount      int
	FnRenderCallerCount     int
	FnEvalCallerCount       int
	LiveInitCallerCount     int
	LiveApplyCallerCount    int
	LiveStatusCallerCount   int
}

func NewMockClient() *mockClient {
	return &mockClient{
		SetLocalPathCallerCount: 0,
		PkgGetCallerCount:       0,
		PkgTreeCallerCount:      0,
		PkgDiffCallerCount:      0,
		FnRenderCallerCount:     0,
		FnEvalCallerCount:       0,
		LiveInitCallerCount:     0,
		LiveApplyCallerCount:    0,
		LiveStatusCallerCount:   0,
	}
}

func (m *mockClient) SetLocalPath(localPath string) {
	m.SetLocalPathCallerCount += 1
}

func (m *mockClient) PkgGet(pkg *kpt.Package) {
	m.PkgGetCallerCount++
}

func (m *mockClient) PkgTree() {
	m.PkgTreeCallerCount += 1
}

func (m *mockClient) PkgDiff() {
	m.PkgDiffCallerCount += 1
}

func (m *mockClient) FnRender() {
	m.FnRenderCallerCount += 1
}

func (m *mockClient) FnEval(image, byPath, byValueRegex, putValue string) {
	m.FnEvalCallerCount += 1
}

func (m *mockClient) LiveInit() {
	m.LiveInitCallerCount += 1
}

func (m *mockClient) LiveApply() {
	m.LiveApplyCallerCount += 1
}

func (m *mockClient) LiveStatus() {
	m.LiveStatusCallerCount += 1
}

func NewNephioRunnerOptions(debug bool, args ...string) *app.NephioRunnerOptions {
	opts := &app.NephioRunnerOptions{Debug: debug}

	switch len(args) {
	case 5:
		opts.WebUIClusterType = args[4]

		fallthrough
	case 4:
		opts.BackendBaseUrl = args[3]

		fallthrough
	case 3:
		opts.GitServiceURI = args[2]

		fallthrough
	case 2:
		opts.NephioRepoURI = args[1]

		fallthrough
	case 1:
		opts.BasePath = args[0]
	}

	return opts
}

func (c *mockClient) checkCallerCountsFromProvider(debug bool, expected int) {
	Expect(c.SetLocalPathCallerCount).Should(Equal(expected))
	Expect(c.PkgGetCallerCount).Should(Equal(expected))
	Expect(c.FnRenderCallerCount).Should(Equal(expected))
	Expect(c.FnEvalCallerCount).Should(Equal(1))
	Expect(c.LiveInitCallerCount).Should(Equal(expected))
	Expect(c.LiveApplyCallerCount).Should(Equal(expected))

	if debug {
		Expect(c.PkgTreeCallerCount).Should(Equal(expected))
		Expect(c.PkgDiffCallerCount).Should(Equal(expected))
		Expect(c.LiveStatusCallerCount).Should(Equal(expected))
	} else {
		Expect(c.PkgTreeCallerCount).Should(Equal(0))
		Expect(c.PkgDiffCallerCount).Should(Equal(0))
		Expect(c.LiveStatusCallerCount).Should(Equal(0))
	}
}

var _ = Describe("Provider Service", func() {
	var provider app.NephioProvider
	var client *mockClient

	BeforeEach(func() {
		client = NewMockClient()
		provider = *app.NewProvider(client, fakeReadResourceFromFile, fakeWriteResourceToFile)
	})

	DescribeTable("initialization execution process", func(debug bool, args ...string) {
		err := provider.Init(NewNephioRunnerOptions(debug, args...))

		Expect(err).NotTo(HaveOccurred())
		client.checkCallerCountsFromProvider(debug, 3)
	},
		Entry("when the no options are provided", true),
		Entry("when Base path option is provided", true, "/opt/nephio"),
		Entry("when Base path and Nephio Repo URI options are provided", true,
			"/opt/nephio", "http://gitea/nephio-sandbox/packages.git"),
		Entry("when Base path, Nephio Repo URI and git service URI options are provided", true,
			"/opt/nephio", "http://gitea/nephio-internal/packages.git",
			"http://gitea/nephio-packages/"),
		Entry("when Base path, Nephio Repo URI, git service URI and backend BaseURL options are provided", true,
			"/opt/nephio", "http://gitea/nephio-internal/packages.git",
			"http://gitea/nephio-packages/", "https://CODESPACE_NAME-7007.preview.app.github.dev"),
		Entry("when Base path, Nephio Repo URI, git service URI, backend BaseURL and webUI Cluster type options are provided",
			true, "/opt/nephio", "http://gitea/nephio-internal/packages.git",
			"http://gitea/nephio-packages/", "https://CODESPACE_NAME-7007.preview.app.github.dev", "NodePort"),
		Entry("when the no options are provided and debug is disable", false),
	)

	DescribeTable("join execution process", func(debug bool, args ...string) {
		err := provider.Join(NewNephioRunnerOptions(debug, args...))

		Expect(err).NotTo(HaveOccurred())
		client.checkCallerCountsFromProvider(debug, 1)
	},
		Entry("when the no options are provided", true),
		Entry("when Base path option is provided", true, "/test/"),
		Entry("when Base path and Nephio Repo URI options are provided", true,
			"/test/", "http://gitea/nephio-sandbox/packages.git"),
		Entry("when Base path, Nephio Repo URI and git service URI options are provided", true,
			"/test/", "http://gitea/nephio-internal/packages.git",
			"http://gitea/nephio-packages/"),
		Entry("when the no options are provided and debug is disable", false),
	)
})
