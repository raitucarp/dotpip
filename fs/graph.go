package fs

import (
	"encoding/json"
	"errors"

	"dotpip"

	"gonum.org/v1/gonum/graph/simple"
)

type Graph struct {
	Nodes []*GraphNode `json:"nodes"`
	Edges []*GraphEdge `json:"edges"`
}

type GraphNode struct {
	ID         int            `json:"id"`
	Labels     []string       `json:"labels"`
	Properties map[string]any `json:"properties"`
}

type GraphEdge struct {
	ID         int            `json:"id"`
	Type       string         `json:"type"`
	SourceNode int            `json:"source"`
	TargetNode int            `json:"target"`
	Properties map[string]any `json:"properties"`
}

func (g *Graph) buildGonumGraph() *simple.DirectedGraph {
	dg := simple.NewDirectedGraph()
	for _, n := range g.Nodes {
		dg.AddNode(simple.Node(n.ID))
	}
	for _, e := range g.Edges {
		dg.SetEdge(simple.Edge{F: simple.Node(e.SourceNode), T: simple.Node(e.TargetNode)})
	}
	return dg
}

// Very basic traversal helper to match chains length
func evaluateMatches(g *Graph, chainLength int) int {
	if chainLength == 0 {
		return len(g.Nodes)
	}

	dg := g.buildGonumGraph()
	pathsCount := 0

	for _, n := range g.Nodes {
		// Just a simple simulation of traversing a specific depth
		pathsCount += countPaths(dg, int64(n.ID), chainLength)
	}

	return pathsCount
}

func countPaths(dg *simple.DirectedGraph, nodeID int64, depth int) int {
	if depth == 0 {
		return 1
	}
	count := 0
	nodes := dg.From(nodeID)
	for nodes.Next() {
		count += countPaths(dg, nodes.Node().ID(), depth-1)
	}
	return count
}

func (f *FileSystem) GraphDelete(key dotpip.Key) (int, error) {
	return f.Del(key), nil
}

func (f *FileSystem) GraphExplain(_ dotpip.Key, query string) ([]string, error) {
	q, err := dotpip.CypherParser.ParseString("", query)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, clause := range q.Clauses {
		switch {
		case clause.Create != nil:
			result = append(result, string(dotpip.GraphKeywordCreate))
		case clause.Match != nil:
			result = append(result, string(dotpip.GraphKeywordMatch))
		case clause.Return != nil:
			result = append(result, string(dotpip.GraphKeywordReturn))
		case clause.Delete != nil:
			result = append(result, string(dotpip.GraphKeywordDelete))
		case clause.Set != nil:
			result = append(result, string(dotpip.GraphKeywordSet))
		}
	}
	return result, nil
}

func (f *FileSystem) GraphList() ([]string, error) {
	keys, err := f.Keys("*")
	if err != nil {
		return nil, err
	}

	var graphKeys []string
	for _, key := range keys {
		val, err := f.Get(key)
		if err == nil && val != "" {
			var graph Graph
			if err := json.Unmarshal([]byte(val), &graph); err == nil && (len(graph.Nodes) > 0 || len(graph.Edges) > 0) {
				graphKeys = append(graphKeys, key[len(key)-1])
			}
		}
	}
	return graphKeys, nil
}

func (f *FileSystem) GraphProfile(_ dotpip.Key, query string) ([]string, error) {
	q, err := dotpip.CypherParser.ParseString("", query)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, clause := range q.Clauses {
		switch {
		case clause.Create != nil:
			result = append(result, string(dotpip.GraphKeywordCreate))
		case clause.Match != nil:
			result = append(result, string(dotpip.GraphKeywordMatch))
		case clause.Return != nil:
			result = append(result, string(dotpip.GraphKeywordReturn))
		case clause.Delete != nil:
			result = append(result, string(dotpip.GraphKeywordDelete))
		case clause.Set != nil:
			result = append(result, string(dotpip.GraphKeywordSet))
		}
	}
	return result, nil
}

