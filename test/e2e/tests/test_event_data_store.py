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

"""Integration tests for the CloudTrail EventDataStore resource
"""

import resource
import boto3
import logging
import time
from typing import Dict

import pytest

from acktest import tags
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_cloudtrail_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.tests.helper import CloudTrailValidator

RESOURCE_PLURAL = "eventdatastores"

CREATE_WAIT_AFTER_SECONDS = 10
MODIFY_WAIT_AFTER_SECONDS = 10
DELETE_WAIT_AFTER_SECONDS = 10

@pytest.fixture
def simple_event_data_store():
    resource_name = random_suffix_name("cloudtrail-eds", 24)


    replacements = REPLACEMENT_VALUES.copy()
    replacements["EVENT_DATA_STORE_NAME"] = resource_name

    resource_data = load_cloudtrail_resource(
        "event_data_store",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    # Create the k8s resource
    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
        resource_name, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)

    time.sleep(CREATE_WAIT_AFTER_SECONDS)

    # Get latest event data store CR
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
class TestEventDataStore:
    def test_crud(self, cloudtrail_client, simple_event_data_store):
        (ref, cr) = simple_event_data_store

        eds_name = cr["spec"]["name"]
        eds_arn = cr["status"]["ackResourceMetadata"]["arn"]

        cloudtrail_validator = CloudTrailValidator(cloudtrail_client)
        # verify that data event store exists
        assert cloudtrail_validator.get_event_data_store(eds_arn) is not None

        # verify that the event data store tags are created
        eds_tags = cloudtrail_validator.list_resource_tags(eds_arn)
        tags.assert_ack_system_tags(
            tags=eds_tags,
        )
        tags_dict = tags.to_dict(
            cr["spec"]["tags"],
            key_member_name = 'key',
            value_member_name = 'value'
        )
        tags.assert_equal_without_ack_tags(
            actual = tags_dict,
            expected = eds_tags,
        )
        resource_new_name = f"{eds_name}-renamed"


        cr = k8s.wait_resource_consumed_by_controller(ref)
        # Update the EDS log file validation
        cr["spec"]["name"] = resource_new_name
        k8s.patch_custom_resource(ref, cr)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)

        eds = cloudtrail_validator.get_event_data_store(eds_arn)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)

        assert eds["Name"] == resource_new_name

        # Delete k8s resource
        _, deleted = k8s.delete_custom_resource(ref)
        assert deleted is True

        time.sleep(DELETE_WAIT_AFTER_SECONDS)

        eds = cloudtrail_validator.get_event_data_store(eds_arn)
        # EventDataStore can't be instantly deleted. They stay on "PENDING DELETION"
        # state for 7 days before being deleted.
        assert eds["Status"] == "PENDING_DELETION"