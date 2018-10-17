/*
Copyright 2018 The OpenEBS Authors.

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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

// Liveness is the route handler for '/' GET.
func Liveness(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprint(w, "OpenEBS scheduler is running!\n")
}

// Filter is the route handler for '/filter' POST.
func Filter(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var buf bytes.Buffer
	var extenderArgs schedulerapi.ExtenderArgs
	var extenderFilterResult *schedulerapi.ExtenderFilterResult
	body := io.TeeReader(r.Body, &buf)
	err := json.NewDecoder(body).Decode(&extenderArgs)
	if err != nil {
		log.Println("Error in decoding extender arguments from scheduler main car:", err)
		extenderFilterResult = &schedulerapi.ExtenderFilterResult{
			Error: err.Error(),
		}
	} else {
		extenderFilterResult = filter(extenderArgs)
	}
	response, err := json.Marshal(extenderFilterResult)
	if err != nil {
		log.Println("Error in marshalling extended filter results:", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		log.Println("Nodes filtered successfully by the filter predicate")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// Prioritize is the route handler for '/prioritize' POST.
func Prioritize(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var buf bytes.Buffer
	var extenderArgs schedulerapi.ExtenderArgs
	var hostPriorityList *schedulerapi.HostPriorityList
	body := io.TeeReader(r.Body, &buf)
	err := json.NewDecoder(body).Decode(&extenderArgs)
	if err != nil {
		log.Println("Error in decoding extender arguments from scheduler main car:", err)
		hostPriorityList = &schedulerapi.HostPriorityList{}
	} else {
		hostPriorityList = prioritize(extenderArgs)
	}

	response, err := json.Marshal(hostPriorityList)
	if err != nil {
		log.Println("Error in marshalling host priority list:", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
