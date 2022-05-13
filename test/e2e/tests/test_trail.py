# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Integration tests for the CloudTrail Trail resource
"""

import boto3
import logging
import time
from typing import Dict

import pytest

from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_cloudtrail_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.tests.helper import CloudTrailValidator
from e2e.common import TRAIL_NAME

RESOURCE_PLURAL = "trails"

CREATE_WAIT_AFTER_SECONDS = 10
MODIFY_WAIT_AFTER_SECONDS = 10
DELETE_WAIT_AFTER_SECONDS = 10

@pytest.fixture
def simple_trail():
    replacements = REPLACEMENT_VALUES.copy()
    replacements["TRAIL_NAME"] = TRAIL_NAME
    replacements["BUCKET_NAME"] = get_bootstrap_resources().TrailLogBucket.name

    resource_data = load_cloudtrail_resource(
        "trail",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    # Create the k8s resource
    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
        TRAIL_NAME, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)

    time.sleep(CREATE_WAIT_AFTER_SECONDS)

    # Get latest trail CR
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr)

    # Try to delete, if doesn't already exist
    try:
        _, deleted = k8s.delete_custom_resource(ref, 3, 10)
        assert deleted
    except:
        pass

@service_marker
@pytest.mark.canary
class TestTrail:
    def test_crud(self, cloudtrail_client, simple_trail):
        (ref, cr) = simple_trail

        trail_name = cr["spec"]["name"]
        trail_arn = cr["status"]["ackResourceMetadata"]["arn"]

        cloudtrail_validator = CloudTrailValidator(cloudtrail_client)
        # verify that trail exists
        cloudtrail_validator.assert_trail(trail_name)
        # verify that trail tags are created
        cloudtrail_validator.assert_trail_tags(trail_arn, [{"Key":"env","Value":"testing"}])

        # Update the trail log file validation
        cr["spec"]["includeGlobalServiceEvents"] = True
        k8s.patch_custom_resource(ref, cr)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)

        trail = cloudtrail_validator.get_trail(trail_name)
        assert trail["IncludeGlobalServiceEvents"]

        # Delete k8s resource
        _, deleted = k8s.delete_custom_resource(ref)
        assert deleted is True

        time.sleep(DELETE_WAIT_AFTER_SECONDS)

        # Check trail doesn't exist
        cloudtrail_validator.assert_trail(trail_name, exists=False)