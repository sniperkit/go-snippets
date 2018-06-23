/*
Copyright 2016 The Kubernetes Authors.

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

package utils

import (
	"fmt"
	"os"

	restclient "k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kclientset "k8s.io/client-go/kubernetes"
)

// Namespace will return the value of KUBECTL_PLUGINS_CURRENT_NAMESPACE env var
func Namespace() string {
	return os.Getenv("KUBECTL_PLUGINS_CURRENT_NAMESPACE")
}

func CheckNamespaceExists(ns string, config *restclient.Config) {
	fmt.Printf("Looking up Namespace %s...\n", Entity(ns))
	kubeClient, err := kclientset.NewForConfig(config)
	if err != nil {
		Exit1(fmt.Sprintf("%v", err))
	}

	if _, err := kubeClient.Core().Namespaces().Get(ns, metav1.GetOptions{}); err != nil {
		Exit1(fmt.Sprintf("%v", err))
	}
}

func Loglevel() (flagName, flagValue string) {
	kubeLoglevel := os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_V")
	otherLoglevel := os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_LOGLEVEL")
	if len(otherLoglevel) > 0 {
		return "--loglevel", otherLoglevel
	}
	if len(kubeLoglevel) == 0 {
		kubeLoglevel = "0"
	}
	return "--v", kubeLoglevel
}

// Exit1 will print the specified error string to the screen and
// then stop the program, with an exit code of 1
func Exit1(errStr string) {
	Error(errStr)
	os.Exit(1)
}
