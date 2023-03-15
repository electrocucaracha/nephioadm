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
	"bytes"
	"errors"
	"io/ioutil"
	"os"

	"github.com/electrocucaracha/nephioadm/internal/app"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func fakeReadYamlFile(filename string) ([]byte, error) {
	testdata := map[string]string{
		"/opt/nephio/webui/config-map.yaml": `apiVersion: v1
kind: ConfigMap
metadata:
  name: nephio-webui-config
data:
  app-config.nephio.yaml: |
    backend:
      baseUrl: http://localhost:7007`,
		"/opt/nephio/webui/service.yaml": `apiVersion: v1
kind: Service
metadata:
  name: nephio-webui
  namespace: nephio-webui
spec:
  selector:
    app: nephio-webui
  ports:
    - name: http
      port: 7007
      targetPort: http`,
	}

	val, ok := testdata[filename]
	if ok {
		return ioutil.ReadAll(bytes.NewBufferString(val))
	}

	return nil, errors.New("non-existing file")
}

func fakeWriteResourceToFile(createFunc func(string) (*os.File, error),
	path string, resource runtime.Object,
) error {
	return nil
}

func fakeReadResourceFromFile(readYamlFunc func(string) ([]byte, error), path string, into interface{}) error {
	data, err := fakeReadYamlFile(path)
	if err != nil {
		return err
	}

	if err := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), len(data)).Decode(&into); err != nil {
		return err
	}

	return nil
}

func (c *mockClient) checkCallerCountsFromRunner(debug bool) {
	Expect(c.SetLocalPathCallerCount).Should(Equal(1))
	Expect(c.PkgGetCallerCount).Should(Equal(1))
	Expect(c.FnRenderCallerCount).Should(Equal(1))
	Expect(c.LiveInitCallerCount).Should(Equal(1))
	Expect(c.LiveApplyCallerCount).Should(Equal(1))

	if debug {
		Expect(c.PkgTreeCallerCount).Should(Equal(1))
		Expect(c.PkgDiffCallerCount).Should(Equal(1))
		Expect(c.LiveStatusCallerCount).Should(Equal(1))
	} else {
		Expect(c.PkgTreeCallerCount).Should(Equal(0))
		Expect(c.PkgDiffCallerCount).Should(Equal(0))
		Expect(c.LiveStatusCallerCount).Should(Equal(0))
	}
}

var _ = Describe("Nephio Runner", func() {
	DescribeTable("install System package", func(debug bool, args ...string) {
		client := NewMockClient()
		app.NewRunner(client, fakeReadResourceFromFile, fakeWriteResourceToFile,
			&app.NephioRunnerOptions{Debug: debug}).InstallSystem()
		client.checkCallerCountsFromRunner(debug)
		Expect(client.FnEvalCallerCount).Should(Equal(0))
	},
		Entry("when the no options are provided", true),
		Entry("when the no options are provided and debug is disabled", false),
	)

	DescribeTable("install Web UI package", func(debug bool, args ...string) {
		client := NewMockClient()
		opts := &app.NephioRunnerOptions{Debug: debug}
		if len(args) > 1 {
			opts.BackendBaseUrl = args[0]
		}
		err := app.NewRunner(client, fakeReadResourceFromFile, fakeWriteResourceToFile, opts).InstallWebUI()
		Expect(err).NotTo(HaveOccurred())
		client.checkCallerCountsFromRunner(debug)
		Expect(client.FnEvalCallerCount).Should(Equal(0))
	},
		Entry("when the no options are provided", true),
		Entry("when the no options are provided and debug is disabled", false),
		Entry("when a backend base URL is provided", false, "https://codespace-7007.preview.app.github.dev"),
	)

	DescribeTable("install ConfigSync package", func(debug bool, args ...string) {
		client := NewMockClient()
		app.NewRunner(client, fakeReadResourceFromFile, fakeWriteResourceToFile,
			&app.NephioRunnerOptions{Debug: debug}).InstallConfigSync()
		client.checkCallerCountsFromRunner(debug)
		Expect(client.FnEvalCallerCount).Should(Equal(1))
	},
		Entry("when the no options are provided", true),
		Entry("when the no options are provided and debug is disabled", false),
	)
})
