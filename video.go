// TODO: 値 or ポインタ？

package main

type YTvideo struct {
	url string
}

func (*YTvideo) Download() ([]byte, error) {
	// Simulate downloading a URL
	return []byte("data"), nil
}

func NewYTvideo(url string) *YTvideo {
	return &YTvideo{url: url}
}
