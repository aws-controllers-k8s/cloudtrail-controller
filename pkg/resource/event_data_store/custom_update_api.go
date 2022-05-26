package event_data_store

import (
	"context"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go/service/cloudtrail"
)

// customUpdateEventDataStore implements a custom logic for handling EDS
// resource updates.
func (rm *resourceManager) customUpdateEventDataStore(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	var err error
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.customUpdateEventDataStore")
	defer func(err error) { exit(err) }(err)

	if delta.DifferentExcept("Spec.Tags") {
		err = rm.updateEventDataStoreFields(ctx, desired)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.Tags") {
		err = rm.syncTags(ctx, latest, desired)
		if err != nil {
			return nil, err
		}
	}
	readOneLatest, err := rm.ReadOne(ctx, desired)
	if err != nil {
		return nil, err
	}
	return rm.concreteResource(readOneLatest), nil
}

// syncTags updates a EDS list of tags.
func (rm *resourceManager) syncTags(
	ctx context.Context,
	latest *resource,
	desired *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncTags")
	defer func(err error) { exit(err) }(err)

	added, removed := computeTagsDelta(latest.ko.Spec.Tags, desired.ko.Spec.Tags)

	// Tags to create/update

	if len(removed) > 0 {
		_, err = rm.sdkapi.RemoveTagsWithContext(
			ctx,
			&svcsdk.RemoveTagsInput{
				ResourceId: (*string)(latest.ko.Status.ACKResourceMetadata.ARN),
				TagsList:   sdkTagsFromResourceTags(removed),
			},
		)
		rm.metrics.RecordAPICall("UPDATE", "RemoveTags", err)
		if err != nil {
			return err
		}
	}

	if len(added) > 0 {
		_, err = rm.sdkapi.AddTagsWithContext(
			ctx,
			&svcsdk.AddTagsInput{
				ResourceId: (*string)(latest.ko.Status.ACKResourceMetadata.ARN),
				TagsList:   sdkTagsFromResourceTags(added),
			},
		)
		rm.metrics.RecordAPICall("UPDATE", "AddTags", err)
		if err != nil {
			return err
		}
	}
	return nil
}

// updateEventDataStoreFields updates a given EventDataStore fields.
func (rm *resourceManager) updateEventDataStoreFields(
	ctx context.Context,
	desired *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.updateEventDataStoreFields")
	defer func(err error) { exit(err) }(err)
	input := &svcsdk.UpdateEventDataStoreInput{
		EventDataStore: (*string)(desired.ko.Status.ACKResourceMetadata.ARN),
	}

	if desired.ko.Spec.Name != nil {
		input.SetName(*desired.ko.Spec.Name)
	}
	if desired.ko.Spec.AdvancedEventSelectors != nil {
		input.SetAdvancedEventSelectors(advancedEventSelectorsFromResource(desired))
	}
	if desired.ko.Spec.MultiRegionEnabled != nil {
		input.SetMultiRegionEnabled(*desired.ko.Spec.MultiRegionEnabled)
	}
	if desired.ko.Spec.OrganizationEnabled != nil {
		input.SetOrganizationEnabled(*desired.ko.Spec.OrganizationEnabled)
	}
	if desired.ko.Spec.RetentionPeriod != nil {
		input.SetRetentionPeriod(*desired.ko.Spec.RetentionPeriod)
	}
	if desired.ko.Spec.TerminationProtectionEnabled != nil {
		input.SetTerminationProtectionEnabled(*desired.ko.Spec.TerminationProtectionEnabled)
	}

	_, err = rm.sdkapi.UpdateEventDataStoreWithContext(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateEventDataStore", err)
	return err
}

func advancedEventSelectorsFromResource(r *resource) []*svcsdk.AdvancedEventSelector {
	aes := make([]*svcsdk.AdvancedEventSelector, len(r.ko.Spec.AdvancedEventSelectors))
	for _, advancedEventSelector := range r.ko.Spec.AdvancedEventSelectors {
		afs := make([]*svcsdk.AdvancedFieldSelector, len(advancedEventSelector.FieldSelectors))
		for _, advancedFieldSelector := range advancedEventSelector.FieldSelectors {
			afs = append(afs, &svcsdk.AdvancedFieldSelector{
				Field:         advancedFieldSelector.Field,
				StartsWith:    advancedFieldSelector.StartsWith,
				NotStartsWith: advancedFieldSelector.NotStartsWith,
				EndsWith:      advancedFieldSelector.EndsWith,
				NotEndsWith:   advancedFieldSelector.NotEndsWith,
				Equals:        advancedFieldSelector.Equals,
				NotEquals:     advancedFieldSelector.NotEquals,
			})
		}

		aes = append(aes, &svcsdk.AdvancedEventSelector{
			Name:           advancedEventSelector.Name,
			FieldSelectors: afs,
		})
	}
	return aes
}
