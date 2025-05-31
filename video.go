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

	// 最大max件の動画IDを取得する

	ytVideos := make([]Resource, 0, max)
	var nextPageToken string

	for restSize := int64(max); restSize > 0; restSize -= 50 {
		var requestSize int64
		if restSize > 50 {
			requestSize = 50
		} else {
			requestSize = restSize
		}

		// 卓球に関する動画を検索するクエリを作成
		call := youtubeService.Search.List([]string{"id"}).
			Q("table tennis").       // 卓球の英語表記
			MaxResults(requestSize). // 1回のリクエストあたりの最大結果数
			Type("video").           // 動画のみ
			RelevanceLanguage("ja"). // 日本語の動画を優先
			VideoEmbeddable("true")  // 埋め込み可能な動画のみ

		// 2ページ目以降の場合、nextPageTokenを設定
		if nextPageToken != "" {
			call = call.PageToken(nextPageToken)
		}

		// リクエスト実行
		response, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("YouTube API検索エラー: %v", err)
		}

		// 結果を追加
		for _, item := range response.Items {
			ytVideos = append(ytVideos, NewYTvideo(item.Id.VideoId))
		}

		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	return ytVideos, nil
}
