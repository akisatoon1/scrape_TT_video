package main

type YTvideos []YTvideo

func (v *YTvideos) Find() ([]Resource, error) {
	// 仮の実装: YTvideoオブジェクトのリストを作成
	videos := []Resource{
		&YTvideo{url: "https://example.com/video1"},
		&YTvideo{url: "https://example.com/video2"},
		&YTvideo{url: "https://example.com/video3"},
	}

	return videos, nil
}

func NewYTvideos() *YTvideos {
	return &YTvideos{}
}
