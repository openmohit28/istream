// Package pivot models the career-pivot guidance flow: a versioned decision
// tree that leads a person from "something needs to change" to a concrete,
// research-grounded outcome with an action plan. Users explore the tree in
// persisted threads that can be forked at any earlier answer.
package pivot

import "fmt"

type Option struct {
	Label string `json:"label"`
	// Exactly one of Next (a node ID) or Outcome (an outcome ID) is set.
	Next    string `json:"next,omitempty"`
	Outcome string `json:"outcome,omitempty"`
}

type Node struct {
	ID       string   `json:"id"`
	Question string   `json:"question"`
	Options  []Option `json:"options"`
}

type Resource struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Outcome struct {
	ID        string     `json:"id"`
	Path      string     `json:"path"` // reduce-hours | within-field | switch-out | consultancy | recharge
	Title     string     `json:"title"`
	Tagline   string     `json:"tagline"`
	WhyNow    string     `json:"whyNow"` // research-grounded context
	Plan      []string   `json:"plan"`
	Resources []Resource `json:"resources"`
}

// Step is one answered question in a thread.
type Step struct {
	NodeID string `json:"nodeId"`
	Option string `json:"option"`
}

// State is the computed position after walking a thread's steps.
type State struct {
	Current *Node    `json:"current,omitempty"` // next question, nil when done
	Outcome *Outcome `json:"outcome,omitempty"` // set when the walk reached an outcome
}

// WalkPath replays steps from the root and returns the resulting state.
// Every step must name the node the walk is actually at and pick one of its
// real options, so stored threads can never desync from the tree.
func WalkPath(steps []Step) (State, error) {
	current := RootID
	for i, step := range steps {
		if current == "" {
			return State{}, fmt.Errorf("step %d: path already reached an outcome", i+1)
		}
		node, ok := NodeByID[step.NodeID]
		if !ok {
			return State{}, fmt.Errorf("step %d: unknown node %q", i+1, step.NodeID)
		}
		if node.ID != current {
			return State{}, fmt.Errorf("step %d: expected node %q, got %q", i+1, current, step.NodeID)
		}
		opt := findOption(node, step.Option)
		if opt == nil {
			return State{}, fmt.Errorf("step %d: node %q has no option %q", i+1, node.ID, step.Option)
		}
		if opt.Outcome != "" {
			if i != len(steps)-1 {
				return State{}, fmt.Errorf("step %d: outcome reached before final step", i+1)
			}
			outcome := OutcomeByID[opt.Outcome]
			return State{Outcome: &outcome}, nil
		}
		current = opt.Next
	}
	node := NodeByID[current]
	return State{Current: &node}, nil
}

func findOption(node Node, label string) *Option {
	for i := range node.Options {
		if node.Options[i].Label == label {
			return &node.Options[i]
		}
	}
	return nil
}
