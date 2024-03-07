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
	"google.golang.org/api/option"
	apihttp "google.golang.org/api/transport/http"
)

var (
	projectId   = flag.String("projectId", "", "Required. The projectID to query.")
	cluster     = flag.String("cluster", "", "The GKE cluster to filter in the query. If not set, does no filtering.")
	namespace   = flag.String("namespace", "", "The k8s namespace to filter in the query. If not set, does no filtering.")
	queryWindow = flag.Duration("queryWindow", time.Minute*5, "Query window for service graph metrics.")
)

func main() {
	flag.Parse()
	if *projectId == "" {
		fmt.Println("You must set the -projectId flag!")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	client, err := createPromApiClient(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Got error creating prometheus api client")
		panic(err)
	}

	graph, err := internal.QueryPrometheus(
		ctx,
		client,
		internal.QueryArgs{
			QueryWindow: *queryWindow,
			Cluster:     *cluster,
			Namespace:   *namespace,
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Got error querying prometheus", "err", err)
		panic(err)
	}
	slog.InfoContext(ctx, "Got graph list", "graph", graph)

	err = graph.Render(os.Stdout)
	if err != nil {
		slog.ErrorContext(ctx, "Got error rendering graph", "err", err)
		panic(err)
	}
}

func createPromApiClient(ctx context.Context) (api.Client, error) {
	roundTripper, err := apihttp.NewTransport(
		ctx,
		http.DefaultTransport,
		option.WithScopes("https://www.googleapis.com/auth/monitoring.read"),
	)
	if err != nil {
		return nil, err
	}
	return api.NewClient(api.Config{
		Address:      fmt.Sprintf("https://monitoring.googleapis.com/v1/projects/%v/location/global/prometheus", *projectId),
		RoundTripper: roundTripper,
	})
}
