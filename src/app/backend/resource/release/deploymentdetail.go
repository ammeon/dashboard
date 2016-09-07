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
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

// ReleaseDetail is a presentation layer view of Kubernetes Release resource.
type ReleaseDetail struct {
	ObjectMeta common.ObjectMeta `json:"objectMeta"`

	Name   string `json:"name"`
	Status string `json:"status"`
}

// GetReleaseDetail returns model object of release and error, if any.
func GetReleaseDetail(client client.Interface, namespace string,
	name string) (*ReleaseDetail, error) {

	log.Printf("Getting details of %s release in %s namespace", name, namespace)

	return getReleaseDetail(release), nil
}

func getReleaseDetail(release *extensions.Release) {
	return &ReleaseDetail{
		ObjectMeta: common.NewObjectMeta(release.ObjectMeta),
		Name:       "happy-panda", // TODO: Releases
		Name:       "DEPLOYED",    // TODO: Releases
	}
}
