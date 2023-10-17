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

var http = require("http");

console.log("serving on port 8080...")

var serviceURL = process.env.NEXT_SERVER

var listener = function (req, res) {

  if (serviceURL) {
    http.get(serviceURL, resp => {
      let data = ''
      resp.on('data', d => { data += d });
      resp.on('end', () => { 
        res.writeHead(200);
        res.end(data);
      });
    });
  } else {
    res.writeHead(200);
    res.end(JSON.stringify({data: "Hello world"}));
  }
};
var server = http.createServer(listener);

server.listen(8080);
