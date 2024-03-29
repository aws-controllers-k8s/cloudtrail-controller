// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package v1alpha1

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EventDataStoreSpec defines the desired state of EventDataStore.
//
// A storage lake of event data against which you can run complex SQL-based
// queries. An event data store can include events that you have logged on your
// account from the last 90 to 2555 days (about three months to up to seven
// years). To select events for an event data store, use advanced event selectors
// (https://docs.aws.amazon.com/awscloudtrail/latest/userguide/logging-data-events-with-cloudtrail.html#creating-data-event-selectors-advanced).
type EventDataStoreSpec struct {

	// The advanced event selectors to use to select the events for the data store.
	// For more information about how to use advanced event selectors, see Log events
	// by using advanced event selectors (https://docs.aws.amazon.com/awscloudtrail/latest/userguide/logging-data-events-with-cloudtrail.html#creating-data-event-selectors-advanced)
	// in the CloudTrail User Guide.
	AdvancedEventSelectors []*AdvancedEventSelector `json:"advancedEventSelectors,omitempty"`
	// Specifies whether the event data store includes events from all regions,
	// or only from the region in which the event data store is created.
	MultiRegionEnabled *bool `json:"multiRegionEnabled,omitempty"`
	// The name of the event data store.
	// +kubebuilder:validation:Required
	Name *string `json:"name"`
	// Specifies whether an event data store collects events logged for an organization
	// in Organizations.
	OrganizationEnabled *bool `json:"organizationEnabled,omitempty"`
	// The retention period of the event data store, in days. You can set a retention
	// period of up to 2555 days, the equivalent of seven years.
	RetentionPeriod *int64 `json:"retentionPeriod,omitempty"`
	Tags            []*Tag `json:"tags,omitempty"`
	// Specifies whether termination protection is enabled for the event data store.
	// If termination protection is enabled, you cannot delete the event data store
	// until termination protection is disabled.
	TerminationProtectionEnabled *bool `json:"terminationProtectionEnabled,omitempty"`
}

// EventDataStoreStatus defines the observed state of EventDataStore
type EventDataStoreStatus struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	// +kubebuilder:validation:Optional
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRS managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	// +kubebuilder:validation:Optional
	Conditions []*ackv1alpha1.Condition `json:"conditions"`
	// The timestamp that shows when the event data store was created.
	// +kubebuilder:validation:Optional
	CreatedTimestamp *metav1.Time `json:"createdTimestamp,omitempty"`
	// The status of event data store creation.
	// +kubebuilder:validation:Optional
	Status *string `json:"status,omitempty"`
	// The timestamp that shows when an event data store was updated, if applicable.
	// UpdatedTimestamp is always either the same or newer than the time shown in
	// CreatedTimestamp.
	// +kubebuilder:validation:Optional
	UpdatedTimestamp *metav1.Time `json:"updatedTimestamp,omitempty"`
}

// EventDataStore is the Schema for the EventDataStores API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type EventDataStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              EventDataStoreSpec   `json:"spec,omitempty"`
	Status            EventDataStoreStatus `json:"status,omitempty"`
}

// EventDataStoreList contains a list of EventDataStore
// +kubebuilder:object:root=true
type EventDataStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventDataStore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EventDataStore{}, &EventDataStoreList{})
}
