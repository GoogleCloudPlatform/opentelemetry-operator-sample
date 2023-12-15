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
	"io"
	"log/slog"
	"os"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type Node struct {
	Ip   string
	Name string
}

type Graph struct {
	// Adjacency list where nodes are IP address strings and values are the nodes. Each key called the value.
	adjacencies map[string][]string
	nodes       map[string]*Node
}

func NewGraph() *Graph {
	return &Graph{
		adjacencies: make(map[string][]string),
		nodes:       make(map[string]*Node),
	}
}

func (g Graph) AddEdge(client, server *Node) {
	g.upsertNode(client)
	g.upsertNode(server)
	g.adjacencies[client.Ip] = append(g.adjacencies[client.Ip], server.Ip)
}

func (g *Graph) upsertNode(node *Node) {
	_, ok := g.nodes[node.Ip]
	if !ok {
		g.nodes[node.Ip] = node
	}

	// only server nodes will have pod name assigned to name, so try and upsert it if there is
	// an existing key that was a client only
	if g.nodes[node.Ip].Name == "" {
		g.nodes[node.Ip].Name = node.Name
	}
}

func (g Graph) Render(writer io.Writer) error {
	gv := graphviz.New()
	graph, err := gv.Graph()
	if err != nil {
		return err
	}
	defer func() {
		if err := graph.Close(); err != nil {
			slog.Error("error closing graph", "err", err)
		}
		gv.Close()
	}()

	nodes := map[*Node]*cgraph.Node{}

	// add nodes and edges
	for clientIp, servers := range g.adjacencies {
		clientNode, err := getCNode(g.nodes[clientIp], graph, nodes)
		if err != nil {
			return err
		}
		for _, serverIp := range servers {
			serverNode, err := getCNode(g.nodes[serverIp], graph, nodes)
			if err != nil {
				return err
			}

			_, err = graph.CreateEdge("", clientNode, serverNode)
			// TODO: set weight in label
			// e.SetLabel("e")
			if err != nil {
				return err
			}
		}
	}

	// print to stdout
	if err := gv.Render(graph, graphviz.SVG, os.Stdout); err != nil {
		return err
	}
	return nil
}

func getCNode(node *Node, graph *cgraph.Graph, nodes map[*Node]*cgraph.Node) (*cgraph.Node, error) {
	if node, ok := nodes[node]; ok {
		return node, nil
	}

	cnode, err := graph.CreateNode(node.Ip)
	if err != nil {
		return nil, err
	}
	if node.Name != "" {
		cnode.SetLabel(node.Name)
	}
	cnode.SetOrdering(cgraph.OutOrdering)

	nodes[node] = cnode
	return cnode, nil
}
