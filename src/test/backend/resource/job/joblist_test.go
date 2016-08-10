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

package job

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kubernetes/dashboard/src/app/backend/resource/common"
	"k8s.io/kubernetes/pkg/api"
	k8serrors "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/batch"
)

func TestGetJobListFromChannels(t *testing.T) {
	var jobCompletions int32 = 21
	cases := []struct {
		k8sRs         batch.JobList
		k8sRsError    error
		pods          *api.PodList
		expected      *JobList
		expectedError error
	}{
		{
			batch.JobList{},
			nil,
			&api.PodList{},
			&JobList{common.ListMeta{}, []Job{}},
			nil,
		},
		{
			batch.JobList{},
			errors.New("MyCustomError"),
			&api.PodList{},
			nil,
			errors.New("MyCustomError"),
		},
		{
			batch.JobList{},
			&k8serrors.StatusError{},
			&api.PodList{},
			nil,
			&k8serrors.StatusError{},
		},
		{
			batch.JobList{},
			&k8serrors.StatusError{ErrStatus: unversioned.Status{}},
			&api.PodList{},
			nil,
			&k8serrors.StatusError{ErrStatus: unversioned.Status{}},
		},
		{
			batch.JobList{},
			&k8serrors.StatusError{ErrStatus: unversioned.Status{Reason: "foo-bar"}},
			&api.PodList{},
			nil,
			&k8serrors.StatusError{ErrStatus: unversioned.Status{Reason: "foo-bar"}},
		},
		{
			batch.JobList{},
			&k8serrors.StatusError{ErrStatus: unversioned.Status{Reason: "NotFound"}},
			&api.PodList{},
			&JobList{
				Jobs: make([]Job, 0),
			},
			nil,
		},
		{
			batch.JobList{
				Items: []batch.Job{{
					ObjectMeta: api.ObjectMeta{
						Name:              "rs-name",
						Namespace:         "rs-namespace",
						Labels:            map[string]string{"key": "value"},
						CreationTimestamp: unversioned.Unix(111, 222),
					},
					Spec: batch.JobSpec{
						Selector:    &unversioned.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
						Completions: &jobCompletions,
					},
					Status: batch.JobStatus{
						Active: 7,
					},
				},
					{
						ObjectMeta: api.ObjectMeta{
							Name:              "rs-name",
							Namespace:         "rs-namespace",
							Labels:            map[string]string{"key": "value"},
							CreationTimestamp: unversioned.Unix(111, 222),
						},
						Spec: batch.JobSpec{
							Selector: &unversioned.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
						},
						Status: batch.JobStatus{
							Active: 7,
						},
					},
				},
			},
			nil,
			&api.PodList{
				Items: []api.Pod{
					{
						ObjectMeta: api.ObjectMeta{
							Namespace: "rs-namespace",
							Labels:    map[string]string{"foo": "bar"},
						},
						Status: api.PodStatus{Phase: api.PodFailed},
					},
					{
						ObjectMeta: api.ObjectMeta{
							Namespace: "rs-namespace",
							Labels:    map[string]string{"foo": "baz"},
						},
						Status: api.PodStatus{Phase: api.PodFailed},
					},
				},
			},
			&JobList{
				common.ListMeta{TotalItems: 2},
				[]Job{{
					ObjectMeta: common.ObjectMeta{
						Name:              "rs-name",
						Namespace:         "rs-namespace",
						Labels:            map[string]string{"key": "value"},
						CreationTimestamp: unversioned.Unix(111, 222),
					},
					TypeMeta: common.TypeMeta{Kind: common.ResourceKindJob},
					Pods: common.PodInfo{
						Current:  7,
						Desired:  21,
						Failed:   1,
						Warnings: []common.Event{},
					},
				}, {
					ObjectMeta: common.ObjectMeta{
						Name:              "rs-name",
						Namespace:         "rs-namespace",
						Labels:            map[string]string{"key": "value"},
						CreationTimestamp: unversioned.Unix(111, 222),
					},
					TypeMeta: common.TypeMeta{Kind: common.ResourceKindJob},
					Pods: common.PodInfo{
						Current:  7,
						Desired:  0,
						Failed:   1,
						Warnings: []common.Event{},
					},
				}},
			},
			nil,
		},
	}

	for _, c := range cases {
		channels := &common.ResourceChannels{
			JobList: common.JobListChannel{
				List:  make(chan *batch.JobList, 1),
				Error: make(chan error, 1),
			},
			NodeList: common.NodeListChannel{
				List:  make(chan *api.NodeList, 1),
				Error: make(chan error, 1),
			},
			ServiceList: common.ServiceListChannel{
				List:  make(chan *api.ServiceList, 1),
				Error: make(chan error, 1),
			},
			PodList: common.PodListChannel{
				List:  make(chan *api.PodList, 1),
				Error: make(chan error, 1),
			},
			EventList: common.EventListChannel{
				List:  make(chan *api.EventList, 1),
				Error: make(chan error, 1),
			},
		}

		channels.JobList.Error <- c.k8sRsError
		channels.JobList.List <- &c.k8sRs

		channels.NodeList.List <- &api.NodeList{}
		channels.NodeList.Error <- nil

		channels.ServiceList.List <- &api.ServiceList{}
		channels.ServiceList.Error <- nil

		channels.PodList.List <- c.pods
		channels.PodList.Error <- nil

		channels.EventList.List <- &api.EventList{}
		channels.EventList.Error <- nil

		actual, err := GetJobListFromChannels(channels, common.NoDataSelect)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetJobListFromChannels() ==\n          %#v\nExpected: %#v", actual, c.expected)
		}
		if !reflect.DeepEqual(err, c.expectedError) {
			t.Errorf("GetJobListFromChannels() ==\n          %#v\nExpected: %#v", err, c.expectedError)
		}
	}
}
