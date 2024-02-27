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
	"fmt"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	// metrics are based on server side, so the peer is the caller (client) and pod is the
	// callee (server)
	clientIpKey  = "client_address"
	serverIpKey  = "k8s_pod_ip"
	serverPodKey = "k8s_pod_name"
)

func QueryPrometheus(ctx context.Context, client api.Client, queryWindow time.Duration) (*Graph, error) {
	promApi := v1.NewAPI(client)
	res, warnings, err := promApi.Query(ctx, getQuery(queryWindow), time.Now())
	if err != nil {
		return nil, err
	}

	if len(warnings) != 0 {
		slog.WarnContext(ctx, "Warnings from promQL query", "warnings", warnings)
	}
	slog.InfoContext(ctx, "Got metrics", "metrics", res)
	slog.InfoContext(ctx, "type", "type", res.Type())

	vec, ok := res.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("couldn't cast %v to vector", res)
	}

	graph := NewGraph()
	for _, sample := range vec {
		labels := sample.Metric
		client := &Node{Ip: string(labels[clientIpKey])}
		server := &Node{Ip: string(labels[serverIpKey]), Name: string(labels[serverPodKey])}
		graph.AddEdge(client, server)
	}
	return graph, nil
}

func getQuery(queryWindow time.Duration) string {
	return fmt.Sprintf(
		`sum by (
			k8s_pod_ip, client_address, k8s_pod_name
		) (
			rate(http_servicegraph_calls_total[%s])
		)`,
		queryWindow,
	)
}
