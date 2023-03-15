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

package kpt

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

type Package struct {
	repoURI string
	path    string
	version string
}

type PackageOptions struct {
	RepoURI string
	Path    string

	// Version is optional
	Version string
}

func (p Package) String() string {
	repo := p.repoURI + "/" + p.path
	if len(p.version) != 0 {
		return repo + "@" + p.version
	}

	return repo
}

func NewPackage(opts *PackageOptions) *Package {
	p := &Package{
		repoURI: opts.RepoURI,
		path:    opts.Path,
	}

	if len(opts.Version) != 0 {
		p.version = opts.Version
	}

	return p
}

type Client interface {
	PkgGet(*Package)
	PkgTree()
	PkgDiff()
	FnRender()
	FnEval(string, string, string, string)
	LiveInit()
	LiveApply()
	LiveStatus()
	SetLocalPath(string)
}

type CommandLine struct {
	localPath string
}

var _ Client = (*CommandLine)(nil)

func (c *CommandLine) SetLocalPath(localPath string) {
	c.localPath = localPath
}

func (c CommandLine) runCmd(args ...string) {
	kptExecPath, err := exec.LookPath("kpt")
	if err != nil {
		log.Println(err)
	}

	command := &exec.Cmd{
		Path:   kptExecPath,
		Args:   append([]string{kptExecPath}, args...),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	// set var to get the output
	var out bytes.Buffer

	// set the output to our variable
	command.Stdout = &out
	if err := command.Run(); err != nil {
		log.Println(err)
	}
}

func (c *CommandLine) PkgGet(pkg *Package) {
	_, err := os.Stat(c.localPath)

	if os.IsExist(err) {
		return
	}

	args := []string{"pkg", "get", pkg.String(), c.localPath, "--for-deployment"}
	c.runCmd(args...)
}

func (c *CommandLine) PkgTree() {
	args := []string{"pkg", "tree", c.localPath}
	c.runCmd(args...)
}

func (c *CommandLine) PkgDiff() {
	args := []string{"pkg", "diff", c.localPath}
	c.runCmd(args...)
}

func (c *CommandLine) FnRender() {
	args := []string{"fn", "render", c.localPath}
	c.runCmd(args...)
}

func (c *CommandLine) FnEval(image, byPath, byValueRegex, putValue string) {
	args := []string{
		"fn", "eval", c.localPath, "--save",
		"--type", "mutator", "--image", image, "--",
		"by-path=" + byPath, "by-value-regex=" + byValueRegex,
		"put-value=" + putValue,
	}
	c.runCmd(args...)
}

func (c *CommandLine) LiveInit() {
	args := []string{"live", "init", c.localPath, "--force"}
	c.runCmd(args...)
}

func (c *CommandLine) LiveApply() {
	args := []string{"live", "apply", c.localPath, "--reconcile-timeout", "15m"}
	c.runCmd(args...)
}

func (c *CommandLine) LiveStatus() {
	args := []string{"live", "status", c.localPath}
	c.runCmd(args...)
}
