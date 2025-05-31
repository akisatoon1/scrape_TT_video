package main

import (
	"os"
	"testing"
)

func TestFind(t *testing.T) {
	apiKey, err := getAPIKey()
	if err != nil {
		t.Fatalf("Error loading API key: %v", err)
	}
	if os.Getenv("RUN_API_TESTS") != "true" {
		t.Skip("APIテストはRUN_API_TESTS=trueで実行します")
	}

	finder := NewYTvideoFinder(apiKey)
	max := 111
	ytVideos, err := finder.Find(max)
	if err != nil {
		t.Fatalf("Failed to find video IDs: %v", err)
	}

	l := len(ytVideos)
	if l != max {
		t.Errorf("Expected %d video IDs, got %d", max, l)
	}
	t.Logf("Output first 10 video IDs:")
	for i, video := range ytVideos {
		if i <= 10 {
			t.Logf("Video ID: %s", video.ID())
		}
	}
}

func TestDownload(t *testing.T) {
	url := "https://www.youtube.com/watch?v=ekr2nIex040"
	t.Log("Testing video download from URL:", url)

	video := &YTvideo{id: url}
	err := video.Download("testvideo/testvideo.mp4")
	if err != nil {
		t.Fatalf("Failed to download video: %v", err)
	}

	t.Logf("Video saved successfully")
}
