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
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

const (
	// LivenessPrefix is the address prefix which can be used to test for liveness or readiness
	// of OpenEBS scheduler .
	LivenessPrefix = "/"
	// FilterPrefix is the address prefix where the default scheduler acting as main car will
	// request for OpenEBS scheduling predicates.
	// Predicates are used to reject node that are unfit for scheduling a pod on it.
	FilterPrefix = "/filter"
	// PriorityPrefix is the address prefix where the default scheduler acting as main car will
	// request for OpenEBS scheduling priorities.
	// Priorities are used to give weightage to nodes for the most preferable scheduling of pod to a node.
	// PriorityPrefix is not called in case of single node cluster.
	PriorityPrefix = "/prioritize"
)

func main() {
	// Create a new router object
	router := httprouter.New()
	// Register LivenessPrefix route with GET verb action.
	router.GET(LivenessPrefix, Liveness)
	// Register PredicatePrefix route with POST verb action.
	router.POST(FilterPrefix, Filter)
	// Register PriorityPrefix route with POST verb action.
	router.POST(PriorityPrefix, Prioritize)
	log.Println("Openebs-scheduler listening on http://localhost:80")
	log.Fatal(http.ListenAndServe(":80", router))
}
