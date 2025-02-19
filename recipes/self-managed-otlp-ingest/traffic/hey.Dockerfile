# Copyright 2025 Google LLC
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

FROM golang:1.21-alpine3.19 AS build

RUN apk add --no-cache wget
# The link for the binary is specified in library's README
# https://github.com/rakyll/hey/blob/master/README.md
RUN wget -O /hey https://hey-release.s3.us-east-2.amazonaws.com/hey_linux_amd64
RUN chmod +x /hey

FROM scratch
COPY --from=build /hey /hey

ENTRYPOINT ["/hey"]
