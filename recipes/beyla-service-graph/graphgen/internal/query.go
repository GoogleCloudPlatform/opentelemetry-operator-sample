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
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	clientIpKey = "client_address"
	clientKey   = "client"
	serverIpKey = "server_address"
	serverKey   = "server"
)

type QueryArgs struct {
	QueryWindow time.Duration
	Cluster     string
	Namespace   string
}

func QueryPrometheus(
	ctx context.Context,
	client api.Client,
	queryArgs QueryArgs,
) (*Graph, error) {
	promApi := v1.NewAPI(client)
	query := getQuery(queryArgs)
	slog.InfoContext(ctx, "logging", "query", query)
	res, warnings, err := promApi.Query(ctx, query, time.Now())
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
		client := &Node{Ip: string(labels[clientIpKey]), Name: string(labels[clientKey])}
		server := &Node{Ip: string(labels[serverIpKey]), Name: string(labels[serverKey])}
		graph.AddEdge(client, server)
	}
	return graph, nil
}

func getQuery(queryArgs QueryArgs) string {
	filters := addFilter(nil, "cluster", queryArgs.Cluster)
	filters = addFilter(filters, "namespace", queryArgs.Namespace)
	return fmt.Sprintf(
		`sum by (
			server, client_address, client, server_address
		) (
			rate(traces_service_graph_request_total{%s}[%s])
		)`,
		strings.Join(filters, ","),
		queryArgs.QueryWindow,
	)
}

func addFilter(filters []string, key, value string) []string {
	if value == "" {
		return filters
	}
	return append(filters, fmt.Sprintf(`%s="%s"`, key, value))
}
