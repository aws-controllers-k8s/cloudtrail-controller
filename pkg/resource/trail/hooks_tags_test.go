package trail

import (
	"context"
	"errors"
	"reflect"
	"testing"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackmetrics "github.com/aws-controllers-k8s/runtime/pkg/metrics"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/smithy-go/middleware"

	svcapitypes "github.com/aws-controllers-k8s/cloudtrail-controller/apis/v1alpha1"
)

func newTrailWithTags(tags ...*svcapitypes.Tag) *resource {
	return &resource{ko: &svcapitypes.Trail{Spec: svcapitypes.TrailSpec{Tags: tags}}}
}

func tag(key, value string) *svcapitypes.Tag {
	return &svcapitypes.Tag{Key: &key, Value: &value}
}

func TestCustomPreCompareTags(t *testing.T) {
	tests := []struct {
		name        string
		a           *resource
		b           *resource
		wantDiffers bool
	}{
		{
			name:        "no-op when tags are identical",
			a:           newTrailWithTags(tag("env", "prod"), tag("team", "core")),
			b:           newTrailWithTags(tag("env", "prod"), tag("team", "core")),
			wantDiffers: false,
		},
		{
			name:        "no-op when tags are identical but reordered",
			a:           newTrailWithTags(tag("env", "prod"), tag("team", "core")),
			b:           newTrailWithTags(tag("team", "core"), tag("env", "prod")),
			wantDiffers: false,
		},
		{
			name:        "add a tag",
			a:           newTrailWithTags(tag("env", "prod")),
			b:           newTrailWithTags(tag("env", "prod"), tag("team", "core")),
			wantDiffers: true,
		},
		{
			name:        "remove a tag",
			a:           newTrailWithTags(tag("env", "prod"), tag("team", "core")),
			b:           newTrailWithTags(tag("env", "prod")),
			wantDiffers: true,
		},
		{
			name:        "value change is a difference",
			a:           newTrailWithTags(tag("env", "prod")),
			b:           newTrailWithTags(tag("env", "dev")),
			wantDiffers: true,
		},
		{
			name:        "nil key does not panic and is ignored",
			a:           newTrailWithTags(&svcapitypes.Tag{Value: strPtr("orphan")}, tag("env", "prod")),
			b:           newTrailWithTags(tag("env", "prod")),
			wantDiffers: false,
		},
		{
			name:        "nil key on both sides does not panic",
			a:           newTrailWithTags(&svcapitypes.Tag{}),
			b:           newTrailWithTags(&svcapitypes.Tag{}),
			wantDiffers: false,
		},
		{
			name:        "add and remove in the same update",
			a:           newTrailWithTags(tag("env", "prod"), tag("team", "core")),
			b:           newTrailWithTags(tag("env", "prod"), tag("region", "us-east-1")),
			wantDiffers: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			delta := ackcompare.NewDelta()
			customPreCompare(delta, tc.a, tc.b)
			if got := delta.DifferentAt("Spec.Tags"); got != tc.wantDiffers {
				t.Errorf("DifferentAt(Spec.Tags) = %v, want %v", got, tc.wantDiffers)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}

type recordedCall struct {
	op   string
	tags map[string]string
}

func newFakeTagClient(calls *[]recordedCall, failOn string) *svcsdk.Client {
	return svcsdk.New(svcsdk.Options{
		Region: "us-east-1",
		APIOptions: []func(*middleware.Stack) error{
			func(s *middleware.Stack) error {
				return s.Initialize.Add(middleware.InitializeMiddlewareFunc(
					"fakeTagAPI",
					func(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (middleware.InitializeOutput, middleware.Metadata, error) {
						var op string
						var tagMap map[string]string
						var out interface{}
						switch p := in.Parameters.(type) {
						case *svcsdk.RemoveTagsInput:
							op = "RemoveTags"
							tagMap = make(map[string]string)
							for _, t := range p.TagsList {
								if t.Key != nil {
									v := ""
									if t.Value != nil {
										v = *t.Value
									}
									tagMap[*t.Key] = v
								}
							}
							out = &svcsdk.RemoveTagsOutput{}
						case *svcsdk.AddTagsInput:
							op = "AddTags"
							tagMap = make(map[string]string)
							for _, t := range p.TagsList {
								if t.Key != nil {
									v := ""
									if t.Value != nil {
										v = *t.Value
									}
									tagMap[*t.Key] = v
								}
							}
							out = &svcsdk.AddTagsOutput{}
						default:
							return middleware.InitializeOutput{}, middleware.Metadata{}, errors.New("unexpected operation")
						}
						*calls = append(*calls, recordedCall{op: op, tags: tagMap})
						if op == failOn {
							return middleware.InitializeOutput{}, middleware.Metadata{}, errors.New("injected failure")
						}
						return middleware.InitializeOutput{Result: out}, middleware.Metadata{}, nil
					},
				), middleware.Before)
			},
		},
	})
}

func newTestRM(calls *[]recordedCall, failOn string) *resourceManager {
	return &resourceManager{
		sdkapi:  newFakeTagClient(calls, failOn),
		metrics: ackmetrics.NewMetrics("cloudtrail"),
	}
}

func withARN(r *resource) *resource {
	arn := ackv1alpha1.AWSResourceName("arn:aws:cloudtrail:us-east-1:123456789012:trail/test")
	if r.ko.Status.ACKResourceMetadata == nil {
		r.ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	r.ko.Status.ACKResourceMetadata.ARN = &arn
	return r
}

func TestSyncTrailTags(t *testing.T) {
	cases := []struct {
		name      string
		latest    []*svcapitypes.Tag
		desired   []*svcapitypes.Tag
		failOn    string
		wantCalls []recordedCall
		wantErr   bool
	}{
		{
			name:    "no-op when identical",
			latest:  []*svcapitypes.Tag{tag("env", "prod"), tag("team", "core")},
			desired: []*svcapitypes.Tag{tag("env", "prod"), tag("team", "core")},
		},
		{
			name:    "no-op when reordered",
			latest:  []*svcapitypes.Tag{tag("env", "prod"), tag("team", "core")},
			desired: []*svcapitypes.Tag{tag("team", "core"), tag("env", "prod")},
		},
		{
			name:      "add only",
			latest:    []*svcapitypes.Tag{tag("env", "prod")},
			desired:   []*svcapitypes.Tag{tag("env", "prod"), tag("team", "core")},
			wantCalls: []recordedCall{{op: "AddTags", tags: map[string]string{"team": "core"}}},
		},
		{
			name:      "remove only",
			latest:    []*svcapitypes.Tag{tag("env", "prod"), tag("team", "core")},
			desired:   []*svcapitypes.Tag{tag("env", "prod")},
			wantCalls: []recordedCall{{op: "RemoveTags", tags: map[string]string{"team": "core"}}},
		},
		{
			name:      "value change is add-only",
			latest:    []*svcapitypes.Tag{tag("env", "prod")},
			desired:   []*svcapitypes.Tag{tag("env", "dev")},
			wantCalls: []recordedCall{{op: "AddTags", tags: map[string]string{"env": "dev"}}},
		},
		{
			name:    "remove before add",
			latest:  []*svcapitypes.Tag{tag("env", "prod"), tag("team", "core")},
			desired: []*svcapitypes.Tag{tag("env", "prod"), tag("region", "us-east-1")},
			wantCalls: []recordedCall{
				{op: "RemoveTags", tags: map[string]string{"team": "core"}},
				{op: "AddTags", tags: map[string]string{"region": "us-east-1"}},
			},
		},
		{
			name:    "RemoveTags error stops AddTags",
			latest:  []*svcapitypes.Tag{tag("env", "prod"), tag("team", "core")},
			desired: []*svcapitypes.Tag{tag("env", "prod"), tag("region", "us-east-1")},
			failOn:  "RemoveTags",
			wantCalls: []recordedCall{
				{op: "RemoveTags", tags: map[string]string{"team": "core"}},
			},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var calls []recordedCall
			rm := newTestRM(&calls, tc.failOn)
			latest := withARN(newTrailWithTags(tc.latest...))
			desired := newTrailWithTags(tc.desired...)
			err := rm.syncTrailTags(context.Background(), latest, desired)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tc.wantCalls) == 0 {
				if len(calls) != 0 {
					t.Errorf("want no API calls, got %v", calls)
				}
				return
			}
			if !reflect.DeepEqual(calls, tc.wantCalls) {
				t.Errorf("calls = %v, want %v", calls, tc.wantCalls)
			}
		})
	}
}
