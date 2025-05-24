// TODO: 並行処理を実装

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var finder Finder = NewYTvideos()

	resources, err := finder.Find()
	if err != nil {
		log.Fatal(err)
	}

	// 逐次処理
	log.Println("動画の逐次ダウンロードを開始します")

	for i, downloader := range resources {
		log.Printf("動画 %d のダウンロードを開始します", i)

		// ダウンロード実行
		data, err := downloader.Download()
		if err != nil {
			log.Printf("Error downloading video %d: %v", i, err)
			continue
		}

		// ファイル名を生成
		filename := filepath.Join("downloads", fmt.Sprintf("video_%d.mp4", i))

		// ファイルに保存
		err = saveToFile(data, filename)
		if err != nil {
			log.Printf("Error saving video to %s: %v", filename, err)
			continue
		}

		log.Printf("Successfully downloaded and saved video %d to %s", i, filename)
	}

	log.Println("All downloads completed")
}

type Resource interface {
	Download() ([]byte, error)
}

type Finder interface {
	Find() ([]Resource, error)
}

// ファイルに保存する関数
func saveToFile(data []byte, filename string) error {
	// ディレクトリを作成
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
