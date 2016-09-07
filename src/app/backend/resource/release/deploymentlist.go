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
	"github.com/kubernetes/dashboard/src/app/backend/resource/event"

	heapster "github.com/kubernetes/dashboard/src/app/backend/client"
	"k8s.io/kubernetes/pkg/api"
	k8serrors "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"

	"github.com/kubernetes/dashboard/src/app/backend/resource/dataselect"
	"github.com/kubernetes/dashboard/src/app/backend/resource/metric"
)

// ReplicationSetList contains a list of Releases in the cluster.
type ReleaseList struct {
	ListMeta common.ListMeta `json:"listMeta"`

	// Unordered list of Releases.
	Releases          []Release       `json:"releases"`
	CumulativeMetrics []metric.Metric `json:"cumulativeMetrics"`
}

// Release is a presentation layer view of Kubernetes Release resource. This means
// it is Release plus additional augumented data we can get from other sources
// (like services that target the same pods).
type Release struct {
	ObjectMeta common.ObjectMeta `json:"objectMeta"`
	TypeMeta   common.TypeMeta   `json:"typeMeta"`

	Name string `json:"name"`
}

// GetReleaseList returns a list of all Releases in the cluster.
func GetReleaseList(client client.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery, heapsterClient *heapster.HeapsterClient) (*ReleaseList, error) {
	log.Printf("Getting list of all releases in the cluster")

	channels := &common.ResourceChannels{
		ReleaseList: common.GetReleaseListChannel(client.Extensions(), nsQuery, 1),
		PodList:     common.GetPodListChannel(client, nsQuery, 1),
		EventList:   common.GetEventListChannel(client, nsQuery, 1),
	}

	return GetReleaseListFromChannels(channels, dsQuery, heapsterClient)
}

// GetReleaseList returns a list of all Releases in the cluster
// reading required resource list once from the channels.
func GetReleaseListFromChannels(channels *common.ResourceChannels,
	dsQuery *dataselect.DataSelectQuery, heapsterClient *heapster.HeapsterClient) (*ReleaseList, error) {

	releases := <-channels.ReleaseList.List
	if err := <-channels.ReleaseList.Error; err != nil {
		statusErr, ok := err.(*k8serrors.StatusError)
		if ok && statusErr.ErrStatus.Reason == "NotFound" {
			// NotFound - this means that the server does not support Release objects, which
			// is fine.
			emptyList := &ReleaseList{
				Releases: make([]Release, 0),
			}
			return emptyList, nil
		}
		return nil, err
	}

	pods := <-channels.PodList.List
	if err := <-channels.PodList.Error; err != nil {
		return nil, err
	}

	events := <-channels.EventList.List
	if err := <-channels.EventList.Error; err != nil {
		return nil, err
	}

	return CreateReleaseList([]string{"happy-panda"}), nil
}

// CreateReleaseList returns a list of all Release model objects in the cluster, based on all
// Kubernetes Release API objects.
func CreateReleaseList(releases []string) *ReleaseList {
	releaseList := &ReleaseList{
		Releases: make([]Release, 0),
		ListMeta: common.ListMeta{TotalItems: len(releases)},
	}

	for _, release := range releases {

		releaseList.Releases = append(releaseList.Releases,
			Release{
				ObjectMeta: common.NewObjectMeta(release.ObjectMeta),
				TypeMeta:   common.NewTypeMeta(common.ResourceKindRelease),
				name:       release,
				Pods:       podInfo,
			})
	}

	releaseList.CumulativeMetrics = []metric.Metric{}
	return releaseList
}
