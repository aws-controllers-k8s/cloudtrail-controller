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

package trail

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	smithy "github.com/aws/smithy-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/cloudtrail-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &svcsdk.Client{}
	_ = &svcapitypes.Trail{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
	_ = &ackrequeue.NoRequeue{}
	_ = &aws.Config{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.GetTrailOutput
	resp, err = rm.sdkapi.GetTrail(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "GetTrail", err)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "TrailNotFoundException" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if resp.Trail.CloudWatchLogsLogGroupArn != nil {
		ko.Spec.CloudWatchLogsLogGroupARN = resp.Trail.CloudWatchLogsLogGroupArn
	} else {
		ko.Spec.CloudWatchLogsLogGroupARN = nil
	}
	if resp.Trail.CloudWatchLogsRoleArn != nil {
		ko.Spec.CloudWatchLogsRoleARN = resp.Trail.CloudWatchLogsRoleArn
	} else {
		ko.Spec.CloudWatchLogsRoleARN = nil
	}
	if resp.Trail.IncludeGlobalServiceEvents != nil {
		ko.Spec.IncludeGlobalServiceEvents = resp.Trail.IncludeGlobalServiceEvents
	} else {
		ko.Spec.IncludeGlobalServiceEvents = nil
	}
	if resp.Trail.IsMultiRegionTrail != nil {
		ko.Spec.IsMultiRegionTrail = resp.Trail.IsMultiRegionTrail
	} else {
		ko.Spec.IsMultiRegionTrail = nil
	}
	if resp.Trail.IsOrganizationTrail != nil {
		ko.Spec.IsOrganizationTrail = resp.Trail.IsOrganizationTrail
	} else {
		ko.Spec.IsOrganizationTrail = nil
	}
	if resp.Trail.KmsKeyId != nil {
		ko.Spec.KMSKeyID = resp.Trail.KmsKeyId
	} else {
		ko.Spec.KMSKeyID = nil
	}
	if resp.Trail.LogFileValidationEnabled != nil {
		ko.Status.LogFileValidationEnabled = resp.Trail.LogFileValidationEnabled
	} else {
		ko.Status.LogFileValidationEnabled = nil
	}
	if resp.Trail.Name != nil {
		ko.Spec.Name = resp.Trail.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.Trail.S3BucketName != nil {
		ko.Spec.S3BucketName = resp.Trail.S3BucketName
	} else {
		ko.Spec.S3BucketName = nil
	}
	if resp.Trail.S3KeyPrefix != nil {
		ko.Spec.S3KeyPrefix = resp.Trail.S3KeyPrefix
	} else {
		ko.Spec.S3KeyPrefix = nil
	}
	if resp.Trail.SnsTopicARN != nil {
		ko.Status.SNSTopicARN = resp.Trail.SnsTopicARN
	} else {
		ko.Status.SNSTopicARN = nil
	}
	if resp.Trail.SnsTopicName != nil {
		ko.Spec.SNSTopicName = resp.Trail.SnsTopicName
	} else {
		ko.Spec.SNSTopicName = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.Trail.TrailARN != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.Trail.TrailARN)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}

	rm.setStatusDefaults(ko)
	if err := rm.setResourceAdditionalFields(ctx, ko); err != nil {
		return nil, err
	}
	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return r.ko.Spec.Name == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.GetTrailInput, error) {
	res := &svcsdk.GetTrailInput{}

	if r.ko.Spec.Name != nil {
		res.Name = r.ko.Spec.Name
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer func() {
		exit(err)
	}()
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateTrailOutput
	_ = resp
	resp, err = rm.sdkapi.CreateTrail(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateTrail", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.CloudWatchLogsLogGroupArn != nil {
		ko.Spec.CloudWatchLogsLogGroupARN = resp.CloudWatchLogsLogGroupArn
	} else {
		ko.Spec.CloudWatchLogsLogGroupARN = nil
	}
	if resp.CloudWatchLogsRoleArn != nil {
		ko.Spec.CloudWatchLogsRoleARN = resp.CloudWatchLogsRoleArn
	} else {
		ko.Spec.CloudWatchLogsRoleARN = nil
	}
	if resp.IncludeGlobalServiceEvents != nil {
		ko.Spec.IncludeGlobalServiceEvents = resp.IncludeGlobalServiceEvents
	} else {
		ko.Spec.IncludeGlobalServiceEvents = nil
	}
	if resp.IsMultiRegionTrail != nil {
		ko.Spec.IsMultiRegionTrail = resp.IsMultiRegionTrail
	} else {
		ko.Spec.IsMultiRegionTrail = nil
	}
	if resp.IsOrganizationTrail != nil {
		ko.Spec.IsOrganizationTrail = resp.IsOrganizationTrail
	} else {
		ko.Spec.IsOrganizationTrail = nil
	}
	if resp.KmsKeyId != nil {
		ko.Spec.KMSKeyID = resp.KmsKeyId
	} else {
		ko.Spec.KMSKeyID = nil
	}
	if resp.LogFileValidationEnabled != nil {
		ko.Status.LogFileValidationEnabled = resp.LogFileValidationEnabled
	} else {
		ko.Status.LogFileValidationEnabled = nil
	}
	if resp.Name != nil {
		ko.Spec.Name = resp.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.S3BucketName != nil {
		ko.Spec.S3BucketName = resp.S3BucketName
	} else {
		ko.Spec.S3BucketName = nil
	}
	if resp.S3KeyPrefix != nil {
		ko.Spec.S3KeyPrefix = resp.S3KeyPrefix
	} else {
		ko.Spec.S3KeyPrefix = nil
	}
	if resp.SnsTopicARN != nil {
		ko.Status.SNSTopicARN = resp.SnsTopicARN
	} else {
		ko.Status.SNSTopicARN = nil
	}
	if resp.SnsTopicName != nil {
		ko.Spec.SNSTopicName = resp.SnsTopicName
	} else {
		ko.Spec.SNSTopicName = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.TrailARN != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.TrailARN)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateTrailInput, error) {
	res := &svcsdk.CreateTrailInput{}

	if r.ko.Spec.CloudWatchLogsLogGroupARN != nil {
		res.CloudWatchLogsLogGroupArn = r.ko.Spec.CloudWatchLogsLogGroupARN
	}
	if r.ko.Spec.CloudWatchLogsRoleARN != nil {
		res.CloudWatchLogsRoleArn = r.ko.Spec.CloudWatchLogsRoleARN
	}
	if r.ko.Spec.EnableLogFileValidation != nil {
		res.EnableLogFileValidation = r.ko.Spec.EnableLogFileValidation
	}
	if r.ko.Spec.IncludeGlobalServiceEvents != nil {
		res.IncludeGlobalServiceEvents = r.ko.Spec.IncludeGlobalServiceEvents
	}
	if r.ko.Spec.IsMultiRegionTrail != nil {
		res.IsMultiRegionTrail = r.ko.Spec.IsMultiRegionTrail
	}
	if r.ko.Spec.IsOrganizationTrail != nil {
		res.IsOrganizationTrail = r.ko.Spec.IsOrganizationTrail
	}
	if r.ko.Spec.KMSKeyID != nil {
		res.KmsKeyId = r.ko.Spec.KMSKeyID
	}
	if r.ko.Spec.Name != nil {
		res.Name = r.ko.Spec.Name
	}
	if r.ko.Spec.S3BucketName != nil {
		res.S3BucketName = r.ko.Spec.S3BucketName
	}
	if r.ko.Spec.S3KeyPrefix != nil {
		res.S3KeyPrefix = r.ko.Spec.S3KeyPrefix
	}
	if r.ko.Spec.SNSTopicName != nil {
		res.SnsTopicName = r.ko.Spec.SNSTopicName
	}
	if r.ko.Spec.Tags != nil {
		f11 := []svcsdktypes.Tag{}
		for _, f11iter := range r.ko.Spec.Tags {
			f11elem := &svcsdktypes.Tag{}
			if f11iter.Key != nil {
				f11elem.Key = f11iter.Key
			}
			if f11iter.Value != nil {
				f11elem.Value = f11iter.Value
			}
			f11 = append(f11, *f11elem)
		}
		res.TagsList = f11
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	return rm.customUpdateTrail(ctx, desired, latest, delta)
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteTrailOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteTrail(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteTrail", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteTrailInput, error) {
	res := &svcsdk.DeleteTrailInput{}

	if r.ko.Spec.Name != nil {
		res.Name = r.ko.Spec.Name
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.Trail,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	if err == nil {
		return false
	}

	var terminalErr smithy.APIError
	if !errors.As(err, &terminalErr) {
		return false
	}
	switch terminalErr.ErrorCode() {
	case "InvalidParameterCombination",
		"InvalidParameterValue",
		"InvalidQueryParameter",
		"MissingAction",
		"MissingParameter":
		return true
	default:
		return false
	}
}
