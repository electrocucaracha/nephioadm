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
	"io/ioutil"
	"os"

	"github.com/electrocucaracha/nephioadm/internal/kpt"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Runner interface {
	InstallSystem()
	InstallWebUI() error
	InstallConfigSync()
}

type NephioRunner struct {
	kpt.Client
	basePath          string
	gitServiceURI     string
	backendBaseUrl    string
	webUIClusterType  string
	packageOptions    kpt.PackageOptions
	debug             bool
	readResourceFunc  func(func(string) ([]byte, error), string, interface{}) error
	writeResourceFunc func(func(string) (*os.File, error), string, runtime.Object) error
}

type NephioRunnerOptions struct {
	BasePath      string
	NephioRepoURI string
	GitServiceURI string
	Debug         bool

	// Optional
	BackendBaseUrl   string
	WebUIClusterType string
}

var _ Runner = (*NephioRunner)(nil)

const (
	DefaultBasePath      = "/opt/nephio"
	DefaultWebUINodePort = 30007
)

func NewRunner(client kpt.Client,
	readResourceFunc func(func(string) ([]byte, error), string, interface{}) error,
	writeResourceFunc func(func(string) (*os.File, error), string, runtime.Object) error,
	opts *NephioRunnerOptions,
) *NephioRunner {
	r := &NephioRunner{
		Client:            client,
		readResourceFunc:  readResourceFunc,
		writeResourceFunc: writeResourceFunc,
		gitServiceURI:     opts.GitServiceURI,
		packageOptions: kpt.PackageOptions{
			RepoURI: opts.NephioRepoURI,
		},
		debug: opts.Debug,
	}

	r.basePath = DefaultBasePath
	if len(opts.BasePath) != 0 {
		r.basePath = opts.BasePath
	}

	if len(opts.BackendBaseUrl) != 0 {
		r.backendBaseUrl = opts.BackendBaseUrl
	}

	if len(opts.WebUIClusterType) != 0 {
		r.webUIClusterType = opts.WebUIClusterType
	}

	return r
}

func (r *NephioRunner) getPackage() {
	pkg := kpt.NewPackage(&r.packageOptions)
	r.PkgGet(pkg)

	if r.debug {
		r.PkgTree()
	}
}

func (r *NephioRunner) installPackage() {
	r.FnRender()

	if r.debug {
		r.PkgDiff()
	}

	r.LiveInit()
	r.LiveApply()

	if r.debug {
		r.LiveStatus()
	}
}

func (r *NephioRunner) InstallSystem() {
	r.SetLocalPath(r.basePath + "/system")
	r.packageOptions.Path = "nephio-system"

	r.getPackage()
	r.installPackage()
}

func (r *NephioRunner) setBackendBaseUrl(filename, backendBaseUrl string) error {
	var configMap v1.ConfigMap
	if err := r.readResourceFunc(ioutil.ReadFile, filename, &configMap); err != nil {
		return err
	}

	config := configMap.Data["app-config.nephio.yaml"]

	var backstageConfig map[string]interface{}
	if err := yaml.Unmarshal([]byte(config), &backstageConfig); err != nil {
		return err
	}

	backend, ok := backstageConfig["backend"].(map[interface{}]interface{})
	if !ok {
		return nil
	}

	backend["baseUrl"] = backendBaseUrl
	backstageConfig["backend"] = backend

	data, err := yaml.Marshal(backstageConfig)
	if err != nil {
		return err
	}

	configMap.Data["app-config.nephio.yaml"] = string(data[:])

	if err := r.writeResourceFunc(os.Create, filename, &configMap); err != nil {
		return err
	}

	return nil
}

func (r *NephioRunner) setClusterType(filename, clusterType string) error {
	var service v1.Service

	if err := r.readResourceFunc(ioutil.ReadFile, filename, &service); err != nil {
		return err
	}

	svcType := v1.ServiceType(clusterType)
	service.Spec.Type = svcType

	if svcType == v1.ServiceTypeNodePort {
		if len(service.Spec.Ports) > 1 {
			service.Spec.Ports[0].NodePort = DefaultWebUINodePort
		} else {
			service.Spec.Ports = append(service.Spec.Ports, v1.ServicePort{NodePort: DefaultWebUINodePort})
		}
	}

	if err := r.writeResourceFunc(os.Create, filename, &service); err != nil {
		return err
	}

	return nil
}

func (r *NephioRunner) InstallWebUI() error {
	r.SetLocalPath(r.basePath + "/webui")
	r.packageOptions.Path = "nephio-webui"

	r.getPackage()

	if len(r.backendBaseUrl) != 0 {
		if err := r.setBackendBaseUrl(r.basePath+"/webui/config-map.yaml", r.backendBaseUrl); err != nil {
			return err
		}
	}

	if r.webUIClusterType != string(v1.ServiceTypeClusterIP) {
		if err := r.setClusterType(r.basePath+"/webui/service.yaml", r.webUIClusterType); err != nil {
			return err
		}
	}

	r.installPackage()

	return nil
}

func (r *NephioRunner) InstallConfigSync() {
	r.SetLocalPath(r.basePath + "/configsync")
	r.packageOptions.Path = "nephio-configsync"

	r.getPackage()
	r.FnEval("gcr.io/kpt-fn/search-replace:v0.2", "spec.git.repo",
		"https://github.com/(.*)/(.*)", r.gitServiceURI+"/${2}")
	r.installPackage()
}
