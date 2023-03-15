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

	"github.com/electrocucaracha/nephioadm/internal/kpt"
	"k8s.io/apimachinery/pkg/runtime"
)

type Provider interface {
	Init(*NephioRunnerOptions) error
	Join(*NephioRunnerOptions) error
}

type NephioProvider struct {
	client        kpt.Client
	readResource  func(func(string) ([]byte, error), string, interface{}) error
	writeResource func(func(string) (*os.File, error), string, runtime.Object) error
}

var _ Provider = (*NephioProvider)(nil)

func NewProvider(client kpt.Client,
	readResourceFunc func(func(string) ([]byte, error), string, interface{}) error,
	writeResourceFunc func(func(string) (*os.File, error), string, runtime.Object) error,
) *NephioProvider {
	return &NephioProvider{
		client:        client,
		readResource:  readResourceFunc,
		writeResource: writeResourceFunc,
	}
}

func (p NephioProvider) Init(opts *NephioRunnerOptions) error {
	runner := NewRunner(p.client, p.readResource, p.writeResource, opts)
	runner.InstallSystem()

	if err := runner.InstallWebUI(); err != nil {
		return err
	}

	runner.InstallConfigSync()

	return nil
}

func (p NephioProvider) Join(opts *NephioRunnerOptions) error {
	runner := NewRunner(p.client, p.readResource, p.writeResource, opts)
	runner.InstallConfigSync()

	return nil
}
