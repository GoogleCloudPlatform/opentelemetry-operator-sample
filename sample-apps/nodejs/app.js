// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

const http = require('http');

var serviceURL = process.env.SERVICE_NAME

console.log("starting app for service at "+serviceURL)

setInterval(function() {
  http.get(serviceURL, resp => {
    let data = ''
    resp.on('data', d => {
      data += d
    })
    resp.on('end', () => {
      console.log(JSON.parse(data))
    })
  }).on('error', err => {
    console.log(err)
  });
}, 1000)
