// Copyright 2015 Google, Inc All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
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
	"log"
	"net/http"
)

func GetClusterState(w http.ResponseWriter, r *http.Request) {
	clusterStatus, err := clusterStateManager.ClusterStatus()
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}
	data, err := json.MarshalIndent(clusterStatus, "", "  ")
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}
	w.Write(data)
}

func SubmitClusterState(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var cs ClusterState
	err := decoder.Decode(&cs)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}
	if err := clusterStateManager.Submit(&cs); err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}
}

func RemoveClusterState(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		log.Println("error removing cluster state: missing name parameter")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	clusterStateManager.Remove(name)
}
