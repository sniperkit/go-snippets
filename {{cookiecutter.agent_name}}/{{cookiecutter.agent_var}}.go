// Copyright (C) 2016 {{cookiecutter.author_name}}
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/crucibuild/sdk-agent-go/agentiface"
	"github.com/crucibuild/sdk-agent-go/agentimpl"
)

// Resources represents an handler on the various data files
// Used by the agent(avro files, manifest, etc...).
var Resources http.FileSystem

// {{cookiecutter.agent_var}} is an implementation over the Agent implementation
// available in sdk-agent-go.
type {{cookiecutter.agent_var}} struct {
	*agentimpl.Agent
}

// nolint: unparam
func mustOpenResources(path string) []byte {
	file, err := Resources.Open(path)

	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadAll(file)

	if err != nil {
		panic(err)
	}

	return content
}

// New{{cookiecutter.agent_var}} creates a new instance of {{cookiecutter.agent_var}}.
func New{{cookiecutter.agent_var}}() (agentiface.Agent, error) {
	var agentSpec map[string]interface{}

	manifest := mustOpenResources("/resources/manifest.json")

	err := json.Unmarshal(manifest, &agentSpec)

	if err != nil {
		return nil, err
	}

	impl, err := agentimpl.NewAgent(agentimpl.NewManifest(agentSpec))

	if err != nil {
		return nil, err
	}

	agent := &{{cookiecutter.agent_var}}{
		impl,
	}

	if err := agent.init(); err != nil {
		return nil, err
	}

	return agent, nil
}

func (a *{{cookiecutter.agent_var}}) init() error {
	return nil
}
