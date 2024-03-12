// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/stretchr/testify/assert"
)

const testJson = `{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": [
      {
          "metric": {
          "client": "root",
          "client_address": "10.4.2.54",
          "server": "child1",
          "server_address": "10.4.1.85"
          },
          "value": [1709790055.51, "9.7341795580096"]
      },
      {
          "metric": {
          "client": "root",
          "client_address": "10.4.2.54",
          "server": "leaf1",
          "server_address": "10.4.1.86"
          },
          "value": [1709790055.51, "9.7341795580096"]
      },
      {
          "metric": {
          "client": "child1",
          "client_address": "10.4.1.85",
          "server": "leaf2",
          "server_address": "10.112.0.1"
          },
          "value": [1709790055.51, "9.7341795580096"]
      }
    ]
  }
}`

func TestQueryPrometheus(t *testing.T) {
	ctx := context.Background()
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, testJson)
		}),
	)
	defer ts.Close()

	client, err := api.NewClient(api.Config{
		Address: ts.URL,
		Client:  ts.Client(),
	})
	assert.NoError(t, err)

	graph, err := QueryPrometheus(
		ctx,
		client,
		QueryArgs{
			QueryWindow: time.Second * 5,
			Cluster:     "cluster",
		},
	)
	assert.NoError(t, err)

	assert.Equal(
		t,
		graph.nodes,
		map[string]*Node{
			"10.4.2.54":  {Ip: "10.4.2.54", Name: "root"},
			"10.4.1.85":  {Ip: "10.4.1.85", Name: "child1"},
			"10.4.1.86":  {Ip: "10.4.1.86", Name: "leaf1"},
			"10.112.0.1": {Ip: "10.112.0.1", Name: "leaf2"},
		},
	)
	assert.Equal(
		t,
		graph.adjacencies,
		map[string][]string{
			"10.4.1.85": {"10.112.0.1"},
			"10.4.2.54": {"10.4.1.85", "10.4.1.86"},
		},
	)
}
