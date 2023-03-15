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

package k8s

import (
	"bytes"
	_ "log"
	"os"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
)

// ReadResourceFromFile decodes the file document in a Kubernetes resource.
func ReadResourceFromFile(readYamlFunc func(string) ([]byte, error), path string, into interface{}) error {
	data, err := readYamlFunc(path)
	if err != nil {
		return errors.Wrapf(err, "failed to read the %s resource file", path)
	}

	if err := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), len(data)).Decode(&into); err != nil {
		return errors.Wrapf(err, "failed to decode the data %v", data)
	}

	return nil
}

// WriteResourceToFile serializes the Kubernetes resource provided into a file document.
func WriteResourceToFile(createFunc func(string) (*os.File, error),
	path string, resource runtime.Object,
) error {
	file, err := createFunc(path)
	if err != nil {
		return errors.Wrapf(err, "failed to write the Kubernetes resource into the %s resource file", path)
	}

	printr := printers.NewTypeSetter(scheme.Scheme).ToPrinter(&printers.YAMLPrinter{})
	if err := printr.PrintObj(resource, file); err != nil {
		return errors.Wrap(err, "failed to marshal the Kubernetes resource")
	}

	return nil
}
