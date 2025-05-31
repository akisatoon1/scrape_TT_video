// TODO: 並行処理を実装

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	// コマンドライン引数のチェック
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run main.go <max>")
	}
	max, err := strconv.Atoi(os.Args[1])
	if err != nil || max <= 0 {
		log.Fatal("Invalid max value. It should be a positive integer.")
	}

	apiKey, err := getAPIKey()
	if err != nil {
		log.Fatalf("Error loading API key: %v", err)
	}
	var finder Finder = NewYTvideoFinder(apiKey)

	resources, err := finder.Find(max)
	if err != nil {
		log.Fatal(err)
	}

	// 逐次処理
	log.Println("動画の逐次ダウンロードを開始します")

	for i, rsrc := range resources {
		id := rsrc.ID()
		log.Printf("%v) 動画 %v のダウンロードを開始します\n", i, id)

		// ダウンロード実行
		filename := id + ".mp4"
		err := rsrc.Download(filename)
		if err != nil {
			log.Printf("Error downloading video %v: %v\n", id, err)
			continue
		}

		log.Printf("Successfully downloaded and saved video %v to %s\n", id, filename)
	}

	log.Println("All downloads completed")
}

type Resource interface {
	ID() string
	Download(filename string) error
}

type Finder interface {
	Find(max int) ([]Resource, error)
}

// TODO: env読み込みを分ける
func getAPIKey() (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", fmt.Errorf("failed to load .env file: %v", err)
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API_KEY is not set in the environment")
	}
	return apiKey, nil
}
