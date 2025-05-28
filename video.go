// TODO: 値 or ポインタ？

package main

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Resource interfaceの実装
type YTvideo struct {
	url string
}

func (*YTvideo) Download() ([]byte, error) {
	// TODO: 未実装
	return []byte("data"), nil
}

func NewYTvideo(url string) *YTvideo {
	return &YTvideo{url: url}
}

// Finder interfaceの実装
type YTvideoFinder struct {
	apiKey string
}

func (v *YTvideoFinder) Find(max int) ([]Resource, error) {
	videoIDs, err := v.findVideoIDs(max)
	if err != nil {
		return nil, err
	}

	var videos []Resource
	for _, id := range videoIDs {
		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)
		videos = append(videos, NewYTvideo(videoURL))
	}

	return videos, nil
}

// YouTube動画を検索し、videoIDのリストを返す
func (v *YTvideoFinder) findVideoIDs(max int) ([]string, error) {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(v.apiKey))
	if err != nil {
		return nil, fmt.Errorf("YouTube APIの初期化エラー: %v", err)
	}

	// 卓球に関する動画を検索するクエリを作成
	// TODO: max>50の時の処理
	// TODO: クエリを最適化したい
	call := youtubeService.Search.List([]string{"id"}).
		Q("table tennis").       // 卓球の英語表記
		MaxResults(int64(max)).  // 最大結果数
		Type("video").           // 動画のみ
		RelevanceLanguage("ja"). // 日本語の動画を優先
		VideoEmbeddable("true")  // 埋め込み可能な動画のみ

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("YouTube API検索エラー: %v", err)
	}

	var videoIDs []string
	for _, item := range response.Items {
		videoIDs = append(videoIDs, item.Id.VideoId)
	}

	return videoIDs, nil
}

func NewYTvideoFinder(apiKey string) *YTvideoFinder {
	return &YTvideoFinder{apiKey: apiKey}
}
