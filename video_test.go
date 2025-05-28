package main

import "testing"

func TestFindVideoIDs(t *testing.T) {
	apiKey, err := getAPIKey()
	if err != nil {
		t.Fatalf("Error loading API key: %v", err)
	}
	finder := NewYTvideoFinder(apiKey)
	max := 10
	videoIDs, err := finder.findVideoIDs(max)
	if err != nil {
		t.Fatalf("Failed to find video IDs: %v", err)
	}

	l := len(videoIDs)
	if l != max {
		t.Errorf("Expected %d video IDs, got %d", max, l)
	}
	t.Logf("Found videos len: %d", l)
	for _, id := range videoIDs {
		t.Logf("Video ID: %s", id)
	}
}
