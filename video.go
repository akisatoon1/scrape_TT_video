// TODO: 値 or ポインタ？

package main

import (
	"context"
	"fmt"

	ytdownloader "github.com/kkdai/youtube/v2/downloader"
	"google.golang.org/api/option"
	ytapi "google.golang.org/api/youtube/v3"
)

// Resource interfaceの実装
type YTvideo struct {
	id string
}

func NewYTvideo(id string) *YTvideo {
	return &YTvideo{id: id}
}

func (v *YTvideo) ID() string {
	return v.id
}

func (v *YTvideo) Download(filename string) error {
	downloader := ytdownloader.Downloader{}

	// 動画を取得
	video, err := downloader.GetVideo(v.id)
	if err != nil {
		return fmt.Errorf("YouTubeビデオの取得エラー: %v", err)
	}

	// 動画のフォーマットを取得
	formats := video.Formats.Type("video/mp4").Quality("medium")
	if len(formats) == 0 {
		return fmt.Errorf("動画のフォーマットが見つかりませんでした")
	}

	// 動画をダウンロード
	// TODO: contextについて
	ctx := context.Background()
	// TODO: フォーマットによる違いがある
	err = downloader.Download(ctx, video, &formats[0], filename)
	if err != nil {
		return fmt.Errorf("動画のダウンロードエラー: %v", err)
	}
	fmt.Printf("動画 %s を %s に保存しました\n", v.id, filename)

	return nil
}

// Finder interfaceの実装
type YTvideoFinder struct {
	apiKey string
}

func NewYTvideoFinder(apiKey string) *YTvideoFinder {
	return &YTvideoFinder{apiKey: apiKey}
}

func (v *YTvideoFinder) Find(max int) ([]Resource, error) {
	ctx := context.Background()
	youtubeService, err := ytapi.NewService(ctx, option.WithAPIKey(v.apiKey))
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

	ytVideos := make([]Resource, max)
	for i, item := range response.Items {
		ytVideos[i] = NewYTvideo(item.Id.VideoId)
	}

	return ytVideos, nil
}
