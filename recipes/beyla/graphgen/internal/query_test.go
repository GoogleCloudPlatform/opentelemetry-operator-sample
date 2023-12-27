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
          "k8s_pod_ip": "10.112.0.0",
          "k8s_pod_name": "child1",
          "net_sock_peer_addr": "10.112.1.100"
        },
        "value": [1703024317.243, "0.09666750444003214"]
      },
      {
        "metric": {
          "k8s_pod_ip": "10.112.0.1",
          "k8s_pod_name": "child2",
          "net_sock_peer_addr": "10.112.1.100"
        },
        "value": [1703024317.243, "0.09666750444003214"]
      },
      {
        "metric": {
          "k8s_pod_ip": "10.112.0.2",
          "k8s_pod_name": "leaf",
          "net_sock_peer_addr": "10.112.0.1"
        },
        "value": [1703024317.243, "0.09666750444003214"]
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

	graph, err := QueryPrometheus(ctx, client, time.Second*5)
	assert.NoError(t, err)

	assert.Equal(
		t,
		graph.nodes,
		map[string]*Node{
			"10.112.0.0":   {Ip: "10.112.0.0", Name: "child1"},
			"10.112.0.1":   {Ip: "10.112.0.1", Name: "child2"},
			"10.112.0.2":   {Ip: "10.112.0.2", Name: "leaf"},
			"10.112.1.100": {Ip: "10.112.1.100"},
		},
	)
	assert.Equal(
		t,
		graph.adjacencies,
		map[string][]string{
			"10.112.1.100": {"10.112.0.0", "10.112.0.1"},
			"10.112.0.1":   {"10.112.0.2"},
		},
	)
}
