// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/pingcap-incubator/tiup/pkg/localdata"
)

// DashboardConfig represent the data to generate Dashboard config
type DashboardConfig struct {
	ClusterName string
	DeployDir   string
}

// NewDashboardConfig returns a DashboardConfig
func NewDashboardConfig(cluster, deployDir string) *DashboardConfig {
	return &DashboardConfig{
		ClusterName: cluster,
		DeployDir:   deployDir,
	}
}

// Config read ${localdata.EnvNameComponentInstallDir}/templates/config/dashboard.yml
// and generate the config by ConfigWithTemplate
func (c *DashboardConfig) Config() (string, error) {
	fp := path.Join(os.Getenv(localdata.EnvNameComponentInstallDir), "templates", "config", "dashboard.yml.tpl")
	tpl, err := ioutil.ReadFile(fp)
	if err != nil {
		return "", err
	}
	return c.ConfigWithTemplate(string(tpl))
}

// ConfigWithTemplate generate the Dashboard config content by tpl
func (c *DashboardConfig) ConfigWithTemplate(tpl string) (string, error) {
	tmpl, err := template.New("dashboard").Parse(tpl)
	if err != nil {
		return "", err
	}

	content := bytes.NewBufferString("")
	if err := tmpl.Execute(content, c); err != nil {
		return "", err
	}

	return content.String(), nil
}