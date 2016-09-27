// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package release

import (
	"log"

	"github.com/kubernetes/dashboard/src/app/backend/resource/common"

	heapster "github.com/kubernetes/dashboard/src/app/backend/client"
	k8serrors "k8s.io/kubernetes/pkg/api/errors"
	client "k8s.io/kubernetes/pkg/client/unversioned"

	"github.com/kubernetes/dashboard/src/app/backend/resource/dataselect"
	"github.com/kubernetes/dashboard/src/app/backend/resource/metric"
)

// ReleaseList contains a list of Releases in the cluster.
type ReleaseList struct {
	ListMeta common.ListMeta `json:"listMeta"`

	// Unordered list of Releases.
	Releases          []common.Release `json:"releases"`
	CumulativeMetrics []metric.Metric  `json:"cumulativeMetrics"`
}

// GetReleaseList returns a list of all Releases in the cluster.
func GetReleaseList(client client.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery, heapsterClient *heapster.HeapsterClient) (*ReleaseList, error) {
	log.Printf("Getting list of all releases in the cluster")

	channels := &common.ResourceChannels{
		ReleaseList: common.GetReleaseListChannel(nsQuery, 1),
		PodList:     common.GetPodListChannel(client, nsQuery, 1),
		EventList:   common.GetEventListChannel(client, nsQuery, 1),
	}

	return GetReleaseListFromChannels(channels, dsQuery, heapsterClient)
}

// GetReleaseList returns a list of all Releases in the cluster
// reading required resource list once from the channels.
func GetReleaseListFromChannels(channels *common.ResourceChannels,
	dsQuery *dataselect.DataSelectQuery, heapsterClient *heapster.HeapsterClient) (*ReleaseList, error) {

	//releases := <-channels.ReleaseList.List
	// TODO: Releases
	if err := <-channels.ReleaseList.Error; err != nil {
		statusErr, ok := err.(*k8serrors.StatusError)
		if ok && statusErr.ErrStatus.Reason == "NotFound" {
			// NotFound - this means that the server does not support Release objects, which
			// is fine.
			emptyList := &ReleaseList{
				Releases: make([]common.Release, 0),
			}
			return emptyList, nil
		}
		return nil, err
	}

	// pods := <-channels.PodList.List
	// TODO: Releases
	if err := <-channels.PodList.Error; err != nil {
		return nil, err
	}

	// events := <-channels.EventList.List
	// TODO: Releases
	if err := <-channels.EventList.Error; err != nil {
		return nil, err
	}

	return CreateReleaseList([]string{"happy-panda"}), nil
}

// CreateReleaseList returns a list of all Release model objects in the cluster, based on all
// Kubernetes Release API objects.
func CreateReleaseList(releases []string) *ReleaseList {
	releaseList := &ReleaseList{
		Releases: make([]common.Release, 0),
		ListMeta: common.ListMeta{TotalItems: len(releases)},
	}

	for _, release := range releases {

		releaseList.Releases = append(releaseList.Releases,
			common.Release{
				Name:      release,
				Namespace: "default", // TODO: Releases
			})
	}

	releaseList.CumulativeMetrics = []metric.Metric{}
	return releaseList
}
