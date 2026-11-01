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

package event_data_store

import (
	"context"
	"fmt"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"

	svcapitypes "github.com/aws-controllers-k8s/cloudtrail-controller/apis/v1alpha1"
)

// setResourceAdditionalFields will set the fields that are not returned by
// GetEventDataStore calls
func (rm *resourceManager) setResourceAdditionalFields(
	ctx context.Context,
	ko *svcapitypes.EventDataStore,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.setResourceAdditionalFields")
	defer func() { exit(err) }()

	if ko.Status.ACKResourceMetadata != nil && ko.Status.ACKResourceMetadata.ARN != nil &&
		*ko.Status.ACKResourceMetadata.ARN != "" {
		// Set event data store tags
		ko.Spec.Tags, err = rm.getTags(ctx, string(*ko.Status.ACKResourceMetadata.ARN))
		if err != nil {
			return err
		}
	}

	return nil
}

// syncTags updates a EDS list of tags.
func (rm *resourceManager) syncTags(
	ctx context.Context,
	latest *resource,
	desired *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncTags")
	defer func() { exit(err) }()

	from, fromOrder := convertToOrderedACKTags(latest.ko.Spec.Tags)
	to, toOrder := convertToOrderedACKTags(desired.ko.Spec.Tags)
	added, _, removed := ackcompare.GetTagsDifference(from, to)
	toAdd := fromACKTags(added, toOrder)
	toDelete := fromACKTags(removed, fromOrder)

	if len(toAdd) == 0 && len(toDelete) == 0 {
		return nil
	}

	// Remove tags no longer in desired state
	if len(toDelete) > 0 {
		_, err = rm.sdkapi.RemoveTags(
			ctx,
			&svcsdk.RemoveTagsInput{
				ResourceId: (*string)(latest.ko.Status.ACKResourceMetadata.ARN),
				TagsList:   sdkTagsFromResourceTags(toDelete),
			},
		)
		rm.metrics.RecordAPICall("UPDATE", "RemoveTags", err)
		if err != nil {
			return err
		}
	}

	// Add/update tags in desired state
	if len(toAdd) > 0 {
		_, err = rm.sdkapi.AddTags(
			ctx,
			&svcsdk.AddTagsInput{
				ResourceId: (*string)(latest.ko.Status.ACKResourceMetadata.ARN),
				TagsList:   sdkTagsFromResourceTags(toAdd),
			},
		)
		rm.metrics.RecordAPICall("UPDATE", "AddTags", err)
		if err != nil {
			return err
		}
	}
	return nil
}

// getTags retrieves a resource list of tags.
func (rm *resourceManager) getTags(
	ctx context.Context,
	resourceARN string,
) ([]*svcapitypes.Tag, error) {
	tags := []*svcapitypes.Tag{}

	listTagsResponse, err := rm.sdkapi.ListTags(
		ctx,
		&svcsdk.ListTagsInput{
			ResourceIdList: []string{resourceARN},
		},
	)
	rm.metrics.RecordAPICall("GET", "ListTags", err)
	if err != nil {
		return nil, err
	}
	if len(listTagsResponse.ResourceTagList) != 1 {
		return nil, fmt.Errorf("expected to receive one exact resource tags in ListTags response")
	}
	for _, tag := range listTagsResponse.ResourceTagList[0].TagsList {
		tags = append(tags, &svcapitypes.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}
	return tags, nil
}

func customPreCompare(
	delta *ackcompare.Delta,
	a *resource,
	b *resource,
) {
	// a = desired, b = latest; GetTagsDifference(from=latest, to=desired)
	desiredTags, _ := convertToOrderedACKTags(a.ko.Spec.Tags)
	latestTags, _ := convertToOrderedACKTags(b.ko.Spec.Tags)
	addedTags, _, removedTags := ackcompare.GetTagsDifference(latestTags, desiredTags)
	if len(addedTags) > 0 || len(removedTags) > 0 {
		delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
	}
}

// sdkTagsFromResourceTags transforms a *svcapitypes.Tag array to a *svcsdk.Tag array.
func sdkTagsFromResourceTags(rTags []*svcapitypes.Tag) []svcsdktypes.Tag {
	tags := make([]svcsdktypes.Tag, len(rTags))
	for i := range rTags {
		tags[i] = svcsdktypes.Tag{
			Key:   rTags[i].Key,
			Value: rTags[i].Value,
		}
	}
	return tags
}
