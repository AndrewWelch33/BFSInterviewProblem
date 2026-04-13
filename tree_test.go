package tree

import "testing"

func ptr(i int64) *int64 { return &i }

// TestBuildTree_Shape verifies that BuildTree correctly wires parents and
// children for a simple hierarchy. This test should pass as-is.
func TestBuildTree_Shape(t *testing.T) {
	locations := []Location{
		{ID: 1, Name: "Property", LocationType: "Property"},
		{ID: 2, Name: "Floor 1", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 3, Name: "Room 101", ParentID: ptr(2), LocationType: "Room"},
	}

	locMap, root, err := BuildTree(locations)
	if err != nil {
		t.Fatalf("BuildTree returned error: %v", err)
	}
	if root == nil {
		t.Fatalf("expected a root node, got nil")
	}
	if root.ID != 1 {
		t.Fatalf("expected root ID 1, got %d", root.ID)
	}
	if len(root.Children) != 1 || root.Children[0].ID != 2 {
		t.Fatalf("expected property to have a single child Floor 1, got %+v", root.Children)
	}

	room := locMap[3]
	if room.Parent == nil || room.Parent.ID != 2 {
		t.Fatalf("expected Room 101's parent to be Floor 1, got %+v", room.Parent)
	}
}
