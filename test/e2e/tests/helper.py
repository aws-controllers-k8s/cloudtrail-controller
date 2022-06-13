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

"""Helper functions for CloudTrail tests
"""

from typing import Union, Dict
import logging
from e2e import tag

class CloudTrailValidator:
    def __init__(self, cloudtrail_client):
        self.cloudtrail_client = cloudtrail_client

    def get_trail(self, trail_name: str) -> dict:
        try:
            resp = self.cloudtrail_client.get_trail(
                Name=trail_name,
            )
        except Exception as e:
            logging.debug(e)
            return None

        return resp["Trail"]

    def list_trail_tags(self, trail_arn: str) -> list:
        try:
            resp = self.cloudtrail_client.list_tags(
                ResourceIdList=[trail_arn],
            )
        except Exception as e:
            logging.debug(e)
            return None

        return resp["ResourceTagList"][0]["TagsList"]

    def assert_trail(self, trail_name: str, exists=True):
        assert (self.get_trail(trail_name) is not None) == exists

    def assert_trail_tags(self, trail_arn: str, tags:list):
        trail_tags = tag.cleaned(self.list_trail_tags(trail_arn))
        assert len(trail_tags) == len(tags)
        for i in range(0, len(tags)):
            assert tags[i]["Key"] == trail_tags[i]["Key"]
            assert tags[i]["Value"] == trail_tags[i]["Value"]