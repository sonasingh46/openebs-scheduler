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
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

// prioritize function will add score to hosts.
// Higher score mean high priority and thus will be preferred for scheduling.
// Note: The score added here will get added to the default scheduler (main car) score.
func prioritize(args schedulerapi.ExtenderArgs) *schedulerapi.HostPriorityList {
	nodes := args.Nodes.Items

	hostPriorityList := make(schedulerapi.HostPriorityList, len(nodes))
	for i, node := range nodes {
		// ToDO: More POC on assigning score to hosts.
		var score int
		score = 0
		if node.Labels["storage-high-priority"] == "true" {
			score = 1
			log.Printf("Node '%s' is high priority storage node", node.Name)
		}
		hostPriorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: score,
		}
	}

	return &hostPriorityList
}
