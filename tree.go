// Package tree builds a hierarchical tree of locations (e.g. a property,
// its floors, sections, rooms, storage areas, etc.) from a flat list.
//
// Each Node exposes GetID, which returns a path-like string identifying the
// node within the tree. Other parts of the system rely on this identifier as a
// unique key (for example, when indexing nodes in a map, caching results, or
// emitting structured logs).
package tree

import (
	"fmt"
	"sort"
)

const RootLocationType = "Property"

// Location is the raw input describing a node in the hierarchy.
type Location struct {
	ID           int64
	Name         string
	ParentID     *int64
	LocationType string
	Sequence     float64
}

// Node is a location placed inside the tree.
type Node struct {
	ID       int64
	Name     string
	Location *Location
	Parent   *Node
	Children []*Node
	Next     *Node
	Previous *Node
}

// BuildTree creates a tree from a flat slice of Locations. It returns:
//   - a map from Location ID to its Node, for quick lookup
//   - the root Node of the tree
//   - an error if the input is invalid
//
// Locations with a ParentID that points to a missing parent are skipped.
func BuildTree(locations []Location) (map[int64]*Node, *Node, error) {
	locationMap := make(map[int64]*Node, len(locations))
	for i := range locations {
		loc := &locations[i]
		locationMap[loc.ID] = &Node{
			ID:       loc.ID,
			Name:     loc.Name,
			Location: loc,
			Children: []*Node{},
		}
	}

	var root *Node
	for i := range locations {
		loc := &locations[i]
		node := locationMap[loc.ID]

		if loc.ParentID == nil {
			if loc.LocationType == RootLocationType || root == nil {
				root = node
			}
			continue
		}

		parent, ok := locationMap[*loc.ParentID]
		if !ok {
			continue
		}
		parent.Children = append(parent.Children, node)
		node.Parent = parent
	}

	sortChildren(root)

	return locationMap, root, nil
}

// sortChildren sorts each node's children by Sequence (when any child has a
// positive sequence) or by Name otherwise. It also wires up sibling Next /
// Previous pointers so callers can traverse a level in order.
func sortChildren(node *Node) {
	if node == nil {
		return
	}

	hasSequence := false
	for _, c := range node.Children {
		if c != nil && c.Location != nil && c.Location.Sequence > 0 {
			hasSequence = true
			break
		}
	}

	if hasSequence {
		sort.Slice(node.Children, func(i, j int) bool {
			return node.Children[i].Location.Sequence < node.Children[j].Location.Sequence
		})
	} else {
		sort.Slice(node.Children, func(i, j int) bool {
			return node.Children[i].Name < node.Children[j].Name
		})
	}

	for i := 0; i < len(node.Children)-1; i++ {
		node.Children[i].Next = node.Children[i+1]
		node.Children[i+1].Previous = node.Children[i]
	}

	for _, c := range node.Children {
		sortChildren(c)
	}
}

// GetID returns a path-style identifier for this node, built by walking up
// to the root and joining the names along the way.
//
// The returned value is used elsewhere as a key in maps and for logging.
func (n *Node) GetID() string {
	key := n.Name

	tempParent := n.Parent
	for tempParent != nil {
		key = fmt.Sprintf("%s/%s", tempParent.Name, key)
		tempParent = tempParent.Parent
	}

	return key
}

