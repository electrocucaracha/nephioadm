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

package k8s_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"

	"github.com/electrocucaracha/nephioadm/internal/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func fakeYamlReadFile(filename string) ([]byte, error) {
	testdata := map[string]string{
		"/tmp/configmap.yml": `apiVersion: v1
kind: ConfigMap
metadata:
  name: nephio-webui-config
data:
  app-config.nephio.yaml: |
    backend:
      baseUrl: http://localhost:7007`,
		"/tmp/invalid.yml": "error=true",
	}

	val, ok := testdata[filename]
	if ok {
		return ioutil.ReadAll(bytes.NewBufferString(val))
	}

	return nil, errors.New("non-existing file")
}

func fakeCreate(name string) (*os.File, error) {
	return os.Stdout, nil
}

var _ = Describe("Kubernetes File Parser", func() {
	DescribeTable("reads resource from a file", func(shouldSucceed bool, path string, resource runtime.Object) {
		var got v1.ConfigMap
		err := k8s.ReadResourceFromFile(fakeYamlReadFile, path, &got)

		if shouldSucceed {
			Expect(err).NotTo(HaveOccurred())
		} else {
			Expect(err).To(HaveOccurred())
			Expect(got).Should(Equal(got))
		}

		switch resource.(type) {
		case *v1.ConfigMap:
			configMap, ok := resource.(*v1.ConfigMap)
			Expect(ok).Should(BeTrue())
			Expect(got.Data).Should(Equal(configMap.Data))
		case *v1.Service:
			service, ok := resource.(*v1.Service)
			Expect(ok).Should(BeTrue())
			Expect(got.Data).Should(Equal(service.Spec.Ports))
		}
	},
		Entry("when a valid ConfigMap file is read", true, "/tmp/configmap.yml",
			&v1.ConfigMap{Data: map[string]string{"app-config.nephio.yaml": "backend:\n  baseUrl: http://localhost:7007\n"}}),
		Entry("when a non existing file is read", false, "/tmp/non-existing.yml", nil),
		Entry("when an invalid YAML file is read", false, "/tmp/invalid.yml", nil),
	)

	DescribeTable("writes resource into a file", func(shouldSucceed bool, path string, resource runtime.Object) {
		err := k8s.WriteResourceToFile(fakeCreate, path, resource)

		if shouldSucceed {
			Expect(err).NotTo(HaveOccurred())
		} else {
			Expect(err).To(HaveOccurred())
		}
	},
		Entry("when a valid ConfigMap resource is written", true, "/tmp/configmap.yml",
			&v1.ConfigMap{Data: map[string]string{"app-config.nephio.yaml": "backend:\n  baseUrl: http://localhost:7007\n"}}),
		Entry("when a valid Service resource is written", true, "/tmp/service.yml",
			&v1.Service{Spec: v1.ServiceSpec{Ports: []v1.ServicePort{{Name: "http", Port: 7007}}}}),
	)
})
