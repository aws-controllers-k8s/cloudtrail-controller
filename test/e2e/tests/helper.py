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

    def get_event_data_store(self, event_data_store_arn: str) -> dict:
        try:
            resp = self.cloudtrail_client.get_event_data_store(
                EventDataStore=event_data_store_arn,
            )
        except Exception as e:
            logging.debug(e)
            return None

        return resp


    def list_resource_tags(self, resource_arn: str) -> list:
        try:
            resp = self.cloudtrail_client.list_tags(
                ResourceIdList=[resource_arn],
            )
        except Exception as e:
            logging.debug(e)
            return None

        return resp["ResourceTagList"][0]["TagsList"]