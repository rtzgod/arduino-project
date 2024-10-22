package httpclient

import (
	"net/http"
	"log"
    "strings"
)

type Client struct {
	Url string
}

func New(url string) *Client {
	return &Client{
		Url: url,
	}
}

func (c *Client) Post (data string) {
	resp, err := http.Post(c.Url, "text/plain", strings.NewReader(data))
	if err != nil {
		log.Println("Ошибка отправки:", err)
		return
	}
	resp.Body.Close()
}

