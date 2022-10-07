# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
import time
import requests

service_url = os.environ['SERVICE_NAME']
print(f"starting app for service at {service_url}")


while True:
    time.sleep(1)
    try:
        json = requests.get(service_url, timeout=10).json()
        print(json)
    except requests.RequestException as e:
        print(e)
