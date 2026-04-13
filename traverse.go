package tree

import "context"

// Traverse walks the tree starting at node, visiting each node's children and
// its same-level Next / Previous siblings. It invokes f exactly once per
// visited node. The visitedNodes map is used to avoid re-visiting a node that
// was already seen through another path.
func Traverse(ctx context.Context, node *Node, visitedNodes map[string]bool, f func(node *Node) error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if node == nil || visitedNodes[node.GetID()] {
		return nil
	}

	visitedNodes[node.GetID()] = true

	if err := f(node); err != nil {
		return err
	}

	for _, child := range node.Children {
		if err := Traverse(ctx, child, visitedNodes, f); err != nil {
			return err
		}
	}

	if node.Next != nil {
		if err := Traverse(ctx, node.Next, visitedNodes, f); err != nil {
			return err
		}
	}

	if node.Previous != nil {
		if err := Traverse(ctx, node.Previous, visitedNodes, f); err != nil {
			return err
		}
	}

	return nil
}
