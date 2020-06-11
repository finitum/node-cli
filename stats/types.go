/*
Copyright 2015 The Kubernetes Authors.

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
package stats

// Summary is a top-level container for holding NodeStats and PodStats.
type Summary struct {
	// Node stats.
	Node NodeStats

	// Per-pod stats.
	Pods []PodStats
}

// NodeStats holds node-level stats.
type NodeStats struct {
	Name string

	UsageNanoCores uint64

	AvailableBytesMemory uint64
	UsageBytesMemory     uint64

	AvailableBytesEphemeral uint64
	CapacityBytesEphemeral  uint64
	UsedBytesEphemeral      uint64

	AvailableBytesStorage uint64
	CapacityBytesStorage  uint64
	UsedBytesStorage      uint64

	Pods uint64
}

// PodStats holds pod-level unprocessed sample stats.
type PodStats struct {
	PodRef PodReference

	UsageNanoCores     uint64
	UsageBytesMemory   uint64
	UsedBytesEphemeral uint64
	UsedBytesStorage   uint64
}

// PodReference contains enough information to locate the referenced pod.
type PodReference struct {
	Name      string
	Namespace string
	UID       string
}
