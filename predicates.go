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
	"strings"

	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

const (
	// StorageNodePredicateKey is the key for StorageNodePredicate function.
	StorageNodePredicateKey = "StorageNode"
)

// Predicates type is the type for predicate functions.
type Predicates func(pod *v1.Pod, node v1.Node) (bool, []string, error)

// predicatesFuncMap is the a map of string(key) to predicate functions(value).
// Key describes the predicate behaviour.
var predicatesFuncMap = map[string]Predicates{
	StorageNodePredicateKey: StorageNodePredicate,
}

// predicatesEvalOrderStore is an array store for predicate function keys.
// The predicates are evaluated in order they are present in the store.
// If a new predicate is defined, it should be registered in the store.
var predicatesEvalOrderStore = []string{StorageNodePredicateKey}

// filter function filters nodes by using the defined predicates in the openebs scheduler extender.
func filter(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	// selectedNodes contains the list of nodes that are capable of scheduling a pod on them.
	var selectedNodes []v1.Node
	// failedNodes is a map which represents the filtered out nodes, with node names and failure messages
	failedNodes := make(schedulerapi.FailedNodesMap)
	pod := args.Pod
	for _, node := range args.Nodes.Items {
		fits, failureReasons, _ := isNodeCapable(pod, node)
		if fits {
			selectedNodes = append(selectedNodes, node)
		} else {
			failedNodes[node.Name] = strings.Join(failureReasons, ",")
		}
	}
	result := schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: selectedNodes,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}
	return &result
}

// isNodeCapable run all the predicates one by one and decides whether the pod can
// be scheduled on the node or not.
// TODO: Make the subsequent predicates operate on the nodes selected by the previous predicate.
// Currently, all predicated operate on all the nodes and if a node is rejected in any one predicate
// it is not selected for scheduling.
func isNodeCapable(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	fits := true
	failReasons := []string{}
	for _, predicateKey := range predicatesEvalOrderStore {
		fit, failures, err := predicatesFuncMap[predicateKey](pod, node)
		if err != nil {
			return false, nil, err
		}
		fits = fits && fit
		failReasons = append(failReasons, failures...)
	}
	return fits, failReasons, nil
}

// StorageNodePredicate selects nodes havin "storage=openebs" label.
func StorageNodePredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	if node.Labels["storage"] == "openebs" {
		log.Printf("Selected in StorageNodePredicate:pod %v/%v can be scheduled on node %v as it is a storage node\n", pod.Name, pod.Namespace, node.Name)
		return true, nil, nil
	}
	log.Printf("Rejected in StorageNodePredicate:pod %v/%v cannot be schediled on node %v as it is not a storage node\n", pod.Name, pod.Namespace, node.Name)
	return false, []string{"Node is not a sotrage node"}, nil
}
