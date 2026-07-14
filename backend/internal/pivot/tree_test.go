package pivot

import "testing"

func TestTreeIntegrity(t *testing.T) {
	if _, ok := NodeByID[RootID]; !ok {
		t.Fatalf("root node %q missing", RootID)
	}

	seenNode := map[string]bool{}
	for _, n := range Nodes {
		if seenNode[n.ID] {
			t.Errorf("duplicate node id %q", n.ID)
		}
		seenNode[n.ID] = true
		if n.Question == "" {
			t.Errorf("node %q has empty question", n.ID)
		}
		if len(n.Options) < 2 {
			t.Errorf("node %q needs at least 2 options", n.ID)
		}
		for _, opt := range n.Options {
			if (opt.Next == "") == (opt.Outcome == "") {
				t.Errorf("node %q option %q must set exactly one of next/outcome", n.ID, opt.Label)
			}
			if opt.Next != "" {
				if _, ok := NodeByID[opt.Next]; !ok {
					t.Errorf("node %q option %q points to unknown node %q", n.ID, opt.Label, opt.Next)
				}
			}
			if opt.Outcome != "" {
				if _, ok := OutcomeByID[opt.Outcome]; !ok {
					t.Errorf("node %q option %q points to unknown outcome %q", n.ID, opt.Label, opt.Outcome)
				}
			}
		}
	}
}

func TestAllNodesReachableAndNoCycles(t *testing.T) {
	// BFS from root; a node visited twice along one path would loop forever,
	// so bound the walk by the node count.
	visited := map[string]bool{}
	queue := []string{RootID}
	for steps := 0; len(queue) > 0; steps++ {
		if steps > len(Nodes)*len(Nodes) {
			t.Fatal("walk did not terminate - probable cycle")
		}
		id := queue[0]
		queue = queue[1:]
		if visited[id] {
			continue
		}
		visited[id] = true
		for _, opt := range NodeByID[id].Options {
			if opt.Next != "" {
				queue = append(queue, opt.Next)
			}
		}
	}
	for _, n := range Nodes {
		if !visited[n.ID] {
			t.Errorf("node %q unreachable from root", n.ID)
		}
	}
}

func TestAllOutcomesReachable(t *testing.T) {
	reachable := map[string]bool{}
	for _, n := range Nodes {
		for _, opt := range n.Options {
			if opt.Outcome != "" {
				reachable[opt.Outcome] = true
			}
		}
	}
	for _, o := range Outcomes {
		if !reachable[o.ID] {
			t.Errorf("outcome %q not reachable from any node", o.ID)
		}
		if o.Title == "" || o.WhyNow == "" || len(o.Plan) < 3 || len(o.Resources) == 0 {
			t.Errorf("outcome %q is missing content (title/whyNow/plan/resources)", o.ID)
		}
	}
}

func TestWalkPathEmpty(t *testing.T) {
	state, err := WalkPath(nil)
	if err != nil {
		t.Fatal(err)
	}
	if state.Current == nil || state.Current.ID != RootID {
		t.Fatalf("empty walk should sit at root, got %+v", state)
	}
	if state.Outcome != nil {
		t.Error("empty walk should have no outcome")
	}
}

func TestWalkPathToOutcome(t *testing.T) {
	steps := []Step{
		{NodeID: "driver", Option: "I'm burnt out - I need more time and energy for life"},
		{NodeID: "hours-fix", Option: "Yes - the job is fine, it's the hours"},
		{NodeID: "employer-flex", Option: "Open - there are precedents for part-time or 4-day weeks"},
	}
	state, err := WalkPath(steps)
	if err != nil {
		t.Fatal(err)
	}
	if state.Outcome == nil || state.Outcome.ID != "reduce-hours" {
		t.Fatalf("want reduce-hours outcome, got %+v", state)
	}
	if state.Current != nil {
		t.Error("finished walk should have no current node")
	}
}

func TestWalkPathMidway(t *testing.T) {
	steps := []Step{
		{NodeID: "driver", Option: "I want more autonomy and ownership over my work"},
	}
	state, err := WalkPath(steps)
	if err != nil {
		t.Fatal(err)
	}
	if state.Current == nil || state.Current.ID != "autonomy-kind" {
		t.Fatalf("want autonomy-kind, got %+v", state)
	}
}

func TestWalkPathRejectsWrongNode(t *testing.T) {
	steps := []Step{{NodeID: "hours-fix", Option: "Yes - the job is fine, it's the hours"}}
	if _, err := WalkPath(steps); err == nil {
		t.Fatal("expected error: first step must be the root node")
	}
}

func TestWalkPathRejectsUnknownOption(t *testing.T) {
	steps := []Step{{NodeID: "driver", Option: "I want a pony"}}
	if _, err := WalkPath(steps); err == nil {
		t.Fatal("expected error for unknown option")
	}
}

func TestWalkPathRejectsStepsAfterOutcome(t *testing.T) {
	steps := []Step{
		{NodeID: "driver", Option: "I'm burnt out - I need more time and energy for life"},
		{NodeID: "hours-fix", Option: "Yes - the job is fine, it's the hours"},
		{NodeID: "employer-flex", Option: "Open - there are precedents for part-time or 4-day weeks"},
		{NodeID: "driver", Option: "The work itself no longer fits me"},
	}
	if _, err := WalkPath(steps); err == nil {
		t.Fatal("expected error for steps after an outcome")
	}
}

// Every leaf of the tree must be walkable end to end: enumerate all paths by
// DFS and verify WalkPath accepts each.
func TestEveryPathWalks(t *testing.T) {
	var dfs func(nodeID string, steps []Step)
	paths := 0
	dfs = func(nodeID string, steps []Step) {
		node := NodeByID[nodeID]
		for _, opt := range node.Options {
			next := append(append([]Step{}, steps...), Step{NodeID: nodeID, Option: opt.Label})
			if opt.Outcome != "" {
				paths++
				state, err := WalkPath(next)
				if err != nil {
					t.Errorf("path via %q failed: %v", opt.Label, err)
				} else if state.Outcome == nil || state.Outcome.ID != opt.Outcome {
					t.Errorf("path via %q: want outcome %q, got %+v", opt.Label, opt.Outcome, state)
				}
			} else {
				dfs(opt.Next, next)
			}
		}
	}
	dfs(RootID, nil)
	if paths < 10 {
		t.Errorf("expected a meaningful tree, only %d complete paths", paths)
	}
	t.Logf("verified %d complete paths", paths)
}
