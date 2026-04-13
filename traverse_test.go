package tree

import (
	"context"
	"testing"
)

// TestTraverse_GrandPlazaHotel reproduces a customer report: when the
// property's location tree is walked with Traverse, at least one node is
// silently skipped — its callback never fires.
//
// With the data below, Traverse should invoke the callback once for every
// location in the fixture. It doesn't.
//
// This test currently fails. Diagnose why and fix it.
func TestTraverse_GrandPlazaHotel(t *testing.T) {
	// Sequence values are set on Floor 1's children to keep the traversal
	// order stable across sort implementations, so the failure is
	// deterministic and easy to reason about.
	locations := []Location{
		// Property
		{ID: 1, Name: "Grand Plaza Hotel", LocationType: "Property"},

		// Floors
		{ID: 100, Name: "Basement", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 200, Name: "Ground Floor", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 300, Name: "Mezzanine", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 1000, Name: "Floor 1", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 2000, Name: "Floor 2", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 3000, Name: "Floor 3", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 4000, Name: "Floor 4", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 5000, Name: "Floor 5", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 6000, Name: "Floor 6", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 7000, Name: "Floor 7", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 8000, Name: "Floor 8", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 9000, Name: "Floor 9", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 10000, Name: "Floor 10", ParentID: ptr(1), LocationType: "Floor"},
		{ID: 50, Name: "Rooftop", ParentID: ptr(1), LocationType: "Floor"},

		// Basement
		{ID: 101, Name: "Loading Dock", ParentID: ptr(100), LocationType: "Section"},
		{ID: 102, Name: "Main Kitchen", ParentID: ptr(100), LocationType: "Section"},
		{ID: 103, Name: "Cold Storage", ParentID: ptr(100), LocationType: "Storage"},
		{ID: 104, Name: "Dry Storage", ParentID: ptr(100), LocationType: "Storage"},
		{ID: 105, Name: "Laundry", ParentID: ptr(100), LocationType: "Section"},
		{ID: 106, Name: "Maintenance Workshop", ParentID: ptr(100), LocationType: "Section"},
		{ID: 107, Name: "Staff Lounge", ParentID: ptr(100), LocationType: "Section"},
		{ID: 108, Name: "HR Office", ParentID: ptr(100), LocationType: "Section"},
		{ID: 109, Name: "Security Office", ParentID: ptr(100), LocationType: "Section"},

		// Ground Floor
		{ID: 201, Name: "Reception", ParentID: ptr(200), LocationType: "Section"},
		{ID: 202, Name: "Concierge", ParentID: ptr(200), LocationType: "Section"},
		{ID: 203, Name: "Main Lobby", ParentID: ptr(200), LocationType: "Section"},
		{ID: 204, Name: "Lobby Bar", ParentID: ptr(200), LocationType: "Section"},
		{ID: 205, Name: "Harbour Restaurant", ParentID: ptr(200), LocationType: "Section"},
		{ID: 206, Name: "Business Center", ParentID: ptr(200), LocationType: "Section"},
		{ID: 207, Name: "Gift Shop", ParentID: ptr(200), LocationType: "Section"},
		{ID: 208, Name: "Cloakroom", ParentID: ptr(200), LocationType: "Storage"},

		// Mezzanine
		{ID: 301, Name: "Ballroom A", ParentID: ptr(300), LocationType: "Section"},
		{ID: 302, Name: "Ballroom B", ParentID: ptr(300), LocationType: "Section"},
		{ID: 303, Name: "Meeting Room 1", ParentID: ptr(300), LocationType: "Section"},
		{ID: 304, Name: "Meeting Room 2", ParentID: ptr(300), LocationType: "Section"},
		{ID: 305, Name: "Meeting Room 3", ParentID: ptr(300), LocationType: "Section"},
		{ID: 306, Name: "Pre-function", ParentID: ptr(300), LocationType: "Section"},
		{ID: 307, Name: "Service Kitchen", ParentID: ptr(300), LocationType: "Section"},

		// Floor 1
		{ID: 1001, Name: "East Wing", ParentID: ptr(1000), LocationType: "Section", Sequence: 1},
		{ID: 1002, Name: "West Wing", ParentID: ptr(1000), LocationType: "Section", Sequence: 3},
		{ID: 1003, Name: "East Wing", ParentID: ptr(1000), LocationType: "Storage", Sequence: 2},
		{ID: 1004, Name: "Housekeeping Closet", ParentID: ptr(1000), LocationType: "Storage", Sequence: 4},
		{ID: 1010, Name: "101", ParentID: ptr(1001), LocationType: "Room"},
		{ID: 1011, Name: "102", ParentID: ptr(1001), LocationType: "Room"},
		{ID: 1012, Name: "103", ParentID: ptr(1001), LocationType: "Room"},
		{ID: 1013, Name: "104", ParentID: ptr(1001), LocationType: "Room"},
		{ID: 1014, Name: "105", ParentID: ptr(1001), LocationType: "Room"},
		{ID: 1015, Name: "106", ParentID: ptr(1002), LocationType: "Room"},
		{ID: 1016, Name: "107", ParentID: ptr(1002), LocationType: "Room"},
		{ID: 1017, Name: "108", ParentID: ptr(1002), LocationType: "Room"},
		{ID: 1018, Name: "109", ParentID: ptr(1002), LocationType: "Room"},
		{ID: 1019, Name: "110", ParentID: ptr(1002), LocationType: "Room"},

		// Floor 2
		{ID: 2001, Name: "East Wing", ParentID: ptr(2000), LocationType: "Section"},
		{ID: 2002, Name: "West Wing", ParentID: ptr(2000), LocationType: "Section"},
		{ID: 2003, Name: "Housekeeping Closet", ParentID: ptr(2000), LocationType: "Storage"},
		{ID: 2010, Name: "201", ParentID: ptr(2001), LocationType: "Room"},
		{ID: 2011, Name: "202", ParentID: ptr(2001), LocationType: "Room"},
		{ID: 2012, Name: "203", ParentID: ptr(2001), LocationType: "Room"},
		{ID: 2013, Name: "204", ParentID: ptr(2001), LocationType: "Room"},
		{ID: 2014, Name: "205", ParentID: ptr(2001), LocationType: "Room"},
		{ID: 2015, Name: "206", ParentID: ptr(2002), LocationType: "Room"},
		{ID: 2016, Name: "207", ParentID: ptr(2002), LocationType: "Room"},
		{ID: 2017, Name: "208", ParentID: ptr(2002), LocationType: "Room"},
		{ID: 2018, Name: "209", ParentID: ptr(2002), LocationType: "Room"},
		{ID: 2019, Name: "210", ParentID: ptr(2002), LocationType: "Room"},

		// Floor 3
		{ID: 3001, Name: "East Wing", ParentID: ptr(3000), LocationType: "Section"},
		{ID: 3002, Name: "West Wing", ParentID: ptr(3000), LocationType: "Section"},
		{ID: 3003, Name: "Housekeeping Closet", ParentID: ptr(3000), LocationType: "Storage"},
		{ID: 3010, Name: "301", ParentID: ptr(3001), LocationType: "Room"},
		{ID: 3011, Name: "302", ParentID: ptr(3001), LocationType: "Room"},
		{ID: 3012, Name: "303", ParentID: ptr(3001), LocationType: "Room"},
		{ID: 3013, Name: "304", ParentID: ptr(3001), LocationType: "Room"},
		{ID: 3014, Name: "305", ParentID: ptr(3001), LocationType: "Room"},
		{ID: 3015, Name: "306", ParentID: ptr(3002), LocationType: "Room"},
		{ID: 3016, Name: "307", ParentID: ptr(3002), LocationType: "Room"},
		{ID: 3017, Name: "308", ParentID: ptr(3002), LocationType: "Room"},
		{ID: 3018, Name: "309", ParentID: ptr(3002), LocationType: "Room"},
		{ID: 3019, Name: "310", ParentID: ptr(3002), LocationType: "Room"},

		// Floor 4
		{ID: 4001, Name: "East Wing", ParentID: ptr(4000), LocationType: "Section"},
		{ID: 4002, Name: "West Wing", ParentID: ptr(4000), LocationType: "Section"},
		{ID: 4003, Name: "Housekeeping Closet", ParentID: ptr(4000), LocationType: "Storage"},
		{ID: 4010, Name: "401", ParentID: ptr(4001), LocationType: "Room"},
		{ID: 4011, Name: "402", ParentID: ptr(4001), LocationType: "Room"},
		{ID: 4012, Name: "403", ParentID: ptr(4001), LocationType: "Room"},
		{ID: 4013, Name: "404", ParentID: ptr(4001), LocationType: "Room"},
		{ID: 4014, Name: "405", ParentID: ptr(4001), LocationType: "Room"},
		{ID: 4015, Name: "406", ParentID: ptr(4002), LocationType: "Room"},
		{ID: 4016, Name: "407", ParentID: ptr(4002), LocationType: "Room"},
		{ID: 4017, Name: "408", ParentID: ptr(4002), LocationType: "Room"},
		{ID: 4018, Name: "409", ParentID: ptr(4002), LocationType: "Room"},
		{ID: 4019, Name: "410", ParentID: ptr(4002), LocationType: "Room"},

		// Floor 5
		{ID: 5001, Name: "East Wing", ParentID: ptr(5000), LocationType: "Section"},
		{ID: 5002, Name: "West Wing", ParentID: ptr(5000), LocationType: "Section"},
		{ID: 5003, Name: "Housekeeping Closet", ParentID: ptr(5000), LocationType: "Storage"},
		{ID: 5010, Name: "501", ParentID: ptr(5001), LocationType: "Room"},
		{ID: 5011, Name: "502", ParentID: ptr(5001), LocationType: "Room"},
		{ID: 5012, Name: "503", ParentID: ptr(5001), LocationType: "Room"},
		{ID: 5013, Name: "504", ParentID: ptr(5001), LocationType: "Room"},
		{ID: 5014, Name: "505", ParentID: ptr(5001), LocationType: "Room"},
		{ID: 5015, Name: "506", ParentID: ptr(5002), LocationType: "Room"},
		{ID: 5016, Name: "507", ParentID: ptr(5002), LocationType: "Room"},
		{ID: 5017, Name: "508", ParentID: ptr(5002), LocationType: "Room"},
		{ID: 5018, Name: "509", ParentID: ptr(5002), LocationType: "Room"},
		{ID: 5019, Name: "510", ParentID: ptr(5002), LocationType: "Room"},

		// Floor 6
		{ID: 6001, Name: "East Wing", ParentID: ptr(6000), LocationType: "Section"},
		{ID: 6002, Name: "West Wing", ParentID: ptr(6000), LocationType: "Section"},
		{ID: 6003, Name: "Housekeeping Closet", ParentID: ptr(6000), LocationType: "Storage"},
		{ID: 6010, Name: "601", ParentID: ptr(6001), LocationType: "Room"},
		{ID: 6011, Name: "602", ParentID: ptr(6001), LocationType: "Room"},
		{ID: 6012, Name: "603", ParentID: ptr(6001), LocationType: "Room"},
		{ID: 6013, Name: "604", ParentID: ptr(6001), LocationType: "Room"},
		{ID: 6014, Name: "605", ParentID: ptr(6001), LocationType: "Room"},
		{ID: 6015, Name: "606", ParentID: ptr(6002), LocationType: "Room"},
		{ID: 6016, Name: "607", ParentID: ptr(6002), LocationType: "Room"},
		{ID: 6017, Name: "608", ParentID: ptr(6002), LocationType: "Room"},
		{ID: 6018, Name: "609", ParentID: ptr(6002), LocationType: "Room"},
		{ID: 6019, Name: "610", ParentID: ptr(6002), LocationType: "Room"},

		// Floor 7
		{ID: 7001, Name: "East Wing", ParentID: ptr(7000), LocationType: "Section"},
		{ID: 7002, Name: "West Wing", ParentID: ptr(7000), LocationType: "Section"},
		{ID: 7003, Name: "Housekeeping Closet", ParentID: ptr(7000), LocationType: "Storage"},
		{ID: 7010, Name: "701", ParentID: ptr(7001), LocationType: "Room"},
		{ID: 7011, Name: "702", ParentID: ptr(7001), LocationType: "Room"},
		{ID: 7012, Name: "703", ParentID: ptr(7001), LocationType: "Room"},
		{ID: 7013, Name: "704", ParentID: ptr(7001), LocationType: "Room"},
		{ID: 7014, Name: "705", ParentID: ptr(7001), LocationType: "Room"},
		{ID: 7015, Name: "706", ParentID: ptr(7002), LocationType: "Room"},
		{ID: 7016, Name: "707", ParentID: ptr(7002), LocationType: "Room"},
		{ID: 7017, Name: "708", ParentID: ptr(7002), LocationType: "Room"},
		{ID: 7018, Name: "709", ParentID: ptr(7002), LocationType: "Room"},
		{ID: 7019, Name: "710", ParentID: ptr(7002), LocationType: "Room"},

		// Floor 8
		{ID: 8001, Name: "East Wing", ParentID: ptr(8000), LocationType: "Section"},
		{ID: 8002, Name: "West Wing", ParentID: ptr(8000), LocationType: "Section"},
		{ID: 8003, Name: "Housekeeping Closet", ParentID: ptr(8000), LocationType: "Storage"},
		{ID: 8010, Name: "801", ParentID: ptr(8001), LocationType: "Room"},
		{ID: 8011, Name: "802", ParentID: ptr(8001), LocationType: "Room"},
		{ID: 8012, Name: "803", ParentID: ptr(8001), LocationType: "Room"},
		{ID: 8013, Name: "804", ParentID: ptr(8001), LocationType: "Room"},
		{ID: 8014, Name: "805", ParentID: ptr(8001), LocationType: "Room"},
		{ID: 8015, Name: "806", ParentID: ptr(8002), LocationType: "Room"},
		{ID: 8016, Name: "807", ParentID: ptr(8002), LocationType: "Room"},
		{ID: 8017, Name: "808", ParentID: ptr(8002), LocationType: "Room"},
		{ID: 8018, Name: "809", ParentID: ptr(8002), LocationType: "Room"},
		{ID: 8019, Name: "810", ParentID: ptr(8002), LocationType: "Room"},

		// Floor 9
		{ID: 9001, Name: "East Wing", ParentID: ptr(9000), LocationType: "Section"},
		{ID: 9002, Name: "West Wing", ParentID: ptr(9000), LocationType: "Section"},
		{ID: 9003, Name: "Housekeeping Closet", ParentID: ptr(9000), LocationType: "Storage"},
		{ID: 9010, Name: "901", ParentID: ptr(9001), LocationType: "Room"},
		{ID: 9011, Name: "902", ParentID: ptr(9001), LocationType: "Room"},
		{ID: 9012, Name: "903", ParentID: ptr(9001), LocationType: "Room"},
		{ID: 9013, Name: "904", ParentID: ptr(9001), LocationType: "Room"},
		{ID: 9014, Name: "905", ParentID: ptr(9001), LocationType: "Room"},
		{ID: 9015, Name: "906", ParentID: ptr(9002), LocationType: "Room"},
		{ID: 9016, Name: "907", ParentID: ptr(9002), LocationType: "Room"},
		{ID: 9017, Name: "908", ParentID: ptr(9002), LocationType: "Room"},
		{ID: 9018, Name: "909", ParentID: ptr(9002), LocationType: "Room"},
		{ID: 9019, Name: "910", ParentID: ptr(9002), LocationType: "Room"},

		// Floor 10
		{ID: 10001, Name: "East Wing", ParentID: ptr(10000), LocationType: "Section"},
		{ID: 10002, Name: "West Wing", ParentID: ptr(10000), LocationType: "Section"},
		{ID: 10003, Name: "Housekeeping Closet", ParentID: ptr(10000), LocationType: "Storage"},
		{ID: 10010, Name: "1001", ParentID: ptr(10001), LocationType: "Room"},
		{ID: 10011, Name: "1002", ParentID: ptr(10001), LocationType: "Room"},
		{ID: 10012, Name: "1003", ParentID: ptr(10001), LocationType: "Room"},
		{ID: 10013, Name: "1004", ParentID: ptr(10001), LocationType: "Room"},
		{ID: 10014, Name: "1005", ParentID: ptr(10001), LocationType: "Room"},
		{ID: 10015, Name: "1006", ParentID: ptr(10002), LocationType: "Room"},
		{ID: 10016, Name: "1007", ParentID: ptr(10002), LocationType: "Room"},
		{ID: 10017, Name: "1008", ParentID: ptr(10002), LocationType: "Room"},
		{ID: 10018, Name: "1009", ParentID: ptr(10002), LocationType: "Room"},
		{ID: 10019, Name: "1010", ParentID: ptr(10002), LocationType: "Room"},

		// Rooftop
		{ID: 51, Name: "Pool Deck", ParentID: ptr(50), LocationType: "Section"},
		{ID: 52, Name: "Sky Bar", ParentID: ptr(50), LocationType: "Section"},
		{ID: 53, Name: "Spa", ParentID: ptr(50), LocationType: "Section"},
		{ID: 54, Name: "Gym", ParentID: ptr(50), LocationType: "Section"},
		{ID: 55, Name: "Pool Equipment", ParentID: ptr(50), LocationType: "Storage"},
	}

	_, root, err := BuildTree(locations)
	if err != nil {
		t.Fatalf("BuildTree returned error: %v", err)
	}

	visited := 0
	err = Traverse(context.Background(), root, map[string]bool{}, func(n *Node) error {
		visited++
		return nil
	})
	if err != nil {
		t.Fatalf("Traverse returned error: %v", err)
	}

	if visited != len(locations) {
		t.Errorf("Traverse did not invoke the callback for every location")
	}
}
