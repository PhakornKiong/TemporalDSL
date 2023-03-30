package dsl

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

type WorkflowData struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Nodes       []*Node      `yaml:"nodes"`
	Connections []Connection `yaml:"connections"`
}

type Node struct {
	Name   string                 `yaml:"name"`
	Type   string                 `yaml:"type"`
	Input  map[string]interface{} `yaml:"input"`
	Output interface{}            `yaml:"output"`
}

type Connection struct {
	From string   `yaml:"from"`
	To   []string `yaml:"to"`
}

type CustomWorkflow struct {
	WF     func(ctx workflow.Context) error
	Name   string
	WfData WorkflowData
}

func (wf *CustomWorkflow) Init() error {
	wf.Name = wf.WfData.Name

	nodes := make(map[string]Node)
	for _, node := range wf.WfData.Nodes {
		nodes[node.Name] = *node
	}

	// Use toposort to determine order of execution
	sortedNodes, err := toposort(wf.WfData.Connections, func(node string) []string {
		toNodes := []string{}
		for _, connection := range wf.WfData.Connections {
			if connection.From == node {
				toNodes = append(toNodes, connection.To...)
			}
		}
		return toNodes
	})

	if err != nil {
		return err
	}

	wf.WF = func(ctx workflow.Context) error {
		ao := workflow.ActivityOptions{
			StartToCloseTimeout: 10 * time.Second,
		}
		ctx = workflow.WithActivityOptions(ctx, ao)

		// Execute each node as an activity.
		for key, nodeKey := range sortedNodes {

			node := nodes[nodeKey]

			prevNode := Node{}
			if key > 0 {
				prevNode = nodes[sortedNodes[key-1]]
			}

			input := node.Input["main"]

			// If previous node has output, use it as input
			if prevNode.Output != nil {
				input = prevNode.Output
			}

			var result interface{}

			err = workflow.ExecuteActivity(ctx, node.Name, input).Get(ctx, &result)
			node.Output = result

			if err != nil {
				return err
			}

			// Update node with output
			nodes[nodeKey] = node
		}

		return nil
	}

	return nil
}

// Adapted from https://pkg.go.dev/github.com/philopon/go-toposort
func toposort(edges []Connection, adjacent func(string) []string) ([]string, error) {
	var sorted []string
	permanent := make(map[string]bool)
	temporary := make(map[string]bool)
	visited := make(map[string]bool)
	for _, e := range edges {
		permanent[e.From] = false
		for _, t := range e.To {
			permanent[t] = false
		}
	}

	var visit func(string) error
	visit = func(node string) error {
		if visited[node] {
			return nil
		}
		if temporary[node] {
			return fmt.Errorf("cycle detected")
		}
		temporary[node] = true
		for _, n := range adjacent(node) {
			if err := visit(n); err != nil {
				return err
			}
		}
		temporary[node] = false
		visited[node] = true
		sorted = append(sorted, node)
		return nil
	}

	for n := range permanent {
		if !visited[n] {
			if err := visit(n); err != nil {
				return nil, err
			}
		}
	}
	for i, j := 0, len(sorted)-1; i < j; i, j = i+1, j-1 {
		sorted[i], sorted[j] = sorted[j], sorted[i]
	}
	return sorted, nil
}