func (f *FileSystem) GraphQuery(key dotpip.Key, query string) ([]map[string]any, error) {
	q, err := dotpip.CypherParser.ParseString("", query)
	if err != nil {
		return nil, err
	}

	var graph Graph
	val, err := f.Get(key)
	if err == nil && val != "" {
		_ = json.Unmarshal([]byte(val), &graph)
	}

	result := []map[string]any{}

	for _, clause := range q.Clauses {
		switch {
		case clause.Create != nil:
			m := make(map[string]any)

			if clause.Create.Pattern != nil && clause.Create.Pattern.Node != nil {
				labels := clause.Create.Pattern.Node.Labels

				node := &GraphNode{
					ID:         len(graph.Nodes) + 1,
					Labels:     labels,
					Properties: make(map[string]any),
				}
				graph.Nodes = append(graph.Nodes, node)

				if len(labels) > 0 {
					m[string(dotpip.GraphKeywordLabelsAdded)] = len(labels)
				}
				m[string(dotpip.GraphKeywordNodesCreated)] = 1
				m[string(dotpip.GraphKeywordPropertiesSet)] = 0

				if clause.Create.Pattern.Node.Properties != nil {
					m[string(dotpip.GraphKeywordPropertiesSet)] = len(clause.Create.Pattern.Node.Properties.Props)
				}

				if len(clause.Create.Pattern.Chain) > 0 {
					m[string(dotpip.GraphKeywordRelationshipsCreated)] = len(clause.Create.Pattern.Chain)
					m[string(dotpip.GraphKeywordNodesCreated)] = m[string(dotpip.GraphKeywordNodesCreated)].(int) + len(clause.Create.Pattern.Chain)

					for _, chain := range clause.Create.Pattern.Chain {
						targetNode := &GraphNode{
							ID:         len(graph.Nodes) + 1,
							Labels:     chain.Node.Labels,
							Properties: make(map[string]any),
						}
						graph.Nodes = append(graph.Nodes, targetNode)

						edgeType := ""
						if chain.Relationship.Details != nil && len(chain.Relationship.Details.Types) > 0 {
							edgeType = chain.Relationship.Details.Types[0]
						}

						edge := &GraphEdge{
							ID:         len(graph.Edges) + 1,
							Type:       edgeType,
							SourceNode: node.ID,
							TargetNode: targetNode.ID,
							Properties: make(map[string]any),
						}
						graph.Edges = append(graph.Edges, edge)
						node = targetNode // chain moves forward
					}
				}
			}

			result = append(result, m)
		case clause.Match != nil:
			m := make(map[string]any)

			chainLen := 0
			if clause.Match.Pattern != nil {
				chainLen = len(clause.Match.Pattern.Chain)
			}

			paths := evaluateMatches(&graph, chainLen)
			m[string(dotpip.GraphKeywordNodesFound)] = len(graph.Nodes)
			m[string(dotpip.GraphKeywordPathsMatched)] = paths

			result = append(result, m)
		case clause.Return != nil:
			// Not implemented yet
		case clause.Delete != nil:
			m := make(map[string]any)
			m[string(dotpip.GraphKeywordNodesDeleted)] = len(graph.Nodes)
			graph.Nodes = []*GraphNode{}
			graph.Edges = []*GraphEdge{}
			result = append(result, m)
		case clause.Set != nil:
			m := make(map[string]any)
			m[string(dotpip.GraphKeywordPropertiesSet)] = len(clause.Set.Items)
			result = append(result, m)
		}
	}

	b, _ := json.Marshal(graph)
	_, _ = f.Set(key, string(b))

	if len(result) == 0 {
		return []map[string]any{{"Query Execution Time": "0ms"}}, nil
	}

	return result, nil
}

func (f *FileSystem) GraphROQuery(key dotpip.Key, query string) ([]map[string]any, error) {
	q, err := dotpip.CypherParser.ParseString("", query)
	if err != nil {
		return nil, err
	}

	var graph Graph
	val, err := f.Get(key)
	if err == nil && val != "" {
		_ = json.Unmarshal([]byte(val), &graph)
	}

	result := []map[string]any{}

	for _, clause := range q.Clauses {
		switch {
		case clause.Match != nil:
			m := make(map[string]any)

			chainLen := 0
			if clause.Match.Pattern != nil {
				chainLen = len(clause.Match.Pattern.Chain)
			}

			paths := evaluateMatches(&graph, chainLen)
			m[string(dotpip.GraphKeywordNodesFound)] = len(graph.Nodes)
			m[string(dotpip.GraphKeywordPathsMatched)] = paths

			result = append(result, m)
		case clause.Return != nil:
			// Not implemented yet
		case clause.Create != nil || clause.Delete != nil || clause.Set != nil:
			return nil, errors.New(string(dotpip.ErrMsgReadOnlyQuery))
		}
	}

	if len(result) == 0 {
		return []map[string]any{{"Query Execution Time": "0ms"}}, nil
	}

	return result, nil
}

// GraphSlowlog simulates returning a slowlog.
// For a filesystem dummy implementation, this just returns an empty list or a mock entry.
func (f *FileSystem) GraphSlowlog(_ dotpip.Key) ([]any, error) {
	// A real implementation would maintain an internal log array for execution times.
	// For now, we mock the return to comply with the required interface execution properly.
	mockLog := []any{
		"1600000000",         // timestamp
		"MATCH (n) RETURN n", // query
		"10",                 // execution time
	}
	return []any{mockLog}, nil
}
