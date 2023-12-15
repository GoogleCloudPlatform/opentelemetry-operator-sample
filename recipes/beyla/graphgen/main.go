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

package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/opentelemetry-operator-sample/recipes/beyla/graphgen/internal"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"google.golang.org/api/option"
	apihttp "google.golang.org/api/transport/http"
)

var (
	projectId   = flag.String("projectId", "", "The projectID to query")
	queryWindow = flag.Duration("queryWindow", time.Minute*5, "Query window for service graph metrics")
)

const (
	// metrics are based on server side, so the peer is the caller (client) and pod is the
	// callee (server)
	clientIpKey  = "net_sock_peer_addr"
	serverIpKey  = "k8s_pod_ip"
	serverPodKey = "k8s_pod_name"
)

func main() {
	flag.Parse()
	if *projectId == "" {
		fmt.Println("You must set the -projectId flag!")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	graph, err := queryPrometheus(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Got error querying prometheus", "err", err)
		panic(err)
	}
	outputGraph(graph)

	slog.InfoContext(ctx, "Got graph list", "graph", graph)
}

func getQuery() string {
	return fmt.Sprintf(
		`sum by (
			k8s_pod_ip, net_sock_peer_addr, k8s_pod_name
		) (
			rate(http_servicegraph_calls_total[%s])
		)
		`,
		*queryWindow,
	)
}

func queryPrometheus(ctx context.Context) (*internal.Graph, error) {
	roundTripper, err := apihttp.NewTransport(
		ctx,
		http.DefaultTransport,
		option.WithScopes("https://www.googleapis.com/auth/monitoring.read"),
	)
	if err != nil {
		return nil, err
	}
	client, err := api.NewClient(api.Config{
		Address:      fmt.Sprintf("https://monitoring.googleapis.com/v1/projects/%v/location/global/prometheus", *projectId),
		RoundTripper: roundTripper,
	})
	if err != nil {
		return nil, err
	}

	promApi := v1.NewAPI(client)
	res, warnings, err := promApi.Query(ctx, getQuery(), time.Now())
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

	graph := internal.NewGraph()
	for _, sample := range vec {
		labels := sample.Metric
		client := &internal.Node{Ip: string(labels[clientIpKey])}
		server := &internal.Node{Ip: string(labels[serverIpKey]), Name: string(labels[serverPodKey])}
		graph.AddEdge(client, server)
	}
	return graph, nil
}

func outputGraph(g *internal.Graph) error {
	return g.Render(os.Stdout)
}
